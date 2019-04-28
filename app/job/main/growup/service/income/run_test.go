package income

import (
	"bytes"
	"context"
	"fmt"
	"go-common/app/job/main/growup/conf"
	"strconv"
	"strings"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/income"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Run(t *testing.T) {
	Convey("Test Run creative one day one record\n", t, func() {
		c := context.Background()
		deleteAll(c)
		date := "2018-06-01"
		dateSec := "2018-06-01 15:02:03"
		d, _ := time.ParseInLocation("2006-01-02", date, time.Local)
		ac := insertAvDailyCharge(c, date, dateSec, 1101, 110, 10)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)

		var totalIncome, taxMoney, count int64 = 10, 0, 1
		checkAllData(c, ac, totalIncome, taxMoney, count, date)
	})
}

func Test_RunWithBubble(t *testing.T) {
	var (
		c      = context.TODO()
		date   = time.Now()
		avid   = int64(2333)
		mid    = int64(233)
		charge = int64(10)
		err    error
	)
	Convey("Test Run ", t, func() {
		deleteAll(c)
		s.conf.Bubble = &conf.BubbleConfig{BRatio: []*conf.BRatio{{BType: 1, Ratio: 0.8}}}
		insertBubbleMeta(c, fmt.Sprintf("(%d,%d,'%s')", avid, 1, date.Format("2006-01-02")))
		dailyCharge := insertAvDailyCharge(c, date.Format("2006-01-02"), date.Format("2006-01-02 15:04:05"), avid, mid, charge)
		err = s.run(c, date)
		So(err, ShouldBeNil)
		err = s.runStatis(c, date)
		So(err, ShouldBeNil)
		checkAllData(c, dailyCharge, 8, 0, 1, date.Format("2006-01-02"))
	})
}

func Test_RunWithMultiBubble(t *testing.T) {
	var (
		c      = context.TODO()
		date   = time.Now()
		avid   = int64(2333)
		mid    = int64(233)
		charge = int64(10)
		err    error
	)
	Convey("Test Run ", t, func() {
		deleteAll(c)
		s.conf.Bubble = &conf.BubbleConfig{BRatio: []*conf.BRatio{{BType: 1, Ratio: 0.8}, {BType: 2, Ratio: 0.7}}}
		insertBubbleMeta(c, fmt.Sprintf("(%d,%d,'%s')", avid, 1, date.Format("2006-01-02")))
		insertBubbleMeta(c, fmt.Sprintf("(%d,%d,'%s')", avid, 2, date.Format("2006-01-02")))
		dailyCharge := insertAvDailyCharge(c, date.Format("2006-01-02"), date.Format("2006-01-02 15:04:05"), avid, mid, charge)
		err = s.run(c, date)
		So(err, ShouldBeNil)
		err = s.runStatis(c, date)
		So(err, ShouldBeNil)
		checkAllData(c, dailyCharge, 7, 0, 1, date.Format("2006-01-02"))
	})
}

func Test_AvIncome(t *testing.T) {
	Convey("Test av income\n", t, func() {
		c := context.Background()
		deleteAll(c)
		ac1 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		ac2 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome1 := checkAvIncome(c, ac1, 10, 0, 1, "2018-06-01")
		avIncome2 := checkAvIncome(c, ac2, 10, 0, 1, "2018-06-01")
		checkAvIncomeStatis(c, avIncome1, 10)
		checkAvIncomeStatis(c, avIncome2, 10)

		ac3 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 20)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome3 := checkAvIncome(c, ac3, 30, 0, 1, "2018-06-02")
		checkAvIncomeStatis(c, avIncome2, 10)
		checkAvIncomeStatis(c, avIncome3, 30)

		ac3 = insertAvDailyCharge(c, "2018-07-01", "2018-06-01 15:02:03", 1101, 110, 20)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-01", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome3 = checkAvIncome(c, ac3, 50, 0, 1, "2018-07-01")
		checkAvIncomeStatis(c, avIncome3, 50)
	})
}

