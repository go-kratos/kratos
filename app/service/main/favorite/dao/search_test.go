package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchArcs(t *testing.T) {
	Convey("SearchArcs", t, func() {
		var (
			keyword = "comic"
			ids     = []int64{1, 2, 3}
			pn      = 1
			ps      = 30
		)
		err, _ := d.SearchArcs(context.Background(), keyword, ids, pn, ps) // todo:debug uat 502
		So(err, ShouldBeNil)
	})
}
