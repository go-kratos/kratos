package income

import (
	"context"
	"testing"
	"time"

	model "go-common/app/job/main/growup/model/income"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetUpIncomeTable(t *testing.T) {
	Convey("GetUpIncomeTable", t, func() {
		_, err := s.income.upIncomeSvr.GetUpIncomeTable(context.Background(), time.Now(), _upIncomeWeekly)
		So(err, ShouldBeNil)
	})
}

func Test_GetUpIncomeWeeklyAndMonthly(t *testing.T) {
	Convey("GetUpIncomeWeeklyAndMonthly", t, func() {
		_, _, err := s.income.upIncomeSvr.GetUpIncomeWeeklyAndMonthly(context.Background(), time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_UpIncomeDBStore(t *testing.T) {
	Convey("UpIncomeDBStore", t, func() {
		err := s.income.upIncomeSvr.UpIncomeDBStore(context.Background(), map[int64]*model.UpIncome{}, map[int64]*model.UpIncome{})
		So(err, ShouldBeNil)
	})
}
