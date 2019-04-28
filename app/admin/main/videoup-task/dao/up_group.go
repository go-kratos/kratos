package dao

import (
	"context"
	"net/url"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

//UPGroups get all the up groups
func (d *Dao) UPGroups(c context.Context, mids []int64) (groups map[int64][]*model.UPGroup, err error) {
	val := url.Values{}
	mid := xstr.JoinInts(mids)
	val.Set("mids", mid)
	val.Set("group_id", "0")

	groups = map[int64][]*model.UPGroup{}
	for _, mid := range mids {
		groups[mid] = []*model.UPGroup{}
	}

	var res struct {
		Code int `json:"code"`
		Data struct {
			Items []map[string]interface{} `json:"items"`
		}
	}
	if err = d.hclient.Get(c, d.upGroupURL, "", val, &res); err != nil {
		log.Error("UPGroups url(%s) error(%v)", d.upGroupURL+"?"+val.Encode(), err)
		return
	}
	if res.Code != 0 || res.Data.Items == nil {
		log.Warn("UPGroups code(%d) !=0 or empty url(%s) error(%v)", res.Code, d.upGroupURL+"?"+val.Encode(), res.Code)
		return
	}
	for _, item := range res.Data.Items {
		g := &model.UPGroup{
			ID:        int64(item["group_id"].(float64)),
			Tag:       item["group_name"].(string),
			ShortTag:  item["group_tag"].(string),
			FontColor: item["font_color"].(string),
			BgColor:   item["bg_color"].(string),
			Note:      item["note"].(string),
		}
		mid := int64(item["mid"].(float64))
		groups[mid] = append(groups[mid], g)
	}
	return
}
