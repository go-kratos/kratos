package service

import (
	"context"
	"time"

	"go-common/app/interface/main/space/model"
	arcmdl "go-common/app/service/main/archive/api"
	coinmdl "go-common/app/service/main/coin/api"
	"go-common/library/log"
)

const (
	_coinVideoLimit = 100
	_businessCoin   = "archive"
)

var _emptyCoinArcList = make([]*model.CoinArc, 0)

// CoinVideo get coin archives
func (s *Service) CoinVideo(c context.Context, mid, vmid int64) (list []*model.CoinArc, err error) {
	var (
		coinReply *coinmdl.ListReply
		aids      []int64
		arcReply  *arcmdl.ArcsReply
	)
	if mid != vmid {
		if err = s.privacyCheck(c, vmid, model.PcyCoinVideo); err != nil {
			return
		}
	}
	if coinReply, err = s.coinClient.List(c, &coinmdl.ListReq{Mid: vmid, Business: _businessCoin, Ts: time.Now().Unix()}); err != nil {
		log.Error("s.coinClinet.List(%d) error(%v)", vmid, err)
		err = nil
		list = _emptyCoinArcList
		return
	}
	existAids := make(map[int64]int64, len(coinReply.List))
	afVideos := make(map[int64]*coinmdl.ModelList, len(coinReply.List))
	for _, v := range coinReply.List {
		if len(aids) > _coinVideoLimit {
			break
		}
		if _, ok := existAids[v.Aid]; ok {
			if v.Aid > 0 {
				afVideos[v.Aid].Number += v.Number
			}
			continue
		}
		if v.Aid > 0 {
			afVideos[v.Aid] = v
			aids = append(aids, v.Aid)
			existAids[v.Aid] = v.Aid
		}
	}
	if len(aids) == 0 {
		list = _emptyCoinArcList
		return
	}
	if arcReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("s.arcClient.Arcs(%v) error(%v)", aids, err)
		return
	}
	for _, aid := range aids {
		if arc, ok := arcReply.Arcs[aid]; ok && arc.IsNormal() {
			if arc.Access >= 10000 {
				arc.Stat.View = -1
			}
			if item, ok := afVideos[aid]; ok {
				list = append(list, &model.CoinArc{Arc: arc, Coins: item.Number, Time: item.Ts})
			}
		}
	}
	if len(list) == 0 {
		list = _emptyCoinArcList
	}
	return
}
