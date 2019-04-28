package service

import (
	"context"
	"errors"
	"time"

	tmod "go-common/app/job/main/videoup-report/model/task"
	account "go-common/app/service/main/account/api"
	upsrpc "go-common/app/service/main/up/api/v1"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"fmt"
	"math"
)

//ERROR
var (
	ErrRPCEmpty = errors.New("rpc reply empty")
)

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

//
func (s *Service) profile(c context.Context, mid int64) (p *account.ProfileStatReply, err error) {
	if p, err = s.accRPC.ProfileWithStat3(c, &account.MidReq{Mid: mid}); err != nil {
		p = nil
		log.Error("s.accRPC.ProfileWithStat3(%d) error(%v)", mid, err)
	}
	return
}

func (s *Service) getUpperFans(c context.Context, mid int64) (fans int64, failed bool) {
	card, err := s.profile(c, mid)
	if err != nil {
		failed = true
		log.Error("s.profile(mid=%d) error(%v)", mid, err)
		return
	}

	fans = card.Follower
	log.Info("s.profile(mid=%d) fans(%d)", mid, fans)
	return
}

func (s *Service) isWhite(mid int64) bool {
	if ups, ok := s.upperCache[tmod.UpperTypeWhite]; ok {
		_, isWhite := ups[mid]
		return isWhite
	}
	return false
}

func (s *Service) isBlack(mid int64) bool {
	if ups, ok := s.upperCache[tmod.UpperTypeBlack]; ok {
		_, isBlack := ups[mid]
		return isBlack
	}
	return false
}

func (s *Service) isPolitices(mid int64) bool {
	if ups, ok := s.upperCache[tmod.UpperTypePolitices]; ok {
		_, isShiZheng := ups[mid]
		return isShiZheng
	}
	return false
}

func (s *Service) isEnterprise(mid int64) bool {
	if ups, ok := s.upperCache[tmod.UpperTypeEnterprise]; ok {
		_, isQiYe := ups[mid]
		return isQiYe
	}
	return false
}

func (s *Service) isSigned(mid int64) bool {
	if ups, ok := s.upperCache[tmod.UpperTypeSigned]; ok {
		_, signed := ups[mid]
		return signed
	}
	return false
}

// Until next day x hours
func nextDay(hour int) time.Duration {
	n := time.Now().Add(24 * time.Hour)
	d := time.Date(n.Year(), n.Month(), n.Day(), hour, 0, 0, 0, n.Location())
	return time.Until(d)
}

func secondsFormat(sec int) (str string) {
	if sec < 0 {
		return "--:--:--"
	}
	if sec == 0 {
		return "00:00:00"
	}
	h := math.Floor(float64(sec) / 3600)
	m := math.Floor((float64(sec) - 3600*h) / 60)
	se := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", int64(h), int64(m), se)
}
