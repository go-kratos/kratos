package service

import (
	"context"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// GenMachinesV2 create multiple machines for v2.
func (s *Service) GenMachinesV2(c context.Context, gmr *model.GenMachinesRequest, u string) (err error) {
	var (
		ins        []*model.CreateInstance
		cluster    *model.Cluster
		hasMachine bool
		ms         []*model.Machine
	)
	if hasMachine, err = s.dao.HasMachine(gmr.Name); err != nil {
		return
	}
	if hasMachine {
		err = ecode.MerlinDuplicateMachineNameErr
		return
	}
	if cluster, err = s.dao.QueryCluster(c, gmr.NetworkID); err != nil {
		return
	}
	if cluster == nil {
		err = ecode.MerlinInvalidClusterErr
		return
	}
	if err = cluster.Verify(); err != nil {
		return
	}
	pgmr := gmr.ToPaasGenMachineRequest(s.c.Paas.MachineLimitRatio)
	if ms, err = s.dao.InsertMachinesV2(u, gmr, pgmr); err != nil {
		return
	}
	go func() {
		var (
			mIDs  []int64
			mim   map[int64]string
			nmm   = make(map[string]*model.Machine)
			cTODO = context.TODO()
		)
		for _, m := range ms {
			mIDs = append(mIDs, m.ID)
			nmm[m.Name] = m
		}
		if mim, err = s.AddTagToMachine(cTODO, u, gmr.Image, mIDs); err != nil {
			if err = s.dao.UpdateStatusForMachines(model.CreateTagFailedMachineInMerlin, mIDs); err != nil {
				log.Error("Update the status of machines(%v) err(%v)", mIDs, err)
			}
			return
		}
		for i, pm := range pgmr.Machines {
			pgmr.Machines[i].Image = mim[nmm[pm.Name].ID]
			pgmr.Machines[i].Snapshot = true
			pgmr.Machines[i].ForcePullImage = true
		}
		if ins, err = s.dao.GenPaasMachines(cTODO, pgmr); err != nil {
			return
		}
		for _, in := range ins {
			if in.InstanceCreateStatus == model.CreateFailedMachineInPaas {
				s.dao.UpdateMachineStatusByName(model.ImmediatelyFailedMachineInMerlin, in.InstanceName)
			}
		}
	}()
	return
}
