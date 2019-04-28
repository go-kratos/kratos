package grpc

import (
	"context"

	pb "go-common/app/interface/main/tag/api"
	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
)

// Channel get a channel info, include tags, and channel synonyms.
func (s *grpcServer) Channel(c context.Context, arg *pb.ChannelReq) (res *pb.ChannelReply, err error) {
	res = &pb.ChannelReply{
		Channel:  new(pb.Channel),
		Synonyms: make([]*pb.ChannelSynonym, 0),
	}
	req := &model.ReqChannelDetail{
		Tid:   arg.Tid,
		TName: arg.Tname,
		Mid:   arg.Mid,
		From:  arg.From,
	}
	channelDetail, err := s.svr.ChannelDetail(c, req)
	if err != nil || channelDetail == nil || channelDetail.Tag == nil {
		return
	}
	res.Channel = &pb.Channel{
		Id:           channelDetail.Tag.ID,
		Name:         channelDetail.Tag.Name,
		Type:         channelDetail.Tag.Type,
		Cover:        channelDetail.Tag.Cover,
		HeadCover:    channelDetail.Tag.HeadCover,
		Content:      channelDetail.Tag.Content,
		ShortContent: channelDetail.Tag.ShortContent,
		Bind:         channelDetail.Tag.Bind,
		Sub:          channelDetail.Tag.Sub,
		Attention:    channelDetail.Tag.Attention,
		Activity:     channelDetail.Tag.Activity,
		IntShield:    channelDetail.Tag.INTShield,
		State:        channelDetail.Tag.State,
		Ctime:        channelDetail.Tag.CTime,
		Mtime:        channelDetail.Tag.MTime,
	}
	for _, synonym := range channelDetail.Synonym {
		res.Synonyms = append(res.Synonyms, &pb.ChannelSynonym{
			Id:    synonym.Id,
			Name:  synonym.Name,
			Alias: synonym.Alias,
			Ctime: synonym.CTime,
			Mtime: synonym.MTime,
		})
	}
	return
}

// ChannelCategory get channels categories by from that international version or normal app version.
func (s *grpcServer) ChannelCategory(c context.Context, arg *pb.ChannelCategoryReq) (res *pb.ChannelCategoryReply, err error) {
	res = &pb.ChannelCategoryReply{
		Categories: make([]*pb.ChannelCategory, 0),
	}
	categories, err := s.svr.ChannelCategories(c, &model.ArgChannelCategories{From: arg.From})
	if err != nil {
		return
	}
	for _, category := range categories {
		res.Categories = append(res.Categories, &pb.ChannelCategory{
			Id:        category.ID,
			Name:      category.Name,
			State:     category.State,
			IntShield: category.INTShield,
			Ctime:     category.CTime,
			Mtime:     category.MTime,
		})
	}
	return
}

// ChanneList get channels by channel category id and the app version.
func (s *grpcServer) ChanneList(c context.Context, arg *pb.ChanneListReq) (res *pb.ChannelsReply, err error) {
	res = &pb.ChannelsReply{
		Channels: make([]*pb.Channel, 0),
	}
	channels, err := s.svr.ChanneList(c, arg.Mid, arg.Id, arg.From)
	if err != nil {
		return
	}
	for _, channel := range channels {
		res.Channels = append(res.Channels, &pb.Channel{
			Id:           channel.ID,
			Name:         channel.Name,
			TypeId:       channel.Type,
			Cover:        channel.Cover,
			HeadCover:    channel.HeadCover,
			Content:      channel.Content,
			ShortContent: channel.ShortContent,
			Bind:         channel.Bind,
			Sub:          channel.Sub,
			Attention:    channel.Attention,
			Activity:     channel.AttrVal(model.ChannelAttrActivity),
			IntShield:    channel.AttrVal(model.ChannelAttrINT),
			State:        channel.State,
			Ctime:        channel.CTime,
			Mtime:        channel.MTime,
		})
	}
	return
}

// ChannelRecommend get a recommend channel list by mid.
func (s *grpcServer) ChannelRecommend(c context.Context, arg *pb.ChannelRecommendReq) (res *pb.ChannelsReply, err error) {
	res = &pb.ChannelsReply{
		Channels: make([]*pb.Channel, 0),
	}
	channels, err := s.svr.RecommandChannel(c, arg.Mid, arg.From)
	if err != nil {
		return
	}
	for _, channel := range channels {
		res.Channels = append(res.Channels, &pb.Channel{
			Id:           channel.ID,
			Name:         channel.Name,
			TypeId:       channel.Type,
			Cover:        channel.Cover,
			HeadCover:    channel.HeadCover,
			Content:      channel.Content,
			ShortContent: channel.ShortContent,
			Bind:         channel.Bind,
			Sub:          channel.Sub,
			Attention:    channel.Attention,
			Activity:     channel.AttrVal(model.ChannelAttrActivity),
			IntShield:    channel.AttrVal(model.ChannelAttrINT),
			State:        channel.State,
			Ctime:        channel.CTime,
			Mtime:        channel.MTime,
		})
	}
	return
}

