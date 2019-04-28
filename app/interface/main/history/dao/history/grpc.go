package history

import (
	"context"

	"go-common/app/interface/main/history/model"
	epmdl "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"
	"go-common/library/log"
)

// Bangumis .
func (d *Dao) Bangumis(c context.Context, mid int64, epid []int64) (map[int64]*model.Bangumi, error) {
	req := &epmdl.EpReq{}
	for _, aid := range epid {
		req.Epids = append(req.Epids, int32(aid))
	}
	resp, err := d.epidGRPC.Cards(c, req)
	if err != nil {
		log.Error("s.epidGRPC.Cards(%v, %v) error(%v)", c, req, err)
		return nil, err
	}
	data := make(map[int64]*model.Bangumi)
	for k, v := range resp.Cards {
		data[int64(k)] = &model.Bangumi{
			Epid:          int64(v.EpisodeId),
			Title:         v.Title,
			LongTitle:     v.LongTitle,
			EpisodeStatus: int(v.EpisodeStatus),
			Cover:         v.Cover,
			Season: &model.Season{
				ID:            int64(v.Season.SeasonId),
				Title:         v.Season.Title,
				SeasonStatus:  int(v.Season.SeasonStatus),
				IsFinish:      int(v.Season.IsFinish),
				TotalCount:    v.Season.TotalCount,
				NewestEpid:    int64(v.Season.NewEpId),
				NewestEpIndex: v.Season.NewEpIndex,
				SeasonType:    int(v.Season.SeasonType),
			},
		}
	}
	return data, nil
}

// BangumisByAids .
func (d *Dao) BangumisByAids(c context.Context, mid int64, aids []int64, realIP string) (map[int64]*model.BangumiSeason, error) {
	req := &epmdl.EpAidReq{}
	for _, aid := range aids {
		req.Aids = append(req.Aids, int32(aid))
	}
	resp, err := d.epidGRPC.ListByAids(c, req)
	if err != nil {
		log.Error("s.epidGRPC.ListByAids(%v, %v) error(%v)", c, req, err)
		return nil, err
	}
	data := make(map[int64]*model.BangumiSeason)
	for _, v := range resp.Infos {
		data[int64(v.Aid)] = &model.BangumiSeason{
			ID:   int64(v.SeasonId),
			Epid: int64(v.EpisodeId),
		}
	}
	return data, nil
}
