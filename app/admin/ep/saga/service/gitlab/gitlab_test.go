package gitlab

import (
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/admin/ep/saga/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	g *Gitlab
)

// TestMain ...
func TestMain(m *testing.M) {
	flag.Set("conf", "../../cmd/saga-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	g = New(conf.Conf.Property.Gitlab.API, conf.Conf.Property.Gitlab.Token)
	os.Exit(m.Run())
}

// TestListProjects ...
func TestListProjects(t *testing.T) {
	Convey("ListProjects", t, func() {

		projects, err := g.ListProjects(1)
		So(err, ShouldBeNil)
		So(len(projects), ShouldBeGreaterThan, 1)
	})
}

func TestListProjectPipelines(t *testing.T) {
	Convey("listProjectPipelines", t, func() {
		_, _, err := g.ListProjectPipelines(1, 682, "")
		So(err, ShouldBeNil)
	})
}

func TestGetPipeline(t *testing.T) {
	Convey("GetPipeline", t, func() {
		_, _, err := g.GetPipeline(682, 166011)
		So(err, ShouldBeNil)
	})
}

func TestListProjectJobs(t *testing.T) {
	Convey("ListJobs", t, func() {
		jobs, resp, err := g.ListProjectJobs(5822, 1)
		So(err, ShouldBeNil)
		So(resp, ShouldNotBeNil)
		So(len(jobs), ShouldBeGreaterThan, 1)
	})
}

func TestGitlab_ListProjectMergeRequests(t *testing.T) {
	Convey("ListProjectMergeRequest", t, func() {
		var (
			project = 682
			until   = time.Now()
			since   = until.AddDate(0, -1, 0)
		)
		mrs, resp, err := g.ListProjectMergeRequests(project, &since, &until, -1)
		So(len(mrs), ShouldBeGreaterThan, 1)
		So(resp, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestGitlab_ListProjectBranch(t *testing.T) {
	Convey("ListProjectBranch", t, func() {
		var (
			project = 682
			page    = 1
		)
		branches, resp, err := g.ListProjectBranch(project, page)
		So(err, ShouldBeNil)
		So(len(branches), ShouldBeGreaterThan, 0)
		So(resp, ShouldNotBeNil)
	})
}

func TestGitlab_ListProjectCommit(t *testing.T) {
	Convey("List Project branch commit", t, func() {
		var (
			project = 682
			page    = 1
		)
		commits, resp, err := g.ListProjectCommit(project, page, nil, nil)
		So(err, ShouldBeNil)
		So(commits, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
	})
}

func TestGitlab_ListProjectRunners(t *testing.T) {
	Convey("test list project runners", t, func() {
		var (
			project = 4928
			page    = 1
		)
		runners, resp, err := g.ListProjectRunners(project, page)
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
		So(runners, ShouldNotBeNil)
	})
}
