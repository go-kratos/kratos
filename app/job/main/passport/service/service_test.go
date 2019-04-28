package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go-common/app/job/main/passport/conf"
	idfgmdl "go-common/app/service/main/identify-game/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	// init log
	log.Init(conf.Conf.Xlog)
	s = New(conf.Conf)
}

func TestNew(t *testing.T) {
	once.Do(startService)
	Convey("new", t, func() {
		So(s.gameAppIDs[0], ShouldEqual, _gameAppID)
		t.Logf("s.gameAppIDs: %v", s.gameAppIDs)

		So(s.c.URI, ShouldNotBeNil)
		So(s.c.URI.DelCache, ShouldNotBeEmpty)
		So(s.c.URI.SetToken, ShouldNotBeEmpty)
		t.Logf("s.c.URI: %+v", s.c.URI)
	})
}

func TestDelCache(t *testing.T) {
	once.Do(startService)
	time.Sleep(time.Second * 1)
	Convey("del cache", t, func() {
		arg := &idfgmdl.CleanCacheArgs{
			Token: "foo",
		}
		err := s.igRPC.DelCache(context.Background(), arg)
		So(err, ShouldBeNil)
	})
}
