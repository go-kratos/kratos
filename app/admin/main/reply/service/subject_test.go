package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpSubjectAttr(t *testing.T) {
	c := context.Background()
	Convey("test update subject attr", t, WithService(func(s *Service) {
		_, err := s.FreezeSub(c, 23, "asd", []int64{1}, 1, 1, "test")
		So(err, ShouldBeNil)
		sub, err := s.Subject(c, 1, 1)
		So(err, ShouldBeNil)
		So(sub.AttrVal(model.SubAttrFrozen), ShouldEqual, 1)
	}))
}
