package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoElecShow(t *testing.T) {
	convey.Convey("ElecShow", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(2089809)
			aid     = int64(10100652)
			loginID = int64(2089809)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("GET", d.elecShowURL).Reply(200).JSON(`{"code":0,"data":{"show":true,"state":1}}`)
			rs, err := d.ElecShow(c, mid, aid, loginID)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}
