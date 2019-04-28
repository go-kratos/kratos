package creative

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCreativeTool(t *testing.T) {
	var (
		c  = context.TODO()
		ty = "'play'"
	)
	convey.Convey("Tool", t, func(ctx convey.C) {
		_, err := d.Tool(c, ty)
		ctx.Convey("Then err should be nil.ops should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCreativeOperations(t *testing.T) {
	var (
		c   = context.TODO()
		tys = []string{"'play'", "'notice'", "'road'", "'creative'", "'collect_arc'"}
	)
	convey.Convey("Operations", t, func(ctx convey.C) {
		_, err := d.Operations(c, tys)
		ctx.Convey("Then err should be nil.ops should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCreativeAllOperByTypeSQL(t *testing.T) {
	var (
		c   = context.TODO()
		tys = []string{"'play'", "'notice'", "'road'", "'creative'", "'collect_arc'"}
	)
	convey.Convey("AllOperByTypeSQL", t, func(ctx convey.C) {
		_, err := d.AllOperByTypeSQL(c, tys)
		ctx.Convey("Then err should be nil.ops should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
