package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_NavNum(t *testing.T) {
	Convey("test nav num", t, WithService(func(s *Service) {
		mid := int64(883968)
		vmid := int64(883968)
		data := s.NavNum(context.Background(), mid, vmid)
		Printf("%+v", data)
	}))
}

func TestService_UpStat(t *testing.T) {
	Convey("test up stat", t, WithService(func(s *Service) {
		mid := int64(883968)
		data, err := s.UpStat(context.Background(), mid)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestService_AccTags(t *testing.T) {
	Convey("test account tags", t, WithService(func(s *Service) {
		mid := int64(15555180)
		data, err := s.AccTags(context.Background(), mid)
		So(err, ShouldBeNil)
		str, _ := json.Marshal(data)
		Println(string(str))
	}))
}

func TestService_SetAccTags(t *testing.T) {
	Convey("test set account tags", t, WithService(func(s *Service) {
		tags := "a,b,c,test"
		ck := ""
		err := s.SetAccTags(context.Background(), tags, ck)
		So(err, ShouldBeNil)
	}))
}

func TestService_MyInfo(t *testing.T) {
	Convey("test my info", t, WithService(func(s *Service) {
		mid := int64(15555180)
		data, err := s.MyInfo(context.Background(), mid)
		So(err, ShouldBeNil)
		str, _ := json.Marshal(data)
		Println(string(str))
	}))
}
