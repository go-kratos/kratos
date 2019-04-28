package http

import (
	"go-common/app/admin/ep/marthe/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

// bugly version
func updateVersion(c *bm.Context) {
	var (
		req      = &model.AddVersionRequest{}
		err      error
		username string
	)

	if err = c.BindWith(req, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.UpdateBuglyVersion(c, username, req))
}

func getVersionAndProjectList(c *bm.Context) {
	c.JSON(srv.BuglyVersionAndProjectList(c))
}

func queryVersions(c *bm.Context) {
	var (
		req = &model.QueryBuglyVersionRequest{}
		err error
	)

	if err = req.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	if err = c.BindWith(req, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryBuglyVersions(c, req))
}

func runVersions(c *bm.Context) {
	var (
		v = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}

	c.JSON(srv.RunVersions(v.ID))
}

// bugly cookie
func updateCookie(c *bm.Context) {
	var (
		req      = &model.AddCookieRequest{}
		err      error
		username string
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(req, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.UpdateCookie(c, username, req))
}

func queryCookies(c *bm.Context) {
	var (
		req = &model.QueryBuglyCookiesRequest{}
		err error
	)

	if err = req.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	if err = c.BindWith(req, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryCookies(c, req))
}

func updateCookieStatus(c *bm.Context) {
	var (
		v = new(struct {
			ID     int64 `form:"id"`
			Status int   `form:"status"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, srv.UpdateCookieStatus(c, v.ID, v.Status))
}

// bugly issue
func queryBuglyIssue(c *bm.Context) {
	var (
		err error
		v   = &model.QueryBuglyIssueRequest{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryBuglyIssueRecords(c, v))
}

func queryBatchRun(c *bm.Context) {
	var (
		err error
		v   = &model.QueryBuglyBatchRunsRequest{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryBatchRuns(c, v))
}

// bugly project
func updateProject(c *bm.Context) {
	var (
		req      = &model.AddProjectRequest{}
		err      error
		username string
	)

	if err = c.BindWith(req, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(srv.UpdateBuglyProject(c, username, req))
}

func queryProjects(c *bm.Context) {
	var (
		req = &model.QueryBuglyProjectRequest{}
		err error
	)

	if err = req.Verify(); err != nil {
		c.JSON(nil, err)
		return
	}

	if err = c.BindWith(req, binding.JSON); err != nil {
		return
	}

	c.JSON(srv.QueryBuglyProjects(c, req))
}

func queryProjectVersions(c *bm.Context) {
	var (
		v = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(srv.QueryBuglyProjectVersions(c, v.ID))
}

func queryProject(c *bm.Context) {
	var (
		v = new(struct {
			ID int64 `form:"id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(srv.QueryBuglyProject(c, v.ID))
}

func queryAllProjects(c *bm.Context) {
	c.JSON(srv.QueryAllBuglyProjects(c))
}

//test
func test(c *bm.Context) {
	srv.BatchRunTask(model.TaskBatchRunUpdateBugInTapd, srv.BatchRunUpdateBugInTapd)
	c.JSON(nil, nil)

}
