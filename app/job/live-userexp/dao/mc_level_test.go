package dao

import (
	"context"
	"sync"
	"testing"
	"time"

	"go-common/app/job/live-userexp/conf"
	"go-common/app/job/live-userexp/model"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	d    *Dao
	ctx  = context.TODO()
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
}

func startService() {
	initConf()
	d = New(conf.Conf)
	time.Sleep(time.Second * 2)
}

func TestSetLevelCache(t *testing.T) {
	Convey("SetLevelCache", t, func() {
		once.Do(startService)
		err := d.SetLevelCache(ctx, &model.Level{Uid: 10001, Uexp: 1000, Rexp: 100, Ulevel: 2, Rlevel: 1, Color: 12345})
		So(err, ShouldBeNil)
	})
}
