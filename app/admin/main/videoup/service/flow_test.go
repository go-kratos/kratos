package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/model/archive"
	"testing"
)

func TestService_GetFlowsByOID(t *testing.T) {
	Convey("getFlowsByOID", t, WithService(func(s *Service) {
		var (
			c   = context.TODO()
			oid = int64(2222)

			res []*archive.FlowData
			err error
		)

		res, err = s.getFlowsByOID(c, oid)
		t.Logf("res(%+v) error(%v)", res, err)
		So(err, ShouldBeNil)
	}))
}

func TestService_TxAddFlow(t *testing.T) {
	var (
		id  int64
		err error
	)
	Convey("txAddFlow", t, WithService(func(s *Service) {
		c := context.TODO()
		tx, _ := s.arc.BeginTran(c)
		id, err = s.txAddFlow(tx, archive.PoolArcForbid, 1, archive.FlowGroupNoChannel, 421, "测试-频道禁止-添加")
		tx.Commit()

		So(id, ShouldBeGreaterThan, 0)
		So(err, ShouldBeNil)
	}))
}

func TestService_TxUpFlowState(t *testing.T) {
	var (
		err error
	)

	Convey("TxUpFlowState", t, WithService(func(s *Service) {
		c := context.TODO()
		f := &archive.FlowData{ID: 540, Pool: archive.PoolArcForbid, OID: 1, GroupID: archive.FlowGroupNoChannel}
		tx, _ := s.arc.BeginTran(c)
		err = s.txUpFlowState(tx, archive.FlowDelete, 421, f)
		tx.Commit()

		So(err, ShouldBeNil)
	}))
}

func TestService_txAddOrUpdateFlowState(t *testing.T) {
	Convey("txAddOrUpdateFlowState", t, WithService(func(s *Service) {
		var (
			c     = context.TODO()
			group = archive.FlowGroupNoChannel
			pool  = archive.PoolArcForbid

			err  error
			diff string
			res  *archive.FlowData
		)

		tx, _ := s.arc.BeginTran(c)
		res, diff, err = s.txAddOrUpdateFlowState(c, tx, 16, group, 421, pool, archive.FlowDelete, "新增测试啦")
		t.Logf("res(%+v) diff(%s) error(%v)", res, diff, err)
		So(err, ShouldBeNil)
		tx.Commit()

		tx, _ = s.arc.BeginTran(c)
		res, diff, err = s.txAddOrUpdateFlowState(c, tx, 16, group, 421, pool, archive.FlowOpen, "修改state测试啦")
		t.Logf("res(%+v) diff(%s) error(%v)", res, diff, err)
		So(err, ShouldBeNil)
		tx.Commit()
	}))
}

func TestService_HitFlowGroups(t *testing.T) {
	Convey("HitFlowGroups", t, WithService(func(s *Service) {
		var (
			c   = context.TODO()
			oid = int64(2222)
			res map[string]int
			err error
		)

		res, err = s.HitFlowGroups(c, oid, []int8{})
		t.Logf("res(%+v) error(%v)", res, err)
		So(err, ShouldBeNil)

		res, err = s.HitFlowGroups(c, oid, []int8{archive.PoolUp})
		t.Logf("res(%+v) error(%v)", res, err)
		So(err, ShouldBeNil)

	}))
}
