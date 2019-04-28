package assist

import (
	"context"
	"net/url"
	"strconv"
	"testing"

	"go-common/app/interface/main/creative/model/assist"

	"github.com/smartystreets/goconvey/convey"
)

func TestAssistLiveStatus(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ip  = ""
		res = struct {
			Code int `json:"code"`
		}{
			Code: 0,
		}
	)
	convey.Convey("LiveStatus", t, func(ctx convey.C) {
		httpMock("GET", d.liveStatusURL+"?uid="+strconv.FormatInt(mid, 10)).Reply(0).JSON(res)
		ok, err := d.LiveStatus(c, mid, ip)
		ctx.Convey("Then err should be nil.ok should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(ok, convey.ShouldNotBeNil)
		})
	})
}

var msgRes = struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}{
	Code:    0,
	Message: "message",
}

func TestAssistLiveAddAssist(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(1)
		assistMid = int64(10)
		cookie    = ""
		ip        = ""
		params    = url.Values{}
	)
	params.Set("admin", strconv.FormatInt(assistMid, 10))
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	convey.Convey("LiveAddAssist", t, func(ctx convey.C) {
		httpMock("POST", d.liveAddAssistURL+"?"+params.Encode()).
			Reply(0).JSON(msgRes)
		err := d.LiveAddAssist(c, mid, assistMid, cookie, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistLiveDelAssist(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		cookie    = ""
		ip        = ""
		params    = url.Values{}
	)
	params.Set("admin", strconv.FormatInt(assistMid, 10))
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	convey.Convey("LiveDelAssist", t, func(ctx convey.C) {
		httpMock("POST", d.liveDelAssistURL+"?"+params.Encode()).Reply(0).JSON(msgRes)
		err := d.LiveDelAssist(c, mid, assistMid, cookie, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistLiveBannedRevoc(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(0)
		banID  = ""
		cookie = ""
		ip     = ""
		params = url.Values{}
	)
	params.Set("uid", strconv.FormatInt(mid, 10))
	params.Set("id", banID)
	convey.Convey("LiveBannedRevoc", t, func(ctx convey.C) {
		httpMock("POST", d.liveRevocBannedURL+"?"+params.Encode()).Reply(200).JSON(msgRes)
		err := d.LiveBannedRevoc(c, mid, banID, cookie, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistLiveAssists(t *testing.T) {
	var (
		c      = context.TODO()
		mid    = int64(0)
		ip     = ""
		params = url.Values{}
		res    = struct {
			Code int                  `json:"code"`
			Data []*assist.LiveAssist `json:"data"`
		}{
			Code: 0,
			Data: []*assist.LiveAssist{
				{
					AssistMid: 0,
					RoomID:    0,
					CTime:     0,
					Datetime:  "",
				},
			},
		}
	)
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	convey.Convey("LiveAssists", t, func(ctx convey.C) {
		httpMock("GET", d.liveAssistsURL+"?"+params.Encode()).Reply(200).JSON(res)
		assists, err := d.LiveAssists(c, mid, ip)
		ctx.Convey("Then err should be nil.assists should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(assists, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistLiveCheckAssist(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		ip        = ""
		params    = url.Values{}
		res       = struct {
			Code int `json:"code"`
		}{
			Code: 0,
		}
	)
	params.Set("uid", strconv.FormatInt(assistMid, 10))
	params.Set("anchor_id", strconv.FormatInt(mid, 10))
	convey.Convey("LiveCheckAssist", t, func(ctx convey.C) {
		httpMock("POST", d.liveCheckAssURL+"?"+params.Encode()).Reply(200).JSON(res)
		isAss, err := d.LiveCheckAssist(c, mid, assistMid, ip)
		ctx.Convey("Then err should be nil.isAss should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(isAss, convey.ShouldNotBeNil)
		})
	})
}
