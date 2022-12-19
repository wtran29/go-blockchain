// Package database handles all the lower level support for maintaining the
// blockchain in storage and maintaining an in-memory databse of account information.
package database

import (
	"sync"

	"github.com/wtran29/go-blockchain/foundation/blockchain/genesis"
)

// Storage interface represents the behavior required to be implemented by any
// package providing support for reading and writing the blockchain.
type Storage interface {
	Write(blockData BlockData) error
	GetBlock(num uint64) (BlockData, error)
	ForEach() Iterator
	Close() error
	Reset() error
}

// Iterator interface represents the behavior required to be implemented by any
// package providing support to iterate over the blocks.
type Iterator interface {
	Next() (BlockData, error)
	Done() bool
}

// =============================================================================

// Database manages data related to accounts who have transacted on the blockchain.
type Database struct {
	mu          sync.RWMutex
	genesis     genesis.Genesis
	latestBlock Block
	accounts    map[AccountID]Account
	storage     Storage
}

// New constructs a new database and applies account genesis information and
// reads/writes the blockchain database on disk if a dbPath is provided.
func New(genesis genesis.Genesis, storage Storage, evHandler func(v string, args ...any)) (*Database, error) {
	db := Database{
		genesis:  genesis,
		accounts: make(map[AccountID]Account),
		storage:  storage,
	}

	// Update the database with account balance information from genesis.
	for accountStr, balance := range genesis.Balances {
		accountID, err := ToAccountID(accountStr)
		if err != nil {
			return nil, err
		}
		db.accounts[accountID] = newAccount(accountID, balance)
	}

	// // Read all the blocks from storage.
	// iter := db.ForEach()
	// for block, err := iter.Next(); !iter.Done(); block, err = iter.Next() {
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	// Validate the block values and cryptographic audit trail.
	// 	if err := block.ValidateBlock(db.latestBlock, db.HashState(), evHandler); err != nil {
	// 		return nil, err
	// 	}

	// 	// Update the database with the transaction information.
	// 	for _, tx := range block.MerkleTree.Values() {
	// 		db.ApplyTransaction(block, tx)
	// 	}
	// 	db.ApplyMiningReward(block)

	// 	// Update the current latest block.
	// 	db.latestBlock = block
	// }

	return &db, nil
}

// Close closes the open blocks database.
func (db *Database) Close() {
	db.storage.Close()
}
