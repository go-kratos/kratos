package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_PingMc(t *testing.T) {
	Convey("ping mc", t, WithDao(func(d *Dao) {
		err := d.pingMC(context.Background())
		So(err, ShouldBeNil)
	}))
}

func TestAddReportsCacheByMids(t *testing.T) {
	Convey("add reports cache by mids", t, WithDao(func(d *Dao) {
		var err error
		mrs := map[int64][]*model.Report{
			910819: {{
				APPID:       1,
				PlatformID:  1,
				Mid:         910819,
				DeviceToken: "dt1",
			}, {
				APPID:       2,
				PlatformID:  2,
				Mid:         910819,
				DeviceToken: "dt2",
			}},
			123456: {{
				APPID:       3,
				PlatformID:  3,
				Mid:         123456,
				DeviceToken: "dt3",
			}},
		}
		err = d.AddReportsCacheByMids(context.Background(), mrs)
		So(err, ShouldBeNil)
	}))
}

func Test_ReportCache(t *testing.T) {
	Convey("reports cache", t, WithDao(func(d *Dao) {
		var err error
		mrs := map[int64][]*model.Report{
			910819: {{
				APPID:       1,
				PlatformID:  1,
				Mid:         910819,
				DeviceToken: "dt1",
			}, {
				APPID:       2,
				PlatformID:  2,
				Mid:         910819,
				DeviceToken: "dt2",
			}},
			123456: {{
				APPID:       3,
				PlatformID:  3,
				Mid:         123456,
				DeviceToken: "dt3",
			}},
		}
		err = d.AddReportsCacheByMids(context.Background(), mrs)
		So(err, ShouldBeNil)

		// add report
		// err = d.AddReportCache(context.Background(), &model.Report{APPID: 3, PlatformID: 3, Mid: 123456, DeviceToken: "dt4"})
		// So(err, ShouldBeNil)
		// err = d.AddReportCache(context.Background(), &model.Report{APPID: 4, PlatformID: 4, Mid: 123456, DeviceToken: "dt5"})
		// So(err, ShouldBeNil)

		// delete report
		err = d.DelReportCache(context.Background(), 910819, 2, "dt2")
		So(err, ShouldBeNil)

		// get report
		rs, missed, err := d.ReportsCacheByMids(context.Background(), []int64{910819, 123456})
		_ = missed
		So(len(rs), ShouldEqual, 2)
		So(err, ShouldBeNil)
		for mid, v := range rs {
			for _, vv := range v {
				fmt.Printf("mid(%d) %+v \n", mid, vv)
			}
		}

		// report miss
		rs, misses, err := d.ReportsCacheByMids(context.Background(), []int64{1000000, 2000000})
		So(len(rs), ShouldEqual, 0)
		So(len(misses), ShouldEqual, 2)
		So(err, ShouldBeNil)
	}))
}

func Test_TokenCache(t *testing.T) {
	Convey("add token cache", t, WithDao(func(d *Dao) {
		token := "testtoken"
		r := &model.Report{
			APPID:       1,
			DeviceToken: token,
		}
		err := d.AddTokenCache(context.Background(), r.DeviceToken, r)
		So(err, ShouldBeNil)
		m := make(map[string]*model.Report, 0)
		m[r.DeviceToken] = r
		d.AddTokensCache(context.Background(), m)
		So(err, ShouldBeNil)
		Convey("token cache", func() {
			r, err := d.TokenCache(context.Background(), token)
			So(err, ShouldBeNil)
			t.Logf("report(%+v)", r)

			Convey("delete token cache", func() {
				err = d.DelTokenCache(context.Background(), token)
				So(err, ShouldBeNil)
			})
		})
	}))
}

func Test_TokensCache(t *testing.T) {
	Convey("tokens cache", t, WithDao(func(d *Dao) {
		r := &model.Report{APPID: 1, DeviceToken: "testtoken1"}
		err := d.AddTokenCache(context.Background(), r.DeviceToken, r)
		So(err, ShouldBeNil)

		r = &model.Report{APPID: 1, DeviceToken: "testtoken2"}
		err = d.AddTokenCache(context.Background(), r.DeviceToken, r)
		So(err, ShouldBeNil)

		res, missed, err := d.TokensCache(context.Background(), []string{"testtoken1", "testtoken2", "testtoken3"})
		So(err, ShouldBeNil)
		t.Logf("tokens cache missed(%v)", missed)
		for token, val := range res {
			t.Logf("token(%s) value(%+v)", token, val)
		}
	}))
}

func Test_ReportsCacheByMid(t *testing.T) {
	Convey("Test_ReportsCacheByMid", t, WithDao(func(d *Dao) {
		_, err := d.ReportsCacheByMid(context.Background(), 123)
		So(err, ShouldBeNil)
	}))
}
