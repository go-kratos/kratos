package service

import (
	"testing"

	"encoding/json"
	"fmt"

	"go-common/app/interface/main/tv/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_HomeRecom(t *testing.T) {
	Convey("TestService_HomeRecom", t, WithService(func(s *Service) {
		homepage, err := s.ModHome()
		So(homepage, ShouldNotBeNil)
		So(err, ShouldBeNil)
	}))
}

func TestService_PageFollow(t *testing.T) {
	Convey("TestService_PageFollow", t, WithService(func(s *Service) {
		res, err := s.PageFollow(ctx, &model.ReqPageFollow{
			AccessKey: "",
			PageID:    1,
			Build:     1011,
		})
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		for _, v := range res {
			data, _ := json.Marshal(v)
			fmt.Println(string(data))
		}
	}))
}
