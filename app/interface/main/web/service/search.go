package service

import (
	"context"

	"go-common/app/interface/main/web/model"
	accmdl "go-common/app/service/main/account/api"
	relmdl "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const _searchEggWebPlat = 6

var _emptyUpRec = make([]*model.UpRecInfo, 0)

// SearchAll search type all.
func (s *Service) SearchAll(c context.Context, mid int64, arg *model.SearchAllArg, buvid, ua, typ string) (data *model.Search, err error) {
	data, err = s.dao.SearchAll(c, mid, arg, buvid, ua, typ)
	return
}

// SearchByType type video,bangumi,pgc,live,live_user,article,special,topic,bili_user,photo
func (s *Service) SearchByType(c context.Context, mid int64, arg *model.SearchTypeArg, buvid, ua string) (res *model.SearchTypeRes, err error) {
	switch arg.SearchType {
	case model.SearchTypeVideo:
		if res, err = s.dao.SearchVideo(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeBangumi:
		if res, err = s.dao.SearchBangumi(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypePGC:
		if res, err = s.dao.SearchPGC(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeLive:
		if res, err = s.dao.SearchLive(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeLiveRoom:
		if res, err = s.dao.SearchLiveRoom(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeLiveUser:
		if res, err = s.dao.SearchLiveUser(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeArticle:
		if res, err = s.dao.SearchArticle(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeSpecial:
		if res, err = s.dao.SearchSpecial(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeTopic:
		if res, err = s.dao.SearchTopic(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypeUser:
		if res, err = s.dao.SearchUser(c, mid, arg, buvid, ua); err != nil {
			return
		}
	case model.SearchTypePhoto:
		if res, err = s.dao.SearchPhoto(c, mid, arg, buvid, ua); err != nil {
			return
		}
	default:
		err = ecode.RequestErr
		return
	}
	return
}

// SearchRec search recommend data.
func (s *Service) SearchRec(c context.Context, mid int64, pn, ps int, keyword, fromSource, buvid, ua string) (data *model.SearchRec, err error) {
	data, err = s.dao.SearchRec(c, mid, pn, ps, keyword, fromSource, buvid, ua)
	return
}

// SearchDefault get search default word.
func (s *Service) SearchDefault(c context.Context, mid int64, fromSource, buvid, ua string) (data *model.SearchDefault, err error) {
	data, err = s.dao.SearchDefault(c, mid, fromSource, buvid, ua)
	return
}

// UpRec get up recommend
func (s *Service) UpRec(c context.Context, mid int64, arg *model.SearchUpRecArg, buvid string) (data *model.UpRecData, err error) {
	var (
		ups        []*model.SearchUpRecRes
		trackID    string
		mids       []int64
		cardsReply *accmdl.CardsReply
		cardErr    error
	)
	if ups, trackID, err = s.dao.UpRecommend(c, mid, arg, buvid); err != nil {
		return
	}
	data = &model.UpRecData{TrackID: trackID}
	if len(ups) == 0 {
		data.List = _emptyUpRec
		return
	}
	for _, v := range ups {
		mids = append(mids, v.UpID)
	}
	relInfos := make(map[int64]*relmdl.Stat, len(mids))
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if cardsReply, cardErr = s.accClient.Cards3(errCtx, &accmdl.MidsReq{Mids: mids}); cardErr != nil {
			log.Error("UpRec s.accClient.Cards3(%v) error(%v)", mids, cardErr)
			return cardErr
		}
		return nil
	})
	group.Go(func() error {
		if relReply, relErr := s.relation.Stats(errCtx, &relmdl.ArgMids{Mids: mids}); relErr != nil {
			log.Error("UpRec s.relation.Stats(%d,%v) error(%v)", mid, mids, relErr)
		} else if relReply != nil {
			relInfos = relReply
		}
		return nil
	})
	if err = group.Wait(); err != nil {
		return
	}
	for _, v := range ups {
		if info, ok := cardsReply.Cards[v.UpID]; ok && info != nil && info.Silence == 0 {
			upInfo := &model.UpRecInfo{
				Mid:       info.Mid,
				Name:      info.Name,
				Face:      info.Face,
				Official:  info.Official,
				RecReason: v.RecReason,
				Tid:       v.Tid,
				SecondTid: v.SecondTid,
				Sign:      info.Sign,
			}
			upInfo.Vip.Type = info.Vip.Type
			upInfo.Vip.Status = info.Vip.Status
			if stat, ok := relInfos[v.UpID]; ok {
				upInfo.Follower = stat.Follower
			}
			if typ, ok := s.typeNames[int32(v.Tid)]; ok {
				upInfo.Tname = typ.Name
			}
			if typ, ok := s.typeNames[int32(v.SecondTid)]; ok {
				upInfo.SecondTname = typ.Name
			}
			data.List = append(data.List, upInfo)
		}
	}
	if len(data.List) == 0 {
		data.List = _emptyUpRec
	}
	return
}

//SearchEgg get search egg by egg id.
func (s *Service) SearchEgg(c context.Context, eggID int64) (data *model.SearchEggRes, err error) {
	if _, ok := s.searchEggs[eggID]; !ok {
		err = ecode.NothingFound
		return
	}
	data = s.searchEggs[eggID]
	return
}
