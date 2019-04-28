package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_BangumiList(t *testing.T) {
	Convey("test bangumi list", t, WithService(func(s *Service) {
		mid := int64(0)
		vmid := int64(883968)
		pn := 1
		ps := 10
		data, cnt, err := s.BangumiList(context.Background(), mid, vmid, pn, ps)
		So(err, ShouldBeNil)
		Printf("%+v,%d", data, cnt)
	}))
}
