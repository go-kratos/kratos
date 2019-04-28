package http

import (
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func limitUserList(c *bm.Context) {
	var (
		err      error
		total    int64
		userList []*model.LimitUser
		param    = new(model.ParamPage)
	)
	if err = c.Bind(param); err != nil {
		return
	}
	if param.Pn <= 0 {
		param.Pn = model.DefaultPageNum
	}
	if param.Ps <= 0 {
		param.Ps = model.DefaultPagesize
	}
	if userList, total, err = svc.LimitUsers(c, param.Pn, param.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	data["users"] = userList
	data["page"] = map[string]interface{}{
		"page":     param.Pn,
		"pagesize": param.Ps,
		"total":    total,
	}
	c.JSON(data, nil)
}

func limitUserAdd(c *bm.Context) {
	var (
		err      error
		username string
		param    = new(struct {
			Mid int64 `form:"mid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	_, username = managerInfo(c)
	if username == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, svc.LimitUserAdd(c, param.Mid, username))
}

func limitUserDel(c *bm.Context) {
	var (
		err   error
		param = new(struct {
			Mid int64 `form:"mid" validate:"required,gt=0"`
		})
	)
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(nil, svc.LimitUserDel(c, param.Mid))
}
