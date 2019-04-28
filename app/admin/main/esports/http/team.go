package http

import (
	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

func teamInfo(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(esSvc.TeamInfo(c, v.ID))
}

func teamList(c *bm.Context) {
	var (
		list []*model.TeamInfo
		cnt  int64
		err  error
	)
	v := new(struct {
		Pn     int64  `form:"pn" validate:"min=0"`
		Ps     int64  `form:"ps" validate:"min=0,max=30"`
		Title  string `form:"title"`
		Status int    `form:"status"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn == 0 {
		v.Pn = 1
	}
	if v.Ps == 0 {
		v.Ps = 20
	}
	if list, cnt, err = esSvc.TeamList(c, v.Pn, v.Ps, v.Title, v.Status); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int64{
		"num":   v.Pn,
		"size":  v.Ps,
		"count": cnt,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}

func addTeam(c *bm.Context) {
	var (
		err           error
		gids, tmpGids []int64
	)
	v := new(model.Team)
	if err = c.Bind(v); err != nil {
		return
	}
	gidStr := c.Request.Form.Get("gids")
	if gidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tmpGids, err = xstr.SplitInts(gidStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, v := range tmpGids {
		if v > 0 {
			gids = append(gids, v)
		}
	}
	if len(gids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.AddTeam(c, v, gids))
}

func editTeam(c *bm.Context) {
	var (
		err           error
		gids, tmpGids []int64
	)
	v := new(model.Team)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	gidStr := c.Request.Form.Get("gids")
	if gidStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if tmpGids, err = xstr.SplitInts(gidStr); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, v := range tmpGids {
		if v > 0 {
			gids = append(gids, v)
		}
	}
	if len(gids) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.EditTeam(c, v, gids))
}

func forbidTeam(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"min=1"`
		State int   `form:"state" validate:"min=0,max=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.ForbidTeam(c, v.ID, v.State))
}