func Test_AvIncomeDateStatis(t *testing.T) {
	Convey("Test av income daily weekly monthly statis\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1103, 110, 3500)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkAvIncomeDateStatis(c, 0, 2, 1, 3494, 12, d, "av_income_daily_statis")
		checkAvIncomeDateStatis(c, 4, 1, 1, 3494, 12, d, "av_income_daily_statis")
		weekD := getStartWeeklyDate(d)
		checkAvIncomeDateStatis(c, 0, 2, 1, 3494, 12, weekD, "av_income_weekly_statis")
		checkAvIncomeDateStatis(c, 4, 1, 1, 3494, 12, weekD, "av_income_weekly_statis")
		monthD := getStartMonthlyDate(d)
		checkAvIncomeDateStatis(c, 0, 2, 1, 3494, 12, monthD, "av_income_monthly_statis")
		checkAvIncomeDateStatis(c, 4, 1, 1, 3494, 12, monthD, "av_income_monthly_statis")

		// insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 1000)
		// insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 110, 1000)
		// insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1103, 110, 4000)
		// d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		// err = s.run(c, d)
		// So(err, ShouldBeNil)
		// checkAvIncomeDateStatis(c, 0, 2, 1, 5801, 12, d, "av_income_daily_statis")
		// checkAvIncomeDateStatis(c, 1, 1, 1, 5801, 12, d, "av_income_daily_statis")
		// weekD = getStartWeeklyDate(d)
		// income := 5801 + 3494
		// checkAvIncomeDateStatis(c, 0, 2, 1, income, 12, weekD, "av_income_weekly_statis")
		// checkAvIncomeDateStatis(c, 1, 0, 1, income, 12, weekD, "av_income_weekly_statis")
		// checkAvIncomeDateStatis(c, 2, 1, 1, income, 12, weekD, "av_income_weekly_statis")
		// monthD = getStartMonthlyDate(d)
		// checkAvIncomeDateStatis(c, 0, 2, 1, income, 12, monthD, "av_income_monthly_statis")
		// checkAvIncomeDateStatis(c, 1, 0, 1, income, 12, monthD, "av_income_monthly_statis")
		// checkAvIncomeDateStatis(c, 2, 1, 1, income, 12, monthD, "av_income_monthly_statis")

		// insertAvDailyCharge(c, "2018-06-10", "2018-06-01 15:02:03", 1101, 110, 15000)
		// insertAvDailyCharge(c, "2018-06-10", "2018-06-01 15:02:03", 1102, 110, 15000)
		// insertAvDailyCharge(c, "2018-06-10", "2018-06-01 15:02:03", 1103, 110, 40000)
		// d, _ = time.ParseInLocation("2006-01-02", "2018-06-10", time.Local)
		// err = s.run(c, d)
		// So(err, ShouldBeNil)
		// checkAvIncomeDateStatis(c, 0, 0, 1, 54651, 12, d, "av_income_daily_statis")
		// checkAvIncomeDateStatis(c, 1, 0, 1, 54651, 12, d, "av_income_daily_statis")
		// checkAvIncomeDateStatis(c, 3, 2, 1, 54651, 12, d, "av_income_daily_statis")
		// checkAvIncomeDateStatis(c, 4, 1, 1, 54651, 12, d, "av_income_daily_statis")
		// weekD = getStartWeeklyDate(d)
		// checkAvIncomeDateStatis(c, 0, 0, 1, 54651, 12, weekD, "av_income_weekly_statis")
		// checkAvIncomeDateStatis(c, 1, 0, 1, 54651, 12, weekD, "av_income_weekly_statis")
		// checkAvIncomeDateStatis(c, 3, 2, 1, 54651, 12, weekD, "av_income_weekly_statis")
		// checkAvIncomeDateStatis(c, 4, 1, 1, 54651, 12, weekD, "av_income_weekly_statis")
		// monthD = getStartMonthlyDate(d)
		// checkAvIncomeDateStatis(c, 0, 0, 1, 63946, 12, monthD, "av_income_monthly_statis")
		// checkAvIncomeDateStatis(c, 1, 0, 1, 63946, 12, monthD, "av_income_monthly_statis")
		// checkAvIncomeDateStatis(c, 3, 2, 1, 63946, 12, monthD, "av_income_monthly_statis")
		// checkAvIncomeDateStatis(c, 4, 1, 1, 63946, 12, monthD, "av_income_monthly_statis")

		// insertAvDailyCharge(c, "2018-07-10", "2018-06-01 15:02:03", 1101, 110, 100)
		// d, _ = time.ParseInLocation("2006-01-02", "2018-07-10", time.Local)
		// err = s.run(c, d)
		// So(err, ShouldBeNil)
		// monthD = getStartMonthlyDate(d)
		// checkAvIncomeDateStatis(c, 0, 1, 1, 100, 12, monthD, "av_income_monthly_statis")
		// checkAvIncomeDateStatis(c, 1, 0, 1, 100, 12, monthD, "av_income_monthly_statis")
	})
}

func Test_UpIncome(t *testing.T) {
	Convey("Test up income\n", t, func() {
		c := context.Background()
		deleteAll(c)
		ac1 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		ac2 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome1 := checkAvIncome(c, ac1, 10, 0, 1, "2018-06-01")
		avIncome2 := checkAvIncome(c, ac2, 10, 0, 1, "2018-06-01")
		avs := []*model.AvIncome{avIncome1, avIncome2}
		up := checkUpIncome(c, avs, 20, "up_income", "", 2)
		checkUpIncomeWeekly(c, avs, 20, 2)
		checkUpIncomeMonthly(c, avs, 20, 2)
		checkUpIncomeStatis(c, up.MID, 20)

		ac3 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10)
		ac4 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 110, 10)
		ac5 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1103, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)

		avIncome3 := checkAvIncome(c, ac3, 20, 0, 1, "2018-06-02")
		avIncome4 := checkAvIncome(c, ac4, 20, 0, 1, "2018-06-02")
		avIncome5 := checkAvIncome(c, ac5, 10, 0, 1, "2018-06-02")
		avs = append(avs, avIncome3)
		avs = append(avs, avIncome4)
		avs = append(avs, avIncome5)
		up = checkUpIncome(c, []*model.AvIncome{avIncome3, avIncome4, avIncome5}, 50, "up_income", "", 3)
		checkUpIncomeWeekly(c, avs, 50, 3)
		checkUpIncomeMonthly(c, avs, 50, 3)
		checkUpIncomeStatis(c, up.MID, 50)
		checkUpAccount(c, 110, 50, 50, 0, "2018-05")
	})
}

