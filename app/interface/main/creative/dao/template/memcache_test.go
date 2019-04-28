package template

import (
	"context"
	"go-common/app/interface/main/creative/model/template"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTemplatekeyTpl(t *testing.T) {
	var (
		mid = int64(2089809)
	)
	convey.Convey("keyTpl", t, func(ctx convey.C) {
		p1 := keyTpl(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestTemplatetplCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("tplCache", t, func(ctx convey.C) {
		tps, err := d.tplCache(c, mid)
		ctx.Convey("Then err should be nil.tps should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tps, convey.ShouldBeNil)
		})
	})
}

func TestTemplateaddTplCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		tps = []*template.Template{}
	)
	convey.Convey("addTplCache", t, func(ctx convey.C) {
		err := d.addTplCache(c, mid, tps)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestTemplatedelTplCache(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("delTplCache", t, func(ctx convey.C) {
		err := d.delTplCache(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
