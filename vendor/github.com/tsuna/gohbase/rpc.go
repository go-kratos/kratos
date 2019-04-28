// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gohbase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/region"
	"github.com/tsuna/gohbase/zk"
)

// Constants
var (
	// Name of the meta region.
	metaTableName = []byte("hbase:meta")

	infoFamily = map[string][]string{
		"info": nil,
	}

	// ErrRegionUnavailable is returned when sending rpc to a region that is unavailable
	ErrRegionUnavailable = errors.New("region unavailable")

	// TableNotFound is returned when attempting to access a table that
	// doesn't exist on this cluster.
	TableNotFound = errors.New("table not found")

	// ErrCannotFindRegion is returned when it took too many tries to find a
	// region for the request. It's likely that hbase:meta has overlaps or some other
	// inconsistency.
	ErrConnotFindRegion = errors.New("cannot find region for the rpc")

	// ErrClientClosed is returned when the gohbase client has been closed
	ErrClientClosed = errors.New("client is closed")

	// errMetaLookupThrottled is returned when a lookup for the rpc's region
	// has been throttled.
	errMetaLookupThrottled = errors.New("lookup to hbase:meta has been throttled")
)

const (
	// maxSendRPCTries is the maximum number of times to try to send an RPC
	maxSendRPCTries = 10

	backoffStart = 16 * time.Millisecond
)

func (c *client) SendRPC(rpc hrpc.Call) (proto.Message, error) {
	var err error
	for i := 0; i < maxSendRPCTries; i++ {
		// Check the cache for a region that can handle this request
		reg := c.getRegionFromCache(rpc.Table(), rpc.Key())
		if reg == nil {
			reg, err = c.findRegion(rpc.Context(), rpc.Table(), rpc.Key())
			if err == ErrRegionUnavailable {
				continue
			} else if err == errMetaLookupThrottled {
				// lookup for region has been throttled, check the cache
				// again but don't count this as SendRPC try as there
				// might be just too many request going on at a time.
				i--
				continue
			} else if err != nil {
				return nil, err
			}
		}

		msg, err := c.sendRPCToRegion(rpc, reg)
		switch err {
		case ErrRegionUnavailable:
			if ch := reg.AvailabilityChan(); ch != nil {
				// The region is unavailable. Wait for it to become available,
				// a new region or for the deadline to be exceeded.
				select {
				case <-rpc.Context().Done():
					return nil, rpc.Context().Err()
				case <-c.done:
					return nil, ErrClientClosed
				case <-ch:
				}
			}
		default:
			return msg, err
		}
	}
	return nil, ErrConnotFindRegion
}

func sendBlocking(rc hrpc.RegionClient, rpc hrpc.Call) (hrpc.RPCResult, error) {
	rc.QueueRPC(rpc)

	var res hrpc.RPCResult
	// Wait for the response
	select {
	case res = <-rpc.ResultChan():
		return res, nil
	case <-rpc.Context().Done():
		return res, rpc.Context().Err()
	}
}

func (c *client) sendRPCToRegion(rpc hrpc.Call, reg hrpc.RegionInfo) (proto.Message, error) {
	if reg.IsUnavailable() {
		return nil, ErrRegionUnavailable
	}
	rpc.SetRegion(reg)

	// Queue the RPC to be sent to the region
	client := reg.Client()
	if client == nil {
		// There was an error queueing the RPC.
		// Mark the region as unavailable.
		if reg.MarkUnavailable() {
			// If this was the first goroutine to mark the region as
			// unavailable, start a goroutine to reestablish a connection
			go c.reestablishRegion(reg)
		}
		return nil, ErrRegionUnavailable
	}
	res, err := sendBlocking(client, rpc)
	if err != nil {
		return nil, err
	}
	// Check for errors
	switch res.Error.(type) {
	case region.RetryableError:
		// There's an error specific to this region, but
		// our region client is fine. Mark this region as
		// unavailable (as opposed to all regions sharing
		// the client), and start a goroutine to reestablish
		// it.
		if reg.MarkUnavailable() {
			go c.reestablishRegion(reg)
		}
		return nil, ErrRegionUnavailable
	case region.UnrecoverableError:
		// If it was an unrecoverable error, the region client is
		// considered dead.
		if reg == c.adminRegionInfo {
			// If this is the admin client, mark the region
			// as unavailable and start up a goroutine to
			// reconnect if it wasn't already marked as such.
			if reg.MarkUnavailable() {
				go c.reestablishRegion(reg)
			}
		} else {
			c.clientDown(client)
		}

		// Fall through to the case of the region being unavailable,
		// which will result in blocking until it's available again.
		return nil, ErrRegionUnavailable
	default:
		// RPC was successfully sent, or an unknown type of error
		// occurred. In either case, return the results.
		return res.Msg, res.Error
	}
}

