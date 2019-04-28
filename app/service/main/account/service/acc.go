package service

import (
	"context"

	v1 "go-common/app/service/main/account/api"
	"go-common/app/service/main/account/model"
	coin "go-common/app/service/main/coin/model"
	member "go-common/app/service/main/member/model"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// Info get info by mid.
func (s *Service) Info(c context.Context, mid int64) (res *v1.Info, err error) {
	if res, err = s.dao.Info(c, mid); err == ecode.Degrade {
		err = nil
	}
	return
}

// Card get card by mid.
func (s *Service) Card(c context.Context, mid int64) (res *v1.Card, err error) {
	if res, err = s.dao.Card(c, mid); err == ecode.Degrade {
		err = nil
	}
	return
}

// Infos get info by mid.
func (s *Service) Infos(c context.Context, mids []int64) (res map[int64]*v1.Info, err error) {
	if res, err = s.dao.Infos(c, mids); err == ecode.Degrade {
		err = nil
	}
	return
}

// Cards get card by mid.
func (s *Service) Cards(c context.Context, mids []int64) (res map[int64]*v1.Card, err error) {
	if res, err = s.dao.Cards(c, mids); err == ecode.Degrade {
		err = nil
	}
	return
}

// Profile get profile by mid.
func (s *Service) Profile(c context.Context, mid int64) (res *v1.Profile, err error) {
	if res, err = s.dao.Profile(c, mid); err == ecode.Degrade {
		err = nil
	}
	return
}

// InfosByName multi get account info by names from cache otherwise account api.
// NOTE not cache used rarely
func (s *Service) InfosByName(c context.Context, names []string) (map[int64]*v1.Info, error) {
	if len(names) > 100 {
		names = names[:100]
	}
	mids, err := s.dao.MidsByName(c, names)
	if err != nil {
		return nil, err
	}
	return s.dao.Infos(c, mids)
}

// ProfileWithStat get profile by mid.
func (s *Service) ProfileWithStat(c context.Context, mid int64) (res *model.ProfileStat, err error) {
	p, err := s.dao.Profile(c, mid)
	if err != nil && err != ecode.Degrade {
		// err = errors.Wrap(err, "service profile with stat")
		// err = errors.WithStack(err)
		return
	}
	err = nil // NOTE: maybe err == ecode.Degrade
	res = &model.ProfileStat{
		Profile: p,
	}
	eg, errCtx := errgroup.WithContext(c)
	var le *member.LevelInfo
	eg.Go(func() (e error) {
		if le, e = s.dao.LevelExp(c, mid); e != nil {
			log.Error("s.dao.LevelExp(%d) error(%v)", mid, e)
			e = nil
		}
		return
	})
	var count float64
	eg.Go(func() (e error) {
		if count, e = s.coinRPC.UserCoins(errCtx, &coin.ArgCoinInfo{Mid: mid}); e != nil {
			log.Error("d.coinRP.UserCoins(%d) err(%v)", mid, e)
			e = nil
		}
		return
	})
	var rs *relation.Stat
	eg.Go(func() (e error) {
		if rs, e = s.relRPC.Stat(c, &relation.ArgMid{Mid: mid}); e != nil {
			log.Error("d.relRPC.Stat(%d) err(%v)", mid, e)
			e = nil
		}
		return
	})
	eg.Wait()
	if le != nil {
		res.LevelExp = *le
	}
	res.Coins = count
	if rs != nil {
		res.Following = rs.Following
		res.Follower = rs.Follower
	}
	return
}

// Privacy get privacy by mid.
func (s *Service) Privacy(c context.Context, mid int64) (res *model.Privacy, err error) {
	var (
		pProfile *model.PassportProfile
		rDetail  *member.RealnameDetail
		eg       errgroup.Group
		ip       = metadata.String(c, metadata.RemoteIP)
	)
	eg.Go(func() (e error) {
		if pProfile, e = s.dao.PassportProfile(c, mid, ip); e != nil {
			log.Error("s.dao.PassportProfile(%d) err(%+v)", mid, e)
			// e = nil
		}
		return
	})
	eg.Go(func() (e error) {
		if rDetail, e = s.dao.RealnameDetail(c, mid); e != nil {
			log.Error("s.dao.RealnameDetail(%d) err(%+v)", mid, e)
			e = nil
		}
		return
	})
	if err = eg.Wait(); err != nil {
		return
	}
	res = &model.Privacy{}
	if rDetail != nil && rDetail.RealnameBrief != nil {
		res.Realname = rDetail.Realname
		if rDetail.CardType == 0 {
			res.IdentityCard = rDetail.Card
			res.IdentitySex = rDetail.Gender
			res.HandIMG = rDetail.HandIMG
		}
	}
	if pProfile != nil {
		res.Tel = pProfile.Telphone
		res.RegTS = pProfile.JoinTime.Time().Unix()
		res.RegIP = pProfile.JoinIP
	}
	return
}
