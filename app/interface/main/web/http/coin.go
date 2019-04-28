package http

import (
	"strconv"
	"time"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
)

func coins(c *bm.Context) {
	var (
		mid, aid int64
		err      error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	aidStr := params.Get("aid")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(webSvc.Coins(c, mid, aid))
}

func addCoin(c *bm.Context) {
	var (
		mid, aid, multiply int64
		err                error
		like               bool
		actLike            = "cointolike"
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	multiStr := params.Get("multiply")
	avTypeStr := params.Get("avtype")
	business := params.Get("business")
	upIDStr := params.Get("upid")
	selectLikeStr := params.Get("select_like")
	selectLike, _ := strconv.Atoi(selectLikeStr)
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	ck := c.Request.Header.Get("Cookie")
	if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil || aid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if multiply, err = strconv.ParseInt(multiStr, 10, 64); err != nil || multiply <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	upID, _ := strconv.ParseInt(upIDStr, 10, 64)
	avtype, _ := strconv.ParseInt(avTypeStr, 10, 64)
	if avtype == 0 {
		avtype = model.CoinAddArcType
	}
	if avtype != model.CoinAddArcType && avtype != model.CoinAddArtType {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if avtype == model.CoinAddArtType && upID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if business != "" {
		if business == model.CoinArcBusiness {
			avtype = model.CoinAddArcType
		} else if business == model.CoinArtBusiness {
			avtype = model.CoinAddArtType
		} else {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	} else {
		switch avtype {
		case model.CoinAddArcType:
			business = model.CoinArcBusiness
		case model.CoinAddArtType:
			business = model.CoinArtBusiness
		default:
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	ua := c.Request.UserAgent()
	refer := c.Request.Referer()
	if like, err = webSvc.AddCoin(c, aid, mid, upID, multiply, avtype, business, ck, ua, refer, time.Now(), selectLike); err == nil {
		if webSvc.CheatInfoc != nil {
			itemType := infoc.ItemTypeAv
			if avtype == model.CoinAddArtType {
				itemType = "article"
			}
			webSvc.CheatInfoc.InfoAntiCheat2(c, strconv.FormatInt(upID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), itemType, "coin", "")
			if like {
				webSvc.CheatInfoc.InfoAntiCheat2(c, strconv.FormatInt(upID, 10), strconv.FormatInt(aid, 10), strconv.FormatInt(mid, 10), strconv.FormatInt(aid, 10), "av", actLike, "")
			}
		}
	}
	c.JSON(struct {
		Like bool `json:"like"`
	}{Like: like}, err)
}

func coinExp(c *bm.Context) {
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	c.JSON(webSvc.CoinExp(c, mid))
}

func coinList(c *bm.Context) {
	var (
		ls    []*model.CoinArc
		count int
		err   error
	)
	v := new(struct {
		Mid int64 `form:"mid" validate:"min=1"`
		Pn  int   `form:"pn" default:"1" validate:"min=1"`
		Ps  int   `form:"ps" default:"20" validate:"min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if ls, count, err = webSvc.CoinList(c, v.Mid, v.Pn, v.Ps); err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSONMap(map[string]interface{}{
		"count": count,
		"data":  ls,
	}, err)
}
