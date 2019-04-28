package service

import (
	"context"
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_GendDiff(t *testing.T) {
	Convey("TestService_GendDiff", t, WithService(func(svf *Service) {
		generated, err := svf.GendDiff(57)
		So(err, ShouldBeNil)
		So(len(generated), ShouldBeGreaterThan, 0)
		fmt.Println(generated)
	}))
}

func TestService_Publish(t *testing.T) {
	Convey("TestService_Publish", t, WithService(func(svf *Service) {
		data, err := svf.Publish(context.Background(), 57)
		So(err, ShouldBeNil)
		fmt.Println(data)
	}))
}
