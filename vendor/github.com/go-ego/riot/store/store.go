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
	"fmt"
	"os"
)

const (
	// DefaultStore default store engine
	DefaultStore = "ldb"
	// DefaultStore = "bad"
	// DefaultStore = "bolt"
)

var supportedStore = map[string]func(path string) (Store, error){
	"ldb":  OpenLeveldb,
	"bg":   OpenBadger, // bad to bg
	"bolt": OpenBolt,
	// "kv":   OpenKV,
	// "ledisdb": Open,
}

// RegisterStore register store engine
func RegisterStore(name string, fn func(path string) (Store, error)) {
	supportedStore[name] = fn
}

// Store is store interface
type Store interface {
	// type KVBatch interface {
	Set(k, v []byte) error
	Get(k []byte) ([]byte, error)
	Delete(k []byte) error
	Has(k []byte) (bool, error)
	ForEach(fn func(k, v []byte) error) error
	Close() error
	WALName() string
}

// OpenStore open store engine
func OpenStore(path string, args ...string) (Store, error) {
	storeName := DefaultStore

	if len(args) > 0 && args[0] != "" {
		storeName = args[0]
	} else {
		storeEnv := os.Getenv("Riot_Store_Engine")
		if storeEnv != "" {
			storeName = storeEnv
		}
	}

	if fn, has := supportedStore[storeName]; has {
		return fn(path)
	}

	return nil, fmt.Errorf("unsupported store engine: %v", storeName)
}
