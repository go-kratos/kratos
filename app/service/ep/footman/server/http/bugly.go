package http

import (
	bm "go-common/library/net/http/blademaster"
)

func queryIssue(c *bm.Context) {
	v := new(struct {
		Version string `form:"version"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	c.JSON(srv.Issue(c, v.Version))
}

func updateToken(c *bm.Context) {
	c.JSON(srv.UpdateToken(c))
}

func saveIssue(c *bm.Context) {
	v := new(struct {
		Version string `form:"version"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	c.JSON(srv.SaveIssue(c, v.Version))
}

func saveIssues(c *bm.Context) {
	c.JSON(srv.SaveIssues(c))
}
