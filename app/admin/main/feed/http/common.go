package http

import (
	"strconv"

	"go-common/app/admin/main/feed/model/common"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func managerInfo(c *bm.Context) (uid int64, username string) {
	if nameInter, ok := c.Get("username"); ok {
		username = nameInter.(string)
	}
	if uidInter, ok := c.Get("uid"); ok {
		uid = uidInter.(int64)
	}
	if username == "" {
		cookie, err := c.Request.Cookie("username")
		if err != nil {
			log.Error("managerInfo get cookie error (%v)", err)
			return
		}
		username = cookie.Value
		c, err := c.Request.Cookie("uid")
		if err != nil {
			log.Error("managerInfo get cookie error (%v)", err)
			return
		}
		uidInt, _ := strconv.Atoi(c.Value)
		uid = int64(uidInt)
	}
	return
}

func cardPreview2(c *bm.Context) {
	var (
		err   error
		title string
		res   = map[string]interface{}{}
	)
	type Card struct {
		Type string `form:"type" validate:"required"`
		ID   int64  `form:"id" validate:"required"`
	}
	card := &Card{}
	if err = c.Bind(card); err != nil {
		return
	}
	if title, err = commonSvc.CardPreview(c, card.Type, card.ID); err != nil {
		res["message"] = err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	titleReturn := common.CardPreview{
		Title: title,
	}
	c.JSON(titleReturn, nil)
}

func actionLog(c *bm.Context) {
	var (
		res = map[string]interface{}{}
	)
	type Log struct {
		Type      int64  `form:"module" validate:"required"`
		Uame      string `form:"uname"`
		CtimeFrom string `form:"ctime_from"`
		CtimeTo   string `form:"ctime_to"`
		Ps        int64  `form:"ps" default:"20"`
		Pn        int64  `form:"pn" default:"1"`
	}
	log := &Log{}
	if err := c.Bind(log); err != nil {
		return
	}
	searchRes, err := commonSvc.LogAction(c, log.Type, log.Ps, log.Pn, log.CtimeFrom, log.CtimeTo, log.Uame)
	if err != nil {
		res["message"] = err.Error()
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	res["data"] = searchRes.Item
	res["pager"] = searchRes.Page
	c.JSONMap(res, nil)
}

func cardType(c *bm.Context) {
	var (
		res = map[string]interface{}{}
	)
	res["data"] = commonSvc.CardType()
	c.JSONMap(res, nil)
}
