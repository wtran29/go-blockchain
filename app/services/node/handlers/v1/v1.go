// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"net/http"

	"github.com/wtran29/go-blockchain/app/services/node/handlers/v1/private"
	"github.com/wtran29/go-blockchain/app/services/node/handlers/v1/public"
	"github.com/wtran29/go-blockchain/foundation/blockchain/state"
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
	// Evts  *events.Events
}

// PublicRoutes binds all the version 1 public routes.
func PublicRoutes(app *web.App, cfg Config) {
	pbl := public.Handlers{
		Log:   cfg.Log,
		State: cfg.State,
		NS:    cfg.NS,
		// WS:    websocket.Upgrader{},
		// Evts:  cfg.Evts,
	}

	// app.Handle(http.MethodPost, version, "/start/mining", pbl.StartMining)
	app.Handle(http.MethodPost, version, "/tx/submit", pbl.SubmitWalletTransaction)
	app.Handle(http.MethodGet, version, "/tx/uncommitted/list/:account", pbl.Mempool)
}

// PrivateRoutes binds all the version 1 private routes.
func PrivateRoutes(app *web.App, cfg Config) {
	prv := private.Handlers{
		Log: cfg.Log,
	}

	app.Handle(http.MethodGet, version, "/node/sample", prv.Sample)
}
