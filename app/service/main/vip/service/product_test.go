package service

import (
	"testing"

	"go-common/app/service/main/vip/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceProductLimit(t *testing.T) {
	Convey(" TestServiceProductLimit ", t, func() {
		err := s.ProductLimit(c, &model.ArgProductLimit{
			Mid:       2,
			PanelType: "ele",
			Months:    1,
		})
		So(err, ShouldBeNil)
	})
}