// clientDown removes client from cache and marks
// all the regions sharing this region's
// client as unavailable, and start a goroutine
// to reconnect for each of them.
func (c *client) clientDown(client hrpc.RegionClient) {
	downregions := c.clients.clientDown(client)
	for downreg := range downregions {
		if downreg.MarkUnavailable() {
			downreg.SetClient(nil)
			go c.reestablishRegion(downreg)
		}
	}
}

func (c *client) lookupRegion(ctx context.Context,
	table, key []byte) (hrpc.RegionInfo, string, error) {
	var reg hrpc.RegionInfo
	var addr string
	var err error
	backoff := backoffStart
	for {
		// If it takes longer than regionLookupTimeout, fail so that we can sleep
		lookupCtx, cancel := context.WithTimeout(ctx, c.regionLookupTimeout)
		if c.clientType == adminClient {
			log.WithField("resource", zk.Master).Debug("looking up master")

			addr, err = c.zkLookup(lookupCtx, zk.Master)
			cancel()
			reg = c.adminRegionInfo
		} else if bytes.Compare(table, metaTableName) == 0 {
			log.WithField("resource", zk.Meta).Debug("looking up region server of hbase:meta")

			addr, err = c.zkLookup(lookupCtx, zk.Meta)
			cancel()
			reg = c.metaRegionInfo
		} else {
			log.WithFields(log.Fields{
				"table": strconv.Quote(string(table)),
				"key":   strconv.Quote(string(key)),
			}).Debug("looking up region")

			reg, addr, err = c.metaLookup(lookupCtx, table, key)
			cancel()
			if err == TableNotFound {
				log.WithFields(log.Fields{
					"table": strconv.Quote(string(table)),
					"key":   strconv.Quote(string(key)),
					"err":   err,
				}).Debug("hbase:meta does not know about this table/key")

				return nil, "", err
			} else if err == errMetaLookupThrottled {
				return nil, "", err
			} else if err == ErrClientClosed {
				return nil, "", err
			}
		}
		if err == nil {
			log.WithFields(log.Fields{
				"table":  strconv.Quote(string(table)),
				"key":    strconv.Quote(string(key)),
				"region": reg,
				"addr":   addr,
			}).Debug("looked up a region")

			return reg, addr, nil
		}

		log.WithFields(log.Fields{
			"table":   strconv.Quote(string(table)),
			"key":     strconv.Quote(string(key)),
			"backoff": backoff,
			"err":     err,
		}).Error("failed looking up region")

		// This will be hit if there was an error locating the region
		backoff, err = sleepAndIncreaseBackoff(ctx, backoff)
		if err != nil {
			return nil, "", err
		}
	}
}

