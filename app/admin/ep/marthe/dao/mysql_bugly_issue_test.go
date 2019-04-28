package dao

import (
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/ep/marthe/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	tmpIssueNoStr = strconv.FormatInt(time.Now().Unix(), 10)
	tmpTapdBugID  = "bug" + tmpIssueNoStr

	buglyIssue = &model.BuglyIssue{
		IssueNo:      tmpIssueNoStr,
		Title:        "Title" + tmpIssueNoStr,
		ExceptionMsg: "ExceptionMsg" + tmpIssueNoStr,
		KeyStack:     "KeyStack" + tmpIssueNoStr,
		Detail:       "Detail" + tmpIssueNoStr,
		Tags:         "Tags" + tmpIssueNoStr,
		LastTime:     time.Now(),
		HappenTimes:  10,
		UserTimes:    20,
		Version:      "Version" + tmpIssueNoStr,
		ProjectID:    "ProjectID" + tmpIssueNoStr,
		IssueLink:    "IssueLink" + tmpIssueNoStr,
	}

	queryBuglyIssueRequest = &model.QueryBuglyIssueRequest{
		Pagination: model.Pagination{
			PageSize: 10,
			PageNum:  1,
		},
		IssueNo: buglyIssue.IssueNo,
	}
)

func Test_Bugly_Issue(t *testing.T) {
	Convey("test insert bugly issue", t, func() {
		err := d.InsertBuglyIssue(buglyIssue)
		So(err, ShouldBeNil)
	})

	Convey("test update Bugly Issue", t, func() {
		buglyIssue.ExceptionMsg = "update exception message"
		err := d.UpdateBuglyIssue(buglyIssue)
		So(err, ShouldBeNil)
	})

	Convey("test update Bugly Issue tapd bug id", t, func() {
		err := d.UpdateBuglyIssueTapdBugID(buglyIssue.ID, tmpTapdBugID)
		So(err, ShouldBeNil)
	})

	Convey("test Get Bugly Issue", t, func() {
		tmpBuglyIssue, err := d.GetBuglyIssue(buglyIssue.IssueNo, buglyIssue.Version)
		So(err, ShouldBeNil)
		So(tmpBuglyIssue.ID, ShouldEqual, buglyIssue.ID)
	})

	Convey("test Get Bugly Issues By Filter SQL", t, func() {
		sql := "select * from bugly_issues where issue_no = '" + buglyIssue.IssueNo + "'"
		tmpBuglyIssues, err := d.GetBuglyIssuesByFilterSQL(sql)
		So(err, ShouldBeNil)
		So(tmpBuglyIssues[0].IssueNo, ShouldEqual, buglyIssue.IssueNo)
	})

	Convey("test Find Bugly Issues", t, func() {
		total, tmpBuglyIssues, err := d.FindBuglyIssues(queryBuglyIssueRequest)
		So(err, ShouldBeNil)
		So(total, ShouldBeGreaterThan, 0)
		So(len(tmpBuglyIssues), ShouldBeGreaterThan, 0)
	})

}
