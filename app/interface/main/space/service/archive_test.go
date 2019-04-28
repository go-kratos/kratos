package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_UpArcStat(t *testing.T) {
	Convey("test up stat", t, WithService(func(s *Service) {
		mid := int64(883968)
		data, err := s.UpArcStat(context.Background(), mid)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestService_SetTopArc(t *testing.T) {
	Convey("test set top arc", t, WithService(func(s *Service) {
		mid := int64(15555180)
		aid := int64(5464686)
		reason := "11123"
		err := s.SetTopArc(context.Background(), mid, aid, reason)
		So(err, ShouldBeNil)
	}))
}

func TestService_TopArc(t *testing.T) {
	Convey("test top arc", t, WithService(func(s *Service) {
		mid := int64(0)
		vmid := int64(15555180)
		data, err := s.TopArc(context.Background(), mid, vmid)
		So(err, ShouldBeNil)
		str, _ := json.Marshal(data)
		Println(string(str))
	}))
}

func TestService_DelTopArc(t *testing.T) {
	Convey("test del top arc", t, WithService(func(s *Service) {
		mid := int64(16913)
		err := s.DelTopArc(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}

func TestService_Masterpiece(t *testing.T) {
	Convey("test masterpiece", t, WithService(func(s *Service) {
		mid := int64(0)
		vmid := int64(15555180)
		data, err := s.Masterpiece(context.Background(), mid, vmid)
		So(err, ShouldBeNil)
		str, _ := json.Marshal(data)
		Println(string(str))
	}))
}

func TestService_EditMasterpiece(t *testing.T) {
	Convey("test masterpiece", t, WithService(func(s *Service) {
		mid := int64(15555180)
		preAid := int64(5464686)
		aid := int64(10098536)
		//aid := int64(5464827)
		reason := "test edit"
		err := s.EditMasterpiece(context.Background(), mid, preAid, aid, reason)
		So(err, ShouldBeNil)
	}))
}

func TestService_AddMasterpiece(t *testing.T) {
	Convey("test masterpiece", t, WithService(func(s *Service) {
		mid := int64(15555180)
		aid := int64(5464827)
		reason := "test add"
		err := s.AddMasterpiece(context.Background(), mid, aid, reason)
		So(err, ShouldBeNil)
	}))
}

func TestService_CancelMasterpiece(t *testing.T) {
	Convey("test masterpiece", t, WithService(func(s *Service) {
		mid := int64(15555180)
		aid := int64(5464827)
		err := s.CancelMasterpiece(context.Background(), mid, aid)
		So(err, ShouldBeNil)
	}))
}
