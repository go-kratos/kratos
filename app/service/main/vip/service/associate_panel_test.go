package service

import (
	"fmt"
	"go-common/app/service/main/vip/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceAssociatePanel(t *testing.T) {
	Convey(" TestServiceAssociatePanel ", t, func() {
		res, err := s.AssociatePanel(c, &model.ArgAssociatePanel{
			Mid:       1,
			Platform:  "android",
			Device:    "android",
			MobiApp:   "android_i",
			PanelType: "normal",
		})
		fmt.Println("----res:", res)
		So(err, ShouldBeNil)
	})
}
