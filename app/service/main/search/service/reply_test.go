package service

import (
	"context"
	"go-common/app/service/main/search/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Reply(t *testing.T) {

	var (
		err error
		c   = context.TODO()
		sp  *model.ReplyRecordParams
	)

	Convey("Reply", t, WithService(func(s *Service) {
		_, err = s.ReplyRecord(c, sp)
		So(err, ShouldBeNil)
	}))
}
