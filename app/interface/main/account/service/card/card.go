package card

import (
	"context"

	"go-common/app/interface/main/account/conf"
	v1 "go-common/app/service/main/card/api/grpc/v1"
)

// Service .
type Service struct {
	// conf
	c *conf.Config
	// card service
	cardRPC v1.CardClient
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
	}
	cardRPC, err := v1.NewClient(c.CardClient)
	if err != nil {
		panic(err)
	}
	s.cardRPC = cardRPC
	return
}

// UserCard user card info.
func (s *Service) UserCard(c context.Context, mid int64) (res *v1.ModelUserCard, err error) {
	var reply *v1.UserCardReply
	if reply, err = s.cardRPC.UserCard(c, &v1.UserCardReq{Mid: mid}); err != nil {
		return
	}
	res = reply.Res
	return
}

// Card get card info by id.
func (s *Service) Card(c context.Context, id int64) (res *v1.ModelCard, err error) {
	var reply *v1.CardReply
	if reply, err = s.cardRPC.Card(c, &v1.CardReq{Id: id}); err != nil {
		return
	}
	res = reply.Data_0
	return
}

// CardHots get all hots cards.
func (s *Service) CardHots(c context.Context) (res []*v1.ModelCard, err error) {
	var reply *v1.CardHotsReply
	if reply, err = s.cardRPC.CardHots(c, &v1.CardHotsReq{}); err != nil {
		return
	}
	res = reply.Data_0
	return
}

// AllGroup all group.
func (s *Service) AllGroup(c context.Context, mid int64) (res *v1.ModelAllGroupResp, err error) {
	var reply *v1.AllGroupReply
	if reply, err = s.cardRPC.AllGroup(c, &v1.AllGroupReq{Mid: mid}); err != nil {
		return
	}
	res = reply.Res
	return
}

// CardsByGid get cards by gid.
func (s *Service) CardsByGid(c context.Context, id int64) (res []*v1.ModelCard, err error) {
	var reply *v1.CardsByGidReply
	if reply, err = s.cardRPC.CardsByGid(c, &v1.CardsByGidReq{Gid: id}); err != nil {
		return
	}
	res = reply.Data_0
	return
}

// Equip card equip.
func (s *Service) Equip(c context.Context, arg *v1.ModelArgEquip) (err error) {
	_, err = s.cardRPC.Equip(c, &v1.EquipReq{Arg: arg})
	return
}

// Demount card demount.
func (s *Service) Demount(c context.Context, mid int64) (err error) {
	_, err = s.cardRPC.DemountEquip(c, &v1.DemountEquipReq{Mid: mid})
	return
}
