package model

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var c *Chall

func init() {
	c = new(Chall)
}

func TestSetState(t *testing.T) {
	convey.Convey("SetState", t, func() {
		c.DispatchState = uint32(0x3a4b5c6d)
		state := uint32(0xf)
		role := uint8(1)
		c.SetState(state, role)
		convey.So(c.DispatchState, convey.ShouldEqual, uint32(0x3a4b5cfd))
	})
}

func TestGetState(t *testing.T) {
	convey.Convey("GetState", t, func() {
		c.DispatchState = uint32(0x3a4b5c6d)
		role := uint8(2)
		result := c.getState(role)
		convey.So(result, convey.ShouldEqual, uint32(0xc))
	})
}

func TestFormatState(t *testing.T) {
	convey.Convey("FormatState", t, func() {
		c.FormatState()
	})
}

func TestFromState(t *testing.T) {
	convey.Convey("FromState", t, func() {
		c.FromState()
	})
}
