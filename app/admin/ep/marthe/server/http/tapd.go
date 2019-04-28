package http

import (
	"go-common/app/admin/ep/marthe/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func updateTapdBugTpl(c *bm.Context) {
	var (
		err      error
		v        = &model.UpdateTapdBugTplRequest{}
		username string
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.UpdateTapdBugTpl(c, username, v))
}

func queryTapdBugTpl(c *bm.Context) {
	var (
		err error
		v   = &model.QueryTapdBugTemplateRequest{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryTapdBugTpl(c, v))
}

func queryAllTapdBugTpl(c *bm.Context) {
	c.JSON(srv.QueryAllTapdBugTpl(c))
}

func updateTapdBugVersionTpl(c *bm.Context) {
	var (
		err      error
		v        = &model.UpdateTapdBugVersionTplRequest{}
		username string
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.UpdateTapdBugVersionTpl(c, username, v))
}

func queryTapdBugVersionTpl(c *bm.Context) {
	var (
		err error
		v   = &model.QueryTapdBugVersionTemplateRequest{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryTapdBugVersionTpl(c, v))
}

func checkFilterSql(c *bm.Context) {
	var (
		err error
		v   = new(struct {
			SQL string `json:"issue_filter_sql"`
		})
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.CheckTapdBugTplSQL(c, v.SQL))
}

func bugInsertTapdWithProject(c *bm.Context) {
	var (
		v = new(struct {
			ID int64 `form:"id"`
		})
		err      error
		username string
	)

	if err = c.Bind(v); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.BugInsertTapdWithProject(c, v.ID, username))
}

func bugInsertTapdWithVersion(c *bm.Context) {
	var (
		v = new(struct {
			ID int64 `form:"id"`
		})
		err      error
		username string
	)

	if err = c.Bind(v); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}
	c.JSON(srv.BugInsertTapdWithVersion(c, v.ID, username))
}

func queryBugRecord(c *bm.Context) {
	var (
		err error
		v   = &model.QueryBugRecordsRequest{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryBugRecords(c, v))
}

func updateTapdBugPriorityConf(c *bm.Context) {
	var (
		err      error
		v        = &model.UpdateTapdBugPriorityConfRequest{}
		username string
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.UpdateTapdBugPriorityConf(c, username, v))
}

func queryTapdBugPriorityConf(c *bm.Context) {
	var (
		err error
		v   = &model.QueryTapdBugPriorityConfsRequest{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryTapdBugPriorityConfsRequest(c, v))
}

func checkAuth(c *bm.Context) {
	var (
		v = new(struct {
			Username    string `form:"username"`
			WorkspaceID string `form:"workspace_id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(srv.HttpAccessToWorkspace(c, v.WorkspaceID, v.Username))
}
