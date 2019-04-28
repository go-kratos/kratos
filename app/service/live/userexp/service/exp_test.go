package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/userexp/conf"
)

var (
	once sync.Once
	s    *Service
	ctx  = context.TODO()
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	s = New(conf.Conf)
}

func TestGetLevel(t *testing.T) {
	Convey("GetLevel", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)

		level, err := s.Level(ctx, 10001)
		So(err, ShouldBeNil)
		t.Logf("level:%v", level)
	})
}

func TestMultiGetLevel(t *testing.T) {
	Convey("MultiGetLevel", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)

		level, err := s.MultiGetLevel(ctx, []int64{10001, 20002, 30003})
		So(err, ShouldBeNil)
		t.Logf("level:%v", level)
	})
}
