package dao

import (
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmpVersion    = time.Now().Format("2006_01_02_15_04_05")
	buglyBatchRun = &model.BuglyBatchRun{
		BuglyVersionID: 1,
		Version:        tmpVersion,
		BatchID:        uuid.NewV4().String(),
		RetryCount:     0,
		Status:         model.BuglyBatchRunStatusRunning,
		ErrorMsg:       "no",
	}

	queryBuglyBatchRunsRequest = &model.QueryBuglyBatchRunsRequest{
		Pagination: model.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
		Version: tmpVersion,
	}
)

func Test_Bugly_batch_run(t *testing.T) {
	Convey("test insert bugly batch run", t, func() {
		err := d.InsertBuglyBatchRun(buglyBatchRun)
		So(err, ShouldBeNil)
	})

	Convey("test update bugly batch run", t, func() {
		buglyBatchRun.Status = model.BuglyBatchRunStatusDone
		err := d.UpdateBuglyBatchRun(buglyBatchRun)
		So(err, ShouldBeNil)
	})

	Convey("test Find Bugly Batch Runs", t, func() {
		buglyBatchRun.Status = model.BuglyBatchRunStatusDone
		total, buglyBatchRuns, err := d.FindBuglyBatchRuns(queryBuglyBatchRunsRequest)
		So(err, ShouldBeNil)
		So(total, ShouldEqual, 1)
		So(buglyBatchRun.BatchID, ShouldEqual, buglyBatchRuns[0].BatchID)
	})

	Convey("test Find Last Success Batch Run By Version", t, func() {
		tmpBuglyBatchRun, err := d.QueryLastSuccessBatchRunByVersion(tmpVersion)
		So(err, ShouldBeNil)
		So(buglyBatchRun.BatchID, ShouldEqual, tmpBuglyBatchRun.BatchID)
	})
}
