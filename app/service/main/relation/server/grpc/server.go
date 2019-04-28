package grpc

import (
	"context"

	pb "go-common/app/service/main/relation/api"
	"go-common/app/service/main/relation/conf"
	"go-common/app/service/main/relation/model"
	"go-common/app/service/main/relation/service"
	"go-common/library/net/rpc/warden"
)

// New warden rpc server
func New(c *conf.Config, s *service.Service) *warden.Server {
	svr := warden.NewServer(c.WardenServer)
	pb.RegisterRelationServer(svr.Server(), &server{as: s})
	svr, err := svr.Start()
	if err != nil {
		panic(err)
	}
	return svr
}

type server struct {
	as *service.Service
}

var _ pb.RelationServer = &server{}

func (s *server) Relation(ctx context.Context, req *pb.RelationReq) (*pb.FollowingReply, error) {
	following, err := s.as.Relation(ctx, req.Mid, req.Fid)
	if err != nil {
		return nil, err
	}
	followingReply := new(pb.FollowingReply)
	followingReply.DeepCopyFromFollowing(following)
	return followingReply, nil
}

func (s *server) Relations(ctx context.Context, req *pb.RelationsReq) (*pb.FollowingMapReply, error) {
	followsing, err := s.as.Relations(ctx, req.Mid, req.Fid)
	if err != nil {
		return nil, err
	}
	followingMap := map[int64]*pb.FollowingReply{}
	for key, value := range followsing {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingMap[key] = followingReply
	}
	return &pb.FollowingMapReply{FollowingMap: followingMap}, nil
}

