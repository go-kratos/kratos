package timemachine

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	model "go-common/app/interface/main/activity/model/timemachine"
	tagmdl "go-common/app/interface/main/tag/api"
	artmdl "go-common/app/interface/openplatform/article/model"
	accmdl "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"go-common/library/sync/errgroup.v2"
)

// StartTmproc start tm proc.
//func (s *Service) StartTmproc(c context.Context) (err error) {
//	s.dao.StartTmProc(context.Background())
//	return
//}

// StopTmproc start tm proc.
//func (s *Service) StopTmproc(c context.Context) (err error) {
//	s.dao.StopTmproc(c)
//	return
//}

// Timemachine2018Raw .
func (s *Service) Timemachine2018Raw(c context.Context, loginMid, mid int64) (data *model.Item, err error) {
	if _, ok := s.tmMidMap[loginMid]; !ok {
		err = ecode.AccessDenied
		return
	}
	if mid == 0 {
		mid = loginMid
	}
	if data, err = s.dao.RawTimemachine(c, mid); err != nil {
		log.Error("Timemachine2018 s.dao.RawTimemachine(%d) error(%v)", mid, err)
	}
	return
}

// Timemachine2018Cache .
func (s *Service) Timemachine2018Cache(c context.Context, loginMid, mid int64) (data *model.Item, err error) {
	if _, ok := s.tmMidMap[loginMid]; !ok {
		err = ecode.AccessDenied
		return
	}
	if mid == 0 {
		mid = loginMid
	}
	if data, err = s.dao.CacheTimemachine(c, mid); err != nil {
		log.Error("Timemachine2018Cache s.dao.CacheTimemachine(%d) error(%v)", mid, err)
	}
	return
}

// Timemachine2018 .
func (s *Service) Timemachine2018(c context.Context, loginMid, mid int64) (data *model.Timemachine, err error) {
	if _, ok := s.tmMidMap[loginMid]; !ok {
		mid = loginMid
	} else {
		if mid == 0 {
			mid = loginMid
		}
	}
	var item *model.Item
	if item, err = s.dao.CacheTimemachine(c, mid); err != nil {
		log.Error("Timemachine2018 s.dao.Timemachine(%d) error(%v)", mid, err)
		err = nil
	}
	if item == nil || item.DurationHour == 0 || (item.LikeUpAvDuration == 0 && item.LikeUpLiveDuration == 0) {
		item = &model.Item{Mid: mid}
	}
	data = s.groupTmData(c, item)
	return
}

