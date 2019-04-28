package realname

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

func TestTelInfo(t *testing.T) {
	convey.Convey("ReplyHistoryList", t, func() {
		res, err := d.TelInfo(context.TODO(), 46333)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
	})
}

func TestAntispam(t *testing.T) {
	convey.Convey("antispam", t, func() {
		anti := new(AlipayAntispamValue)
		anti.IncreaseCount()
		c := anti.Count()
		convey.So(c, convey.ShouldEqual, 1)
		p := anti.Pass()
		convey.So(p, convey.ShouldBeFalse)

		anti.IncreaseCount()
		c = anti.Count()
		convey.So(c, convey.ShouldEqual, 2)

		anti.SetPass(true)
		p = anti.Pass()
		convey.So(p, convey.ShouldBeTrue)
	})
}
