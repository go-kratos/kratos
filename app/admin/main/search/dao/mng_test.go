package dao

import (
	"context"
	"go-common/app/admin/main/search/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoBusinessList(t *testing.T) {
	convey.Convey("BusinessList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			name   = ""
			offset = int(0)
			limit  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.BusinessList(c, name, offset, limit)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBusinessTotal(t *testing.T) {
	convey.Convey("BusinessTotal", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.BusinessTotal(c, name)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBusinessAll(t *testing.T) {
	convey.Convey("BusinessAll", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.BusinessAll(c)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddBusiness(t *testing.T) {
	convey.Convey("AddBusiness", t, func(ctx convey.C) {
		//var (
		//	c = context.Background()
		//	b = &model.MngBusiness{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	id, _ := d.AddBusiness(c, b)
		//	ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
		//		//ctx.So(err, convey.ShouldBeNil)
		//		ctx.So(id, convey.ShouldNotBeNil)
		//	})
		//})
	})
}

func TestDaoUpdateBusiness(t *testing.T) {
	convey.Convey("UpdateBusiness", t, func(ctx convey.C) {
		//var (
		//	c = context.Background()
		//	b = &model.MngBusiness{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	err := d.UpdateBusiness(c, b)
		//	ctx.Convey("Then err should be nil.", func(ctx convey.C) {
		//		ctx.So(err, convey.ShouldBeNil)
		//	})
		//})
	})
}

func TestDaoBusinessInfo(t *testing.T) {
	convey.Convey("BusinessInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.BusinessInfo(c, id)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBusinessInfoByName(t *testing.T) {
	convey.Convey("BusinessInfoByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "log"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.BusinessInfoByName(c, name)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAssetList(t *testing.T) {
	convey.Convey("AssetList", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			typ    = int(0)
			name   = ""
			offset = int(0)
			limit  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.AssetList(c, typ, name, offset, limit)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAssetTotal(t *testing.T) {
	convey.Convey("AssetTotal", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			typ  = int(0)
			name = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			total, err := d.AssetTotal(c, typ, name)
			ctx.Convey("Then err should be nil.total should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAssetAll(t *testing.T) {
	convey.Convey("AssetAll", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.AssetAll(c)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAssetInfo(t *testing.T) {
	convey.Convey("AssetInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.AssetInfo(c, id)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAssetInfoByName(t *testing.T) {
	convey.Convey("AssetInfoByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			info, err := d.AssetInfoByName(c, name)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddAsset(t *testing.T) {
	convey.Convey("AddAsset", t, func(ctx convey.C) {
		//var (
		//	c = context.Background()
		//	b = &model.MngAsset{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	//id, err :=
		//	d.AddAsset(c, b)
		//	ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
		//		//ctx.So(err, convey.ShouldBeNil)
		//		//ctx.So(id, convey.ShouldNotBeNil)
		//	})
		//})
	})
}

func TestDaoUpdateAsset(t *testing.T) {
	convey.Convey("UpdateAsset", t, func(ctx convey.C) {
		//var (
		//	c = context.Background()
		//	b = &model.MngAsset{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	err := d.UpdateAsset(c, b)
		//	ctx.Convey("Then err should be nil.", func(ctx convey.C) {
		//		ctx.So(err, convey.ShouldBeNil)
		//	})
		//})
	})
}

func TestDaoAppList(t *testing.T) {
	convey.Convey("AppList", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			business = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.AppList(c, business)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAppInfo(t *testing.T) {
	convey.Convey("AppInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//a, err :=
			d.AppInfo(c, id)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				//ctx.So(err, convey.ShouldBeNil)
				//ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAppInfoByAppid(t *testing.T) {
	convey.Convey("AppInfoByAppid", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			appid = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.AppInfoByAppid(c, appid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddApp(t *testing.T) {
	convey.Convey("AddApp", t, func(ctx convey.C) {
		//var (
		//	c = context.Background()
		//	a = &model.MngApp{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	id, err :=
		//		d.AddApp(c, a)
		//	ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
		//		ctx.So(err, convey.ShouldBeNil)
		//		ctx.So(id, convey.ShouldNotBeNil)
		//	})
		//})
	})
}

func TestDaoUpdateApp(t *testing.T) {
	convey.Convey("UpdateApp", t, func(ctx convey.C) {
		//var (
		//	c = context.Background()
		//	a = &model.MngApp{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	err := d.UpdateApp(c, a)
		//	ctx.Convey("Then err should be nil.", func(ctx convey.C) {
		//		ctx.So(err, convey.ShouldBeNil)
		//	})
		//})
	})
}

func TestDaoUpdateAppAssetTable(t *testing.T) {
	convey.Convey("UpdateAppAssetTable", t, func(ctx convey.C) {
		//var (
		//	c    = context.Background()
		//	name = ""
		//	no   = &model.MngAssetTable{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	err := d.UpdateAppAssetTable(c, name, no)
		//	ctx.Convey("Then err should be nil.", func(ctx convey.C) {
		//		ctx.So(err, convey.ShouldBeNil)
		//	})
		//})
	})
}

func TestDaoUpdateAppAssetDatabus(t *testing.T) {
	convey.Convey("UpdateAppAssetDatabus", t, func(ctx convey.C) {
		//var (
		//	c    = context.Background()
		//	name = ""
		//	v    = &model.MngAssetDatabus{}
		//)
		//ctx.Convey("When everything goes positive", func(ctx convey.C) {
		//	err := d.UpdateAppAssetDatabus(c, name, v)
		//	ctx.Convey("Then err should be nil.", func(ctx convey.C) {
		//		ctx.So(err, convey.ShouldBeNil)
		//	})
		//})
	})
}

func TestDaoMngCount(t *testing.T) {
	convey.Convey("MngCount", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.MngCount{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.MngCount(c, v)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoMngPercent(t *testing.T) {
	convey.Convey("MngPercent", t, func(ctx convey.C) {
		var (
			c = context.Background()
			v = &model.MngCount{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			list, err := d.MngPercent(c, v)
			ctx.Convey("Then err should be nil.list should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUnames(t *testing.T) {
	convey.Convey("Unames", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			uids = []string{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Unames(c, uids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
