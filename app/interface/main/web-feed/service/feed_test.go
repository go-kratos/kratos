package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/web-feed/conf"
	"go-common/app/interface/main/web-feed/model"

	. "github.com/smartystreets/goconvey/convey"
)

func initConf(t *testing.T) {
	if err := conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
}

func TestFeed(t *testing.T) {
	Convey("feed", t, func() {
		var (
			mid int64 = 27515256
			pn        = 1
			ps        = 20
			c         = context.TODO()
			res []*model.Feed
			err error
		)
		initConf(t)
		svr := New(conf.Conf)
		if res, err = svr.Feed(c, mid, pn, ps); err != nil {
			t.Error(err)
		}
		t.Logf("result length:%d", len(res))
	})
}

func TestUnreadCount(t *testing.T) {
	var (
		mid   int64 = 27515256
		c           = context.TODO()
		count int
		err   error
	)
	initConf(t)
	svr := New(conf.Conf)
	if count, err = svr.UnreadCount(c, mid); err != nil {
		t.Error(err)
	}
	t.Logf("count:%d", count)
}
