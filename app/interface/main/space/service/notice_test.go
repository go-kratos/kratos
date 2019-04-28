package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_Notice(t *testing.T) {
	Convey("test notice", t, WithService(func(s *Service) {
		mid := int64(883968)
		data, err := s.Notice(context.Background(), mid)
		So(err, ShouldBeNil)
		Printf("%v", data)
	}))
}

func TestService_SetNotice(t *testing.T) {
	Convey("test set notice", t, WithService(func(s *Service) {
		mid := int64(883968)
		notice := ""
		err := s.SetNotice(context.Background(), mid, notice)
		So(err, ShouldBeNil)
	}))
}
