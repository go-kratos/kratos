# cedar 
[![Build Status](https://travis-ci.org/go-ego/cedar.svg)](https://travis-ci.org/go-ego/cedar)
[![CircleCI Status](https://circleci.com/gh/go-ego/cedar.svg?style=shield)](https://circleci.com/gh/go-ego/cedar)
[![codecov](https://codecov.io/gh/go-ego/cedar/branch/master/graph/badge.svg)](https://codecov.io/gh/go-ego/cedar)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-ego/cedar)](https://goreportcard.com/report/github.com/go-ego/cedar)
[![GoDoc](https://godoc.org/github.com/go-ego/cedar?status.svg)](https://godoc.org/github.com/go-ego/cedar)
[![Release](https://github-release-version.herokuapp.com/github/go-ego/cedar/release.svg?style=flat)](https://github.com/go-ego/cedar/releases/latest)
[![Join the chat at https://gitter.im/go-ego/ego](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-ego/ego?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Package `cedar` implementes double-array trie, base on [cedar-go](https://github.com/adamzy/cedar-go).

It is a [Golang](https://golang.org/) port of [cedar](http://www.tkl.iis.u-tokyo.ac.jp/~ynaga/cedar) which is written in C++ by Naoki Yoshinaga. `cedar-go` currently implements the `reduced` verion of cedar. 
This package is not thread safe if there is one goroutine doing insertions or deletions. 

## Install
```
go get github.com/go-ego/cedar
```

## Usage
```go
package main

import (
	"fmt"

	"github.com/go-ego/cedar"
)

func main() {
	// create a new cedar trie.
	trie := cedar.New()

	// a helper function to print the id-key-value triple given trie node id
	printIdKeyValue := func(id int) {
		// the key of node `id`.
		key, _ := trie.Key(id)
		// the value of node `id`.
		value, _ := trie.Value(id)
		fmt.Printf("%d\t%s:%v\n", id, key, value)
	}

	// Insert key-value pairs.
    // The order of insertion is not important.
	trie.Insert([]byte("How many"), 0)
	trie.Insert([]byte("How many loved"), 1)
	trie.Insert([]byte("How many loved your moments"), 2)
	trie.Insert([]byte("How many loved your moments of glad grace"), 3)
	trie.Insert([]byte("姑苏"), 4)
	trie.Insert([]byte("姑苏城外"), 5)
	trie.Insert([]byte("姑苏城外寒山寺"), 6)

	// Get the associated value of a key directly.
	value, _ := trie.Get([]byte("How many loved your moments of glad grace"))
	fmt.Println(value)

	// Or, jump to the node first,
	id, _ := trie.Jump([]byte("How many loved your moments"), 0)
	// then get the key and the value
	printIdKeyValue(id)

	fmt.Println("\nPrefixMatch\nid\tkey:value")
	for _, id := range trie.PrefixMatch([]byte("How many loved your moments of glad grace"), 0) {
		printIdKeyValue(id)
	}

	fmt.Println("\nPrefixPredict\nid\tkey:value")
	for _, id := range trie.PrefixPredict([]byte("姑苏"), 0) {
		printIdKeyValue(id)
	}
}
```
will produce
```
3
281	How many loved your moments:2

PrefixMatch
id	key:value
262	How many:0
268	How many loved:1
281	How many loved your moments:2
296	How many loved your moments of glad grace:3

PrefixPredict
id	key:value
303	姑苏:4
309	姑苏城外:5
318	姑苏城外寒山寺:6
```
## License

Under the GPL-3.0 License.