package goblin

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestGoblinUgcPlayurl(t *testing.T) {
	var (
		ctx = context.Background()
		p   = &model.PlayURLReq{
			Cid: fmt.Sprintf("%d", 10131156),
		}
	)
	convey.Convey("UgcPlayurl", t, func(c convey.C) {
		defer gock.OffAll()
		c.Convey("Normal Situation, Then err should be nil.res,resp should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.UgcPlayURL).Reply(200).JSON(`{
				"result": "succ",
				"message": "succ",
				"code": 0
			}`)
			res, resp, err := d.UgcPlayurl(ctx, p)
			fmt.Println(resp)
			cx.So(err, convey.ShouldBeNil)
			cx.So(resp, convey.ShouldNotBeNil)
			cx.So(res, convey.ShouldNotBeNil)
		})
		c.Convey("Request Error", func(cx convey.C) {
			httpMock("GET", d.conf.Host.UgcPlayURL).Reply(404).JSON(``)
			_, _, err := d.UgcPlayurl(ctx, p)
			cx.So(err, convey.ShouldNotBeNil)
		})
		c.Convey("Code Error", func(cx convey.C) {
			httpMock("GET", d.conf.Host.UgcPlayURL).Reply(200).JSON(`{"code":-400}`)
			_, _, err := d.UgcPlayurl(ctx, p)
			cx.So(err, convey.ShouldNotBeNil)
		})
		c.Convey("Json Error", func(cx convey.C) {
			httpMock("GET", d.conf.Host.UgcPlayURL).Reply(200).JSON(`{"code":-400:}`)
			_, _, err := d.UgcPlayurl(ctx, p)
			cx.So(err, convey.ShouldNotBeNil)
		})
	})
}
