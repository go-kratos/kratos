package http

import (
	"strconv"
	"strings"

	"go-common/app/interface/main/space/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func settingInfo(c *bm.Context) {
	var (
		mid int64
		err error
	)
	midStr := c.Request.Form.Get("mid")
	if mid, err = strconv.ParseInt(midStr, 10, 64); err != nil || mid <= 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(spcSvc.SettingInfo(c, mid))
}

func privacyModify(c *bm.Context) {
	var (
		mid        int64
		field      string
		value      int
		fieldCheck bool
		err        error
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	if field = params.Get("field"); field == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, v := range model.PrivacyFields {
		if v == field {
			fieldCheck = true
			break
		}
	}
	if !fieldCheck {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	valueStr := params.Get("value")
	if value, err = strconv.Atoi(valueStr); err != nil || (value != 0 && value != 1) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.PrivacyModify(c, mid, field, value))
}

func privacyBatchModify(c *bm.Context) {
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid := midStr.(int64)
	batch := make(map[string]int)
	for _, v := range model.PrivacyFields {
		if valueStr := params.Get(v); valueStr != "" {
			if value, err := strconv.Atoi(params.Get(v)); err == nil && (value == 0 || value == 1) {
				batch[v] = value
			}
		}
	}
	c.JSON(nil, spcSvc.PrivacyBatchModify(c, mid, batch))
}

func indexOrderModify(c *bm.Context) {
	var (
		mid           int64
		indexOrderStr string
		indexOrder    []string
	)
	params := c.Request.Form
	midStr, _ := c.Get("mid")
	mid = midStr.(int64)
	if indexOrderStr = params.Get("index_order"); indexOrderStr == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	indexOrders := strings.Split(indexOrderStr, ",")
	for _, v := range indexOrders {
		i, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if _, ok := model.IndexOrderMap[i]; !ok {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		indexOrder = append(indexOrder, v)
	}
	if len(indexOrder) != len(model.DefaultIndexOrder) {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, spcSvc.IndexOrderModify(c, mid, indexOrder))
}
