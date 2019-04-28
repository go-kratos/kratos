package service

import (
	"context"
	"testing"

	pb "go-common/app/service/main/coin/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTodayExp(t *testing.T) {
	Convey("todayexp", t, func() {
		arg := &pb.TodayExpReq{
			Mid: 1,
		}
		_, err := s.TodayExp(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}
