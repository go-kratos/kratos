package http

import (
	"encoding/json"

	"go-common/app/service/main/search/dao"
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// pgcUpdate .
func pgcUpdate(c *bm.Context) {
	params := c.Request.Form
	appid := params.Get("appid")
	if appid == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	switch appid {
	case "pgc_media":
		go pgcMediaUpdate(c)
	default:
		c.JSON(nil, ecode.RequestErr)
	}
}

func pgcMediaUpdate(c *bm.Context) {
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
	key := []string{"season_id", "copyright", "latest_time", "dm_count", "play_count", "fav_count", "area_id", "score",
		"is_finish", "season_version", "season_status", "release_date", "pub_time", "season_month", "copyright_info"}
	for _, v := range arr {
		item := &model.PgcMediaUptParams{}
		var (
			ok bool
			id float64
		)
		if _, ok = v["media_id"]; !ok {
			continue
		}
		if id, ok = v["media_id"].(float64); !ok {
			continue
		}
		item.MediaID = int64(id)
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
	if err = svr.UpdateMap(c, "externalPublic", bulkItem); err != nil {
		log.Error("srv.Update erro(%v)", err)
	}
}
