package income

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/income"
	//xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvDailyCharge(t *testing.T) {
	Convey("AvDailyCharge", t, func() {
		_, err := d.AvDailyCharge(context.Background(), time.Now(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_AvWeeklyCharge(t *testing.T) {
	Convey("AvWeeklyCharge", t, func() {
		_, err := d.AvWeeklyCharge(context.Background(), time.Now(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_AvMonthlyCharge(t *testing.T) {
	Convey("AvMonthlyCharge", t, func() {
		_, err := d.AvMonthlyCharge(context.Background(), time.Now(), 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_InsertAvChargeTable(t *testing.T) {
	Convey("InsertAvChargeTable", t, func() {
		c := context.Background()
		d.db.Exec(c, "DELETE FROM av_weekly_charge where av_id = 11")

		avCharge := []*model.AvCharge{}
		value := &model.AvCharge{
			AvID:  11,
			MID:   11,
			TagID: 11,
		}
		avCharge = append(avCharge, value)
		vals := assembleAvCharge(avCharge)
		count, err := d.InsertAvChargeTable(c, vals, "av_weekly_charge")
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 1)

		d.db.Exec(c, "DELETE FROM av_weekly_charge where av_id = 11")
	})
}

func assembleAvCharge(avCharge []*model.AvCharge) (vals string) {
	var buf bytes.Buffer
	for _, row := range avCharge {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(row.AvID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TagID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(row.IsOriginal))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.DanmakuCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CommentCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CollectCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.CoinCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.ShareCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.ElecPayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalPlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.WebPlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.AppPlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.H5PlayCount, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.LvUnknown, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv0, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv1, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv2, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv3, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv4, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv5, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.Lv6, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.VScore, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.IncCharge, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(row.TotalCharge, 10))
		buf.WriteByte(',')
		buf.WriteByte('\'')
		buf.WriteString(row.Date.Time().Format(_layout))
		buf.WriteByte('\'')
		buf.WriteByte(',')
		buf.WriteByte('\'')
		buf.WriteString(row.UploadTime.Time().Format(_layout))
		buf.WriteByte('\'')
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
