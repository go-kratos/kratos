package http

import (
	"go-common/library/ecode"

	bm "go-common/library/net/http/blademaster"
)

//staff only
func webApplySubmit(c *bm.Context) {
	v := new(struct {
		ID           int64 `form:"id" validate:"required"`
		AID          int64 `form:"aid" validate:"required"`
		Type         int64 `form:"type" validate:"required"`
		MID          int64 `form:"mid" validate:"required"`
		FlagAddBlack int   `form:"flag_add_black" `
		FlagRejuse   int   `form:"flag_rejuse"`
		State        int64 `form:"state" validate:"min=1,max=5"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(nil, arcSvc.StaffApplySubmit(c, v.ID, v.AID, v.MID, v.State, v.Type, v.FlagAddBlack, v.FlagRejuse))
}

//staff only
func webApplyCreate(c *bm.Context) {
	v := new(struct {
		AID int64 `form:"aid" validate:"required"`
		MID int64 `form:"mid" validate:"required"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	// check user
	_, ok := c.Get("mid")
	if !ok {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	c.JSON(nil, arcSvc.StaffApplySubmit(c, 0, v.AID, v.MID, 0, 4, 0, 0))
}
