package dao

import (
	"context"
	"sync"
	"testing"
	"time"

	"go-common/app/service/live/userexp/conf"
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

func TestInitExp(t *testing.T) {
	Convey("Init Exp", t, func() {
		once.Do(startService)
		_, err := d.InitExp(ctx, 10001, 0, 0)
		So(err, ShouldBeNil)
	})
}

func TestExp(t *testing.T) {
	Convey("Init Exp", t, func() {
		once.Do(startService)
		rs, err := d.Exp(ctx, 10001)
		So(err, ShouldBeNil)
		t.Logf("QueryExp %v", rs)
	})
}

func TestMultiExp(t *testing.T) {
	Convey("Multi Exp", t, func() {
		once.Do(startService)
		rs, err := d.MultiExp(ctx, []int64{10001, 10002})
		So(err, ShouldBeNil)
		t.Logf("QueryExp rs=%v", rs)
	})
}

func TestAddUexp(t *testing.T) {
	Convey("Add Uexp", t, func() {
		once.Do(startService)
		_, err := d.AddUexp(ctx, 10001, 111)
		So(err, ShouldBeNil)
	})
}

func TestAddRexp(t *testing.T) {
	Convey("Add Uexp", t, func() {
		once.Do(startService)
		_, err := d.AddRexp(ctx, 11111, 111)
		So(err, ShouldBeNil)
	})
}
