// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package hrpc

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/pb"
)

// CheckAndPut performs a provided Put operation if the value specified
// by condition equals to the one set in the HBase.
type CheckAndPut struct {
	*Mutate

	family    []byte
	qualifier []byte

	comparator *pb.Comparator
}

// NewCheckAndPut creates a new CheckAndPut request that will compare provided
// expectedValue with the on in HBase located at put's row and provided family:qualifier,
// and if they are equal, perform the provided put request on the row
func NewCheckAndPut(put *Mutate, family string,
	qualifier string, expectedValue []byte) (*CheckAndPut, error) {
	if put.mutationType != pb.MutationProto_PUT {
		return nil, fmt.Errorf("'CheckAndPut' only takes 'Put' request")
	}

	// The condition that needs to match for the edit to be applied.
	exp := filter.NewByteArrayComparable(expectedValue)
	cmp, err := filter.NewBinaryComparator(exp).ConstructPBComparator()
	if err != nil {
		return nil, err
	}

	// CheckAndPut is not batchable as MultiResponse doesn't return Processed field
	// for Mutate Action
	put.setSkipBatch(true)

	return &CheckAndPut{
		Mutate:     put,
		family:     []byte(family),
		qualifier:  []byte(qualifier),
		comparator: cmp,
	}, nil
}

// ToProto converts the RPC into a protobuf message
func (cp *CheckAndPut) ToProto() proto.Message {
	mutateRequest := cp.toProto()
	mutateRequest.Condition = &pb.Condition{
		Row:         cp.key,
		Family:      cp.family,
		Qualifier:   cp.qualifier,
		CompareType: pb.CompareType_EQUAL.Enum(),
		Comparator:  cp.comparator,
	}
	return mutateRequest
}
