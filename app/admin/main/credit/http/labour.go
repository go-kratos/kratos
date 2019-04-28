package http

import (
	"go-common/app/admin/main/credit/model/blocked"
	bm "go-common/library/net/http/blademaster"
)

func operQuestion(c *bm.Context) {
	v := new(struct {
		IDS    []int64 `form:"ids,split" validate:"min=1,max=20"`
		Status int8    `form:"status" validate:"min=0,max=1" default:"1"`
		OID    int     `form:"oper_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if err := creSvc.DB.Model(&blocked.LabourQuestion{}).Where(v.IDS).Updates(map[string]interface{}{"status": v.Status, "oper_id": v.OID}).Error; err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}

func delQuestion(c *bm.Context) {
	v := new(struct {
		IDS    []int64 `form:"ids,split" validate:"min=1,max=20"`
		Status int8    `form:"status" validate:"min=0,max=1" default:"1"`
		OID    int     `form:"oper_id" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if err := creSvc.DB.Model(&blocked.LabourQuestion{}).Where(v.IDS).Updates(map[string]interface{}{"isdel": v.Status, "oper_id": v.OID}).Error; err != nil {
		httpCode(c, err)
		return
	}
	httpCode(c, nil)
}
