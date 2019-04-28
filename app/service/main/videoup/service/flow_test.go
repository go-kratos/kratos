package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestService_AddByMid(t *testing.T) {
	var (
		c = context.TODO()
	)
	Convey("AddByMid", t, WithService(func(s *Service) {
		err := svr.AddByMid(c, 3, 25438, 2222, 0)
		So(err, ShouldBeNil)
	}))
}
