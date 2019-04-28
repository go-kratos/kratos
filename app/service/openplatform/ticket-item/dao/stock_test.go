package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestStockChanged
func TestDao_StockChanged(t *testing.T) {
	Convey("StockChanged", t, func() {
		once.Do(startService)
		res := d.StockChanged([]int64{1015, 1016})

		So(res, ShouldBeTrue)
	})
}
