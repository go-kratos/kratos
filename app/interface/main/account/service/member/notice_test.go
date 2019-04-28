package member

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_NoticeV2(t *testing.T) {
	Convey("TestService_NoticeV2", func() {
		time.Sleep(time.Second * 2)
		var (
			mid = int64(1)
			u   = "foo"
			//ip          = "127.0.0.1"
			pf          = "ios"
			build int64 = 123
		)
		res, err := s.NoticeV2(context.TODO(), mid, u, pf, build)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestService_CloseNoticeV2(t *testing.T) {
	Convey("TestService_CloseNoticeV2", func() {
		time.Sleep(time.Second * 2)
		var (
			mid = int64(1)
			u   = "foo"
			//ip  = "127.0.0.1"
		)
		err := s.CloseNoticeV2(context.TODO(), mid, u)
		So(err, ShouldBeNil)
	})
}
