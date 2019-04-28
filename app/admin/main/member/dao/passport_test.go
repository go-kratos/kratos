package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoUpdateUname(t *testing.T) {
	convey.Convey("UpdateUname", t, func(ctx convey.C) {
		var (
			mid  = int64(321)
			name = fmt.Sprintf("321testName%v%v", time.Now().Minute(), time.Now().Second())
		)
		ctx.Convey("UpdateUname success", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.upUnameURL).Reply(200).JSON(`{"code":0}`)
			err := d.UpdateUname(context.Background(), mid, name)
			ctx.So(err, convey.ShouldBeNil)
		})

		ctx.Convey("UpdateUname failed", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.upUnameURL).Reply(200).JSON(`{"code":500}`)
			err := d.UpdateUname(context.Background(), mid, name)
			ctx.So(err, convey.ShouldNotBeNil)
		})

	})
}

func TestDaoPassportQueryByMids(t *testing.T) {
	convey.Convey("PassportQueryByMids", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p, err := d.PassportQueryByMids(context.Background(), []int64{1, 2, 3})
			ctx.Convey("Then err should be nil.p should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPassportQueryByMidsChunked(t *testing.T) {
	convey.Convey("PassportQueryByMidsChunked", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p, err := d.PassportQueryByMidsChunked(context.Background(), []int64{1, 2, 3}, 50)
			ctx.Convey("Then err should be nil.p should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(p, convey.ShouldNotBeNil)
			})
		})
	})
}