func (s *server) Stat(ctx context.Context, req *pb.MidReq) (*pb.StatReply, error) {
	stat, err := s.as.Stat(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	statReply := new(pb.StatReply)
	statReply.DeepCopyFromStat(stat)
	return statReply, nil
}

func (s *server) Stats(ctx context.Context, req *pb.MidsReq) (*pb.StatsReply, error) {
	stat, err := s.as.Stats(ctx, req.Mids)
	if err != nil {
		return nil, err
	}
	statMap := map[int64]*pb.StatReply{}
	for key, value := range stat {
		statReply := new(pb.StatReply)
		statReply.DeepCopyFromStat(value)
		statMap[key] = statReply
	}
	return &pb.StatsReply{StatReplyMap: statMap}, nil
}

func (s *server) Attentions(ctx context.Context, req *pb.MidReq) (*pb.FollowingsReply, error) {
	followings, err := s.as.Attentions(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	followingList := make([]*pb.FollowingReply, len(followings))
	for index, value := range followings {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingList[index] = followingReply
	}
	return &pb.FollowingsReply{FollowingList: followingList}, nil
}

func (s *server) Followings(ctx context.Context, req *pb.MidReq) (*pb.FollowingsReply, error) {
	followings, err := s.as.Followings(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	followingList := make([]*pb.FollowingReply, len(followings))
	for index, value := range followings {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingList[index] = followingReply
	}
	return &pb.FollowingsReply{FollowingList: followingList}, nil
}

func (s *server) AddFollowing(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.AddFollowing(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) DelFollowing(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.DelFollowing(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) Whispers(ctx context.Context, req *pb.MidReq) (*pb.FollowingsReply, error) {
	followings, err := s.as.Whispers(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	followingList := make([]*pb.FollowingReply, len(followings))
	for index, value := range followings {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingList[index] = followingReply
	}
	return &pb.FollowingsReply{FollowingList: followingList}, nil
}

func (s *server) AddWhisper(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.AddWhisper(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) DelWhisper(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.DelWhisper(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}

	return &pb.EmptyReply{}, nil
}

func (s *server) Blacks(ctx context.Context, req *pb.MidReq) (*pb.FollowingsReply, error) {
	followings, err := s.as.Blacks(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	followingList := make([]*pb.FollowingReply, len(followings))
	for index, value := range followings {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingList[index] = followingReply
	}
	return &pb.FollowingsReply{FollowingList: followingList}, nil
}

func (s *server) AddBlack(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.AddBlack(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) DelBlack(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.DelBlack(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) Followers(ctx context.Context, req *pb.MidReq) (*pb.FollowingsReply, error) {
	followings, err := s.as.Followers(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	followingList := make([]*pb.FollowingReply, len(followings))
	for index, value := range followings {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingList[index] = followingReply
	}
	return &pb.FollowingsReply{FollowingList: followingList}, nil
}

func (s *server) DelFollower(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.DelFollower(ctx, req.Mid, req.Fid, req.Source, req.Infoc); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) Tag(ctx context.Context, req *pb.TagIdReq) (*pb.TagReply, error) {
	mids, err := s.as.Tag(ctx, req.Mid, req.TagId, req.RealIp)
	if err != nil {
		return nil, err
	}
	return &pb.TagReply{Mids: mids}, nil
}

func (s *server) Tags(ctx context.Context, req *pb.MidReq) (*pb.TagsCountReply, error) {
	tagCount, err := s.as.Tags(ctx, req.Mid, req.RealIp)
	if err != nil {
		return nil, err
	}
	tagCountList := make([]*pb.TagCountReply, len(tagCount))
	for index, value := range tagCount {
		tagCountReply := new(pb.TagCountReply)
		tagCountReply.DeepCopyFromTagCount(value)
		tagCountList[index] = tagCountReply
	}
	return &pb.TagsCountReply{TagCountList: tagCountList}, nil
}

func (s *server) UserTag(ctx context.Context, req *pb.RelationReq) (*pb.UserTagReply, error) {
	res, err := s.as.UserTag(ctx, req.Mid, req.Fid, req.RealIp)
	if err != nil {
		return nil, err
	}
	return &pb.UserTagReply{Tags: res}, nil
}

func (s *server) CreateTag(ctx context.Context, req *pb.TagReq) (*pb.CreateTagReply, error) {
	res, err := s.as.CreateTag(ctx, req.Mid, req.Tag, req.RealIp)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTagReply{TagId: res}, nil
}

func (s *server) UpdateTag(ctx context.Context, req *pb.TagUpdateReq) (*pb.EmptyReply, error) {
	if err := s.as.UpdateTag(ctx, req.Mid, req.TagId, req.New, req.RealIp); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) DelTag(ctx context.Context, req *pb.TagDelReq) (*pb.EmptyReply, error) {
	if err := s.as.DelTag(ctx, req.Mid, req.TagId, req.RealIp); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) TagsAddUsers(ctx context.Context, req *pb.TagsMoveUsersReq) (*pb.EmptyReply, error) {
	if err := s.as.TagsAddUsers(ctx, req.Mid, req.AfterTagIds, req.Fids, req.RealIp); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) TagsCopyUsers(ctx context.Context, req *pb.TagsMoveUsersReq) (*pb.EmptyReply, error) {
	if err := s.as.TagsMoveUsers(ctx, req.Mid, req.BeforeId, req.AfterTagIds, req.Fids, req.RealIp); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) TagsMoveUsers(ctx context.Context, req *pb.TagsMoveUsersReq) (*pb.EmptyReply, error) {
	if err := s.as.TagsMoveUsers(ctx, req.Mid, req.BeforeId, req.AfterTagIds, req.Fids, req.RealIp); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) Prompt(ctx context.Context, req *pb.PromptReq) (*pb.PromptReply, error) {
	argPrompt := &model.ArgPrompt{Mid: req.Mid, Fid: req.Fid, Btype: req.Btype}
	success, err := s.as.Prompt(ctx, argPrompt)
	if err != nil {
		return nil, err
	}
	return &pb.PromptReply{Success: success}, nil
}

func (s *server) ClosePrompt(ctx context.Context, req *pb.PromptReq) (*pb.EmptyReply, error) {
	argPrompt := &model.ArgPrompt{Mid: req.Mid, Fid: req.Fid, Btype: req.Btype}
	if err := s.as.ClosePrompt(ctx, argPrompt); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) AddSpecial(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.AddSpecial(ctx, req.Mid, req.Fid); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) DelSpecial(ctx context.Context, req *pb.FollowingReq) (*pb.EmptyReply, error) {
	if err := s.as.DelSpecial(ctx, req.Mid, req.Fid); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) Special(ctx context.Context, req *pb.MidReq) (*pb.SpecialReply, error) {
	mids, err := s.as.Special(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.SpecialReply{Mids: mids}, nil
}

func (s *server) FollowersUnread(ctx context.Context, req *pb.MidReq) (*pb.FollowersUnreadReply, error) {
	hasUnread, err := s.as.Unread(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.FollowersUnreadReply{HasUnread: hasUnread}, nil
}

func (s *server) FollowersUnreadCount(ctx context.Context, req *pb.MidReq) (*pb.FollowersUnreadCountReply, error) {
	count, err := s.as.UnreadCount(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.FollowersUnreadCountReply{UnreadCount: count}, nil
}

func (s *server) AchieveGet(ctx context.Context, req *pb.AchieveGetReq) (*pb.AchieveGetReply, error) {
	argAchieveGet := &model.ArgAchieveGet{Mid: req.Mid, Award: req.Award}
	achieveGetReply, err := s.as.AchieveGet(ctx, argAchieveGet)
	if err != nil {
		return nil, err
	}
	return &pb.AchieveGetReply{AwardToken: achieveGetReply.AwardToken}, nil
}

func (s *server) Achieve(ctx context.Context, req *pb.AchieveReq) (*pb.AchieveReply, error) {
	argAchieve := &model.ArgAchieve{AwardToken: req.AwardToken}
	achieveReply, err := s.as.Achieve(ctx, argAchieve)
	if err != nil {
		return nil, err
	}
	return &pb.AchieveReply{Award: achieveReply.Award, Mid: achieveReply.Mid}, nil
}

func (s *server) ResetFollowersUnread(ctx context.Context, req *pb.MidReq) (*pb.EmptyReply, error) {
	if err := s.as.ResetUnread(ctx, req.Mid); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) ResetFollowersUnreadCount(ctx context.Context, req *pb.MidReq) (*pb.EmptyReply, error) {
	if err := s.as.ResetUnreadCount(ctx, req.Mid); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) DisableFollowerNotify(ctx context.Context, req *pb.MidReq) (*pb.EmptyReply, error) {
	if err := s.as.DisableFollowerNotify(ctx, &model.ArgMid{Mid: req.Mid}); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) EnableFollowerNotify(ctx context.Context, req *pb.MidReq) (*pb.EmptyReply, error) {
	if err := s.as.EnableFollowerNotify(ctx, &model.ArgMid{Mid: req.Mid}); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}

func (s *server) FollowerNotifySetting(ctx context.Context, req *pb.MidReq) (*pb.FollowerNotifySettingReply, error) {
	followerNotifySetting, err := s.as.FollowerNotifySetting(ctx, &model.ArgMid{Mid: req.Mid})
	if err != nil {
		return nil, err
	}
	return &pb.FollowerNotifySettingReply{Mid: followerNotifySetting.Mid, Enabled: followerNotifySetting.Enabled}, nil
}

func (s *server) SameFollowings(ctx context.Context, req *pb.SameFollowingReq) (*pb.FollowingsReply, error) {
	followings, err := s.as.SameFollowings(ctx, &model.ArgSameFollowing{Mid1: req.Mid, Mid2: req.Mid2})
	if err != nil {
		return nil, err
	}
	followingList := make([]*pb.FollowingReply, len(followings))
	for index, value := range followings {
		followingReply := new(pb.FollowingReply)
		followingReply.DeepCopyFromFollowing(value)
		followingList[index] = followingReply
	}
	return &pb.FollowingsReply{FollowingList: followingList}, nil
}
