package assist

import (
	"context"
	"go-common/app/interface/main/creative/model/assist"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAssistAssists(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
		ip  = ""
	)
	convey.Convey("Assists", t, func(ctx convey.C) {
		var res = struct {
			Code int              `json:"code"`
			Data []*assist.Assist `json:"data"`
		}{
			Code: 0,
			Data: []*assist.Assist{
				{
					AssistMid: 1,
					Banned:    1,
				},
			},
		}
		httpMock("GET", d.assistListURL).Reply(200).JSON(res)
		assists, err := d.Assists(c, mid, ip)
		ctx.Convey("Then err should be nil.assists should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(assists, convey.ShouldNotBeNil)
		})
	})
}

var assLogRes = struct {
	Code int               `json:"code"`
	Data *assist.AssistLog `json:"data"`
}{
	Code: 0,
	Data: &assist.AssistLog{
		ID:           1,
		Mid:          0,
		AssistMid:    0,
		AssistAvatar: "",
		AssistName:   "",
		Type:         0,
		Action:       0,
		SubjectID:    0,
		ObjectID:     "",
		Detail:       "",
		State:        0,
		CTime:        0,
	},
}

func TestAssistAssistLog(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		logID     = int64(0)
		ip        = ""
	)
	convey.Convey("AssistLog", t, func(ctx convey.C) {
		httpMock("GET", d.assistLogInfoURL).Reply(200).JSON(assLogRes)
		assistLog, err := d.AssistLog(c, mid, assistMid, logID, ip)
		ctx.Convey("Then err should be nil.assistLog should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(assistLog, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistAssistLogs(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		pn        = int64(0)
		ps        = int64(0)
		stime     = int64(0)
		etime     = int64(0)
		ip        = ""
		res       = struct {
			Code  int                 `json:"code"`
			Data  []*assist.AssistLog `json:"data"`
			Pager map[string]int64    `json:"pager"`
		}{
			Code: 0,
			Data: []*assist.AssistLog{
				{
					ID:           0,
					Mid:          0,
					AssistMid:    0,
					AssistAvatar: "",
					AssistName:   "",
					Type:         0,
					Action:       0,
					SubjectID:    0,
					ObjectID:     "",
					Detail:       "",
					State:        0,
					CTime:        0,
				},
			},
			Pager: map[string]int64{
				"pn": 1,
				"ps": 10,
			},
		}
	)
	convey.Convey("AssistLogs", t, func(ctx convey.C) {
		httpMock("GET", d.assistLogsURL).Reply(200).JSON(res)
		logs, pager, err := d.AssistLogs(c, mid, assistMid, pn, ps, stime, etime, ip)
		ctx.Convey("Then err should be nil.logs,pager should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(pager, convey.ShouldNotBeNil)
			ctx.So(logs, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistAddAssist(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		ip        = "127.0.0.1"
		upUname   = "12"
	)
	convey.Convey("AddAssist", t, func(ctx convey.C) {
		httpMock("POST", d.assistAddURL).Reply(200).JSON(`{"code":0}`)
		err := d.AddAssist(c, mid, assistMid, ip, upUname)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistDelAssist(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		ip        = ""
		upUname   = ""
	)
	convey.Convey("DelAssist", t, func(ctx convey.C) {
		httpMock("POST", d.assistDelURL).Reply(200).JSON(`{"code":0}`)
		err := d.DelAssist(c, mid, assistMid, ip, upUname)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistRevocAssistLog(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		logID     = int64(0)
		ip        = ""
	)
	convey.Convey("RevocAssistLog", t, func(ctx convey.C) {
		httpMock("POST", d.assistLogRevocURL).Reply(200).JSON(`{"code":0}`)
		err := d.RevocAssistLog(c, mid, assistMid, logID, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAssistStat(t *testing.T) {
	var (
		c          = context.TODO()
		mid        = int64(0)
		assistMids = []int64{}
		ip         = ""
		res        = struct {
			Code int                             `json:"code"`
			Data map[int64]map[int8]map[int8]int `json:"data"`
		}{
			Code: 0,
			Data: map[int64]map[int8]map[int8]int{},
		}
	)
	convey.Convey("Stat", t, func(ctx convey.C) {
		httpMock("GET", d.assistStatURL).Reply(200).JSON(res)
		stat, err := d.Stat(c, mid, assistMids, ip)
		ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(stat, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistInfo(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(0)
		assistMid = int64(0)
		ip        = ""
		res       = struct {
			Code int `json:"code"`
			Data struct {
				Assist int8 `json:"assist"`
			} `json:"data"`
		}{
			Code: 0,
			Data: struct {
				Assist int8 `json:"assist"`
			}{
				Assist: 1,
			},
		}
	)
	convey.Convey("Info", t, func(ctx convey.C) {
		httpMock("GET", d.assistInfoURL).Reply(200).JSON(res)
		assist, err := d.Info(c, mid, assistMid, ip)
		ctx.Convey("Then err should be nil.assist should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(assist, convey.ShouldNotBeNil)
		})
	})
}

func TestAssistAssistLogObj(t *testing.T) {
	var (
		c     = context.TODO()
		tp    = int8(0)
		act   = int8(0)
		mid   = int64(0)
		objID = int64(0)
	)
	convey.Convey("AssistLogObj", t, func(ctx convey.C) {
		httpMock("GET", d.assistLogObjURL).Reply(200).JSON(assLogRes)
		assLog, err := d.AssistLogObj(c, tp, act, mid, objID)
		ctx.Convey("Then err should be nil.assLog should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(assLog, convey.ShouldNotBeNil)
		})
	})
}
