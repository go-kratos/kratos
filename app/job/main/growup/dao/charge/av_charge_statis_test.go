package charge

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/charge"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvChargeStatis(t *testing.T) {
	Convey("AvChargeStatis", t, func() {
		_, err := d.AvChargeStatis(context.Background(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_InsertAvChargeStatisBatch(t *testing.T) {
	Convey("InsertAvChargeStatisBatch", t, func() {
		c := context.Background()
		d.db.Exec(c, "DELETE FROM av_charge_statis where av_id = 11")

		avChargeStatis := []*model.AvChargeStatis{}
		value := &model.AvChargeStatis{
			AvID:  11,
			MID:   11,
			TagID: 11,
		}
		avChargeStatis = append(avChargeStatis, value)
		vals := assembleAvChargeStatis(avChargeStatis)
		count, err := d.InsertAvChargeStatisBatch(c, vals)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 1)

		d.db.Exec(c, "DELETE FROM av_charge_statis where av_id = 11")
	})
}

func benchmarkInsertAvChargeStatisBatch(size int64, b *testing.B) {
	avChargeStatis := make([]*model.AvChargeStatis, size)
	var i int64
	for i = 0; i < size; i++ {
		avChargeStatis[i] = &model.AvChargeStatis{
			AvID:        i,
			MID:         i,
			TagID:       i,
			IsOriginal:  int(i),
			UploadTime:  xtime.Time(time.Now().Unix()),
			TotalCharge: i,
		}
	}

	vals := assembleAvChargeStatis(avChargeStatis)
	for n := 0; n < b.N; n++ {
		d.InsertAvChargeStatisBatch(context.Background(), vals)
	}
}

func BenchmarkInsertAvChargeStatisBatch100(b *testing.B)  { benchmarkInsertAvChargeStatisBatch(100, b) }
func BenchmarkInsertAvChargeStatisBatch1000(b *testing.B) { benchmarkInsertAvChargeStatisBatch(1000, b) }
func BenchmarkInsertAvChargeStatisBatch10000(b *testing.B) {
	benchmarkInsertAvChargeStatisBatch(10000, b)
}

func assembleAvChargeStatis(avChargeStatis []*model.AvChargeStatis) (vals string) {
	var buf bytes.Buffer
	for _, row := range avChargeStatis {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.IsOriginal))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteString(row.UploadTime.Time().Format(_layout))
		buf.WriteString(")")
		buf.WriteByte(',')
	}

	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}

	vals = buf.String()
	buf.Reset()
	return
}
