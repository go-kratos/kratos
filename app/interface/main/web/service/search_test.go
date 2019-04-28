package service

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/interface/main/web/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_UpRec(t *testing.T) {
	Convey("test up rec", t, WithService(func(s *Service) {
		mid := int64(908085)
		arg := &model.SearchUpRecArg{ServiceArea: "reg_ok", Platform: "h5"}
		data, err := s.UpRec(context.Background(), mid, arg, "")
		So(err, ShouldBeNil)
		str, _ := json.Marshal(data)
		Printf("%+v", string(str))
	}))
}
