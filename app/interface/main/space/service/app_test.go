package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/interface/main/space/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AppIndex(t *testing.T) {
	Convey("test app acc info", t, WithService(func(s *Service) {
		arg := &model.AppIndexArg{
			Vmid:     15555180,
			Mid:      0,
			Platform: "ios",
			Qn:       16,
			Device:   _devicePad,
		}
		data, err := s.AppIndex(context.Background(), arg)
		So(err, ShouldBeNil)
		str, _ := json.Marshal(data)
		Println(string(str))
	}))
}
