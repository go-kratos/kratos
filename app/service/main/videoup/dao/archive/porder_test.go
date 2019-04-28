package archive

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/model/archive"
)

func TestDao_TxUpPorder(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		ap    = &archive.ArcParam{
			Porder: &archive.Porder{},
		}
	)
	Convey("TxUpPorder", t, func(ctx C) {
		_, err := d.TxUpPorder(tx, 23333, ap)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestArchivePorder(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(23333)
	)
	Convey("Porder", t, func(ctx C) {
		_, err := d.Porder(c, aid)
		ctx.Convey("Then err should be nil.p should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchivePorderCfgList(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("PorderCfgList", t, func(ctx C) {
		_, err := d.PorderCfgList(c)
		ctx.Convey("Then err should be nil.pcfgs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchivePorderArcList(t *testing.T) {
	var (
		c     = context.Background()
		begin = time.Now()
		end   = time.Now()
	)
	Convey("PorderArcList", t, func(ctx C) {
		_, err := d.PorderArcList(c, begin, end)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}
