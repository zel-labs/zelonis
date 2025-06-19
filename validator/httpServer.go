/*
Copyright (C) 2025 Zelonis Contributors

This file is part of Zelonis.

Zelonis is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Zelonis is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Zelonis. If not, see <https://www.gnu.org/licenses/>.
*/

package validator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"zelonis/external"
	"zelonis/gossip"
	"zelonis/utils/maths"
	"zelonis/validator/domain"
	"zelonis/wallet"
)

type RpcServer struct {
	timeouts time.Duration
	mux      http.ServeMux // registered handlers go here

	mu       sync.Mutex
	server   *http.Server
	listener net.Listener // non-nil when server is running

	// HTTP RPC handler things.

	httpHandler atomic.Value // *rpcHandler

	// These are set by setListenAddr.
	endpoint string
	host     string
	port     int

	handlerNames  map[string]string
	databases     map[*DbTrackers]struct{}
	domain        *domain.Domain
	gossipManager *gossip.Manager
}

func NewHTTPServer(timeouts time.Duration, domain *domain.Domain, port int) *RpcServer {
	return &RpcServer{
		timeouts:     timeouts,
		handlerNames: make(map[string]string),
		port:         port,
		host:         DefaultHTTPHost,
		domain:       domain,
	}
}
func (s *RpcServer) Start(flowManager *gossip.Manager) {
	s.gossipManager = flowManager
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma-separated
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT",
	}))
	app.Get("/createWallet", s.createWallet)
	app.Get("/recoverWallet/:seed/keys/:keys", s.recoverWallet)
	app.Post("/sendTx/", s.sendTx)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Hello, Fiber!"))
	})
	s.webserver(app)
	err := app.Listen(":" + strconv.Itoa(s.port))
	if err != nil {
		panic(err)
	}
}

func (s *RpcServer) webserver(app *fiber.App) {
	//app.Get("/block/:hash", s.getBlock)
	app.Get("/blockById/:id", s.getBlockById)
	app.Get("/latestBlocks", s.getLatestBlocks)
	app.Get("/latestTx/", s.getLatestTx)
	app.Get("/blockByHash/:hash", s.getBlockByHash)
	app.Get("/tx/:hash", s.getTxByHash)
	app.Get("/account/:account", s.getAccountBalance)
	app.Get("/accountTx/:account", s.getAccountTx)
	app.Get("/accountTxList/:account", s.getAccountTxList)
	app.Get("/accountTxLimit/:account/:from/:limit", s.getAccountTxWithLimit)
	app.Get("/currentStatus/", s.getCurrentStatus)

}

func (s *RpcServer) getLatestTx(c *fiber.Ctx) error {
	blockHeight, _ := s.domain.StatsManager().GetHighestBlockHeight()
	txs := make([]*JsonTx, 0)
txLimit:
	for blockHeight >= 0 {

		block, _ := s.domain.BlockManager().GetBlockById(strconv.FormatUint(blockHeight, 10))
		for _, tx := range block.Transactions {
			jsonTx := s.externalTxHashToJsonTx(tx)

			jsonTx.BlockHeight = strconv.FormatUint(blockHeight, 10)

			jsonTx.Timstamp = block.Header.BlockTime

			txs = append(txs, jsonTx)
			if len(txs) == 10 {
				break txLimit
			}
		}
		if blockHeight == 0 {
			break
		}
		blockHeight--

	}
	JsonTxList := &jsonTxList{
		TxList: txs,
	}
	return c.JSON(JsonTxList)

}

type nodeInfo struct {
	LatestBlock        string          `json:"latestBlock"`
	LastestBlockHeight uint64          `json:"lastestBlockHeight"`
	TotalTransactions  uint64          `json:"totalTransactions"`
	TotalStaked        string          `json:"totalStaked"`
	TotalSupply        string          `json:"totalSupply"`
	TotalCirculating   string          `json:"totalCirculating"`
	EpochInfo          *external.Epoch `json:"epoch"`
}

func (s *RpcServer) getCurrentStatus(c *fiber.Ctx) error {
	lastBlockHeight, _ := s.domain.StatsManager().GetHighestBlockHeight()
	lastBlockHash, _ := s.domain.GetHighestBlockHash()
	totalTx, _ := s.domain.StatsManager().GetTotalTransactions()
	totalStaked, _ := s.domain.StatsManager().GetTotalStake()
	totalSupply, _ := s.domain.StatsManager().GetTotalSupply()
	circulating, _ := s.domain.StatsManager().GetCirculating()
	epoch := s.domain.StatsManager().GetEpoch()
	nodeinfo := &nodeInfo{
		LastestBlockHeight: lastBlockHeight,
		LatestBlock:        fmt.Sprintf("%x", lastBlockHash),
		TotalTransactions:  totalTx,
		TotalStaked:        string(totalStaked),
		TotalSupply:        string(totalSupply),
		TotalCirculating:   string(circulating),
		EpochInfo:          epoch,
	}

	c.JSON(nodeinfo)
	return nil
}

