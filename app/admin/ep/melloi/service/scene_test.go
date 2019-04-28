package service

import (
	"context"
	"go-common/app/admin/ep/melloi/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceQueryDraft(t *testing.T) {
	convey.Convey("QueryDraft", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				UserName: "chenmeng",
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.QueryDraft(scene)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateScene(t *testing.T) {
	convey.Convey("UpdateScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID: 66,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			fusing, err := s.UpdateScene(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(fusing, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceAddScene(t *testing.T) {
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
			id, err := s.AddScene(scene)
			convCtx.Convey("Then err should be nil.id should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceQueryAPI(t *testing.T) {
	convey.Convey("QueryAPI", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID:        325,
				SceneType: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.QueryAPI(scene)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceAddConfig(t *testing.T) {
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
			err := s.AddConfig(script)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceQueryTree(t *testing.T) {
	convey.Convey("QueryTree", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				SceneID: 325,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.QueryTree(script)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceQueryScenesByPage(t *testing.T) {
	convey.Convey("QueryScenesByPage", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			sessionID = ""
			scene     = &model.Scene{
				Department: "test",
				Project:    "ep",
				APP:        "melloi",
				SceneName:  "场景压测cmtest01",
				UserName:   "chenmeng",
			}
			page = &model.Pagination{
				PageNum:  int32(1),
				PageSize: int32(20),
			}
			qsrq = &model.QuerySceneRequest{
				Scene:      *scene,
				Pagination: *page,
				Executor:   "chenmeng",
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rsp, err := s.QueryScenesByPage(c, sessionID, qsrq)
			convCtx.Convey("Then err should be nil.rsp should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rsp, convey.ShouldNotBeNil)
			})
		})
	})
}

//func TestServiceAddAndExecuScene(t *testing.T) {
//	convey.Convey("AddAndExecuScene", t, func(convCtx convey.C) {
//		var (
//			c     = context.Background()
//			scene model.Scene
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			resp, err := s.AddAndExecuScene(c, scene)
//			convCtx.Convey("Then err should be nil.resp should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(resp, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestServiceSaveScene(t *testing.T) {
	convey.Convey("SaveScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID:        66,
				IsDraft:   0,
				SceneType: 2,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.SaveScene(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceSaveOrder(t *testing.T) {
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
			req = model.SaveOrderReq{
				GroupOrderList: reqList,
			}
			scene = &model.Scene{
				SceneType: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.SaveOrder(req, scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceQueryRelation(t *testing.T) {
	convey.Convey("QueryRelation", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				ID: 1568,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.QueryRelation(script)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDeleteAPI(t *testing.T) {
	convey.Convey("DeleteAPI", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				ID: 1568,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.DeleteAPI(script)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

//func TestServiceDoScenePtest(t *testing.T) {
//	convey.Convey("DoScenePtest", t, func(convCtx convey.C) {
//		var (
//			c          = context.Background()
//			ptestScene model.DoPtestSceneParam
//			addPtest   bool
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			resp, err := s.DoScenePtest(c, ptestScene, addPtest)
//			convCtx.Convey("Then err should be nil.resp should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(resp, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestServiceDoScenePtestBatch(t *testing.T) {
	convey.Convey("DoScenePtestBatch", t, func(convCtx convey.C) {
		var (
			c           = context.Background()
			ptestScenes model.DoPtestSceneParams
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.DoScenePtestBatch(c, ptestScenes, "")
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

//func TestServiceGetSceneInfo(t *testing.T) {
//	convey.Convey("GetSceneInfo", t, func(convCtx convey.C) {
//		var (
//			scene = &model.Scene{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			sceneInfo := GetSceneInfo(scene)
//			convCtx.Convey("Then sceneInfo should not be nil.", func(convCtx convey.C) {
//				convCtx.So(sceneInfo, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestServiceQueryExistAPI(t *testing.T) {
	convey.Convey("QueryExistAPI", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			sessionID = "24133ff6e561c0cdd0bfb1a8e7802f1c"
			script    = &model.Script{
				SceneID: 325,
			}
			page = &model.Pagination{
				PageNum:  int32(1),
				PageSize: int32(10),
			}
			req = &model.APIInfoRequest{
				Script:     *script,
				Pagination: *page,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.QueryExistAPI(c, sessionID, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceQueryPreview(t *testing.T) {
	convey.Convey("QueryPreview", t, func(convCtx convey.C) {
		var (
			req = &model.Script{
				SceneID: 325,
				GroupID: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			preRes, err := s.QueryPreview(req)
			convCtx.Convey("Then err should be nil.preRes should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(preRes, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceQueryParams(t *testing.T) {
	convey.Convey("QueryParams", t, func(convCtx convey.C) {
		var (
			req = &model.Script{
				SceneID: 325,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, tempRes, err := s.QueryParams(req)
			convCtx.Convey("Then err should be nil.res,tempRes should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(tempRes, convey.ShouldNotBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateBindScene(t *testing.T) {
	convey.Convey("UpdateBindScene", t, func(convCtx convey.C) {
		var (
			bindScene = &model.BindScene{
				SceneID: 325,
				ID:      "2000",
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.UpdateBindScene(bindScene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceQueryDrawRelation(t *testing.T) {
	convey.Convey("QueryDrawRelation", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID: 980,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, tempRes, err := s.QueryDrawRelation(scene)
			convCtx.Convey("Then err should be nil.res,tempRes should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(tempRes, convey.ShouldNotBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDeleteDraft(t *testing.T) {
	convey.Convey("DeleteDraft", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				UserName: "chenmeng",
				ID:       980,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.DeleteDraft(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceQueryConfig(t *testing.T) {
	convey.Convey("QueryConfig", t, func(convCtx convey.C) {
		var (
			script = &model.Script{
				SceneID: 325,
				GroupID: 1,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1, err := s.QueryConfig(script)
			convCtx.Convey("Then err should be nil.p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceDeleteScene(t *testing.T) {
	convey.Convey("DeleteScene", t, func(convCtx convey.C) {
		var (
			scene = &model.Scene{
				ID: 325,
			}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := s.DeleteScene(scene)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
