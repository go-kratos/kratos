package http

import (
	"encoding/json"

	"go-common/app/service/main/search/dao"
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// dmUpdate .
func dmUpdate(c *bm.Context) {
	params := c.Request.Form
	appid := params.Get("appid")
	if appid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	switch appid {
	case "dm_search":
		dmMediaUpdate(c)
	default:
		c.JSON(nil, ecode.RequestErr)
	}
}

func dmMediaUpdate(c *bm.Context) {
	var (
		err      error
		bulkItem []dao.BulkMapItem
	)
	params := c.Request.Form
	data := params.Get("data")
	if data == "" {
		return
	}
	var arr []map[string]interface{}
	if err = json.Unmarshal([]byte(data), &arr); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		return
	}
	key := []string{"attr", "attr_format", "ctime", "mid", "mode", "msg", "mtime",
		"pool", "progress", "state", "type", "oidstr"}
	for _, v := range arr {
		item := &model.DmUptParams{}
		var (
			ok  bool
			id  float64
			oid float64
		)
		if _, ok = v["id"]; !ok {
			continue
		}
		if id, ok = v["id"].(float64); !ok {
			continue
		}
		if _, ok = v["oid"]; !ok {
			continue
		}
		if oid, ok = v["oid"].(float64); !ok {
			continue
		}
		item.ID = int64(id)
		item.Oid = int64(oid)
		itemField := make(map[string]interface{})
		for _, k := range key {
			var it interface{}
			if it, ok = v[k]; ok && v[k] != nil {
				itemField[k] = it
			}
		}
		item.Field = itemField
		bulkItem = append(bulkItem, item)
	}
	if err = svr.UpdateMap(c, "dmExternal", bulkItem); err != nil {
		log.Error("srv.Update erro(%v)", err)
	}
	c.JSON(nil, nil)
}
