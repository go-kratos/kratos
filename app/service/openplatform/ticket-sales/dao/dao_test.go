package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/service/openplatform/ticket-sales/conf"

	"go-common/library/conf/paladin"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d   *Dao
	ctx context.Context
)

func init() {
	dir, _ := filepath.Abs("../cmd/ticket-sales.toml")
	flag.Set("conf", dir)
	flag.Set("appid", "open.ticket.sales")
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("ticket-sales.toml", conf.Conf); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	ctx = context.TODO()
}

func TestDatabusPub(t *testing.T) {
	convey.Convey("DatabusPub", t, func() {
		err := d.DatabusPub(context.TODO(), "aaa", []string{"masaji"})
		convey.So(err, convey.ShouldEqual, nil)
	})
}
