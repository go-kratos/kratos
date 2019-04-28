package http

import (
	"go-common/app/interface/main/creative/model/app"
	mMdl "go-common/app/interface/main/creative/model/music"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
	"strconv"
)

// prepare ext data for bgm-detail pager
func appBgmExt(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	sidStr := params.Get("sid")
	var (
		sid int64
		err error
		ext *mMdl.BgmExt
	)
	sid, err = strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	ext, err = musicSvc.BgmExt(c, mid, sid)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(ext, nil)
}

func appBgmPre(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	fromStr := params.Get("from")
	from, _ := strconv.Atoi(fromStr)
	if _, ok := app.BgmFrom[from]; !ok {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(map[string]interface{}{
		"typelist": musicSvc.PreByFrom(c, from),
		"hotword":  musicSvc.Hotwords,
	}, nil)
}

func appBgmList(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	tidStr := params.Get("tid")
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(musicSvc.BgmList(c, tid), nil)
}

func arcBgmList(c *bm.Context) {
	v := new(struct {
		AID   int64 `form:"aid" validate:"required"`
		CID   int64 `form:"cid" validate:"required"`
		Cache bool  `form:"cache"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	res, err := pubSvc.BgmBindList(c, v.AID, v.CID, 3, true)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(res, nil)
}

func appBgmSearch(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	pnStr := params.Get("pn")
	psStr := params.Get("ps")
	kw := params.Get("kw")
	pn, _ := strconv.Atoi(pnStr)
	if pn < 1 {
		pn = 1
	}
	ps, _ := strconv.Atoi(psStr)
	if ps > 20 || ps < 0 {
		ps = 20
	}
	c.JSON(musicSvc.BgmSearch(c, kw, mid, pn, ps), nil)
}

func appBgmView(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	sidStr := params.Get("sid")
	sid, err := strconv.ParseInt(sidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(musicSvc.BgmView(c, sid), nil)
}

func appMaterialPre(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	buildStr := params.Get("build")
	build, _ := strconv.Atoi(buildStr)
	plat := params.Get("platform")
	c.JSON(musicSvc.MaterialPre(c, mid, plat, build), nil)
}

func appMaterial(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	idStr := params.Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	tpStr := params.Get("type")
	tp, _ := strconv.Atoi(tpStr)
	c.JSONMap(map[string]interface{}{
		"data": musicSvc.Material(c, id, int8(tp), mid),
	}, nil)
}

func appMissionByType(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	params := c.Request.Form
	tidStr := params.Get("tid")
	tid, _ := strconv.Atoi(tidStr)
	actWithTP, _ := arcSvc.MissionOnlineByTid(c, int16(tid), 1)
	c.JSON(actWithTP, nil)
}

func appH5BgmFeedback(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if mid <= 0 {
		c.JSON(nil, ecode.CreativeNotLogin)
		return
	}
	var err error
	v := new(struct {
		Name      string `form:"name"`
		Musicians string `form:"musicians"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if len(v.Name) == 0 && len(v.Musicians) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(musicSvc.AddBgmFeedBack(c, v.Name, v.Musicians, "h5", mid), nil)
}
