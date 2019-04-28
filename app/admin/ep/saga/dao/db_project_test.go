package dao

import (
	"testing"

	"go-common/app/admin/ep/saga/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestProjectExist(t *testing.T) {
	convey.Convey("ProjectExist", t, func(ctx convey.C) {
		var (
			projectID  = int(682)
			projectID2 = int(6820)
		)
		ctx.Convey("When projectID exist", func(ctx convey.C) {
			b, err := d.ProjectExist(projectID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(b, convey.ShouldEqual, true)
			})
		})
		ctx.Convey("When projectID not exist", func(ctx convey.C) {
			bb, err := d.ProjectExist(projectID2)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bb, convey.ShouldEqual, false)
			})
		})
	})
}

func TestFavorite(t *testing.T) {
	convey.Convey("Favorite project", t, func(ctx convey.C) {
		var (
			userName  = "zhangsan"
			projectID = 333
		)

		ctx.Convey("add favorite project", func(ctx convey.C) {
			err := d.AddFavorite(userName, projectID)
			ctx.Convey("The err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("get favorite project", func(ctx convey.C) {
			favorites, err := d.FavoriteProjects(userName)
			ctx.Convey("The err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(favorites), convey.ShouldEqual, 1)
				ctx.So(favorites[0].UserName, convey.ShouldEqual, "zhangsan")
				ctx.So(favorites[0].ProjID, convey.ShouldEqual, 333)
			})
		})
		ctx.Convey("delete favorite project", func(ctx convey.C) {
			err := d.DelFavorite(userName, projectID)
			ctx.Convey("delete err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			favorites, err := d.FavoriteProjects(userName)
			ctx.Convey("get err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(favorites), convey.ShouldEqual, 0)
			})
		})
	})
}

func TestProjectInfoByID(t *testing.T) {
	convey.Convey("ProjectInfoByID", t, func(ctx convey.C) {

		ctx.Convey("add project info", func(ctx convey.C) {
			addInfo := &model.ProjectInfo{
				ProjectID:     11111,
				Name:          "myProject",
				Description:   "myProject des",
				WebURL:        "git@git.bilibili.test/test.git",
				Repo:          "git@git.bilibili.test/test.git",
				DefaultBranch: "master",
				Saga:          true,
				Runner:        false,
			}
			err := d.AddProjectInfo(addInfo)
			ctx.Convey("Add Project Info. Then err should be nil. ", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			info, err := d.ProjectInfoByID(11111)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info.Name, convey.ShouldEqual, "myProject")
				ctx.So(info.Repo, convey.ShouldEqual, "git@git.bilibili.test/test.git")
				ctx.So(info.Saga, convey.ShouldEqual, true)
			})
		})

		ctx.Convey("update project info", func(ctx convey.C) {
			updateInfo := &model.ProjectInfo{
				ProjectID:     11111,
				Name:          "lisi",
				Description:   "lisi des",
				WebURL:        "git@git.bilibili.test/test.git",
				Repo:          "git@git.bilibili.test/test.git",
				DefaultBranch: "dev",
				Saga:          true,
				Runner:        false,
			}
			err := d.UpdateProjectInfo(11111, updateInfo)
			ctx.Convey("update Project Info. Then err should be nil. ", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			info, err := d.ProjectInfoByID(11111)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info.Name, convey.ShouldEqual, "lisi")
				ctx.So(info.Description, convey.ShouldEqual, "lisi des")
				ctx.So(info.DefaultBranch, convey.ShouldEqual, "dev")
			})
		})

		ctx.Convey("query not exist project info", func(ctx convey.C) {
			_, err := d.ProjectInfoByID(6820)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHasProjectInfo(t *testing.T) {
	convey.Convey("HasProjectInfo", t, func(ctx convey.C) {
		var (
			projectID  = int(682)
			projectID2 = int(6820)
		)
		ctx.Convey("When projectID exist", func(ctx convey.C) {
			b, err := d.HasProjectInfo(projectID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(b, convey.ShouldEqual, true)
			})
		})
		ctx.Convey("When projectID not exist", func(ctx convey.C) {
			b, err := d.HasProjectInfo(projectID2)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(b, convey.ShouldEqual, false)
			})
		})
	})
}

func TestQueryProjectInfo(t *testing.T) {
	convey.Convey("QueryProjectInfo", t, func(ctx convey.C) {
		var (
			projectName = "andruid"
			projectUrl  = "git@git-test.bilibili.co:android/andruid.git"

			projectInfoRequest = &model.ProjectInfoRequest{
				Name: projectName,
			}
		)
		ctx.Convey("query project exist", func(ctx convey.C) {
			var url string
			total, projectInfo, err := d.QueryProjectInfo(false, projectInfoRequest)
			for _, project := range projectInfo {
				url = project.Repo
			}
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldEqual, 1)
				ctx.So(url, convey.ShouldEqual, projectUrl)
			})
		})
	})
}

func TestQueryConfigInfo(t *testing.T) {
	convey.Convey("QueryConfigInfo", t, func(ctx convey.C) {
		var (
			projectName  = "Kratos"
			projectName2 = "android-blue"
			queryObject  = "saga"
		)
		ctx.Convey("query saga config exist", func(ctx convey.C) {
			total, err := d.QueryConfigInfo(projectName, "", "", queryObject)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldEqual, 1)
			})
		})

		ctx.Convey("query saga config not exist", func(ctx convey.C) {
			total, err := d.QueryConfigInfo(projectName2, "", "", queryObject)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldEqual, 0)
			})
		})
	})
}
