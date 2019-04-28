// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package hrpc

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/pb"
)

// Get represents a Get HBase call.
type Get struct {
	base
	baseQuery
	// Don't return any KeyValue, just say whether the row key exists in the
	// table or not.
	existsOnly bool
	skipbatch  bool
}

// baseGet returns a Get struct with default values set.
func baseGet(ctx context.Context, table []byte, key []byte,
	options ...func(Call) error) (*Get, error) {
	g := &Get{
		base: base{
			key:      key,
			table:    table,
			ctx:      ctx,
			resultch: make(chan RPCResult, 1),
		},
		baseQuery: newBaseQuery(),
	}
	err := applyOptions(g, options...)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// NewGet creates a new Get request for the given table and row key.
func NewGet(ctx context.Context, table, key []byte,
	options ...func(Call) error) (*Get, error) {
	return baseGet(ctx, table, key, options...)
}

// NewGetStr creates a new Get request for the given table and row key.
func NewGetStr(ctx context.Context, table, key string,
	options ...func(Call) error) (*Get, error) {
	return NewGet(ctx, []byte(table), []byte(key), options...)
}

// Name returns the name of this RPC call.
func (g *Get) Name() string {
	return "Get"
}

// SkipBatch returns true if the Get request shouldn't be batched,
// but should be sent to Region Server right away.
func (g *Get) SkipBatch() bool {
	return g.skipbatch
}

func (g *Get) setSkipBatch(v bool) {
	g.skipbatch = v
}

// ExistsOnly makes this Get request not return any KeyValue, merely whether
// or not the given row key exists in the table.
func (g *Get) ExistsOnly() {
	g.existsOnly = true
}

// ToProto converts this RPC into a protobuf message.
func (g *Get) ToProto() proto.Message {
	get := &pb.GetRequest{
		Region: g.regionSpecifier(),
		Get: &pb.Get{
			Row:       g.key,
			Column:    familiesToColumn(g.families),
			TimeRange: &pb.TimeRange{},
		},
	}

	/* added support for limit number of cells per row */
	if g.storeLimit != DefaultMaxResultsPerColumnFamily {
		get.Get.StoreLimit = &g.storeLimit
	}
	if g.storeOffset != 0 {
		get.Get.StoreOffset = &g.storeOffset
	}

	if g.maxVersions != DefaultMaxVersions {
		get.Get.MaxVersions = &g.maxVersions
	}
	if g.fromTimestamp != MinTimestamp {
		get.Get.TimeRange.From = &g.fromTimestamp
	}
	if g.toTimestamp != MaxTimestamp {
		get.Get.TimeRange.To = &g.toTimestamp
	}
	if g.existsOnly {
		get.Get.ExistenceOnly = proto.Bool(true)
	}
	get.Get.Filter = g.filter
	return get
}

// NewResponse creates an empty protobuf message to read the response of this
// RPC.
func (g *Get) NewResponse() proto.Message {
	return &pb.GetResponse{}
}

// DeserializeCellBlocks deserializes get result from cell blocks
func (g *Get) DeserializeCellBlocks(m proto.Message, b []byte) (uint32, error) {
	resp := m.(*pb.GetResponse)
	if resp.Result == nil {
		// TODO: is this possible?
		return 0, nil
	}
	cells, read, err := deserializeCellBlocks(b, uint32(resp.Result.GetAssociatedCellCount()))
	if err != nil {
		return 0, err
	}
	resp.Result.Cell = append(resp.Result.Cell, cells...)
	return read, nil
}

// familiesToColumn takes a map from strings to lists of strings, and converts
// them into protobuf Columns
func familiesToColumn(families map[string][]string) []*pb.Column {
	cols := make([]*pb.Column, len(families))
	counter := 0
	for family, qualifiers := range families {
		bytequals := make([][]byte, len(qualifiers))
		for i, qual := range qualifiers {
			bytequals[i] = []byte(qual)
		}
		cols[counter] = &pb.Column{
			Family:    []byte(family),
			Qualifier: bytequals,
		}
		counter++
	}
	return cols
}
