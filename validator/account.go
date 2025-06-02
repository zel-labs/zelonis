package validator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type jsonAccountInfo struct {
	Balance string `json:"balance"`
	Stake   string `json:"stake"`
	Reward  string `json:"reward"`
}

func (s *RpcServer) getAccountBalance(c *fiber.Ctx) error {
	account := c.Params("account")
	accountInfo, status := s.domain.GetAccountBalance([]byte(account))
	if !status {
		return fmt.Errorf("Account not found")
	}

	jsonAccount := &jsonAccountInfo{
		Balance: s.valToJsonVal(accountInfo.Balance),
		Stake:   s.valToJsonVal(accountInfo.Stake),
		Reward:  s.valToJsonVal(accountInfo.Reward),
	}
	c.JSON(jsonAccount)
	return nil
}

func (s *RpcServer) valToJsonVal(val []byte) string {
	valStr := fmt.Sprintf("%s", val)
	//bigVal, _ := big.NewFloat(0).SetString(valStr)
	valFloat, _ := strconv.ParseFloat(valStr, 64)
	return fmt.Sprintf("%.7f", valFloat)
}

func (s *RpcServer) getAccountTx(c *fiber.Ctx) error {
	account := c.Params("account")

	list, err := s.domain.AccountManager().GetTxHistoryList([]byte(account), 0, 10)
	if err != nil {
		return err
	}
	c.JSON(list)
	return nil
}

func (s *RpcServer) getAccountTxWithLimit(c *fiber.Ctx) error {
	account := c.Params("account")
	from, _ := strconv.Atoi(c.Params("from"))
	limit, _ := strconv.Atoi(c.Params("limit"))

	list, err := s.domain.AccountManager().GetTxHistoryList([]byte(account), from, limit)
	if err != nil {
		return err
	}
	c.JSON(list)
	return nil
}
