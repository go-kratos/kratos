package mission

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/videoup/model/mission"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_msAllURL           = "/activity/list/videoall"
	_actOnlineByTypeURI = "/activity/online/by/type"
)

// Missions get missions.
func (d *Dao) Missions(c context.Context) (mm map[int]*mission.Mission, err error) {
	var res struct {
		Code int `json:"code"`
		Data []*struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			ETime string `json:"etime"`
			Tags  string `json:"tags"`
		} `json:"data"`
	}
	if err = d.httpR.Get(c, d.missAllURL, "", nil, &res); err != nil {
		log.Error("videoup mission list error(%v) | missAllURL(%s)", err, d.missAllURL)
		return
	}
	if res.Code != 0 {
		log.Error("videoup mission list res.Code nq zero error(%v) | missAllURL(%s) res(%v)", res.Code, d.missAllURL, res)
		err = ecode.CreativeActivityErr
		return
	}
	mm = make(map[int]*mission.Mission, len(res.Data))
	for _, m := range res.Data {
		miss := &mission.Mission{}
		miss.ID = m.ID
		miss.Name = m.Name
		miss.ETime, _ = time.Parse("2006-01-02 15:04:05", m.ETime)
		miss.Tags = m.Tags
		mm[miss.ID] = miss
	}
	return
}

// MissionOnlineByTid fn, 这里默认会返回所有无投稿分区限制要求的通用活动，在做校验的时候允许此类活动投稿到任意分区
func (d *Dao) MissionOnlineByTid(c context.Context, tid int16) (mm map[int]*mission.Mission, err error) {
	var res struct {
		Code int `json:"code"`
		Data []*struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			ETime string `json:"etime"`
			Tags  string `json:"tags"`
		} `json:"data"`
	}
	mm = make(map[int]*mission.Mission)
	params := url.Values{}
	params.Set("type", strconv.Itoa(int(tid)))
	params.Set("plat", "1")
	if err = d.httpR.Get(c, d.actOnlineByTypeURL, "", params, &res); err != nil {
		log.Error("videoup actOnlineByTypeURL error(%v) | actOnlineByTypeURL(%s)", err, d.actOnlineByTypeURL+"?"+params.Encode())
		err = ecode.CreativeActivityErr
		return
	}
	if res.Code != 0 {
		log.Error("videoup actOnlineByTypeURL res.Code nq zero error(%v) | actOnlineByTypeURL(%s) res(%v)", res.Code, d.actOnlineByTypeURL+"?"+params.Encode(), res)
		err = ecode.CreativeActivityErr
		return
	}
	for _, m := range res.Data {
		miss := &mission.Mission{}
		miss.ID = m.ID
		miss.Name = m.Name
		miss.ETime, _ = time.Parse("2006-01-02 15:04:05", m.ETime)
		miss.Tags = m.Tags
		mm[miss.ID] = miss
	}
	return
}
