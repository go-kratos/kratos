package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateWithdraw(t *testing.T) {
	Convey("growup-job UpdateWithdraw", t, WithService(func(s *Service) {
		newDate, oldDate := "2018-01", "2018-02"
		count := int64(0)

		err := s.UpdateWithdraw(context.Background(), oldDate, newDate, count)
		So(err, ShouldBeNil)
	}))
}
