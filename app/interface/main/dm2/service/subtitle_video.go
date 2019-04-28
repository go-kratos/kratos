package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

func keySubtitleSingle(oid int64, tp int32) string {
	return fmt.Sprintf("subtitle_%d_%d", oid, tp)
}

// GetWebVideoSubtitle .
func (s *Service) GetWebVideoSubtitle(c context.Context, aid, oid int64, tp int32) (res *model.VideoSubtitles, err error) {
	var (
		subtitles       []*model.VideoSubtitle
		subtitleSubject *model.SubtitleSubject
		allowSubmit     bool
		closed          bool
		lan, lanDoc     string
	)
	if subtitleSubject, err = s.subtitleSubject(c, aid); err != nil {
		log.Error("params(aid:%v).err(%v)", aid, err)
		err = nil
	}
	if subtitleSubject != nil {
		allowSubmit = subtitleSubject.Allow
		closed = subtitleSubject.AttrVal(model.AttrSubtitleClose) == model.AttrYes
		lan, lanDoc = s.subtitleLans.GetByID(int64(subtitleSubject.Lan))
	}
	res = &model.VideoSubtitles{
		AllowSubmit: allowSubmit,
		Lan:         lan,
		LanDoc:      lanDoc,
	}
	if closed {
		res.AllowSubmit = false
		return
	}
	if subtitles, err = s.getVideoSubtitles(c, oid, tp); err != nil {
		return
	}
	res.Subtitles = subtitles
	return
}

// singleGetVideoSubtitle use singleflight, but not cache sub item
func (s *Service) singleGetVideoSubtitle(c context.Context, oid int64, tp int32) (res []*model.VideoSubtitle, err error) {
	var (
		v           interface{}
		subtitleIds []int64
		subtitles   map[int64]*model.Subtitle
	)
	v, err, _ = s.subtitleSingleGroup.Do(keySubtitleSingle(oid, tp), func() (reply interface{}, err error) {
		if subtitleIds, err = s.dao.GetSubtitleIds(c, oid, tp); err != nil {
			log.Error("params(oid:%v, tp:%v).err(%v)", oid, tp, err)
			return
		}
		if len(subtitleIds) == 0 {
			return
		}
		if subtitles, err = s.getSubtitles(c, oid, subtitleIds); err != nil {
			log.Error("params(oid:%v, subtitleIds:%v).err(%v)", oid, subtitleIds, err)
			return
		}
		result := make([]*model.VideoSubtitle, 0, len(subtitles))
		for _, subtitle := range subtitles {
			lan, lanDoc := s.subtitleLans.GetByID(int64(subtitle.Lan))
			vs := &model.VideoSubtitle{
				ID:          subtitle.ID,
				IsLock:      subtitle.IsLock,
				Lan:         lan,
				LanDoc:      lanDoc,
				SubtitleURL: subtitle.SubtitleURL,
			}
			if subtitle.IsSign {
				vs.AuthorMid = subtitle.AuthorID
			}
			result = append(result, vs)
		}
		reply = result
		return
	})
	if err != nil {
		log.Error("params(oid:%v, tp:%v).err(%v)", oid, tp, err)
		return
	}
	res, _ = v.([]*model.VideoSubtitle)
	return
}

// getVideoSubtitles get from cache
func (s *Service) getVideoSubtitles(c context.Context, oid int64, tp int32) (subtitles []*model.VideoSubtitle, err error) {
	var (
		cacheErr           bool
		videoSubtitleCache *model.VideoSubtitleCache
	)
	if videoSubtitleCache, err = s.dao.VideoSubtitleCache(c, oid, tp); err != nil {
		cacheErr = true
		err = nil
	}
	if videoSubtitleCache != nil {
		subtitles = videoSubtitleCache.VideoSubtitles
		return
	}
	if subtitles, err = s.singleGetVideoSubtitle(c, oid, tp); err != nil {
		log.Error("params(oid:%v,tp:%v).err(%v)", oid, tp, err)
		return
	}
	videoSubtitleCache = &model.VideoSubtitleCache{
		VideoSubtitles: subtitles,
	}
	if !cacheErr {
		temp := videoSubtitleCache
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetVideoSubtitleCache(ctx, oid, tp, temp)
		})
	}
	return
}
