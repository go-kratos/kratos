// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package hrpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/pb"
)

// GetProcedureState represents a call to HBase to check status of a procedure
type GetProcedureState struct {
	base

	procID uint64
}

// NewGetProcedureState creates a new GetProcedureState request. For use by the admin client.
func NewGetProcedureState(ctx context.Context, procID uint64) *GetProcedureState {
	return &GetProcedureState{
		base: base{
			ctx:      ctx,
			resultch: make(chan RPCResult, 1),
		},
		procID: procID,
	}
}

// Name returns the name of this RPC call.
func (ps *GetProcedureState) Name() string {
	return "getProcedureResult"
}

// ToProto converts the RPC into a protobuf message
func (ps *GetProcedureState) ToProto() proto.Message {
	return &pb.GetProcedureResultRequest{ProcId: &ps.procID}
}

// NewResponse creates an empty protobuf message to read the response of this RPC.
func (ps *GetProcedureState) NewResponse() proto.Message {
	return &pb.GetProcedureResultResponse{}
}
