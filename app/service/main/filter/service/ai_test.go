package service

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_FilterAiScore(t *testing.T) {
	Convey("TestDao_AiWhite", t, func() {
		err := service.FilterAiScore(ctx, "111", 1, 1, 1, 1, 1)
		fmt.Printf("err:%+v \n", err)
		So(err, ShouldBeNil)
	})
}
