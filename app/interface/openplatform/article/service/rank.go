package service

import (
	"context"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
)

// RankCategories rank categoires
func (s *Service) RankCategories(c context.Context) (res []*artmdl.RankCategory) {
	res = s.c.RankCategories
	return
}

// Ranks get ranks
func (s *Service) Ranks(c context.Context, cid int64, mid int64, ip string) (res []*artmdl.RankMeta, note string, err error) {
	var (
		exist    bool
		addCache = true
		aids     []int64
		rank     artmdl.RankResp
		metas    map[int64]*artmdl.Meta
	)
	if !s.ranksMap[cid] {
		err = ecode.RequestErr
		return
	}
	if exist, err = s.dao.ExpireRankCache(c, cid); err != nil {
		addCache = false
		err = nil
	}
	if exist {
		if rank, err = s.dao.RankCache(c, cid); err != nil {
			exist = false
			err = nil
			addCache = false
		}
	}
	if !exist {
		if rank, err = s.dao.Rank(c, cid, ip); err != nil {
			if rank, err = s.dao.RankCache(c, cid); err != nil {
				return
			}
		} else {
			if addCache && len(rank.List) > 0 {
				cache.Save(func() {
					s.dao.AddRankCache(context.TODO(), cid, rank)
				})
			}
		}
	}
	if len(rank.List) == 0 {
		return
	}
	for _, a := range rank.List {
		aids = append(aids, a.Aid)
	}
	if metas, err = s.ArticleMetas(c, aids); err != nil {
		return
	}
	var ups []int64
	for _, r := range rank.List {
		if metas[r.Aid] != nil {
			res = append(res, &artmdl.RankMeta{Meta: metas[r.Aid], Score: r.Score})
			ups = append(ups, metas[r.Aid].Author.Mid)

		}
	}
	if (len(ups) > 0) && (mid != 0) {
		if attentions, e := s.isAttentions(c, mid, ups); e == nil {
			for _, r := range res {
				r.Attention = attentions[r.Author.Mid]
			}
		}
	}
	if s.setting.ShowRankNote {
		note = rank.Note
	}
	return
}

func (s *Service) loadRanks() {
	for _, rank := range s.c.RankCategories {
		s.ranksMap[rank.ID] = true
	}
}
