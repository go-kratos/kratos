package service

import (
	"context"
	"testing"

	"go-common/app/service/main/videoup/model/archive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddByPGC(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("AddByPGC", t, WithService(func(s *Service) {
		_, err := svr.AddByPGC(c, ap)
		So(err, ShouldNotBeNil)
	}))
}

func TestService_EditByPGC(t *testing.T) {
	var (
		c  = context.TODO()
		ap = &archive.ArcParam{}
	)
	Convey("EditByPGC", t, WithService(func(s *Service) {
		err := svr.EditByPGC(c, ap)
		So(err, ShouldNotBeNil)
	}))
}
