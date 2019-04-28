package archive

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAddit_InnerAttrSet(t *testing.T) {
	convey.Convey("InnerAttrSet", t, func() {
		add := &Addit{
			InnerAttr: 0,
		}
		add.InnerAttrSet(1, InnerAttrChannelReview)
		convey.So(add.InnerAttr, convey.ShouldEqual, 1)

		add.InnerAttr = 16
		add.InnerAttrSet(1, 3)
		convey.So(add.InnerAttr, convey.ShouldEqual, 24)
	})
}
