package dao

import (
	"go-common/app/admin/ep/melloi/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAddScene(t *testing.T) {
	convey.Convey("AddScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				SceneName:  "CMTest001",
				UserName:   "chenmeng",
				Department: "test",
				Project:    "ep",
				APP:        "melloi",
				IsDraft:    1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			sceneId, err := d.AddScene(scene)
			convCtx.Convey("Then err should be nil.sceneId should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(sceneId, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryDraft(t *testing.T) {
	convey.Convey("QueryDraft", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				UserName: "chenmeng",
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryDraft(scene)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateScene(t *testing.T) {
	convey.Convey("UpdateScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID: 66,
			}
			scriptIDList = []int{1831, 1832, 1833, 1834, 1835}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			fusing, err := d.UpdateScene(scene, scriptIDList)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(fusing, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSaveScene(t *testing.T) {
	convey.Convey("SaveScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID:        66,
				IsDraft:   0,
				SceneType: 2,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SaveScene(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSaveOrderAuto(t *testing.T) {
	convey.Convey("SaveOrder", t, func(convCtx convey.C) {
		var (
			reqList = []*model.GroupOrder{
				{
					GroupID:  1,
					RunOrder: 2,
					ID:       980,
					TestName: "cm-test",
				},
				{
					GroupID:  1,
					RunOrder: 1,
					ID:       973,
					TestName: "status111",
				},
			}
			sceneAuto = &model.Scene{
				SceneType: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SaveOrder(reqList, sceneAuto)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSaveOrderGroup(t *testing.T) {
	convey.Convey("SaveOrder", t, func(convCtx convey.C) {
		var (
			reqList = []*model.GroupOrder{
				{
					GroupID:  1,
					RunOrder: 2,
					ID:       980,
					TestName: "cm-test",
				},
				{
					GroupID:  1,
					RunOrder: 1,
					ID:       973,
					TestName: "status111",
				},
			}
			sceneGroup = &model.Scene{
				SceneType: 2,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SaveOrder(reqList, sceneGroup)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryGroupId(t *testing.T) {
	convey.Convey("QueryGroupId", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				ID: 1568,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, groupId, err := d.QueryGroupId(script)
			convCtx.Convey("Then err should be nil.res,groupId should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(groupId, convey.ShouldNotBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryRelation(t *testing.T) {
	convey.Convey("QueryRelation", t, func(convCtx convey.C) {
		var (
			groupId = 11
			script  = &model.Script{
				ID: 1568,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryRelation(groupId, script)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryAPI(t *testing.T) {
	convey.Convey("QueryAPI", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID:        325,
				SceneType: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryAPI(scene)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeleteAPI(t *testing.T) {
	d.DB.Model(&model.Script{}).Where("ID = 1568").Update("active", 1)
	convey.Convey("DeleteAPI", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				ID: 1568,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DeleteAPI(script)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoAddConfig(t *testing.T) {
	d.DB.Exec("update script set threads_sum = 10, load_time = 10, ready_time = 10 where scene_id = 325 and group_id = 1")
	convey.Convey("AddConfig", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				ThreadsSum: 120,
				LoadTime:   300,
				ReadyTime:  10,
				SceneID:    325,
				GroupID:    1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddConfig(script)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryTree(t *testing.T) {
	convey.Convey("QueryTree", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				SceneID: 325,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryTree(script)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryScenesByPage(t *testing.T) {
	convey.Convey("QueryScenesByPage", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				Department: "test",
				Project:    "ep",
				APP:        "melloi",
				SceneName:  "场景压测cmtest01",
				UserName:   "chenmeng",
			}
			pn        = int32(1)
			ps        = int32(20)
			treeNodes = []string{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			qsr, err := d.QueryScenesByPage(scene, pn, ps, treeNodes)
			convCtx.Convey("Then err should be nil.qsr should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(qsr, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryScenesByPageWhiteName(t *testing.T) {
	convey.Convey("QueryScenesByPageWhiteName", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				Department: "test",
				Project:    "ep",
				APP:        "melloi",
				SceneName:  "场景压测cmtest01",
				UserName:   "chenmeng",
			}
			pn = int32(1)
			ps = int32(20)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			qsr, err := d.QueryScenesByPageWhiteName(scene, pn, ps)
			convCtx.Convey("Then err should be nil.qsr should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(qsr, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryScenes(t *testing.T) {
	convey.Convey("QueryScenes", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				UserName: "chenmeng",
			}
			pn = 1
			ps = 10
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			scenes, err := d.QueryScenes(scene, pn, ps)
			convCtx.Convey("Then err should be nil.scenes should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(scenes, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryExistAPI(t *testing.T) {
	convey.Convey("QueryExistAPI", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				SceneID: 325,
			}
			pageNum   = int32(1)
			pageSize  = int32(10)
			sceneId   = int(325)
			treeNodes = []string{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryExistAPI(script, pageNum, pageSize, sceneId, treeNodes)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryGroup(t *testing.T) {
	convey.Convey("QueryGroup", t, func(convCtx convey.C) {
		var (
			sceneId = 325
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryGroup(sceneId)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryPreview(t *testing.T) {
	convey.Convey("QueryPreview", t, func(convCtx convey.C) {
		var (
			sceneId = 325
			groupId = 1
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryPreview(sceneId, groupId)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoQueryUsefulParams(t *testing.T) {
	convey.Convey("QueryUsefulParams", t, func(convCtx convey.C) {
		var (
			sceneId = 325
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryUsefulParams(sceneId)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateBindScene(t *testing.T) {
	convey.Convey("UpdateBindScene", t, func(convCtx convey.C) {
		var (
			bindScene = &model.BindScene{
				SceneID: 325,
				ID:      "2000",
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateBindScene(bindScene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryDrawRelation(t *testing.T) {
	convey.Convey("QueryDrawRelation", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID: 980,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryDrawRelation(scene)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeleteDraft(t *testing.T) {
	convey.Convey("DeleteDraft", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				UserName: "chenmeng",
				ID:       980,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DeleteDraft(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoQueryConfig(t *testing.T) {
	convey.Convey("QueryConfig", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				SceneID: 325,
				GroupID: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.QueryConfig(script)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDeleteScene(t *testing.T) {
	convey.Convey("DeleteScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID: 325,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.DeleteScene(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
