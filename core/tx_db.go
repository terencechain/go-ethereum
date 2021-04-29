// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"io/ioutil"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	// maxMemoryItems is the number of transactions that should be
	// held in memory maximally
	maxMemoryItems = 4096
	// minMemoryItems is the number of transactions that should be
	// held in memory after a rebalance
	minMemoryItems = 3072
)

type txDB struct {
	items map[common.Hash]*types.Transaction // Hash map storing the transaction data
	db    *leveldb.DB
	mu    *sync.RWMutex
}

func newTxDB(db *leveldb.DB) *txDB {
	if db == nil {
		var err error
		g, _ := ioutil.TempDir("", "")
		db, err = leveldb.OpenFile(g, nil)
		if err != nil {
			panic(err)
		}
	}
	return &txDB{
		items: make(map[common.Hash]*types.Transaction),
		db:    db,
		mu:    new(sync.RWMutex),
	}
}

func (t *txDB) Add(tx *types.Transaction) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.items) > maxMemoryItems {
		t.rebalance()
	}
	t.items[tx.Hash()] = tx
}

func (t *txDB) Get(hash common.Hash) (*types.Transaction, error) {
	t.mu.RLock()
	// retrieve from cache
	item, ok := t.items[hash]
	t.mu.RUnlock()
	if ok {
		return item, nil
	}

	// retrieve from db
	var tx *types.Transaction
	val, err := t.db.Get(hash[:], nil)
	if err != nil {
		return nil, err
	}
	if tx.UnmarshalBinary(val) != nil {
		return nil, err
	}
	return tx, nil
}

func (t *txDB) Remove(hash common.Hash) (*types.Transaction, error) {
	t.mu.Lock()
	// retrieve from cache
	item, ok := t.items[hash]
	if ok {
		delete(t.items, hash)
		t.mu.Unlock()
		return item, nil
	}
	t.mu.Unlock()
	// retrieve from db
	val, err := t.db.Get(hash[:], nil)
	if err != nil {
		return nil, err
	}
	// found, delete from db
	defer func() {
		t.db.Delete(hash[:], nil)
	}()
	var tx *types.Transaction
	if tx.UnmarshalBinary(val) != nil {
		return nil, err
	}
	return tx, nil
}

// rebalance rebalances the transactions between the DB and memory.
// assumes a write lock is held.
func (t *txDB) rebalance() {
	txs := make(types.Transactions, len(t.items))
	for _, item := range t.items {
		txs = append(txs, item)
	}
	sort.Sort(types.TxByPriceAndTime(txs))
	for i := minMemoryItems; i < len(t.items); i++ {
		tx := txs[i]
		marshalled, _ := tx.MarshalBinary()
		if err := t.db.Put(tx.Hash().Bytes(), marshalled, nil); err != nil {
			panic(err)
		}
		delete(t.items, tx.Hash())
	}
}
