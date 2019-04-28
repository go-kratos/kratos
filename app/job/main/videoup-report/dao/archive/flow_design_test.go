package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
)

func TestDao_TxAddFlow(t *testing.T) {
	Convey("TxAddFlow", t, func() {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		fid, err := d.TxAddFlow(tx, archive.PoolArcForbid, 1, 0, archive.FLowGroupIDChannel, "测试添加")
		tx.Commit()
		So(err, ShouldBeNil)
		Println(fid)
	})
}

func TestDao_TxAddFlowLog(t *testing.T) {
	Convey("TxAddFlowLog", t, func() {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		id, err := d.TxAddFlowLog(tx, archive.PoolArcForbid, archive.FlowLogAdd, 1, 0, archive.FLowGroupIDChannel, "测试添加")
		tx.Commit()
		So(err, ShouldBeNil)
		Println(id)
	})
}

func TestDao_TxUpFlowState(t *testing.T) {
	Convey("TxUpFlowState", t, func() {
		c := context.TODO()
		tx, _ := d.BeginTran(c)
		id, err := d.TxUpFlowState(tx, 551, archive.FlowDelete)
		tx.Commit()
		So(err, ShouldBeNil)
		Println(id)
	})
}

func TestDao_FlowUnique(t *testing.T) {
	Convey("FlowUnique", t, func() {
		c := context.TODO()
		f, err := d.FlowUnique(c, 1, archive.FLowGroupIDChannel, archive.PoolArcForbid)
		So(err, ShouldBeNil)
		Println(f)
	})
}