// ChannelDiscovery get a channel list by mid and channel state (3>2>1).
func (s *grpcServer) ChannelDiscovery(c context.Context, arg *pb.ChannelDiscoveryReq) (res *pb.ChannelsReply, err error) {
	res = &pb.ChannelsReply{
		Channels: make([]*pb.Channel, 0, model.DiscoveryChannelNum),
	}
	channels, err := s.svr.DiscoveryChannel(c, arg.Mid, arg.From)
	if err != nil {
		return
	}
	for _, channel := range channels {
		res.Channels = append(res.Channels, &pb.Channel{
			Id:           channel.ID,
			Name:         channel.Name,
			TypeId:       channel.Type,
			Cover:        channel.Cover,
			HeadCover:    channel.HeadCover,
			Content:      channel.Content,
			ShortContent: channel.ShortContent,
			Bind:         channel.Bind,
			Sub:          channel.Sub,
			Attention:    channel.Attention,
			Activity:     channel.AttrVal(model.ChannelAttrActivity),
			IntShield:    channel.AttrVal(model.ChannelAttrINT),
			State:        channel.State,
			Ctime:        channel.CTime,
			Mtime:        channel.MTime,
		})
	}
	return
}

// ChannelSquare get channel infos and archives.
func (s *grpcServer) ChannelSquare(c context.Context, arg *pb.ChannelSquareReq) (res *pb.ChannelSquareReply, err error) {
	res = &pb.ChannelSquareReply{
		Squares: make([]*pb.ChannelSquare, 0, arg.TagNumber),
	}
	if arg.TagNumber <= 0 {
		err = ecode.RequestErr
		return
	}
	req := &model.ReqChannelSquare{
		Mid:        arg.Mid,
		TagNumber:  arg.TagNumber,
		OidNumber:  arg.ResourceNumber,
		Type:       arg.Type,
		Buvid:      arg.Buvid,
		Build:      arg.Build,
		LoginEvent: arg.LoginEvent,
		DisplayID:  arg.DisplayId,
		Plat:       arg.Plat,
		From:       arg.From,
	}
	css, err := s.svr.ChannelSquare(c, req)
	if err != nil {
		return
	}
	for _, cs := range css.Channels {
		channel := &pb.Channel{
			Id:           cs.ID,
			Name:         cs.Name,
			TypeId:       cs.Type,
			Cover:        cs.Cover,
			HeadCover:    cs.HeadCover,
			Content:      cs.Content,
			ShortContent: cs.ShortContent,
			Bind:         cs.Bind,
			Sub:          cs.Sub,
			Attention:    cs.Attention,
			Activity:     cs.AttrVal(model.ChannelAttrActivity),
			IntShield:    cs.AttrVal(model.ChannelAttrINT),
			State:        cs.State,
			Ctime:        cs.CTime,
			Mtime:        cs.MTime,
		}
		k, ok := css.Oids[cs.ID]
		if !ok {
			k = make([]int64, 0)
		}
		res.Squares = append(res.Squares, &pb.ChannelSquare{
			Channel: channel,
			Oids:    k,
		})
	}
	return
}

// ChannelResources resource feed under channel.
func (s *grpcServer) ChannelResources(c context.Context, arg *pb.ChannelResourcesReq) (res *pb.ChannelResourcesReply, err error) {
	res = &pb.ChannelResourcesReply{
		Oids: make([]int64, 0),
	}
	req := &model.ArgChannelResource{
		Tid:        arg.Tid,
		Mid:        arg.Mid,
		Plat:       arg.Plat,
		LoginEvent: arg.LoginEvent,
		RequestCNT: arg.RequestCnt,
		DisplayID:  arg.DisplayId,
		From:       arg.From,
		Type:       arg.Type,
		Build:      arg.Build,
		Name:       arg.Tname,
		Buvid:      arg.Buvid,
	}
	cr, err := s.svr.ChannelResources(c, req)
	if err != nil {
		return
	}
	res.Oids = cr.Oids
	res.Failover = cr.Failover
	res.WhetherChannel = cr.IsChannel
	return
}

// ChannelCheckBack resource channel checkback.
func (s *grpcServer) ChannelCheckBack(c context.Context, arg *pb.ChannelCheckBackReq) (res *pb.ChannelCheckBackReply, err error) {
	res = &pb.ChannelCheckBackReply{
		Checkbacks: make(map[int64]*pb.ChannelCheckBack, len(arg.Oids)),
	}
	if len(arg.Oids) <= 0 {
		err = ecode.RequestErr
		return
	}
	checkInfoMap, err := s.svr.ResChannelCheckBack(c, arg.Oids, arg.Type)
	if err != nil {
		return
	}
	for _, oid := range arg.Oids {
		checkInfo, ok := checkInfoMap[oid]
		if !ok || checkInfo == nil {
			res.Checkbacks[oid] = new(pb.ChannelCheckBack)
			continue
		}
		channelHitMap := make(map[int64]*pb.ChannelHit, len(checkInfo.Channels))
		for _, v := range checkInfo.Channels {
			if v.Tid <= 0 || v.TName == "" {
				continue
			}
			channelHitMap[v.Tid] = &pb.ChannelHit{
				Tid:       v.Tid,
				Tname:     v.TName,
				HitRules:  v.HitRules,
				HitTnames: v.HitTNames,
			}
		}
		res.Checkbacks[oid] = &pb.ChannelCheckBack{
			Hits:      channelHitMap,
			Checkback: checkInfo.CheckBack,
		}
	}
	return
}
