package lrucache

import (
	"hash/crc32"
	"sync"
	"time"
)

// hashCode hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func hashCode(s string) (hc int) {
	hc = int(crc32.ChecksumIEEE([]byte(s)))
	if hc >= 0 {
		return hc
	}
	if -hc >= 0 {
		return -hc
	}
	// hc == MinInt
	return hc
}

// SyncCache - concurrent cache structure
type SyncCache struct {
	locks   []sync.Mutex
	caches  []*LRUCache
	mask    int
	timeout int64
}

type scValue struct {
	Value interface{}
	ts    int64
}

func nextPowOf2(cap int) int {
	if cap < 2 {
		return 2
	}
	if cap&(cap-1) == 0 {
		return cap
	}
	cap |= cap >> 1
	cap |= cap >> 2
	cap |= cap >> 4
	cap |= cap >> 8
	cap |= cap >> 16
	return cap + 1
}

// NewSyncCache - create sync cache
// `capacity` is lru cache length of each bucket
// store `capacity * bucket` count of element in SyncCache at most
// `timeout` is in seconds
func NewSyncCache(capacity int, bucket int, timeout int64) *SyncCache {
	size := nextPowOf2(bucket)
	sc := SyncCache{make([]sync.Mutex, size), make([]*LRUCache, size), size - 1, timeout}
	for i := range sc.caches {
		sc.caches[i] = New(capacity)
	}
	return &sc
}

// Put - put a cache item into sync cache
func (sc *SyncCache) Put(key string, value interface{}) {
	idx := hashCode(key) & sc.mask
	sc.locks[idx].Lock()
	sc.caches[idx].Put(key, &scValue{value, time.Now().Unix()})
	sc.locks[idx].Unlock()
}

// Get - get value of key from sync cache with result
func (sc *SyncCache) Get(key string) (interface{}, bool) {
	idx := hashCode(key) & sc.mask
	sc.locks[idx].Lock()
	v, b := sc.caches[idx].Get(key)
	if !b {
		sc.locks[idx].Unlock()
		return nil, false
	}
	if time.Now().Unix()-v.(*scValue).ts >= sc.timeout {
		sc.caches[idx].Delete(key)
		sc.locks[idx].Unlock()
		return nil, false
	}
	sc.locks[idx].Unlock()
	return v.(*scValue).Value, b
}

// Delete - delete item by key from sync cache
func (sc *SyncCache) Delete(key string) {
	idx := hashCode(key) & sc.mask
	sc.locks[idx].Lock()
	sc.caches[idx].Delete(key)
	sc.locks[idx].Unlock()
}
