package vip

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/xuser/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	xtime "go-common/library/time"
	"math/rand"
	"testing"
	"time"
)

func Test_getUserLevelTable(t *testing.T) {
	initd()
	Convey("test get user_x table name by uid", t, func() {
		var uid = int64(123)
		table := getUserLevelTable(uid)
		So(table, ShouldEqual, "user_2")
	})
}

func Test_getUserVipRecordTable(t *testing.T) {
	initd()
	Convey("test get user_vip_record_x table name by uid", t, func() {
		var uid = rand.Int63()
		log.Info("Test_getUserVipRecordTable uid(%d)", uid)
		table := getUserVipRecordTable(uid)
		t := fmt.Sprintf(_userVipRecordPrefix, uid%_userVipRecordCount)
		So(table, ShouldEqual, t)
	})
}

func TestDao_GetVipFromDB(t *testing.T) {
	initd()
	Convey("test get vip from db", t, testWithTestUser(func(u *TestUser) {
		log.Info("TestDao_GetVipFromDB uid(%d), table(%s)", u.Uid, getUserLevelTable(u.Uid))
		var (
			ctx  = context.Background()
			err  error
			info *model.VipInfo
		)

		// delete random uid at begin
		err = d.deleteVip(ctx, u.Uid)
		So(err, ShouldBeNil)

		// get nil result from db
		Convey("get nil result from db", func() {
			info, err = d.GetVipFromDB(ctx, u.Uid)
			So(err, ShouldResemble, sql.ErrNoRows)
			So(info, ShouldNotBeNil)
			So(info.Vip, ShouldEqual, 0)
			So(info.VipTime, ShouldEqual, "")
			So(info.Svip, ShouldEqual, 0)
			So(info.SvipTime, ShouldEqual, "")
		})

		// insert and then get
		Convey("insert and then get", func() {
			var info2 *model.VipInfo
			info = &model.VipInfo{
				Vip:      1,
				VipTime:  time.Now().Add(time.Hour * 12).Format(model.TimeNano),
				Svip:     1,
				SvipTime: time.Now().Add(time.Hour * 12).Format(model.TimeNano),
			}
			err = d.createVip(ctx, u.Uid, info)
			So(err, ShouldBeNil)

			info2, err = d.GetVipFromDB(ctx, u.Uid)
			log.Info("TestDao_GetVipFromDB info2(%v)", info2)
			So(err, ShouldBeNil)
			So(info, ShouldResemble, info2)
		})

		// valid info but expired vip time
		Convey("insert valid but expired vip time", func() {
			info = &model.VipInfo{
				Vip:      1,
				VipTime:  time.Now().AddDate(0, -1, 0).Format(model.TimeNano),
				Svip:     1,
				SvipTime: time.Now().AddDate(-1, 0, 0).Format(model.TimeNano),
			}
			err = d.createVip(ctx, u.Uid, info)
			So(err, ShouldBeNil)

			info2, err := d.GetVipFromDB(ctx, u.Uid)
			log.Info("TestDao_GetVipFromDB info2(%v)", info2)
			So(err, ShouldBeNil)
			So(info2.Vip, ShouldEqual, 0)
			So(info2.Svip, ShouldEqual, 0)
			So(info2.VipTime, ShouldEqual, info.VipTime)
			So(info2.SvipTime, ShouldEqual, info.SvipTime)
		})
	}))
}

func TestDao_AddVip(t *testing.T) {
	initd()
	Convey("test add vip", t, testWithTestUser(func(u *TestUser) {
		log.Info("TestDao_GetVipFromDB uid(%d), table(%s)", u.Uid, getUserLevelTable(u.Uid))
		var (
			ctx  = context.Background()
			err  error
			info *model.VipInfo
			row  int64
		)

		// create one row at begin
		info = &model.VipInfo{
			Vip:      1,
			VipTime:  time.Now().Add(time.Hour * 12).Format(model.TimeNano),
			Svip:     1,
			SvipTime: time.Now().Add(time.Hour * 12).Format(model.TimeNano),
		}
		err = d.createVip(ctx, u.Uid, info)
		So(err, ShouldBeNil)

		// empty vip time, should return err
		Convey("empty vip time, should return err", func() {
			row, err = d.AddVip(ctx, u.Uid, 0, 0)
			So(row, ShouldEqual, 0)
			So(err, ShouldResemble, errUpdateVipTimeInvalid)
		})

		// add vip and svip time
		Convey("add vip and svip time", func() {
			// add one month vip
			dt := xtime.Time(30 * 86400)
			vtime, err := time.Parse(model.TimeNano, info.VipTime)
			So(err, ShouldBeNil)
			newvtime := xtime.Time(vtime.Unix()) + dt
			log.Info("TestDao_AddVip info(%v), oldvt(%v), newvtime(%v), dt(%v)", info, vtime.Unix(), newvtime, dt)
			row, err := d.AddVip(ctx, u.Uid, newvtime, 0)
			So(row, ShouldEqual, 1)
			So(err, ShouldBeNil)

			info2, err := d.GetVipFromDB(ctx, u.Uid)
			log.Info("TestDao_AddVip info2(%v)", info2)
			So(err, ShouldBeNil)
			So(info2.Vip, ShouldEqual, info.Vip)
			So(info2.Svip, ShouldEqual, info.Svip)
			So(info2.SvipTime, ShouldEqual, info.SvipTime)
			So(info2.VipTime, ShouldEqual, newvtime.Time().Format(model.TimeNano))
		})
	}))
}
