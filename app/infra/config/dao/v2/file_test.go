package v2

import (
	"go-common/app/infra/config/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestV2SetFile(t *testing.T) {
	var (
		name = "main.common-arch.apm-admin_460"
		conf = &model.Content{}
	)
	convey.Convey("SetFile", t, func(ctx convey.C) {
		err := d.SetFile(name, conf)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2File(t *testing.T) {
	var (
		name = "main.common-arch.apm-admin_460"
	)
	convey.Convey("File", t, func(ctx convey.C) {
		res, err := d.File(name)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestV2DelFile(t *testing.T) {
	var (
		name = "main.common-arch.apm-admin_460"
	)
	convey.Convey("DelFile", t, func(ctx convey.C) {
		err := d.DelFile(name)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2SetFileStr(t *testing.T) {
	var (
		name = "main.common-arch.apm-admin_460"
		val  = "test"
	)
	convey.Convey("SetFileStr", t, func(ctx convey.C) {
		err := d.SetFileStr(name, val)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestV2FileStr(t *testing.T) {
	var (
		name = "main.common-arch.apm-admin_460"
	)
	convey.Convey("FileStr", t, func(ctx convey.C) {
		file, err := d.FileStr(name)
		ctx.Convey("Then err should be nil.file should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(file, convey.ShouldNotBeNil)
		})
	})
}
