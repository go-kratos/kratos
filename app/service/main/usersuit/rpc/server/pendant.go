package server

import (
	"go-common/app/service/main/usersuit/model"
	"go-common/library/net/rpc/context"
)

// Equipment obtain Equipment by mid .
func (r *RPC) Equipment(c context.Context, mid int64, res *model.PendantEquip) (err error) {
	var pe *model.PendantEquip
	if pe, err = r.s.Equipment(c, mid); err == nil && pe != nil {
		*res = *pe
	}
	return
}

// Equipments obtain equipments by mids .
func (r *RPC) Equipments(c context.Context, mids []int64, res *map[int64]*model.PendantEquip) (err error) {
	*res, err = r.s.Equipments(c, mids)
	return
}
