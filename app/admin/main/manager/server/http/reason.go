package http

import (
	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// cateSecExtList .
func cateSecExtList(c *bm.Context) {
	arg := new(model.CateSecExt)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(mngSvc.CateSecExtList(c, arg))
}

// AssociationList .
func associationList(c *bm.Context) {
	// Display all record
	arg := new(model.Association)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.BusinessID < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.AssociationList(c, arg))
}

// addCateSecExt .
func addCateSecExt(c *bm.Context) {
	arg := new(model.CateSecExt)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AddCateSecExt(c, arg))
}

// updateCateSecExt .
func updateCateSecExt(c *bm.Context) {
	arg := new(model.CateSecExt)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateCateSecExt(c, arg))
}

// banCateSecExt .
func banCateSecExt(c *bm.Context) {
	arg := new(model.CateSecExt)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.BanCateSecExt(c, arg))
}

// addAssociation .
func addAssociation(c *bm.Context) {
	arg := new(model.Association)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AddAssociation(c, arg))
}

// updateAssociation .
func updateAssociation(c *bm.Context) {
	arg := new(model.Association)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateAssociation(c, arg))
}

// banAssocaition .
func banAssociation(c *bm.Context) {
	arg := new(model.Association)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.BanAssociation(c, arg))
}

// addReason .
func addReason(c *bm.Context) {
	arg := new(model.Reason)
	if err := c.Bind(arg); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		arg.UID = uid.(int64)
	}
	c.JSON(nil, mngSvc.AddReason(c, arg))
}

// updateReason .
func updateReason(c *bm.Context) {
	arg := new(model.Reason)
	if err := c.Bind(arg); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		arg.UID = uid.(int64)
	}
	if arg.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateReason(c, arg))
}

// reasonList .
func reasonList(c *bm.Context) {
	arg := new(model.SearchReasonParams)
	if err := c.Bind(arg); err != nil {
		return
	}
	data, total, err := mngSvc.ReasonList(c, arg)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	page := map[string]int64{
		"num":   arg.PN,
		"size":  arg.PS,
		"total": total,
	}
	c.JSON(map[string]interface{}{
		"page": page,
		"data": data,
	}, err)
}

// batchUpdateReasonState .
func batchUpdateReasonState(c *bm.Context) {
	arg := new(model.BatchUpdateReasonState)
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(nil, mngSvc.BatchUpdateReasonState(c, arg))
}

// dropList .
func dropDownList(c *bm.Context) {
	arg := new(model.Association)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.BusinessID < 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.DropDownList(c, arg))
}

// businessAttr .
func businessAttr(c *bm.Context) {
	arg := &model.BusinessAttr{}
	if err := c.Bind(arg); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.BusinessAttr(c, arg))
}
