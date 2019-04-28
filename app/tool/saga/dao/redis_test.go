package dao

import (
	"context"
	"encoding/json"
	"testing"

	"go-common/app/tool/saga/model"
	"go-common/library/cache/redis"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMergeTaskKey(t *testing.T) {
	convey.Convey("mergeTaskKey", t, func(ctx convey.C) {
		var (
			taskType = int(111)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mergeTaskKey(taskType)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_task_111")
			})
		})
	})
}

func TestDaoMrIIDKey(t *testing.T) {
	convey.Convey("mrIIDKey", t, func(ctx convey.C) {
		var (
			mrIID = int(222)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mrIIDKey(mrIID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_mrIID_222")
			})
		})
	})
}

func TestDaoPingRedis(t *testing.T) {
	convey.Convey("pingRedis", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingRedis(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddMRIID(t *testing.T) {
	convey.Convey("AddMRIID", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mrIID  = int(333)
			expire = int(1800)
			conn   = d.redis.Get(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMRIID(c, mrIID, expire)
			mrID, err := redis.Int(conn.Do("GET", "saga_mrIID_333"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// 检查redis存储结果
				ctx.So(mrID, convey.ShouldEqual, mrIID)
			})
		})
	})
}

func TestDaoExistMRIID(t *testing.T) {
	convey.Convey("ExistMRIID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mrIID = int(444)
			conn  = d.redis.Get(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			conn.Do("DEL", "saga_mrIID_444")
			ok, err := d.ExistMRIID(c, mrIID)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldBeFalse)
			})
		})
		ctx.Convey("When MRIID exist", func(ctx convey.C) {
			conn.Do("SET", "saga_mrIID_444", 1)
			ok, err := d.ExistMRIID(c, mrIID)
			ctx.Convey("Then ok should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ok, convey.ShouldBeTrue)
			})
		})
	})
}

func TestDaoDeleteMRIID(t *testing.T) {
	convey.Convey("DeleteMRIID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mrIID = int(555)
			conn  = d.redis.Get(c)
		)
		ctx.Convey("When data not exist", func(ctx convey.C) {
			err := d.DeleteMRIID(c, mrIID)
			value, _ := redis.Int(conn.Do("GET", "saga_mrIID_555"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, 0)
			})
		})
		ctx.Convey("When data exist", func(ctx convey.C) {
			redis.String(conn.Do("SET", "saga_mrIID_555", "Test"))
			err := d.DeleteMRIID(c, mrIID)
			value, _ := redis.Int(conn.Do("GET", "saga_mrIID_555"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, 0)
			})
		})
	})
}

func TestDaoPushMergeTask(t *testing.T) {
	convey.Convey("PushMergeTask", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskType = int(666)
			taskInfo = &model.TaskInfo{NoteID: 111, Event: nil, Repo: nil}
			conn     = d.redis.Get(c)
		)
		bs, _ := json.Marshal(taskInfo)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.PushMergeTask(c, taskType, taskInfo)
			value, _ := redis.Bytes(conn.Do("LPOP", "saga_task_666"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// push 运行函数 pop 检验结果
				ctx.So(string(value), convey.ShouldEqual, string(bs))
			})
		})
	})
}

func TestDaoDeleteMergeTask(t *testing.T) {
	convey.Convey("DeleteMergeTask", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskType = int(777)
			taskInfo = &model.TaskInfo{NoteID: 777, Event: nil, Repo: nil}
			conn     = d.redis.Get(c)
		)
		ctx.Convey("When data is not exist", func(ctx convey.C) {
			err := d.DeleteMergeTask(c, taskType, taskInfo)
			value, _ := redis.Int(conn.Do("LPOP", "saga_task_777"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, 0)
			})
		})
		ctx.Convey("When data is exist", func(ctx convey.C) {
			taskInfoJSON, _ := json.Marshal(taskInfo)
			redis.String(conn.Do("LPUSH", "saga_task_777", taskInfoJSON))
			err := d.DeleteMergeTask(c, taskType, taskInfo)
			value, _ := redis.Int(conn.Do("LPOP", "saga_task_777"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, 0)
			})
		})
	})
}

