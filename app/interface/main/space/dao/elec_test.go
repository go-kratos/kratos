package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_ElecInfo(t *testing.T) {
	convey.Convey("test elec info", t, func(ctx convey.C) {
		mid := int64(28272030)
		paymid := int64(0)
		data, err := d.ElecInfo(context.Background(), mid, paymid)
		convey.So(err, convey.ShouldBeNil)
		convey.Printf("%v", data)
	})
}
