package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func add(c *bm.Context) {
	v := new(struct {
		MID         int64 `form:"mid"`
		AccountType int   `form:"account_type"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.AddUp(c, v.MID, v.AccountType)
	if err != nil {
		log.Error("growup svr.AddUp error(%v)", err)
	}
	c.JSON(nil, err)
}

func recovery(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.Recovery(c, v.MID)
	if err != nil {
		log.Error("growup svr.Recovery error(%v)", err)
	}
	c.JSON(nil, err)
}

// pgc移除不区分业务
func deleteUp(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	err := svr.DeleteUp(c, v.MID)
	if err != nil {
		log.Error("growup svr.DeleteUp error(%v)", err)
	}
	c.JSON(nil, err)
}

func queryForUps(c *bm.Context) {
	v := new(struct {
		BusinessType int     `form:"business_type"`
		AccountType  int     `form:"account_type"`
		States       []int64 `form:"account_states,split"`
		MID          int64   `form:"mid"`
		Category     int     `form:"category"`
		SignType     int     `form:"sign_type"`
		Nickname     string  `form:"nickname"`
		Lower        int     `form:"lower"`
		Upper        int     `form:"upper"`
		From         int     `form:"from" default:"0" validate:"min=0"`
		Limit        int     `form:"limit" default:"20" validate:"min=1"`
		Sort         string  `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	ups, total, err := svr.QueryFromUpInfo(c, v.BusinessType, v.AccountType, v.States, v.MID, v.Category, v.SignType, v.Nickname, v.Lower, v.Upper, v.From, v.Limit, v.Sort)
	if err != nil {
		log.Error("growup svr.QueryFromUpInfo error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    ups,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func reject(c *bm.Context) {
	v := new(struct {
		Type   int     `form:"type"`
		MIDs   []int64 `form:"mids,split"`
		Reason string  `form:"reason" validate:"required"`
		Days   int     `form:"days"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.Reject(c, v.Type, v.MIDs, v.Reason, v.Days)
	if err != nil {
		log.Error("growup svr.Reject error(%v)", err)
	}
	c.JSON(nil, err)
}

func pass(c *bm.Context) {
	v := new(struct {
		Type int     `form:"type"`
		MIDs []int64 `form:"mids,split"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.Pass(c, v.MIDs, v.Type)
	if err != nil {
		log.Error("growup svr.Pass error(%v)", err)
	}
	c.JSON(nil, err)
}

func dismiss(c *bm.Context) {
	v := new(struct {
		Type   int    `form:"type"`
		MID    int64  `form:"mid" validate:"required"`
		Reason string `form:"reason" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	u, err := c.Request.Cookie("username")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	username := u.Value
	err = svr.Dismiss(c, username, v.Type, 3, v.MID, v.Reason)
	if err != nil {
		log.Error("growup svr.Dismiss error(%v)", err)
	}
	c.JSON(nil, err)
}

func forbid(c *bm.Context) {
	v := new(struct {
		Type   int    `form:"type"`
		MID    int64  `form:"mid" validate:"required"`
		Reason string `form:"reason" validate:"required"`
		Days   int    `form:"days"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	u, err := c.Request.Cookie("username")
	if err != nil {
		c.JSON(nil, err)
		return
	}
	username := u.Value
	err = svr.Forbid(c, username, v.Type, 3, v.MID, v.Reason, v.Days, v.Days*86400)
	if err != nil {
		log.Error("growup svr.Forbid error(%v)", err)
	}
	c.JSON(nil, err)
}

func addToBlocked(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.Block(c, v.MID)
	if err != nil {
		log.Error("growup svr.Block error(%v)", err)
	}
	c.JSON(nil, err)
}

func queryFromBlocked(c *bm.Context) {
	v := new(struct {
		MID      int64  `form:"mid"`
		Category int    `form:"category"`
		Nickname string `form:"nickname"`
		Lower    int    `form:"lower"`
		Upper    int    `form:"upper"`
		From     int    `form:"from" default:"0" validate:"min=0"`
		Limit    int    `form:"limit" default:"20" validate:"min=1"`
		Sort     string `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	ups, total, err := svr.QueryFromBlocked(c, v.MID, v.Category, v.Nickname, v.Lower, v.Upper, v.From, v.Limit, v.Sort)
	if err != nil {
		log.Error("growup svr.QueryFromBlocked error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    ups,
		"paging": map[string]int{
			"page_size": v.Limit,
			"total":     total,
		},
	}))
}

func deleteFromBlocked(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.DeleteFromBlocked(c, v.MID)
	if err != nil {
		log.Error("growup svr.DeleteFromBlocked error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateAccountState(c *bm.Context) {
	v := new(struct {
		MID   int64 `form:"mid"`
		State int   `form:"state"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateUpAccountState(c, "up_info_video", v.MID, v.State)
	if err != nil {
		log.Error("growup svr.UpdateUpAccountState error(%v)", err)
	}
	c.JSON(nil, err)
}

func delUpAccount(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.DelUpAccount(c, v.MID)
	if err != nil {
		log.Error("growup svr.DelUpAccount error(%v)", err)
	}
	c.JSON(nil, err)
}

func updateUpAccount(c *bm.Context) {
	v := new(struct {
		MID          int64  `form:"mid" validate:"required"`
		IsDeleted    int    `form:"is_deleted"`
		WithdrawDate string `form:"withdraw_date" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.UpdateUpAccount(c, v.MID, v.IsDeleted, v.WithdrawDate)
	if err != nil {
		log.Error("growup svr.UpdateUpAccount error(%v)", err)
	}
	c.JSON(nil, err)
}

func recoverCredit(c *bm.Context) {
	v := new(struct {
		Type int   `form:"type"`
		ID   int64 `form:"id" validate:"required"`
		MID  int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	err := svr.RecoverCreditScore(c, v.Type, v.ID, v.MID)
	if err != nil {
		log.Error("growup svr.RecoverCreditScore error(%v)", err)
	}
	c.JSON(nil, err)
}

func creditRecords(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	crs, err := svr.CreditRecords(c, v.MID)
	if err != nil {
		log.Error("growup svr.CreditRecords error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(crs, nil)
}

func exportUps(c *bm.Context) {
	v := new(struct {
		BusinessType int     `form:"business_type"`
		AccountType  int     `form:"account_type"`
		States       []int64 `form:"account_states,split"`
		MID          int64   `form:"mid"`
		Category     int     `form:"category"`
		SignType     int     `form:"sign_type"`
		Nickname     string  `form:"nickname"`
		Lower        int     `form:"lower"`
		Upper        int     `form:"upper"`
		From         int     `form:"from" default:"0" validate:"min=0"`
		Limit        int     `form:"limit" default:"20" validate:"min=1"`
		Sort         string  `form:"sort"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	content, err := svr.ExportUps(c, v.BusinessType, v.AccountType, v.States, v.MID, v.Category, v.SignType, v.Nickname, v.Lower, v.Upper, v.From, v.Limit, v.Sort)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.ExportUps error(%v)", err)
		return
	}

	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "ups"),
	})
}

func upState(c *bm.Context) {
	v := new(struct {
		Type int   `form:"type"`
		MID  int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	data, err := svr.UpState(c, v.MID, v.Type)
	if err != nil {
		log.Error("growup svr.UpState error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(data, nil)
}
