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
	"time"

	"github.com/coreos/bbolt"
	// "github.com/boltdb/bolt"
)

var gdocs = []byte("gdocs")

// Bolt bolt store struct
type Bolt struct {
	db *bolt.DB
}

// OpenBolt open Bolt store
func OpenBolt(dbPath string) (Store, error) {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 3600 * time.Second})
	// db, err := bolt.Open(dbPath, 0600, &bolt.Options{})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(gdocs)
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Bolt{db}, nil
}

// WALName returns the path to currently open database file.
func (s *Bolt) WALName() string {
	return s.db.Path()
}

// Set executes a function within the context of a read-write managed
// transaction. If no error is returned from the function then the transaction
// is committed. If an error is returned then the entire transaction is rolled back.
// Any error that is returned from the function or returned from the commit is returned
// from the Update() method.
func (s *Bolt) Set(k []byte, v []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(gdocs).Put(k, v)
	})
}

// Get executes a function within the context of a managed read-only transaction.
// Any error that is returned from the function is returned from the View() method.
func (s *Bolt) Get(k []byte) (b []byte, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(gdocs).Get(k)
		return nil
	})
	return
}

// Delete deletes a key. Exposing this so that user does not
// have to specify the Entry directly.
func (s *Bolt) Delete(k []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(gdocs).Delete(k)
	})
}

// Has returns true if the DB does contains the given key.
func (s *Bolt) Has(k []byte) (bool, error) {
	// return s.db.Exists(k)
	var b []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(gdocs).Get(k)
		return nil
	})

	// b == nil
	if err != nil || string(b) == "" {
		return false, err
	}

	return true, nil
}

// ForEach get all key and value
func (s *Bolt) ForEach(fn func(k, v []byte) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(gdocs)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
		}
		return nil
	})
}

// Close releases all database resources. All transactions
// must be closed before closing the database.
func (s *Bolt) Close() error {
	return s.db.Close()
}
