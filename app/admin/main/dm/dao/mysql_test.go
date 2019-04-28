package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUpSubjectAttr(t *testing.T) {
	Convey("test update subject attr", t, func() {
		_, err := testDao.UpSubjectAttr(context.TODO(), 1, 1221, 16)
		So(err, ShouldBeNil)
	})
}

func TestUpSubjectCount(t *testing.T) {
	Convey("update count in subject", t, func() {
		_, err := testDao.UpSubjectCount(context.TODO(), 1, 1221, 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestIncrSubjectCount(t *testing.T) {
	Convey("update count in subject", t, func() {
		_, err := testDao.IncrSubjectCount(context.TODO(), 1, 1221, 1)
		So(err, ShouldBeNil)
	})
}

func TestIncrSubMoveCount(t *testing.T) {
	Convey("update move count in subject", t, func() {
		_, err := testDao.IncrSubMoveCount(context.TODO(), 1, 1221, 1)
		So(err, ShouldBeNil)
	})
}

func TestIncrSubCount(t *testing.T) {
	Convey("update count in subject", t, func() {
		_, err := testDao.IncrSubjectCount(context.TODO(), 1, 1221, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpSubjectPool(t *testing.T) {
	Convey("update childpool in subject", t, func() {
		_, err := testDao.UpSubjectPool(context.TODO(), 1, 1221, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpSubjectState(t *testing.T) {
	Convey("update state in subject", t, func() {
		_, err := testDao.UpSubjectState(context.TODO(), 1, 1221, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpSubjectMaxlimit(t *testing.T) {
	Convey("update maxlimit in subject", t, func() {
		_, err := testDao.UpSubjectMaxlimit(context.TODO(), 1, 1221, 8000)
		So(err, ShouldBeNil)
	})
}

func TestSubject(t *testing.T) {
	var (
		tp  int32 = 1
		oid int64 = 1508
		c         = context.TODO()
	)
	Convey("subject test", t, func() {
		_, err := testDao.Subject(c, tp, oid)
		So(err, ShouldBeNil)
	})
}

func TestSubjects(t *testing.T) {
	oids := []int64{1221, 1321}
	Convey("subject test", t, func() {
		res, err := testDao.Subjects(context.TODO(), 1, oids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestChangeReportStat(t *testing.T) {
	var (
		c           = context.TODO()
		cid   int64 = 1
		dmids       = []int64{1, 2, 3, 4, 5}
		state       = model.StatFirstDelete
	)
	Convey("update dm report state and err shoule be nil", t, func() {
		err := testDao.ChangeReportStat(c, cid, dmids, state)
		So(err, ShouldBeNil)
	})
}

func TestIgnoreReport(t *testing.T) {
	var (
		cid   int64 = 1
		c           = context.TODO()
		dmids       = []int64{1, 2, 3, 4, 5}
		state int8  = 4
	)
	Convey("ignore dm report state and err shoule be nil", t, func() {
		err := testDao.IgnoreReport(c, cid, dmids, state)
		So(err, ShouldBeNil)
	})
}

func TestReports(t *testing.T) {
	var (
		cid   int64 = 10109027
		c           = context.TODO()
		dmids       = []int64{719218595}
	)
	Convey("", t, func() {
		res, err := testDao.Reports(c, cid, dmids)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		for _, r := range res {
			t.Logf("===:%+v", r)
		}
	})
}

func TestReportUsers(t *testing.T) {
	var (
		tableID int64 = 1
		c             = context.TODO()
		state         = model.NoticeUnsend
		dmids         = []int64{1, 2, 3, 4, 5}
	)
	Convey("", t, func() {
		_, err := testDao.ReportUsers(c, tableID, dmids, state)
		So(err, ShouldBeNil)
	})
}

func TestUpReportUserState(t *testing.T) {
	var (
		tableID int64 = 1
		c             = context.TODO()
		state         = model.NoticeSend
		dmids         = []int64{1, 2, 3, 4, 5}
	)
	Convey("", t, func() {
		_, err := testDao.UpReportUserState(c, tableID, dmids, state)
		So(err, ShouldBeNil)
	})
}

func TestAddReportLog(t *testing.T) {
	var (
		tableID int64 = 1
		c             = context.TODO()
		lg            = &model.ReportLog{
			ID:      1,
			Did:     1234,
			AdminID: 1,
		}
	)
	Convey("", t, func() {
		err := testDao.AddReportLog(c, tableID, []*model.ReportLog{lg})
		So(err, ShouldBeNil)
	})
}

func TestReportLog(t *testing.T) {
	var (
		dmid int64 = 719918888
		c          = context.TODO()
	)
	Convey("", t, func() {
		res, err := testDao.ReportLog(c, dmid)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
