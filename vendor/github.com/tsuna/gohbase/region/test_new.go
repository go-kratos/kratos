// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// +build testing

package region

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
)

type testClient struct {
	addr    string
	numNSRE int32
}

var nsreRegion = &pb.Result{Cell: []*pb.Cell{
	&pb.Cell{
		Row:       []byte("nsre,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("regioninfo"),
		Value: []byte("PBUF\b\xc4\xcd\xe9\x99\xe0)\x12\x0f\n\adefault\x12\x04nsre" +
			"\x1a\x00\"\x00(\x000\x008\x00"),
	},
	&pb.Cell{
		Row:       []byte("nsre,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("seqnumDuringOpen"),
		Value:     []byte("\x00\x00\x00\x00\x00\x00\x00\x02"),
	},
	&pb.Cell{
		Row:       []byte("nsre,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("server"),
		Value:     []byte("regionserver:1"),
	},
	&pb.Cell{
		Row:       []byte("nsre,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("serverstartcode"),
		Value:     []byte("\x00\x00\x01N\x02\x92R\xb1"),
	},
}}

// makeRegionResult returns a region that spans the whole table
// and uses name of the table as the hostname of the regionserver
func makeRegionResult(key []byte) *pb.ScanResponse {
	s := bytes.SplitN(key, []byte(","), 2)
	fqtable := s[0]

	row := append(fqtable, []byte(",,1434573235908.56f833d5569a27c7a43fbf547b4924a4.")...)
	t := bytes.SplitN(fqtable, []byte{':'}, 2)
	var namespace, table []byte
	if len(t) == 2 {
		namespace = t[0]
		table = t[1]
	} else {
		namespace = []byte("default")
		table = fqtable
	}
	regionInfo := &pb.RegionInfo{
		RegionId: proto.Uint64(1434573235908),
		TableName: &pb.TableName{
			Namespace: namespace,
			Qualifier: table,
		},
		Offline: proto.Bool(false),
	}
	regionInfoValue, err := proto.Marshal(regionInfo)
	if err != nil {
		panic(err)
	}
	regionInfoValue = append([]byte("PBUF"), regionInfoValue...)

	return &pb.ScanResponse{Results: []*pb.Result{
		&pb.Result{Cell: []*pb.Cell{
			&pb.Cell{
				Row:       row,
				Family:    []byte("info"),
				Qualifier: []byte("regioninfo"),
				Value:     regionInfoValue,
			},
			&pb.Cell{
				Row:       row,
				Family:    []byte("info"),
				Qualifier: []byte("seqnumDuringOpen"),
				Value:     []byte("\x00\x00\x00\x00\x00\x00\x00\x02"),
			},
			&pb.Cell{
				Row:       row,
				Family:    []byte("info"),
				Qualifier: []byte("server"),
				Value:     fqtable,
			},
			&pb.Cell{
				Row:       row,
				Family:    []byte("info"),
				Qualifier: []byte("serverstartcode"),
				Value:     []byte("\x00\x00\x01N\x02\x92R\xb1"),
			},
		}}}}
}

var metaRow = &pb.Result{Cell: []*pb.Cell{
	&pb.Cell{
		Row:       []byte("test,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("regioninfo"),
		Value: []byte("PBUF\b\xc4\xcd\xe9\x99\xe0)\x12\x0f\n\adefault\x12\x04test" +
			"\x1a\x00\"\x00(\x000\x008\x00"),
	},
	&pb.Cell{
		Row:       []byte("test,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("seqnumDuringOpen"),
		Value:     []byte("\x00\x00\x00\x00\x00\x00\x00\x02"),
	},
	&pb.Cell{
		Row:       []byte("test,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("server"),
		Value:     []byte("regionserver:2"),
	},
	&pb.Cell{
		Row:       []byte("test,,1434573235908.56f833d5569a27c7a43fbf547b4924a4."),
		Family:    []byte("info"),
		Qualifier: []byte("serverstartcode"),
		Value:     []byte("\x00\x00\x01N\x02\x92R\xb1"),
	},
}}

var test1SplitA = &pb.Result{Cell: []*pb.Cell{
	&pb.Cell{
		Row:       []byte("test1,,1480547738107.825c5c7e480c76b73d6d2bad5d3f7bb8."),
		Family:    []byte("info"),
		Qualifier: []byte("regioninfo"),
		Value: []byte("PBUF\b\xfbÖ\xbc\x8b+\x12\x10\n\adefault\x12\x05" +
			"test1\x1a\x00\"\x03baz(\x000\x008\x00"),
	},
	&pb.Cell{
		Row:       []byte("test1,,1480547738107.825c5c7e480c76b73d6d2bad5d3f7bb8."),
		Family:    []byte("info"),
		Qualifier: []byte("seqnumDuringOpen"),
		Value:     []byte("\x00\x00\x00\x00\x00\x00\x00\v"),
	},
	&pb.Cell{
		Row:       []byte("test1,,1480547738107.825c5c7e480c76b73d6d2bad5d3f7bb8."),
		Family:    []byte("info"),
		Qualifier: []byte("server"),
		Value:     []byte("regionserver:1"),
	},
	&pb.Cell{
		Row:       []byte("test1,,1480547738107.825c5c7e480c76b73d6d2bad5d3f7bb8."),
		Family:    []byte("info"),
		Qualifier: []byte("serverstartcode"),
		Value:     []byte("\x00\x00\x01X\xb6\x83^3"),
	},
}}

var test1SplitB = &pb.Result{Cell: []*pb.Cell{
	&pb.Cell{
		Row:       []byte("test1,baz,1480547738107.3f2483f5618e1b791f58f83a8ebba6a9."),
		Family:    []byte("info"),
		Qualifier: []byte("regioninfo"),
		Value: []byte("PBUF\b\xfbÖ\xbc\x8b+\x12\x10\n\adefault\x12\x05" +
			"test1\x1a\x03baz\"\x00(\x000\x008\x00"),
	},
	&pb.Cell{
		Row:       []byte("test1,baz,1480547738107.3f2483f5618e1b791f58f83a8ebba6a9."),
		Family:    []byte("info"),
		Qualifier: []byte("seqnumDuringOpen"),
		Value:     []byte("\x00\x00\x00\x00\x00\x00\x00\f"),
	},
	&pb.Cell{
		Row:       []byte("test1,baz,1480547738107.3f2483f5618e1b791f58f83a8ebba6a9."),
		Family:    []byte("info"),
		Qualifier: []byte("server"),
		Value:     []byte("regionserver:3"),
	},
	&pb.Cell{
		Row:       []byte("test1,baz,1480547738107.3f2483f5618e1b791f58f83a8ebba6a9."),
		Family:    []byte("info"),
		Qualifier: []byte("serverstartcode"),
		Value:     []byte("\x00\x00\x01X\xb6\x83^3"),
	},
}}

var m sync.RWMutex
var clients map[string]uint32

func init() {
	clients = make(map[string]uint32)
}

// NewClient creates a new test region client.
func NewClient(ctx context.Context, addr string, ctype ClientType,
	queueSize int, flushInterval time.Duration, effectiveUser string,
	readTimeout time.Duration) (hrpc.RegionClient, error) {
	m.Lock()
	clients[addr]++
	m.Unlock()
	return &testClient{addr: addr}, nil
}

func (c *testClient) Addr() string {
	return c.addr
}

func (c *testClient) String() string {
	return fmt.Sprintf("RegionClient{Addr: %s}", c.addr)
}

func (c *testClient) QueueRPC(call hrpc.Call) {
	// ignore timed out rpcs to mock the region client
	select {
	case <-call.Context().Done():
		return
	default:
	}
	if !bytes.Equal(call.Table(), []byte("hbase:meta")) {
		_, ok := call.(*hrpc.Get)
		if !ok || !bytes.HasSuffix(call.Key(), bytes.Repeat([]byte{0}, 17)) {
			// not a get and not a region probe
			// just return as the mock call should just populate the ResultChan in test
			return
		}
		// region probe, fail for the nsre region 3 times to force retry
		if bytes.Equal(call.Table(), []byte("nsre")) {
			i := atomic.AddInt32(&c.numNSRE, 1)
			if i <= 3 {
				call.ResultChan() <- hrpc.RPCResult{Error: RetryableError{}}
				return
			}
		}
		m.RLock()
		i := clients[c.addr]
		m.RUnlock()

		// if we are connected to this client the first time,
		// pretend it's down to fail the probe and start a reconnect
		if bytes.Equal(call.Table(), []byte("down")) {
			if i <= 1 {
				call.ResultChan() <- hrpc.RPCResult{Error: UnrecoverableError{}}
			} else {
				// otherwise, the region is fine
				call.ResultChan() <- hrpc.RPCResult{}
			}
			return
		}
	}
	if bytes.HasSuffix(call.Key(), bytes.Repeat([]byte{0}, 17)) {
		// meta region probe, return empty to signify that region is online
		call.ResultChan() <- hrpc.RPCResult{}
	} else if bytes.HasPrefix(call.Key(), []byte("test,")) {
		call.ResultChan() <- hrpc.RPCResult{Msg: &pb.ScanResponse{
			Results: []*pb.Result{metaRow}}}
	} else if bytes.HasPrefix(call.Key(), []byte("test1,,")) {
		call.ResultChan() <- hrpc.RPCResult{Msg: &pb.ScanResponse{
			Results: []*pb.Result{test1SplitA}}}
	} else if bytes.HasPrefix(call.Key(), []byte("nsre,,")) {
		call.ResultChan() <- hrpc.RPCResult{Msg: &pb.ScanResponse{
			Results: []*pb.Result{nsreRegion}}}
	} else if bytes.HasPrefix(call.Key(), []byte("tablenotfound,")) {
		call.ResultChan() <- hrpc.RPCResult{Msg: &pb.ScanResponse{
			Results:     []*pb.Result{},
			MoreResults: proto.Bool(false),
		}}
	} else {
		call.ResultChan() <- hrpc.RPCResult{Msg: makeRegionResult(call.Key())}
	}
}

func (c *testClient) Close() {}