func Test_UpIncomeDailyStatis(t *testing.T) {
	Convey("Test up income daily statis\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1103, 111, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeDailyStatis(c, 0, 2, 30, d)

		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10000)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 111, 30000)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1103, 110, 70000)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeDailyStatis(c, 7, 1, 87300, d)
		checkUpIncomeDailyStatis(c, 8, 1, 87300, d)
	})
}

func Test_UpIncomeStatis(t *testing.T) {
	Convey("Test up income statis\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1103, 111, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeStatis(c, 110, 20)
		checkUpIncomeStatis(c, 111, 10)

		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1103, 111, 10)
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1104, 112, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeStatis(c, 110, 40)
		checkUpIncomeStatis(c, 111, 20)
		checkUpIncomeStatis(c, 112, 10)

		insertAvDailyCharge(c, "2018-07-02", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-07-02", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-07-02", "2018-06-01 15:02:03", 1103, 111, 10)
		insertAvDailyCharge(c, "2018-07-02", "2018-06-01 15:02:03", 1104, 112, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeStatis(c, 110, 60)
		checkUpIncomeStatis(c, 111, 30)
		checkUpIncomeStatis(c, 112, 20)
	})
}

func Test_AvIncomeTag(t *testing.T) {
	Convey("Test av income tag\n", t, func() {
		c := context.Background()
		deleteAll(c)
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,200,0)")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkAvIncomeTag(c, 1101, 20, 0, 20, "2018-06-01")

		s.dao.Exec(c, "delete from av_charge_ratio")
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,20,1)")
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkAvIncomeTag(c, 1101, 30, 0, 50, "2018-06-02")

		// tax
		s.dao.Exec(c, "delete from av_charge_ratio")
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,20000,1)")
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1101, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-03", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkAvIncomeTag(c, 1101, 20010, 0, 20060, "2018-06-03")
	})
}

func Test_IncomeTax(t *testing.T) {
	Convey("Test income tax\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 3000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 111, 4000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1103, 112, 8000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1104, 113, 12000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1105, 114, 18000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1106, 115, 40000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1107, 116, 80000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1108, 117, 150000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1109, 118, 250000)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 11010, 119, 500000)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkAvIncomeTag(c, 1101, 3000, 0, 3000, "2018-06-01")
		checkAvIncomeTag(c, 1102, 3950, 50, 3950, "2018-06-01")
		checkAvIncomeTag(c, 1103, 7600, 400, 7600, "2018-06-01")
		checkAvIncomeTag(c, 1104, 11100, 900, 11100, "2018-06-01")
		checkAvIncomeTag(c, 1105, 16050, 1950, 16050, "2018-06-01")
		checkAvIncomeTag(c, 1106, 33150, 6850, 33150, "2018-06-01")
		checkAvIncomeTag(c, 1107, 61650, 18350, 61650, "2018-06-01")
		checkAvIncomeTag(c, 1108, 105650, 44350, 105650, "2018-06-01")
		checkAvIncomeTag(c, 1109, 160650, 89350, 160650, "2018-06-01")
		checkAvIncomeTag(c, 11010, 265650, 234350, 265650, "2018-06-01")
	})
}

func Test_UpIncomeTag(t *testing.T) {
	Convey("Test up income tag\n", t, func() {
		c := context.Background()
		deleteAll(c)
		s.dao.Exec(c, "insert into up_charge_ratio(mid,ratio,adjust_type) values(110,200,0)")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeTag(c, 110, 20, 0, 20, "2018-06-01", "up_income")

		s.dao.Exec(c, "delete from up_charge_ratio")
		s.dao.Exec(c, "insert into up_charge_ratio(mid,ratio,adjust_type) values(110,20000,1)")
		insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		checkUpIncomeTag(c, 110, 20010, 0, 20030, "2018-06-02", "up_income")

		// av ratio + up ratio float + float
		s.dao.Exec(c, "delete from av_charge_ratio")
		s.dao.Exec(c, "delete from up_charge_ratio")
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,200,0)")
		s.dao.Exec(c, "insert into up_charge_ratio(mid,ratio,adjust_type) values(110,200,0)")
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-03", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		var income int64 = (10*2 + 10) * 2
		checkUpIncomeTag(c, 110, income, 0, 20030+income, "2018-06-03", "up_income")

		// av ratio + up ratio float + fixed
		s.dao.Exec(c, "delete from av_charge_ratio")
		s.dao.Exec(c, "delete from up_charge_ratio")
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,200,0)")
		s.dao.Exec(c, "insert into up_charge_ratio(mid,ratio,adjust_type) values(110,200,1)")
		insertAvDailyCharge(c, "2018-06-04", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-04", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-04", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		var income1 int64 = (10*2 + 10) + 200
		checkUpIncomeTag(c, 110, income1, 0, 20030+income+income1, "2018-06-04", "up_income")

		// av ratio + up ratio fixed + float
		s.dao.Exec(c, "delete from av_charge_ratio")
		s.dao.Exec(c, "delete from up_charge_ratio")
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,200,1)")
		s.dao.Exec(c, "insert into up_charge_ratio(mid,ratio,adjust_type) values(110,200,0)")
		insertAvDailyCharge(c, "2018-06-05", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-05", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-05", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		var income2 int64 = 10*2 + 200 + 10*2
		checkUpIncomeTag(c, 110, income2, 0, 20030+income+income1+income2, "2018-06-05", "up_income")

		// av ratio + up ratio fixed + fixed
		s.dao.Exec(c, "delete from av_charge_ratio")
		s.dao.Exec(c, "delete from up_charge_ratio")
		s.dao.Exec(c, "insert into av_charge_ratio(av_id,ratio,adjust_type) values(1101,200,1)")
		s.dao.Exec(c, "insert into up_charge_ratio(mid,ratio,adjust_type) values(110,200,1)")
		insertAvDailyCharge(c, "2018-06-06", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-06", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-06", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		var income3 int64 = 10 + 10 + 200 + 200
		checkUpIncomeTag(c, 110, income3, 0, 20030+income+income1+income2+income3, "2018-06-06", "up_income")
	})
}

