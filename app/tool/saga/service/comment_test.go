package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/tool/saga/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceHandleGitlabComment(t *testing.T) {
	convey.Convey("HandleGitlabComment", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			event = &model.HookComment{}
			err   error
			s     Service
		)
		err = json.Unmarshal(GitlabHookCommentTest, event)
		convey.So(err, convey.ShouldBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.HandleGitlabComment(c, event)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
