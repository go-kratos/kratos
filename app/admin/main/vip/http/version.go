package http

import (
	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func updateVersion(c *bm.Context) {
	var (
		err error
	)
	arg := new(struct {
		ID      int64  `form:"id" validate:"required"`
		Version string `form:"version" `
		Tip     string `form:"tip"`
		Link    string `form:"link"`
	})
	if err = c.Bind(arg); err != nil {
		return
	}
	operator, ok := c.Get("username")
	if !ok {
		c.JSON(nil, ecode.AccessDenied)
		return
	}
	if err = vipSvc.UpdateVersion(c, &model.VipAppVersion{
		ID:       arg.ID,
		Version:  arg.Version,
		Tip:      arg.Tip,
		Operator: operator.(string),
		Link:     arg.Link,
	}); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func versions(c *bm.Context) {
	var (
		res []*model.VipAppVersion
		err error
	)
	if res, err = vipSvc.AllVersion(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
