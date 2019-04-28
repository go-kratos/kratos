package dao

import (
	"context"
	"net/url"
	"strconv"

	"go-common/app/interface/main/web/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// OnlineCount get online count
func (d *Dao) OnlineCount(c context.Context) (data *model.OnlineCount, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	var res struct {
		Code int                `json:"code"`
		Data *model.OnlineCount `json:"data"`
	}
	if err = d.httpR.Get(c, d.onlineURL, ip, params, &res); err != nil {
		log.Error(" d.httpW.Get.Get(%s) error(%v)", d.onlineURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpW.Get(%s) code error(%d)", d.onlineURL, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}

// LiveOnlineCount get live online count
func (d *Dao) LiveOnlineCount(c context.Context) (data *model.LiveOnlineCount, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	var res struct {
		Code int                    `json:"code"`
		Data *model.LiveOnlineCount `json:"data"`
	}
	if err = d.httpR.Get(c, d.liveOnlineURL, ip, params, &res); err != nil {
		log.Error(" d.httpW.Get.Get(%s) error(%v)", d.liveOnlineURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpW.Get(%s) code error(%d)", d.liveOnlineURL, res.Code)
		err = ecode.Int(res.Code)
		return
	}
	data = res.Data
	return
}

// OnlineList get online list
func (d *Dao) OnlineList(c context.Context, num int64) (data []*model.OnlineAid, err error) {
	var (
		params = url.Values{}
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("num", strconv.FormatInt(num, 10))
	var res struct {
		Code int                `json:"code"`
		Data []*model.OnlineAid `json:"data"`
	}
	if err = d.httpR.Get(c, d.onlineListURL, ip, params, &res); err != nil {
		log.Error(" d.httpR.Get.Get(%s) error(%v)", d.onlineListURL, err)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("d.httpR.Get(%s) code error(%d)", d.onlineListURL, res.Code)
		return
	}
	data = res.Data
	return
}
