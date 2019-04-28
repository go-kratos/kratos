package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_NewList1(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, count, err := svf.NewList(context.Background(), 129, 1, 1, 10)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestService_NewList2(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		res, count, err := svf.NewList(context.Background(), 1, 1, 1, 10)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}