// black_list
func Test_AvBlackList(t *testing.T) {
	Convey("Test av black list\n", t, func() {
		c := context.Background()
		deleteAll(c)
		s.dao.Exec(c, "insert into av_black_list(av_id,mid,reason) values(1101,110,1)")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		var count int
		s.dao.QueryRow(c, "select count(*) from av_income").Scan(&count)
		So(count, ShouldEqual, 0)
		s.dao.QueryRow(c, "select count(*) from up_income").Scan(&count)
		So(count, ShouldEqual, 0)
	})
}

func Test_UpInfoVideo(t *testing.T) {
	Convey("Test up_info_video\n", t, func() {
		c := context.Background()
		deleteAll(c)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")
		s.dao.Exec(c, "insert into up_info_video(mid,account_type,account_state,signed_at) values(1115,1,4,'2018-05-30')")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 20:02:03", 1101, 1115, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		var count int
		s.dao.QueryRow(c, "select count(*) from av_income").Scan(&count)
		So(count, ShouldEqual, 0)
		s.dao.QueryRow(c, "select count(*) from up_income").Scan(&count)
		So(count, ShouldEqual, 0)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")

		deleteAll(c)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")
		s.dao.Exec(c, "insert into up_info_video(mid,account_type,account_state,signed_at) values(1115,1,3,'2018-05-30')")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 20:02:03", 1101, 1115, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		s.dao.QueryRow(c, "select count(*) from av_income").Scan(&count)
		So(count, ShouldEqual, 1)
		s.dao.QueryRow(c, "select count(*) from up_income").Scan(&count)
		So(count, ShouldEqual, 1)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")

		deleteAll(c)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")
		s.dao.Exec(c, "insert into up_info_video(mid,account_type,account_state,quit_at,signed_at) values(1115,1,5,'2018-06-01','2018-05-30')")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 20:02:03", 1101, 1115, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		s.dao.QueryRow(c, "select count(*) from av_income").Scan(&count)
		So(count, ShouldEqual, 0)
		s.dao.QueryRow(c, "select count(*) from up_income").Scan(&count)
		So(count, ShouldEqual, 0)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")

		deleteAll(c)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")
		s.dao.Exec(c, "insert into up_info_video(mid,account_type,account_state,quit_at,signed_at) values(1115,1,5,'2018-06-02','2018-05-30')")
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 20:02:03", 1101, 1115, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		s.dao.QueryRow(c, "select count(*) from av_income").Scan(&count)
		So(count, ShouldEqual, 1)
		s.dao.QueryRow(c, "select count(*) from up_income").Scan(&count)
		So(count, ShouldEqual, 1)
		s.dao.Exec(c, "delete from up_info_video where mid = 1115")
	})
}

