package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServicestartNet(t *testing.T) {
	var (
		rid   = int64(1)
		biz   = int64(1)
		nid   = int64(1)
		tx, _ = s.gorm.BeginTx(context.TODO())
	)
	convey.Convey("startNet", t, func(ctx convey.C) {
		result, err := s.startNet(context.TODO(), biz, nid)
		if err == nil {
			t.Logf("result(%+v) result.resulttoken(%+v)", result, result.ResultToken)
		} else {
			return
		}

		result.RID = rid
		err = s.reachNewFlowDB(context.TODO(), tx, result)
		if err == nil {
			tx.Commit()
		}

		s.afterReachNewFlow(context.TODO(), result, biz)

		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})

}

func TestServicecancelNet(t *testing.T) {
	var (
		rid   = []int64{1, 3, 4}
		tx, _ = s.gorm.BeginTx(context.TODO())
	)
	convey.Convey("cancelNet", t, func(ctx convey.C) {
		_, err := s.cancelNet(context.TODO(), tx, rid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

/**
  批量
*/
func TestServicefetchBatchOperations(t *testing.T) {
	convey.Convey("fetchBatchOperations", t, func(ctx convey.C) {
		result, err := s.fetchBatchOperations(cntx, 1, 0)
		for _, item := range result {
			t.Logf("operations(%+v)", item)
		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicecomputeBatchTriggerResult(t *testing.T) {
	var (
		rid   = int64(1)   //running at flow(1), flow(1)->t1->flow(3)
		binds = []int64{1} //bind by transition 3 & 1
	)
	convey.Convey("computeBatchTriggerResult", t, func(ctx convey.C) {
		result, err := s.computeBatchTriggerResult(cntx, 1, rid, binds)
		if err != nil {
			return
		}
		tx, _ := s.gorm.BeginTx(cntx)
		s.reachNewFlowDB(cntx, tx, result)
		err = tx.Commit().Error
		t.Logf("result(%+v)  submittoken(%+v) resulttoken(%+v)", result, result.SubmitToken, result.ResultToken)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

/**
 *资源单个处理
 */
func TestServicefetchResourceTranInfo(t *testing.T) {
	convey.Convey("fetchResourceTranInfo", t, func(ctx convey.C) {
		result, err := s.fetchResourceTranInfo(cntx, 1, 1, 0)
		t.Logf("result(%+v)", result)
		for _, item := range result.Operations {
			t.Logf("operations(%+v)", item)
		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicecomputeTriggerResult(t *testing.T) {
	var (
		rid    = int64(1)
		flowID = int64(1)
		binds  = []int64{23}
	)
	convey.Convey("computeTriggerResult", t, func(ctx convey.C) {
		result, err := s.computeTriggerResult(cntx, rid, flowID, binds)
		if err == nil {
			t.Logf("result(%+v)  submittoken(%+v) resulttoken(%+v)", result, result.SubmitToken, result.ResultToken)
			s.sendNetTriggerLog(context.TODO(), result)
		}

		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})

	})
}

/**
 * 任务单个处理
 */
func TestServicefetchTaskTransitionInfo(t *testing.T) {
	convey.Convey("fetchTaskTranInfo", t, func(ctx convey.C) {
		result, err := s.fetchTaskTranInfo(cntx, 1, 1, 1)
		t.Logf("result(%+v)", result)
		for _, item := range result.Operations {
			t.Logf("operations(%+v)", item)

		}
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

/**
 * 跳流程
 */
func TestServicejumpFlow(t *testing.T) {
	var (
		tx, _ = s.gorm.BeginTx(context.TODO())
	)
	convey.Convey("jumpFlow", t, func(ctx convey.C) {
		result, err := s.jumpFlow(cntx, tx, 1, 1, 2, []int64{10})
		if err != nil {
			return
		}

		t.Logf("result(%+v) resulttoken(%+v) submittoken(%+v)", result, result.ResultToken, result.SubmitToken)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicefetchJumpFlowInfo(t *testing.T) {
	convey.Convey("FetchJumpFlowInfo", t, func(ctx convey.C) {
		res, err := s.FetchJumpFlowInfo(cntx, 1)
		for _, item := range res.Operations {
			t.Logf("operations(%+v)", item)
		}
		for _, item := range res.Flows {
			t.Logf("flowArr(%+v)", item)
		}

		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
