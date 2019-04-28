package http

import (
	"go-common/app/admin/ep/tapd/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

const (
	_sessUnKey = "username"
)

func tapdCallback(c *bm.Context) {
	c.JSON(nil, svc.TapdCallBack(c, c.Request.Body))
}

func updateHook(c *bm.Context) {
	var (
		err      error
		v        = &model.HookURLUpdateReq{}
		username string
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.UpdateHookURL(c, username, v))
}

func queryHook(c *bm.Context) {
	var (
		v   = &model.QueryHookURLReq{}
		err error
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if err = v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(svc.QueryHookURL(c, v))
}

func queryURLEvent(c *bm.Context) {
	var (
		v = new(struct {
			ID int64 `form:"url_id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(svc.QueryURLEvent(c, v.ID))
}

func queryEventLog(c *bm.Context) {
	var (
		v   = &model.QueryEventLogReq{}
		err error
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if err = v.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(svc.QueryEventLog(c, v))
}

func saveHookUrlInCache(c *bm.Context) {
	c.JSON(svc.SaveEnableHookURL(c))
}

func queryHookUrlInCache(c *bm.Context) {
	c.JSON(svc.QueryEnableHookURLInCache(c))
}

func test(c *bm.Context) {
	var (
		v = new(struct {
			Data interface{} `json:"data"`
			Code int         `json:"code"`
		})
		err error
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	log.Info("WorkspaceID [%d]", v.Data)
	c.JSON(v, err)
}

func testform(c *bm.Context) {
	var (
		v = new(struct {
			ID string `form:"event"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}

	log.Info("fengyifeng WorkspaceID", v.ID)
	c.JSON(v, err)
}

func getUsername(c *bm.Context) (username string, err error) {
	user, exist := c.Get(_sessUnKey)
	if !exist {
		err = ecode.AccessKeyErr
		c.JSON(nil, err)
		return
	}
	username = user.(string)
	return
}
