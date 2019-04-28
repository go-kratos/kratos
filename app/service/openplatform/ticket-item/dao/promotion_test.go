package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestHasPromotion
func TestDao_HasPromotion(t *testing.T) {
	Convey("HasPromotion", t, func() {
		once.Do(startService)
		res := d.HasPromotion([]int64{78}, 1)

		So(res, ShouldBeFalse)
	})
}
