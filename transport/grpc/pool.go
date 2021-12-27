package grpc

import (
	"sync"
)

type pools struct {
	// Transport's pool
	tr sync.Pool
}

var pool *pools

func init() {
	pool = &pools{}
	pool.tr.New = func() interface{} {
		return &Transport{}
	}
}
