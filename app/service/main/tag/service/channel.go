package service

import (
	"context"

	"go-common/app/service/main/tag/model"
)

// ChannelCategories get channel categories.
func (s *Service) ChannelCategories(c context.Context, arg *model.ArgChannelCategory) (res []*model.ChannelCategory, err error) {
	return s.dao.ChannelCategories(c, arg)
}

// Channels get channels by ids.
func (s *Service) Channels(c context.Context, arg *model.ArgChannels) (res []*model.Channel, err error) {
	return s.dao.Channels(c, arg)
}

// ChannelRule get channel rules by pn,ps.
func (s *Service) ChannelRule(c context.Context, arg *model.ArgChannelRule) (res []*model.ChannelRule, err error) {
	return s.dao.ChannelRule(c, arg)
}

// ChannelGroup ChannelGroup.
func (s *Service) ChannelGroup(c context.Context, tid int64) (res []*model.ChannelGroup, err error) {
	if res, err = s.dao.ChannelGroupCache(c, tid); err != nil || res != nil {
		return
	}
	if res, err = s.dao.ChannelGroup(c, tid); err != nil {
		return
	}
	cgs := make([]*model.ChannelGroup, 0, len(res))
	for _, v := range res {
		cg := &model.ChannelGroup{}
		*cg = *v
		cgs = append(cgs, cg)
	}
	s.cacheCh.Save(func() {
		s.dao.AddChannelGroupCache(context.Background(), tid, cgs)
	})
	return
}
