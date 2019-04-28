package kfc

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSearchList(t *testing.T) {
	convey.Convey("SearchList", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			code = ""
			mid  = int64(1505589)
			pn   = 1
			ps   = 15
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rp, err := d.SearchList(c, code, mid, pn, ps)
			ctx.Convey("Then err should be nil.rp should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				for _, v := range rp {
					fmt.Printf("%+v", v)
				}

			})
		})
	})
}
