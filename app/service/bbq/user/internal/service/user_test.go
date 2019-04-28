package service

import (
	"context"
	"encoding/json"
	"go-common/app/service/bbq/user/api"
	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/log"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceLogin(t *testing.T) {
	convey.Convey("Login", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			in = &api.UserBase{Mid: 88895104}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			guard := s.dao.MockRawUserCard(&model.UserCard{Sex: "女", Name: "哈哈"}, nil)
			defer guard.Unpatch()
			res, err := s.Login(c, in)
			data, _ := json.Marshal(res)
			log.V(1).Infow(c, "res", string(data))
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

//
//func TestServiceListUserInfo(t *testing.T) {
//	convey.Convey("ListUserInfo", t, func(convCtx convey.C) {
//		var (
//			c  = context.Background()
//			in = &api.ListUserInfoReq{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.ListUserInfo(c, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceUserEdit(t *testing.T) {
//	convey.Convey("UserEdit", t, func(convCtx convey.C) {
//		var (
//			c  = context.Background()
//			in = &api.UserBase{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.UserEdit(c, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceUserCmsTagEdit(t *testing.T) {
//	convey.Convey("UserCmsTagEdit", t, func(convCtx convey.C) {
//		var (
//			c  = context.Background()
//			in = &api.CmsTagRequest{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.UserCmsTagEdit(c, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceUserFieldEdit(t *testing.T) {
//	convey.Convey("UserFieldEdit", t, func(convCtx convey.C) {
//		var (
//			c  = context.Background()
//			in = &api.UserBase{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.UserFieldEdit(c, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServicebatchUserInfo(t *testing.T) {
//	convey.Convey("batchUserInfo", t, func(convCtx convey.C) {
//		var (
//			c          = context.Background()
//			visitorMID = int64(0)
//			upMIDs     = []int64{}
//			conf       = &api.ListUserInfoConf{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.batchUserInfo(c, visitorMID, upMIDs, conf)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServicegenUserDesc(t *testing.T) {
//	convey.Convey("genUserDesc", t, func(convCtx convey.C) {
//		var (
//			c    = context.Background()
//			base = &api.UserBase{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			userDesc, regionName := s.genUserDesc(c, base)
//			convCtx.Convey("Then userDesc,regionName should not be nil.", func(convCtx convey.C) {
//				convCtx.So(regionName, convey.ShouldNotBeNil)
//				convCtx.So(userDesc, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServicePhoneCheck(t *testing.T) {
//	convey.Convey("PhoneCheck", t, func(convCtx convey.C) {
//		var (
//			c  = context.Background()
//			in = &api.PhoneCheckReq{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.PhoneCheck(c, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceForbidUser(t *testing.T) {
//	convey.Convey("ForbidUser", t, func(convCtx convey.C) {
//		var (
//			c   = context.Background()
//			req = &api.ForbidRequest{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.ForbidUser(c, req)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceReleaseUser(t *testing.T) {
//	convey.Convey("ReleaseUser", t, func(convCtx convey.C) {
//		var (
//			c   = context.Background()
//			req = &api.ReleaseRequest{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.ReleaseUser(c, req)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceUpdateUserVideoView(t *testing.T) {
//	convey.Convey("UpdateUserVideoView", t, func(convCtx convey.C) {
//		var (
//			c   = context.Background()
//			req = &api.UserVideoView{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.UpdateUserVideoView(c, req)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
