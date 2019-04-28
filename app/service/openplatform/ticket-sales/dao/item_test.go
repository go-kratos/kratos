package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestItemBillInfo(t *testing.T) {
	convey.Convey("ItemBillInfo", t, func() {
		itemID := int64(676)
		scID := int64(870)
		tkID := int64(2843)
		infos, _ := d.ItemBillInfo(context.TODO(), []int64{itemID}, []int64{scID}, []int64{tkID})
		convey.So(infos.BaseInfo, convey.ShouldContainKey, itemID)
		convey.So(infos.BaseInfo[itemID].Screen, convey.ShouldContainKey, scID)
		convey.So(infos.BaseInfo[itemID].Screen[scID].Ticket, convey.ShouldContainKey, tkID)
	})
}
