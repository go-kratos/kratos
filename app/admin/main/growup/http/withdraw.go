package http

import (
	"fmt"
	"net/http"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
)

func upWithdraw(c *bm.Context) {
	v := new(struct {
		MIDs      []int64 `form:"mids,split"`
		IsDeleted int     `form:"is_deleted"`
		Page      int     `form:"page" validate:"min=1" default:"1"`
		Size      int     `form:"size" validtae:"min=1" default:"20"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	from := (v.Page - 1) * v.Size
	data, total, err := incomeSvr.UpWithdraw(c, v.MIDs, v.IsDeleted, from, v.Size)
	if err != nil {
		log.Error("growup incomeSvr.UpWithdraw error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result": map[string]interface{}{
			"data":        data,
			"page":        v.Page,
			"total_count": total,
		},
		"status": "success",
	}))
}

func upWithdrawExport(c *bm.Context) {
	v := new(struct {
		MIDs      []int64 `form:"mids,split"`
		IsDeleted int     `form:"is_deleted"`
		Page      int     `form:"page" validate:"min=1" default:"1"`
		Size      int     `form:"size" validtae:"min=1" default:"20"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	from := (v.Page - 1) * v.Size
	content, err := incomeSvr.UpWithdrawExport(c, v.MIDs, v.IsDeleted, from, v.Size)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup incomeSvr.UpWithdrawExport error(%v)", err)
		return
	}

	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "up_withdraw"),
	})
}

func upWithdrawStatis(c *bm.Context) {
	v := new(struct {
		FromTime  int64 `form:"from_time"`
		ToTime    int64 `form:"to_time"`
		IsDeleted int   `form:"is_deleted"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	data, err := incomeSvr.UpWithdrawStatis(c, v.FromTime, v.ToTime, v.IsDeleted)
	if err != nil {
		log.Error("growup incomeSvr.UpWithdrawStatis error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func upWithdrawDetail(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})

	if err := c.Bind(v); err != nil {
		return
	}

	data, err := incomeSvr.UpWithdrawDetail(c, v.MID)
	if err != nil {
		log.Error("growup incomeSvr.UpWithdrawStatis error(%v)", err)
		c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
			"code":   err,
			"status": "fail",
		}))
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"result":  data,
		"status":  "success",
	}))
}

func upWithdrawDetailExport(c *bm.Context) {
	v := new(struct {
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	content, err := incomeSvr.UpWithdrawDetailExport(c, v.MID)
	if err != nil {
		c.JSON(nil, err)
		log.Error("growup svr.UpWithdrawDetailExport error(%v)", err)
		return
	}

	c.Render(http.StatusOK, CSV{
		Content: content,
		Title:   fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), "up_withdraw_detail"),
	})
}
