package dao

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaowaveFormKey(t *testing.T) {
	convey.Convey("waveFormKey", t, func(ctx convey.C) {
		var (
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := testDao.waveFormKey(oid, tp)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetWaveFormCache(t *testing.T) {
	convey.Convey("SetWaveFormCache", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			waveForm = &model.WaveForm{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.SetWaveFormCache(c, waveForm)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWaveFormCache(t *testing.T) {
	convey.Convey("WaveFormCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			waveForm, err := testDao.WaveFormCache(c, oid, tp)
			ctx.Convey("Then err should be nil.waveForm should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(waveForm, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelWaveFormCache(t *testing.T) {
	convey.Convey("DelWaveFormCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			oid = int64(0)
			tp  = int32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := testDao.DelWaveFormCache(c, oid, tp)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
