package cms

import (
	"context"
	"fmt"
	"testing"

	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsValidSns(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("ValidSns", t, func(cx convey.C) {
		res, err := d.ValidSns(ctx, false)
		cx.Convey("Consider Audited, Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
			fmt.Println(len(res))
		})
		res, err = d.ValidSns(ctx, true)
		cx.Convey("Consider Audited and Free, Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
			fmt.Println(len(res))
		})
	})
}

func TestCmsShelveOp(t *testing.T) {
	var (
		ctx         = context.Background()
		validSns, _ = d.ValidSns(ctx, true)
	)
	convey.Convey("ShelveOp", t, func(cx convey.C) {
		onIDs, offIDs, err := d.ShelveOp(ctx, validSns)
		cx.Convey("Then err should be nil.onIDs,offIDs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(offIDs, convey.ShouldNotBeNil)
			ctx.So(onIDs, convey.ShouldNotBeNil)
			fmt.Println(offIDs)
			fmt.Println(onIDs)
		})
	})
}

func TestCmsActOps(t *testing.T) {
	var (
		ctx = context.Background()
		sid int64
	)
	d.DB.QueryRow(ctx, "SELECT id FROM tv_ep_season WHERE valid = 1 and is_deleted = 0 AND `check` = 1 LIMIT 1").Scan(&sid)
	convey.Convey("ActOps", t, func(cx convey.C) {
		err := d.ActOps(ctx, []int64{sid}, false)
		cx.Convey("Action 0 Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		err = d.ActOps(ctx, []int64{sid}, true)
		cx.Convey("Action 1 Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		fmt.Println(sid)
	})
}

func TestCmsOffArcs(t *testing.T) {
	var (
		ctx = context.Background()
		aid int64
	)
	convey.Convey("OffArcs", t, func(cx convey.C) {
		cx.Convey("Then err should be nil.offAids should not be nil.", func(cx convey.C) {
			if err := d.DB.QueryRow(ctx, "select aid from ugc_archive where deleted = 0 and valid = 0 and result = 1 limit 1").Scan(&aid); err != nil {
				offAids, err := d.OffArcs(context.Background(), []int64{1, 2, 3})
				cx.So(err, convey.ShouldBeNil)
				cx.So(offAids, convey.ShouldBeNil)
			} else {
				fmt.Println("Have Aid ", aid)
				offAids, err := d.OffArcs(context.Background(), []int64{1, 2, 3, aid})
				cx.So(err, convey.ShouldBeNil)
				cx.So(offAids, convey.ShouldNotBeNil)
			}
		})
		cx.Convey("Arg Error", func(cx convey.C) {
			_, err := d.OffArcs(context.Background(), []int64{})
			cx.So(err, convey.ShouldNotBeNil)
		})
		cx.Convey("DB Error", func(cx convey.C) {
			d.DB.Close()
			_, err := d.OffArcs(context.Background(), []int64{1, 2, 3})
			cx.So(err, convey.ShouldNotBeNil)
			d.DB = sql.NewMySQL(d.conf.Mysql)
		})
	})
}

func TestCmsReshelfArcs(t *testing.T) {
	convey.Convey("OffArcs", t, func(cx convey.C) {
		cx.Convey("Then err should be nil.offAids should not be nil.", func(cx convey.C) {
			err := d.ReshelfArcs(context.Background(), []int64{1, 2, 3})
			cx.So(err, convey.ShouldBeNil)
		})
		cx.Convey("Arg Error", func(cx convey.C) {
			err := d.ReshelfArcs(context.Background(), []int64{})
			cx.So(err, convey.ShouldNotBeNil)
		})
		cx.Convey("DB Error", func(cx convey.C) {
			d.DB.Close()
			err := d.ReshelfArcs(context.Background(), []int64{1, 2, 3})
			cx.So(err, convey.ShouldNotBeNil)
			d.DB = sql.NewMySQL(d.conf.Mysql)
		})
	})
}
