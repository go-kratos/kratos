package search

import (
	"context"
	"testing"

	searchMdl "go-common/app/interface/main/tv/model/search"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchSearchSug(t *testing.T) {
	var (
		ctx = context.Background()
		req = &searchMdl.ReqSug{
			MobiApp:  "android_tv_yst",
			Build:    "1011",
			Platform: "android",
			Term:     "test",
		}
	)
	convey.Convey("SearchSug", t, func(c convey.C) {
		result, err := d.SearchSug(ctx, req)
		c.Convey("Then err should be nil.result should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(result, convey.ShouldNotBeNil)
		})
	})
}
