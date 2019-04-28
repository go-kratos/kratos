// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gohbase

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/cznic/b"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
	"github.com/tsuna/gohbase/region"
	"github.com/tsuna/gohbase/zk"
	"golang.org/x/time/rate"
)

const (
	standardClient = iota
	adminClient
	defaultRPCQueueSize  = 100
	defaultFlushInterval = 20 * time.Millisecond
	defaultZkRoot        = "/hbase"
	defaultZkTimeout     = 30 * time.Second
	defaultEffectiveUser = "root"
	// metaBurst is maximum number of request allowed at once.
	metaBurst = 10
	// metaLimit is rate at which to throttle requests to hbase:meta table.
	metaLimit = rate.Limit(100)
)

// Client a regular HBase client
type Client interface {
	Scan(s *hrpc.Scan) hrpc.Scanner
	Get(g *hrpc.Get) (*hrpc.Result, error)
	Put(p *hrpc.Mutate) (*hrpc.Result, error)
	Delete(d *hrpc.Mutate) (*hrpc.Result, error)
	Append(a *hrpc.Mutate) (*hrpc.Result, error)
	Increment(i *hrpc.Mutate) (int64, error)
	CheckAndPut(p *hrpc.Mutate, family string, qualifier string,
		expectedValue []byte) (bool, error)
	Close()
}

// RPCClient is core client of gohbase. It's exposed for testing.
type RPCClient interface {
	SendRPC(rpc hrpc.Call) (proto.Message, error)
}

// Option is a function used to configure optional config items for a Client.
type Option func(*client)

// A Client provides access to an HBase cluster.
type client struct {
	clientType int

	regions keyRegionCache

	// Maps a hrpc.RegionInfo to the *region.Client that we think currently
	// serves it.
	clients clientRegionCache

	metaRegionInfo hrpc.RegionInfo

	adminRegionInfo hrpc.RegionInfo

	// The maximum size of the RPC queue in the region client
	rpcQueueSize int

	// zkClient is zookeeper for retrieving meta and admin information
	zkClient zk.Client

	// The root zookeeper path for Hbase. By default, this is usually "/hbase".
	zkRoot string

	// The zookeeper session timeout
	zkTimeout time.Duration

	// The timeout before flushing the RPC queue in the region client
	flushInterval time.Duration

	// The user used when accessing regions.
	effectiveUser string

	// metaLookupLimiter is used to throttle lookups to hbase:meta table
	metaLookupLimiter *rate.Limiter

	// How long to wait for a region lookup (either meta lookup or finding
	// meta in ZooKeeper).  Should be greater than or equal to the ZooKeeper
	// session timeout.
	regionLookupTimeout time.Duration

	// regionReadTimeout is the maximum amount of time to wait for regionserver reply
	regionReadTimeout time.Duration

	done      chan struct{}
	closeOnce sync.Once
}

// NewClient creates a new HBase client.
func NewClient(zkquorum string, options ...Option) Client {
	return newClient(zkquorum, options...)
}

func newClient(zkquorum string, options ...Option) *client {
	log.WithFields(log.Fields{
		"Host": zkquorum,
	}).Debug("Creating new client.")
	c := &client{
		clientType: standardClient,
		regions:    keyRegionCache{regions: b.TreeNew(region.CompareGeneric)},
		clients: clientRegionCache{
			regions: make(map[hrpc.RegionClient]map[hrpc.RegionInfo]struct{}),
		},
		rpcQueueSize:  defaultRPCQueueSize,
		flushInterval: defaultFlushInterval,
		metaRegionInfo: region.NewInfo(
			0,
			[]byte("hbase"),
			[]byte("meta"),
			[]byte("hbase:meta,,1"),
			nil,
			nil),
		zkRoot:              defaultZkRoot,
		zkTimeout:           defaultZkTimeout,
		effectiveUser:       defaultEffectiveUser,
		metaLookupLimiter:   rate.NewLimiter(metaLimit, metaBurst),
		regionLookupTimeout: region.DefaultLookupTimeout,
		regionReadTimeout:   region.DefaultReadTimeout,
		done:                make(chan struct{}),
	}
	for _, option := range options {
		option(c)
	}

	//Have to create the zkClient after the Options have been set
	//since the zkTimeout could be changed as an option
	c.zkClient = zk.NewClient(zkquorum, c.zkTimeout)

	return c
}

