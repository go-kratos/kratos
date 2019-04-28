package dao

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"testing"
	"time"
)

func TestDao_Tx1(t *testing.T) {
	Convey("commit wrong sample", t, func() {
		once.Do(startService)

		uid := int64(1)
		var wallet *model.Melonseed
		wallet, err := d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		og := wallet.Gold
		So(err, ShouldBeNil)
		So(1, ShouldEqual, 1)
		tx, err := d.db.Begin(ctx)
		So(err, ShouldBeNil)
		affect, err := d.AddGold(ctx, uid, 1)
		So(affect, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
		wallet, err = d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		So(err, ShouldBeNil)
		So(wallet.Gold-og, ShouldEqual, 1)
	})

	Convey("rollback wrong sample use db ins instead tx", t, func() {
		once.Do(startService)

		uid := int64(1)
		var wallet *model.Melonseed
		wallet, err := d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		og := wallet.Gold
		So(err, ShouldBeNil)
		So(1, ShouldEqual, 1)
		tx, err := d.db.Begin(ctx)
		So(err, ShouldBeNil)
		affect, err := d.AddGold(ctx, uid, 1)
		So(affect, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
		err = tx.Rollback()
		So(err, ShouldBeNil)
		wallet, err = d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		So(err, ShouldBeNil)
		So(wallet.Gold-og, ShouldEqual, 1)
	})

	Convey("commit", t, func() {
		once.Do(startService)

		uid := int64(1)
		var wallet *model.Melonseed
		wallet, err := d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		og := wallet.Gold
		So(err, ShouldBeNil)
		So(1, ShouldEqual, 1)
		tx, err := d.db.Begin(ctx)
		So(err, ShouldBeNil)
		res, err := tx.Exec(fmt.Sprintf("update user_wallet_%d set gold = gold + ? where uid = ?", uid%10), 1, uid)
		So(err, ShouldBeNil)
		affect, err := res.RowsAffected()
		So(affect, ShouldEqual, 1)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
		wallet, err = d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		So(err, ShouldBeNil)
		So(wallet.Gold-og, ShouldEqual, 1)
	})

	Convey("rollback", t, func() {
		once.Do(startService)

		uid := int64(1)
		var wallet *model.Melonseed
		wallet, err := d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		og := wallet.Gold
		So(err, ShouldBeNil)
		So(1, ShouldEqual, 1)
		tx, err := d.db.Begin(ctx)
		So(err, ShouldBeNil)
		res, err := tx.Exec(fmt.Sprintf("update user_wallet_%d set gold = gold + ? where uid = ?", uid%10), 1, uid)
		So(err, ShouldBeNil)
		affect, err := res.RowsAffected()
		So(affect, ShouldEqual, 1)
		So(err, ShouldBeNil)
		err = tx.Rollback()
		So(err, ShouldBeNil)
		wallet, err = d.Melonseed(ctx, uid)
		t.Logf("uid gold : %d", wallet.Gold)
		So(err, ShouldBeNil)
		So(wallet.Gold-og, ShouldEqual, 0)
	})

}

func TestDao_DoubleCoin(t *testing.T) {
	Convey("check1", t, func() {
		d := &model.DetailWithSnapShot{}
		So(model.NeedSnapshot(d, time.Now()), ShouldBeTrue)
		d.SnapShotTime = "2018-09-22 16:39:05"
		So(model.NeedSnapshot(d, time.Now()), ShouldBeTrue)

		now, _ := time.Parse("2006-01-02 15:04:05", "2018-09-22 16:39:04")
		t.Logf("today:%+v", now)
		So(model.NeedSnapshot(d, now), ShouldBeFalse)

		now = model.GetTodayTime(now)
		t.Logf("today:%+v", now)
		d.SnapShotTime = "2018-09-22 04:39:04"
		So(model.NeedSnapshot(d, now), ShouldBeFalse)

		d.SnapShotTime = "2018-09-21 04:39:04"
		So(model.NeedSnapshot(d, now), ShouldBeTrue)

		d.SnapShotTime = time.Now().Format("2006-01-02 15:04:05")
		t.Logf("today:%+v", d.SnapShotTime)
		So(model.TodayNeedSnapShot(d), ShouldBeFalse)

		d.SnapShotTime = time.Now().Add(-time.Second * 86400).Format("2006-01-02 15:04:05")
		t.Logf("today:%+v", d.SnapShotTime)
		So(model.TodayNeedSnapShot(d), ShouldBeTrue)

	})
}
