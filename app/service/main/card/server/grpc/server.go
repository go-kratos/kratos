// Package server generate by warden_gen
package server

import (
	"context"

	"go-common/app/service/main/card/api/grpc/v1"
	"go-common/app/service/main/card/model"
	service "go-common/app/service/main/card/service"
	"go-common/library/net/rpc/warden"
)

// New Card warden rpc server
func New(c *warden.ServerConfig, svr *service.Service) *warden.Server {
	ws := warden.NewServer(c)
	v1.RegisterCardServer(ws.Server(), &server{svr})
	ws, err := ws.Start()
	if err != nil {
		panic(err)
	}
	return ws
}

type server struct {
	svr *service.Service
}

var _ v1.CardServer = &server{}

// Card get card info.
func (s *server) Card(c context.Context, req *v1.CardReq) (*v1.CardReply, error) {
	var res *v1.ModelCard
	if cd := s.svr.Card(c, req.Id); cd != nil {
		res = convertModelCard(cd)
	}
	return &v1.CardReply{Data_0: res}, nil
}

// CardHots get card hots.
func (s *server) CardHots(c context.Context, req *v1.CardHotsReq) (*v1.CardHotsReply, error) {
	cs := s.svr.CardHots(c)
	ls := make([]*v1.ModelCard, len(cs))
	for i, v := range cs {
		ls[i] = convertModelCard(v)
	}
	return &v1.CardHotsReply{Data_0: ls}, nil
}

// CardsByGid get card by gid.
func (s *server) CardsByGid(c context.Context, req *v1.CardsByGidReq) (*v1.CardsByGidReply, error) {
	cs := s.svr.CardsByGid(c, req.Gid)
	ls := make([]*v1.ModelCard, len(cs))
	for i, v := range cs {
		ls[i] = convertModelCard(v)
	}
	return &v1.CardsByGidReply{Data_0: ls}, nil
}

// Equip user equip card.
func (s *server) Equip(c context.Context, req *v1.EquipReq) (*v1.EquipReply, error) {
	return nil, s.svr.Equip(c, &model.ArgEquip{
		Mid:    req.Arg.Mid,
		CardID: req.Arg.CardId,
	})
}

// DemountEquip delete equip.
func (s *server) DemountEquip(c context.Context, req *v1.DemountEquipReq) (*v1.DemountEquipReply, error) {
	return nil, s.svr.DemountEquip(c, req.Mid)
}

// AllGroup all group.
func (s *server) AllGroup(c context.Context, req *v1.AllGroupReq) (reply *v1.AllGroupReply, err error) {
	var gs *model.AllGroupResp
	if gs, err = s.svr.AllGroup(c, req.Mid); err != nil {
		return
	}
	if gs == nil {
		return
	}
	rs := new(v1.ModelAllGroupResp)
	if gs.UserCard != nil {
		rs.UserCard = convertModelUserCard(gs.UserCard)
	}
	ls := make([]*v1.ModelGroupInfo, len(gs.List))
	for i, v := range gs.List {
		cs := make([]*v1.ModelCard, len(v.Cards))
		for ci, cv := range v.Cards {
			cs[ci] = convertModelCard(cv)
		}
		ls[i] = &v1.ModelGroupInfo{
			GroupId:   v.GroupID,
			GroupName: v.GroupName,
			Cards:     cs,
		}
	}
	rs.List = ls
	return &v1.AllGroupReply{
		Res: rs,
	}, nil
}

// UserCard get user card info.
func (s *server) UserCard(c context.Context, req *v1.UserCardReq) (res *v1.UserCardReply, err error) {
	var cd *model.UserCard
	res = new(v1.UserCardReply)
	if cd, err = s.svr.UserCard(c, req.Mid); err != nil {
		return
	}
	if cd == nil {
		return
	}
	res.Res = convertModelUserCard(cd)
	return
}

// UserCards get user card infos.
func (s *server) UserCards(c context.Context, req *v1.UserCardsReq) (res *v1.UserCardsReply, err error) {
	var cs map[int64]*model.UserCard
	res = new(v1.UserCardsReply)
	if cs, err = s.svr.UserCards(c, req.Mids); err != nil {
		return
	}
	if len(cs) <= 0 {
		return
	}
	ls := make(map[int64]*v1.ModelUserCard, len(cs))
	for k, v := range cs {
		ls[k] = convertModelUserCard(v)
	}
	res.Res = ls
	return
}

func convertModelUserCard(v *model.UserCard) *v1.ModelUserCard {
	return &v1.ModelUserCard{
		Mid:          v.Mid,
		Id:           v.ID,
		CardUrl:      v.CardURL,
		CardType:     v.CardType,
		Name:         v.Name,
		ExpireTime:   v.ExpireTime,
		CardTypeName: v.CardTypeName,
		BigCardUrl:   v.BigCradURL,
	}
}

func convertModelCard(cv *model.Card) *v1.ModelCard {
	return &v1.ModelCard{
		Id:           cv.ID,
		Name:         cv.Name,
		State:        cv.State,
		Deleted:      cv.Deleted,
		IsHot:        cv.IsHot,
		CardUrl:      cv.CardURL,
		BigCardUrl:   cv.BigCradURL,
		CardType:     cv.CardType,
		CardTypeName: cv.CardTypeName,
	}
}
