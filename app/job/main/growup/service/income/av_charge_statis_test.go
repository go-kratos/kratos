package income

import (
	"context"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/income"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetAvChargeStatisMap(t *testing.T) {
	Convey("GetAvChargeStatisMap", t, func() {
		_, err := charge.GetAvChargeStatisMap(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_GetAvChargeStatis(t *testing.T) {
	Convey("GetAvChargeStatis", t, func() {
		_, err := charge.GetAvChargeStatis(context.Background())
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
		err := charge.AvChargeStatisDBStore(context.Background(), chargeStatisMap)
		So(err, ShouldBeNil)
	})
}

func benchmarkAvChargeStatisDBStore(bsize int, size int64, b *testing.B) {
	batchSize = bsize
	var i int64
	chargeStatisMap := make(map[int64]*model.AvChargeStatis)
	for i = 0; i < size; i++ {
		chargeStatisMap[i] = &model.AvChargeStatis{
			AvID:        i,
			MID:         i,
			TagID:       i,
			IsOriginal:  int(i),
			UploadTime:  xtime.Time(time.Now().Unix()),
			TotalCharge: i,
			DBState:     int(i % 2),
		}
	}

	for n := 0; n < b.N; n++ {
		charge.AvChargeStatisDBStore(context.Background(), chargeStatisMap)
	}
}

func BenchmarkAvChargeStatisDBStore100(b *testing.B) {
	benchmarkAvChargeStatisDBStore(100, 100000, b)
}
func BenchmarkAvChargeStatisDBStore1000(b *testing.B) {
	benchmarkAvChargeStatisDBStore(1000, 100000, b)
}
func BenchmarkAvChargeStatisDBStore2000(b *testing.B) {
	benchmarkAvChargeStatisDBStore(2000, 100000, b)
}
func BenchmarkAvChargeStatisDBStore10000(b *testing.B) {
	benchmarkAvChargeStatisDBStore(10000, 100000, b)
}
