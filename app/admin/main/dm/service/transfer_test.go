package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddTransferJob(t *testing.T) {
	var (
		mid    int64 = 27515615
		from   int64 = 10108765
		to     int64 = 10108763
		offset       = 1.01
		state  int8
	)
	Convey("test TransferJob", t, func() {
		err := svr.AddTransferJob(context.TODO(), from, to, mid, offset, state)
		So(err, ShouldBeNil)
	})
}

func TestTransferList(t *testing.T) {
	Convey("test transfer list", t, func() {
		res, _, err := svr.TransferList(context.TODO(), 10109082, 3, 1, 20)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestReTransferJob(t *testing.T) {
	Convey("test transfer retry", t, func() {
		err := svr.ReTransferJob(context.TODO(), 256, 1)
		So(err, ShouldBeNil)
	})
	Convey("test transfer retry fail", t, func() {
		err := svr.ReTransferJob(context.TODO(), 256, 1)
		So(err, ShouldNotBeNil)
	})
}
