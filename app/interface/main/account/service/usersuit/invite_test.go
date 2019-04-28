package usersuit

import (
	"context"
	"testing"
	"time"

	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_FetchMultiInfo(t *testing.T) {
	time.Sleep(time.Second * 2)
	Convey("Fetch multi info", t, func() {
		mids := []int64{88888970}
		Convey("when not timeout", func() {
			res, err := s.fetchInfos(context.Background(), mids, "127.0.0.1", time.Second)
			So(err, ShouldBeNil)
			So(len(res), ShouldEqual, len(mids))
		})
		Convey("when timeout", func() {
			_, err := s.fetchInfos(context.Background(), mids, "127.0.0.1", time.Millisecond)
			So(err, ShouldEqual, ecode.Deadline.Error())
		})
	})
}
