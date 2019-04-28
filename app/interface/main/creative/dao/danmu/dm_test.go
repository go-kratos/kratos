package danmu

import (
	"context"
	"encoding/json"
	"go-common/app/interface/main/creative/model/danmu"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDanmuList(t *testing.T) {
	var (
		c      = context.TODO()
		cid    = int64(0)
		mid    = int64(2089809)
		page   = int(1)
		size   = int(10)
		order  = "ctime"
		pool   = "0"
		midStr = ""
		ip     = "127.0.0.1"
	)
	convey.Convey("List", t, func(ctx convey.C) {
		dmList, err := d.List(c, cid, mid, page, size, order, pool, midStr, ip)
		ctx.Convey("Then err should be nil.dmList should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(dmList, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuEdit(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		cid   = int64(0)
		state = int8(0)
		dmids = []int64{}
		ip    = "127.0.0.1"
	)
	convey.Convey("Edit", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.dmEditURL).Reply(200).JSON(`{"code":0,"data":""}`)
		err := d.Edit(c, mid, cid, state, dmids, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDanmuTransfer(t *testing.T) {
	var (
		c       = context.TODO()
		mid     = int64(2089809)
		fromCID = int64(1)
		toCID   = int64(2)
		offset  = float64(10.0)
		ak      = "ak"
		ck      = "ck"
		ip      = "127.0.0.1"
	)
	convey.Convey("Transfer", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.dmTransferURL).Reply(200).JSON(`{"code":0,"data":""}`)
		err := d.Transfer(c, mid, fromCID, toCID, offset, ak, ck, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDanmuUpPool(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		cid   = int64(0)
		dmids = []int64{}
		pool  = int8(0)
	)
	convey.Convey("UpPool", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.dmPoolURL).Reply(200).JSON(`{"code":0,"data":""}`)
		err := d.UpPool(c, mid, cid, dmids, pool)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDanmuDistri(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		cid = int64(1)
		ip  = "127.0.0.1"
	)
	convey.Convey("Distri", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Get", d.dmDistriURL).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":{"1":1}}`)
		distri, err := d.Distri(c, mid, cid, ip)
		ctx.Convey("Then err should be nil.distri should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(distri, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuRecent(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		pn  = int64(1)
		ps  = int64(10)
		ip  = "127.0.0.1"
	)
	convey.Convey("Recent", t, func(ctx convey.C) {
		var res struct {
			Code         int                 `json:"code"`
			ResNewRecent *danmu.ResNewRecent `json:"data"`
		}
		res.ResNewRecent = &danmu.ResNewRecent{
			Page: &danmu.RecentPage{
				Pn:    1,
				Ps:    10,
				Total: 20,
			},
		}
		res.ResNewRecent.Result = append(res.ResNewRecent.Result, &danmu.DMMember{
			ID:  1,
			Aid: 99,
		})
		defer gock.OffAll()
		js, _ := json.Marshal(res)
		httpMock("Get", d.dmRecentURL).Reply(200).JSON(string(js))
		dmRecent, aids, err := d.Recent(c, mid, pn, ps, ip)
		ctx.Convey("Recent", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(aids, convey.ShouldNotBeNil)
			ctx.So(dmRecent, convey.ShouldNotBeNil)
		})
	})
}

func TestDanmuisProtect(t *testing.T) {
	var (
		attrs = ""
		num   = int64(0)
	)
	convey.Convey("isProtect", t, func(ctx convey.C) {
		p1 := d.isProtect(attrs, num)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
