package service

import (
	"fmt"

	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// loadTypes is used for gettting archive data from rpc
func (s *Service) loadTypes() (err error) {
	var (
		res       map[int32]*arcmdl.Tp
		resRel    = make(map[int32][]int32)
		typeReply *arcmdl.TypesReply
	)
	if typeReply, err = s.arcClient.Types(ctx, &arcmdl.NoArgRequest{}); err != nil {
		log.Error("arcRPC loadType Error %v", err)
		return
	}
	res = typeReply.Types
	if len(res) == 0 {
		log.Error("arcRPC loadType Empty")
		return
	}
	for _, tInfo := range res {
		if _, ok := resRel[tInfo.Pid]; !ok {
			resRel[tInfo.Pid] = []int32{tInfo.ID}
			continue
		}
		resRel[tInfo.Pid] = append(resRel[tInfo.Pid], tInfo.ID)
	}
	s.ArcTypes = res
	s.arcPTids = resRel
	return
}

// arcPName is used for get arc first partition with typeID(second partition)
func (s *Service) arcPName(cID int32) (name string, pid int32, err error) {
	var (
		c, p *arcmdl.Tp
		ok   bool
		code = ecode.RequestErr.Code()
	)
	if c, ok = s.ArcTypes[cID]; !ok {
		err = errors.Wrap(ecode.Int(code), fmt.Sprintf("can't find type for ID: %d ", cID))
		return
	}
	if p, ok = s.ArcTypes[c.Pid]; !ok {
		err = errors.Wrap(ecode.Int(code), fmt.Sprintf("can't find type for ID: %d, parent id: %d", cID, c.Pid))
		return
	}
	return p.Name, c.Pid, nil
}

//Contains is used for check string in array
func (s *Service) Contains(tid int32) (contain bool) {
	var (
		name string
		err  error
	)
	if name, _, err = s.arcPName(tid); err != nil {
		log.Warn("s.CheckArc.arcPName Tid %d, error(%v)", tid, err)
		return
	}
	for _, v := range s.c.Cfg.PGCTypes {
		if v == name {
			return true
		}
	}
	return
}