func (s *Service) groupTmData(c context.Context, item *model.Item) (data *model.Timemachine) {
	var (
		aids, viewAids, upAids, artIDs, mids                          []int64
		totalView, ugcView, pgcView                                   []*model.AidView
		upBestLiveFanMid, liveMinute, upBestFanMid, favVv, subTidBest int64
	)
	data = &model.Timemachine{
		Mid:                item.Mid,
		IsUp:               item.IsUp,
		DurationHour:       item.DurationHour,
		ArchiveVv:          item.ArchiveVv,
		BrainwashCirVv:     item.BrainwashCirVv,
		FirstSubmitType:    item.FirstSubmitType,
		LikeSubtidVv:       item.LikeSubtidVv,
		LikeUpAvDuration:   item.LikeUpAvDuration,
		LikeUpLiveDuration: item.LikeUpLiveDuration,
		LikeUpDuration:     item.LikeUpAvDuration + item.LikeUpLiveDuration,
		LikeLiveUpSubTname: item.LikeLiveUpSubTname,
		BestAvidType:       item.BestAvidType,
		AllVv:              item.AllVv,
		BestAvidOldType:    item.BestAvidOldType,
		OldAvVv:            item.OldAvVv,
		UpLiveDuration:     item.UpLiveDuration,
		IsLiveUp:           item.IsLiveUp,
		ValidLiveDays:      item.ValidLiveDays,
		AddAttentions:      item.Attentions,
		WinRatio:           item.WinRatio,
		SubmitAvsRds:       item.SubmitAvsRds,
	}
	// win ratio fix
	if data.WinRatio == "100%" {
		data.WinRatio = "99%"
	}
	if data.ValidLiveDays > 365 {
		data.ValidLiveDays = 365
	}
	if item.Like2SubTids != "" {
		subTidBest = s.groupTagDesc(item, data)
	}
	if tagDesc, ok := s.tagDescMap[item.LikeTagID]; ok && tagDesc != nil {
		data.LikeTagDescFirst = tagDesc.Desc1
		data.LikeTagDescSecond = tagDesc.Desc2Line1
		data.LikeTagDescSecond2 = tagDesc.Desc2Line2
	} else if subTidBest > 0 {
		if tagRegDesc, ok := s.tagRegionDescMap[subTidBest]; ok && tagRegDesc != nil {
			data.LikeTagDescFirst = tagRegDesc.Desc1
			data.LikeTagDescSecond = tagRegDesc.Desc2Line1
			data.LikeTagDescSecond2 = tagRegDesc.Desc2Line2
		}
	}
	if braTime, e := time.Parse("20060102", item.BrainwashCirTime); e != nil {
		log.Warn("groupTmData BrainwashCirTime time.Parse(%s) warn(%v)", item.BrainwashCirTime, e)
	} else {
		data.BrainwashCirTime = braTime.Format("2006.01.02")
	}
	// LikeUpDuration data fix
	if data.LikeUpDuration > item.DurationHour*60 {
		if data.LikeUpDuration < (item.AvDurationHour*60 + item.LikeUpLiveDuration) {
			if (item.AvDurationHour*60 + item.LikeUpLiveDuration) > 0 {
				data.LikeUpDuration = int64(float64(item.DurationHour*60) * (float64(data.LikeUpDuration) / float64(item.AvDurationHour*60+item.LikeUpLiveDuration)))
			}
		} else {
			if (item.PlayDurationHourRep*60 + item.LikeUpLiveDuration) > 0 {
				data.LikeUpDuration = int64(float64(item.DurationHour*60) * (float64(item.LikeUpAvDurationRep+item.LikeUpLiveDuration) / float64(item.PlayDurationHourRep*60+item.LikeUpLiveDuration)))
			}
		}
	}
	realIP := metadata.String(c, metadata.RemoteIP)
	group := errgroup.WithCancel(c)
	// tag data
	if item.LikeTagID > 0 {
		group.Go(func(ctx context.Context) error {
			if tag, e := s.tagClient.Tag(ctx, &tagmdl.TagReq{Tid: item.LikeTagID}); e != nil || tag == nil {
				log.Error("s.tagClient.Tag tid(%d) error(%v)", item.LikeTagID, e)
			} else {
				data.LikeTagID = item.LikeTagID
				data.LikeTagName = tag.Tag.Name
			}
			return nil
		})
	}
	// group aids
	if item.BestAvid > 0 {
		aids = append(aids, item.BestAvid)
	}
	if item.LikesUgc3Avids != "" {
		list := strings.Split(item.LikesUgc3Avids, ",")
		for _, v := range list {
			items := strings.Split(v, ":")
			if len(items) == 2 {
				aid, e := strconv.ParseInt(items[0], 10, 64)
				if e != nil {
					continue
				}
				vv, e := strconv.ParseInt(items[1], 10, 64)
				if e != nil {
					continue
				}
				ugcView = append(ugcView, &model.AidView{Aid: aid, View: vv})
			}
		}
	}
	if item.LikePgc3Avids != "" {
		list := strings.Split(item.LikePgc3Avids, ",")
		for _, v := range list {
			items := strings.Split(v, "@")
			if len(items) == 2 {
				aid, e := strconv.ParseInt(items[0], 10, 64)
				if e != nil {
					continue
				}
				vv, e := strconv.ParseInt(items[1], 10, 64)
				if e != nil {
					continue
				}
				pgcView = append(pgcView, &model.AidView{Aid: aid, View: vv})
			}
		}
	}
	ugcViewLen := len(ugcView)
	if ugcViewLen > 3 {
		ugcView = ugcView[:_totalViewLen]
	}
	sort.Slice(ugcView, func(i, j int) bool {
		return ugcView[i].View > ugcView[j].View
	})
	pgcViewLen := len(pgcView)
	if pgcViewLen > 3 {
		pgcView = pgcView[:_totalViewLen]
	}
	sort.Slice(pgcView, func(i, j int) bool {
		return pgcView[i].View > pgcView[j].View
	})
	totalView = append(totalView, ugcView...)
	switch {
	case ugcViewLen == 0:
		totalView = append(totalView, pgcView...)
	case ugcViewLen == 1:
		switch {
		case pgcViewLen == 1:
			totalView = append(totalView, pgcView[0])
		case pgcViewLen > 1:
			totalView = append(totalView, pgcView[:2]...)
		}
	case ugcViewLen >= 2:
		if pgcViewLen > 0 {
			totalView = append(totalView, pgcView[0])
		}
	}
	sort.Slice(totalView, func(i, j int) bool {
		return totalView[i].View > totalView[j].View
	})
	if len(totalView) > _totalViewLen {
		totalView = totalView[:_totalViewLen]
	}
	for _, v := range totalView {
		viewAids = append(viewAids, v.Aid)
	}
	aids = append(aids, viewAids...)
	if item.LikeUp3Avs != "" {
		aidsStr := strings.Split(item.LikeUp3Avs, ",")
		for _, aidStr := range aidsStr {
			if aid, e := strconv.ParseInt(aidStr, 10, 64); e != nil {
				continue
			} else {
				upAids = append(upAids, aid)
			}
		}
	}
	aids = append(aids, upAids...)
	if item.BrainwashCirAvid > 0 {
		aids = append(aids, item.BrainwashCirAvid)
	}
	if item.BestAvid > 0 {
		switch item.BestAvidType {
		case _typeArticle:
			artIDs = append(artIDs, item.BestAvid)
		case _typeArchive:
			aids = append(aids, item.BestAvid)
		}
	}
	if item.BestAvidOld > 0 {
		switch item.BestAvidOldType {
		case _typeArticle:
			artIDs = append(artIDs, item.BestAvidOld)
		case _typeArchive:
			aids = append(aids, item.BestAvidOld)
		}
	}
	if item.FirstSubmitAvid > 0 {
		switch item.FirstSubmitType {
		case _typeArticle:
			artIDs = append(artIDs, item.FirstSubmitAvid)
		case _typeArchive:
			aids = append(aids, item.FirstSubmitAvid)
		}
	}
	if len(aids) > 0 {
		group.Go(func(ctx context.Context) error {
			s.groupArcData(ctx, aids, upAids, totalView, item, data)
			return nil
		})
	}
	// article
	if len(artIDs) > 0 {
		group.Go(func(ctx context.Context) error {
			s.groupArtData(ctx, artIDs, item, data)
			return nil
		})
	}
	group.Go(func(ctx context.Context) error {
		if acc, e := s.accClient.ProfileWithStat3(ctx, &accmdl.MidReq{Mid: item.Mid}); e != nil {
			log.Error("groupTmData s.arcClient.ProfileWithStat3(%v) error(%v)", item.Mid, e)
		} else {
			data.Uname = acc.Profile.Name
			data.Face = acc.Profile.Face
			data.Fans = acc.Follower
			data.RegTime = time.Unix(int64(acc.Profile.JoinTime), 0).Format("2006.01.02")
			data.RegDay = (time.Now().Unix() - int64(acc.Profile.JoinTime)) / 86400
		}
		return nil
	})
	// group mids
	if item.LikeBestUpID > 0 {
		mids = append(mids, item.LikeBestUpID)
	}
	if item.IsUp == 1 {
		if firstSubTime, e := time.Parse("2006-01-02 15:04:05", item.FirstSubmitTime); e != nil {
			log.Warn("groupTmData FirstSubmitTime time.Parse(%s) warn(%v)", item.FirstSubmitTime, e)
		} else {
			data.FirstSubmitTime = firstSubTime.Format("2006.01.02")
		}
		if item.UpBestFanVv != "" {
			list := strings.Split(item.UpBestFanVv, "@")
			if len(list) == 2 {
				if mid, e := strconv.ParseInt(list[0], 10, 64); e != nil {
					log.Error("UpBestFanVv parse(%s) error(%v)", list[0], e)
				} else {
					upBestFanMid = mid
					mids = append(mids, mid)
				}
				if vv, e := strconv.ParseInt(list[1], 10, 64); e != nil {
					log.Error("UpBestFanLiveMinute parse(%s) error(%v)", list[1], e)
				} else {
					favVv = vv
				}
			}
		}
	}
	if item.IsLiveUp == 1 {
		if maxCdnTime, e := time.Parse("20060102", item.MaxCdnNumDate); e != nil {
			log.Warn("groupTmData MaxCdnNumDate time.Parse(%s) warn(%v)", item.MaxCdnNumDate, e)
		} else {
			data.MaxCdnNumDate = maxCdnTime.Format("2006.01.02")
		}
		data.MaxCdnNum = item.MaxCdnNum
		if item.UpBestFanLiveMinute != "" {
			list := strings.Split(item.UpBestFanLiveMinute, "@")
			if len(list) == 2 {
				if mid, e := strconv.ParseInt(list[0], 10, 64); e != nil {
					log.Error("UpBestFanLiveMinute parse(%s) error(%v)", list[0], e)
				} else {
					upBestLiveFanMid = mid
					mids = append(mids, mid)
				}
				if minute, e := strconv.ParseInt(list[1], 10, 64); e != nil {
					log.Error("UpBestFanLiveMinute parse(%s) error(%v)", list[1], e)
				} else {
					liveMinute = minute
				}
			}
		}
	}
	if len(mids) > 0 {
		group.Go(func(ctx context.Context) error {
			if accs, e := s.accClient.Infos3(ctx, &accmdl.MidsReq{Mids: mids, RealIp: realIP}); e != nil || accs.Infos == nil {
				log.Error("groupTmData s.accClient.Cards3(%v) error(%v)", mids, e)
			} else {
				if item.LikeBestUpID > 0 {
					if info, ok := accs.Infos[item.LikeBestUpID]; ok {
						data.LikeBestUpID = item.LikeBestUpID
						data.LikeBestUpName = info.Name
						data.LikeBestUpFace = info.Face
					}
				}
				if item.IsUp == 1 && upBestFanMid > 0 {
					if info, ok := accs.Infos[upBestFanMid]; ok {
						data.UpBestFanVv = &model.FavVv{Mid: upBestFanMid, Name: info.Name, Face: info.Face, Vv: favVv}
					}
				}
				if item.IsLiveUp == 1 && upBestLiveFanMid > 0 {
					if info, ok := accs.Infos[upBestLiveFanMid]; ok {
						data.UpBestFanLiveMinute = &model.FanMinute{Mid: upBestLiveFanMid, Name: info.Name, Face: info.Face, Minute: liveMinute}
					}
				}
			}
			return nil
		})
	}
	if e := group.Wait(); e != nil {
		log.Error("groupTmData group.Wait error(%v)", e)
	}
	if len(data.LikeUp3Arcs) == 0 {
		data.LikeUp3Arcs = make([]*model.TmArc, 0)
	}
	return
}

