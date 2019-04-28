package service

import (
	"context"
	"fmt"

	"go-common/app/interface/openplatform/article/dao"
	"go-common/app/interface/openplatform/article/model"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_serviceArea = "article_info"
	_recType     = "up_rec"
	_rapagesize  = 20
)

// UpdateAuthorCache update author cache
func (s *Service) UpdateAuthorCache(c context.Context, mid int64) (err error) {
	err = s.dao.DelCacheAuthor(c, mid)
	dao.PromInfo("author:更新作者状态")
	return
}

// AddAuthor add author to db
func (s *Service) AddAuthor(c context.Context, mid int64) (err error) {
	if len(s.activities) == 0 {
		err = ecode.ArtNoActivity
		return
	}
	res, forbid, _ := s.IsAuthor(c, mid)
	if res {
		return
	}
	if forbid {
		err = ecode.ArtUserDisabled
		return
	}
	author, err := s.dao.RawAuthor(c, mid)
	if err != nil {
		return
	}
	if author != nil {
		if (author.State == model.AuthorStateReject) || (author.State == model.AuthorStateClose) {
			err = ecode.ArtAuthorReject
			return
		}
	}
	err = s.dao.AddAuthor(c, mid)
	if err == nil {
		s.dao.DelCacheAuthor(c, mid)
	}
	return
}

// IsAuthor check that whether user has permission to write article.
func (s *Service) IsAuthor(c context.Context, mid int64) (res bool, forbid bool, err error) {
	var level int
	forbid, level, _ = s.UserDisabled(c, mid)
	if forbid {
		return
	}
	var limit *model.AuthorLimit
	limit, _ = s.dao.Author(c, mid)
	if limit.Pass() {
		res = true
		return
	}
	if limit.Forbid() {
		return
	}
	if level >= 2 {
		res = true
		return
	}
	res, err = s.isUpper(c, mid)
	return
}

// Authors recommends similar authors by categories
func (s *Service) Authors(c context.Context, mid int64, author int64) (res []*model.AccountCard) {
	var (
		categories    []int64
		authors       []int64
		attentions    map[int64]bool
		filterAuthors []int64
		blacks        map[int64]struct{}
		mids          []int64
		err           error
	)
	if categories, err = s.dao.AuthorMostCategories(c, author); err != nil {
		return
	}
	for _, category := range categories {
		if authors, err = s.dao.CategoryAuthors(c, category, s.c.Article.RecommendAuthors); err != nil {
			return
		}
		if attentions, err = s.isAttentions(c, mid, authors); err != nil {
			return
		}
		if blacks, err = s.isBlacks(c, mid, authors); err != nil {
			return
		}
		for _, a := range authors {
			if _, ok := blacks[a]; !attentions[a] && !ok {
				filterAuthors = append(filterAuthors, a)
			}
		}
	}
	for _, filterAuthor := range filterAuthors {
		if filterAuthor == mid || filterAuthor == author {
			continue
		}
		if forbid, _, _ := s.UserDisabled(c, filterAuthor); !forbid {
			mids = append(mids, filterAuthor)
		}
	}
	if len(mids) < 2 {
		return
	}
	for _, m := range mids {
		var (
			accountCard = &model.AccountCard{}
			profile     *account.ProfileStat
		)
		if profile, err = s.accountRPC.ProfileWithStat3(c, &account.ArgMid{Mid: m}); err != nil {
			dao.PromError("article:ProfileWithStat3")
			log.Error("s.acc.ProfileWithStat3(%d) error %v", m, err)
			return
		}
		if profile != nil {
			accountCard.FromProfileStat(profile)
		}
		res = append(res, accountCard)
	}
	if len(res) > 3 {
		res = res[0:3]
	}
	return
}

// LevelRequired .
func (s *Service) LevelRequired(c context.Context, mid int64) (ok bool, err error) {
	var (
		card *account.Card
		arg  = account.ArgMid{Mid: mid}
	)
	if card, err = s.accountRPC.Card3(c, &arg); err != nil {
		dao.PromError("accountRPC.Card3")
		log.Error("s.LevelRequired.accountRPC.Card3(%d) error %v", mid, err)
		return
	}
	if card.Level >= 4 {
		ok = true
	}
	return
}

// RecommendAuthors get recommends from search.
func (s *Service) RecommendAuthors(c context.Context, platform string, mobiApp string, device string, build int, clientIP string, userID int64, buvid string, mid int64) (res *model.RecommendAuthors, err error) {
	var ra []*model.RecommendAuthor
	if ra, err = s.dao.RecommendAuthors(c, platform, mobiApp, device, build, clientIP, userID, buvid, _recType, _serviceArea, _rapagesize, mid); err != nil {
		return
	}
	if ra == nil {
		return
	}
	res = &model.RecommendAuthors{}
	for _, r := range ra {
		var (
			author  = &model.RecAuthor{}
			profile *account.ProfileStat
		)
		if profile, err = s.accountRPC.ProfileWithStat3(c, &account.ArgMid{Mid: r.UpID}); err != nil {
			dao.PromError("article:ProfileWithStat3")
			log.Error("s.acc.ProfileWithStat3(%d) error %v", r.UpID, err)
			return
		}
		if profile != nil {
			author.AccountCard.FromProfileStat(profile)
		}
		if r.RecReason == "" {
			var (
				view             int64
				fans             int
				viewStr, fansStr string
			)
			fans = author.Fans
			if view, err = s.readCount(c, r.UpID); err != nil {
				continue
			}
			if fans < 1e4 {
				fansStr = fmt.Sprintf("粉丝：%d", fans)
			} else if fans < 1e8 {
				fansStr = fmt.Sprintf("粉丝：%.1f万", float64(fans/1e4))
			} else {
				fansStr = fmt.Sprintf("粉丝：%.1f亿", float64(fans/1e8))
			}
			if view < 1e4 {
				viewStr = fmt.Sprintf("阅读：%d", view)
			} else if view < 1e8 {
				viewStr = fmt.Sprintf("阅读：%.1f万", float64(view/1e4))
			} else {
				viewStr = fmt.Sprintf("阅读：%.1f亿", float64(view/1e8))
			}
			r.RecReason = fansStr + "  " + viewStr
		}
		author.RecReason = r.RecReason
		res.Authors = append(res.Authors, author)
	}
	res.Count = len(res.Authors)
	return
}

func (s *Service) readCount(c context.Context, mid int64) (res int64, err error) {
	var (
		stat *model.UpStat
		st   model.UpStat
	)
	if stat, err = s.dao.CacheUpStatDaily(c, mid); err != nil {
		dao.PromError("article:CacheUpStatDaily")
	}
	if stat != nil {
		res = stat.View
		return
	}
	if st, err = s.dao.UpStat(c, mid); err != nil {
		dao.PromError("article:获取作者文章数")
	}
	res = st.View
	return
}
