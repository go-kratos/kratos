# proto

[![Build Status](https://travis-ci.org/emicklei/proto.png)](https://travis-ci.org/emicklei/proto)
[![Go Report Card](https://goreportcard.com/badge/github.com/emicklei/proto)](https://goreportcard.com/report/github.com/emicklei/proto)
[![GoDoc](https://godoc.org/github.com/emicklei/proto?status.svg)](https://godoc.org/github.com/emicklei/proto)

Package in Go for parsing Google Protocol Buffers [.proto files version 2 + 3] (https://developers.google.com/protocol-buffers/docs/reference/proto3-spec)

### install

    go get -u -v github.com/emicklei/proto

### usage

	package main

	import (
		"fmt"
		"os"

		"github.com/emicklei/proto"
	)

	func main() {
		reader, _ := os.Open("test.proto")
		defer reader.Close()

		parser := proto.NewParser(reader)
		definition, _ := parser.Parse()

		proto.Walk(definition,
			proto.WithService(handleService),
			proto.WithMessage(handleMessage))
	}

	func handleService(s *proto.Service) {
		fmt.Println(s.Name)
	}

	func handleMessage(m *proto.Message) {
		fmt.Println(m.Name)
	}


### contributions

See (https://github.com/emicklei/proto-contrib) for other contributions on top of this package such as protofmt, proto2xsd and proto2gql.


© 2017, [ernestmicklei.com](http://ernestmicklei.com).  MIT License. Contributions welcome.