package thirdp

import (
	"fmt"
	"testing"

	"encoding/json"

	"go-common/app/interface/main/tv/dao/thirdp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_PickDBeiPage(t *testing.T) {
	Convey("TestService_PickDBeiPage", t, WithService(func(s *Service) {
		data, err := s.PickDBeiPage(0, thirdp.DBeiUGC)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		fmt.Println(data)
	}))
}

func TestService_MangoSns(t *testing.T) {
	Convey("TestService_MangoSns", t, WithService(func(s *Service) {
		data, err := s.MangoSns(ctx, 7)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		str, _ := json.Marshal(data)
		fmt.Println(string(str))
	}))
}

func TestService_MangoArcs(t *testing.T) {
	Convey("TestService_MangoArcs", t, WithService(func(s *Service) {
		data, err := s.MangoArcs(ctx, 3)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
		str, _ := json.Marshal(data)
		fmt.Println(string(str))
	}))
}
