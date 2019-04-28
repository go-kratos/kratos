package http

import (
	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

// addType .
func addType(c *bm.Context) {
	tt := &model.TagType{}
	if err := c.Bind(tt); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AddType(c, tt))
}

// updateType .
func updateType(c *bm.Context) {
	tt := &model.TagType{}
	if err := c.Bind(tt); err != nil {
		return
	}
	c.JSON(nil, mngSvc.UpdateType(c, tt))
}

// deleteType .
func deleteType(c *bm.Context) {
	td := &model.TagTypeDel{}
	res := map[string]interface{}{}
	if err := c.Bind(td); err != nil {
		return
	}
	if td.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := mngSvc.DeleteType(c, td); err != nil {
		if err == ecode.ManagerTagTypeDelErr {
			res["message"] = "这个类型下可是关联了Tag的哦，你真的考虑清楚了吗???"
			c.JSONMap(res, ecode.ManagerTagTypeDelErr)
			return
		}
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)

}

// addTag .
func addTag(c *bm.Context) {
	t := &model.Tag{}
	if err := c.Bind(t); err != nil {
		return
	}
	if uid, exists := c.Get("uid"); exists {
		t.UID = uid.(int64)
	}
	c.JSON(nil, mngSvc.AddTag(c, t))
}

// updateTag .
func updateTag(c *bm.Context) {
	t := &model.Tag{}
	if err := c.Bind(t); err != nil {
		return
	}
	if t.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateTag(c, t))
}

// batchUpdateState .
func batchUpdateState(c *bm.Context) {
	b := &model.BatchUpdateState{}
	if err := c.Bind(b); err != nil {
		return
	}
	c.JSON(nil, mngSvc.BatchUpdateState(c, b))
}

// tagList .
func tagList(c *bm.Context) {
	t := &model.SearchTagParams{}
	if err := c.Bind(t); err != nil {
		return
	}
	data, total, err := mngSvc.TagList(c, t)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	page := map[string]int64{
		"num":   t.PN,
		"size":  t.PS,
		"total": total,
	}

	c.JSON(map[string]interface{}{
		"page": page,
		"data": data,
	}, err)
}

// typeList .
func typeList(c *bm.Context) {
	tl := &model.TagTypeList{}
	if err := c.Bind(tl); err != nil {
		return
	}
	if tl.BID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.TypeList(c, tl))
}

// attrList .
func attrList(c *bm.Context) {
	tba := &model.TagBusinessAttr{}
	if err := c.Bind(tba); err != nil {
		return
	}
	if tba.Bid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(mngSvc.AttrList(c, tba))
}

// attrUpdate .
func attrUpdate(c *bm.Context) {
	tba := &model.TagBusinessAttr{}
	if err := c.Bind(tba); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AttrUpdate(c, tba))
}

// addControl .
func addControl(c *bm.Context) {
	tc := &model.TagControl{}
	if err := c.Bind(tc); err != nil {
		return
	}
	c.JSON(nil, mngSvc.AddControl(c, tc))
}

// updateControl .
func updateControl(c *bm.Context) {
	tc := &model.TagControl{}
	if err := c.Bind(tc); err != nil {
		return
	}
	if tc.ID <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, mngSvc.UpdateControl(c, tc))
}

// tagControl .
func tagControl(c *bm.Context) {
	tc := &model.TagControlParam{}
	if err := c.Bind(tc); err != nil {
		return
	}
	c.JSON(mngSvc.TagControl(c, tc))
}
