package income

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpIncomeStatis(t *testing.T) {
	Convey("UpIncomeStatis", t, WithService(func(s *Service) {
		mids := []int64{}
		groupType := 1
		fromTime := time.Now().AddDate(0, -1, 0).Unix() * 1000
		toTime := time.Now().Unix() * 1000
		_, err := s.UpIncomeStatis(context.Background(), mids, 0, groupType, fromTime, toTime)
		So(err, ShouldBeNil)
	}))
}

func Test_GetUpIncome(t *testing.T) {
	Convey("GetUpIncome", t, WithService(func(s *Service) {
		mids := []int64{}
		fromTime := time.Now().AddDate(0, -1, 0)
		toTime := time.Now()
		query := formatUpQuery(mids, fromTime, toTime, "income")
		_, err := s.GetUpIncome(context.Background(), "up_income", "income", query)
		So(err, ShouldBeNil)
	}))
}

func BenchmarkUpIncomeStatis(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mids := []int64{}
		groupType := 1
		fromTime := time.Now().AddDate(0, -1, 0).Unix() * 1000
		toTime := time.Now().Unix() * 1000
		s.UpIncomeStatis(context.Background(), mids, 0, groupType, fromTime, toTime)
	}
}

func BenchmarkGetUpIncome(b *testing.B) {
	for n := 0; n < b.N; n++ {
		mids := []int64{}
		fromTime := time.Now().AddDate(0, -1, 0)
		toTime := time.Now()
		query := formatUpQuery(mids, fromTime, toTime, "income")
		s.GetUpIncome(context.Background(), "up_income", "income", query)
	}
}
