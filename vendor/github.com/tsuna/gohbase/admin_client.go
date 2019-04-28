// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gohbase

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
	"github.com/tsuna/gohbase/region"
	"github.com/tsuna/gohbase/zk"
)

// AdminClient to perform admistrative operations with HMaster
type AdminClient interface {
	CreateTable(t *hrpc.CreateTable) error
	DeleteTable(t *hrpc.DeleteTable) error
	EnableTable(t *hrpc.EnableTable) error
	DisableTable(t *hrpc.DisableTable) error
	ClusterStatus() (*pb.ClusterStatus, error)
}

// NewAdminClient creates an admin HBase client.
func NewAdminClient(zkquorum string, options ...Option) AdminClient {
	return newAdminClient(zkquorum, options...)
}

func newAdminClient(zkquorum string, options ...Option) AdminClient {
	log.WithFields(log.Fields{
		"Host": zkquorum,
	}).Debug("Creating new admin client.")
	c := &client{
		clientType:    adminClient,
		rpcQueueSize:  defaultRPCQueueSize,
		flushInterval: defaultFlushInterval,
		// empty region in order to be able to set client to it
		adminRegionInfo:     region.NewInfo(0, nil, nil, nil, nil, nil),
		zkTimeout:           defaultZkTimeout,
		zkRoot:              defaultZkRoot,
		effectiveUser:       defaultEffectiveUser,
		regionLookupTimeout: region.DefaultLookupTimeout,
		regionReadTimeout:   region.DefaultReadTimeout,
	}
	for _, option := range options {
		option(c)
	}
	c.zkClient = zk.NewClient(zkquorum, c.zkTimeout)
	return c
}

//Get the status of the cluster
func (c *client) ClusterStatus() (*pb.ClusterStatus, error) {
	pbmsg, err := c.SendRPC(hrpc.NewClusterStatus())
	if err != nil {
		return nil, err
	}

	r, ok := pbmsg.(*pb.GetClusterStatusResponse)
	if !ok {
		return nil, fmt.Errorf("sendRPC returned not a ClusterStatusResponse")
	}

	return r.GetClusterStatus(), nil
}

func (c *client) CreateTable(t *hrpc.CreateTable) error {
	pbmsg, err := c.SendRPC(t)
	if err != nil {
		return err
	}

	r, ok := pbmsg.(*pb.CreateTableResponse)
	if !ok {
		return fmt.Errorf("sendRPC returned not a CreateTableResponse")
	}

	return c.checkProcedureWithBackoff(t.Context(), r.GetProcId())
}

func (c *client) DeleteTable(t *hrpc.DeleteTable) error {
	pbmsg, err := c.SendRPC(t)
	if err != nil {
		return err
	}

	r, ok := pbmsg.(*pb.DeleteTableResponse)
	if !ok {
		return fmt.Errorf("sendRPC returned not a DeleteTableResponse")
	}

	return c.checkProcedureWithBackoff(t.Context(), r.GetProcId())
}

func (c *client) EnableTable(t *hrpc.EnableTable) error {
	pbmsg, err := c.SendRPC(t)
	if err != nil {
		return err
	}

	r, ok := pbmsg.(*pb.EnableTableResponse)
	if !ok {
		return fmt.Errorf("sendRPC returned not a EnableTableResponse")
	}

	return c.checkProcedureWithBackoff(t.Context(), r.GetProcId())
}

func (c *client) DisableTable(t *hrpc.DisableTable) error {
	pbmsg, err := c.SendRPC(t)
	if err != nil {
		return err
	}

	r, ok := pbmsg.(*pb.DisableTableResponse)
	if !ok {
		return fmt.Errorf("sendRPC returned not a DisableTableResponse")
	}

	return c.checkProcedureWithBackoff(t.Context(), r.GetProcId())
}

func (c *client) checkProcedureWithBackoff(ctx context.Context, procID uint64) error {
	backoff := backoffStart
	for {
		pbmsg, err := c.SendRPC(hrpc.NewGetProcedureState(ctx, procID))
		if err != nil {
			return err
		}

		res := pbmsg.(*pb.GetProcedureResultResponse)
		switch res.GetState() {
		case pb.GetProcedureResultResponse_NOT_FOUND:
			return fmt.Errorf("procedure not found")
		case pb.GetProcedureResultResponse_FINISHED:
			if fe := res.Exception; fe != nil {
				ge := fe.GenericException
				if ge == nil {
					return errors.New("got unexpected empty exception")
				}
				return fmt.Errorf("procedure exception: %s: %s", ge.GetClassName(), ge.GetMessage())
			}
			return nil
		default:
			backoff, err = sleepAndIncreaseBackoff(ctx, backoff)
			if err != nil {
				return err
			}
		}
	}
}
