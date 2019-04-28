package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/rank/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_patch(t *testing.T) {
	Convey("patch", t, func() {
		var (
			begin = "2018-01-10 03:04:05"
			end   = "2018-08-25 03:04:05"
		)
		timeBegin, _ := time.Parse(model.TimeFormat, begin)
		timeEnd, _ := time.Parse(model.TimeFormat, end)
		err := s.patch(context.Background(), timeBegin, timeEnd)
		t.Logf("err:%+v", err)
		So(err, ShouldBeNil)
	})
}
