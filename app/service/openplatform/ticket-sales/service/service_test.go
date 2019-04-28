package service

import (
	"context"
	"flag"
	"testing"

	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/library/conf/paladin"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
	ctx = context.TODO()
)

func init() {
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("ticket-sales.toml", conf.Conf); err != nil {
		panic(err)
	}
	svr = New(conf.Conf)
}

func TestPing(t *testing.T) {
	Convey("TestPing: ", t, func() {
		err := svr.Ping(context.TODO())
		So(err, ShouldBeNil)
	})
}