func TestDaoMergeTasks(t *testing.T) {
	convey.Convey("MergeTasks", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			taskType = int(888)
			conn     = d.redis.Get(c)
			taskInfo = &model.TaskInfo{NoteID: 888, Event: nil, Repo: nil}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			taskInfoJSON, _ := json.Marshal(taskInfo)
			conn.Do("LPUSH", "saga_task_888", taskInfoJSON)
			count, taskInfos, err := d.MergeTasks(c, taskType)
			taskInfoFirst, _ := json.Marshal(taskInfos[0])
			ctx.Convey("Then err should be nil.count,taskInfos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(string(taskInfoFirst), convey.ShouldEqual, string(taskInfoJSON))
				ctx.So(count, convey.ShouldNotEqual, 0)
			})
		})
	})
}

func TestDaoMergeInfoKey(t *testing.T) {
	convey.Convey("mergeInfoKey", t, func(ctx convey.C) {
		var (
			projID = int(999)
			branch = "test"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := mergeInfoKey(projID, branch)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_mergeInfo_999_test")
			})
		})
	})
}

func TestDaoPathOwnerKey(t *testing.T) {
	convey.Convey("pathOwnerKey", t, func(ctx convey.C) {
		var (
			projID = int(1111)
			branch = "test"
			path   = "."
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := pathOwnerKey(projID, branch, path)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_PathOwner_1111_test_.")
			})
		})
	})
}

func TestDaoPathReviewerKey(t *testing.T) {
	convey.Convey("pathReviewerKey", t, func(ctx convey.C) {
		var (
			projID = int(2222)
			branch = "test"
			path   = "."
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := pathReviewerKey(projID, branch, path)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_PathReviewer_2222_test_.")
			})
		})
	})
}

func TestDaoAuthInfoKey(t *testing.T) {
	convey.Convey("authInfoKey", t, func(ctx convey.C) {
		var (
			projID = int(3333)
			mrIID  = int(3333)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := authInfoKey(projID, mrIID)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_auth_3333_3333")
			})
		})
	})
}

func TestDaoSetMergeInfo(t *testing.T) {
	convey.Convey("SetMergeInfo", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			projID    = int(4444)
			branch    = "test"
			mergeInfo = &model.MergeInfo{}
			conn      = d.redis.Get(c)
		)
		redis.String(conn.Do(""))
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetMergeInfo(c, projID, branch, mergeInfo)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoMergeInfo(t *testing.T) {
	convey.Convey("MergeInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			projID = int(5555)
			branch = "test"
			conn   = d.redis.Get(c)
			info   = []byte(`{"NoteID":111,"Event":null,"Repo":null}`)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			conn.Do("SET", "saga_mergeInfo_5555_test", info)
			ok, mergeInfo, err := d.MergeInfo(c, projID, branch)
			ctx.Convey("Then err should be nil.ok,mergeInfo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mergeInfo, convey.ShouldNotBeNil)
				ctx.So(ok, convey.ShouldBeTrue)
			})
		})
	})
}

func TestDaoDeleteMergeInfo(t *testing.T) {
	convey.Convey("DeleteMergeInfo", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			projID = int(6666)
			branch = "test"
			conn   = d.redis.Get(c)
			info   = []byte(`{"NoteID":111,"Event":null,"Repo":null}`)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			conn.Do("SET", "saga_mergeInfo_6666_test", info)
			err := d.DeleteMergeInfo(c, projID, branch)
			value, _ := redis.Int(conn.Do("GET", "saga_mergeInfo_6666_test"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, 0)
			})
		})
	})
}

