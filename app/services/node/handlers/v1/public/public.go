// Package public maintains the group of handlers for public access.
package public

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/wtran29/go-blockchain/business/web/v1"
	"github.com/wtran29/go-blockchain/foundation/blockchain/database"
	"github.com/wtran29/go-blockchain/foundation/blockchain/state"
	"github.com/wtran29/go-blockchain/foundation/nameservice"
	"github.com/wtran29/go-blockchain/foundation/web"
	"go.uber.org/zap"
)

// Handlers manages the set of bar ledger endpoints.
type Handlers struct {
	Log   *zap.SugaredLogger
	State *state.State
	NS    *nameservice.NameService
	// WS    websocket.Upgrader
	// Evts  *events.Events
}

// // Events handles a web socket to provide events to a client.
// func (h Handlers) Events(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	v, err := web.GetValues(ctx)
// 	if err != nil {
// 		return web.NewShutdownError("web value missing from context")
// 	}

// 	// Need this to handle CORS on the websocket.
// 	h.WS.CheckOrigin = func(r *http.Request) bool { return true }

// 	// This upgrades the HTTP connection to a websocket connection.
// 	c, err := h.WS.Upgrade(w, r, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer c.Close()

// 	// This provides a channel for receiving events from the blockchain.
// 	ch := h.Evts.Acquire(v.TraceID)
// 	defer h.Evts.Release(v.TraceID)

// 	// Starting a ticker to send a ping message over the websocket.
// 	ticker := time.NewTicker(time.Second)

// 	// Block waiting for events from the blockchain or ticker.
// 	for {
// 		select {
// 		case msg, wd := <-ch:

// 			// If the channel is closed, release the websocket.
// 			if !wd {
// 				return nil
// 			}

// 			if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
// 				return err
// 			}

// 		case <-ticker.C:
// 			if err := c.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
// 				return nil
// 			}
// 		}
// 	}
// }

// SubmitWalletTransaction adds new transactions to the mempool.
func (h Handlers) SubmitWalletTransaction(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	// Decode the JSON in the post call into a Signed transaction.
	var signedTx database.SignedTx
	if err := web.Decode(r, &signedTx); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	h.Log.Infow("add tran", "traceid", v.TraceID, "sig:nonce", signedTx, "from", signedTx.FromID, "to", signedTx.ToID, "value", signedTx.Value, "tip", signedTx.Tip)

	// Ask the state package to add this transaction to the mempool. Only the
	// checks are the transaction signature and the recipient account format.
	// It's up to the wallet to make sure the account has a proper balance and
	// nonce. Fees will be taken if this transaction is mined into a block.
	if err := h.State.UpsertWalletTransaction(signedTx); err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "transactions added to mempool",
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

// Mempool returns the set of uncommitted transactions.
func (h Handlers) Mempool(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	acct := web.Param(r, "account")

	mempool := h.State.Mempool()

	trans := []tx{}
	for _, tran := range mempool {
		if acct != "" && ((acct != string(tran.FromID)) && (acct != string(tran.ToID))) {
			continue
		}

		trans = append(trans, tx{
			FromAccount: tran.FromID,
			FromName:    h.NS.Lookup(tran.FromID),
			To:          tran.ToID,
			ToName:      h.NS.Lookup(tran.ToID),
			ChainID:     tran.ChainID,
			Nonce:       tran.Nonce,
			Value:       tran.Value,
			Tip:         tran.Tip,
			Data:        tran.Data,
			TimeStamp:   tran.TimeStamp,
			GasPrice:    tran.GasPrice,
			GasUnits:    tran.GasUnits,
			Sig:         tran.SignatureString(),
		})
	}

	return web.Respond(ctx, w, trans, http.StatusOK)
}
