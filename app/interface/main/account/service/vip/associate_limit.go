package vip

import (
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

// ActivityTimeLimit activity time limit.
func (s *Service) ActivityTimeLimit(mid int64) error {
	if len(s.c.Vipproperty.AssociateWhiteMidMap) > 0 && mid != 0 {
		for _, v := range s.c.Vipproperty.AssociateWhiteMidMap {
			if v == mid {
				return nil
			}
		}
	}
	now := time.Now().Unix()
	if s.c.Vipproperty.ActStartTime > now {
		return ecode.VipActivityNotStart
	}
	if s.c.Vipproperty.ActEndTime < now {
		return ecode.VipActivityHadEnd
	}
	return nil
}

// ActivityWhiteIPLimit act ip limit.
func (s *Service) ActivityWhiteIPLimit(appkey string, ip string) error {
	var (
		whiteips []string
		ok       bool
	)
	if whiteips, ok = s.c.Vipproperty.AssociateWhiteIPMap[appkey]; !ok {
		log.Error("act ip limit appkey(%s) ip(%s)", appkey, ip)
		return ecode.VipWhiteIPListErr
	}
	for _, v := range whiteips {
		if v == ip {
			return nil
		}
	}
	log.Error("act ip limit appkey(%s) ip(%s)", appkey, ip)
	return ecode.VipWhiteIPListErr
}

// ActivityWhiteOutOpenIDLimit act out open id limit.
func (s *Service) ActivityWhiteOutOpenIDLimit(openid string) error {
	if len(s.c.Vipproperty.AssociateWhiteOutOpenIDMap) > 0 && openid != "" {
		for _, v := range s.c.Vipproperty.AssociateWhiteOutOpenIDMap {
			if v == openid {
				return nil
			}
		}
	}
	now := time.Now().Unix()
	if s.c.Vipproperty.ActStartTime > now {
		return ecode.VipActivityNotStart
	}
	if s.c.Vipproperty.ActEndTime < now {
		return ecode.VipActivityHadEnd
	}
	return nil
}
