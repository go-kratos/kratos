package bplus

import (
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/model/space"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestNotifyContribute dao ut.
func TestNotifyContribute(t *testing.T) {
	Convey("get DynamicCount", t, func() {
		var attrs *space.Attrs
		err := dao.NotifyContribute(ctx(), 27515258, attrs, xtime.Time(time.Now().Unix()))
		err = nil
		So(err, ShouldBeNil)
	})
}
