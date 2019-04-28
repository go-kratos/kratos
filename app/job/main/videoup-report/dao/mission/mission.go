package mission

import (
	"context"
	"time"

	"go-common/app/job/main/videoup-report/model/mission"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_msAllURL = "/activity/list/videoall"
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
