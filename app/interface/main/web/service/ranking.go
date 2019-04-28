package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/web/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/log"
)

const (
	_rankIndexLen   = 8
	_rankLen        = 100
	_rankRegionLen  = 10
	_rankOtherLimit = 10
	_avType         = "av"
)

var (
	_emptyRankArchive = make([]*model.RankArchive, 0)
)

// Ranking get ranking data.
func (s *Service) Ranking(c context.Context, rid int16, rankType, day, arcType int) (res *model.RankData, err error) {
	var (
		rankArc  *model.RankNew
		addCache = true
	)
	if res, err = s.dao.RankingCache(c, rid, rankType, day, arcType); err != nil {
		err = nil
		addCache = false
	} else if res != nil && len(res.List) > 0 {
		return
	}
	if rankArc, err = s.dao.Ranking(c, rid, rankType, day, arcType); err != nil {
		err = nil
	} else if rankArc != nil && len(rankArc.List) > s.c.Rule.MinRankCount {
		res = &model.RankData{Note: rankArc.Note}
		if res.List, err = s.fmtRankArcs(c, rankArc.List, _rankLen); err != nil {
			err = nil
		} else if len(res.List) > 0 {
			if addCache {
				s.cache.Do(c, func(c context.Context) {
					s.dao.SetRankingCache(c, rid, rankType, day, arcType, res)
				})
			}
			return
		}
	} else {
		log.Error("s.dao.RankingNew(%d,%d,%d) len(aids) (%d)", rid, day, arcType, len(rankArc.List))
	}
	res, err = s.dao.RankingBakCache(c, rid, rankType, day, arcType)
	if res == nil || len(res.List) == 0 {
		res = &model.RankData{List: _emptyRankArchive}
	}
	return
}

// RankingIndex get index ranking data
func (s *Service) RankingIndex(c context.Context, day int) (res []*model.IndexArchive, err error) {
	var (
		addCache = true
		rs       []*model.NewArchive
		arcs     map[int64]*arcmdl.Arc
	)
	if res, err = s.dao.RankingIndexCache(c, day); err != nil {
		err = nil
		addCache = false
	} else if len(res) > 0 {
		return
	}
	if rs, err = s.dao.RankingIndex(c, day); err != nil {
		err = nil
	} else if len(rs) > s.c.Rule.MinRankIndexCount {
		if arcs, err = s.fillArcs(c, rs); err != nil {
			err = nil
		} else if len(arcs) > 0 {
			res = fmtIndexArcs(rs, arcs)
			if addCache {
				s.cache.Do(c, func(c context.Context) {
					s.dao.SetRankingIndexCache(c, day, res)
				})
			}
			return
		}
	} else {
		log.Error("s.dao.RankingIndexCache(%d) len(aids) (%d)", day, len(res))
	}
	if res, err = s.dao.RankingIndexBakCache(c, day); err != nil {
		return
	}
	if len(res) == 0 {
		res = []*model.IndexArchive{}
	}
	return
}

// RankingRegion get region ranking data
func (s *Service) RankingRegion(c context.Context, rid int16, day, original int) (res []*model.RegionArchive, err error) {
	var (
		addCache = true
		rs       []*model.NewArchive
		arcs     map[int64]*arcmdl.Arc
	)
	defer func() {
		if len(res) > 0 {
			s.fmtRegionStats(c, res)
		}
	}()
	if res, err = s.dao.RankingRegionCache(c, rid, day, original); err != nil {
		err = nil
		addCache = false
	} else if len(res) > 0 {
		return
	}
	if rs, err = s.dao.RankingRegion(c, rid, day, original); err != nil {
		err = nil
	} else if len(rs) > s.c.Rule.MinRankRegionCount {
		if arcs, err = s.fillArcs(c, rs); err != nil {
			err = nil
		} else if len(arcs) > 0 {
			res = fmtRegionArcs(rs, arcs)
			if addCache {
				s.cache.Do(c, func(c context.Context) {
					s.dao.SetRankingRegionCache(c, rid, day, original, res)
				})
			}
			return
		}
	} else {
		log.Error("s.dao.RankingRegion(%d,%d,%d) len(aids) (%d)", rid, day, original, len(rs))
	}
	res, err = s.dao.RankingRegionBakCache(c, rid, day, original)
	if len(res) == 0 {
		res = []*model.RegionArchive{}
	}
	return
}

