package http

import (
	"go-common/app/service/bbq/user/api"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

//
////userBase .
//func userBase(c *bm.Context) {
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	c.JSON(srv.UserBase(c, mid.(int64)))
//}
//
////spaceUserProfile ...
//func spaceUserProfile(c *bm.Context) {
//	arg := new(v1.SpaceUserProfileRequest)
//	mid := int64(0)
//	midValue, exists := c.Get("mid")
//	if exists {
//		mid = midValue.(int64)
//	}
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	if arg.Upmid == 0 {
//		c.JSON(nil, ecode.ReqParamErr)
//		return
//	}
//	c.JSON(srv.SpaceUserProfile(c, mid, arg))
//}
//
////userEdit
//func userBaseEdit(c *bm.Context) {
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	arg := new(v1.UserBaseEditRequest)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	c.JSON(srv.UserEdit(c, mid.(int64), arg))
//}
//
////addUserLike .
//func addUserLike(c *bm.Context) {
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	arg := new(v1.UserLikeAddRequest)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	resp, err := srv.AddUserLike(c, mid.(int64), arg.SVID)
//	c.JSON(resp, err)
//
//	// 埋点
//	if err == nil {
//		uiLog(c, model.ActionLike, nil)
//	}
//}
//
////cancelUserLike .
//func cancelUserLike(c *bm.Context) {
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	arg := new(v1.UserLikeCancelRequest)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	resp, err := srv.CancelUserLike(c, mid.(int64), arg.SVID)
//	c.JSON(resp, err)
//
//	// 埋点
//	if err == nil {
//		uiLog(c, model.ActionCancelLike, nil)
//	}
//}
//
////userLikeList .
//func userLikeList(c *bm.Context) {
//	arg := &v1.SpaceSvListRequest{}
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	mid, exists := c.Get("mid")
//	if exists {
//		arg.MID = mid.(int64)
//	}
//	dev, _ := c.Get("device")
//	arg.Device = dev.(*bm.Device)
//	arg.Size = model.SpaceListLen
//
//	c.JSON(srv.UserLikeList(c, arg))
//}
//
//
////userFollowList .
//func userFollowList(c *bm.Context) {
//	arg := new(api.ListRelationUserInfoReq)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	c.JSON(svc.ListFollowUserInfo(c, arg))
//}
//
////userFanList .
//func userFanList(c *bm.Context) {
//	arg := new(api.ListRelationUserInfoReq)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	c.JSON(svc.ListFanUserInfo(c, arg))
//}
//
////userInfoBlackList .
//func userInfoBlackList(c *bm.Context) {
//	arg := new(api.ListRelationUserInfoReq)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	c.JSON(svc.ListBlackUserInfo(c, arg))
//}

//userBlackList .
func userBlackList(c *bm.Context) {
	arg := new(api.ListRelationReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(svc.ListBlack(c, arg))
}

// userInfo used for audit
func userInfoList(c *bm.Context) {
	arg := new(api.ListUserInfoReq)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(svc.ListUserInfo(c, arg))
}

// userInfo used for audit
func auditModifyUser(c *bm.Context) {
	arg := new(api.UserBase)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(svc.UserFieldEdit(c, arg))
}

// userInfo used for audit
func auditModifyCmsTag(c *bm.Context) {
	arg := new(api.CmsTagRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(svc.UserCmsTagEdit(c, arg))
}

func forbidUser(c *bm.Context) {
	arg := new(api.ForbidRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数校验失败")
		return
	}
	c.JSON(svc.ForbidUser(c, arg))
}

func releaseUser(c *bm.Context) {
	arg := new(api.ReleaseRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数校验失败")
		return
	}
	c.JSON(svc.ReleaseUser(c, arg))
}

////userRelationModify .
//func userRelationModify(c *bm.Context) {
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	arg := new(v1.UserRelationRequest)
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	var res *v1.RelationResponse
//	var err error
//	var reportAction int
//	switch arg.Action {
//	case v1.FollowAdd:
//		res, err = srv.AddUserFollow(c, mid.(int64), arg.UPMID)
//		reportAction = model.ActionFollow
//	case v1.FollowCancel:
//		res, err = srv.CancelUserFollow(c, mid.(int64), arg.UPMID)
//		reportAction = model.ActionCancelFollow
//	case v1.BlackAdd:
//		res, err = srv.AddUserBlack(c, mid.(int64), arg.UPMID)
//	case v1.BlackCancel:
//		res, err = srv.CancelUserBlack(c, mid.(int64), arg.UPMID)
//	default:
//		errors.Wrap(ecode.ReqParamErr, "参数验证失败")
//		return
//	}
//	c.JSON(res, err)
//
//	// 埋点
//	if err == nil && reportAction != 0 {
//		ext := struct {
//			UPMID int64 `json:"up_mid"`
//		}{
//			UPMID: arg.UPMID,
//		}
//		uiLog(c, reportAction, ext)
//	}
//}
//
////login 登陆
//func login(c *bm.Context) {
//	arg := new(v1.LoginRequest)
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	if err := c.Bind(arg); err != nil {
//		errors.Wrap(err, "参数验证失败")
//		return
//	}
//	c.JSON(srv.Login(c, arg.NewTag, mid.(int64)))
//}
//
////手机验证
//func phoneCheck(c *bm.Context) {
//	mid, exists := c.Get("mid")
//	if !exists {
//		c.JSON(nil, ecode.NoLogin)
//		return
//	}
//	c.JSON(srv.PhoneCheck(c, mid.(int64)))
//}

// updateUserVideoView 更新用户作品播放量
func updateUserVideoView(c *bm.Context) {
	args := &api.UserVideoView{}
	if err := c.Bind(args); err != nil {
		errors.Wrap(err, "参数校验失败")
		return
	}
	c.JSON(svc.UpdateUserVideoView(c, args))
}
