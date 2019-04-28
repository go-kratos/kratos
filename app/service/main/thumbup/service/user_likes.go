package service

import (
	"context"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/sync/errgroup"
)

// UserLikes .
func (s *Service) UserLikes(c context.Context, business string, mid int64, pn, ps int) (res []*model.ItemLikeRecord, err error) {
	return s.userLikes(c, business, mid, pn, ps, model.StateLike)
}

// UserDislikes user dislikes
func (s *Service) UserDislikes(c context.Context, business string, mid int64, pn, ps int) (res []*model.ItemLikeRecord, err error) {
	return s.userLikes(c, business, mid, pn, ps, model.StateDislike)
}

func (s *Service) userLikes(c context.Context, business string, mid int64, pn, ps int, state int8) (res []*model.ItemLikeRecord, err error) {
	var businessID int64
	if businessID, err = s.CheckBusiness(business); err != nil {
		return
	}
	if err = s.checkUserLikeType(businessID, state); err != nil {
		return
	}
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1 // from cache, end-1
	)
	res, err = s.dao.UserLikeList(c, mid, businessID, state, start, end)
	return
}

// UserTotalLike user like list with total
func (s *Service) UserTotalLike(c context.Context, business string, mid int64, pn, ps int) (res *model.UserTotalLike, err error) {
	group, ctx := errgroup.WithContext(c)
	res = &model.UserTotalLike{}
	group.Go(func() (err error) {
		res.List, err = s.UserLikes(ctx, business, mid, pn, ps)
		return
	})
	group.Go(func() (err error) {
		res.Total, err = s.userTotal(ctx, business, mid)
		return
	})
	if err = group.Wait(); err != nil {
		res = nil
		return
	}
	return
}

// userTotal user total like
func (s *Service) userTotal(c context.Context, business string, mid int64) (res int, err error) {
	var businessID int64
	if businessID, err = s.CheckBusiness(business); err != nil {
		return
	}
	if exist, _ := s.dao.ExpireUserLikesCache(c, mid, businessID, model.StateLike); exist {
		res, err = s.dao.UserLikesCountCache(c, businessID, mid)
		return
	}
	res, err = s.dao.UserLikeCount(c, businessID, mid, model.StateLike)
	return
}
