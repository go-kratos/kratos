package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/space/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_FavNav(t *testing.T) {
	Convey("test fav nav", t, WithService(func(s *Service) {
		vmid := int64(24598781)
		mid := int64(0)
		data, err := s.FavNav(context.Background(), mid, vmid)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestService_FavArchive(t *testing.T) {
	Convey("test fav archive", t, WithService(func(s *Service) {
		mid := int64(0)
		arg := &model.FavArcArg{
			Vmid: 908085,
			Fid:  629658,
			Pn:   1,
			Ps:   20,
		}
		data, err := s.FavArchive(context.Background(), mid, arg)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
