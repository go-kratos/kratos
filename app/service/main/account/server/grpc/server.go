package grpc

import (
	"context"

	pb "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/conf"
	"go-common/app/service/main/account/service"
	"go-common/library/net/rpc/warden"
)

// New warden rpc server
func New(c *conf.Config, s *service.Service) (svr *warden.Server) {
	svr = warden.NewServer(c.WardenServer)
	pb.RegisterAccountServer(svr.Server(), &server{as: s})
	return svr
}

// Start create and start warden rpc server
func Start(c *conf.Config, s *service.Service) (svr *warden.Server, err error) {
	svr = warden.NewServer(c.WardenServer)
	pb.RegisterAccountServer(svr.Server(), &server{as: s})
	if svr, err = svr.Start(); err != nil {
		return
	}
	return
}

type server struct {
	as *service.Service
}

var _ pb.AccountServer = &server{}

func (s *server) Info3(ctx context.Context, req *pb.MidReq) (*pb.InfoReply, error) {
	info, err := s.as.Info(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.InfoReply{Info: info}, nil
}

func (s *server) Infos3(ctx context.Context, req *pb.MidsReq) (*pb.InfosReply, error) {
	infos, err := s.as.Infos(ctx, req.Mids)
	if err != nil {
		return nil, err
	}
	return &pb.InfosReply{Infos: infos}, nil
}

func (s *server) InfosByName3(ctx context.Context, req *pb.NamesReq) (*pb.InfosReply, error) {
	infos, err := s.as.InfosByName(ctx, req.Names)
	if err != nil {
		return nil, err
	}
	return &pb.InfosReply{Infos: infos}, nil
}

func (s *server) Card3(ctx context.Context, req *pb.MidReq) (*pb.CardReply, error) {
	card, err := s.as.Card(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.CardReply{Card: card}, nil
}

func (s *server) Cards3(ctx context.Context, req *pb.MidsReq) (*pb.CardsReply, error) {
	cards, err := s.as.Cards(ctx, req.Mids)
	if err != nil {
		return nil, err
	}
	return &pb.CardsReply{Cards: cards}, nil
}

func (s *server) Profile3(ctx context.Context, req *pb.MidReq) (*pb.ProfileReply, error) {
	profile, err := s.as.Profile(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.ProfileReply{Profile: profile}, nil
}

func (s *server) ProfileWithStat3(ctx context.Context, req *pb.MidReq) (*pb.ProfileStatReply, error) {
	profileStat, err := s.as.ProfileWithStat(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	level := pb.LevelInfo{}
	level.DeepCopyFromLevelInfo(&profileStat.LevelExp)
	return &pb.ProfileStatReply{
		Profile:   profileStat.Profile,
		LevelInfo: level,
		Coins:     profileStat.Coins,
		Follower:  profileStat.Follower,
		Following: profileStat.Following,
	}, nil
}

func (s *server) AddExp3(ctx context.Context, req *pb.ExpReq) (*pb.ExpReply, error) {
	return &pb.ExpReply{}, s.as.AddExp(ctx, req.Mid, req.Exp, req.Operater, req.Operate, req.Reason)
}

func (s *server) AddMoral3(ctx context.Context, req *pb.MoralReq) (*pb.MoralReply, error) {
	return &pb.MoralReply{}, s.as.AddMoral(ctx, req.Mid, req.Moral, req.Oper, req.Reason, req.Remark)
}

func (s *server) Relation3(ctx context.Context, req *pb.RelationReq) (*pb.RelationReply, error) {
	relation, err := s.as.Relation(ctx, req.Mid, req.Owner)
	if err != nil {
		return nil, err
	}
	return &pb.RelationReply{Following: relation.Following}, nil
}

func (s *server) Attentions3(ctx context.Context, req *pb.MidReq) (*pb.AttentionsReply, error) {
	attentions, err := s.as.Attentions(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	return &pb.AttentionsReply{Attentions: attentions}, nil
}

func (s *server) Blacks3(ctx context.Context, req *pb.MidReq) (*pb.BlacksReply, error) {
	blackList, err := s.as.Blacks(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	blackListBool := make(map[int64]bool, len(blackList))
	for k := range blackList {
		blackListBool[k] = true
	}
	return &pb.BlacksReply{BlackList: blackListBool}, nil
}

func (s *server) Relations3(ctx context.Context, req *pb.RelationsReq) (*pb.RelationsReply, error) {
	relations, err := s.as.Relations(ctx, req.Mid, req.Owners)
	if err != nil {
		return nil, err
	}
	newRelations := make(map[int64]*pb.RelationReply, len(relations))
	for k, v := range relations {
		newRelations[k] = &pb.RelationReply{Following: v.Following}
	}
	return &pb.RelationsReply{Relations: newRelations}, nil
}

func (s *server) RichRelations3(ctx context.Context, req *pb.RichRelationReq) (*pb.RichRelationsReply, error) {
	richRelations, err := s.as.RichRelations2(ctx, req.Owner, req.Mids)
	if err != nil {
		return nil, err
	}
	newRichRelations := make(map[int64]int32, len(richRelations))
	for k, v := range richRelations {
		newRichRelations[k] = int32(v)
	}
	return &pb.RichRelationsReply{RichRelations: newRichRelations}, nil
}

func (s *server) Vip3(ctx context.Context, req *pb.MidReq) (*pb.VipReply, error) {
	vip, err := s.as.Vip(ctx, req.Mid)
	if err != nil {
		return nil, err
	}
	reply := new(pb.VipReply)
	reply.DeepCopyFromVipInfo(vip)
	return reply, nil
}

func (s *server) Vips3(ctx context.Context, req *pb.MidsReq) (*pb.VipsReply, error) {
	vips, err := s.as.Vips(ctx, req.Mids)
	if err != nil {
		return nil, err
	}
	pvips := make(map[int64]*pb.VipReply, len(vips))
	for mid, vip := range vips {
		pvip := new(pb.VipReply)
		pvip.DeepCopyFromVipInfo(vip)
		pvips[mid] = pvip
	}
	return &pb.VipsReply{Vips: pvips}, nil
}
