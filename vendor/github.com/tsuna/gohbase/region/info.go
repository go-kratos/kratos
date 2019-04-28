// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package region contains data structures to represent HBase regions.
package region

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
)

var defaultNamespace = []byte("default")

// OfflineRegionError is returned if region is offline
type OfflineRegionError struct {
	n string
}

func (e OfflineRegionError) Error() string {
	return fmt.Sprintf("region %s is offline", e.n)
}

// info describes a region.
type info struct {
	id        uint64 // A timestamp when the region is created
	namespace []byte
	table     []byte
	name      []byte
	startKey  []byte
	stopKey   []byte
	ctx       context.Context
	cancel    context.CancelFunc

	// The attributes before this mutex are supposed to be immutable.
	// The attributes defined below can be changed and accesses must
	// be protected with this mutex.
	m sync.RWMutex

	client hrpc.RegionClient

	// Once a region becomes unreachable, this channel is created, and any
	// functions that wish to be notified when the region becomes available
	// again can read from this channel, which will be closed when the region
	// is available again
	available chan struct{}
}

// NewInfo creates a new region info
func NewInfo(id uint64, namespace, table, name, startKey, stopKey []byte) hrpc.RegionInfo {
	ctx, cancel := context.WithCancel(context.Background())
	return &info{
		id:        id,
		ctx:       ctx,
		cancel:    cancel,
		namespace: namespace,
		table:     table,
		name:      name,
		startKey:  startKey,
		stopKey:   stopKey,
	}
}

// infoFromCell parses a KeyValue from the meta table and creates the
// corresponding Info object.
func infoFromCell(cell *hrpc.Cell) (hrpc.RegionInfo, error) {
	value := cell.Value
	if len(value) == 0 {
		return nil, fmt.Errorf("empty value in %q", cell)
	} else if value[0] != 'P' {
		return nil, fmt.Errorf("unsupported region info version %d in %q", value[0], cell)
	}
	const pbufMagic = 1346524486 // 4 bytes: "PBUF"
	magic := binary.BigEndian.Uint32(value[:4])
	if magic != pbufMagic {
		return nil, fmt.Errorf("invalid magic number in %q", cell)
	}
	var regInfo pb.RegionInfo
	err := proto.UnmarshalMerge(value[4:], &regInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q: %s", cell, err)
	}
	if regInfo.GetOffline() {
		return nil, OfflineRegionError{n: string(cell.Row)}
	}
	var namespace []byte
	if !bytes.Equal(regInfo.TableName.Namespace, defaultNamespace) {
		// if default namespace, pretend there's no namespace
		namespace = regInfo.TableName.Namespace
	}

	return NewInfo(
		regInfo.GetRegionId(),
		namespace,
		regInfo.TableName.Qualifier,
		cell.Row,
		regInfo.StartKey,
		regInfo.EndKey,
	), nil
}

// ParseRegionInfo parses the contents of a row from the meta table.
// It's guaranteed to return a region info and a host:port OR return an error.
func ParseRegionInfo(metaRow *hrpc.Result) (hrpc.RegionInfo, string, error) {
	var reg hrpc.RegionInfo
	var addr string

	for _, cell := range metaRow.Cells {
		switch string(cell.Qualifier) {
		case "regioninfo":
			var err error
			reg, err = infoFromCell(cell)
			if err != nil {
				return nil, "", err
			}
		case "server":
			value := cell.Value
			if len(value) == 0 {
				continue // Empty during NSRE.
			}
			addr = string(value)
		default:
			// Other kinds of qualifiers: ignore them.
			// TODO: If this is the parent of a split region, there are two other
			// KVs that could be useful: `info:splitA' and `info:splitB'.
			// Need to investigate whether we can use those as a hint to update our
			// regions_cache with the daughter regions of the split.
		}
	}

	if reg == nil {
		// There was no region in the row in meta, this is really not expected.
		return nil, "", fmt.Errorf("meta seems to be broken, there was no region in %v", metaRow)
	}
	if len(addr) == 0 {
		return nil, "", fmt.Errorf("meta doesn't have a server location in %v", metaRow)
	}
	return reg, addr, nil
}

// IsUnavailable returns true if this region has been marked as unavailable.
func (i *info) IsUnavailable() bool {
	i.m.RLock()
	res := i.available != nil
	i.m.RUnlock()
	return res
}

// AvailabilityChan returns a channel that can be used to wait on for
// notification that a connection to this region has been reestablished.
// If this region is not marked as unavailable, nil will be returned.
func (i *info) AvailabilityChan() <-chan struct{} {
	i.m.RLock()
	ch := i.available
	i.m.RUnlock()
	return ch
}

// MarkUnavailable will mark this region as unavailable, by creating the struct
// returned by AvailabilityChan. If this region was marked as available
// before this, true will be returned.
func (i *info) MarkUnavailable() bool {
	created := false
	i.m.Lock()
	if i.available == nil {
		i.available = make(chan struct{})
		created = true
	}
	i.m.Unlock()
	return created
}

