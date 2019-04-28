package dao

import (
	"context"

	"go-common/app/interface/main/tag/model"
	taGrpcModel "go-common/app/service/main/tag/api"
	"go-common/library/log"
)

// ChannelCategories get channel categories.
func (d *Dao) ChannelCategories(c context.Context, lastID int64, pageSize, state int32) (res []*taGrpcModel.ChannelCategory, err error) {
	var (
		arg = &taGrpcModel.ChannelCategoriesReq{
			StartId: lastID,
			Num:     pageSize,
			State:   state,
		}
		reply *taGrpcModel.ChannelCategoriesReply
	)
	if reply, err = d.tagRPC.ChannelCategories(c, arg); err != nil {
		log.Error("d.dao.ChannelCategories(%d,%d,%d) error: %v", lastID, pageSize, state, err)
		return
	}
	return reply.Categories, nil
}

// Channels get channels.
func (d *Dao) Channels(c context.Context, lastID int64, pageSize int32) (res []*taGrpcModel.Channel, err error) {
	var (
		arg = &taGrpcModel.ChannelsReq{
			StartId: lastID,
			Num:     pageSize,
		}
		reply *taGrpcModel.ChannelsReply
	)
	if reply, err = d.tagRPC.Channels(c, arg); err != nil {
		log.Error("d.dao.Channels(%d,%d) error: %v", lastID, pageSize, err)
		return
	}
	return reply.Channels, nil
}

// ChannelRules get channel rules.
func (d *Dao) ChannelRules(c context.Context, lastID int64, pageSize, state int32) (res []*taGrpcModel.ChannelRule, err error) {
	var (
		arg = &taGrpcModel.ChannelRulesReq{
			StartId: lastID,
			Num:     pageSize,
			State:   state,
		}
		reply *taGrpcModel.ChannelRulesReply
	)
	if reply, err = d.tagRPC.ChannelRules(c, arg); err != nil {
		log.Error("d.dao.ChannelRules(%d,%d,%d) error: %v", lastID, pageSize, state, err)
		return
	}
	return reply.Rules, nil
}

// ChannelGroup get channel groups.
func (d *Dao) ChannelGroup(c context.Context, tid int64) (res []*model.ChannelSynonym, err error) {
	var (
		arg = &taGrpcModel.ChannelGroupReq{
			Tid: tid,
		}
		reply *taGrpcModel.ChannelGroupReply
	)
	res = make([]*model.ChannelSynonym, 0)
	if reply, err = d.tagRPC.ChannelGroup(c, arg); err != nil {
		log.Error("d.dao.ChannelGroup(%d) error: %v", tid, err)
		return
	}
	for _, v := range reply.Groups {
		cg := &model.ChannelSynonym{
			Id:    v.Tid,
			Name:  v.Tname,
			Alias: v.Alias,
			Rank:  v.Rank,
			CTime: v.CTime,
			MTime: v.MTime,
		}
		res = append(res, cg)
	}
	return
}