func (s *Service) groupArcData(c context.Context, aids, upAids []int64, totalView []*model.AidView, item *model.Item, data *model.Timemachine) {
	if arcs, e := s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); e != nil {
		log.Error("groupTmData s.arcClient.Arcs(%v) error(%v)", aids, e)
	} else if arcs != nil {
		if item.BrainwashCirAvid > 0 {
			if arc, ok := arcs.Arcs[item.BrainwashCirAvid]; ok && arc.IsNormal() {
				data.BrainwashCirArc = &model.TmArc{Aid: item.BrainwashCirAvid, Title: arc.Title, Cover: arc.Pic, Author: arc.Author}
			}
		}
		if item.BestAvid > 0 && item.BestAvidType == _typeArchive {
			if arc, ok := arcs.Arcs[item.BestAvid]; ok && arc.IsNormal() {
				data.BestArc = &model.TmArc{Aid: item.BestAvid, Title: arc.Title, Cover: arc.Pic, Author: arc.Author}
			}
		}
		if item.BestAvidOld > 0 && item.BestAvidOldType == _typeArchive {
			if arc, ok := arcs.Arcs[item.BestAvidOld]; ok && arc.IsNormal() {
				data.BestArcOld = &model.TmArc{Aid: item.BestAvidOld, Title: arc.Title, Cover: arc.Pic, Author: arc.Author}
			}
		}
		if item.FirstSubmitAvid > 0 && item.FirstSubmitType == _typeArchive {
			if arc, ok := arcs.Arcs[item.FirstSubmitAvid]; ok && arc.IsNormal() {
				data.FirstSubmitArc = &model.TmArc{Aid: item.FirstSubmitAvid, Title: arc.Title, Cover: arc.Pic, Author: arc.Author}
			}
		}
		for _, aid := range upAids {
			if arc, ok := arcs.Arcs[aid]; ok && arc.IsNormal() {
				data.LikeUp3Arcs = append(data.LikeUp3Arcs, &model.TmArc{Aid: aid, Title: arc.Title, Cover: arc.Pic, Author: arc.Author})
			}
		}
		for _, item := range totalView {
			if arc, ok := arcs.Arcs[item.Aid]; ok && arc.IsNormal() {
				data.Likes3Arcs = append(data.Likes3Arcs, &model.TmArc{Aid: item.Aid, Title: arc.Title, Cover: arc.Pic, Author: arc.Author})
			}
		}
	}
}