func Test_UpIncomeWeekAndMonthly(t *testing.T) {
	Convey("Test up weekly and monthly income\n", t, func() {
		c := context.Background()
		deleteAll(c)
		ac1 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		ac2 := insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 111, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome1 := checkAvIncome(c, ac1, 10, 0, 1, "2018-06-01")
		avIncome2 := checkAvIncome(c, ac2, 10, 0, 1, "2018-06-01")
		checkUpIncome(c, []*model.AvIncome{avIncome1}, 10, "up_income", "", 1)
		checkUpIncome(c, []*model.AvIncome{avIncome2}, 10, "up_income", "", 1)
		checkUpIncomeWeekly(c, []*model.AvIncome{avIncome1}, 10, 1)
		checkUpIncomeMonthly(c, []*model.AvIncome{avIncome2}, 10, 1)

		ac11 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1101, 110, 10)
		ac21 := insertAvDailyCharge(c, "2018-06-02", "2018-06-01 15:02:03", 1102, 111, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome11 := checkAvIncome(c, ac11, 20, 0, 1, "2018-06-02")
		avIncome21 := checkAvIncome(c, ac21, 20, 0, 1, "2018-06-02")
		checkUpIncome(c, []*model.AvIncome{avIncome11}, 20, "up_income", "", 1)
		checkUpIncome(c, []*model.AvIncome{avIncome21}, 20, "up_income", "", 1)
		checkUpIncomeWeekly(c, []*model.AvIncome{avIncome1, avIncome11}, 20, 1)
		checkUpIncomeMonthly(c, []*model.AvIncome{avIncome2, avIncome21}, 20, 1)

		ac12 := insertAvDailyCharge(c, "2018-06-05", "2018-06-01 15:02:03", 1101, 110, 10)
		ac22 := insertAvDailyCharge(c, "2018-06-05", "2018-06-01 15:02:03", 1102, 111, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-05", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome12 := checkAvIncome(c, ac12, 30, 0, 1, "2018-06-05")
		avIncome22 := checkAvIncome(c, ac22, 30, 0, 1, "2018-06-05")
		checkUpIncome(c, []*model.AvIncome{avIncome12}, 30, "up_income", "", 1)
		checkUpIncome(c, []*model.AvIncome{avIncome22}, 30, "up_income", "", 1)
		checkUpIncomeWeekly(c, []*model.AvIncome{avIncome12}, 30, 1)
		checkUpIncomeMonthly(c, []*model.AvIncome{avIncome2, avIncome21, avIncome22}, 30, 1)

		ac13 := insertAvDailyCharge(c, "2018-07-01", "2018-06-01 15:02:03", 1101, 110, 10)
		ac23 := insertAvDailyCharge(c, "2018-07-01", "2018-06-01 15:02:03", 1102, 111, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-01", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		err = s.runStatis(c, d)
		So(err, ShouldBeNil)
		avIncome13 := checkAvIncome(c, ac13, 40, 0, 1, "2018-07-01")
		avIncome23 := checkAvIncome(c, ac23, 40, 0, 1, "2018-07-01")
		checkUpIncome(c, []*model.AvIncome{avIncome13}, 40, "up_income", "", 1)
		checkUpIncome(c, []*model.AvIncome{avIncome23}, 40, "up_income", "", 1)
		checkUpIncomeWeekly(c, []*model.AvIncome{avIncome13}, 40, 1)
		checkUpIncomeMonthly(c, []*model.AvIncome{avIncome23}, 40, 1)
	})
}

func Test_UpAccount(t *testing.T) {
	Convey("Test up_account\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-30", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-30", "2018-06-01 15:02:03", 1102, 110, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-30", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		checkUpAccount(c, 110, 20, 20, 0, "2018-05")

		insertAvDailyCharge(c, "2018-07-01", "2018-06-01 15:02:03", 1103, 110, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-01", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		checkUpAccount(c, 110, 30, 20, 0, "2018-05")

		insertAvDailyCharge(c, "2018-07-02", "2018-06-01 15:02:03", 1101, 110, 50)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-02", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		checkUpAccount(c, 110, 80, 20, 0, "2018-05")

		s.dao.Exec(c, "update up_account set total_unwithdraw_income = 60, total_withdraw_income = 20, withdraw_date_version = '2018-06' where mid = 110")
		insertAvDailyCharge(c, "2018-07-03", "2018-06-01 15:02:03", 1101, 110, 50)
		d, _ = time.ParseInLocation("2006-01-02", "2018-07-03", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		checkUpAccount(c, 110, 130, 110, 20, "2018-06")
	})
}

func Test_MuchUpAccount(t *testing.T) {
	Convey("Test much up account\n", t, func() {
		c := context.Background()
		deleteAll(c)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1102, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1103, 110, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1104, 111, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1105, 111, 10)
		insertAvDailyCharge(c, "2018-06-01", "2018-06-01 15:02:03", 1106, 112, 10)
		d, _ := time.ParseInLocation("2006-01-02", "2018-06-01", time.Local)
		err := s.run(c, d)
		So(err, ShouldBeNil)
		checkUpAccount(c, 110, 30, 30, 0, "2018-05")
		checkUpAccount(c, 111, 20, 20, 0, "2018-05")
		checkUpAccount(c, 112, 10, 10, 0, "2018-05")

		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1101, 110, 10)
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1103, 110, 10)
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1104, 111, 10)
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1105, 111, 10)
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1106, 112, 10)
		insertAvDailyCharge(c, "2018-06-03", "2018-06-01 15:02:03", 1107, 112, 10)
		d, _ = time.ParseInLocation("2006-01-02", "2018-06-03", time.Local)
		err = s.run(c, d)
		So(err, ShouldBeNil)
		checkUpAccount(c, 110, 50, 50, 0, "2018-05")
		checkUpAccount(c, 111, 40, 40, 0, "2018-05")
		checkUpAccount(c, 112, 30, 30, 0, "2018-05")
	})
}

func checkAllData(c context.Context, ac *model.AvCharge, totalIncome, taxMoney, count int64, date string) {
	d, _ := time.ParseInLocation(_layout, date, time.Local)
	// av
	avIncome := checkAvIncome(c, ac, totalIncome, taxMoney, count, date)
	checkAvIncomeStatis(c, avIncome, totalIncome) // Background
	checkAvIncomeDateStatis(c, int(0), int(1), int(ac.TagID), int(totalIncome), int(12), d, "av_income_daily_statis")
	weekD := getStartWeeklyDate(d)
	checkAvIncomeDateStatis(c, int(0), int(1), int(ac.TagID), int(totalIncome), int(12), weekD, "av_income_weekly_statis")
	monthD := getStartMonthlyDate(d)
	checkAvIncomeDateStatis(c, int(0), int(1), int(ac.TagID), int(totalIncome), int(12), monthD, "av_income_monthly_statis")

	// up
	up := checkUpIncome(c, []*model.AvIncome{avIncome}, avIncome.Income, "up_income", "", 1)
	checkUpIncomeWeekly(c, []*model.AvIncome{avIncome}, avIncome.Income, 1)
	checkUpIncomeMonthly(c, []*model.AvIncome{avIncome}, avIncome.Income, 1)
	checkUpIncomeStatis(c, up.MID, avIncome.Income)
	checkUpIncomeDailyStatis(c, 0, 1, avIncome.Income, d)
	checkUpAccount(c, up.MID, avIncome.Income, avIncome.Income, 0, "2018-05")
}

