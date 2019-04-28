package operation

import (
	"context"
	"regexp"
	"strconv"

	opmdl "go-common/app/interface/main/web-show/model/operation"
	"go-common/app/service/main/archive/api"
	comarcmdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emptyPromoteMap = make(map[string][]*opmdl.Promote)
	_avReg           = regexp.MustCompile(`video\/av[0-9]+`)
)

// Promote Service
func (s *Service) Promote(c context.Context, arg *opmdl.ArgPromote) (res map[string][]*opmdl.Promote, err error) {
	var (
		ok   bool
		arcs map[int64]*api.Arc
		arc  *api.Arc
		aid  int64
		aids []int64
	)
	opMap := s.operation(arg.Tp, arg.Rank, arg.Count)
	for _, ops := range opMap {
		for _, op := range ops {
			if aid, err = s.regAid(op.Link); err != nil {
				log.Error("service.regAid error(%v)", err)
				continue
			}
			op.Aid = aid
			aids = append(aids, aid)
		}
	}
	argAids := &comarcmdl.ArgAids2{
		Aids: aids,
	}
	if arcs, err = s.arcRPC.Archives3(c, argAids); err != nil {
		log.Error("s.arcRPC.Archives2(arcAids:(%v), arcs), err(%v)", aids, err)
		res = _emptyPromoteMap
		return
	}
	res = make(map[string][]*opmdl.Promote)
	for rk, ops := range opMap {
		promotes := make([]*opmdl.Promote, 0, len(ops))
		for _, op := range ops {
			if arc, ok = arcs[op.Aid]; !ok {
				continue
			}
			promote := &opmdl.Promote{
				IsAd:    int8(op.Ads),
				Archive: arc,
			}
			promotes = append(promotes, promote)
		}
		res[rk] = promotes
	}
	return
}

// regAid Service
func (s *Service) regAid(link string) (aid int64, err error) {
	avStr := _avReg.FindString(link)
	if avStr != "" {
		aidStr := avStr[8:]
		if aid, err = strconv.ParseInt(aidStr, 10, 64); err != nil {
			log.Error("strconv.ParseInt error(%v)", err)
			return
		}
	} else {
		err = ecode.ArchiveNotExist
	}
	return
}
