package dao

import (
	"context"
	"go-common/app/job/main/figure/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	mid   int64 = 15555180
	score int8  = 60
)

//go test  -test.v -test.run TestPutSpyScore
func TestPutSpyScore(t *testing.T) {
	Convey("TestPutSpyScore no err", t, WithDao(func(d *Dao) {
		err := d.PutSpyScore(context.TODO(), mid, score)
		So(err, ShouldBeNil)
	}))
}

//go test  -test.v -test.run TestPutReplyAct
func TestPutReplyAct(t *testing.T) {
	Convey("TestPutReplyAct no err", t, WithDao(func(d *Dao) {
		err := d.PutReplyAct(context.TODO(), mid, model.ACColumnReplyLiked, int64(-1))
		So(err, ShouldBeNil)
	}))
}

//go test  -test.v -test.run TestPutCoinUnusual
func TestPutCoinUnusual(t *testing.T) {
	Convey("TestPutCoinUnusual no err", t, WithDao(func(d *Dao) {
		err := d.PutCoinUnusual(context.TODO(), mid, model.ACColumnLowRisk)
		So(err, ShouldBeNil)
	}))
}

//go test  -test.v -test.run TestPutCoinCount
func TestPutCoinCount(t *testing.T) {
	Convey("TestPutCoinCount no err", t, WithDao(func(d *Dao) {
		err := d.PutCoinCount(context.TODO(), mid)
		So(err, ShouldBeNil)
	}))
}

//go test  -test.v -test.run TestPayOrderInfo
func TestPayOrderInfo(t *testing.T) {
	Convey("PayOrderInfo no err", t, WithDao(func(d *Dao) {
		err := d.PayOrderInfo(context.TODO(), "", mid, 1253)
		So(err, ShouldBeNil)
	}))
}
