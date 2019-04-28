package service

import (
	"context"
	"math"
	"time"

	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/log"
)

//FlushHot for compute archive hot value by time.
func (s *Service) FlushHot(bs int) {
	var (
		c     = context.TODO()
		id    int64
		limit = 30
	)
	for {
		arcs, err := s.aca.Archives(c, id, bs, limit)
		if err != nil {
			log.Error("s.aca.Archives id(%d) error(%v)", id, err)
			return
		}
		count := len(arcs)
		if count == 0 {
			id = 0 //重新按定时跑数据
			time.Sleep(time.Hour * 1)
			continue
		}
		oids := make([]int64, 0)
		for _, a := range arcs {
			oids = append(oids, a.OID)
			id = a.ID //标记当前处理到的id
		}
		hots, err := s.computeHotByOIDs(c, oids, bs)
		if err != nil {
			log.Error("s.aca.Archives id(%d) error(%v)", id, err)
			return
		}
		if s.aca.UPHotByAIDs(c, hots) != nil {
			log.Error("s.aca.UPHotByAIDs hots(%+v) error(%v)", hots, err)
			return
		}
	}
}

func (s *Service) computeHotByOIDs(c context.Context, oids []int64, bs int) (res map[int64]int64, err error) {
	res = make(map[int64]int64)
	if bs == 1 {
		arcs, err := s.arc.Archives(c, oids)
		if err != nil {
			return nil, err
		}
		stat, err := s.arc.Stats(c, oids)
		if err != nil {
			log.Error("s.arc.Stats oids(%+v)|business(%d)|error(%v)", oids, bs, err)
			return nil, err
		}
		for _, oid := range oids {
			if v, ok := arcs[oid]; ok && v != nil {
				if t, ok := stat[oid]; ok && t != nil {
					res[oid] = countArcHot(t, int64(v.PubDate))
				}
			}
		}
	} else if bs == 2 {
		arts, err := s.arc.ArticleMetas(c, oids)
		if err != nil {
			return nil, err
		}
		for _, oid := range oids {
			if v, ok := arts[oid]; ok && v != nil {
				res[oid] = countArtHot(v)
			}
		}
	}
	return
}

//countArcHot 视频=硬币*0.4+收藏*0.3+弹幕*0.4+评论*0.4+播放*0.25+点赞*0.4+分享*0.6 最新视频（一天内发布）提权[总值*1.5]
func countArcHot(t *api.Stat, ptime int64) int64 {
	if t == nil {
		return 0
	}
	hot := float64(t.Coin)*0.4 +
		float64(t.Fav)*0.3 +
		float64(t.Danmaku)*0.4 +
		float64(t.Reply)*0.4 +
		float64(t.View)*0.25 +
		float64(t.Like)*0.4 +
		float64(t.Share)*0.6
	if ptime >= time.Now().AddDate(0, 0, -1).Unix() && ptime <= time.Now().Unix() {
		hot *= 1.5
	}
	return int64(math.Floor(hot))
}

// countArtHot 专栏=硬币*0.4+收藏*0.3+评论*0.4+阅读*0.25+点赞*0.4+分享*0.6 最新专栏（一天内发布）提权[总值*1.5]
func countArtHot(t *model.Meta) int64 {
	if t.Stats == nil {
		return 0
	}
	hot := float64(t.Stats.Coin)*0.4 +
		float64(t.Stats.Favorite)*0.3 +
		float64(t.Stats.Reply)*0.4 +
		float64(t.Stats.View)*0.25 +
		float64(t.Stats.Like)*0.4 +
		float64(t.Stats.Share)*0.6
	if int64(t.PublishTime) >= time.Now().AddDate(0, 0, -1).Unix() && int64(t.PublishTime) <= time.Now().Unix() {
		hot *= 1.5
	}
	return int64(math.Floor(hot))
}
