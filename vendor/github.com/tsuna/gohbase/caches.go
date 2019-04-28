// Copyright (C) 2016  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gohbase

import (
	"bytes"
	"io"
	"sync"

	"github.com/cznic/b"
	log "github.com/sirupsen/logrus"
	"github.com/tsuna/gohbase/hrpc"
)

// clientRegionCache is client -> region cache. Used to quickly
// look up all the regioninfos that map to a specific client
type clientRegionCache struct {
	m sync.RWMutex

	regions map[hrpc.RegionClient]map[hrpc.RegionInfo]struct{}
}

// put caches client and associates a region with it. Returns a client that is in cache.
// TODO: obvious place for optimization (use map with address as key to lookup exisiting clients)
func (rcc *clientRegionCache) put(c hrpc.RegionClient, r hrpc.RegionInfo) hrpc.RegionClient {
	rcc.m.Lock()
	for existingClient, regions := range rcc.regions {
		// check if client already exists, checking by host and port
		// because concurrent callers might try to put the same client
		if c.Addr() == existingClient.Addr() {
			// check client already knows about the region, checking
			// by pointer is enough because we make sure that there are
			// no regions with the same name around
			if _, ok := regions[r]; !ok {
				regions[r] = struct{}{}
			}
			rcc.m.Unlock()

			log.WithFields(log.Fields{
				"existingClient": existingClient,
				"client":         c,
			}).Debug("region client is already in client's cache")
			return existingClient
		}
	}

	// no such client yet
	rcc.regions[c] = map[hrpc.RegionInfo]struct{}{r: struct{}{}}
	rcc.m.Unlock()

	log.WithField("client", c).Info("added new region client")
	return c
}

func (rcc *clientRegionCache) del(r hrpc.RegionInfo) {
	rcc.m.Lock()
	c := r.Client()
	if c != nil {
		r.SetClient(nil)
		regions := rcc.regions[c]
		delete(regions, r)
	}
	rcc.m.Unlock()
}

func (rcc *clientRegionCache) closeAll() {
	rcc.m.Lock()
	for client, regions := range rcc.regions {
		for region := range regions {
			region.MarkUnavailable()
			region.SetClient(nil)
		}
		client.Close()
	}
	rcc.m.Unlock()
}

func (rcc *clientRegionCache) clientDown(c hrpc.RegionClient) map[hrpc.RegionInfo]struct{} {
	rcc.m.Lock()
	downregions, ok := rcc.regions[c]
	delete(rcc.regions, c)
	rcc.m.Unlock()

	if ok {
		log.WithField("client", c).Info("removed region client")
	}
	return downregions
}

// TODO: obvious place for optimization (use map with address as key to lookup exisiting clients)
func (rcc *clientRegionCache) checkForClient(addr string) hrpc.RegionClient {
	rcc.m.RLock()

	for client := range rcc.regions {
		if client.Addr() == addr {
			rcc.m.RUnlock()
			return client
		}
	}

	rcc.m.RUnlock()
	return nil
}

// key -> region cache.
type keyRegionCache struct {
	m sync.RWMutex

	// Maps a []byte of a region start key to a hrpc.RegionInfo
	regions *b.Tree
}

func (krc *keyRegionCache) get(key []byte) ([]byte, hrpc.RegionInfo) {
	krc.m.RLock()

	enum, ok := krc.regions.Seek(key)
	if ok {
		krc.m.RUnlock()
		log.Fatalf("WTF: got exact match for region search key %q", key)
		return nil, nil
	}
	k, v, err := enum.Prev()
	enum.Close()

	krc.m.RUnlock()

	if err == io.EOF {
		// we are the beginning of the tree
		return nil, nil
	}
	return k.([]byte), v.(hrpc.RegionInfo)
}

func isRegionOverlap(regA, regB hrpc.RegionInfo) bool {
	// if region's stop key is empty, it's assumed to be the greatest key
	return bytes.Equal(regA.Namespace(), regB.Namespace()) &&
		bytes.Equal(regA.Table(), regB.Table()) &&
		(len(regB.StopKey()) == 0 || bytes.Compare(regA.StartKey(), regB.StopKey()) < 0) &&
		(len(regA.StopKey()) == 0 || bytes.Compare(regA.StopKey(), regB.StartKey()) > 0)
}

