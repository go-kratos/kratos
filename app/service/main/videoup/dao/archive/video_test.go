package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/model/archive"
	"math/rand"
	"time"
)

func TestDao_TxAddVideo(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		v     = &archive.Video{
			ID:    123,
			Aid:   2333,
			Title: "UT测试",
		}
	)
	rand.Seed(time.Now().Unix())
	v.ID = int64(rand.Intn(999999999) + 1000000000)
	v.Aid = int64(rand.Intn(999999999) + 1000000000)
	Convey("TxAddVideo", t, func(ctx C) {
		_, err := d.TxAddVideo(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideo(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		v     = &archive.Video{
			ID:    123,
			Aid:   2333,
			Title: "UT测试",
		}
	)
	Convey("TxUpVideo", t, func(ctx C) {
		_, err := d.TxUpVideo(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoStatus(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoStatus", t, func(ctx C) {
		_, err := d.TxUpVideoStatus(tx, 2333, "sadasdadsds", 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoXcode(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoXcode", t, func(ctx C) {
		_, err := d.TxUpVideoXcode(tx, 2333, "sadasdadsds", 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoAttr(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoAttr", t, func(ctx C) {
		_, err := d.TxUpVideoAttr(tx, 2333, "sadasdadsds", 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpVideoCid(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoCid", t, func(ctx C) {
		_, err := d.TxUpVideoCid(tx, 2333, "sadasdadsds", 1213)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxAddAudit(t *testing.T) {
	rand.Seed(time.Now().Unix())
	vid := int64(rand.Intn(999999) + 1000000000)
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		vs    = []*archive.Video{{
			ID:    vid,
			Aid:   2333,
			Title: "UT测试",
		}}
	)

	Convey("TxAddAudit", t, func(ctx C) {
		_, err := d.TxAddAudit(tx, vs)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
