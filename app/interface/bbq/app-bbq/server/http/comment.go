package http

import (
	v1 "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

var (
	cmType int64
	cmID   int64
)

func commentInit(c *bm.Context) {
	if cmType = cfg.Comment.Type; cmType == 0 {
		cmType = model.DefaultCmType
	}
	//debug conf
	if deid := cfg.Comment.DebugID; deid > 0 {
		cmID = deid
	}
}

func commentCursor(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &v1.CommentCursorReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	device := dev.(*bm.Device)
	if mid != nil {
		arg.MID = mid.(int64)
	} else {
		arg.MID = 0
	}
	//评论区类型overwrite
	arg.Type = cmType
	if cmID > 0 {
		arg.SvID = cmID
	}
	c.JSON(srv.CommentCursor(c, arg, device))
}

func commentAdd(c *bm.Context) {
	dev, _ := c.Get("device")
	arg := &v1.CommentAddReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	if len([]rune(arg.Message)) > 96 {
		err := ecode.CommentLengthIllegal
		c.JSON(nil, err)
		return
	}
	device := dev.(*bm.Device)
	arg.AccessKey = c.Request.Form.Get("access_key")
	arg.Type = cmType
	if cmID > 0 {
		arg.SvID = cmID
	}
	midVal, _ := c.Get("mid")
	resp, err := srv.CommentAdd(c, midVal.(int64), arg, device)
	c.JSON(resp, err)

	// 埋点
	if err != nil {
		return
	}
	ext := &struct {
		SVID   int64
		Root   int64
		Parent int64
		Type   int64
	}{
		SVID:   arg.SvID,
		Root:   arg.Root,
		Parent: arg.Parent,
		Type:   arg.Type,
	}
	uiLog(c, model.ActionCommentAdd, ext)
}

func commentLike(c *bm.Context) {
	dev, _ := c.Get("device")
	arg := &v1.CommentLikeReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	arg.AccessKey = c.Request.Form.Get("access_key")
	device := dev.(*bm.Device)
	arg.Type = cmType
	if cmID > 0 {
		arg.SvID = cmID
	}
	mid, _ := c.Get("mid")
	err := srv.CommentLike(c, mid.(int64), arg, device)
	c.JSON(nil, err)

	// 埋点
	if err != nil {
		return
	}
	ext := &struct {
		SVID   int64
		RPID   int64
		Action int16
		Type   int64
	}{
		SVID:   arg.SvID,
		RPID:   arg.RpID,
		Action: arg.Action,
		Type:   arg.Type,
	}
	uiLog(c, model.ActionCommentLike, ext)
}

func commentList(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &v1.CommentListReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	device := dev.(*bm.Device)
	if mid != nil {
		arg.MID = mid.(int64)
	} else {
		arg.MID = 0
	}
	arg.Type = cmType
	if cmID > 0 {
		arg.SvID = cmID
	}

	// 这里是转换成评论那边的平台
	if device.RawPlatform == "ios" {
		arg.Plat = 3
	} else if device.RawPlatform == "android" {
		arg.Plat = 2
		arg.Pn++
	}

	c.JSON(srv.CommentList(c, arg, device))
}

func commentSubCursor(c *bm.Context) {
	dev, _ := c.Get("device")
	arg := &v1.CommentSubCursorReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	device := dev.(*bm.Device)
	arg.Type = cmType
	if cmID > 0 {
		arg.SvID = cmID
	}

	var mid int64
	midValue, _ := c.Get("mid")
	if midValue != nil {
		mid = midValue.(int64)
	}
	c.JSON(srv.CommentSubCursor(c, mid, arg, device))
}
