package dao

import (
	"context"
	"strconv"
	"testing"

	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/utils"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoCommit(t *testing.T) {
	convey.Convey("test dao commit", t, func(ctx convey.C) {

		var (
			commit *model.StatisticsCommits
			err    error
			total  int
		)

		ctx.Convey("When project not exist", func(ctx convey.C) {
			total, err = d.HasCommit(55555, "55555")
			ctx.Convey("HasCommit err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldEqual, 0)
			})
			err = d.DelCommit(55555, "55555")
			ctx.Convey("DelCommit err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})

		ctx.Convey("When project exist", func(ctx convey.C) {
			var (
				commitID = "11111"
				PorjID   = 11111
			)

			ctx.Convey("create commit", func(ctx convey.C) {
				info := &model.StatisticsCommits{
					CommitID:  commitID,
					ProjectID: PorjID,
					Title:     "test",
				}
				err = d.CreateCommit(info)
				ctx.Convey("CreateCommit err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})

			ctx.Convey("has commit after create", func(ctx convey.C) {
				total, err = d.HasCommit(PorjID, commitID)
				ctx.Convey("HasCommit err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(total, convey.ShouldEqual, 1)
				})
			})

			ctx.Convey("query commit", func(ctx convey.C) {
				commit, err = d.QueryCommitByID(PorjID, commitID)
				ctx.Convey("Then err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(commit.Title, convey.ShouldEqual, "test")
				})
			})

			ctx.Convey("create second temp commit", func(ctx convey.C) {
				temp := &model.StatisticsCommits{
					CommitID:  "33333",
					ProjectID: 33333,
					Title:     "temp",
				}
				err = d.CreateCommit(temp)
				ctx.Convey("CreateCommit again err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})

			ctx.Convey("has commit when exist two rows", func(ctx convey.C) {
				total, err = d.HasCommit(PorjID, commitID)
				ctx.Convey("HasCommit err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(total, convey.ShouldEqual, 1)
				})
			})

			ctx.Convey("update commit", func(ctx convey.C) {
				updateInfo := &model.StatisticsCommits{
					CommitID:  commitID,
					ProjectID: PorjID,
					Title:     "update",
				}
				err = d.UpdateCommit(PorjID, commitID, updateInfo)
				ctx.Convey("UpdateCommit err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})

			ctx.Convey("query commit after update", func(ctx convey.C) {
				commit, err = d.QueryCommitByID(PorjID, commitID)
				ctx.Convey("Then err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(commit.Title, convey.ShouldEqual, "update")
				})
			})

			ctx.Convey("update temp commit", func(ctx convey.C) {
				updateInfo := &model.StatisticsCommits{
					CommitID:  "33333",
					ProjectID: 33333,
					Title:     "update temp",
				}
				err = d.UpdateCommit(33333, "33333", updateInfo)
				ctx.Convey("UpdateCommit temp err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})

			ctx.Convey("query temp commit after update", func(ctx convey.C) {
				commit, err = d.QueryCommitByID(33333, "33333")
				ctx.Convey("Then err should be nil.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(commit.Title, convey.ShouldEqual, "update temp")
				})
			})

			ctx.Convey("delete commit", func(ctx convey.C) {
				err = d.DelCommit(PorjID, commitID)
				ctx.Convey("DelCommit commit exist.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})

			ctx.Convey("has commit after delete", func(ctx convey.C) {
				total, err = d.HasCommit(PorjID, commitID)
				ctx.Convey("HasCommit after delete .", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(total, convey.ShouldEqual, 0)
				})
			})

			ctx.Convey("has temp commit when not delete", func(ctx convey.C) {
				total, err = d.HasCommit(33333, "33333")
				ctx.Convey("HasCommit after delete .", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(total, convey.ShouldEqual, 1)
				})
			})

			ctx.Convey("delete temp commit", func(ctx convey.C) {
				err = d.DelCommit(33333, "33333")
				ctx.Convey("DelCommit commit exist.", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
				})
			})

			ctx.Convey("has tem commit after delete", func(ctx convey.C) {
				total, err = d.HasCommit(33333, "33333")
				ctx.Convey("HasCommit after delete .", func(ctx convey.C) {
					ctx.So(err, convey.ShouldBeNil)
					ctx.So(total, convey.ShouldEqual, 0)
				})
			})
		})
	})
}

func TestDaoCommitSpe(t *testing.T) {
	convey.Convey("test dao commit include special char", t, func(ctx convey.C) {
		var err error

		commitSpe := &model.StatisticsCommits{
			CommitID:  "22222",
			ProjectID: 22222,
			Title:     "🇭相关",
		}

		ctx.Convey("create commit", func(ctx convey.C) {
			err = d.CreateCommit(commitSpe)
			ctx.Convey("CreateCommit err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(err.Error(), convey.ShouldContainSubstring, "Incorrect string value")
			})
		})

		ctx.Convey("create commit after to ascii", func(ctx convey.C) {
			commitSpe.Title = strconv.QuoteToASCII(commitSpe.Title)
			err = d.CreateCommit(commitSpe)
			ctx.Convey("CreateCommit err should be nil after to ascii.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})

		ctx.Convey("delete commit", func(ctx convey.C) {
			err = d.DelCommit(22222, "22222")
			ctx.Convey("DelCommit err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestQueryPipelinesByTime(t *testing.T) {
	convey.Convey("QueryPipelinesByTime", t, func(ctx convey.C) {
		var (
			req = &model.PipelineDataReq{
				Branch: "master",
				State:  "success",
				User:   "wangweizhen",
			}
			since = `2018-12-01 00:00:00`
			until = `2018-12-02 00:00:00`
		)
		ctx.Convey("query pipelines info", func(ctx convey.C) {

			total, statNum, _, err := d.QueryPipelinesByTime(682, req, since, until)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldBeGreaterThan, 1)
				ctx.So(statNum, convey.ShouldBeGreaterThan, 1)
			})
		})
	})
}

func TestMRByProjectID(t *testing.T) {
	convey.Convey("QueryMRByProjectID", t, func(ctx convey.C) {
		since, until := utils.CalSyncTime()
		ctx.Convey("query MRByProjectID info", func(ctx convey.C) {

			mrs, err := d.MRByProjectID(context.TODO(), 682, since, until)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(len(mrs), convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestDaoBranchSpe(t *testing.T) {
	convey.Convey("test dao branch include special char", t, func(ctx convey.C) {
		var (
			branch = &model.StatisticsBranches{
				ProjectID:          666,
				ProjectName:        "六六大顺",
				CommitID:           "6661",
				BranchName:         "666",
				Protected:          false,
				Merged:             false,
				DevelopersCanPush:  false,
				DevelopersCanMerge: false,
				IsDeleted:          true,
			}
			total int
			err   error
		)

		ctx.Convey("create branch", func(ctx convey.C) {
			err = d.CreateBranch(context.TODO(), branch)
			ctx.Convey("CreateBranch err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("has branch", func(ctx convey.C) {
			total, err = d.HasBranch(context.TODO(), 666, "666")
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(total, convey.ShouldBeGreaterThan, 0)
			})
		})
		ctx.Convey("update branch", func(ctx convey.C) {
			err = d.UpdateBranch(context.TODO(), 666, "666", branch)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
