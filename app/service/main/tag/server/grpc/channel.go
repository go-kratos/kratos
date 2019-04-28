package grpc

import (
	"context"

	v1 "go-common/app/service/main/tag/api"
	"go-common/app/service/main/tag/model"
	"go-common/library/ecode"
)

func (s *grpcServer) ChannelCategories(c context.Context, arg *v1.ChannelCategoriesReq) (res *v1.ChannelCategoriesReply, err error) {
	var (
		categories []*model.ChannelCategory
		req        = &model.ArgChannelCategory{
			LastID: arg.StartId,
			Size:   arg.Num,
			State:  arg.State,
		}
	)
	res = &v1.ChannelCategoriesReply{
		Categories: make([]*v1.ChannelCategory, 0, arg.Num),
	}
	if categories, err = s.svr.ChannelCategories(c, req); err != nil {
		return
	}
	for _, category := range categories {
		t := &v1.ChannelCategory{
			Id:        category.ID,
			Name:      category.Name,
			Order:     category.Order,
			State:     category.State,
			Ctime:     category.CTime,
			Mtime:     category.MTime,
			IntShield: category.AttrVal(model.ChannelCategoryAttrINT),
		}
		res.Categories = append(res.Categories, t)
	}
	return
}

func (s *grpcServer) Channels(c context.Context, arg *v1.ChannelsReq) (res *v1.ChannelsReply, err error) {
	var (
		channels []*model.Channel
		req      = &model.ArgChannels{
			LastID: arg.StartId,
			Size:   arg.Num,
		}
	)
	res = &v1.ChannelsReply{
		Channels: make([]*v1.Channel, 0, arg.Num),
	}
	if channels, err = s.svr.Channels(c, req); err != nil {
		return
	}
	for _, channel := range channels {
		t := &v1.Channel{
			Tid:      channel.ID,
			Type:     channel.Type,
			Rank:     channel.Rank,
			Operator: channel.Operator,
			State:    channel.State,
			Ctime:    channel.CTime,
			Mtime:    channel.MTime,
			Attr:     channel.Attr,
			TopRank:  channel.TopRank,
		}
		res.Channels = append(res.Channels, t)
	}
	return
}

func (s *grpcServer) ChannelRules(c context.Context, arg *v1.ChannelRulesReq) (res *v1.ChannelRulesReply, err error) {
	var (
		rules []*model.ChannelRule
		req   = &model.ArgChannelRule{
			LastID: arg.StartId,
			Size:   arg.Num,
			State:  arg.State,
		}
	)
	res = &v1.ChannelRulesReply{
		Rules: make([]*v1.ChannelRule, 0, arg.Num),
	}
	if rules, err = s.svr.ChannelRule(c, req); err != nil {
		return
	}
	for _, rule := range rules {
		r := &v1.ChannelRule{
			Id:        rule.Id,
			Tid:       rule.Tid,
			InRule:    rule.InRule,
			NotinRule: rule.NotinRule,
			State:     rule.State,
		}
		res.Rules = append(res.Rules, r)
	}
	return
}

func (s *grpcServer) ChannelGroup(c context.Context, arg *v1.ChannelGroupReq) (res *v1.ChannelGroupReply, err error) {
	res = &v1.ChannelGroupReply{
		Groups: make([]*v1.ChannelGroup, 0, model.ChannelMaxGroups),
	}
	if arg.Tid <= 0 {
		err = ecode.RequestErr
		return
	}
	var (
		groups []*model.ChannelGroup
		tids   = make([]int64, 0, model.ChannelMaxGroups)
	)
	if groups, err = s.svr.ChannelGroup(c, arg.Tid); err != nil {
		return
	}
	for _, group := range groups {
		if group.Tid > 0 {
			tids = append(tids, group.Tid)
		}
	}
	if len(tids) == 0 {
		return
	}
	tagMap, err := s.svr.InfoMap(c, 0, tids)
	if err != nil {
		return
	}
	for _, v := range groups {
		tag, ok := tagMap[v.Tid]
		if !ok || tag == nil || tag.State != model.TagStateNormal {
			continue
		}
		cg := &v1.ChannelGroup{
			Id:    v.Id,
			Ptid:  v.Ptid,
			Tid:   v.Tid,
			Tname: tag.Name,
			Alias: v.Alias,
			Rank:  v.Rank,
			CTime: v.CTime,
			MTime: v.MTime,
		}
		res.Groups = append(res.Groups, cg)
	}
	return
}
