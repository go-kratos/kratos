// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package hrpc

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/pb"
)

// RegionInfo represents HBase region.
type RegionInfo interface {
	IsUnavailable() bool
	AvailabilityChan() <-chan struct{}
	MarkUnavailable() bool
	MarkAvailable()
	MarkDead()
	Context() context.Context
	String() string
	ID() uint64
	Name() []byte
	StartKey() []byte
	StopKey() []byte
	Namespace() []byte
	Table() []byte
	SetClient(RegionClient)
	Client() RegionClient
}

// RegionClient represents HBase region client.
type RegionClient interface {
	Close()
	Addr() string
	QueueRPC(Call)
	String() string
}

// Call represents an HBase RPC call.
type Call interface {
	Table() []byte
	Name() string
	Key() []byte
	Region() RegionInfo
	SetRegion(region RegionInfo)
	ToProto() proto.Message
	// Returns a newly created (default-state) protobuf in which to store the
	// response of this call.
	NewResponse() proto.Message
	ResultChan() chan RPCResult
	Context() context.Context
}

type withOptions interface {
	Options() []func(Call) error
	setOptions([]func(Call) error)
}

// Batchable interface should be implemented by calls that can be batched into MultiRequest
type Batchable interface {
	// SkipBatch returns true if a call shouldn't be batched into MultiRequest and
	// should be sent right away.
	SkipBatch() bool

	setSkipBatch(v bool)
}

// SkipBatch is an option for batchable requests (Get and Mutate) to tell
// the client to skip batching and just send the request to Region Server
// right away.
func SkipBatch() func(Call) error {
	return func(c Call) error {
		if b, ok := c.(Batchable); ok {
			b.setSkipBatch(true)
			return nil
		}
		return errors.New("'SkipBatch' option only works with Get and Mutate requests")
	}
}

// hasQueryOptions is interface that needs to be implemented by calls
// that allow to provide Families and Filters options.
type hasQueryOptions interface {
	setFamilies(families map[string][]string)
	setFilter(filter *pb.Filter)
	setTimeRangeUint64(from, to uint64)
	setMaxVersions(versions uint32)
	setMaxResultsPerColumnFamily(maxresults uint32)
	setResultOffset(offset uint32)
}

// RPCResult is struct that will contain both the resulting message from an RPC
// call, and any errors that may have occurred related to making the RPC call.
type RPCResult struct {
	Msg   proto.Message
	Error error
}

type base struct {
	ctx     context.Context
	table   []byte
	key     []byte
	options []func(Call) error

	region   RegionInfo
	resultch chan RPCResult
}

func (b *base) Context() context.Context {
	return b.ctx
}

func (b *base) Region() RegionInfo {
	return b.region
}

func (b *base) SetRegion(region RegionInfo) {
	b.region = region
}

func (b *base) regionSpecifier() *pb.RegionSpecifier {
	return &pb.RegionSpecifier{
		Type:  pb.RegionSpecifier_REGION_NAME.Enum(),
		Value: []byte(b.region.Name()),
	}
}

func (b *base) setOptions(options []func(Call) error) {
	b.options = options
}

// Options returns all the options passed to this call
func (b *base) Options() []func(Call) error {
	return b.options
}

func applyOptions(call Call, options ...func(Call) error) error {
	call.(withOptions).setOptions(options)
	for _, option := range options {
		err := option(call)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *base) Table() []byte {
	return b.table
}

func (b *base) Key() []byte {
	return b.key
}

func (b *base) ResultChan() chan RPCResult {
	return b.resultch
}

// Cell is the smallest level of granularity in returned results.
// Represents a single cell in HBase (a row will have one cell for every qualifier).
type Cell pb.Cell

func (c *Cell) String() string {
	return (*pb.Cell)(c).String()
}

// cellFromCellBlock deserializes a cell from a reader
func cellFromCellBlock(b []byte) (*pb.Cell, uint32, error) {
	if len(b) < 4 {
		return nil, 0, fmt.Errorf(
			"buffer is too small: expected %d, got %d", 4, len(b))
	}

	kvLen := binary.BigEndian.Uint32(b[0:4])
	if len(b) < int(kvLen)+4 {
		return nil, 0, fmt.Errorf(
			"buffer is too small: expected %d, got %d", int(kvLen)+4, len(b))
	}

	rowKeyLen := binary.BigEndian.Uint32(b[4:8])
	valueLen := binary.BigEndian.Uint32(b[8:12])
	keyLen := binary.BigEndian.Uint16(b[12:14])
	b = b[14:]

	key := b[:keyLen]
	b = b[keyLen:]

	familyLen := uint8(b[0])
	b = b[1:]

	family := b[:familyLen]
	b = b[familyLen:]

	qualifierLen := rowKeyLen - uint32(keyLen) - uint32(familyLen) - 2 - 1 - 8 - 1
	if 4 /*rowKeyLen*/ +4 /*valueLen*/ +2 /*keyLen*/ +
		uint32(keyLen)+1 /*familyLen*/ +uint32(familyLen)+qualifierLen+
		8 /*timestamp*/ +1 /*cellType*/ +valueLen != kvLen {
		return nil, 0, fmt.Errorf("HBase has lied about KeyValue length: expected %d, got %d",
			kvLen, 4+4+2+uint32(keyLen)+1+uint32(familyLen)+qualifierLen+8+1+valueLen)
	}
	qualifier := b[:qualifierLen]
	b = b[qualifierLen:]

	timestamp := binary.BigEndian.Uint64(b[:8])
	b = b[8:]

	cellType := uint8(b[0])
	b = b[1:]

	value := b[:valueLen]

	return &pb.Cell{
		Row:       key,
		Family:    family,
		Qualifier: qualifier,
		Timestamp: &timestamp,
		Value:     value,
		CellType:  pb.CellType(cellType).Enum(),
	}, kvLen + 4, nil
}

func deserializeCellBlocks(b []byte, cellsLen uint32) ([]*pb.Cell, uint32, error) {
	cells := make([]*pb.Cell, cellsLen)
	var readLen uint32
	for i := 0; i < int(cellsLen); i++ {
		c, l, err := cellFromCellBlock(b[readLen:])
		if err != nil {
			return nil, readLen, err
		}
		cells[i] = c
		readLen += l
	}
	return cells, readLen, nil
}

// Result holds a slice of Cells as well as miscellaneous information about the response.
type Result struct {
	Cells   []*Cell
	Stale   bool
	Partial bool
	// Exists is only set if existance_only was set in the request query.
	Exists *bool
}

func extractBool(v *bool) bool {
	return v != nil && *v
}

// ToLocalResult takes a protobuf Result type and converts it to our own
// Result type in constant time.
func ToLocalResult(pbr *pb.Result) *Result {
	if pbr == nil {
		return &Result{}
	}
	return &Result{
		// Should all be O(1) operations.
		Cells:   toLocalCells(pbr),
		Stale:   extractBool(pbr.Stale),
		Partial: extractBool(pbr.Partial),
		Exists:  pbr.Exists,
	}
}

func toLocalCells(pbr *pb.Result) []*Cell {
	return *(*[]*Cell)(unsafe.Pointer(pbr))
}

// We can now define any helper functions on Result that we want.
