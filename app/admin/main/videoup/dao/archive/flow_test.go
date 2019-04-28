package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"go-common/app/admin/main/videoup/model/archive"
)

func TestDao_TxAddFlowLog(t *testing.T) {
	var (
		id  int64
		err error
	)
	Convey("TxAddFlowLog", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		id, err = d.TxAddFlowLog(tx, archive.PoolPrivateOrder, archive.FlowLogAdd, 10, 421, 1, "测试添加流量日志-私单-其他")
		tx.Commit()
		So(id, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)

		tx, _ = d.BeginTran(c)
		id, err = d.TxAddFlowLog(tx, archive.PoolArcForbid, archive.FlowLogAdd, 10, 421, archive.FlowGroupNoChannel, "测试添加流量日志-回查-频道禁止")
		tx.Commit()
		So(id, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	}))
}

func TestDao_TxUpFlowState(t *testing.T) {
	var (
		id, rows   int64
		err1, err2 error
	)
	Convey("TxUpFlowState", t, WithDao(func(d *Dao) {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		id, err1 = d.TxAddFlow(tx, archive.PoolArcForbid, 1, 421, archive.FlowGroupNoChannel, "测试添加频道禁止流量套餐")
		rows, err2 = d.TxUpFlowState(tx, id, archive.FlowOpen)
		tx.Commit()
		So(err1, ShouldBeNil)
		So(id, ShouldBeGreaterThan, 0)
		So(err2, ShouldBeNil)
		So(rows, ShouldEqual, 0)

		tx, _ = d.BeginTran(c)
		rows, err2 = d.TxUpFlowState(tx, id, archive.FlowDelete)
		tx.Commit()
		So(err2, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)
	}))
}

func TestDao_FlowsByOID(t *testing.T) {
	var (
		flows []*archive.FlowData
		err   error
	)
	Convey("FlowsByOID", t, WithDao(func(d *Dao) {
		c := context.TODO()
		flows, err = d.FlowsByOID(c, 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_FlowUnique(t *testing.T) {
	var (
		err error
	)
	Convey("FlowUnique", t, WithDao(func(d *Dao) {
		c := context.TODO()
		_, err = d.FlowUnique(c, 1, archive.FlowGroupNoChannel, archive.PoolArcForbid)
		So(err, ShouldBeNil)
	}))
}

func TestDao_FlowGroupPools(t *testing.T) {
	Convey("FlowGroupPools", t, WithDao(func(d *Dao) {
		c := context.TODO()
		pools, err := d.FlowGroupPools(c, []int64{23, 24, 1})
		So(err, ShouldBeNil)
		So(pools, ShouldNotBeNil)
		t.Logf("pools(%+v)", pools)
	}))
}

func TestDao_TxUpFlow(t *testing.T) {
	var (
		c = context.Background()
	)
	Convey("TxUpFlow", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(c)
		_, err := d.TxUpFlow(tx, 0, 0, 0)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func TestDao_FindGroupIDByScope(t *testing.T) {
	Convey("FindGroupIDByScope", t, WithDao(func(d *Dao) {
		c := context.TODO()
		_, err := d.FindGroupIDByScope(c, 0, 0, 0, 0)
		So(err, ShouldBeNil)
	}))
}

func TestDao_Flows(t *testing.T) {
	Convey("Flows", t, WithDao(func(d *Dao) {
		c := context.TODO()
		_, err := d.Flows(c)
		So(err, ShouldBeNil)
	}))
}

func TestDao_FlowByPool(t *testing.T) {
	Convey("FlowByPool", t, WithDao(func(d *Dao) {
		_, err := d.FlowByPool(0, 0)
		So(err, ShouldBeNil)
	}))
}
