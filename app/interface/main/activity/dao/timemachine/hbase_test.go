package timemachine

import (
	"context"
	"testing"

	"go-common/app/interface/main/activity/model/timemachine"

	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestTimemachinehbaseRowKey(t *testing.T) {
	convey.Convey("hbaseRowKey", t, func(ctx convey.C) {
		var (
			mid = int64(121212)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := hbaseRowKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTimemachineRawTimemachine(t *testing.T) {
	convey.Convey("RawTimemachine", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			data, err := d.RawTimemachine(c, mid)
			ctx.Convey("Then err should be nil.data should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.Printf("%+v", data)
			})
		})
	})
}

func TestTimemachinetmFillFields(t *testing.T) {
	convey.Convey("tmFillFields", t, func(ctx convey.C) {
		var (
			data = &timemachine.Item{}
			c    = &hrpc.Cell{Qualifier: []byte("is_up"), Value: []byte("1")}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("is_up"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("dh"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("adh"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("a_vv"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("tag_id"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("ls_vv"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("ugc_avs"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("pgc_avs"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("up"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("up_ad"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("up_ld"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("up_avs"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("up_st"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("cir_tm"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("cir_av"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("cir_v"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("fs_av"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("fs_tm"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("fs_ty"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("s_av_rd"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("bt_av"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("bt_ty"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("bt_av_o"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("bt_av_ty"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("o_vv"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("all_vv"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("live_d"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("is_live"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("ld"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("md"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("mc"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("att"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("fan_vv"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("fan_live"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("like_tid"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("like_st"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("wr"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("pd_hr"), Value: []byte("1")}
			tmFillFields(data, c)
			c = &hrpc.Cell{Qualifier: []byte("lu_adr"), Value: []byte("1")}
			tmFillFields(data, c)
		})
	})
}
