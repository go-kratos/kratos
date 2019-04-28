package http

import (
	match "go-common/app/interface/main/activity/model/like"
	bm "go-common/library/net/http/blademaster"
)

func matchs(c *bm.Context) {
	p := new(match.ParamSid)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(matchSvc.Match(c, p.Sid))
}

func addGuess(c *bm.Context) {
	var (
		mid int64
		err error
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	p := new(match.ParamAddGuess)
	if err = c.Bind(p); err != nil {
		return
	}
	if p.Stake < 1 {
		p.Stake = 1
	}
	_, err = matchSvc.AddGuess(c, mid, p)
	c.JSON(nil, err)
}

func listGuess(c *bm.Context) {
	var (
		mid int64
	)
	p := new(match.ParamSid)
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	c.JSON(matchSvc.ListGuess(c, p.Sid, mid))
}

func guess(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	p := new(match.ParamSid)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(matchSvc.Guess(c, mid, p))
}

func clearCache(c *bm.Context) {
	p := new(match.ParamMsg)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, matchSvc.ClearCache(c, p.Msg))
}

func addFollow(c *bm.Context) {
	var (
		mid int64
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	p := new(match.ParamTeams)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, matchSvc.AddFollow(c, mid, p.Teams))
}

func follow(c *bm.Context) {
	var (
		mid int64
	)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	c.JSON(matchSvc.Follow(c, mid))
}

func unStart(c *bm.Context) {
	var (
		mid   int64
		err   error
		total int
		list  []*match.Object
	)
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	p := new(match.ParamObject)
	if err = c.Bind(p); err != nil {
		return
	}
	if list, total, err = matchSvc.ObjectsUnStart(c, mid, p); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   p.Pn,
		"size":  p.Ps,
		"total": total,
	}
	data["page"] = page
	data["list"] = list
	c.JSON(data, nil)
}
