package service

import (
	"context"
	"errors"

	tmod "go-common/app/admin/main/videoup-task/model"
	account "go-common/app/service/main/account/api"
	upsrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

//ERROR
var (
	ErrRPCEmpty = errors.New("rpc reply empty")
)

func (s *Service) profile(c context.Context, mid int64) (profile *account.ProfileStatReply, err error) {
	if profile, err = s.accRPC.ProfileWithStat3(c, &account.MidReq{Mid: mid}); err != nil {
		log.Error("s.accRPC.ProfileWithStat3(%d) error(%v)", mid, err)
	}
	return
}

func (s *Service) upSpecial(c context.Context) (ups map[int8]map[int64]struct{}, err error) {
	var (
		g                                                                  errgroup.Group
		whitegroup, blackgroup, policesgroup, enterprisegroup, signedgroup map[int64]struct{}
	)
	ups = make(map[int8]map[int64]struct{})

	f := func(gid int8) (map[int64]struct{}, error) {
		group := make(map[int64]struct{})
		mids, e := s.upGroupMids(c, int64(gid))
		if e != nil {
			return group, e
		}
		for _, mid := range mids {
			group[mid] = struct{}{}
		}
		return group, nil
	}
	g.Go(func() error {
		whitegroup, err = f(tmod.UpperTypeWhite)
		return err
	})
	g.Go(func() error {
		blackgroup, err = f(tmod.UpperTypeBlack)
		return err
	})
	g.Go(func() error {
		policesgroup, err = f(tmod.UpperTypePolitices)
		return err
	})
	g.Go(func() error {
		enterprisegroup, err = f(tmod.UpperTypeEnterprise)
		return err
	})
	g.Go(func() error {
		signedgroup, err = f(tmod.UpperTypeSigned)
		return err
	})
	if err = g.Wait(); err != nil {
		return
	}

	ups[tmod.UpperTypeWhite] = whitegroup
	ups[tmod.UpperTypeBlack] = blackgroup
	ups[tmod.UpperTypePolitices] = policesgroup
	ups[tmod.UpperTypeEnterprise] = enterprisegroup
	ups[tmod.UpperTypeSigned] = signedgroup
	return
}

func (s *Service) upGroupMids(c context.Context, gid int64) (mids []int64, err error) {
	var (
		total int
		maxps = 10000
		req   = &upsrpc.UpGroupMidsReq{
			Pn:      1,
			GroupID: gid,
			Ps:      maxps,
		}
		reply *upsrpc.UpGroupMidsReply
	)

	for {
		reply, err = s.upsRPC.UpGroupMids(c, req)
		if err == nil && (reply == nil || reply.Mids == nil) {
			err = ErrRPCEmpty
		}
		if err != nil {
			log.Error("UpGroupMids req(%+v) error(%v)", req, err)
			return
		}
		total = reply.Total
		mids = append(mids, reply.Mids...)
		if reply.Size() != maxps {
			break
		}
		req.Pn++
	}
	log.Info("upGroupMids(%d) reply total(%d) len(%d)", gid, total, len(mids))
	return
}
