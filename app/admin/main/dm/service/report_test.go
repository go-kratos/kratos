package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeReportStat(t *testing.T) {
	var (
		c        = context.TODO()
		cidDmids = map[int64][]int64{
			9968618: {
				719923090, 719923092,
			},
		}
		state             = model.StatSecondDelete
		reason      int8  = 2
		notice      int8  = 3
		adminID     int64 = 222
		remark            = "二审删除"
		block       int64 = 3
		blockReason int64
		moral       int64 = 10
		uname             = "zzz delete"
	)
	Convey("test change report stat", t, func() {
		affect, err := svr.ChangeReportStat(c, cidDmids, state, reason, notice, adminID, block, blockReason, moral, remark, uname)
		So(err, ShouldNotBeNil)
		So(affect, ShouldBeGreaterThan, 0)
	})
}

func TestReportLog(t *testing.T) {
	var (
		c          = context.TODO()
		dmid int64 = 2
	)
	Convey("test report log", t, func() {
		lg, err := svr.ReportLog(c, dmid)
		So(err, ShouldBeNil)
		So(lg, ShouldNotBeEmpty)
	})
}

func TestReportList(t *testing.T) {
	var (
		c             = context.TODO()
		page    int64 = 1
		size    int64 = 100
		start         = "2017-05-10 00:00:00"
		end           = "2017-12-13 00:00:00"
		order         = "rp_time"
		sort          = "asc"
		keyword       = ""
		tid           = []int64{}
		rpID          = []int64{}
		state         = []int64{0, 1, 2, 3, 4, 5, 6, 7}
		upOp          = []int64{0, 1, 2}
		rt            = &model.Report{}
	)
	Convey("test report list", t, func() {
		list, err := svr.ReportList(c, page, size, start, end, order, sort, keyword, tid, rpID, state, upOp, rt)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
	})
}

func TestDMReportJudge(t *testing.T) {
	var (
		err      error
		c        = context.TODO()
		cidDmids = map[int64][]int64{
			9968618: {719923090}}
	)
	Convey("test report judge", t, func() {
		err = svr.DMReportJudge(c, cidDmids, 122, "zhang1111")
		So(err, ShouldBeNil)
	})
}

func TestJudgeResult(t *testing.T) {
	var (
		c               = context.TODO()
		cid, dmid int64 = 10109084, 719213118
	)
	Convey("test judge result", t, func() {
		err := svr.JudgeResult(c, cid, dmid, 1)
		So(err, ShouldBeNil)
	})
}
