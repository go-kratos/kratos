package charge

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/charge"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvCharge(t *testing.T) {
	Convey("Test av weekly and monthly charge\n", t, func() {
		c := context.Background()
		deleteAll(c)
		ac1 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		ac2 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvWeeklyCharge(c, ac1, 30, 10, 10)
		checkAvWeeklyCharge(c, ac2, 30, 10, 10)
		checkAvMonthlyCharge(c, ac1, 30, 10, 10)
		checkAvMonthlyCharge(c, ac2, 30, 10, 10)

		ac1 = insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10)
		ac2 = insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 110, 20)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvWeeklyCharge(c, ac1, 60, 20, 20)
		checkAvWeeklyCharge(c, ac2, 90, 30, 30)
		checkAvMonthlyCharge(c, ac1, 60, 20, 20)
		checkAvMonthlyCharge(c, ac2, 90, 30, 30)

		ac1 = insertAvDailyCharge(c, "2018-06-05", "2018-06-01 15:02:03", 1101, 110, 15)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-05", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvWeeklyCharge(c, ac1, 45, 15, 15)
		checkAvMonthlyCharge(c, ac1, 105, 35, 35)
	})
}

func Test_AvChargeDateStatis(t *testing.T) {
	Convey("Test av charge date statis\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1103, 110, 3500)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvChargeDateStatis(c, 0, 2, 1, 3520, 12, d, "av_charge_daily_statis")
		checkAvChargeDateStatis(c, 1, 1, 1, 3520, 12, d, "av_charge_daily_statis")
		weekD := getStartWeeklyDate(d)
		checkAvChargeDateStatis(c, 0, 2, 1, 3520, 12, weekD, "av_charge_weekly_statis")
		checkAvChargeDateStatis(c, 1, 1, 1, 3520, 12, weekD, "av_charge_weekly_statis")
		monthD := getStartMonthlyDate(d)
		checkAvChargeDateStatis(c, 0, 2, 1, 3520, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 1, 1, 1, 3520, 12, monthD, "av_charge_monthly_statis")

		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 1000)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 110, 1000)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1103, 110, 4000)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvChargeDateStatis(c, 0, 2, 1, 6000, 12, d, "av_charge_daily_statis")
		checkAvChargeDateStatis(c, 1, 1, 1, 6000, 12, d, "av_charge_daily_statis")
		weekD = getStartWeeklyDate(d)
		charge := 6000 + 3520
		checkAvChargeDateStatis(c, 0, 2, 1, charge, 12, weekD, "av_charge_weekly_statis")
		checkAvChargeDateStatis(c, 1, 0, 1, charge, 12, weekD, "av_charge_weekly_statis")
		checkAvChargeDateStatis(c, 2, 1, 1, charge, 12, weekD, "av_charge_weekly_statis")
		monthD = getStartMonthlyDate(d)
		checkAvChargeDateStatis(c, 0, 2, 1, charge, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 1, 0, 1, charge, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 2, 1, 1, charge, 12, monthD, "av_charge_monthly_statis")

		insertAvDailyCharge(c, "2018-06-10", "2018-06-01 15:02:03", 1101, 110, 15000)
		insertAvDailyCharge(c, "2018-06-10", "2018-06-01 15:02:03", 1102, 110, 15000)
		insertAvDailyCharge(c, "2018-06-10", "2018-06-01 15:02:03", 1103, 110, 40000)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-10", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvChargeDateStatis(c, 0, 0, 1, 70000, 12, d, "av_charge_daily_statis")
		checkAvChargeDateStatis(c, 1, 0, 1, 70000, 12, d, "av_charge_daily_statis")
		checkAvChargeDateStatis(c, 3, 2, 1, 70000, 12, d, "av_charge_daily_statis")
		checkAvChargeDateStatis(c, 4, 1, 1, 70000, 12, d, "av_charge_daily_statis")
		weekD = getStartWeeklyDate(d)
		checkAvChargeDateStatis(c, 0, 0, 1, 70000, 12, weekD, "av_charge_weekly_statis")
		checkAvChargeDateStatis(c, 1, 0, 1, 70000, 12, weekD, "av_charge_weekly_statis")
		checkAvChargeDateStatis(c, 3, 2, 1, 70000, 12, weekD, "av_charge_weekly_statis")
		checkAvChargeDateStatis(c, 4, 1, 1, 70000, 12, weekD, "av_charge_weekly_statis")
		monthD = getStartMonthlyDate(d)
		checkAvChargeDateStatis(c, 0, 0, 1, 79520, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 1, 0, 1, 79520, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 3, 2, 1, 79520, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 4, 1, 1, 79520, 12, monthD, "av_charge_monthly_statis")

		insertAvDailyCharge(c, "2018-07-10", "2018-06-01 15:02:03", 1101, 110, 100)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-10", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		monthD = getStartMonthlyDate(d)
		checkAvChargeDateStatis(c, 0, 1, 1, 100, 12, monthD, "av_charge_monthly_statis")
		checkAvChargeDateStatis(c, 1, 0, 1, 100, 12, monthD, "av_charge_monthly_statis")
	})
}

