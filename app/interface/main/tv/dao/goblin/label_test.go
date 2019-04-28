package goblin

import (
	"context"
	"testing"

	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestGoblinLabel(t *testing.T) {
	var (
		c        = context.Background()
		category = int(1)
		catType  = int(1)
	)
	convey.Convey("Label", t, func(ctx convey.C) {
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			res, err := d.Label(c, category, catType)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
		ctx.Convey("db closed", func(ctx convey.C) {
			d.db.Close()
			_, err := d.Label(c, category, catType)
			ctx.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}