func (c *client) findRegion(ctx context.Context, table, key []byte) (hrpc.RegionInfo, error) {
	// The region was not in the cache, it
	// must be looked up in the meta table
	reg, addr, err := c.lookupRegion(ctx, table, key)
	if err != nil {
		return nil, err
	}

	// We are the ones that looked up the region, so we need to
	// mark in unavailable and find a client for it.
	reg.MarkUnavailable()

	if reg != c.metaRegionInfo && reg != c.adminRegionInfo {
		// Check that the region wasn't added to
		// the cache while we were looking it up.
		overlaps, replaced := c.regions.put(reg)
		if !replaced {
			// the same or younger regions are already in cache, retry looking up in cache
			return nil, ErrRegionUnavailable
		}

		// otherwise, new region in cache, delete overlaps from client's cache
		for _, r := range overlaps {
			c.clients.del(r)
		}
	}

	// Start a goroutine to connect to the region
	go c.establishRegion(reg, addr)

	// Wait for the new region to become
	// available, and then send the RPC
	return reg, nil
}

// Searches in the regions cache for the region hosting the given row.
func (c *client) getRegionFromCache(table, key []byte) hrpc.RegionInfo {
	if c.clientType == adminClient {
		return c.adminRegionInfo
	} else if bytes.Equal(table, metaTableName) {
		return c.metaRegionInfo
	}
	regionName := createRegionSearchKey(table, key)
	_, region := c.regions.get(regionName)
	if region == nil {
		return nil
	}

	// make sure the returned region is for the same table
	if !bytes.Equal(fullyQualifiedTable(region), table) {
		// not the same table, can happen if we got the last region
		return nil
	}

	if len(region.StopKey()) != 0 &&
		// If the stop key is an empty byte array, it means this region is the
		// last region for this table and this key ought to be in that region.
		bytes.Compare(key, region.StopKey()) >= 0 {
		return nil
	}

	return region
}

// Creates the META key to search for in order to locate the given key.
func createRegionSearchKey(table, key []byte) []byte {
	metaKey := make([]byte, 0, len(table)+len(key)+3)
	metaKey = append(metaKey, table...)
	metaKey = append(metaKey, ',')
	metaKey = append(metaKey, key...)
	metaKey = append(metaKey, ',')
	// ':' is the first byte greater than '9'.  We always want to find the
	// entry with the greatest timestamp, so by looking right before ':'
	// we'll find it.
	metaKey = append(metaKey, ':')
	return metaKey
}

// lookupLimit throttles lookups to hbase:meta to metaLookupLimit requests
// per metaLookupInterval. It returns nil if we were lucky enough to
// reserve right away and errMetaLookupThrottled or context's error otherwise.
func (c *client) metaLookupLimit(ctx context.Context) error {
	r := c.metaLookupLimiter.Reserve()
	if !r.OK() {
		panic("wtf: cannot reserve a meta lookup")
	}

	delay := r.Delay()
	if delay <= 0 {
		return nil
	}

	// We've been rate limitted
	t := time.NewTimer(delay)
	defer t.Stop()
	select {
	case <-t.C:
		return errMetaLookupThrottled
	case <-ctx.Done():
		r.Cancel()
		return ctx.Err()
	}
}

// metaLookup checks meta table for the region in which the given row key for the given table is.
func (c *client) metaLookup(ctx context.Context,
	table, key []byte) (hrpc.RegionInfo, string, error) {
	metaKey := createRegionSearchKey(table, key)
	rpc, err := hrpc.NewScanRange(ctx, metaTableName, metaKey, table,
		hrpc.Families(infoFamily),
		hrpc.Reversed(),
		hrpc.CloseScanner(),
		hrpc.NumberOfRows(1))
	if err != nil {
		return nil, "", err
	}

	scanner := c.Scan(rpc)
	resp, err := scanner.Next()
	if err == io.EOF {
		return nil, "", TableNotFound
	}
	if err != nil {
		return nil, "", err
	}

	reg, addr, err := region.ParseRegionInfo(resp)
	if err != nil {
		return nil, "", err
	}
	if !bytes.Equal(table, fullyQualifiedTable(reg)) {
		// This would indicate a bug in HBase.
		return nil, "", fmt.Errorf("wtf: meta returned an entry for the wrong table!"+
			"  Looked up table=%q key=%q got region=%s", table, key, reg)
	} else if len(reg.StopKey()) != 0 &&
		bytes.Compare(key, reg.StopKey()) >= 0 {
		// This would indicate a hole in the meta table.
		return nil, "", fmt.Errorf("wtf: meta returned an entry for the wrong region!"+
			"  Looked up table=%q key=%q got region=%s", table, key, reg)
	}
	return reg, addr, nil
}

