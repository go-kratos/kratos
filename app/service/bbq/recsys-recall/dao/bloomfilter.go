package dao

import (
	"context"
	"go-common/library/cache/redis"

	"github.com/Dai0522/go-hash/bloomfilter"
)

// LoadBloomFilter .
func (d *Dao) LoadBloomFilter(ctx *context.Context, key string) (*bloomfilter.BloomFilter, error) {
	conn := d.bfredis.Get(*ctx)
	defer conn.Close()

	var bf *bloomfilter.BloomFilter
	// 获取mid维度
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err != redis.ErrNil && raw != nil {
		bf, err = bloomfilter.Load(&raw)
	}

	return bf, err
}