func checkAvIncome(c context.Context, ac *model.AvCharge, totalIncome, taxMoney, count int64, date string) *model.AvIncome {
	ai := &model.AvIncome{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select av_id,mid,tag_id,is_original,upload_time,play_count,total_income,income,tax_money,date from av_income where av_id = %d and date = '%s'", ac.AvID, date)).Scan(
		&ai.AvID, &ai.MID, &ai.TagID, &ai.IsOriginal, &ai.UploadTime, &ai.PlayCount, &ai.TotalIncome, &ai.Income, &ai.TaxMoney, &ai.Date)
	So(err, ShouldBeNil)
	So(ac.AvID, ShouldEqual, ai.AvID)
	So(ac.MID, ShouldEqual, ai.MID)
	So(ac.TagID, ShouldEqual, ai.TagID)
	So(ac.IsOriginal, ShouldEqual, ai.IsOriginal)
	So(ac.UploadTime, ShouldEqual, ai.UploadTime)
	So(ac.TotalPlayCount, ShouldEqual, ai.PlayCount)
	So(totalIncome, ShouldEqual, ai.TotalIncome)
	So(taxMoney, ShouldEqual, ai.TaxMoney)
	So(ac.IncCharge, ShouldEqual, ai.Income)
	So(ac.Date, ShouldEqual, ai.Date)
	return ai
}

func checkAvIncomeTag(c context.Context, avID, income, taxMoney, totalIncome int64, date string) {
	ai := &model.AvIncome{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select total_income,income,tax_money from av_income where av_id = %d and date = '%s'", avID, date)).Scan(
		&ai.TotalIncome, &ai.Income, &ai.TaxMoney)
	So(err, ShouldBeNil)
	So(totalIncome, ShouldEqual, ai.TotalIncome)
	So(taxMoney, ShouldEqual, ai.TaxMoney)
	So(income, ShouldEqual, ai.Income)
}

func checkAvIncomeStatis(c context.Context, ac *model.AvIncome, totalIncome int64) {
	ai := &model.AvIncomeStat{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select av_id,mid,tag_id,is_original,upload_time,total_income from av_income_statis where av_id = %d", ac.AvID)).Scan(
		&ai.AvID, &ai.MID, &ai.TagID, &ai.IsOriginal, &ai.UploadTime, &ai.TotalIncome)
	So(err, ShouldBeNil)
	So(ac.AvID, ShouldEqual, ai.AvID)
	So(ac.MID, ShouldEqual, ai.MID)
	So(ac.TagID, ShouldEqual, ai.TagID)
	So(ac.IsOriginal, ShouldEqual, ai.IsOriginal)
	So(ac.UploadTime, ShouldEqual, ai.UploadTime)
	So(ac.TotalIncome, ShouldEqual, ai.TotalIncome)
}

func checkAvIncomeDateStatis(c context.Context, section, avs, categoryID, income, count int, d time.Time, table string) {
	xd := xtime.Time(d.Unix())
	ads := &model.DateStatis{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select avs,income,cdate from %s where money_section = %d and category_id = %d and cdate = '%s'", table, section, categoryID, d.Format(_layout))).Scan(
		&ads.Count, &ads.Income, &ads.CDate)
	So(err, ShouldBeNil)
	So(income, ShouldEqual, ads.Income)
	So(xd, ShouldEqual, ads.CDate)
	So(avs, ShouldEqual, ads.Count)
	var ccount int64
	err = s.dao.QueryRow(c, fmt.Sprintf("select count(*) from %s where category_id = %d and cdate = '%s'", table, categoryID, d.Format(_layout))).Scan(&ccount)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, ccount)
}

func checkUpIncome(c context.Context, avs []*model.AvIncome, totalIncome int64, table string, date string, count int64) (up *model.UpIncome) {
	if len(avs) == 0 {
		return
	}
	mid := avs[0].MID
	if date == "" {
		date = avs[0].Date.Time().Format(_layout)
	}
	var playCount, avIncome, taxMoney int64
	for _, av := range avs {
		playCount += av.PlayCount
		avIncome += av.Income
		taxMoney += av.TaxMoney
	}

	up = &model.UpIncome{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select mid,av_count,play_count,av_income,audio_income,column_income,tax_money,income,total_income from %s where mid = %d and date = '%s'", table, mid, date)).Scan(
		&up.MID, &up.AvCount, &up.PlayCount, &up.AvIncome, &up.AudioIncome, &up.ColumnIncome, &up.TaxMoney, &up.Income, &up.TotalIncome)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, up.AvCount)
	// So(playCount, ShouldEqual, up.PlayCount)
	So(avIncome, ShouldEqual, up.AvIncome)
	So(0, ShouldEqual, up.AudioIncome)
	So(0, ShouldEqual, up.ColumnIncome)
	So(taxMoney, ShouldEqual, up.TaxMoney)
	So(totalIncome, ShouldEqual, up.TotalIncome)
	return up
}

