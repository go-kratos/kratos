package bnj

import "context"

func timeFinishKey() string {
	return "time_finish"
}

func lessTimeKey() string {
	return "time_less"
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// mc: -key=timeFinishKey
	CacheTimeFinish(c context.Context) (int64, error)
	// mc: -key=timeFinishKey -expire=d.timeFinishExpire -encode=raw
	AddCacheTimeFinish(c context.Context, value int64) error
	// mc: -key=lessTimeKey
	CacheLessTime(c context.Context) (int64, error)
	// mc: -key=lessTimeKey -expire=d.lessTimeExpire -encode=raw
	AddCacheLessTime(c context.Context, value int64) error
}
