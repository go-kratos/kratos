package http

import (
	"time"

	"go-common/app/interface/main/esports/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func addFav(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	v := new(struct {
		Cid int64 `form:"cid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, switchCode(eSvc.AddFav(c, mid, v.Cid)))
}

func delFav(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	v := new(struct {
		Cid int64 `form:"cid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(nil, switchCode(eSvc.DelFav(c, mid, v.Cid)))
}

func listFav(c *bm.Context) {
	var (
		mid     int64
		total   int
		contest []*model.Contest
		err     error
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	v := new(struct {
		VMID int64 `form:"vmid"`
		Pn   int   `form:"pn" default:"1" validate:"min=1"`
		Ps   int   `form:"ps" default:"5" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if contest, total, err = eSvc.ListFav(c, mid, v.VMID, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": total,
	}
	data["page"] = page
	data["list"] = contest
	c.JSON(data, nil)
}
func appListFav(c *bm.Context) {
	var (
		mid     int64
		total   int
		contest []*model.Contest
		err     error
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	v := new(model.ParamFav)
	if err = c.Bind(v); err != nil {
		return
	}
	if mid == 0 && v.VMID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if v.Stime != "" {
		if _, err = time.Parse("2006-01-02", v.Stime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if v.Etime != "" {
		if _, err = time.Parse("2006-01-02", v.Etime); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	if contest, total, err = eSvc.ListAppFav(c, mid, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": total,
	}
	data["page"] = page
	data["list"] = contest
	c.JSON(data, nil)
}

func seasonFav(c *bm.Context) {
	var (
		mid     int64
		total   int
		seasons []*model.Season
		err     error
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	v := new(model.ParamSeason)
	if err = c.Bind(v); err != nil {
		return
	}
	if mid == 0 && v.VMID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if seasons, total, err = eSvc.SeasonFav(c, mid, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": total,
	}
	data["page"] = page
	data["list"] = seasons
	c.JSON(data, nil)
}

func stimeFav(c *bm.Context) {
	var (
		mid    int64
		total  int
		stimes []string
		err    error
	)
	if midStr, ok := c.Get("mid"); ok {
		mid = midStr.(int64)
	}
	v := new(model.ParamSeason)
	if err = c.Bind(v); err != nil {
		return
	}
	if mid == 0 && v.VMID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if stimes, total, err = eSvc.StimeFav(c, mid, v); err != nil {
		c.JSON(nil, err)
		return
	}
	data := make(map[string]interface{}, 2)
	page := map[string]int{
		"num":   v.Pn,
		"size":  v.Ps,
		"total": total,
	}
	data["page"] = page
	data["list"] = stimes
	c.JSON(data, nil)
}

func switchCode(err error) error {
	if err == nil {
		return err
	}
	switch ecode.Cause(err) {
	case ecode.FavResourceOverflow:
		err = ecode.EsportsContestMaxCount
	case ecode.FavResourceAlreadyDel:
		err = ecode.EsportsContestFavDel
	case ecode.FavResourceExist:
		err = ecode.EsportsContestFavExist
	case ecode.FavFolderNotExist:
		err = ecode.EsportsContestNotExist
	}
	return err
}
