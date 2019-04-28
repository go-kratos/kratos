package goblin

import (
	"context"
	"go-common/app/interface/main/tv/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGoblinVerUpdate(t *testing.T) {
	var (
		ctx = context.Background()
		ver = &model.VerUpdate{}
	)
	convey.Convey("VerUpdate", t, func(c convey.C) {
		c.Convey("Then err should be nil.result,errCode should not be nil.", func(c convey.C) {
			result, errCode, err := d.VerUpdate(ctx, ver)
			httpMock("GET", d.conf.Host.ReqURL).Reply(200).JSON(`{
				"Data": {"ver":123},
				"message": "succ",
				"code": 0
			}`)
			c.So(err, convey.ShouldBeNil)
			c.So(errCode, convey.ShouldBeNil)
			c.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("http error", func(c convey.C) {
			_, _, err := d.VerUpdate(ctx, ver)
			httpMock("GET", d.conf.Host.ReqURL).Reply(404).JSON(``)
			c.So(err, convey.ShouldNotBeNil)
		})
		c.Convey("code error", func(c convey.C) {
			_, _, err := d.VerUpdate(ctx, ver)
			httpMock("GET", d.conf.Host.ReqURL).Reply(200).JSON(`{"code":400}`)
			c.So(err, convey.ShouldNotBeNil)
		})
	})
}
