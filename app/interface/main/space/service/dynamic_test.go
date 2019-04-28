package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/interface/main/space/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_DynamicList(t *testing.T) {
	convey.Convey("test dynamic list", t, WithService(func(s *Service) {
		arg := &model.DyListArg{
			Vmid: 908085,
			Pn:   1,
			Qn:   16,
		}
		list, err := s.DynamicList(context.Background(), arg)
		convey.So(err, convey.ShouldBeNil)
		bs, _ := json.Marshal(list)
		convey.Println(string(bs))
	}))
}
