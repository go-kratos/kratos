package income

import (
	"context"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/income"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetAvCharge(t *testing.T) {
	Convey("GetAvCharge", t, func() {
		_, err := charge.GetAvCharge(context.Background(), time.Now(), charge.dao.AvMonthlyCharge)
		So(err, ShouldBeNil)
	})
}

func Test_AvChargeDBStore(t *testing.T) {
	Convey("AvChargeDBStore", t, func() {
		monthlyChargeMap := make(map[int64]*model.AvCharge)
		value := &model.AvCharge{
			AvID:    11,
			MID:     11,
			TagID:   11,
			DBState: 1,
		}
		monthlyChargeMap[11] = value
		err := charge.AvChargeDBStore(context.Background(), monthlyChargeMap, monthlyChargeMap)
		So(err, ShouldBeNil)
	})
}

func benchmarkAvChargeDBStore(bsize int, size int64, b *testing.B) {
	batchSize = bsize
	for n := 0; n < b.N; n++ {
		var i int64
		weeklyChargeMap := make(map[int64]*model.AvCharge, size)
		for i = 0; i < size; i++ {
			v := int64(n+1) * i
			weeklyChargeMap[v] = &model.AvCharge{
				AvID:           v,
				MID:            v,
				TagID:          v,
				IsOriginal:     int(v),
				DanmakuCount:   v,
				CommentCount:   v,
				CollectCount:   v,
				CoinCount:      v,
				ShareCount:     v,
				ElecPayCount:   v,
				TotalPlayCount: v,
				WebPlayCount:   v,
				AppPlayCount:   v,
				H5PlayCount:    v,
				LvUnknown:      v,
				Lv0:            v,
				Lv1:            v,
				Lv2:            v,
				Lv3:            v,
				Lv4:            v,
				Lv5:            v,
				Lv6:            v,
				VScore:         v,
				IncCharge:      v,
				TotalCharge:    v,
				UploadTime:     xtime.Time(time.Now().Unix()),
				Date:           xtime.Time(time.Now().Unix()),
				DBState:        1,
			}
		}

		charge.AvChargeDBStore(context.Background(), weeklyChargeMap, weeklyChargeMap)
	}
}

func BenchmarkAvChargeDBStore100(b *testing.B) {
	benchmarkAvChargeDBStore(100, 100000, b)
}

func BenchmarkAvChargeDBStore1000(b *testing.B) {
	benchmarkAvChargeDBStore(1000, 100000, b)
}

func BenchmarkAvChargeDBStore2000(b *testing.B) {
	benchmarkAvChargeDBStore(2000, 100000, b)
}

func BenchmarkAvChargeDBStore10000(b *testing.B) {
	benchmarkAvChargeDBStore(10000, 100000, b)
}
