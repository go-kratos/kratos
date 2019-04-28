package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDMS(t *testing.T) {
	Convey("test dm list by dmids", t, func() {
		rs, err := svr.dms(context.TODO(), 1, 10108013, []int64{719925639, 719925638})
		So(err, ShouldBeNil)
		So(rs, ShouldNotBeEmpty)
		for _, r := range rs {
			t.Logf("%+v", r)
			t.Logf("%+v", r.Content)
		}
	})
}

func TestMidHash(t *testing.T) {
	Convey("test mid hash", t, func() {
		hash, err := svr.MidHash(context.TODO(), 27515256)
		So(err, ShouldBeNil)
		So(hash, ShouldNotBeBlank, hash)
	})
}

func TestTransferJob(t *testing.T) {
	Convey("test TransferJob", t, func() {
		err := svr.TransferJob(context.TODO(), 27515615, 10108765, 10108763, 1.00)
		So(err, ShouldBeNil)
	})
}

func TestTransferList(t *testing.T) {
	var (
		cid int64 = 10109082
	)
	Convey("test transfer list", t, func() {
		l, err := svr.TransferList(c, cid)
		So(err, ShouldBeNil)
		So(l, ShouldNotBeEmpty)
	})
}

func TestTransferRetry(t *testing.T) {
	var (
		id, mid int64 = 265, 1
	)
	Convey("test transfer retry", t, func() {
		err := svr.TransferRetry(c, id, mid)
		So(err, ShouldBeNil)
	})

	Convey("test transfer retry fail", t, func() {
		err := svr.TransferRetry(c, id, mid)
		So(err, ShouldNotBeNil)
	})
}
