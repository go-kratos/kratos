package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusinessAllV2(t *testing.T) {
	convey.Convey("BusinessAllV2", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, err := d.BusinessAllV2(c)
			convCtx.Convey("Then err should be nil.list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBusinessInfoV2(t *testing.T) {
	convey.Convey("BusinessInfoV2", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			name = "dm"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			b, err := d.BusinessInfoV2(c, name)
			convCtx.Convey("Then err should be nil.b should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(b, convey.ShouldNotBeNil)
			})
		})
	})
}

//func TestDaoBusinessIns(t *testing.T) {
//	convey.Convey("BusinessIns", t, func(convCtx convey.C) {
//		var (
//			c           = context.Background()
//			pid         = int64(0)
//			name        = ""
//			description = ""
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			rows, err := d.BusinessIns(c, pid, name, description)
//			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(rows, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

//func TestDaoBusinessUpdate(t *testing.T) {
//	convey.Convey("BusinessUpdate", t, func(convCtx convey.C) {
//		var (
//			c     = context.Background()
//			name  = ""
//			field = ""
//			value = ""
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			rows, err := d.BusinessUpdate(c, name, field, value)
//			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(rows, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestDaoAssetDBTables(t *testing.T) {
	convey.Convey("AssetDBTables", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, err := d.AssetDBTables(c)
			convCtx.Convey("Then err should be nil.list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

//
//func TestDaoAssetDBIns(t *testing.T) {
//	convey.Convey("AssetDBIns", t, func(convCtx convey.C) {
//		var (
//			c           = context.Background()
//			name        = ""
//			description = ""
//			dsn         = ""
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			rows, err := d.AssetDBIns(c, name, description, dsn)
//			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(rows, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

//func TestDaoAssetTableIns(t *testing.T) {
//	convey.Convey("AssetTableIns", t, func(convCtx convey.C) {
//		var (
//			c           = context.Background()
//			name        = ""
//			db          = ""
//			regex       = ""
//			fields      = ""
//			description = ""
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			rows, err := d.AssetTableIns(c, name, db, regex, fields, description)
//			convCtx.Convey("Then err should be nil.rows should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(rows, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestDaoAsset(t *testing.T) {
	convey.Convey("Asset", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			name = "bilibili_article"
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			r, err := d.Asset(c, name)
			convCtx.Convey("Then err should be nil.r should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(r, convey.ShouldNotBeNil)
			})
		})
	})
}