// RpcQueueSize will return an option that will set the size of the RPC queues
// used in a given client
func RpcQueueSize(size int) Option {
	return func(c *client) {
		c.rpcQueueSize = size
	}
}

// ZookeeperRoot will return an option that will set the zookeeper root path used in a given client.
func ZookeeperRoot(root string) Option {
	return func(c *client) {
		c.zkRoot = root
	}
}

// ZookeeperTimeout will return an option that will set the zookeeper session timeout.
func ZookeeperTimeout(to time.Duration) Option {
	return func(c *client) {
		c.zkTimeout = to
	}
}

// RegionLookupTimeout will return an option that sets the region lookup timeout
func RegionLookupTimeout(to time.Duration) Option {
	return func(c *client) {
		c.regionLookupTimeout = to
	}
}

// RegionReadTimeout will return an option that sets the region read timeout
func RegionReadTimeout(to time.Duration) Option {
	return func(c *client) {
		c.regionReadTimeout = to
	}
}

// EffectiveUser will return an option that will set the user used when accessing regions.
func EffectiveUser(user string) Option {
	return func(c *client) {
		c.effectiveUser = user
	}
}

// FlushInterval will return an option that will set the timeout for flushing
// the RPC queues used in a given client
func FlushInterval(interval time.Duration) Option {
	return func(c *client) {
		c.flushInterval = interval
	}
}

// Close closes connections to hbase master and regionservers
func (c *client) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		if c.clientType == adminClient {
			if ac := c.adminRegionInfo.Client(); ac != nil {
				ac.Close()
			}
		}
		c.clients.closeAll()
	})
}

func (c *client) Scan(s *hrpc.Scan) hrpc.Scanner {
	return newScanner(c, s)
}

func (c *client) Get(g *hrpc.Get) (*hrpc.Result, error) {
	pbmsg, err := c.SendRPC(g)
	if err != nil {
		return nil, err
	}

	r, ok := pbmsg.(*pb.GetResponse)
	if !ok {
		return nil, fmt.Errorf("sendRPC returned not a GetResponse")
	}

	return hrpc.ToLocalResult(r.Result), nil
}

func (c *client) Put(p *hrpc.Mutate) (*hrpc.Result, error) {
	return c.mutate(p)
}

func (c *client) Delete(d *hrpc.Mutate) (*hrpc.Result, error) {
	return c.mutate(d)
}

func (c *client) Append(a *hrpc.Mutate) (*hrpc.Result, error) {
	return c.mutate(a)
}

func (c *client) Increment(i *hrpc.Mutate) (int64, error) {
	r, err := c.mutate(i)
	if err != nil {
		return 0, err
	}

	if len(r.Cells) != 1 {
		return 0, fmt.Errorf("increment returned %d cells, but we expected exactly one",
			len(r.Cells))
	}

	val := binary.BigEndian.Uint64(r.Cells[0].Value)
	return int64(val), nil
}

func (c *client) mutate(m *hrpc.Mutate) (*hrpc.Result, error) {
	pbmsg, err := c.SendRPC(m)
	if err != nil {
		return nil, err
	}

	r, ok := pbmsg.(*pb.MutateResponse)
	if !ok {
		return nil, fmt.Errorf("sendRPC returned not a MutateResponse")
	}

	return hrpc.ToLocalResult(r.Result), nil
}

func (c *client) CheckAndPut(p *hrpc.Mutate, family string,
	qualifier string, expectedValue []byte) (bool, error) {
	cas, err := hrpc.NewCheckAndPut(p, family, qualifier, expectedValue)
	if err != nil {
		return false, err
	}

	pbmsg, err := c.SendRPC(cas)
	if err != nil {
		return false, err
	}

	r, ok := pbmsg.(*pb.MutateResponse)
	if !ok {
		return false, fmt.Errorf("sendRPC returned a %T instead of MutateResponse", pbmsg)
	}

	if r.Processed == nil {
		return false, fmt.Errorf("protobuf in the response didn't contain the field "+
			"indicating whether the CheckAndPut was successful or not: %s", r)
	}

	return r.GetProcessed(), nil
}
