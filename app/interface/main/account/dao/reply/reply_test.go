package reply

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

func TestReplyHistoryList(t *testing.T) {
	convey.Convey("ReplyHistoryList", t, func() {
		res, err := d.ReplyHistoryList(context.TODO(), 1, "2017-3-21", "2017-3-21", "time", "desc", 1, 50, "access_key", "cookie", "0.0.0.0")
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