// RankingRecommend get rank recommend data.
func (s *Service) RankingRecommend(c context.Context, rid int16) (res []*model.IndexArchive, err error) {
	var (
		addCache = true
		rs       []*model.NewArchive
		arcs     map[int64]*arcmdl.Arc
	)
	if res, err = s.dao.RankingRecommendCache(c, rid); err != nil {
		err = nil
		addCache = false
	} else if len(res) > 0 {
		return
	}
	if rs, err = s.dao.RankingRecommend(c, rid); err != nil {
		err = nil
	} else if len(rs) > s.c.Rule.MinRankRecCount {
		if arcs, err = s.fillArcs(c, rs); err != nil {
			err = nil
		} else if len(arcs) > 0 {
			res = fmtIndexArcs(rs, arcs)
			if addCache {
				s.cache.Do(c, func(c context.Context) {
					s.dao.SetRankingRecommendCache(c, rid, res)
				})
			}
			return
		}
	} else {
		log.Error("s.dao.RankingRecommend(%d) len(aids) (%d)", rid, len(res))
	}
	res, err = s.dao.RankingRecommendBakCache(c, rid)
	if len(res) == 0 {
		res = []*model.IndexArchive{}
	}
	return
}

// RankingTag get tag ranking data
func (s *Service) RankingTag(c context.Context, rid int16, tagID int64) (res []*model.TagArchive, err error) {
	var (
		addCache  = true
		tagArcs   []*model.NewArchive
		arcsReply *arcmdl.ArcsReply
	)
	defer func() {
		s.fmtTagStats(c, res)
	}()
	if res, err = s.dao.RankingTagCache(c, rid, tagID); err != nil {
		err = nil
		addCache = false
	} else if len(res) > 0 {
		return
	}
	if tagArcs, err = s.dao.RankingTag(c, rid, tagID); err != nil {
		err = nil
	} else if len(tagArcs) > s.c.Rule.MinRankTagCount {
		var aids []int64
		for _, v := range tagArcs {
			aids = append(aids, v.Aid)
		}
		if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
			err = nil
		} else if len(arcsReply.Arcs) > 0 {
			res = fmtTagArchives(tagArcs, arcsReply.Arcs)
			if addCache {
				s.cache.Do(c, func(c context.Context) {
					s.dao.SetRankingTagCache(c, rid, tagID, res)
				})
			}
		}
		return
	} else {
		log.Error("s.dao.RankingRecommend(%d) len(aids) (%d)", rid, len(res))
	}
	res, err = s.dao.RankingTagBakCache(c, rid, tagID)
	if len(res) == 0 {
		res = []*model.TagArchive{}
	}
	return
}

// RegionCustom region custom data
func (s *Service) RegionCustom(c context.Context) (res []*model.Custom, err error) {
	var (
		addCache  = true
		aids      []int64
		arcsReply *arcmdl.ArcsReply
	)
	if res, err = s.dao.RegionCustomCache(c); err != nil {
		err = nil
		addCache = false
	} else if len(res) > 0 {
		return
	}
	if res, err = s.dao.RegionCustom(c); err != nil {
		log.Error("s.dao.RegionCustom error(%v)", err)
		res = []*model.Custom{}
		return
	} else if len(res) > 0 {
		for _, item := range res {
			if item.Type == _avType && item.Aid > 0 {
				aids = append(aids, item.Aid)
			}
		}
		if len(aids) > 0 {
			archivesArgLog("RegionCustom", aids)
			if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
				log.Error("s.arcClient.Arcs(%v) error(%v)", aids, err)
				err = nil
			} else {
				for _, v := range res {
					if arc, ok := arcsReply.Arcs[v.Aid]; ok && v.Type == _avType {
						v.Pic = arc.Pic
						v.Title = arc.Title
					}
				}
			}
		}
		if addCache {
			s.cache.Do(c, func(c context.Context) {
				s.dao.SetRegionCustomCache(c, res)
			})
		}
		return
	}
	res, err = s.dao.RegionCustomBakCache(c)
	if len(res) == 0 {
		res = []*model.Custom{}
	}
	return
}

