package dao

import (
	"context"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
)

const (
	_dataRankURI = "/data/rank/recent_region-%d-%d.json"
)

// RankList get data rank by tid
func (d *Dao) RankList(c context.Context, tid int64, day int32) (resp *model.RankRecentResp, err error) {
	if err = d.httpCli.RESTfulGet(c, d.conf.Host.DataRank+_dataRankURI, "", nil, &resp, tid, day); err != nil {
		log.Error("RankList(tid:%v,day:%v),error(%v)", tid, day, err)
		return
	}
	return
}
