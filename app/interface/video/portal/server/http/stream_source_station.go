package http

import (
	"go-common/app/interface/video/portal/conf"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

func init() {
	conf.Init()
}

// StreamSourceStation 提供给第三方的源站地址
func streamLevelTwoSourceStation(c *bm.Context) {
	//从配置中心读取
	//conf.Config.LevelTwoSourceStation
	thisconf := *conf.Conf
	c.JSONMap(map[string]interface{}{"message": "ok", "data": thisconf.LevelTwoSourceStation}, nil)
}

//LPL全明星赛

func streamLplAllStar(c *bm.Context) {
	params := c.Request.URL.Query()
	rid := params.Get("room_id")

	if rid == "" {
		c.Set("output_data", "some fields are empty")
		c.JSONMap(map[string]interface{}{"message": "some fields are empty"}, ecode.RequestErr)
		c.Abort()
		return
	}

	urls := []map[string]interface{}{}
	url := map[string]interface{}{}
	if rid == "1111" || rid == "11090072" {
		url["ext"] = "flv"
		url["rate_level"] = 4
		url["url"] = "http://nbvc.live-play.acgvideo.com/live-bvc/946862/live_325164925_5324520_800.flv?wsSecret=a65b9dd9a5a04e298ebce381673d8a77&wsTime=1546571918&trid=f9516d154cf54f47bd9b329bb34c8b25&sig=no"
		// allstarresult.default_rate_level = 4
		urls = append(urls, url)
		//allstarresult.urls = urls
		c.JSONMap(map[string]interface{}{"message": "ok", "data": map[string]interface{}{"default_rate_level": 4, "urls": urls}}, nil)
	} else {
		c.JSONMap(map[string]interface{}{"message": "ok", "data": map[string]interface{}{"default_rate_level": 4, "urls": urls}}, nil)
	}

}
