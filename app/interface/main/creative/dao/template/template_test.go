package template

import (
	"context"
	"go-common/app/interface/main/creative/model/template"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTemplatetemplates(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("templates", t, func(ctx convey.C) {
		tps, err := d.templates(c, mid)
		ctx.Convey("Then err should be nil.tps should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tps, convey.ShouldNotBeNil)
		})
	})
}

func TestTemplateTemplate(t *testing.T) {
	var (
		c   = context.TODO()
		id  = int64(1)
		mid = int64(2089809)
	)
	convey.Convey("Template", t, func(ctx convey.C) {
		no, err := d.Template(c, id, mid)
		ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(no, convey.ShouldBeNil)
		})
	})
}

func TestTemplateTemplates(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("Templates", t, func(ctx convey.C) {
		tps, err := d.Templates(c, mid)
		ctx.Convey("Then err should be nil.tps should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(tps, convey.ShouldNotBeNil)
		})
	})
}

func TestTemplateAddTemplate(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		tp  = &template.Template{}
	)
	convey.Convey("AddTemplate", t, func(ctx convey.C) {
		id, err := d.AddTemplate(c, mid, tp)
		ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestTemplateUpTemplate(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		tp  = &template.Template{}
	)
	convey.Convey("UpTemplate", t, func(ctx convey.C) {
		rows, err := d.UpTemplate(c, mid, tp)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestTemplateDelTemplate(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
		tp  = &template.Template{}
	)
	convey.Convey("DelTemplate", t, func(ctx convey.C) {
		rows, err := d.DelTemplate(c, mid, tp)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rows, convey.ShouldNotBeNil)
		})
	})
}

func TestTemplateCount(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(2089809)
	)
	convey.Convey("Count", t, func(ctx convey.C) {
		count, err := d.Count(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}
