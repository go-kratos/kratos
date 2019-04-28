package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Feedback(t *testing.T) {
	Convey("test feedback Feedback", t, WithService(func(s *Service) {
		feedParams := &model.Feedback{Aid: 5464686, Mid: 27515256, TagID: 10101, Buvid: "", Browser: "",
			Version: "1.1", Content: &model.Content{Reason: ""}, Email: "", QQ: "", Other: ""}
		err := s.Feedback(context.Background(), feedParams)
		So(err, ShouldBeNil)
	}))
}
