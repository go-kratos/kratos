package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/workflow/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestLastLog(t *testing.T) {
	convey.Convey("LastLog", t, func() {
		logs, err := s.LastLog(context.Background(), []int64{2038}, []int{model.WLogModuleGroup})
		convey.So(err, convey.ShouldBeNil)
		fmt.Printf("%+v", logs[2038])
	})
}
