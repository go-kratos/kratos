package service

import (
	"context"
	"testing"
	"time"

	pb "go-common/app/service/main/history/api/grpc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_UserHide(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(100)
	)
	Convey("set data", t, func() {
		_, err := s.UpdateUserHide(c, &pb.UpdateUserHideReq{Mid: mid, Hide: true})
		So(err, ShouldBeNil)
		time.Sleep(time.Millisecond * 50)
		Convey("get data", func() {
			gotReply, err := s.UserHide(c, &pb.UserHideReq{Mid: mid})
			So(err, ShouldBeNil)
			So(gotReply.Hide, ShouldBeTrue)
		})
	})
}
