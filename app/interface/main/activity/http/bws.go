package http

import (
	"strconv"

	"go-common/app/interface/main/activity/model/bws"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func user(c *bm.Context) {
	var loginMid int64
	v := new(struct {
		Bid int64  `form:"bid" validate:"min=1"`
		Mid int64  `form:"mid"`
		Key string `form:"key"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if midInter, ok := c.Get("mid"); ok {
		loginMid = midInter.(int64)
	}
	if v.Mid == 0 {
		v.Mid = loginMid
	}
	if v.Mid == 0 && v.Key == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(bwsSvc.User(c, v.Bid, v.Mid, v.Key))
}

func points(c *bm.Context) {
	p := new(bws.ParamPoints)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(bwsSvc.Points(c, p))
}

func point(c *bm.Context) {
	p := new(bws.ParamID)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(bwsSvc.Point(c, p))
}

func achievements(c *bm.Context) {
	p := new(bws.ParamID)
	if err := c.Bind(p); err != nil {
		return
	}
	if p.Day != "" {
		var (
			day int64
			err error
		)
		if day, err = strconv.ParseInt(p.Day, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if day < 20180719 || day > 20180722 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	c.JSON(bwsSvc.Achievements(c, p))
}

func achievement(c *bm.Context) {
	p := new(bws.ParamID)
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(bwsSvc.Achievement(c, p))
}

func unlock(c *bm.Context) {
	v := new(bws.ParamUnlock)
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	if v.Mid == 0 && v.Key == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, bwsSvc.Unlock(c, mid, v))
}

func binding(c *bm.Context) {
	p := new(bws.ParamBinding)
	if err := c.Bind(p); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, bwsSvc.Binding(c, mid, p))
}

func award(c *bm.Context) {
	p := new(bws.ParamAward)
	if err := c.Bind(p); err != nil {
		return
	}
	if p.Mid == 0 && p.Key == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(nil, bwsSvc.Award(c, mid, p))
}

func lottery(c *bm.Context) {
	p := new(struct {
		Bid int64  `form:"bid" validate:"min=1"`
		Aid int64  `form:"aid" validate:"min=1"`
		Day string `form:"day"`
	})
	if err := c.Bind(p); err != nil {
		return
	}
	if p.Day != "" {
		var (
			day int64
			err error
		)
		if day, err = strconv.ParseInt(p.Day, 10, 64); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if day < 20180719 || day > 20180722 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(bwsSvc.Lottery(c, p.Bid, mid, p.Aid, p.Day))
}

func redisInfo(c *bm.Context) {
	v := new(struct {
		Mid  int64  `form:"mid"`
		Key  string `form:"key"`
		Type string `form:"type" validate:"required"`
		Day  string `form:"day"`
		Del  int    `form:"del"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	loginMid := midStr.(int64)
	c.JSON(bwsSvc.RedisInfo(c, loginMid, v.Mid, v.Key, v.Day, v.Type, v.Del))
}

func keyInfo(c *bm.Context) {
	v := new(struct {
		ID   int64  `form:"id"`
		Mid  int64  `form:"mid"`
		Key  string `form:"key"`
		Type string `form:"type" validate:"required"`
		Del  int    `form:"del"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	loginMid := midStr.(int64)
	c.JSON(bwsSvc.KeyInfo(c, loginMid, v.ID, v.Mid, v.Key, v.Type, v.Del))
}

func lotteryCheck(c *bm.Context) {
	v := new(struct {
		Aid int64  `form:"aid" validate:"min=1"`
		Day string `form:"day"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(bwsSvc.LotteryCheck(c, mid, v.Aid, v.Day))
}

func adminInfo(c *bm.Context) {
	v := new(struct {
		Bid int64 `form:"bid" validate:"min=1"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(bwsSvc.AdminInfo(c, v.Bid, mid))
}
