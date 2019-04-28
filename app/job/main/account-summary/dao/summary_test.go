package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSave(t *testing.T) {
	var (
		// ctx  = context.Background()·
		key  = ""
		data map[string][]byte
	)
	convey.Convey("Save", t, func(ctx convey.C) {
		err := d.Save(context.TODO(), key, data)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoGetByKey(t *testing.T) {
	var (
		// ctx = context.Background()
		key = ""
	)
	convey.Convey("GetByKey", t, func(ctx convey.C) {
		p1, err := d.GetByKey(context.TODO(), key)
		ctx.Convey("Then err should be nil.p1 should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}
