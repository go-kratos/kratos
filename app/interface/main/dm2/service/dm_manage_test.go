package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEditDMState(t *testing.T) {
	var (
		c              = context.TODO()
		mid, oid int64 = 0, 5
		state    int32 = 1
		dmids          = []int64{719149905, 719149906, 719149907}
	)
	Convey("test edit dm state", t, func() {
		err := svr.EditDMState(c, 1, mid, oid, state, dmids, 0, 0)
		So(err, ShouldBeNil)
	})
}

func TestEditDMPool(t *testing.T) {
	var (
		c              = context.TODO()
		mid, oid int64 = 0, 5
		dmids          = []int64{719149905, 719149906, 719149907}
	)
	Convey("test edit dm pool", t, func() {
		err := svr.EditDMPool(c, 1, mid, oid, 1, dmids, 0, 0)
		So(err, ShouldBeNil)
	})
}
