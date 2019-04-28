package upper

import (
	"context"
	"testing"

	ugcMdl "go-common/app/job/main/tv/model/ugc"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpperCountUP(t *testing.T) {
	var c = context.Background()
	convey.Convey("CountUP", t, func(ctx convey.C) {
		count, err := d.CountUP(c)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperSendUpper(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("SendUpper", t, func(ctx convey.C) {
		err := d.SendUpper(context.Background(), mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_SetUpperV(t *testing.T) {
	convey.Convey("TestDao_LoadUpMeta", t, WithDao(func(d *Dao) {
		var (
			query = "SELECT mid FROM ugc_uploader WHERE deleted = 0 LIMIT 1"
			upper = &ugcMdl.Upper{}
		)
		if err := d.DB.QueryRow(ctx, query).Scan(&upper.MID); err != nil || upper.MID == 0 {
			fmt.Println("DB Error ", err)
			return
		}
		req := &ugcMdl.ReqSetUp{
			Value:  "test",
			MID:    upper.MID,
			UpType: _upName,
		}
		oriUp, err2 := d.LoadUpMeta(ctx, req.MID)
		if err2 != nil || oriUp == nil {
			return
		}
		fmt.Println(req.MID)
		err := d.setUpperV(context.Background(), req, oriUp)
		convey.So(err, convey.ShouldBeNil)
	}))
}

func TestUpperImportUp(t *testing.T) {
	var (
		up = &ugcMdl.EasyUp{}
	)
	convey.Convey("ImportUp", t, func(ctx convey.C) {
		err := d.ImportUp(context.Background(), up)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