func (s *Service) fillArcs(c context.Context, rankArchives []*model.NewArchive) (res map[int64]*arcmdl.Arc, err error) {
	var (
		aids      []int64
		arcsReply *arcmdl.ArcsReply
	)
	for _, arc := range rankArchives {
		if arc == nil {
			continue
		}
		aids = append(aids, arc.Aid)
	}
	archivesArgLog("fillArcs", aids)
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("s.arcClient.Arcs(%v) error(%v)", aids, err)
		return
	}
	res = arcsReply.Arcs
	return
}

func fmtRegionArcs(aids []*model.NewArchive, arcs map[int64]*arcmdl.Arc) (res []*model.RegionArchive) {
	for _, arc := range aids {
		if arc == nil {
			continue
		}
		if len(res) > _rankRegionLen {
			break
		}
		if _, ok := arcs[arc.Aid]; !ok {
			continue
		}
		regionArchive := &model.RegionArchive{
			Aid:         strconv.FormatInt(arcs[arc.Aid].Aid, 10),
			Typename:    arcs[arc.Aid].TypeName,
			Title:       arcs[arc.Aid].Title,
			Play:        fmtArcView(arcs[arc.Aid]),
			Review:      arcs[arc.Aid].Stat.Reply,
			VideoReview: arcs[arc.Aid].Stat.Danmaku,
			Favorites:   arcs[arc.Aid].Stat.Fav,
			Mid:         arcs[arc.Aid].Author.Mid,
			Author:      arcs[arc.Aid].Author.Name,
			Description: arcs[arc.Aid].Desc,
			Create:      time.Unix(int64(arcs[arc.Aid].PubDate), 0).Format("2006-01-02 15:04"),
			Pic:         arcs[arc.Aid].Pic,
			Coins:       arcs[arc.Aid].Stat.Coin,
			Duration:    fmtDuration(arcs[arc.Aid].Duration),
			Pts:         arc.Score,
			Rights:      arcs[arc.Aid].Rights,
		}
		res = append(res, regionArchive)
	}
	return
}

func fmtIndexArcs(aids []*model.NewArchive, arcs map[int64]*arcmdl.Arc) (res []*model.IndexArchive) {
	var (
		typeName string
		ok       bool
	)
	for _, arc := range aids {
		if arc == nil {
			continue
		}
		if len(res) > _rankIndexLen {
			break
		}
		if _, ok = arcs[arc.Aid]; !ok {
			continue
		}
		if typeName, ok = model.RecSpecTypeName[arcs[arc.Aid].TypeID]; !ok {
			typeName = arcs[arc.Aid].TypeName
		}
		indexArchive := &model.IndexArchive{
			Aid:         strconv.FormatInt(arcs[arc.Aid].Aid, 10),
			Typename:    typeName,
			Title:       arcs[arc.Aid].Title,
			Play:        fmtArcView(arcs[arc.Aid]),
			Review:      arcs[arc.Aid].Stat.Reply,
			VideoReview: arcs[arc.Aid].Stat.Danmaku,
			Favorites:   arcs[arc.Aid].Stat.Fav,
			Mid:         arcs[arc.Aid].Author.Mid,
			Author:      arcs[arc.Aid].Author.Name,
			Description: arcs[arc.Aid].Desc,
			Create:      time.Unix(int64(arcs[arc.Aid].PubDate), 0).Format("2006-01-02 15:04"),
			Pic:         arcs[arc.Aid].Pic,
			Coins:       arcs[arc.Aid].Stat.Coin,
			Duration:    fmtDuration(arcs[arc.Aid].Duration),
			Rights:      arcs[arc.Aid].Rights,
		}
		res = append(res, indexArchive)
	}
	return
}

