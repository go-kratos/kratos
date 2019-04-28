package dao

import (
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmp5 = strconv.FormatInt(time.Now().Unix(), 10)

	tapdBugTemplate = &model.TapdBugTemplate{
		WorkspaceID:    tmp5,
		IssueFilterSQL: "SELECT * FROM bugly_issues WHERE issue_no = '265'",
		SeverityKey:    "SeverityKey" + tmp5,
		UpdateBy:       "fengyifeng",

		TapdProperty: model.TapdProperty{
			Title:            "Title" + tmp5,
			Description:      "Description" + tmp5,
			CurrentOwner:     "CurrentOwner" + tmp5,
			Platform:         "Platform" + tmp5,
			Module:           "Module" + tmp5,
			IterationID:      "IterationID" + tmp5,
			ReleaseID:        "ReleaseID" + tmp5,
			Priority:         "Priority" + tmp5,
			Severity:         "Severity" + tmp5,
			Source:           "Source" + tmp5,
			CustomFieldFour:  "CustomFieldFour" + tmp5,
			BugType:          "BugType" + tmp5,
			OriginPhase:      "OriginPhase" + tmp5,
			CustomFieldThree: "CustomFieldThree" + tmp5,
			Reporter:         "Reporter" + tmp5,
			Status:           "Status" + tmp5,
		},
	}

	queryTapdBugTemplateRequest = &model.QueryTapdBugTemplateRequest{
		Pagination: model.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
	}
)

func Test_Tapd_bug_template(t *testing.T) {
	Convey("test Insert Tapd Bug Template", t, func() {
		err := d.InsertTapdBugTemplate(tapdBugTemplate)
		So(err, ShouldBeNil)
	})

	Convey("test Update Tapd Bug Template", t, func() {
		tapdBugTemplate.UpdateBy = "xuneng"
		err := d.UpdateTapdBugTemplate(tapdBugTemplate)
		So(err, ShouldBeNil)
	})

	Convey("test Query Tapd Bug Template", t, func() {
		tmpTapdBugTemplate, err := d.QueryTapdBugTemplate(tapdBugTemplate.ID)
		So(err, ShouldBeNil)
		So(tmpTapdBugTemplate.ID, ShouldEqual, tapdBugTemplate.ID)
	})

	Convey("test Find Tapd Bug Templates", t, func() {
		total, tapdBugTemplates, err := d.FindTapdBugTemplates(queryTapdBugTemplateRequest)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(tapdBugTemplates), ShouldBeGreaterThan, 0)
	})

}
