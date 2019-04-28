package model

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDispatchState(t *testing.T) {
	convey.Convey("UpState", t, func() {
		dispatchState := int32(0x3a4b5c6d)
		role := int8(1)
		result, _ := DispatchState(dispatchState, role)
		convey.So(result, convey.ShouldEqual, uint32(0x6))
	})
}

func TestSetDispatchState(t *testing.T) {
	convey.Convey("UpState", t, func() {
		dispatchState := int32(0x3a4b5c6d)
		state := int8(0x1)
		role := int8(1)
		result, err := SetDispatchState(dispatchState, role, state)
		convey.ShouldBeNil(err)
		convey.So(result, convey.ShouldEqual, uint32(0x3a4b5c1d))
	})
}
