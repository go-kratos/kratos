package testdata

import (
	"context"

	"go-common/library/sync/pipeline/fanout"
)

// Article test struct
type Article struct {
	ID    int64
	Title string
}

// Dao .
type Dao struct {
	cache *fanout.Fanout
}

// New .
func New() *Dao {
	return &Dao{cache: fanout.New("cache")}
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -batch=10 -max_group=10 -batch_err=break -nullcache=&Article{ID:-1} -check_null_code=$.ID==-1
	Articles(c context.Context, keys []int64) (map[int64]*Article, error)
	// cache: -sync=true -nullcache=&Article{ID:-1} -check_null_code=$.ID==-1
	Article(c context.Context, key int64) (*Article, error)
	// cache: -paging=true
	Article1(c context.Context, key int64, pn, ps int) (*Article, error)
	// cache: -nullcache=&Article{ID:-1} -check_null_code=$.ID==-1
	None(c context.Context) (*Article, error)
}
