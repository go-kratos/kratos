package http

import (
	"strconv"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func sagaUserList(c *bm.Context) {
	c.JSON(srv.SagaUserList(c))
}

func runnerConfig(c *bm.Context) {
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryAllConfigFile(c, session.Value, false))
}

func sagaConfig(c *bm.Context) {
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryAllConfigFile(c, session.Value, true))
}

func publicSagaConfig(c *bm.Context) {
	req := new(model.TagUpdate)
	if err := c.Bind(req); err != nil {
		return
	}
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var user string
	if user, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.PublicConfig(c, session.Value, user, req.Names, req.Mark, true))
}

func existConfigSaga(c *bm.Context) {
	var (
		err       error
		projectID int
	)

	if projectID, err = strconv.Atoi(c.Request.Form.Get("project_id")); err != nil {
		return
	}
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.QueryProjectSagaConfig(c, session.Value, projectID))
}

func releaseSagaConfig(c *bm.Context) {
	var (
		err  error
		user string
	)

	v := new(model.ConfigList)
	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	if user, err = getUsername(c); err != nil {
		c.JSON(nil, err)
		return
	}
	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(srv.ReleaseSagaConfig(c, user, session.Value, v))
}

func optionSaga(c *bm.Context) {

	session, err := c.Request.Cookie("_AJSESSIONID")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	projectID := c.Request.Form.Get("project_id")
	log.Info("=====optionSaga projectID: %s", projectID)
	c.JSON(srv.OptionSaga(c, projectID, session.Value))
}
