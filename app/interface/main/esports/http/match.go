package http

import (
	"time"

	"go-common/app/interface/main/esports/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func filterMatch(c *bm.Context) {
	p := new(model.ParamFilter)
	if err := c.Bind(p); err != nil {
		return
	}
	if p.Stime != "" {
		if _, err := time.Parse("2006-01-02", p.Stime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(eSvc.FilterMatch(c, p))
}

func listContest(c *bm.Context) {
	var (
		mid   int64
		err   error
		total int
		list  []*model.Contest
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	p := new(model.ParamContest)
	if err = c.Bind(p); err != nil {
		return
	}
	if p.Stime != "" {
		if _, err = time.Parse("2006-01-02", p.Stime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Etime != "" {
		if _, err = time.Parse("2006-01-02", p.Etime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if list, total, err = eSvc.ListContest(c, mid, p); err != nil {
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

func appContest(c *bm.Context) {
	var (
		mid   int64
		err   error
		total int
		list  []*model.Contest
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	p := new(model.ParamContest)
	if err = c.Bind(p); err != nil {
		return
	}
	if p.Stime != "" {
		if _, err = time.Parse("2006-01-02", p.Stime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		p.Etime = p.Stime
	}
	if list, total, err = eSvc.ListContest(c, mid, p); err != nil {
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

func calendar(c *bm.Context) {
	var err error
	p := new(model.ParamFilter)
	if err = c.Bind(p); err != nil {
		return
	}
	if _, err = time.Parse("2006-01-02", p.Stime); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = time.Parse("2006-01-02", p.Etime); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(eSvc.Calendar(c, p))
}

func filterVideo(c *bm.Context) {
	p := new(model.ParamFilter)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(eSvc.FilterVideo(c, p))
}

func listVideo(c *bm.Context) {
	var (
		err   error
		total int
		list  []*arcmdl.Arc
	)
	p := new(model.ParamVideo)
	if err = c.Bind(p); err != nil {
		return
	}
	if list, total, err = eSvc.ListVideo(c, p); err != nil {
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

func actVideos(c *bm.Context) {
	var (
		err error
	)
	param := new(struct {
		MmID int64 `form:"mm_id"  validate:"gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(eSvc.ActModules(c, param.MmID))
}

func active(c *bm.Context) {
	var (
		err error
	)
	param := new(struct {
		Aid int64 `form:"aid"  validate:"gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(eSvc.ActPage(c, param.Aid))
}

func actPoints(c *bm.Context) {
	var (
		mid   int64
		err   error
		total int
		list  []*model.Contest
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	p := new(model.ParamActPoint)
	if err = c.Bind(p); err != nil {
		return
	}
	if list, total, err = eSvc.ActPoints(c, mid, p); err != nil {
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

func actKnockout(c *bm.Context) {
	var (
		err error
	)
	param := new(struct {
		MdID int64 `form:"md_id"  validate:"gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(eSvc.ActKnockout(c, param.MdID))
}

func actTop(c *bm.Context) {
	var (
		mid   int64
		err   error
		total int
		list  []*model.Contest
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	p := new(model.ParamActTop)
	if err = c.Bind(p); err != nil {
		return
	}
	if p.Stime != "" {
		if _, err = time.Parse("2006-01-02", p.Stime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if p.Etime != "" {
		if _, err = time.Parse("2006-01-02", p.Etime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if list, total, err = eSvc.ActTop(c, mid, p); err != nil {
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

func search(c *bm.Context) {
	var (
		mid   int64
		buvid string
		err   error
	)
	p := new(model.ParamSearch)
	if err = c.Bind(p); err != nil {
		return
	}
	if ck, err := c.Request.Cookie("buvid3"); err == nil {
		buvid = ck.Value
	}
	if midInter, ok := c.Get("mid"); ok {
		mid = midInter.(int64)
	}
	c.JSON(eSvc.Search(c, mid, p, buvid))
}
func season(c *bm.Context) {
	var (
		err   error
		total int
		list  []*model.Season
	)
	p := new(model.ParamSeason)
	if err = c.Bind(p); err != nil {
		return
	}
	if list, total, err = eSvc.Season(c, p); err != nil {
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

func appSeason(c *bm.Context) {
	var (
		err   error
		total int
		list  []*model.Season
	)
	p := new(model.ParamSeason)
	if err = c.Bind(p); err != nil {
		return
	}
	if list, total, err = eSvc.AppSeason(c, p); err != nil {
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
func contest(c *bm.Context) {
	var (
		mid int64
		err error
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	param := new(struct {
		Cid int64 `form:"cid"  validate:"gt=0"`
	})
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(eSvc.Contest(c, mid, param.Cid))
}

func recent(c *bm.Context) {
	var (
		mid int64
		err error
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	param := &model.ParamCDRecent{}
	if err = c.Bind(param); err != nil {
		return
	}
	c.JSON(eSvc.Recent(c, mid, param))
}
