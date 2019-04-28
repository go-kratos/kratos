package service

import (
	"context"

	"go-common/app/interface/main/tag/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	figModel "go-common/app/service/main/figure/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const _archiveInterval = 50

func (s *Service) batchArchives(c context.Context, aids []int64) (res map[int64]*api.Arc, err error) {
	var (
		tmpRes map[int64]*api.Arc
		n      = _archiveInterval
	)
	res = make(map[int64]*api.Arc, len(aids))
	for len(aids) > 0 {
		if n > len(aids) {
			n = len(aids)
		}
		arg := &archive.ArgAids2{Aids: aids[:n]}
		aids = aids[n:]
		if tmpRes, err = s.arcRPC.Archives3(c, arg); err != nil {
			log.Error("s.arcRPC.Archives3(%v) error(%v)", arg.Aids, err)
			err = nil
			continue
		}
		for k, v := range tmpRes {
			if _, ok := res[k]; !ok {
				res[k] = v
			}
		}
	}
	return
}

func (s *Service) archive(c context.Context, aid int64) (a *api.Arc, err error) {
	arg := &archive.ArgAid2{
		Aid:    aid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if a, err = s.arcRPC.Archive3(c, arg); err != nil {
		if ecode.Cause(err).Code() == ecode.NothingFound.Code() {
			err = nil
		} else {
			log.Error("s.arcRPC.Archive3(%d) error(%v)", aid, err)
		}
	}
	if a == nil {
		err = ecode.ArchiveNotExist
		return
	}
	if a.Aid == 0 {
		err = ecode.ArchiveNotExist
	}
	return
}

// check archive attribute normalArchive
func (s *Service) normalArchive(c context.Context, aid int64) (a *api.Arc, err error) {
	if a, err = s.archive(c, aid); err != nil {
		return
	}
	if !a.IsNormal() {
		log.Error("archive is not normal")
		err = ecode.TagOperateFail
	}
	return
}

// use figure(credit) info
func (s *Service) userFigure(c context.Context, mid int64) (score int8, err error) {
	arg := &figModel.ArgUserFigure{
		Mid: mid,
	}
	res, err := s.figRPC.UserFigure(c, arg)
	if err != nil {
		log.Error("s.figRPC.UserFigure()Arg:+%v error(%v)", arg, err)
		return
	}
	score = model.TotalScore - res.Percentage
	return
}