func checkUpIncomeTag(c context.Context, mid, income, taxMoney, totalIncome int64, date, table string) {
	up := &model.UpIncome{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select tax_money,income,total_income from %s where mid = %d and date = '%s'", table, mid, date)).Scan(
		&up.TaxMoney, &up.Income, &up.TotalIncome)
	So(err, ShouldBeNil)
	So(income, ShouldEqual, up.Income)
	So(taxMoney, ShouldEqual, up.TaxMoney)
	So(totalIncome, ShouldEqual, up.TotalIncome)
}

func checkUpIncomeWeekly(c context.Context, avs []*model.AvIncome, totalIncome, count int64) {
	m := make(map[string][]*model.AvIncome)
	for _, av := range avs {
		d := getStartWeeklyDate(av.Date.Time()).Format(_layout)
		if _, ok := m[d]; !ok {
			m[d] = make([]*model.AvIncome, 0)
		}
		m[d] = append(m[d], av)
	}
	for date, avs := range m {
		checkUpIncome(c, avs, totalIncome, "up_income_weekly", date, count)
	}
}

func checkUpIncomeMonthly(c context.Context, avs []*model.AvIncome, totalIncome, count int64) {
	m := make(map[string][]*model.AvIncome)
	for _, av := range avs {
		d := getStartMonthlyDate(av.Date.Time()).Format(_layout)
		if _, ok := m[d]; !ok {
			m[d] = make([]*model.AvIncome, 0)
		}
		m[d] = append(m[d], av)
	}
	for date, avs := range m {
		checkUpIncome(c, avs, totalIncome, "up_income_monthly", date, count)
	}
}

func checkUpIncomeStatis(c context.Context, mid, totalIncome int64) {
	var income int64
	err := s.dao.QueryRow(c, fmt.Sprintf("select total_income from up_income_statis where mid = %d", mid)).Scan(&income)
	So(err, ShouldBeNil)
	So(totalIncome, ShouldEqual, income)
}

func checkUpIncomeDailyStatis(c context.Context, section, ups, income int64, d time.Time) {
	xd := xtime.Time(d.Unix())
	ads := &model.DateStatis{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select ups,income,cdate from up_income_daily_statis where money_section = %d and cdate = '%s'", section, d.Format(_layout))).Scan(
		&ads.Count, &ads.Income, &ads.CDate)
	So(err, ShouldBeNil)
	So(income, ShouldEqual, ads.Income)
	So(xd, ShouldEqual, ads.CDate)
	So(ups, ShouldEqual, ads.Count)
	var ccount int64
	err = s.dao.QueryRow(c, fmt.Sprintf("select count(*) from up_income_daily_statis where cdate = '%s'", d.Format(_layout))).Scan(&ccount)
	So(err, ShouldBeNil)
	So(12, ShouldEqual, ccount)
}

func checkUpAccount(c context.Context, mid, total, unwithdraw, withdraw int64, dateVersion string) {
	up := &model.UpAccount{}
	err := s.dao.QueryRow(c, fmt.Sprintf("select total_income,total_unwithdraw_income,total_withdraw_income,withdraw_date_version from up_account where mid = %d", mid)).Scan(
		&up.TotalIncome, &up.TotalUnwithdrawIncome, &up.TotalWithdrawIncome, &up.WithdrawDateVersion)
	So(err, ShouldBeNil)
	So(up.TotalIncome, ShouldEqual, total)
	So(up.TotalUnwithdrawIncome, ShouldEqual, unwithdraw)
	So(up.TotalWithdrawIncome, ShouldEqual, withdraw)
	So(up.WithdrawDateVersion, ShouldEqual, dateVersion)
}

// func prepareTest(c context.Context) {
// 	s.dao.Exec(c, "delete from av_daily_charge_06")
// 	s.dao.Exec(c, "delete from av_daily_charge_07")
// 	s.dao.Exec(c, "delete from av_income where av_id = 1101")
// 	s.dao.Exec(c, "delete from av_income_statis where av_id = 1101")
// 	s.dao.Exec(c, "delete from av_income_daily_statis")
// 	s.dao.Exec(c, "delete from av_income_weekly_statis")
// 	s.dao.Exec(c, "delete from av_income_monthly_statis")
// 	s.dao.Exec(c, "delete from av_weekly_charge where av_id = 1101")
// 	s.dao.Exec(c, "delete from av_monthly_charge where av_id = 1101")
// 	s.dao.Exec(c, "delete from av_charge_statis where av_id = 1101")
// 	s.dao.Exec(c, "delete from up_income where mid = 110")
// 	s.dao.Exec(c, "delete from up_income_weekly where mid = 110")
// 	s.dao.Exec(c, "delete from up_income_monthly where mid = 110")
// 	s.dao.Exec(c, "delete from up_income_statis where mid = 110")
// 	s.dao.Exec(c, "delete from up_income_daily_statis")
// 	s.dao.Exec(c, "delete from up_account where mid = 110")
// 	s.dao.Exec(c, "delete from up_daily_charge where mid = 110")
// 	s.dao.Exec(c, "delete from up_weekly_charge where mid = 110")
// 	s.dao.Exec(c, "delete from up_monthly_charge where mid = 110")
// 	s.dao.Exec(c, "delete from up_av_statis where mid = 110")
// }

