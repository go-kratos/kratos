package service

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AvIncomes(t *testing.T) {
	Convey("growup-job AvIncomes", t, WithService(func(s *Service) {
		mid := int64(1)
		date := time.Now().Add(-24 * time.Hour)
		t := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		res, err := s.AvIncomes(context.Background(), mid, t.Format("2006-01-02 15:04:05"))
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
