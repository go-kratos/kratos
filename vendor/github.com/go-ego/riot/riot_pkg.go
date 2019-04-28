// Copyright 2013 Hui Chen
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

/*

Package riot full text search engine
*/
package riot

import (
	// _ "github.com/cznic/kv"
	_ "github.com/coreos/bbolt"
	// _ "github.com/boltdb/bolt"
	_ "github.com/dgraph-io/badger"
	_ "github.com/go-ego/gse"
	_ "github.com/go-ego/murmur"
	_ "github.com/syndtr/goleveldb/leveldb"
)
