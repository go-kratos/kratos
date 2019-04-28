package service

import (
	"context"
	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

// getSubtitlDraft get a subtitle
// cache throught
func (s *Service) getSubtitlDraft(c context.Context, oid int64, tp int32, mid int64, lanCode uint8) (draft *model.Subtitle, err error) {
	var (
		cacheErr bool
	)
	if draft, err = s.dao.SubtitleDraftCache(c, oid, tp, mid, lanCode); err != nil {
		cacheErr = true
		err = nil
	}
	if draft != nil {
		if draft.ID <= 0 {
			draft = nil
			err = nil
		}
		return
	}
	if draft, err = s.dao.GetSubtitleDraft(c, oid, tp, mid, lanCode); err != nil {
		log.Error("params(oid:%v,tp:%v,mid:%v,lanCode:%v).error(%v)", oid, tp, mid, lanCode, err)
		return
	}
	if draft == nil {
		draft = &model.Subtitle{
			Oid:  oid,
			Type: tp,
			Mid:  mid,
			Lan:  lanCode,
		}
	}
	if !cacheErr {
		temp := draft
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetSubtitleDraftCache(ctx, temp)
		})
	}
	if draft.ID <= 0 {
		draft = nil
		err = nil
	}
	return
}

// GetSubtitle get a subtitle
func (s *Service) getSubtitle(c context.Context, oid int64, subtitleID int64) (subtitle *model.Subtitle, err error) {
	var (
		cacheErr bool
	)
	if subtitle, err = s.dao.SubtitleCache(c, oid, subtitleID); err != nil {
		cacheErr = true
		err = nil
	}
	if subtitle != nil {
		if subtitle.Empty {
			subtitle = nil
			err = nil
		}
		return
	}
	if subtitle, err = s.dao.GetSubtitle(c, oid, subtitleID); err != nil {
		log.Error("params(oid:%v, subtitleID:%v).error(%v)", oid, subtitleID, err)
		return
	}
	if subtitle == nil {
		subtitle = &model.Subtitle{
			Oid:   oid,
			ID:    subtitleID,
			Empty: true,
		}
	}
	if !cacheErr {
		temp := subtitle
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetSubtitleCache(ctx, temp)
		})
	}
	if subtitle.Empty {
		subtitle = nil
		err = nil
	}
	return
}

// getSubtitles 不保证顺序
func (s *Service) getSubtitles(c context.Context, oid int64, subtitleIds []int64) (subtitles map[int64]*model.Subtitle, err error) {
	var (
		hits            map[int64]*model.Subtitle
		missed          []int64
		cacheErr        bool
		missedSubtitles []*model.Subtitle
	)
	if hits, missed, err = s.dao.SubtitlesCache(c, oid, subtitleIds); err != nil {
		cacheErr = true
		err = nil
	}
	subtitles = make(map[int64]*model.Subtitle)
	for _, subtitle := range hits {
		if subtitle.Empty {
			missed = append(missed, subtitle.ID)
			continue
		}
		subtitles[subtitle.ID] = subtitle
	}
	if len(missed) > 0 {
		if missedSubtitles, err = s.dao.GetSubtitles(c, oid, missed); err != nil {
			log.Error("getSubtitles(oid:%v,subtitleIds:%v),error(%v)", oid, subtitleIds, err)
			return
		}
	}
	for _, subtitle := range missedSubtitles {
		subtitles[subtitle.ID] = subtitle
	}
	if !cacheErr {
		for _, subtitle := range missedSubtitles {
			temp := subtitle
			s.cache.Do(c, func(ctx context.Context) {
				s.dao.SetSubtitleCache(ctx, temp)
			})
		}
	}
	return
}
