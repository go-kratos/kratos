package http

import (
	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func codeInfo(c *bm.Context) {
	arg := new(struct {
		Code string `form:"code" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}

	c.JSON(vipSvc.CodeInfo(c, arg.Code))
}

func codeInfos(c *bm.Context) {
	arg := new(struct {
		Codes []string `form:"codes,split" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.CodeInfos(c, arg.Codes))
}

func openCode(c *bm.Context) {
	var (
		token *model.TokenResq
		code  *model.VipResourceCode
		err   error
	)
	arg := new(struct {
		Token  string `form:"token" validate:"required"`
		Code   string `form:"code" validate:"required"`
		Verify string `form:"verify" validate:"required"`
		Mid    int64  `form:"mid" validate:"required"`
	})

	if err = c.Bind(arg); err != nil {
		return
	}
	if token, err = vipSvc.Verify(c, arg.Token, arg.Verify); err != nil {
		c.JSON(nil, ecode.CreativeGeetestErr)
		return
	}
	if token.Code != 0 {
		c.JSON(nil, ecode.CreativeGeetestErr)
		return
	}
	if code, err = vipSvc.OpenCode(c, arg.Code, arg.Mid); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(code, nil)
}

func belong(c *bm.Context) {
	arg := new(struct {
		Mid int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.Belong(c, arg.Mid))
}

func actives(c *bm.Context) {
	arg := new(struct {
		RelationIDs []string `form:"relationIds,split" `
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.Actives(c, arg.RelationIDs))
}

func webToken(c *bm.Context) {
	c.JSON(vipSvc.WebToken(c))
}

func codeOpened(c *bm.Context) {
	arg := new(model.ArgCodeOpened)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(vipSvc.CodeOpened(c, arg))
}
