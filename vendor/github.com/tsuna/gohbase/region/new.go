// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// +build !testing

package region

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tsuna/gohbase/hrpc"
)

// NewClient creates a new RegionClient.
func NewClient(ctx context.Context, addr string, ctype ClientType,
	queueSize int, flushInterval time.Duration, effectiveUser string,
	readTimeout time.Duration) (hrpc.RegionClient, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the RegionServer at %s: %s", addr, err)
	}
	c := &client{
		addr:          addr,
		conn:          conn,
		rpcs:          make(chan hrpc.Call),
		done:          make(chan struct{}),
		sent:          make(map[uint32]hrpc.Call),
		rpcQueueSize:  queueSize,
		flushInterval: flushInterval,
		effectiveUser: effectiveUser,
		readTimeout:   readTimeout,
	}
	// time out send hello if it take long
	// TODO: do we even need to bother, we are going to retry anyway?
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetWriteDeadline(deadline)
	}
	if err := c.sendHello(ctype); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to send hello to the RegionServer at %s: %s", addr, err)
	}
	// reset write deadline
	conn.SetWriteDeadline(time.Time{})

	if ctype == RegionClient {
		go c.processRPCs() // Batching goroutine
	}
	go c.receiveRPCs() // Reader goroutine
	return c, nil
}
