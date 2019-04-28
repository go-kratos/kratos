package gitlab

import (
	"flag"
	"testing"

	"go-common/app/tool/saga/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	g *Gitlab
)

func init() {
	var (
		err error
	)
	g = New("http://gitlab.bilibili.co/api/v4", "z3nN4s4BVX5oNYXKbEPL")
	flag.Set("conf", "/Users/bilibili/go/src/go-common/app/tool/saga/cmd/saga-test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}

}

// Test Function HostToken
func TestHostToken(t *testing.T) {
	Convey("Test HostToken", t, func() {
		var (
			host  string
			token string
			err   error
		)
		host, token, err = g.HostToken()
		So(err, ShouldBeNil)
		So(host, ShouldEqual, "gitlab.bilibili.co")
		So(token, ShouldEqual, "z3nN4s4BVX5oNYXKbEPL")
	})
}

// TODO AutoPush
// Test MR Create/Close
// Test MRNote Create/Update/Delete
//func TestCreateMRNote(t *testing.T) {
//	Convey("Test CreateMRNote", t, func() {
//		var (
//			err    error
//			noteID int
//
//			mr *gitlab.MergeRequest
//			//res *gitlab.Response
//
//			// MR CreateMergeRequestOptions Info
//			title        = "test"
//			description  = "test"
//			sourceBranch = "chy-saga-test"
//			targetBranch = "saga-test-1"
//			assigneeID   = 15
//			projectID    = 23
//
//			//MR Info
//			state  string
//			status string
//		)
//
//		// Create CreateMergeRequestOptions Instance
//		opt := new(gitlab.CreateMergeRequestOptions)
//		opt.Title = &title
//		opt.Description = &description
//		opt.SourceBranch = &sourceBranch
//		opt.TargetBranch = &targetBranch
//		opt.AssigneeID = &assigneeID
//		opt.TargetProjectID = &projectID
//
//		// Create MergeRequest
//		mr, _, err = g.client.MergeRequests.CreateMergeRequest(35, opt)
//		So(err, ShouldBeNil)
//
//		// Query MR Status
//		state = mr.State
//		status = mr.MergeStatus
//		So(state, ShouldEqual, "opened")
//		So(status, ShouldEqual, "cannot_be_merged")
//
//		// Create MRNote
//		noteID, err = g.CreateMRNote(projectID, mr.IID, "CreateMRNote: PASS!")
//		So(err, ShouldBeNil)
//		So(noteID, ShouldNotBeNil)
//
//		// Update MRNote
//		err = g.UpdateMRNote(projectID, mr.IID, noteID, "CreateMRNote: PASS!\nUpdateMRNote: PASS!")
//		So(err, ShouldBeNil)
//
//		// Delete MRNote
//		err = g.DeleteMRNote(projectID, mr.IID, noteID)
//		So(err, ShouldBeNil)
//
//		// Accept MR
//		err = g.AcceptMR(projectID, mr.IID, "Accept MR: PASS!")
//		So(err, ShouldBeNil)
//
//		//Close MR
//		//err = g.CloseMR(projectID, mr.IID)
//		//So(err, ShouldBeNil)
//
//		// Report
//		noteID, err = g.CreateMRNote(projectID, mr.IID, "- MR Create Successfully!\n"+
//			"- MRNote Create Successfully\n"+
//			"- MRNote Update Successfully\n"+
//			"- MRNote Delete Successfully\n"+
//			"- MR Accept Successfully\n"+
//			"- MR Close Successfully")
//		So(err, ShouldBeNil)
//	})
//}

func TestProjectID(t *testing.T) {
	Convey("Test ProjectID", t, func() {
		var (
			projID int
			err    error
		)
		projID, err = g.ProjectID("git@gitlab.bilibili.co:platform/go-common.git")
		So(err, ShouldBeNil)
		So(projID, ShouldEqual, 23)
	})
}

func TestCommitDiff(t *testing.T) {
	Convey("Test CommitDiff", t, func() {
		var (
			err   error
			files []string
		)
		files, err = g.CommitDiff(23, "f5c9bfa037771b7f8179db7a245a25695c90be9b")
		So(err, ShouldBeNil)
		So(files[0], ShouldEqual, ".gitlab-ci.yml")
	})
}
