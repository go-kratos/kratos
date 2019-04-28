package goblin

import (
	"context"
	"testing"

	"go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestGoblinChlInfo(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("ChlInfo", t, func(c convey.C) {
		c.Convey("Then err should be nil.chls should not be nil.", func(c convey.C) {
			chls, err := d.ChlInfo(ctx)
			c.So(err, convey.ShouldBeNil)
			c.So(chls, convey.ShouldNotBeNil)
		})
		c.Convey("db closed", func(c convey.C) {
			d.db.Close()
			_, err := d.ChlInfo(ctx)
			c.So(err, convey.ShouldNotBeNil)
			d.db = sql.NewMySQL(d.conf.Mysql)
		})
	})
}
