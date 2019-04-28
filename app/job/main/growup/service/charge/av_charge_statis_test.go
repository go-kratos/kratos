package charge

import (
	"context"
	"testing"

	model "go-common/app/job/main/growup/model/charge"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetAvChargeStatisMap(t *testing.T) {
	Convey("GetAvChargeStatisMap", t, func() {
		_, err := s.GetAvChargeStatisMap(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_GetAvChargeStatis(t *testing.T) {
	Convey("GetAvChargeStatis", t, func() {
		_, err := s.GetAvChargeStatis(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_AvChargeStatisDBStore(t *testing.T) {
	Convey("AvChargeStatisDBStore", t, func() {
		chargeStatisMap := make(map[int64]*model.AvChargeStatis)
		value := &model.AvChargeStatis{
			AvID:    11,
			MID:     11,
			TagID:   11,
			DBState: 1,
		}
		chargeStatisMap[11] = value
		err := s.AvChargeStatisDBStore(context.Background(), chargeStatisMap)
		So(err, ShouldBeNil)
	})
}
