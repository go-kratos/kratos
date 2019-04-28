package history

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	hmdl "go-common/app/interface/main/history/model"
	"go-common/app/interface/main/tv/model/history"

	"github.com/smartystreets/goconvey/convey"
)

var (
	testMid   = int64(320773689)
	testAid   = int64(10110670)
	testHisMC = &history.HisMC{
		MID: testMid,
		Res: []*history.HisRes{
			{
				Mid: testMid,
				Oid: testAid,
			},
		},
		LastViewAt: time.Now().Unix(),
	}
)

func TestDao_Cursor(t *testing.T) {
	convey.Convey("TestDao_Cursor", t, WithDao(func(d *Dao) {
		var (
			mcDataHis = &hmdl.ArgHistory{
				Mid:      testMid,
				Realtime: time.Now().Unix(),
				History: &hmdl.History{
					Mid:      testMid,
					Aid:      testAid,
					Business: "archive",
				},
			}
		)
		res, err := d.Cursor(ctx, testMid, 0, 100, 0, []string{"pgc", "archive"})
		if len(res) == 0 {
			if errAdd := d.hisRPC.Add(ctx, mcDataHis); errAdd != nil {
				fmt.Println(errAdd)
				return
			}
			res, err = d.Cursor(ctx, testMid, 0, 100, 0, []string{"pgc", "archive"})
		}
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(res), convey.ShouldBeGreaterThan, 0)
		data, _ := json.Marshal(res)
		fmt.Println(string(data))
	}))
}

func TestHistorykeyHis(t *testing.T) {
	var (
		mid = int64(320773689)
	)
	convey.Convey("keyHis", t, func(c convey.C) {
		p1 := keyHis(mid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestHistorySaveHisCache(t *testing.T) {
	var (
		filtered = []*history.HisRes{
			{
				Mid: testMid,
				Oid: testAid,
			},
		}
	)
	convey.Convey("SaveHisCache", t, func(c convey.C) {
		d.SaveHisCache(ctx, filtered)
		c.Convey("No return values", func(c convey.C) {
		})
	})
}

func TestHistoryaddHisCache(t *testing.T) {
	convey.Convey("addHisCache", t, func(c convey.C) {
		d.addHisCache(ctx, testHisMC)
		c.Convey("No return values", func(c convey.C) {
		})
	})
}

func TestHistorysetHisCache(t *testing.T) {
	convey.Convey("setHisCache", t, func(c convey.C) {
		err := d.setHisCache(ctx, testHisMC)
		c.Convey("Then err should be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
		})
	})
}

func TestHistoryHisCache(t *testing.T) {
	convey.Convey("HisCache", t, func(c convey.C) {
		s, err := d.HisCache(ctx, testMid)
		c.Convey("Then err should be nil.s should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(s, convey.ShouldNotBeNil)
		})
	})
}
