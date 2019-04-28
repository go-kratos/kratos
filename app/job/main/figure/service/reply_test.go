package service

import (
	"context"
	"testing"

	"go-common/app/job/main/figure/model"
	repmol "go-common/app/job/main/reply/model/reply"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testReplyMid    int64 = 120
	testRefReplyMid int64 = 121
)

// go test -test.v -test.run TestPutReplyInfo
func TestPutReplyInfo(t *testing.T) {
	Convey("TestPutReplyInfo put add reply", t, WithService(func(s *Service) {
		So(s.PutReplyInfo(context.TODO(), &model.ReplyEvent{
			Mid:    testReplyMid,
			Action: model.EventAdd,
			Reply: &repmol.Reply{
				Mid: testRefReplyMid,
			},
		}), ShouldBeNil)
	}))
	Convey("TestPutReplyInfo put add reply", t, WithService(func(s *Service) {
		So(s.PutReplyInfo(context.TODO(), &model.ReplyEvent{
			Mid:    testReplyMid,
			Action: model.EventLike,
			Reply: &repmol.Reply{
				Mid: testRefReplyMid,
			},
		}), ShouldBeNil)
		So(s.PutReplyInfo(context.TODO(), &model.ReplyEvent{
			Mid:    testReplyMid,
			Action: model.EventLike,
			Reply: &repmol.Reply{
				Mid: testRefReplyMid,
			},
		}), ShouldBeNil)
	}))
	Convey("TestPutReplyInfo put add reply", t, WithService(func(s *Service) {
		So(s.PutReplyInfo(context.TODO(), &model.ReplyEvent{
			Mid:    testReplyMid,
			Action: model.EventLikeCancel,
			Reply: &repmol.Reply{
				Mid: testRefReplyMid,
			},
		}), ShouldBeNil)
	}))
	Convey("TestPutReplyInfo put hate", t, WithService(func(s *Service) {
		So(s.PutReplyInfo(context.TODO(), &model.ReplyEvent{
			Mid:    testReplyMid,
			Action: model.EventHate,
			Reply: &repmol.Reply{
				Mid: testRefReplyMid,
			},
		}), ShouldBeNil)
	}))
	Convey("TestPutReplyInfo put hate cancel", t, WithService(func(s *Service) {
		So(s.PutReplyInfo(context.TODO(), &model.ReplyEvent{
			Mid:    testReplyMid,
			Action: model.EventHateCancel,
			Reply: &repmol.Reply{
				Mid: testRefReplyMid,
			},
		}), ShouldBeNil)
	}))
}