func (s *RpcServer) getBlockByHash(c *fiber.Ctx) error {
	hash := c.Params("hash")
	if hash == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	hashBytes, err := external.NewDomainHashFromString(hash)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	block, err := s.domain.BlockManager().GetBlockByHash(hashBytes.ByteSlice())
	if err != nil {
		return err
	}
	c.JSON(s.blockToJsonBlock(block))
	return nil
}

type jsonBlock struct {
	Header        *jsonHeader
	Transactions  []*JsonTx
	ValidatorInfo *jsonValidator
	Signature     string
	RecivedAt     time.Time
}
type jsonValidator struct {
	Addr    string
	PrevKey string
}
type jsonHeader struct {
	Blockheight uint64
	Blockhash   string
	Blocktime   int64
	ParentSlot  uint64
	ParentHash  string
	Version     int8
}

func (s *RpcServer) getBlockById(c *fiber.Ctx) error {
	blockId := c.Params("id")
	block, err := s.domain.BlockManager().GetBlockById(blockId)
	if err != nil {
		return err
	}
	c.JSON(s.blockToJsonBlock(block))
	return nil
}

type jsonLatestBlocks struct {
	Blocks []*jsonBlock `json:"blocks"`
}

func (s *RpcServer) getLatestBlocks(c *fiber.Ctx) error {
	currentBlockHeight, _ := s.domain.StatsManager().GetHighestBlockHeight()
	i := uint64(0)
	blocks := make([]*jsonBlock, 0)
	for i < 10 {

		ci := currentBlockHeight - i
		block, err := s.domain.BlockManager().GetBlockById(strconv.FormatUint(ci, 10))
		if err != nil {
			return fmt.Errorf("%s : Key val %v", err, ci)
		}
		blocks = append(blocks, s.blockToJsonBlock(block))
		i++
		if ci == 0 {
			break
		}
	}
	latestBlocks := &jsonLatestBlocks{
		Blocks: blocks,
	}

	c.JSON(latestBlocks)
	return nil
}

func (s *RpcServer) blockToJsonBlock(block *external.Block) *jsonBlock {
	jsonTxs := make([]*JsonTx, 0)
	for _, tx := range block.Transactions {
		jsonTxs = append(jsonTxs, s.externalTxHashToJsonTx(tx))
	}
	blockJson := &jsonBlock{
		Header: &jsonHeader{
			Blockheight: block.Header.BlockHeight,
			Blockhash:   fmt.Sprintf("%x", block.Header.BlockHash),
			Blocktime:   block.Header.BlockTime,
			ParentSlot:  block.Header.ParentSlot,
			ParentHash:  fmt.Sprintf("%x", block.Header.ParentHash),
			Version:     block.Header.Version,
		},
		Transactions: jsonTxs,
		ValidatorInfo: &jsonValidator{
			Addr:    fmt.Sprintf("%s", block.Validator.Pubkey),
			PrevKey: fmt.Sprintf("%x", block.Validator.PreviousBlockHash),
		},
		Signature: fmt.Sprintf("%x", block.Signature),
	}
	return blockJson
}

func (s *RpcServer) createWallet(c *fiber.Ctx) error {

	return c.JSON(wallet.CreateWallet())
}

func (s *RpcServer) sendTx(c *fiber.Ctx) error {
	seed := c.FormValue("seed")
	keys := c.FormValue("keys")
	receiverVal := c.FormValue("receiver")
	if len(receiverVal) < 43 || len(receiverVal) > 44 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	receiver := []byte(receiverVal)
	valStr := c.FormValue("val")
	val := []byte(valStr)
	rWallet := (wallet.RecoverWallet(keys, seed))

	if !maths.IsValidNumber(valStr) {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	walletAddr := []byte(rWallet.Address)
	hash, _ := s.domain.GetHighestBlockHash()

	tx := s.domain.TxManager().BuildTxFromType(walletAddr, receiver, val, hash, external.TxTransfer)
	//Build Transaction
	sig := s.domain.TxManager().SignTxAndVerify(tx, rWallet.PrivateKey)
	tx.Signature = sig

	s.gossipManager.BroadcastTransaction(tx)
	c.JSON(fmt.Sprintf("%x", tx.TxHash))
	return nil
}

func (s *RpcServer) recoverWallet(c *fiber.Ctx) error {
	seed := c.Params("seed")
	encryptKey := c.Params("keys")
	oSeed, _ := url.QueryUnescape(seed)
	return c.JSON(wallet.RecoverWallet(encryptKey, oSeed))
}
