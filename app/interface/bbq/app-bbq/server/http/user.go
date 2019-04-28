package http

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	user "go-common/app/service/bbq/user/api"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

//userBase .
func userBase(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(srv.UserBase(c, mid.(int64)))
}

//spaceUserProfile ...
func spaceUserProfile(c *bm.Context) {
	arg := new(v1.SpaceUserProfileRequest)
	mid := int64(0)
	midValue, exists := c.Get("mid")
	if exists {
		mid = midValue.(int64)
	}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.SpaceUserProfile(c, mid, arg.Upmid))
}

//userEdit
func userBaseEdit(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(user.UserBase)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	arg.Mid = mid.(int64)
	c.JSON(srv.UserEdit(c, arg))
}

//addUserLike .
func addUserLike(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(v1.UserLikeAddRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := srv.AddUserLike(c, mid.(int64), arg.SVID)
	c.JSON(resp, err)

	// 埋点
	if err == nil {
		uiLog(c, model.ActionLike, nil)
	}
}

//cancelUserLike .
func cancelUserLike(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(v1.UserLikeCancelRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	resp, err := srv.CancelUserLike(c, mid.(int64), arg.SVID)
	c.JSON(resp, err)

	// 埋点
	if err == nil {
		uiLog(c, model.ActionCancelLike, nil)
	}
}

//userLikeList .
func userLikeList(c *bm.Context) {
	arg := &v1.SpaceSvListRequest{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	mid, exists := c.Get("mid")
	if exists {
		arg.MID = mid.(int64)
	}
	dev, _ := c.Get("device")
	arg.Device = dev.(*bm.Device)
	arg.Size = model.SpaceListLen

	c.JSON(srv.UserLikeList(c, arg))
}

//userFollowList .
func userFollowList(c *bm.Context) {
	arg := new(user.ListRelationUserInfoReq)
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.UserFollowList(c, arg))
}

//userFanList .
func userFanList(c *bm.Context) {
	arg := new(user.ListRelationUserInfoReq)
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.UserFanList(c, arg))
}

//userBlackList .
func userBlackList(c *bm.Context) {
	arg := new(user.ListRelationUserInfoReq)
	mid, exists := c.Get("mid")
	if exists {
		arg.Mid = mid.(int64)
	}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.UserBlackList(c, arg))
}

//userRelationModify .
func userRelationModify(c *bm.Context) {
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(v1.UserRelationRequest)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	var reportAction int
	switch arg.Action {
	case user.FollowAdd:
		reportAction = model.ActionFollow
	case user.FollowCancel:
		reportAction = model.ActionCancelFollow
	case user.BlackAdd:
		reportAction = model.ActionBlack
	case user.BlackCancel:
		reportAction = model.ActionCancelBlack
	default:
		errors.Wrap(ecode.ReqParamErr, "参数验证失败")
		return
	}
	res, err := srv.ModifyRelation(c, mid.(int64), arg.UPMID, arg.Action)
	c.JSON(res, err)

	// 埋点
	if err == nil && reportAction != 0 {
		ext := struct {
			UPMID int64 `json:"up_mid"`
		}{
			UPMID: arg.UPMID,
		}
		uiLog(c, reportAction, ext)
	}
}

//login 登陆
func login(c *bm.Context) {
	arg := new(user.UserBase)
	mid, exists := c.Get("mid")
	if !exists {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	arg.Mid = mid.(int64)
	c.JSON(srv.Login(c, arg))
}

// userUnLike 不感兴趣
func userUnLike(c *bm.Context) {
	tmp, exists := c.Get("mid")
	if !exists || tmp == nil {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	arg := new(v1.UnLikeReq)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.MID = tmp.(int64)

	c.JSON(new(interface{}), nil)

	// 埋点
	uiLog(c, model.ActionUserUnLike, arg)
}
