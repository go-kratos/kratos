package client

import (
	"context"
	"go-common/app/service/main/usersuit/model"
)

const (
	_pendantEquip  = "RPC.Equipment"
	_pendantEquips = "RPC.Equipments"
)

// Equipment obtian equipment by mid
func (s *Service2) Equipment(c context.Context, arg *model.ArgEquipment) (res *model.PendantEquip, err error) {
	res = new(model.PendantEquip)
	err = s.client.Call(c, _pendantEquip, arg.Mid, res)
	return
}

// Equipments obtian equipment by mids
func (s *Service2) Equipments(c context.Context, arg *model.ArgEquipments) (res map[int64]*model.PendantEquip, err error) {
	err = s.client.Call(c, _pendantEquips, arg.Mids, &res)
	return
}