func TestDaoPathAuthKey(t *testing.T) {
	convey.Convey("pathAuthKey", t, func(ctx convey.C) {
		var (
			projID = int(7777)
			branch = "test"
			path   = "."
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := pathAuthKey(projID, branch, path)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_path_auth_7777_test_.")
			})
		})
	})
}

func TestDaoPathAuthR(t *testing.T) {
	convey.Convey("PathAuthR", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			projID = int(8888)
			branch = "test"
			path   = "."
			conn   = d.redis.Get(c)
		)
		redis.Int(conn.Do("SET", "saga_path_auth_8888_test_.", `{"Owners": ["c"]}`))
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authUsers, err := d.PathAuthR(c, projID, branch, path)
			authUsersJSON, _ := json.Marshal(authUsers)
			ctx.Convey("Then err should be nil.authUsers should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(string(authUsersJSON), convey.ShouldEqual, `{"Owners":["c"],"Reviewers":null}`)
				ctx.So(authUsers, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetPathAuthR(t *testing.T) {
	convey.Convey("SetPathAuthR", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			projID    = int(9999)
			branch    = "test"
			path      = "."
			authUsers = &model.AuthUsers{Owners: []string{"testOwner"}, Reviewers: []string{"testReviewers"}}
			conn      = d.redis.Get(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetPathAuthR(c, projID, branch, path, authUsers)
			value, _ := redis.String(conn.Do("GET", "saga_path_auth_9999_test_."))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, `{"Owners":["testOwner"],"Reviewers":["testReviewers"]}`)
			})
		})
	})
}

func TestDaoDeletePathAuthR(t *testing.T) {
	convey.Convey("DeletePathAuthR", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			projID    = int(88)
			branch    = "test"
			path      = "."
			authUsers = &model.AuthUsers{Owners: []string{"testOwner"}, Reviewers: []string{"testReviewers"}}
			conn      = d.redis.Get(c)
		)
		ctx.Convey("When data not exist", func(ctx convey.C) {
			err := d.DeletePathAuthR(c, projID, branch, path)
			value, _ := redis.String(conn.Do("GET", "saga_path_auth_88_test_."))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, "")
			})
		})
		ctx.Convey("When data exist", func(ctx convey.C) {
			err := d.SetPathAuthR(c, projID, branch, path, authUsers)
			value, _ := redis.String(conn.Do("GET", "saga_path_auth_88_test_."))
			ctx.Convey("Then err should be nil 1.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, `{"Owners":["testOwner"],"Reviewers":["testReviewers"]}`)
			})
			err = d.DeletePathAuthR(c, projID, branch, path)
			value, _ = redis.String(conn.Do("GET", "saga_path_auth_88_test_."))
			ctx.Convey("Then err should be nil 2.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, "")
			})
		})
	})
}

func TestDaoSetReportStatus(t *testing.T) {
	convey.Convey("SetReportStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			projID = int(11111)
			mrIID  = int(11111)
			result bool
			conn   = d.redis.Get(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetReportStatus(c, projID, mrIID, result)
			value, _ := redis.String(conn.Do("GET", "saga_auth_11111_11111"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, "false")
			})
		})
	})
}

func TestDaoReportStatus(t *testing.T) {
	convey.Convey("ReportStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			projID = int(22222)
			mrIID  = int(22222)
			conn   = d.redis.Get(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			conn.Do("SET", "saga_auth_22222_22222", "true")
			result, err := d.ReportStatus(c, projID, mrIID)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldEqual, true)
			})
		})
	})
}

func TestDaoDeleteReportStatus(t *testing.T) {
	convey.Convey("DeleteReportStatus", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			projID = int(33333)
			mrIID  = int(33333)
			conn   = d.redis.Get(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			conn.Do("SET", "saga_auth_33333_33333", "true")
			err := d.DeleteReportStatus(c, projID, mrIID)
			value, _ := redis.Int(conn.Do("GET", "saga_auth_33333_33333"))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(value, convey.ShouldEqual, 0)
			})
		})
	})
}