func (s *Service) fmtRankArcs(c context.Context, rankArchives []*model.RankNewArchive, arcLen int) (res []*model.RankArchive, err error) {
	var (
		aids      []int64
		arcsReply *arcmdl.ArcsReply
	)
	for _, arc := range rankArchives {
		if arc == nil {
			continue
		}
		aids = append(aids, arc.Aid)
		if len(arc.Others) > 0 {
			i := 0
			for _, a := range arc.Others {
				if a == nil {
					continue
				}
				aids = append(aids, a.Aid)
				i++
				if i >= _rankOtherLimit {
					break
				}
			}
		}
	}
	archivesArgLog("fmtRankArcs", aids)
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("s.arcClient.Arcs(%v) error(%v)", aids, err)
		return
	}
	for _, arc := range rankArchives {
		if arc == nil {
			continue
		}
		if len(res) > arcLen {
			break
		}
		if _, ok := arcsReply.Arcs[arc.Aid]; !ok {
			continue
		}
		var coin, danmu int32
		if arc.RankStat == nil {
			coin = arcsReply.Arcs[arc.Aid].Stat.Coin
			danmu = arcsReply.Arcs[arc.Aid].Stat.Danmaku
		} else {
			coin = arc.RankStat.Coin
			danmu = arc.RankStat.Danmu
			arcsReply.Arcs[arc.Aid].Stat.View = arc.RankStat.Play
		}
		rankArchive := &model.RankArchive{
			Aid:         strconv.FormatInt(arcsReply.Arcs[arc.Aid].Aid, 10),
			Author:      arcsReply.Arcs[arc.Aid].Author.Name,
			Coins:       coin,
			Duration:    fmtDuration(arcsReply.Arcs[arc.Aid].Duration),
			Mid:         arcsReply.Arcs[arc.Aid].Author.Mid,
			Pic:         arcsReply.Arcs[arc.Aid].Pic,
			Play:        fmtArcView(arcsReply.Arcs[arc.Aid]),
			Pts:         arc.Score,
			Title:       arcsReply.Arcs[arc.Aid].Title,
			VideoReview: danmu,
			Rights:      arcsReply.Arcs[arc.Aid].Rights,
		}
		if len(arc.Others) > 0 {
			for _, a := range arc.Others {
				if a == nil {
					continue
				}
				if _, ok := arcsReply.Arcs[a.Aid]; !ok {
					continue
				}
				archive := &model.Other{
					Aid:         a.Aid,
					Play:        fmtArcView(arcsReply.Arcs[a.Aid]),
					VideoReview: arcsReply.Arcs[a.Aid].Stat.Danmaku,
					Coins:       arcsReply.Arcs[a.Aid].Stat.Coin,
					Pts:         a.Score,
					Title:       arcsReply.Arcs[a.Aid].Title,
					Pic:         arcsReply.Arcs[a.Aid].Pic,
					Duration:    fmtDuration(arcsReply.Arcs[a.Aid].Duration),
					Rights:      arcsReply.Arcs[a.Aid].Rights,
				}
				rankArchive.Others = append(rankArchive.Others, archive)
			}
		}
		res = append(res, rankArchive)
	}
	return
}