func deleteAll(c context.Context) {
	// s.dao.Exec(c, "truncate up_info_video")
	// s.dao.Exec(c, "truncate av_daily_charge_07")
	// s.dao.Exec(c, "truncate av_daily_charge_06")
	s.dao.Exec(c, "truncate av_black_list")
	s.dao.Exec(c, "truncate av_charge_ratio")
	s.dao.Exec(c, "truncate up_charge_ratio")
	s.dao.Exec(c, "truncate av_income")
	s.dao.Exec(c, "truncate av_income_statis")
	s.dao.Exec(c, "truncate av_income_daily_statis")
	s.dao.Exec(c, "truncate av_income_weekly_statis")
	s.dao.Exec(c, "truncate av_income_monthly_statis")
	s.dao.Exec(c, "truncate up_income")
	s.dao.Exec(c, "truncate up_income_weekly")
	s.dao.Exec(c, "truncate up_income_monthly")
	s.dao.Exec(c, "truncate up_income_statis")
	s.dao.Exec(c, "truncate up_income_daily_statis")
	s.dao.Exec(c, "truncate up_account")
	s.dao.Exec(c, "truncate up_av_daily_statis")
	s.dao.Exec(c, "truncate up_column_daily_statis")
	s.dao.Exec(c, "truncate column_income_daily_statis")
	s.dao.Exec(c, "truncate column_income_weekly_statis")
	s.dao.Exec(c, "truncate column_income_monthly_statis")

	// s.dao.Exec(c, "truncate up_info_column")
	s.dao.Exec(c, "truncate column_income")
	s.dao.Exec(c, "truncate column_income_statis")
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

func getColumnDailyChargeStruct(date string, uploadDate string, aid, mid, charge int64) *model.ColumnCharge {
	d, _ := time.ParseInLocation(_layout, date, time.Local)
	upD, _ := time.ParseInLocation(_layoutSec, uploadDate, time.Local)
	cc := &model.ColumnCharge{
		ArticleID:    aid,
		Title:        "sssss",
		MID:          mid,
		TagID:        1,
		IncCharge:    charge,
		IncViewCount: 10,
		Date:         xtime.Time(d.Unix()),
		UploadTime:   xtime.Time(upD.Unix()),
	}
	return cc
}

func insertAvDailyCharge(c context.Context, date string, uploadDate string, avID, mid, charge int64) *model.AvCharge {
	s.dao.Exec(c, fmt.Sprintf("insert into task_status(type,date,status) values(1,'%s',1) on duplicate key update date=values(date)", date))
	ac := getAvDailyChargeStruct(date, uploadDate, avID, mid, charge)
	_, err := s.avCharge.avChargeBatchInsert(c, []*model.AvCharge{ac}, fmt.Sprintf("av_daily_charge_%s", strings.Split(date, "-")[1]))
	So(err, ShouldBeNil)
	return ac
}

func insertBubbleMeta(c context.Context, values string) {
	s.dao.Exec(c, fmt.Sprintf("insert into lottery_av_info(av_id,b_type,date) values %s on duplicate key update b_type=values(b_type) date=values(date)", values))
}

func insertAvDailyChargeBatch(c context.Context, date string, acs []*model.AvCharge) {
	s.avCharge.avChargeBatchInsert(c, acs, fmt.Sprintf("av_daily_charge_%s", strings.Split(date, "-")[1]))
}

func insertColumnDailyChargeBatch(c context.Context, date string, acs []*model.ColumnCharge) {
	s.avCharge.columnChargeBatchInsert(c, acs, "column_daily_charge")
}

func batchInsertUpInfoVideo(c context.Context, mids []int64) {
	var buf bytes.Buffer
	for _, mid := range mids {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteByte(',')
		buf.WriteString("1")
		buf.WriteByte(',')
		buf.WriteString("3")
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values := buf.String()
	buf.Reset()
	s.dao.Exec(c, "delete from up_info_video")
	s.dao.Exec(c, "delete from up_info_column")
	s.dao.Exec(c, fmt.Sprintf("insert into up_info_video(mid,account_type,account_state) values %s", values))
	s.dao.Exec(c, fmt.Sprintf("insert into up_info_column(mid,account_type,account_state) values %s", values))
}

// test 2018-09
func BenchmarkRunOneDay(b *testing.B) {
	c := context.Background()
	deleteAll(c)
	date := "2018-09-02"
	s.dao.Exec(c, "truncate av_daily_charge_09")
	s.dao.Exec(c, "truncate column_daily_charge")
	uploadDate := "2018-09-01 21:20:30"
	acs := make([]*model.AvCharge, 20000)
	ccs := make([]*model.ColumnCharge, 20000)
	j := 0
	for i := 1; i <= 1000000; i++ {
		if j >= 20000 {
			insertAvDailyChargeBatch(c, date, acs)
			insertColumnDailyChargeBatch(c, date, ccs)
			j = 0
		}
		acs[j] = getAvDailyChargeStruct(date, uploadDate, int64(i), int64(j+1), 100)
		ccs[j] = getColumnDailyChargeStruct(date, uploadDate, int64(i), int64(j+1), 100)
		j++
	}
	if j > 0 {
		insertAvDailyChargeBatch(c, date, acs[:j])
		insertColumnDailyChargeBatch(c, date, ccs[:j])
	}

	mids := make([]int64, 20000)
	for i := 1; i <= 20000; i++ {
		mids[i-1] = int64(i)
	}
	batchInsertUpInfoVideo(c, mids)
	fmt.Println("start run...........")
	d, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	for n := 0; n < b.N; n++ {
		s.run(c, d)
	}
}