func (krc *keyRegionCache) getOverlaps(reg hrpc.RegionInfo) []hrpc.RegionInfo {
	var overlaps []hrpc.RegionInfo
	var v interface{}
	var err error

	// deal with empty tree in the beginning so that we don't have to check
	// EOF errors for enum later
	if krc.regions.Len() == 0 {
		return overlaps
	}

	// check if key created from new region falls into any cached regions
	key := createRegionSearchKey(fullyQualifiedTable(reg), reg.StartKey())
	enum, ok := krc.regions.Seek(key)
	if ok {
		log.Fatalf("WTF: found a region with exact name as the search key %q", key)
	}

	// case 1: landed before the first region in cache
	// enum.Prev() returns io.EOF
	// enum.Next() returns io.EOF
	// SeekFirst() + enum.Next() returns the first region, which has larger start key

	// case 2: landed before the second region in cache
	// enum.Prev() returns the first region X and moves pointer to -infinity
	// enum.Next() returns io.EOF
	// SeekFirst() + enum.Next() returns first region X, which has smaller start key

	// case 3: landed anywhere after the second region
	// enum.Prev() returns the region X before it landed, moves pointer to the region X - 1
	// enum.Next() returns X - 1 and move pointer to X, which has smaller start key

	enum.Prev()
	_, _, err = enum.Next()
	if err == io.EOF {
		// we are in the beginning of tree, get new enum starting
		// from first region
		enum.Close()
		enum, err = krc.regions.SeekFirst()
		if err != nil {
			log.Fatalf(
				"error seeking first region when getting  overlaps for region %v: %v", reg, err)
		}
	}

	_, v, err = enum.Next()
	if isRegionOverlap(v.(hrpc.RegionInfo), reg) {
		overlaps = append(overlaps, v.(hrpc.RegionInfo))
	}
	_, v, err = enum.Next()

	// now append all regions that overlap until the end of the tree
	// or until they don't overlap
	for err != io.EOF && isRegionOverlap(v.(hrpc.RegionInfo), reg) {
		overlaps = append(overlaps, v.(hrpc.RegionInfo))
		_, v, err = enum.Next()
	}
	enum.Close()
	return overlaps
}

// put looks up if there's already region with this name in regions cache
// and if there's, returns it in overlaps and doesn't modify the cache.
// Otherwise, it puts the region and removes all overlaps in case all of
// them are older. Returns a slice of overlapping regions and whether
// passed region was put in the cache.
func (krc *keyRegionCache) put(reg hrpc.RegionInfo) (overlaps []hrpc.RegionInfo, replaced bool) {
	krc.m.Lock()
	krc.regions.Put(reg.Name(), func(v interface{}, exists bool) (interface{}, bool) {
		if exists {
			// region is already in cache,
			// note: regions with the same name have the same age
			overlaps = []hrpc.RegionInfo{v.(hrpc.RegionInfo)}
			return nil, false
		}
		// find all entries that are overlapping with the range of the new region.
		overlaps = krc.getOverlaps(reg)
		for _, o := range overlaps {
			if o.ID() > reg.ID() {
				// overlapping region is younger,
				// don't replace any regions
				// TODO: figure out if there can a case where we might
				// have both older and younger overlapping regions, for
				// now we only replace if all overlaps are older
				return nil, false
			}
		}
		// all overlaps are older, put the new region
		replaced = true
		return reg, true
	})
	if !replaced {
		krc.m.Unlock()

		log.WithFields(log.Fields{
			"region":   reg,
			"overlaps": overlaps,
			"replaced": replaced,
		}).Debug("region is already in cache")
		return
	}
	// delete overlapping regions
	// TODO: in case overlaps are always either younger or older,
	// we can just greedily remove them in Put function
	for _, o := range overlaps {
		krc.regions.Delete(o.Name())
		// let region establishers know that they can give up
		o.MarkDead()
	}
	krc.m.Unlock()

	log.WithFields(log.Fields{
		"region":   reg,
		"overlaps": overlaps,
		"replaced": replaced,
	}).Info("added new region")
	return
}

func (krc *keyRegionCache) del(reg hrpc.RegionInfo) bool {
	krc.m.Lock()
	success := krc.regions.Delete(reg.Name())
	krc.m.Unlock()
	// let region establishers know that they can give up
	reg.MarkDead()

	log.WithFields(log.Fields{
		"region": reg,
	}).Debug("removed region")
	return success
}
