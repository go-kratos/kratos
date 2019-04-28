murmur
======
[![CircleCI Status](https://circleci.com/gh/go-ego/murmur.svg?style=shield)](https://circleci.com/gh/go-ego/murmur)
[![codecov](https://codecov.io/gh/go-ego/murmur/branch/master/graph/badge.svg)](https://codecov.io/gh/go-ego/murmur)
[![Build Status](https://travis-ci.org/go-ego/murmur.svg)](https://travis-ci.org/go-ego/murmur)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-ego/murmur)](https://goreportcard.com/report/github.com/go-ego/murmur)
[![GoDoc](https://godoc.org/github.com/go-ego/murmur?status.svg)](https://godoc.org/github.com/go-ego/murmur)
[![Release](https://github-release-version.herokuapp.com/github/go-ego/murmur/release.svg?style=flat)](https://github.com/go-ego/murmur/releases/latest)
[![Join the chat at https://gitter.im/go-ego/ego](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-ego/ego?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Go Murmur3 hash implementation

Based on [MurmurHash](http://en.wikipedia.org/wiki/MurmurHash), [murmur](https://github.com/huichen/murmur).

## Installing
```Go
go get -u github.com/go-ego/murmur
```

## Use

```Go
package main

import (
	"log"

	"github.com/go-ego/murmur"
)

func main() {
	var str = "github.com"
	
	hash32 := murmur.Murmur3([]byte(str))
	log.Println("hash32...", hash32)

	sum32 := murmur.Sum32(str)
	log.Println("hash32...", sum32)
}
```
