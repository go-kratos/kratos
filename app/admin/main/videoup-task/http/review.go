package http

import (
	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

func checkReview(c *bm.Context) {
	sf := &model.SubmitForm{}
	if err := c.Bind(sf); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if ok, err := srv.CheckReview(c, sf); err != nil {
		c.JSON(nil, err)
	} else {
		c.JSON(ok, nil)
	}
}

// list
func listreviews(c *bm.Context) {
	v := &model.ListParser{}
	if err := c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	rcs, count, err := srv.ListReviewConfs(c, v.Unames, v.Bt, v.Et, v.Sort, v.Pn, v.Ps)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSONMap(map[string]interface{}{
		"data":  rcs,
		"pager": &model.Pager{Pn: int(v.Pn), Ps: int(v.Ps), Sum: count},
	}, nil)

}

func addreview(c *bm.Context) {
	uid, uname := getUIDName(c)

	trc := &model.ReviewConf{}
	if err := c.Bind(trc); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	trc.UID = uid
	trc.Uname = uname
	c.JSON(nil, srv.AddReviewConf(c, trc))
}

func editreview(c *bm.Context) {
	uid, uname := getUIDName(c)

	trc := &model.ReviewConf{}
	if err := c.Bind(trc); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}

	trc.UID = uid
	trc.Uname = uname
	c.JSON(nil, srv.EditReviewConf(c, trc))
}

func delreview(c *bm.Context) {
	idStr := c.Request.Form.Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		c.JSON(nil, ecode.RequestErr)
	}

	c.JSON(nil, srv.DelReviewConf(c, id))
}
