// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package zk encapsulates our interactions with ZooKeeper.
package zk

import (
	"encoding/binary"
	"fmt"
	"net"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/tsuna/gohbase/pb"
)

type logger struct{}

func (l *logger) Printf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func init() {
	zk.DefaultLogger = &logger{}
}

// ResourceName is a type alias that is used to represent different resources
// in ZooKeeper
type ResourceName string

// Prepend creates a new ResourceName with prefix prepended to the former ResourceName.
func (r ResourceName) Prepend(prefix string) ResourceName {
	return ResourceName(path.Join(prefix, string(r)))
}

const (
	// Meta is a ResourceName that indicates that the location of the Meta
	// table is what will be fetched
	Meta = ResourceName("/meta-region-server")

	// Master is a ResourceName that indicates that the location of the Master
	// server is what will be fetched
	Master = ResourceName("/master")
)

// Client is an interface of client that retrieves meta infomation from zookeeper
type Client interface {
	LocateResource(ResourceName) (string, error)
}

type client struct {
	zks            []string
	sessionTimeout time.Duration
}

// NewClient establishes connection to zookeeper and returns the client
func NewClient(zkquorum string, st time.Duration) Client {
	return &client{
		zks:            strings.Split(zkquorum, ","),
		sessionTimeout: st,
	}
}

// LocateResource returns address of the server for the specified resource.
func (c *client) LocateResource(resource ResourceName) (string, error) {
	conn, _, err := zk.Connect(c.zks, c.sessionTimeout)
	if err != nil {
		return "", fmt.Errorf("error connecting to ZooKeeper at %v: %s", c.zks, err)
	}
	defer conn.Close()

	buf, _, err := conn.Get(string(resource))
	if err != nil {
		return "", fmt.Errorf("failed to read the %s znode: %s", resource, err)
	}
	if len(buf) == 0 {
		log.Fatalf("%s was empty!", resource)
	} else if buf[0] != 0xFF {
		return "", fmt.Errorf("the first byte of %s was 0x%x, not 0xFF", resource, buf[0])
	}
	metadataLen := binary.BigEndian.Uint32(buf[1:])
	if metadataLen < 1 || metadataLen > 65000 {
		return "", fmt.Errorf("invalid metadata length for %s: %d", resource, metadataLen)
	}
	buf = buf[1+4+metadataLen:]
	magic := binary.BigEndian.Uint32(buf)
	const pbufMagic = 1346524486 // 4 bytes: "PBUF"
	if magic != pbufMagic {
		return "", fmt.Errorf("invalid magic number for %s: %d", resource, magic)
	}
	buf = buf[4:]
	var server *pb.ServerName
	if resource == Meta {
		meta := &pb.MetaRegionServer{}
		err = proto.UnmarshalMerge(buf, meta)
		if err != nil {
			return "",
				fmt.Errorf("failed to deserialize the MetaRegionServer entry from ZK: %s", err)
		}
		server = meta.Server
	} else {
		master := &pb.Master{}
		err = proto.UnmarshalMerge(buf, master)
		if err != nil {
			return "",
				fmt.Errorf("failed to deserialize the Master entry from ZK: %s", err)
		}
		server = master.Master
	}
	return net.JoinHostPort(*server.HostName, fmt.Sprint(*server.Port)), nil
}
