package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testVipStatus int32 = 2
	testVipMid    int64 = 110
)

//go test  -test.v -test.run TestUpdateVipStatus
func TestUpdateVipStatus(t *testing.T) {
	Convey("TestUpdateVipStatus update vip status", t, WithService(func(s *Service) {
		So(s.UpdateVipStatus(context.TODO(), testVipMid, testVipStatus), ShouldBeNil)
	}))
}
