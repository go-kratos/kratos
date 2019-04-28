package archive

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_diff(t *testing.T) {
	e1 := &EditHistory{
		ArcHistory: &ArcHistory{Title: "haha"},
		VHistory:   []*VideoHistory{&VideoHistory{Filename: "hahah1"}, &VideoHistory{Filename: "hahah2"}},
	}
	var e2 *EditHistory
	convey.Convey("diff between one and nil", t, func() {
		diff, _ := e1.Diff(e2)
		convey.So(diff, convey.ShouldEqual, e1)
	})

	e2 = &EditHistory{
		ArcHistory: nil,
		VHistory:   []*VideoHistory{&VideoHistory{Filename: "hahah1"}},
	}
	convey.Convey("diff between 2 history", t, func() {
		diff, _ := e1.Diff(e2)
		convey.So(diff.ArcHistory.Title, convey.ShouldEqual, e1.ArcHistory.Title)
		convey.So(len(diff.VHistory), convey.ShouldEqual, 1)
	})
}