func (s *Service) groupArtData(c context.Context, artIDs []int64, item *model.Item, data *model.Timemachine) {
	if arts, e := s.article.ArticleMetas(c, &artmdl.ArgAids{Aids: artIDs}); e != nil {
		log.Error("groupTmData s.article.ArticleMetas(%v) error(%v)", artIDs, e)
	} else {
		if item.BestAvid > 0 && item.BestAvidType == _typeArticle {
			if art, ok := arts[item.BestAvid]; ok {
				data.BestArc = &model.TmArc{Aid: item.BestAvid, Title: art.Title}
			}
		}
		if item.BestAvidOld > 0 && item.BestAvidOldType == _typeArticle {
			if art, ok := arts[item.BestAvidOld]; ok {
				data.BestArcOld = &model.TmArc{Aid: item.BestAvidOld, Title: art.Title}
			}
		}
		if item.FirstSubmitAvid > 0 && item.FirstSubmitType == _typeArticle {
			if art, ok := arts[item.FirstSubmitAvid]; ok {
				data.FirstSubmitArc = &model.TmArc{Aid: item.FirstSubmitAvid, Title: art.Title}
			}
		}
	}
}

func (s *Service) groupTagDesc(item *model.Item, data *model.Timemachine) (subTidBest int64) {
	var subTid int64
	subList := strings.Split(item.Like2SubTids, ",")
	if len(subList) == 2 {
		if subID, e := strconv.ParseInt(subList[0], 10, 64); e != nil {
			log.Warn("groupTmData Like2SubTids time.Parse(%s) warn(%v)", subList[0], e)
		} else if subID > 0 {
			subTidBest = subID
		}
		if subID, e := strconv.ParseInt(subList[1], 10, 64); e != nil {
			log.Warn("groupTmData Like2SubTids strconv.ParseInt(%s) warn(%v)", subList[1], e)
		} else if subID > 0 {
			subTid = subID
		}
		if regionDesc, ok := s.regionDescMap[subTidBest]; ok && regionDesc != nil {
			data.LikeSubDesc2 = regionDesc.Desc2
			data.LikeSubDesc3 = regionDesc.Desc3
		}
		if regionDesc, ok := s.regionDescMap[subTid]; ok && regionDesc != nil {
			data.LikeSubDesc1 = regionDesc.Desc1
		}
	}
	return
}