func Test_AvChargeStatis(t *testing.T) {
	Convey("Test av charge statis\n", t, func() {
		c := context.Background()
		deleteAll(c)
		ac1 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 1000)
		ac2 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 1000)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvChargeStatis(c, ac1, 1000)
		checkAvChargeStatis(c, ac2, 1000)

		ac1 = insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 1000)
		ac2 = insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 110, 1000)
		ac3 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1103, 110, 1000)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvChargeStatis(c, ac1, 2000)
		checkAvChargeStatis(c, ac2, 2000)
		checkAvChargeStatis(c, ac3, 1000)

		ac2 = insertAvDailyCharge(c, "2018-07-03", "2018-06-01 15:02:03", 1102, 110, 1000)
		ac3 = insertAvDailyCharge(c, "2018-07-03", "2018-06-01 15:02:03", 1103, 110, 1000)
		ac4 := insertAvDailyCharge(c, "2018-07-03", "2018-06-01 15:02:03", 1104, 110, 1000)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-03", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkAvChargeStatis(c, ac1, 2000)
		checkAvChargeStatis(c, ac2, 3000)
		checkAvChargeStatis(c, ac3, 2000)
		checkAvChargeStatis(c, ac4, 1000)
	})
}

// up_charge
func Test_UpCharge(t *testing.T) {
	Convey("Test up daily weekly monthly charge\n", t, func() {
		c := context.Background()
		deleteAll(c)
		ac1 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 20:02:03", 1101, 110, 10)
		ac2 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 20:02:03", 1102, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		acs := []*model.AvCharge{ac1, ac2}
		checkUpDailyCharge(c, acs, 2, 20, "up_daily_charge", "2018-06-01")
		checkUpWeeklyCharge(c, acs, 2, 20)
		checkUpMonthlyCharge(c, acs, 2, 20)

		ac11 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 20:02:03", 1101, 110, 10)
		ac21 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 20:02:03", 1102, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		acs = append(acs, []*model.AvCharge{ac11, ac21}...)
		checkUpDailyCharge(c, []*model.AvCharge{ac11, ac21}, 2, 20, "up_daily_charge", "2018-06-01")
		checkUpWeeklyCharge(c, acs, 2, 40)
		checkUpMonthlyCharge(c, acs, 2, 40)

		ac12 := insertAvDailyCharge(c, "2018-06-05", "2018-06-01 20:02:03", 1101, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-05", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		acs = append(acs, []*model.AvCharge{ac12}...)
		checkUpDailyCharge(c, []*model.AvCharge{ac12}, 1, 10, "up_daily_charge", "2018-06-05")
		checkUpWeeklyCharge(c, []*model.AvCharge{ac12}, 1, 10)
		checkUpMonthlyCharge(c, acs, 2, 50)

		ac13 := insertAvDailyCharge(c, "2018-07-05", "2018-06-01 20:02:03", 1101, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-05", time.Local)
		err = s.runVideo(c, d, nil)
		So(err, ShouldBeNil)
		checkUpDailyCharge(c, []*model.AvCharge{ac13}, 1, 10, "up_daily_charge", "2018-07-05")
		checkUpWeeklyCharge(c, []*model.AvCharge{ac13}, 1, 10)
		checkUpMonthlyCharge(c, []*model.AvCharge{ac13}, 1, 10)
	})
}

func checkAvWeeklyCharge(c context.Context, ac *model.AvCharge, totalPlayCount, incCharge, totalCharge int64) {
	avWeeklyCharge, err := s.GetAvCharge(c, startWeeklyDate, s.dao.AvWeeklyCharge)
	So(err, ShouldBeNil)
	var aw *model.AvCharge
	for _, av := range avWeeklyCharge {
		if av.AvID == ac.AvID {
			aw = av
			break
		}
	}
	So(ac.AvID, ShouldEqual, aw.AvID)
	So(ac.MID, ShouldEqual, aw.MID)
	So(ac.TagID, ShouldEqual, aw.TagID)
	So(ac.IsOriginal, ShouldEqual, aw.IsOriginal)
	So(ac.UploadTime, ShouldEqual, aw.UploadTime)
	So(totalPlayCount, ShouldEqual, aw.TotalPlayCount)
	So(totalCharge, ShouldEqual, aw.TotalCharge)
	So(incCharge, ShouldEqual, aw.IncCharge)
	So(xtime.Time(getStartWeeklyDate(ac.Date.Time()).Unix()), ShouldEqual, aw.Date)
}

func checkAvMonthlyCharge(c context.Context, ac *model.AvCharge, totalPlayCount, incCharge, totalCharge int64) {
	avCharge, err := s.GetAvCharge(c, startMonthlyDate, s.dao.AvMonthlyCharge)
	So(err, ShouldBeNil)
	var aw *model.AvCharge
	for _, av := range avCharge {
		if av.AvID == ac.AvID {
			aw = av
			break
		}
	}
	So(ac.AvID, ShouldEqual, aw.AvID)
	So(ac.MID, ShouldEqual, aw.MID)
	So(ac.TagID, ShouldEqual, aw.TagID)
	So(ac.IsOriginal, ShouldEqual, aw.IsOriginal)
	So(ac.UploadTime, ShouldEqual, aw.UploadTime)
	So(totalPlayCount, ShouldEqual, aw.TotalPlayCount)
	So(totalCharge, ShouldEqual, aw.TotalCharge)
	So(incCharge, ShouldEqual, aw.IncCharge)
	So(xtime.Time(getStartMonthlyDate(ac.Date.Time()).Unix()), ShouldEqual, aw.Date)
}

func checkAvChargeStatis(c context.Context, ac *model.AvCharge, totalCharge int64) {
	ai := &model.AvChargeStatis{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select av_id,mid,tag_id,is_original,upload_time,total_charge from av_charge_statis where av_id = %d", ac.AvID)).Scan(
		&ai.AvID, &ai.MID, &ai.TagID, &ai.IsOriginal, &ai.UploadTime, &ai.TotalCharge)
	So(err, ShouldBeNil)
	So(ac.AvID, ShouldEqual, ai.AvID)
	So(ac.MID, ShouldEqual, ai.MID)
	So(ac.TagID, ShouldEqual, ai.TagID)
	So(ac.IsOriginal, ShouldEqual, ai.IsOriginal)
	So(ac.UploadTime, ShouldEqual, ai.UploadTime)
	So(totalCharge, ShouldEqual, ai.TotalCharge)
}

func checkUpDailyCharge(c context.Context, acs []*model.AvCharge, count, total int64, table, date string) {
	if len(acs) == 0 {
		return
	}
	if date == "" {
		date = acs[0].Date.Time().Format(_layout)
	}
	var incCharge int64
	for _, ac := range acs {
		incCharge += ac.IncCharge
	}
	up := &model.UpCharge{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select av_count,inc_charge, total_charge from %s where mid = %d and date = '%s'", table, acs[0].MID, date)).Scan(
		&up.AvCount, &up.IncCharge, &up.TotalCharge)
	So(err, ShouldBeNil)
	So(total, ShouldEqual, up.TotalCharge)
	So(incCharge, ShouldEqual, up.IncCharge)
}

func checkUpWeeklyCharge(c context.Context, avs []*model.AvCharge, count, total int64) {
	m := make(map[string][]*model.AvCharge)
	for _, av := range avs {
		d := getStartWeeklyDate(av.Date.Time()).Format(_layout)
		if _, ok := m[d]; !ok {
			m[d] = make([]*model.AvCharge, 0)
		}
		m[d] = append(m[d], av)
	}
	for date, avs := range m {
		checkUpDailyCharge(c, avs, count, total, "up_weekly_charge", date)
	}
}

func checkUpMonthlyCharge(c context.Context, avs []*model.AvCharge, count, total int64) {
	m := make(map[string][]*model.AvCharge)
	for _, av := range avs {
		d := getStartMonthlyDate(av.Date.Time()).Format(_layout)
		if _, ok := m[d]; !ok {
			m[d] = make([]*model.AvCharge, 0)
		}
		m[d] = append(m[d], av)
	}
	for date, avs := range m {
		checkUpDailyCharge(c, avs, count, total, "up_monthly_charge", date)
	}
}

func checkAvChargeDateStatis(c context.Context, section, avs, categoryID, charge, count int, d time.Time, table string) {
	xd := xtime.Time(d.Unix())
	ads := &model.DateStatis{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select avs,charge,cdate from %s where money_section = %d and category_id = %d and cdate = '%s'", table, section, categoryID, d.Format(_layout))).Scan(
		&ads.Count, &ads.Charge, &ads.CDate)
	So(err, ShouldBeNil)
	So(charge, ShouldEqual, ads.Charge)
	So(xd, ShouldEqual, ads.CDate)
	So(avs, ShouldEqual, ads.Count)
	var ccount int64
	err = s.dao.QueryRow(c, fmt.Sprintf("select count(*) from %s where category_id = %d and cdate = '%s'", table, categoryID, d.Format(_layout))).Scan(&ccount)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, ccount)
}

func deleteAll(c context.Context) {
	s.dao.Exec(c, "truncate av_daily_charge_06")
	s.dao.Exec(c, "truncate av_daily_charge_07")
	s.dao.Exec(c, "truncate av_weekly_charge")
	s.dao.Exec(c, "truncate av_monthly_charge")
	s.dao.Exec(c, "truncate av_charge_statis")
	s.dao.Exec(c, "truncate up_daily_charge")
	s.dao.Exec(c, "truncate up_weekly_charge")
	s.dao.Exec(c, "truncate up_monthly_charge")
	s.dao.Exec(c, "truncate av_charge_daily_statis")
	s.dao.Exec(c, "truncate av_charge_weekly_statis")
	s.dao.Exec(c, "truncate av_charge_monthly_statis")
}

func getAvDailyChargeStruct(date string, uploadDate string, avID, mid, charge int64) *model.AvCharge {
	d, _ := time.ParseInLocation(_layout, date, time.Local)
	upD, _ := time.ParseInLocation(_layoutSec, uploadDate, time.Local)
	ac := &model.AvCharge{
		AvID:           avID,
		MID:            mid,
		TagID:          1,
		IsOriginal:     1,
		DanmakuCount:   charge,
		CommentCount:   charge,
		CollectCount:   charge,
		CoinCount:      charge,
		ShareCount:     charge,
		ElecPayCount:   charge,
		TotalPlayCount: charge * int64(3),
		WebPlayCount:   charge,
		AppPlayCount:   charge,
		H5PlayCount:    charge,
		LvUnknown:      charge,
		Lv0:            charge,
		Lv1:            charge,
		Lv2:            charge,
		Lv3:            charge,
		Lv4:            charge,
		Lv5:            charge,
		Lv6:            charge,
		VScore:         charge,
		IncCharge:      charge,
		TotalCharge:    0,
		Date:           xtime.Time(d.Unix()),
		UploadTime:     xtime.Time(upD.Unix()),
		DBState:        _dbInsert,
	}
	return ac
}

func insertAvDailyCharge(c context.Context, date string, uploadDate string, avID, mid, charge int64) *model.AvCharge {
	s.dao.Exec(c, fmt.Sprintf("insert into task_status(type,date,status) values(1,'%s',1) on duplicate key update date=values(date)", date))
	ac := getAvDailyChargeStruct(date, uploadDate, avID, mid, charge)
	_, err := s.avChargeBatchInsert(c, []*model.AvCharge{ac}, fmt.Sprintf("av_daily_charge_%s", strings.Split(date, "-")[1]))
	So(err, ShouldBeNil)
	return ac
}
