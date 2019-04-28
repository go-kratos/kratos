package service

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/xanzy/go-gitlab"
)

var (
	srv *Service
)

func TestMain(m *testing.M) {
	var err error
	flag.Set("conf", "../cmd/saga-admin-test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}

	srv = New()
	os.Exit(m.Run())
}

// TestInsertDB ...
func TestInsertDB(t *testing.T) {

	Convey("insertDB", t, func() {
		p := &gitlab.Project{
			ID:            11,
			Name:          "test",
			Description:   "[主站 android java] test",
			WebURL:        "http://gitlab.bilibili.co/platform/go-common",
			SSHURLToRepo:  "git@gitlab.bilibili.co:platform/go-common.git",
			DefaultBranch: "master",
			Namespace: &gitlab.ProjectNamespace{
				Name: "mytest",
				Kind: "group",
			},
		}

		err := srv.insertDB(p)
		So(err, ShouldBeNil)

		projectInfo, err := srv.dao.ProjectInfoByID(p.ID)
		So(err, ShouldBeNil)
		So(projectInfo.ProjectID, ShouldEqual, 11)
		So(projectInfo.Name, ShouldEqual, "test")
		So(projectInfo.WebURL, ShouldEqual, "http://gitlab.bilibili.co/platform/go-common")
		So(projectInfo.Repo, ShouldEqual, "git@gitlab.bilibili.co:platform/go-common.git")
		So(projectInfo.DefaultBranch, ShouldEqual, "master")
		So(projectInfo.Department, ShouldEqual, "主站")
		So(projectInfo.Business, ShouldEqual, "android")
		So(projectInfo.Language, ShouldEqual, "java")
		So(projectInfo.SpaceName, ShouldEqual, "mytest")
		So(projectInfo.SpaceKind, ShouldEqual, "group")
	})

}

// TestResolveDes ...
func TestResolveDes(t *testing.T) {

	Convey("resolveDes", t, func() {

		s := "[主站 android java] test"
		department, business, language, parseFail := parseDes(s)
		So(department, ShouldEqual, "主站")
		So(business, ShouldEqual, "android")
		So(language, ShouldEqual, "java")
		So(parseFail, ShouldEqual, false)
	})

}

// TestCollectProject ...
func TestCollectProject(t *testing.T) {

	Convey("CollectProject", t, func() {

		err := srv.CollectProject(context.Background())
		So(err, ShouldBeNil)
	})

}

// TestPushMsg ...
func TestPushMsg(t *testing.T) {

	Convey("TestPushMsg", t, func() {

		var err error

		result := &model.SyncResult{}

		for i := 0; i < 10; i++ {
			errData := &model.FailData{
				ChildID: i,
			}
			result.FailData = append(result.FailData, errData)
		}

		err = srv.WechatFailData(model.DataTypeJob, 888, result, nil)
		So(err, ShouldBeNil)
	})
}

func TestServiceSaveAggregateBranchDatabase(t *testing.T) {
	Convey("test service aggregate branch include special char", t, func(ctx C) {
		var (
			branch = &model.AggregateBranches{
				ID:               4,
				ProjectID:        666,
				ProjectName:      "六六大顺",
				BranchName:       "666",
				BranchUserName:   "吴维",
				BranchMaster:     "wuwei",
				Behind:           1111,
				Ahead:            2222,
				LatestSyncTime:   nil,
				LatestUpdateTime: nil,
				IsDeleted:        true,
			}
			//total int
			err error
		)
		Convey("SaveAggregateBranchDatabase", func(ctx C) {
			err = srv.SaveAggregateBranchDatabase(context.TODO(), branch)
			Convey("Then err should be nil.", func(ctx C) {
				So(err, ShouldBeNil)
			})
		})
	})
}
