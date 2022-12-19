// Package state is the core API for the blockchain and implements all the
// business rules and processing.
package state

import (
	"sync"

	"github.com/wtran29/go-blockchain/foundation/blockchain/database"
	"github.com/wtran29/go-blockchain/foundation/blockchain/genesis"
	"github.com/wtran29/go-blockchain/foundation/blockchain/mempool"
)

/*
	-- Blockchain
	On chain fork, only remove the block need to be removed and reset.

	-- Testing
	Fork Test
	Mining Test
*/

// =============================================================================

// EventHandler defines a function that is called when events
// occur in the processing of persisting blocks.
type EventHandler func(v string, args ...any)

// For logging purposes and foundation use, this function was built to be used to decouple items
// between production items and development

// =============================================================================

// Config represents the configuration required to start
// the blockchain node.
type Config struct {
	BeneficiaryID  database.AccountID
	Host           string
	Storage        database.Storage
	Genesis        genesis.Genesis
	SelectStrategy string
	// KnownPeers     *peer.PeerSet
	EvHandler EventHandler
	// Consensus      string
}

// State manages the blockchain database.
type State struct {
	mu sync.RWMutex
	// resyncWG    sync.WaitGroup
	// allowMining bool

	beneficiaryID database.AccountID
	host          string
	evHandler     EventHandler
	// consensus     string

	// knownPeers *peer.PeerSet
	storage database.Storage
	genesis genesis.Genesis
	mempool *mempool.Mempool
	db      *database.Database

	// Worker Worker
}

// New constructs a new blockchain for data management.
func New(cfg Config) (*State, error) {

	// Build a safe event handler function for use.
	ev := func(v string, args ...any) {
		if cfg.EvHandler != nil {
			cfg.EvHandler(v, args...)
		}
	}

	// Access the storage for the blockchain.
	db, err := database.New(cfg.Genesis, cfg.Storage, ev)
	if err != nil {
		return nil, err
	}

	// Construct a mempool with the specified sort strategy.
	mempool, err := mempool.NewWithStrategy(cfg.SelectStrategy)
	if err != nil {
		return nil, err
	}

	// Create the State to provide support for managing the blockchain.
	state := State{
		beneficiaryID: cfg.BeneficiaryID,
		host:          cfg.Host,
		storage:       cfg.Storage,
		evHandler:     ev,
		// consensus:     cfg.Consensus,
		// allowMining:   true,

		// knownPeers: cfg.KnownPeers,
		genesis: cfg.Genesis,
		mempool: mempool,
		db:      db,
	}

	// The Worker is not set here. The call to worker.Run will assign itself
	// and start everything up and running for the node.

	return &state, nil
}
