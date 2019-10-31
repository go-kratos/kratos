package testdata

import (
	"context"

	"github.com/bilibili/kratos/pkg/sync/pipeline/fanout"
)

// Demo test struct
type Demo struct {
	ID    int64
	Title string
}

// Dao .
type dao struct {
	cache *fanout.Fanout
}

// New .
func New() *dao {
	return &dao{cache: fanout.New("cache")}
}

//go:generate kratos tool genbts
type _bts interface {
	// bts: -batch=2 -max_group=20 -batch_err=break -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	Demos(c context.Context, keys []int64) (map[int64]*Demo, error)
	// bts: -batch=2 -max_group=20 -batch_err=continue -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	Demos1(c context.Context, keys []int64) (map[int64]*Demo, error)
	// bts: -sync=true -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	Demo(c context.Context, key int64) (*Demo, error)
	// bts: -paging=true
	Demo1(c context.Context, key int64, pn int, ps int) (*Demo, error)
	// bts: -nullcache=&Demo{ID:-1} -check_null_code=$.ID==-1
	None(c context.Context) (*Demo, error)
}
