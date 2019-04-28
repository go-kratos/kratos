package service

import (
	"testing"

	"context"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_TagSubList(t *testing.T) {
	Convey("test sub list", t, WithService(func(s *Service) {
		mid := int64(0)
		vmid := int64(908085)
		pn := 1
		ps := 10
		data, count, err := s.TagSubList(context.Background(), mid, vmid, pn, ps)
		So(err, ShouldBeNil)
		Printf("%+v,%d", data, count)
	}))
}
