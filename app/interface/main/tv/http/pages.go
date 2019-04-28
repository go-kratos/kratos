package http

import (
	"strconv"

	"go-common/app/interface/main/tv/dao/thirdp"
	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func zonePage(c *bm.Context) {
	var (
		t   int
		err error
		req = c.Request.Form
	)
	takeBuild(req) // take build number
	seasonType := req.Get("season_type")
	if seasonType == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if t, err = strconv.Atoi(seasonType); err != nil || t < 1 {
		t = 1
	}
	zone, ok := tvSvc.RankData[t]
	if !ok {
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(zone, nil)
}

func homepage(c *bm.Context) {
	var err error
	params := c.Request.Form
	accessKey := params.Get("access_key")
	homeData := *tvSvc.HomeData
	takeBuild(params) // take build number
	if accessKey != "" {
		_, ok := c.Get("mid")
		if !ok { // if not login, we don't call follow data
			err = ecode.NoLogin
		} else {
			homeData.Follow = tvSvc.FollowData(c, accessKey)
		}
	}
	if tvSvc.HomeData == nil {
		log.Error("HomeData is Nil")
		c.JSON(nil, ecode.ServerErr)
		return
	}
	c.JSON(homeData, err)
}

func modpage(c *bm.Context) {
	var (
		err    error
		params = c.Request.Form
	)
	v := new(model.ReqPageFollow)
	if err = c.Bind(v); err != nil {
		return
	}
	takeBuild(params)
	// check login
	if v.AccessKey != "" {
		_, ok := c.Get("mid")
		if !ok { // if not logged in, not request follow
			err = ecode.NoLogin
			v.AccessKey = ""
		}
	}
	c.JSON(tvSvc.PageFollow(c, v))
}

// get dangbei pgc data by page
func dbeiPage(c *bm.Context) {
	v := new(struct {
		Page  int64  `form:"page" validate:"required"`
		TypeC string `form:"type_c" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	// check the typeC ( type of the content ) value
	if _, err = thirdp.KeyThirdp(v.TypeC); err != nil {
		c.JSON(nil, ecode.TvDangbeiWrongType)
		return
	}
	c.JSON(thirdpSvc.PickDBeiPage(v.Page, v.TypeC))
}

func mangoRecom(c *bm.Context) {
	c.JSON(thirdpSvc.MangoRecom(), nil)
}

func mangoSnPage(c *bm.Context) {
	v := new(struct {
		Page int64 `form:"page" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(thirdpSvc.MangoSns(c, v.Page))
}

func mangoEpPage(c *bm.Context) {
	v := new(struct {
		SID  int64 `form:"sid" validate:"required"`
		Page int   `form:"page" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(thirdpSvc.MangoEps(c, v.SID, v.Page))
}

func mangoArcPage(c *bm.Context) {
	v := new(struct {
		Page int64 `form:"page" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(thirdpSvc.MangoArcs(c, v.Page))
}

func mangoVideoPage(c *bm.Context) {
	v := new(struct {
		AVID int64 `form:"avid" validate:"required"`
		Page int   `form:"page" validate:"required"`
	})
	err := c.Bind(v)
	if err != nil {
		return
	}
	c.JSON(thirdpSvc.MangoVideos(c, v.AVID, v.Page))
}