func fullyQualifiedTable(reg hrpc.RegionInfo) []byte {
	namespace := reg.Namespace()
	table := reg.Table()
	if namespace == nil {
		return table
	}
	// non-default namespace table
	fqTable := make([]byte, 0, len(namespace)+1+len(table))
	fqTable = append(fqTable, namespace...)
	fqTable = append(fqTable, byte(':'))
	fqTable = append(fqTable, table...)
	return fqTable
}

func (c *client) reestablishRegion(reg hrpc.RegionInfo) {
	select {
	case <-c.done:
		return
	default:
	}

	log.WithField("region", reg).Debug("reestablishing region")
	c.establishRegion(reg, "")
}

// probeKey returns a key in region that is unlikely to have data at it
// in order to test if the region is online. This prevents the Get request
// to actually fetch the data from the storage which consumes resources
// of the region server
func probeKey(reg hrpc.RegionInfo) []byte {
	// now we create a probe key: reg.StartKey() + 17 zeros
	probe := make([]byte, len(reg.StartKey())+17)
	copy(probe, reg.StartKey())
	return probe
}

// isRegionEstablished checks whether regionserver accepts rpcs for the region.
// Returns the cause if not established.
func isRegionEstablished(rc hrpc.RegionClient, reg hrpc.RegionInfo) error {
	probe, err := hrpc.NewGet(context.Background(), fullyQualifiedTable(reg), probeKey(reg),
		hrpc.SkipBatch())
	if err != nil {
		panic(fmt.Sprintf("should not happen: %s", err))
	}
	probe.ExistsOnly()

	probe.SetRegion(reg)
	res, err := sendBlocking(rc, probe)
	if err != nil {
		panic(fmt.Sprintf("should not happen: %s", err))
	}

	switch res.Error.(type) {
	case region.RetryableError, region.UnrecoverableError:
		return res.Error
	default:
		return nil
	}
}

