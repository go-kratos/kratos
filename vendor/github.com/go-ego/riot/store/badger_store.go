// Copyright 2016 ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package store

import (
	"log"

	"github.com/dgraph-io/badger"
)

// Badger badger.KV db store
type Badger struct {
	db *badger.DB
}

// OpenBadger open Badger store
func OpenBadger(dbPath string) (Store, error) {
	// err := os.MkdirAll(dbPath, 0777)
	// if err != nil {
	// 	log.Fatal("os.MkdirAll: ", err)
	// 	os.Exit(1)
	// }
	// os.MkdirAll(path.Dir(dbPath), os.ModePerm)

	opt := badger.DefaultOptions
	opt.Dir = dbPath
	opt.ValueDir = dbPath
	opt.SyncWrites = true
	kv, err := badger.Open(opt)
	if err != nil {
		log.Fatal("badger NewKV: ", err)
	}

	return &Badger{kv}, err
}

// WALName is useless for this kv database
func (s *Badger) WALName() string {
	return "" // 对于此数据库，本函数没用~
}

// Set sets the provided value for a given key.
// If key is not present, it is created. If it is present,
// the existing value is overwritten with the one provided.
func (s *Badger) Set(k, v []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		// return txn.Set(k, v, 0x00)
		return txn.Set(k, v)
	})

	return err
}

// Get looks for key and returns a value.
// If key is not found, value is nil.
func (s *Badger) Get(k []byte) ([]byte, error) {
	var ival []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}

		ival, err = item.Value()
		return err
	})

	return ival, err
}

// Delete deletes a key. Exposing this so that user does not
// have to specify the Entry directly. For example, BitDelete
// seems internal to badger.
func (s *Badger) Delete(k []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(k)
	})

	return err
}

// Has returns true if the DB does contains the given key.
func (s *Badger) Has(k []byte) (bool, error) {
	// return s.db.Exists(k)
	val, err := s.Get(k)
	if string(val) == "" && err != nil {
		return false, err
	}

	return true, err
}

// Len returns the size of lsm and value log files in bytes.
// It can be used to decide how often to call RunValueLogGC.
func (s *Badger) Len() (int64, int64) {
	return s.db.Size()
}

// ForEach get all key and value
func (s *Badger) ForEach(fn func(k, v []byte) error) error {
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 1000
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			val, err := item.Value()
			if err != nil {
				return err
			}

			if err := fn(key, val); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// Close closes a KV. It's crucial to call it to ensure
// all the pending updates make their way to disk.
func (s *Badger) Close() error {
	return s.db.Close()
}
