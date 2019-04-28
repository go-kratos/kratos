// Copyright (C) 2017  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package hrpc

import (
	"errors"
	"math"
	"time"

	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/pb"
)

// baseQuery bundles common fields that can be provided for quering requests: Scans and Gets
type baseQuery struct {
	families      map[string][]string
	filter        *pb.Filter
	fromTimestamp uint64
	toTimestamp   uint64
	maxVersions   uint32
	storeLimit    uint32
	storeOffset   uint32
}

// newBaseQuery return baseQuery with all default values
func newBaseQuery() baseQuery {
	return baseQuery{
		storeLimit:    DefaultMaxResultsPerColumnFamily,
		fromTimestamp: MinTimestamp,
		toTimestamp:   MaxTimestamp,
		maxVersions:   DefaultMaxVersions,
	}
}

func (bq *baseQuery) setFamilies(families map[string][]string) {
	bq.families = families
}
func (bq *baseQuery) setFilter(filter *pb.Filter) {
	bq.filter = filter
}
func (bq *baseQuery) setTimeRangeUint64(from, to uint64) {
	bq.fromTimestamp = from
	bq.toTimestamp = to
}
func (bq *baseQuery) setMaxVersions(versions uint32) {
	bq.maxVersions = versions
}
func (bq *baseQuery) setMaxResultsPerColumnFamily(maxresults uint32) {
	bq.storeLimit = maxresults
}
func (bq *baseQuery) setResultOffset(offset uint32) {
	bq.storeOffset = offset
}

// Families option adds families constraint to a Scan or Get request.
func Families(f map[string][]string) func(Call) error {
	return func(hc Call) error {
		if c, ok := hc.(hasQueryOptions); ok {
			c.setFamilies(f)
			return nil
		}
		return errors.New("'Families' option can only be used with Get or Scan request")
	}
}

// Filters option adds filters constraint to a Scan or Get request.
func Filters(f filter.Filter) func(Call) error {
	return func(hc Call) error {
		if c, ok := hc.(hasQueryOptions); ok {
			pbF, err := f.ConstructPBFilter()
			if err != nil {
				return err
			}
			c.setFilter(pbF)
			return nil
		}
		return errors.New("'Filters' option can only be used with Get or Scan request")
	}
}

// TimeRange is used as a parameter for request creation. Adds TimeRange constraint to a request.
// It will get values in range [from, to[ ('to' is exclusive).
func TimeRange(from, to time.Time) func(Call) error {
	return TimeRangeUint64(uint64(from.UnixNano()/1e6), uint64(to.UnixNano()/1e6))
}

// TimeRangeUint64 is used as a parameter for request creation.
// Adds TimeRange constraint to a request.
// from and to should be in milliseconds
// // It will get values in range [from, to[ ('to' is exclusive).
func TimeRangeUint64(from, to uint64) func(Call) error {
	return func(hc Call) error {
		if c, ok := hc.(hasQueryOptions); ok {
			if from >= to {
				// or equal is becuase 'to' is exclusive
				return errors.New("'from' timestamp is greater or equal to 'to' timestamp")
			}
			c.setTimeRangeUint64(from, to)
			return nil
		}
		return errors.New("'TimeRange' option can only be used with Get or Scan request")
	}
}

// MaxVersions is used as a parameter for request creation.
// Adds MaxVersions constraint to a request.
func MaxVersions(versions uint32) func(Call) error {
	return func(hc Call) error {
		if c, ok := hc.(hasQueryOptions); ok {
			if versions > math.MaxInt32 {
				return errors.New("'MaxVersions' exceeds supported number of versions")
			}
			c.setMaxVersions(versions)
			return nil
		}
		return errors.New("'MaxVersions' option can only be used with Get or Scan request")
	}
}

// MaxResultsPerColumnFamily is an option for Get or Scan requests that sets the maximum
// number of cells returned per column family in a row.
func MaxResultsPerColumnFamily(maxresults uint32) func(Call) error {
	return func(hc Call) error {
		if c, ok := hc.(hasQueryOptions); ok {
			if maxresults > math.MaxInt32 {
				return errors.New(
					"'MaxResultsPerColumnFamily' exceeds supported number of value results")
			}
			c.setMaxResultsPerColumnFamily(maxresults)
			return nil
		}
		return errors.New(
			"'MaxResultsPerColumnFamily' option can only be used with Get or Scan request")
	}
}

// ResultOffset is a option for Scan or Get requests that sets the offset for cells
// within a column family.
func ResultOffset(offset uint32) func(Call) error {
	return func(hc Call) error {
		if c, ok := hc.(hasQueryOptions); ok {
			if offset > math.MaxInt32 {
				return errors.New("'ResultOffset' exceeds supported offset value")
			}
			c.setResultOffset(offset)
			return nil
		}
		return errors.New("'ResultOffset' option can only be used with Get or Scan request")
	}
}
