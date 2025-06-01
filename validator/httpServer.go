package validator

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"zelonis/validator/domain"
	"zelonis/wallet"
)

type httpServer struct {
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

	handlerNames map[string]string
	databases    map[*DbTrackers]struct{}
	domain       *domain.Domain
}

func (vn *Validator) newHTTPServer(timeouts time.Duration) *httpServer {
	return &httpServer{
		timeouts:     timeouts,
		handlerNames: make(map[string]string),
		port:         DefaultHTTPPort,
		host:         DefaultHTTPHost,
		domain:       vn.domain,
	}
}
func (s *httpServer) start() {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // comma-separated
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT",
	}))
	app.Get("/createWallet", s.createWallet)
	app.Get("/recoverWallet/:seed/keys/:keys", s.recoverWallet)
	app.Get("/sendTx/", s.sendTx)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Hello, Fiber!"))
	})
	s.webserver(app)
	err := app.Listen(":" + strconv.Itoa(s.port))
	if err != nil {
		panic(err)
	}
}

func (s *httpServer) webserver(app *fiber.App) {
	//app.Get("/block/:hash", s.getBlock)
	app.Get("/blockById/:id", s.getBlockById)

}

func (s *httpServer) getBlockById(c *fiber.Ctx) error {
	blockId := c.Params("id")
	block, err := s.domain.BlockManager().GetBlockById(blockId)
	if err != nil {
		return err
	}
	c.JSON(block)
	return nil
}

func (s *httpServer) createWallet(c *fiber.Ctx) error {

	return c.JSON(wallet.CreateWallet())
}

func (s *httpServer) sendTx(c *fiber.Ctx) error {
	/*seed := c.FormValue("seed")
	keys := c.FormValue("keys")
	reciver := c.FormValue("reciver")
	val := c.FormValue("val")
	*/
	return nil
}

func (s *httpServer) recoverWallet(c *fiber.Ctx) error {
	seed := c.Params("seed")
	encryptKey := c.Params("keys")
	oSeed, _ := url.QueryUnescape(seed)
	return c.JSON(wallet.RecoverWallet(encryptKey, oSeed))
}
