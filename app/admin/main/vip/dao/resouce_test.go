package dao

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/admin/main/vip/model"
	xsql "go-common/library/database/sql"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSelBatchCodeCount(t *testing.T) {
	convey.Convey("SelBatchCodeCount", t, func() {
		n, err := d.SelBatchCodeCount(context.TODO(), &model.ArgBatchCode{})
		convey.So(err, convey.ShouldBeNil)
		convey.So(n, convey.ShouldNotBeNil)
	})
}

func TestDaoselBatchCodeIDs(t *testing.T) {
	convey.Convey("selBatchCodeIDs", t, func() {
		ids, err := d.selBatchCodeIDs(context.TODO(), 0, 0, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ids, convey.ShouldNotBeNil)
	})
}

func TestDaocodeAutoArgSQL(t *testing.T) {
	convey.Convey("codeAutoArgSQL", t, func() {
		p1 := d.codeAutoArgSQL(&model.ArgCode{})
		convey.So(p1, convey.ShouldNotBeNil)
	})
}

func TestDaoSelCode(t *testing.T) {
	convey.Convey("SelCode", t, func() {
		_, err := d.SelCode(context.TODO(), &model.ArgCode{}, 0, 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSelBatchCodes(t *testing.T) {
	convey.Convey("SelBatchCodes", t, func() {
		_, err := d.SelBatchCodes(context.TODO(), []int64{})
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSelBatchCode(t *testing.T) {
	convey.Convey("SelBatchCode", t, func() {
		_, err := d.SelBatchCode(context.TODO(), &model.ArgBatchCode{}, 0, 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoTxAddBatchCode(t *testing.T) {
	var (
		id  int64
		err error
	)
	convey.Convey("TxAddBatchCode", t, func() {
		var tx *xsql.Tx
		tx, err = d.BeginTran(context.Background())
		convey.So(err, convey.ShouldBeNil)
		id, err = d.TxAddBatchCode(tx, &model.BatchCode{BusinessID: 1, BatchName: "ut_test"})
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldNotBeNil)
		tx.Commit()
	})
	convey.Convey("SelCodeID", t, func() {
		r, err := d.SelCodeID(context.TODO(), id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(r, convey.ShouldNotBeNil)
	})
	convey.Convey("SelBatchCodeID", t, func() {
		r, err := d.SelBatchCodeID(context.TODO(), id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(r, convey.ShouldNotBeNil)
	})
	convey.Convey("SelBatchCodeName", t, func() {
		_, err := d.SelBatchCodeName(context.TODO(), "ut_test")
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("UpdateBatchCode", t, func() {
		eff, err := d.UpdateBatchCode(context.TODO(), &model.BatchCode{ID: id, BusinessID: 11})
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
	convey.Convey("UpdateCode", t, func() {
		eff, err := d.UpdateCode(context.TODO(), id, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(eff, convey.ShouldNotBeNil)
	})
}

func TestDaoBatchAddCode(t *testing.T) {
	codes := []*model.ResourceCode{
		{BatchCodeID: int64(rand.Int31())},
	}
	convey.Convey("BatchAddCode", t, func() {
		tx, err := d.BeginTran(context.Background())
		convey.So(err, convey.ShouldBeNil)
		err = d.BatchAddCode(tx, codes)
		convey.So(err, convey.ShouldBeNil)
	})
}
