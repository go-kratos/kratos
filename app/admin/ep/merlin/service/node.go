package service

import (
	"context"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
)

// UpdateMachineNode update the node of machine
func (s *Service) UpdateMachineNode(c context.Context, umnr *model.UpdateMachineNodeRequest) (err error) {
	var (
		m   *model.Machine
		res string
	)
	if m, err = s.dao.QueryMachine(umnr.MachineID); err != nil {
		return
	}
	if res, err = s.dao.UpdatePaasMachineNode(c, model.NewPaasUpdateMachineNodeRequest(m.ToPaasQueryAndDelMachineRequest(), umnr.Nodes)); err != nil {
		return
	}
	if res != model.Success {
		err = ecode.MerlinUpdateNodeErr
		return
	}
	if err = s.dao.DelMachineNodeByMachineID(umnr.MachineID); err != nil {
		return
	}
	if err = s.dao.GenMachineNodes(umnr.ToMachineNodes()); err != nil {
		return
	}
	return
}

// QueryMachineNodes query nodes by mID
func (s *Service) QueryMachineNodes(mID int64) ([]*model.MachineNode, error) {
	return s.dao.QueryMachineNodes(mID)
}
