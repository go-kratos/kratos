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
	"github.com/syndtr/goleveldb/leveldb"
)

// Leveldb leveldb store
type Leveldb struct {
	db *leveldb.DB
}

// OpenLeveldb opens or creates a DB for the given store. The DB
// will be created if not exist, unless ErrorIfMissing is true.
// Also, if ErrorIfExist is true and the DB exist Open will
// returns os.ErrExist error.
func OpenLeveldb(dbPath string) (Store, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	return &Leveldb{db}, nil
}

// WALName is useless for this kv database
func (s *Leveldb) WALName() string {
	return "" // 对于此数据库，本函数没用~
}

// Set sets the provided value for a given key.
// If key is not present, it is created. If it is present,
// the existing value is overwritten with the one provided.
func (s *Leveldb) Set(k, v []byte) error {
	return s.db.Put(k, v, nil)
}

// Get gets the value for the given key. It returns
// ErrNotFound if the DB does not contains the key.
//
// The returned slice is its own copy, it is safe to modify
// the contents of the returned slice. It is safe to modify the contents
// of the argument after Get returns.
func (s *Leveldb) Get(k []byte) ([]byte, error) {
	return s.db.Get(k, nil)
}

// Delete deletes the value for the given key. Delete will not
// returns error if key doesn't exist. Write merge also applies
// for Delete, see Write.
//
// It is safe to modify the contents of the arguments after Delete
// returns but not before.
func (s *Leveldb) Delete(k []byte) error {
	return s.db.Delete(k, nil)
}

// Has returns true if the DB does contains the given key.
// It is safe to modify the contents of the argument after Has returns.
func (s *Leveldb) Has(k []byte) (bool, error) {
	return s.db.Has(k, nil)
}

// Len calculates approximate sizes of the given key ranges.
// The length of the returned sizes are equal with the length of
// the given ranges. The returned sizes measure store space usage,
// so if the user data compresses by a factor of ten, the returned
// sizes will be one-tenth the size of the corresponding user data size.
// The results may not include the sizes of recently written data.
func (s *Leveldb) Len() (leveldb.Sizes, error) {
	return s.db.SizeOf(nil)
}

// ForEach get all key and value
func (s *Leveldb) ForEach(fn func(k, v []byte) error) error {
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		val := iter.Value()
		if err := fn(key, val); err != nil {
			return err
		}
	}
	iter.Release()
	return iter.Error()
}

// Close closes the DB. This will also releases any outstanding snapshot,
// abort any in-flight compaction and discard open transaction.
func (s *Leveldb) Close() error {
	return s.db.Close()
}
