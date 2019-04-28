package dao

import (
	"testing"

	"go-common/app/admin/main/tv/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSavePanel(t *testing.T) {
	convey.Convey("SavePanel", t, func(ctx convey.C) {
		var (
			panel = &model.TvPriceConfig{ID: 100000000}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SavePanel(panel)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetById(t *testing.T) {
	convey.Convey("GetById", t, func(ctx convey.C) {
		var (
			id = int64(100000000)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			panelInfo, err := d.GetById(id)
			ctx.Convey("Then err should be nil.panelInfo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(panelInfo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPanelStatus(t *testing.T) {
	convey.Convey("PanelStatus", t, func(ctx convey.C) {
		var (
			id     = int64(100000000)
			status = int64(2)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.PanelStatus(id, status)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExistProduct(t *testing.T) {
	convey.Convey("ExistProduct", t, func(ctx convey.C) {
		var (
			productID = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			flag := d.ExistProduct(productID)
			ctx.Convey("Then flag should not be nil.", func(ctx convey.C) {
				ctx.So(flag, convey.ShouldNotBeNil)
			})
		})
	})
}
