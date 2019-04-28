package passport

import (
	"context"
	"flag"
	"testing"

	"go-common/app/interface/main/account/conf"

	"github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	flag.Parse()

	flag.Set("conf", "../../cmd/account-interface-example.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func TestDao_UserID(t *testing.T) {
	convey.Convey("ReplyHistoryList", t, func() {
		id, err := d.UserID(context.TODO(), 1, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldNotBeEmpty)
	})
}

func TestDao_FastReg(t *testing.T) {
	convey.Convey("ReplyHistoryList", t, func() {
		fastReg, err := d.FastReg(context.TODO(), 111001347, "")
		convey.So(err, convey.ShouldBeNil)
		convey.So(fastReg, convey.ShouldNotBeEmpty)
	})
}
