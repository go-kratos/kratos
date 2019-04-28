package http

import (
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func sepcialAwardInfo(c *bm.Context) {
	var mid int64
	midI, ok := c.Get("mid")
	if ok {
		mid, _ = midI.(int64)
	}
	data, err := svc.SpecialAwardInfo(c, mid)
	if err != nil {
		log.Error("svc.SpecialAwardInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func specialAwardDetail(c *bm.Context) {
	v := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := svc.AwardDetail(c, v.AwardID)
	if err != nil {
		log.Error("svc.AwardDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func listSpecialAward(c *bm.Context) {
	data, err := svc.AwardList(c)
	if err != nil {
		log.Error("svc.AwardDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func specialAwardRecord(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	data, err := svc.GetWinningRecord(c, mid)
	if err != nil {
		log.Error("svc.AwardDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func specialAwardPoster(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	v := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
		PrizeID int64 `form:"prize_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := svc.GetWinningPoster(c, mid, v.AwardID, v.PrizeID)
	if err != nil {
		log.Error("svc.AwardDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func joinSpecialAward(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)
	v := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svc.JoinAward(c, mid, v.AwardID)
	if err != nil {
		log.Error("svc.AwardDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func specialAwardWinners(c *bm.Context) {
	v := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := svc.Winners(c, v.AwardID)
	if err != nil {
		log.Error("svc.AwardDetail error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func specialAwardUpStatus(c *bm.Context) {
	midI, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	mid, _ := midI.(int64)

	v := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, err := svc.GetAwardUpStatus(c, v.AwardID, mid)
	if err != nil {
		log.Error("svc.GetAwardUpStatus error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
