package http

import (
	"go-common/app/admin/main/esports/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const _special = 1

func contestInfo(c *bm.Context) {
	v := new(struct {
		ID int64 `form:"id" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(esSvc.ContestInfo(c, v.ID))
}

func contestList(c *bm.Context) {
	var (
		list []*model.ContestInfo
		cnt  int64
		err  error
	)
	v := new(struct {
		Pn   int64 `form:"pn" validate:"min=0"`
		Ps   int64 `form:"ps" validate:"min=0,max=50"`
		Mid  int64 `form:"mid"`
		Sid  int64 `form:"sid"`
		Sort int64 `form:"sort"`
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
	if list, cnt, err = esSvc.ContestList(c, v.Mid, v.Sid, v.Pn, v.Ps, v.Sort); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int64{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": cnt,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}

func addContest(c *bm.Context) {
	var (
		err           error
		tmpGids, gids []int64
	)
	v := new(model.Contest)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Special == _special && v.SpecialName == "" {
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
	if v.DataType > 0 && v.MatchID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.AddContest(c, v, gids))
}

func editContest(c *bm.Context) {
	var (
		err           error
		tmpGids, gids []int64
	)
	v := new(model.Contest)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Special == _special && v.SpecialName == "" {
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
	if v.DataType > 0 && v.MatchID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, esSvc.EditContest(c, v, gids))
}

func forbidContest(c *bm.Context) {
	v := new(struct {
		ID    int64 `form:"id" validate:"min=1"`
		State int   `form:"state" validate:"min=0,max=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, esSvc.ForbidContest(c, v.ID, v.State))
}
