package service

import (
	"context"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	thumbupmdl "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) isLike(c context.Context, mid, aid int64) (res int8, err error) {
	r, err := s.HadLikesByMid(c, mid, []int64{aid})
	if err != nil {
		return
	}
	res = r[aid]
	return
}

// HadLikesByMid .
func (s *Service) HadLikesByMid(c context.Context, mid int64, aids []int64) (res map[int64]int8, err error) {
	if mid == 0 || len(aids) == 0 {
		return
	}
	arg := &thumbupmdl.ArgHasLike{Business: "article", MessageIDs: aids, Mid: mid}
	res, err = s.thumbupRPC.HasLike(c, arg)
	return
}

// Like like article
func (s *Service) Like(c context.Context, mid, aid int64, likeType int) (err error) {
	var art *artmdl.Meta
	if (likeType < 0) || (likeType > 4) {
		err = ecode.RequestErr
		return
	}
	if art, err = s.ArticleMeta(c, aid); err != nil || art == nil {
		err = ecode.NothingFound
		return
	}
	arg := &thumbupmdl.ArgLike{
		Mid:       mid,
		UpMid:     art.Author.Mid,
		Business:  "article",
		MessageID: aid,
		Type:      int8(likeType),
	}
	if err = s.thumbupRPC.Like(c, arg); err != nil {
		dao.PromError("like:thumbup-service")
		log.Error("s.thumbupRPC.Like(%+v) err: %+v", arg, err)
	}
	return
}

// RecommendsWithLike recommends with like state
func (s *Service) RecommendsWithLike(c context.Context, cid int64, pn, ps int, lastAids []int64, sort int, mid int64) (res []*artmdl.RecommendArtWithLike, err error) {
	var recs []*artmdl.RecommendArt
	if recs, err = s.Recommends(c, cid, pn, ps, lastAids, sort); err != nil {
		return
	}
	var aids []int64
	for _, rec := range recs {
		aids = append(aids, rec.ID)
	}
	states, _ := s.HadLikesByMid(c, mid, aids)
	for _, rec := range recs {
		r := &artmdl.RecommendArtWithLike{RecommendArt: *rec}
		if states != nil {
			r.LikeState = int(states[rec.ID])
		}
		res = append(res, r)
	}
	return
}
