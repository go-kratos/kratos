package dao

import (
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmp6 = strconv.FormatInt(time.Now().Unix(), 10)

	tapdBugVersionTemplate = &model.TapdBugVersionTemplate{
		Version:           "version" + tmp6,
		ProjectTemplateID: 10,
		IssueFilterSQL:    "SELECT * FROM bugly_issues WHERE issue_no = '265'",
		SeverityKey:       "SeverityKey" + tmp6,
		UpdateBy:          "fengyifeng",

		TapdProperty: model.TapdProperty{
			Title:            "Title" + tmp6,
			Description:      "Description" + tmp6,
			CurrentOwner:     "CurrentOwner" + tmp6,
			Platform:         "Platform" + tmp6,
			Module:           "Module" + tmp6,
			IterationID:      "IterationID" + tmp6,
			ReleaseID:        "ReleaseID" + tmp6,
			Priority:         "Priority" + tmp6,
			Severity:         "Severity" + tmp6,
			Source:           "Source" + tmp6,
			CustomFieldFour:  "CustomFieldFour" + tmp6,
			BugType:          "BugType" + tmp6,
			OriginPhase:      "OriginPhase" + tmp6,
			CustomFieldThree: "CustomFieldThree" + tmp6,
			Reporter:         "Reporter" + tmp6,
			Status:           "Status" + tmp6,
		},
	}

	queryTapdBugVersionTemplateRequest = &model.QueryTapdBugVersionTemplateRequest{
		Pagination: model.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
		Version: tapdBugVersionTemplate.Version,
	}
)

func Test_Tapd_bug_version_template(t *testing.T) {
	Convey("test Insert Tapd Bug version Template", t, func() {
		err := d.InsertTapdBugVersionTemplate(tapdBugVersionTemplate)
		So(err, ShouldBeNil)
	})

	Convey("test Update Tapd Bug version Template", t, func() {
		tapdBugTemplate.UpdateBy = "xuneng"
		err := d.UpdateTapdBugVersionTemplate(tapdBugVersionTemplate)
		So(err, ShouldBeNil)
	})

	Convey("test Query Tapd Bug version Template", t, func() {
		tmpTapdBugVersionTemplate, err := d.QueryTapdBugVersionTemplate(tapdBugVersionTemplate.ID)
		So(err, ShouldBeNil)
		So(tmpTapdBugVersionTemplate.ID, ShouldEqual, tapdBugVersionTemplate.ID)
	})

	Convey("test Query Tapd Bug Version Template By version", t, func() {
		tmpTapdBugVersionTemplate, err := d.QueryTapdBugVersionTemplateByVersion(tapdBugVersionTemplate.Version)
		So(err, ShouldBeNil)
		So(tmpTapdBugVersionTemplate.ID, ShouldEqual, tapdBugVersionTemplate.ID)
	})

	Convey("test Find Tapd Bug version Templates", t, func() {
		total, tapdBugVersionTemplate, err := d.FindTapdBugVersionTemplates(queryTapdBugVersionTemplateRequest)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(tapdBugVersionTemplate), ShouldBeGreaterThan, 0)
	})
}
