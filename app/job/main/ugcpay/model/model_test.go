package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestElecUserSetting(t *testing.T) {
	Convey("", t, func() {
		data := &DBOldElecUserSetting{
			MID:       46333,
			SettingID: 1,
			Status:    1,
		}
		So(data.BitValue(), ShouldEqual, 0x01)
		data = &DBOldElecUserSetting{
			MID:       46333,
			SettingID: 2,
			Status:    1,
		}
		So(data.BitValue(), ShouldEqual, 0x02)
	})
}
