package service

import (
	"context"
	"sort"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// GroupPendantMid .
func (s *Service) GroupPendantMid(c context.Context, arg *model.ArgGPMID) (res []*model.GroupPendantList, err error) {
	var (
		ok   bool
		pids []int64
		pe   = &model.PendantEquip{}
		pps  = []*model.PendantPackage{}
		mpp  = map[int64]*model.PendantPackage{}
	)
	if pids, ok = s.gidMap[arg.GID]; !ok {
		log.Warn("mid(%d) gid(%d) ip(%s)  group pendants is empty", arg.MID, arg.GID, metadata.String(c, metadata.RemoteIP))
		return
	}
	if arg.MID != 0 {
		if pps, err = s.PackageInfo(c, arg.MID); err != nil {
			return
		}
		if pe, err = s.Equipment(c, arg.MID); err != nil {
			return
		}
		mpp = make(map[int64]*model.PendantPackage, len(pps))
		for _, v := range pps {
			mpp[v.Pid] = v
		}
	}
	for _, v := range pids {
		var (
			status  int32
			expires int64
		)
		if pp, exists := mpp[v]; exists {
			status = pp.Status
			expires = pp.Expires
			if pe != nil && pe.Pid == pp.Pid {
				status = model.EquipPendantPKG
			}
		}
		pendant, exists := s.pendantMap[v]
		if !exists {
			log.Warn("gid(%d) pid(%d) pendant is empty", arg.GID, v)
			continue
		}
		psg := &model.GroupPendantList{
			PkgStatus:  status,
			PkgExpires: expires,
			Pendant:    pendant,
		}
		res = append(res, psg)
	}
	sort.Slice(res, func(i int, j int) bool {
		return res[i].Pendant.Rank < res[j].Pendant.Rank
	})
	return
}
