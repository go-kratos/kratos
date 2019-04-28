package service

import (
	"context"
	"fmt"

	"go-common/app/service/main/archive/api"
	arcmdl "go-common/app/service/main/archive/model/archive"
	rankmdl "go-common/app/service/main/rank/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ArcRPC find archive by rpc
func (s *Service) ArcRPC(c context.Context, aid int64) (a *api.Arc, err error) {
	argAid := &arcmdl.ArgAid2{
		Aid: aid,
	}
	if a, err = s.arcRPC.Archive3(c, argAid); err != nil {
		log.Error("s.arcRPC.Archive3(%v), error(%v)", argAid, err)
	}
	if !a.IsNormal() {
		err = ecode.ArchiveNotExist
	}
	return
}

// ArcsRPC find archives by rpc
func (s *Service) ArcsRPC(c context.Context, aids []int64) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	argAids := &arcmdl.ArgAids2{
		Aids: aids,
	}
	if as, err = s.arcRPC.Archives3(c, argAids); err != nil {
		log.Error("s.arcRPC.Archives3(%v, archives), err(%v)", argAids, err)
	}
	return
}

// TypeidsRPC find tids by rpc
func (s *Service) TypeidsRPC(c context.Context, oids []int64) (res *rankmdl.GroupResp, err error) {
	if len(oids) == 0 {
		return
	}
	arg := &rankmdl.GroupReq{
		Business: rankmdl.BusinessArchive,
		Oids:     oids,
		Field:    "pid",
	}
	if res, err = s.rankRPC.Group(c, arg); err != nil {
		log.Error("s.rankRPC.Group(%+v), error(%v)", arg, err)
	}
	return
}

// SortArcsRPC sort oids by rpc
func (s *Service) SortArcsRPC(c context.Context, tid, tv, pn, ps int, field string, oids []int64) (res *rankmdl.SortResp, err error) {
	if len(oids) == 0 {
		return
	}
	fmap := make(map[string]string)
	if tid != 0 {
		fmap["pid"] = fmt.Sprintf("%d", tid)
	}
	if tv != 0 {
		fmap["result"] = "1"
		fmap["deleted"] = "0"
		fmap["valid"] = "1"
	}
	arg := &rankmdl.SortReq{
		Business: rankmdl.BusinessArchive,
		Oids:     oids,
		Order:    rankmdl.RankOrderByDesc,
		Field:    field,
		Filters:  fmap,
		Pn:       pn,
		Ps:       ps,
	}
	if res, err = s.rankRPC.Sort(c, arg); err != nil {
		log.Error("s.rankRPC.Sort(%+v), error(%v)", arg, err)
	}
	return
}
