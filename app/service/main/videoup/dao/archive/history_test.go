package archive

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/model/archive"
)

func TestDao_TxAddArcHistory(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxAddArcHistory", t, func(ctx C) {
		_, err := d.TxAddArcHistory(tx, 23333, 123, "ssss", "content", "", "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAddVideoHistory(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		v     = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
	)
	Convey("TxAddVideoHistory", t, func(ctx C) {
		_, err := d.TxAddVideoHistory(tx, 23333, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoHistory(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoHistory", t, func(ctx C) {
		_, err := d.TxUpVideoHistory(tx, 23333, 1212, "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAddVideoHistorys(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		vs    = []*archive.Video{{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}}
	)
	Convey("TxAddVideoHistorys", t, func(ctx C) {
		err := d.TxAddVideoHistorys(tx, 23333, vs)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestArchiveArcHistory(t *testing.T) {
	var (
		c   = context.Background()
		hid = int64(23333)
	)
	Convey("ArcHistory", t, func(ctx C) {
		_, err := d.ArcHistory(c, hid)
		ctx.Convey("Then err should be nil.ah should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveArcHistorys(t *testing.T) {
	var (
		c     = context.Background()
		aid   = int64(23333)
		stime = time.Now()
	)
	Convey("ArcHistorys", t, func(ctx C) {
		_, err := d.ArcHistorys(c, aid, stime)
		ctx.Convey("Then err should be nil.ahs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveVideoHistory(t *testing.T) {
	var (
		c   = context.Background()
		hid = int64(23333)
	)
	Convey("VideoHistory", t, func(ctx C) {
		_, err := d.VideoHistory(c, hid)
		ctx.Convey("Then err should be nil.vhs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}
