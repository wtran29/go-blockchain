// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wtran29/go-blockchain/app/services/node/handlers/v1/private"
	"github.com/wtran29/go-blockchain/app/services/node/handlers/v1/public"
	"github.com/wtran29/go-blockchain/foundation/blockchain/state"
	"github.com/wtran29/go-blockchain/foundation/events"
	"github.com/wtran29/go-blockchain/foundation/nameservice"
	"github.com/wtran29/go-blockchain/foundation/web"

	"go.uber.org/zap"
)

const version = "v1"

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *zap.SugaredLogger
	State *state.State
	NS    *nameservice.NameService
	Evts  *events.Events
}

// PublicRoutes binds all the version 1 public routes.
func PublicRoutes(app *web.App, cfg Config) {
	pbl := public.Handlers{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
		WS:    websocket.Upgrader{},
		Evts:  cfg.Evts,
	}

	app.Handle(http.MethodPost, version, "/start/mining", pbl.StartMining)
	app.Handle(http.MethodPost, version, "/tx/submit", pbl.SubmitWalletTransaction)
	app.Handle(http.MethodGet, version, "/tx/uncommitted/list", pbl.Mempool)
	app.Handle(http.MethodGet, version, "/tx/uncommitted/list/:account", pbl.Mempool)
	app.Handle(http.MethodGet, version, "/accounts/list", pbl.Accounts)
	app.Handle(http.MethodGet, version, "/accounts/list/:account", pbl.Accounts)
	app.Handle(http.MethodGet, version, "/blocks/list", pbl.BlocksByAccount)
	app.Handle(http.MethodGet, version, "/blocks/list/:account", pbl.BlocksByAccount)
}

// PrivateRoutes binds all the version 1 private routes.
func PrivateRoutes(app *web.App, cfg Config) {
	prv := private.Handlers{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
	}

	app.Handle(http.MethodPost, version, "/node/peers", prv.SubmitPeer)
	app.Handle(http.MethodGet, version, "/node/status", prv.Status)
	app.Handle(http.MethodGet, version, "/node/block/list/:from/:to", prv.BlocksByNumber)
	app.Handle(http.MethodPost, version, "/node/block/propose", prv.ProposeBlock)
	app.Handle(http.MethodPost, version, "/node/tx/submit", prv.SubmitNodeTransaction)
	app.Handle(http.MethodGet, version, "/node/tx/list", prv.Mempool)
}
