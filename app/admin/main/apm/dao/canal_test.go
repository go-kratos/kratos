package dao

import (
	"fmt"
	"reflect"
	"testing"

	cml "go-common/app/admin/main/apm/model/canal"

	"github.com/bouk/monkey"
	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetConfigID(t *testing.T) {
	convey.Convey("SetConfigID", t, func(ctx convey.C) {
		var (
			id   = int64(0)
			addr = "127.0.0.1:8000"
			db   = &gorm.DB{
				Error: fmt.Errorf("test"),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetConfigID(id, addr)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When DB update return err", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Updates", func(_ *gorm.DB, _ interface{}, _ ...bool) *gorm.DB {
				return db
			})
			err := d.SetConfigID(id, addr)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, 70014)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestDaoCanalInfoCounts(t *testing.T) {
	convey.Convey("CanalInfoCounts", t, func(ctx convey.C) {
		var (
			v = &cml.ConfigReq{
				Addr:      "127.0.0.1:3308",
				User:      "admin",
				Password:  "admin",
				Project:   "main.web-svr",
				Leader:    "fss",
				Databases: "ada",
				Mark:      "test",
			}
			db = &gorm.DB{
				Error: fmt.Errorf("test"),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cnt, err := d.CanalInfoCounts(v)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When count error", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Count", func(_ *gorm.DB, _ interface{}) *gorm.DB {
				return db
			})
			cnt, err := d.CanalInfoCounts(v)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, -400)
				ctx.So(cnt, convey.ShouldEqual, 0)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestDaoCanalInfoEdit(t *testing.T) {
	convey.Convey("CanalInfoEdit", t, func(ctx convey.C) {
		var (
			v = &cml.ConfigReq{
				Addr:      "127.0.0.1:3308",
				User:      "admin",
				Password:  "admin",
				Project:   "main.web-svr",
				Leader:    "fss",
				Databases: "ada",
				Mark:      "test",
			}
			db = &gorm.DB{
				Error: fmt.Errorf("test"),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.CanalInfoEdit(v)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When edit return err", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Updates", func(_ *gorm.DB, _ interface{}, _ ...bool) *gorm.DB {
				return db
			})
			err := d.CanalInfoEdit(v)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, 70005)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestDaoCanalApplyCounts(t *testing.T) {
	convey.Convey("CanalApplyCounts", t, func(ctx convey.C) {
		var (
			v = &cml.ConfigReq{
				Addr:      "127.0.0.1:3308",
				User:      "admin",
				Password:  "admin",
				Project:   "main.web-svr",
				Leader:    "fss",
				Databases: "ada",
				Mark:      "test",
			}
			db = &gorm.DB{
				Error: fmt.Errorf("test"),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cnt, err := d.CanalApplyCounts(v)
			ctx.Convey("Then err should be nil.cnt should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cnt, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When count error", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Count", func(_ *gorm.DB, _ interface{}) *gorm.DB {
				return db
			})
			cnt, err := d.CanalApplyCounts(v)
			ctx.Convey("Then err should not be nil.cnt should  be 0.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, -400)
				ctx.So(cnt, convey.ShouldEqual, 0)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestDaoCanalApplyEdit(t *testing.T) {
	convey.Convey("CanalApplyEdit", t, func(ctx convey.C) {
		var (
			v = &cml.ConfigReq{
				Addr:      "127.0.0.1:3308",
				User:      "admin",
				Password:  "admin",
				Databases: "test",
				Mark:      "test",
			}
			db = &gorm.DB{
				Error: fmt.Errorf("test"),
			}
			username = "fengshanshan"
		)

		ctx.Convey("When project and leader is null", func(ctx convey.C) {
			err := d.CanalApplyEdit(v, username)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			v.Leader = "fengshanshan"
			v.Project = "main.web-svr"
			err := d.CanalApplyEdit(v, username)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When edit error", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Updates", func(_ *gorm.DB, _ interface{}, _ ...bool) *gorm.DB {
				return db
			})
			err := d.CanalApplyEdit(v, username)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, 70005)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestDaoCanalApplyCreate(t *testing.T) {
	convey.Convey("CanalApplyCreate", t, func(ctx convey.C) {
		var (
			v = &cml.ConfigReq{
				Addr:      "127.0.0.1:3309",
				User:      "admin",
				Password:  "admin",
				Project:   "main.web-svr",
				Leader:    "fss",
				Databases: "ada",
				Mark:      "test",
			}
			db = &gorm.DB{
				Error: fmt.Errorf("test"),
			}
			username = "fengshanshan"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Create", func(_ *gorm.DB, _ interface{}) *gorm.DB {
				return &gorm.DB{
					Error: nil,
				}
			})
			err := d.CanalApplyCreate(v, username)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When creater error", func(ctx convey.C) {
			monkey.PatchInstanceMethod(reflect.TypeOf(d.DBCanal), "Create", func(_ *gorm.DB, _ interface{}) *gorm.DB {
				return db
			})
			err := d.CanalApplyCreate(v, username)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, 70006)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}