// MarkAvailable will mark this region as available again, by closing the struct
// returned by AvailabilityChan
func (i *info) MarkAvailable() {
	i.m.Lock()
	ch := i.available
	i.available = nil
	close(ch)
	i.m.Unlock()
}

// MarkDead will mark this region as not useful anymore to notify everyone
// who's trying to use it that there's no point
func (i *info) MarkDead() {
	i.cancel()
}

// Context to check if the region is dead
func (i *info) Context() context.Context {
	return i.ctx
}

func (i *info) String() string {
	return fmt.Sprintf(
		"RegionInfo{Name: %q, ID: %d, Namespace: %q, Table: %q, StartKey: %q, StopKey: %q}",
		i.name, i.id, i.namespace, i.table, i.startKey, i.stopKey)
}

// ID returns region's age
func (i *info) ID() uint64 {
	return i.id
}

// Name returns region name
func (i *info) Name() []byte {
	return i.name
}

// StopKey return region stop key
func (i *info) StopKey() []byte {
	return i.stopKey
}

// StartKey return region start key
func (i *info) StartKey() []byte {
	return i.startKey
}

// Namespace returns region table
func (i *info) Namespace() []byte {
	return i.namespace
}

// Table returns region table
func (i *info) Table() []byte {
	return i.table
}

// Client returns region client
func (i *info) Client() hrpc.RegionClient {
	i.m.RLock()
	c := i.client
	i.m.RUnlock()
	return c
}

// SetClient sets region client
func (i *info) SetClient(c hrpc.RegionClient) {
	i.m.Lock()
	i.client = c
	i.m.Unlock()
}

// CompareGeneric is the same thing as Compare but for interface{}.
func CompareGeneric(a, b interface{}) int {
	return Compare(a.([]byte), b.([]byte))
}

// Compare compares two region names.
// We can't just use bytes.Compare() because it doesn't play nicely
// with the way META keys are built as the first region has an empty start
// key.  Let's assume we know about those 2 regions in our cache:
//   .META.,,1
//   tableA,,1273018455182
// We're given an RPC to execute on "tableA", row "\x00" (1 byte row key
// containing a 0).  If we use Compare() to sort the entries in the cache,
// when we search for the entry right before "tableA,\000,:"
// we'll erroneously find ".META.,,1" instead of the entry for first
// region of "tableA".
//
// Since this scheme breaks natural ordering, we need this comparator to
// implement a special version of comparison to handle this scenario.
func Compare(a, b []byte) int {
	var length int
	if la, lb := len(a), len(b); la < lb {
		length = la
	} else {
		length = lb
	}
	// Reminder: region names are of the form:
	//   table_name,start_key,timestamp[.MD5.]
	// First compare the table names.
	var i int
	for i = 0; i < length; i++ {
		ai := a[i]    // Saves one pointer deference every iteration.
		bi := b[i]    // Saves one pointer deference every iteration.
		if ai != bi { // The name of the tables differ.
			if ai == ',' {
				return -1001 // `a' has a smaller table name.  a < b
			} else if bi == ',' {
				return 1001 // `b' has a smaller table name.  a > b
			}
			return int(ai) - int(bi)
		}
		if ai == ',' { // Remember: at this point ai == bi.
			break // We're done comparing the table names.  They're equal.
		}
	}

	// Now find the last comma in both `a' and `b'.  We need to start the
	// search from the end as the row key could have an arbitrary number of
	// commas and we don't know its length.
	aComma := findCommaFromEnd(a, i)
	bComma := findCommaFromEnd(b, i)
	// If either `a' or `b' is followed immediately by another comma, then
	// they are the first region (it's the empty start key).
	i++ // No need to check against `length', there MUST be more bytes.

	// Compare keys.
	var firstComma int
	if aComma < bComma {
		firstComma = aComma
	} else {
		firstComma = bComma
	}
	for ; i < firstComma; i++ {
		ai := a[i]
		bi := b[i]
		if ai != bi { // The keys differ.
			return int(ai) - int(bi)
		}
	}
	if aComma < bComma {
		return -1002 // `a' has a shorter key.  a < b
	} else if bComma < aComma {
		return 1002 // `b' has a shorter key.  a > b
	}

	// Keys have the same length and have compared identical.  Compare the
	// rest, which essentially means: use start code as a tie breaker.
	for ; /*nothing*/ i < length; i++ {
		ai := a[i]
		bi := b[i]
		if ai != bi { // The start codes differ.
			return int(ai) - int(bi)
		}
	}

	return len(a) - len(b)
}

// Because there is no `LastIndexByte()' in the standard `bytes' package.
func findCommaFromEnd(b []byte, offset int) int {
	for i := len(b) - 1; i > offset; i-- {
		if b[i] == ',' {
			return i
		}
	}
	panic(fmt.Errorf("no comma found in %q after offset %d", b, offset))
}
