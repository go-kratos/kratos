package mobile

import (
	"context"
	"time"

	"go-common/app/interface/main/app-wall/model"
	"go-common/app/interface/main/app-wall/model/mobile"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_mobileKey = "mobile"
)

// InOrdersSync insert OrdersSync
func (s *Service) InOrdersSync(c context.Context, ip string, u *mobile.MobileXML, now time.Time) (err error) {
	if !s.iplimit(_mobileKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.InOrdersSync(c, u); err != nil || result == 0 {
		log.Error("mobile_s.dao.OrdersSync (%v) error(%v) or result==0", u, err)
	}
	return
}

// FlowSync update OrdersSync
func (s *Service) FlowSync(c context.Context, u *mobile.MobileXML, ip string) (err error) {
	if !s.iplimit(_mobileKey, ip) {
		err = ecode.AccessDenied
		return
	}
	var result int64
	if result, err = s.dao.FlowSync(c, u); err != nil || result == 0 {
		log.Error("mobile_s.dao.OrdersSync(%v) error(%v) or result==0", u, err)
	}
	return
}

// Activation
func (s *Service) Activation(c context.Context, usermob string, now time.Time) (msg string, err error) {
	rows := s.mobileInfo(c, usermob, now)
	res, ok := rows[usermob]
	if !ok {
		err = ecode.NothingFound
		msg = "该卡号尚未开通哔哩哔哩专属免流服务"
		return
	}
	for _, u := range res {
		if u.Actionid == 1 || (u.Actionid == 2 && now.Unix() <= int64(u.Expiretime)) {
			return
		}
	}
	err = ecode.NotModified
	msg = "该卡号哔哩哔哩专属免流服务已退订且已过期"
	return
}

// Activation
func (s *Service) MobileState(c context.Context, usermob string, now time.Time) (res *mobile.Mobile) {
	data := s.mobileInfo(c, usermob, now)
	res = s.userState(data, usermob, now)
	return
}

// UserFlowState
func (s *Service) UserMobileState(c context.Context, usermob string, now time.Time) (res *mobile.Mobile) {
	data := s.mobileInfo(c, usermob, now)
	if rows, ok := data[usermob]; ok {
		for _, res = range rows {
			if res.Actionid == 1 || (res.Actionid == 2 && now.Unix() <= int64(res.Expiretime)) {
				res.MobileType = 2
				return
			}
		}
		res = &mobile.Mobile{MobileType: 1}
	}
	res = &mobile.Mobile{MobileType: 1}
	return
}

// userState
func (s *Service) userState(user map[string][]*mobile.Mobile, usermob string, now time.Time) (res *mobile.Mobile) {
	if rows, ok := user[usermob]; !ok {
		res = &mobile.Mobile{MobileType: 1}
	} else {
		for _, res = range rows {
			if res.Actionid == 2 && now.Unix() <= int64(res.Expiretime) {
				res.MobileType = 4
				break
			} else if res.Actionid == 1 {
				res.MobileType = 2
				break
			}
		}
		if res.MobileType == 0 {
			res.MobileType = 3
		}
	}
	log.Info("mobile_state_type:%v mobile_state_usermob:%v", res.MobileType, usermob)
	return
}

// mobileInfo
func (s *Service) mobileInfo(c context.Context, usermob string, now time.Time) (res map[string][]*mobile.Mobile) {
	var (
		err  error
		m    []*mobile.Mobile
		tmps []*mobile.Mobile
		row  = map[string][]*mobile.Mobile{}
	)
	res = map[string][]*mobile.Mobile{}
	if m, err = s.dao.MobileCache(c, usermob); err == nil && len(m) > 0 {
		row[usermob] = m
		s.pHit.Incr("mobile_cache")
	} else {
		row, err = s.dao.OrdersUserFlow(c, usermob, now)
		if err != nil {
			log.Error("mobile_s.dao.OrdersUserFlow error(%v)", err)
			return
		}
		s.pMiss.Incr("mobile_cache")
		if user, ok := row[usermob]; ok && len(user) > 0 {
			if err = s.dao.AddMobileCache(c, usermob, user); err != nil {
				log.Error("s.dao.AddMobileCache error(%v)", err)
				return
			}
		}
	}
	if ms, ok := row[usermob]; ok && len(ms) > 0 {
		for _, m := range ms {
			tmp := &mobile.Mobile{}
			*tmp = *m
			tmp.Productid = ""
			tmps = append(tmps, tmp)
		}
		res[usermob] = tmps
	}
	return
}

// IsMobileIP is mobile ip
func (s *Service) IsMobileIP(ipUint uint32, ipStr, usermob string) (res *mobile.MobileUserIP) {
	res = &mobile.MobileUserIP{
		IPStr:    ipStr,
		IsValide: false,
	}
	if !model.IsIPv4(ipStr) {
		return
	}
	for _, u := range s.mobileIpCache {
		if u.IPStartUint <= ipUint && u.IPEndUint >= ipUint {
			res.IsValide = true
			return
		}
	}
	log.Error("mobile_user_ip:%v mobile_ip_usermob:%v", ipStr, usermob)
	return
}

// mobileIp ip limit
func (s *Service) iplimit(k, ip string) bool {
	return true
}
