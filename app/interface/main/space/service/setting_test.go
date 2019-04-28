package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SettingInfo(t *testing.T) {
	Convey("test setting info", t, WithService(func(s *Service) {
		mid := int64(88889018)
		data, err := s.SettingInfo(context.Background(), mid)
		So(err, ShouldBeNil)
		var str []byte
		str, err = json.Marshal(data)
		So(err, ShouldBeNil)
		Printf("%+v", string(str))
	}))
}

func TestService_IndexOrderModify(t *testing.T) {
	Convey("test index order modify", t, WithService(func(s *Service) {
		mid := int64(88889018)
		orderNum := []string{"1", "8", "7", "2", "3", "4", "5", "6", "9", "21", "22", "23", "24", "25"}
		err := s.IndexOrderModify(context.Background(), mid, orderNum)
		So(err, ShouldBeNil)
	}))
}

func TestService_PrivacyModify(t *testing.T) {
	Convey("test index order modify", t, WithService(func(s *Service) {
		mid := int64(88889018)
		field := "bangumi"
		value := 0
		err := s.PrivacyModify(context.Background(), mid, field, value)
		So(err, ShouldBeNil)
	}))
}

func TestService_fixIndexOrder(t *testing.T) {
	Convey("test fixIndexOrder", t, WithService(func(s *Service) {
		mid := int64(88889018)
		indexOrder := `["1","3","2","5","6","4",7,21,22,23,24,25]`
		s.fixIndexOrder(context.Background(), mid, indexOrder)
	}))
}
