package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertTransferJob(t *testing.T) {
	var (
		from, to, mid int64 = 1, 2, 3
		offset              = 1.11
	)
	Convey("insert a transfer job to mysql", t, func() {
		_, err := testDao.InsertTransferJob(context.TODO(), from, to, mid, offset, 0)
		So(err, ShouldBeNil)
	})
}

func TestTransferList(t *testing.T) {
	var cid, state int64 = 2, 3
	Convey("test transfer job list ", t, func() {
		res, _, err := testDao.TransferList(context.TODO(), cid, state, 1, 20)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestCheckTransferID(t *testing.T) {
	var (
		c        = context.TODO()
		id int64 = 265
	)
	Convey("test check trans by id", t, func() {
		_, err := testDao.CheckTransferID(c, id)
		So(err, ShouldBeNil)
	})
}

func TestSetTransferState(t *testing.T) {
	var (
		c           = context.TODO()
		id    int64 = 265
		state int8
	)
	Convey("test change transfer job state", t, func() {
		_, err := testDao.SetTransferState(c, id, state)
		So(err, ShouldBeNil)
	})
}
