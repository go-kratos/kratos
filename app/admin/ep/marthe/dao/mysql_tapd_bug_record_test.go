package dao

import (
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmpID = time.Now().Unix()

	tapdBugRecord = &model.TapdBugRecord{
		ProjectTemplateID: tmpID,
		VersionTemplateID: tmpID + 1,
		Operator:          "fengyifeng",
		Count:             10,
		Status:            model.InsertBugStatusRunning,
		IssueFilterSQL:    "SELECT * FROM bugly_issues WHERE issue_no = '265'",
	}
)

func Test_Tapd_bug_record(t *testing.T) {
	Convey("test Insert Tapd Bug Record", t, func() {
		err := d.InsertTapdBugRecord(tapdBugRecord)
		So(err, ShouldBeNil)
	})

	Convey("test Update Tapd Bug Record", t, func() {
		tapdBugRecord.Status = model.InsertBugStatusDone
		err := d.UpdateTapdBugRecord(tapdBugRecord)
		So(err, ShouldBeNil)
	})

	Convey("test Query Tapd Bug Record By Project ID And Status", t, func() {
		tapdBugRecords, err := d.QueryTapdBugRecordByProjectIDAndStatus(tapdBugRecord.ProjectTemplateID, model.InsertBugStatusDone)
		So(err, ShouldBeNil)
		So(len(tapdBugRecords), ShouldEqual, 1)
	})
}
