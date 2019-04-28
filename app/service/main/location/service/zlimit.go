package service

import (
	"context"

	accoutCli "go-common/app/service/main/account/api"
	"go-common/app/service/main/location/model"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// PgcZone get group_id by zoneid .
func (s *Service) PgcZone(c context.Context, zoneIDs []int64) (res map[string][]int64, err error) {
	var (
		ok       bool
		fb, al   map[int64]int64
		gidNotIn = []int64{}
		gidIn    = []int64{}
	)
	res = make(map[string][]int64)
	for gid, authZid := range s.groupAuthZone {
		// forbidden priority.
		for k, zoneID := range zoneIDs {
			if fb, ok = authZid[model.Forbidden]; ok {
				if _, ok = fb[zoneID]; ok {
					gidNotIn = append(gidNotIn, gid)
					break
				} else {
					// final value
					if k == len(zoneIDs)-1 {
						gidIn = append(gidIn, gid)
					}
				}
			} else if al, ok = authZid[model.Allow]; ok {
				if _, ok = al[zoneID]; ok {
					gidIn = append(gidIn, gid)
					break
				} else {
					// final value
					if k == len(zoneIDs)-1 {
						gidNotIn = append(gidNotIn, gid)
					}
				}
			}
		}
	}
	res["not_in"] = gidNotIn
	res["in"] = gidIn
	return
}

func (s *Service) checkLimit(c context.Context, mid, code int64, ip string) int64 {
	switch code {
	case model.Allow, model.Forbidden:
		return code
	case model.Formal:
		if mid == 0 {
			return model.Forbidden
		}
		if card, err := s.accCard(c, mid, ip); err != nil {
			return model.Forbidden
		} else if card != nil && card.Level < 1 && (card.Vip.Type == 0 || card.Vip.Status == 0 || card.Vip.Status == 2 || card.Vip.Status == 3) {
			return model.Forbidden
		} else {
			return model.Allow
		}

	case model.Pay:
		// TODO VIP member
	}
	return model.Forbidden
}

func (s *Service) accCard(c context.Context, mid int64, ip string) (ai *accoutCli.Card, err error) {
	arg := &accoutCli.MidReq{
		Mid: mid,
	}
	cr := &accoutCli.CardReply{}
	if cr, err = s.accountSvc.Card3(c, arg); err != nil {
		err = errors.Wrapf(err, "mid(%d)", mid)
		return
	}
	return cr.Card, nil
}

// Auth get auth by aid and ipaddr & check mid.
func (s *Service) Auth(c context.Context, aid, mid int64, ipaddr, cdnip string) (ret int64, err error) {
	var retdown int64
	if ret, retdown, err = s.Auth2(c, aid, ipaddr, cdnip); err != nil {
		err = errors.Wrapf(err, "s.Auth2(%d, %s, %s)", aid, ipaddr, cdnip)
		return
	}
	s.authCodePorm.Incr(model.PlayAuth[ret])
	s.authCodePorm.Incr(model.DownAuth[retdown])
	ret = s.checkLimit(c, mid, ret, ipaddr)
	return
}

// Archive2 get auth by aid and ipaddr & check mid.
func (s *Service) Archive2(c context.Context, aid, mid int64, ipaddr, cdnip string) (res *model.Auth, err error) {
	var ret, retdown int64
	if ret, retdown, err = s.Auth2(c, aid, ipaddr, cdnip); err != nil {
		err = errors.Wrapf(err, "s.Auth2(%d, %s, %s)", aid, ipaddr, cdnip)
		return
	}
	s.authCodePorm.Incr(model.PlayAuth[ret])
	s.authCodePorm.Incr(model.DownAuth[retdown])
	ret = s.checkLimit(c, mid, ret, ipaddr)
	res = &model.Auth{Play: ret, Down: retdown}
	return
}

// Auth2 get auth by aid and ipaddr.
func (s *Service) Auth2(c context.Context, aid int64, ipaddr, cdnip string) (ret, retdown int64, err error) {
	var (
		ok                  bool
		auth, pid, zid, gid int64
		rules, pids         []int64
		zids                map[int64]int64
		ipInfo              *model.InfoComplete
	)
	ipInfo, _ = s.InfoComplete(c, ipaddr)
	if (ipInfo == nil || (ipInfo != nil && s.filterInnerIP(ipInfo))) && cdnip != "" {
		log.Info("ip(%v) aid(%v) cdnip(%v), ipaddr is nil or is filter zone", ipaddr, aid, cdnip)
		s.innerIPPorm.Add("archive", 1)
		ipInfo, _ = s.InfoComplete(c, cdnip)
	}
	if ipInfo == nil {
		ret = model.Allow
		retdown = model.AllowDown
		return
	}
	uz := ipInfo.ZoneID // country, state, city
	if ok, err = s.zdb.ExistsAuth(c, aid); err != nil {
		log.Error("s.zdb.ExistsAuth error(%+v)", err)
		err = nil
	} else if ok {
		if rules, err = s.zdb.Auth(c, aid, uz); err != nil {
			log.Error("s.zdb.Auth(%d) error(%+v) ", aid, err)
			err = nil
		} else {
			for _, auth = range rules {
				retdown = 0xff & auth
				ret = auth >> 8
				if ret != 0 {
					break
				}
			}
			if ret == 0 {
				ret = model.Allow
				retdown = model.AllowDown
			}
			return
		}
	}
	s.missedPorm.Incr("redis_missed")
	if gid, err = s.zdb.Groupid(c, aid); err != nil {
		return
	} else if gid != 0 {
		if pids, ok = s.groupPolicy[gid]; ok {
			for _, pid = range pids {
				if zids, ok = s.policy[pid]; !ok {
					continue
				}
				if ret == 0 {
					//  ret already set skip check
					for _, zid = range uz {
						if auth, ok = zids[zid]; ok {
							if ret == 0 {
								retdown = 0xff & auth
								ret = auth >> 8 // ret must not be zero
								break
							}
						}
					}
				}
				tmpZids := map[int64]map[int64]int64{
					aid: zids,
				}
				s.addCache(tmpZids)
			}
			if ret == 0 {
				s.missedPorm.Incr("local_policy_cache_missed")
				ret = model.Allow
				retdown = model.AllowDown
			}
			return
		}
		s.missedPorm.Incr("local_group_cache_missed")
	}
	s.missedPorm.Incr("db_missed")
	ret = model.Allow
	retdown = model.AllowDown
	zids = make(map[int64]int64)
	zids[0] = ret<<8 | retdown
	tmpZids := map[int64]map[int64]int64{
		aid: zids,
	}
	s.addCache(tmpZids)
	return
}

// AuthGID auth by group_id and ipaddr(or cdnip) & check mid.
func (s *Service) AuthGID(c context.Context, gid, mid int64, ipaddr, cdnip string) (res *model.Auth) {
	ret, retdown := s.AuthGID2(c, gid, ipaddr, cdnip)
	s.authCodePorm.Incr(model.PlayAuth[ret])
	s.authCodePorm.Incr(model.DownAuth[retdown])
	ret = s.checkLimit(c, mid, ret, ipaddr)
	res = &model.Auth{
		Play: ret,
		Down: retdown,
	}
	return
}

// AuthGIDs auth by group_id and ipaddr(or cdnip) & check mid.
func (s *Service) AuthGIDs(c context.Context, gids []int64, mid int64, ipaddr, cdnip string) (res map[int64]*model.Auth) {
	res = make(map[int64]*model.Auth)
	for _, gid := range gids {
		ret, retdown := s.AuthGID2(c, gid, ipaddr, cdnip)
		s.authCodePorm.Incr(model.PlayAuth[ret])
		s.authCodePorm.Incr(model.DownAuth[retdown])
		ret = s.checkLimit(c, mid, ret, ipaddr)
		res[gid] = &model.Auth{
			Play: ret,
			Down: retdown,
		}
	}
	return
}

// AuthGID2 auth by group_id and ipaddr(or cdnip).
func (s *Service) AuthGID2(c context.Context, gid int64, ipaddr, cdnip string) (ret, retdown int64) {
	var (
		ipInfo  *model.InfoComplete
		hitMark string
	)
	ret = model.Allow
	retdown = model.AllowDown
	ipInfo, _ = s.InfoComplete(c, ipaddr)
	if (ipInfo == nil || (ipInfo != nil && s.filterInnerIP(ipInfo))) && cdnip != "" {
		log.Info("ip(%v) gid(%v) cdnip(%v), ipaddr is nil or is filter zone", ipaddr, gid, cdnip)
		s.innerIPPorm.Add("group", 1)
		ipInfo, _ = s.InfoComplete(c, cdnip)
	}
	if ipInfo == nil {
		return
	}
	if pids, ok := s.groupPolicy[gid]; ok {
		for _, pid := range pids {
			if _, ok := s.policy[pid]; !ok {
				continue
			}
			for _, zoneid := range ipInfo.ZoneID {
				if auth, ok := s.policy[pid][zoneid]; ok {
					retdown = 0xff & auth
					ret = auth >> 8
					hitMark = "hit"
					break
				}
			}
		}
	}
	if len(hitMark) == 0 {
		s.missedPorm.Incr("local_cache_missed")
	}
	return
}

// AuthPIDs check by policy_ids and ipaddr
func (s *Service) AuthPIDs(c context.Context, pidStr, ipaddr, cdnip string) (res map[int64]*model.Auth, err error) {
	var (
		pids   []int64
		ipInfo *model.InfoComplete
	)
	res = make(map[int64]*model.Auth, len(pids))
	if pids, err = xstr.SplitInts(pidStr); err != nil || len(pids) == 0 {
		log.Error("xstr.SplitInts(%v) error(%v) or pids is empty", pidStr, err)
		return
	}
	ipInfo, _ = s.InfoComplete(c, ipaddr)
	if (ipInfo == nil || (ipInfo != nil && s.filterInnerIP(ipInfo))) && cdnip != "" {
		ipInfo, _ = s.InfoComplete(c, cdnip)
	}
	if ipInfo == nil {
		return
	}
	for _, pid := range pids {
		res[pid] = &model.Auth{
			Play: model.Allow,
			Down: model.AllowDown,
		}
		if _, ok := s.policy[pid]; !ok {
			continue
		}
		for _, zoneid := range ipInfo.ZoneID {
			if auth, ok := s.policy[pid][zoneid]; ok {
				res[pid].Down = 0xff & auth
				res[pid].Play = auth >> 8
				break
			}
		}
	}
	return
}

func (s *Service) filterInnerIP(ip *model.InfoComplete) bool {
	for _, zone := range s.c.FilterZone {
		if ip.Country == zone || ip.Province == zone || ip.City == zone {
			return true
		}
	}
	return false
}