// fmtRegionStats get real time region archive stat
func (s *Service) fmtRegionStats(c context.Context, res []*model.RegionArchive) {
	var (
		aids       []int64
		aid        int64
		err        error
		statsReply *arcmdl.StatsReply
	)
	for _, arc := range res {
		if arc == nil {
			continue
		}
		if aid, err = strconv.ParseInt(arc.Aid, 10, 64); err != nil {
			continue
		}
		aids = append(aids, aid)
	}
	if len(aids) == 0 {
		return
	}
	if statsReply, err = s.arcClient.Stats(c, &arcmdl.StatsRequest{Aids: aids}); err != nil {
		log.Error("s.arcClient.Stats(%v) error(%v)", aids, err)
		return
	}
	arcStats := statsReply.Stats
	for _, arc := range res {
		if aid, err = strconv.ParseInt(arc.Aid, 10, 64); err != nil {
			continue
		}
		if arcStat, ok := arcStats[aid]; ok {
			arc.Play = arcStat.View
			arc.VideoReview = arcStat.Danmaku
			arc.Favorites = arcStat.Fav
			arc.Coins = arcStat.Coin
		}
	}
}

func fmtTagArchives(tagArcs []*model.NewArchive, arcs map[int64]*arcmdl.Arc) (res []*model.TagArchive) {
	for _, tagArc := range tagArcs {
		if arc, ok := arcs[tagArc.Aid]; ok {
			res = append(res, &model.TagArchive{
				Title:       arc.Title,
				Author:      arc.Author.Name,
				Description: arc.Desc,
				Pic:         arc.Pic,
				Play:        strconv.FormatInt(int64(arc.Stat.View), 10),
				Favorites:   strconv.FormatInt(int64(arc.Stat.Fav), 10),
				Mid:         strconv.FormatInt(arc.Author.Mid, 10),
				Review:      strconv.FormatInt(int64(arc.Stat.Reply), 10),
				CreatedAt:   time.Unix(int64(arcs[arc.Aid].PubDate), 0).Format("2006-01-02 15:04"),
				VideoReview: strconv.FormatInt(int64(arc.Stat.Danmaku), 10),
				Coins:       strconv.FormatInt(int64(arc.Stat.Coin), 10),
				Duration:    strconv.FormatInt(arc.Duration, 10),
				Aid:         arc.Aid,
				Pts:         tagArc.Score,
				Rights:      arc.Rights,
			})
		}
	}
	return
}

// fmtTagStats get real time tag archive stat
func (s *Service) fmtTagStats(c context.Context, res []*model.TagArchive) {
	var (
		aids       []int64
		err        error
		statsReply *arcmdl.StatsReply
	)
	for _, arc := range res {
		if arc == nil {
			continue
		}
		aids = append(aids, arc.Aid)
	}
	if len(aids) == 0 {
		return
	}
	if statsReply, err = s.arcClient.Stats(c, &arcmdl.StatsRequest{Aids: aids}); err != nil {
		log.Error("s.arcClient.Stats(%v) error(%v)", aids, err)
		return
	}
	arcStats := statsReply.Stats
	for _, arc := range res {
		if arcStat, ok := arcStats[arc.Aid]; ok {
			arc.Play = strconv.FormatInt(int64(arcStat.View), 10)
			arc.VideoReview = strconv.FormatInt(int64(arcStat.Danmaku), 10)
			arc.Favorites = strconv.FormatInt(int64(arcStat.Fav), 10)
			arc.Coins = strconv.FormatInt(int64(arcStat.Coin), 10)
		}
	}
}

func fmtDuration(duration int64) (du string) {
	if duration == 0 {
		du = "00:00"
	} else {
		var min, sec string
		min = strconv.Itoa(int(duration / 60))
		if int(duration%60) < 10 {
			sec = "0" + strconv.Itoa(int(duration%60))
		} else {
			sec = strconv.Itoa(int(duration % 60))
		}
		du = min + ":" + sec
	}
	return
}

func fmtArcView(a *arcmdl.Arc) interface{} {
	var view interface{} = a.Stat.View
	if a.Access > 0 {
		view = "--"
	}
	return view
}
