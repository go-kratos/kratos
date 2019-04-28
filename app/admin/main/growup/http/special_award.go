package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/app/admin/main/growup/model"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func awardAdd(c *blademaster.Context) {
	username, _, err := checkCookie(c)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	arg := new(model.AddAwardArg)
	if err = c.BindWith(arg, binding.JSON); err != nil {
		return
	}
	awardID, err := svr.AddAward(c, arg, username)
	if err != nil {
		log.Error("svr.AddAward err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(awardID, nil)
}

func awardUpdate(c *blademaster.Context) {
	arg := new(model.SaveAwardArg)
	if err := c.BindWith(arg, binding.JSON); err != nil {
		return
	}
	err := svr.UpdateAward(c, arg)
	if err != nil {
		log.Error("svr.UpdateAward err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(true, nil)
}

func awardList(c *blademaster.Context) {
	arg := new(struct {
		From  int `form:"from" validate:"min=0" default:"0"`
		Limit int `form:"limit" validate:"min=1" default:"20"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	total, list, err := svr.ListAward(c, arg.From, arg.Limit)
	if err != nil {
		log.Error("svr.ListAward err(%v)", err)
		c.JSON(nil, err)
		return
	}
	renderPagRes(list, total, arg.Limit)(c)
}

func awardDetail(c *blademaster.Context) {
	arg := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	data, err := svr.DetailAward(c, arg.AwardID)
	if err != nil {
		log.Error("svr.DetailAward err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func awardWinnerList(c *blademaster.Context) {
	arg := new(model.QueryAwardWinnerArg)
	if err := c.Bind(arg); err != nil {
		return
	}
	total, list, err := svr.ListAwardWinner(c, arg)
	if err != nil {
		log.Error("svr.ListAwardRecord err(%v)", err)
		c.JSON(nil, err)
		return
	}
	renderPagRes(list, total, arg.Limit)(c)
}

func awardWinnerExport(c *blademaster.Context) {
	arg := new(model.QueryAwardWinnerArg)
	if err := c.Bind(arg); err != nil {
		return
	}
	content, err := svr.ExportAwardWinner(c, arg)
	if err != nil {
		log.Error("svr.ExportAwardWinner err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "award_winner"),
	})
}

func awardResult(c *blademaster.Context) {
	arg := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	data, err := svr.AwardResult(c, arg.AwardID)
	if err != nil {
		log.Error("svr.AwardResult err(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}

func awardResultSave(c *blademaster.Context) {
	arg := new(model.AwardResult)
	if err := c.BindWith(arg, binding.JSON); err != nil {
		return
	}
	if arg.AwardID == 0 {
		c.JSON(false, ecode.Error(ecode.RequestErr, "illegal award_id"))
		return
	}
	err := svr.SaveAwardResult(c, arg)
	if err != nil {
		log.Error("svr.SaveAwardResult err(%v)", err)
		c.JSON(false, err)
		return
	}
	c.JSON(true, nil)
}

func awardWinnerReplace(c *blademaster.Context) {
	arg := new(struct {
		AwardID int64 `form:"award_id" validate:"required"`
		PrevMID int64 `form:"prev_mid" validate:"required"`
		MID     int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	err := svr.ReplaceAwardWinner(c, arg.AwardID, arg.PrevMID, arg.MID)
	if err != nil {
		log.Error("svr.ReplaceAwardWinner err(%v)", err)
		c.JSON(false, err)
		return
	}
	c.JSON(true, nil)
}