func (c *client) establishRegion(reg hrpc.RegionInfo, addr string) {
	var backoff time.Duration
	var err error
	for {
		backoff, err = sleepAndIncreaseBackoff(reg.Context(), backoff)
		if err != nil {
			// region is dead
			reg.MarkAvailable()
			return
		}
		if addr == "" {
			// need to look up region and address of the regionserver
			originalReg := reg
			// lookup region forever until we get it or we learn that it doesn't exist
			reg, addr, err = c.lookupRegion(originalReg.Context(),
				fullyQualifiedTable(originalReg), originalReg.StartKey())

			if err == TableNotFound {
				// region doesn't exist, delete it from caches
				c.regions.del(originalReg)
				c.clients.del(originalReg)
				originalReg.MarkAvailable()

				log.WithFields(log.Fields{
					"region":  originalReg.String(),
					"err":     err,
					"backoff": backoff,
				}).Info("region does not exist anymore")

				return
			} else if originalReg.Context().Err() != nil {
				// region is dead
				originalReg.MarkAvailable()

				log.WithFields(log.Fields{
					"region":  originalReg.String(),
					"err":     err,
					"backoff": backoff,
				}).Info("region became dead while establishing client for it")

				return
			} else if err == errMetaLookupThrottled {
				// We've been throttled, backoff and retry the lookup
				// TODO: backoff might be unnecessary
				reg = originalReg
				continue
			} else if err == ErrClientClosed {
				// client has been closed
				return
			} else if err != nil {
				log.WithFields(log.Fields{
					"region":  originalReg.String(),
					"err":     err,
					"backoff": backoff,
				}).Fatal("unknown error occured when looking up region")
			}
			if !bytes.Equal(reg.Name(), originalReg.Name()) {
				// put new region and remove overlapping ones.
				// Should remove the original region as well.
				reg.MarkUnavailable()
				overlaps, replaced := c.regions.put(reg)
				if !replaced {
					// a region that is the same or younger is already in cache
					reg.MarkAvailable()
					originalReg.MarkAvailable()
					return
				}
				// otherwise delete the overlapped regions in cache
				for _, r := range overlaps {
					c.clients.del(r)
				}
				// let rpcs know that they can retry and either get the newly
				// added region from cache or lookup the one they need
				originalReg.MarkAvailable()
			} else {
				// same region, discard the looked up one
				reg = originalReg
			}
		}

		// connect to the region's regionserver
		client, err := c.establishRegionClient(reg, addr)
		if err == nil {
			if reg == c.adminRegionInfo {
				reg.SetClient(client)
				reg.MarkAvailable()
				return
			}

			if existing := c.clients.put(client, reg); existing != client {
				// a client for this regionserver is already in cache, discard this one.
				client.Close()
				client = existing
			}

			if err = isRegionEstablished(client, reg); err == nil {
				// set region client so that as soon as we mark it available,
				// concurrent readers are able to find the client
				reg.SetClient(client)
				reg.MarkAvailable()
				return
			} else if _, ok := err.(region.UnrecoverableError); ok {
				// the client we got died
				c.clientDown(client)
			}
		} else if err == context.Canceled {
			// region is dead
			reg.MarkAvailable()
			return
		}
		log.WithFields(log.Fields{
			"region":  reg,
			"backoff": backoff,
			"err":     err,
		}).Debug("region was not established, retrying")
		// reset address because we weren't able to connect to it
		// or regionserver says it's still offline, should look up again
		addr = ""
	}
}

func sleepAndIncreaseBackoff(ctx context.Context, backoff time.Duration) (time.Duration, error) {
	if backoff == 0 {
		return backoffStart, nil
	}
	select {
	case <-time.After(backoff):
	case <-ctx.Done():
		return 0, ctx.Err()
	}
	// TODO: Revisit how we back off here.
	if backoff < 5000*time.Millisecond {
		return backoff * 2, nil
	}
	return backoff + 5000*time.Millisecond, nil
}

func (c *client) establishRegionClient(reg hrpc.RegionInfo,
	addr string) (hrpc.RegionClient, error) {
	if c.clientType != adminClient {
		// if rpc is not for hbasemaster, check if client for regionserver
		// already exists
		if client := c.clients.checkForClient(addr); client != nil {
			// There's already a client
			return client, nil
		}
	}

	var clientType region.ClientType
	if c.clientType == standardClient {
		clientType = region.RegionClient
	} else {
		clientType = region.MasterClient
	}
	clientCtx, cancel := context.WithTimeout(reg.Context(), c.regionLookupTimeout)
	defer cancel()

	return region.NewClient(clientCtx, addr, clientType,
		c.rpcQueueSize, c.flushInterval, c.effectiveUser,
		c.regionReadTimeout)
}

// zkResult contains the result of a ZooKeeper lookup (when we're looking for
// the meta region or the HMaster).
type zkResult struct {
	addr string
	err  error
}

// zkLookup asynchronously looks up the meta region or HMaster in ZooKeeper.
func (c *client) zkLookup(ctx context.Context, resource zk.ResourceName) (string, error) {
	// We make this a buffered channel so that if we stop waiting due to a
	// timeout, we won't block the zkLookupSync() that we start in a
	// separate goroutine.
	reschan := make(chan zkResult, 1)
	go func() {
		addr, err := c.zkClient.LocateResource(resource.Prepend(c.zkRoot))
		// This is guaranteed to never block as the channel is always buffered.
		reschan <- zkResult{addr, err}
	}()
	select {
	case res := <-reschan:
		return res.addr, res.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
