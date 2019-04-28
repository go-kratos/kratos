package unicom

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/model/unicom"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	log "go-common/library/log"
)

const (
	_initIPUnicomKey = "ipunicom_%v_%v"
)

// loadUnicomIP load unicom ip
func (s *Service) loadUnicomIP() {
	unicomIP, err := s.dao.IPSync(context.TODO())
	if err != nil {
		log.Error("s.dao.IPSync error(%v)", err)
		return
	}
	s.unicomIpCache = unicomIP
	tmp := map[string]*unicom.UnicomIP{}
	for _, u := range s.unicomIpCache {
		key := fmt.Sprintf(_initIPUnicomKey, u.Ipbegin, u.Ipend)
		tmp[key] = u
	}
	s.unicomIpSQLCache = tmp
	log.Info("loadUnicomIPCache success")
}

// loadIPlimit ip limit
func (s *Service) loadIPlimit(c *conf.Config) {
	hosts := make(map[string]struct{}, len(c.IPLimit.Addrs))
	for k, v := range c.IPLimit.Addrs {
		for _, ipStr := range v {
			key := fmt.Sprintf(_initIPlimitKey, k, ipStr)
			if _, ok := hosts[key]; !ok {
				hosts[key] = struct{}{}
			}
		}
	}
	s.operationIPlimit = hosts
	log.Info("load operationIPlimit success")
}

// loadUnicomIPOrder load unciom ip order update
func (s *Service) loadUnicomIPOrder(now time.Time) {
	unicomIP, err := s.dao.UnicomIP(context.TODO(), now)
	if err != nil {
		log.Error("s.dao.UnicomIP(%v)", err)
		return
	}
	if len(unicomIP) == 0 {
		log.Info("unicom ip orders is null")
		return
	}
	tx, err := s.dao.BeginTran(context.TODO())
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}
	for _, uip := range unicomIP {
		key := fmt.Sprintf(_initIPUnicomKey, uip.Ipbegin, uip.Ipend)
		if _, ok := s.unicomIpSQLCache[key]; ok {
			delete(s.unicomIpSQLCache, key)
			continue
		}
		var (
			result int64
		)
		if result, err = s.dao.InUnicomIPSync(tx, uip, time.Now()); err != nil || result == 0 {
			tx.Rollback()
			log.Error("s.dao.InUnicomIPSync error(%v)", err)
			return
		}
	}
	for _, uold := range s.unicomIpSQLCache {
		var (
			result int64
		)
		if result, err = s.dao.UpUnicomIP(tx, uold.Ipbegin, uold.Ipend, 0, time.Now()); err != nil || result == 0 {
			tx.Rollback()
			log.Error("s.dao.UpUnicomIP error(%v)", err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error(%v)", err)
		return
	}
	log.Info("update unicom ip success")
}

func (s *Service) loadUnicomPacks() {
	pack, err := s.dao.UserPacks(context.TODO())
	if err != nil {
		log.Error("s.dao.UserPacks error(%v)", err)
		return
	}
	s.unicomPackCache = pack
}

func (s *Service) loadUnicomFlow() {
	var (
		list map[string]*unicom.UnicomUserFlow
		err  error
	)
	if list, err = s.dao.UserFlowListCache(context.TODO()); err != nil {
		log.Error("load unicom s.dao.UserFlowListCache error(%v)", err)
		return
	}
	for key, u := range list {
		var (
			c           = context.TODO()
			requestNo   int64
			orderstatus string
			msg         string
		)
		if err = s.dao.UserFlowCache(c, key); err != nil {
			if err == memcache.ErrNotFound {
				if err = s.returnPoints(c, u); err != nil {
					if err != ecode.NothingFound {
						log.Error("load unicom s.returnPoints error(%v)", err)
						continue
					}
					err = nil
				}
				log.Info("load unicom userbind timeout flow(%v)", u)
			} else {
				log.Error("load unicom s.dao.UserFlowCache error(%v)", err)
				continue
			}
		} else {
			if requestNo, err = s.seqdao.SeqID(c); err != nil {
				log.Error("load unicom s.seqdao.SeqID error(%v)", err)
				continue
			}
			if orderstatus, msg, err = s.dao.FlowQry(c, u.Phone, requestNo, u.Outorderid, u.Orderid, time.Now()); err != nil {
				log.Error("load unicom s.dao.FlowQry error(%v) msg(%s)", err, msg)
				continue
			}
			log.Info("load unicom userbind flow(%v) orderstatus(%s)", u, orderstatus)
			if orderstatus == "00" {
				continue
			} else if orderstatus != "01" {
				if err = s.returnPoints(c, u); err != nil {
					if err != ecode.NothingFound {
						log.Error("load unicom s.returnPoints error(%v)", err)
						continue
					}
					err = nil
				}
			}
		}
		delete(list, key)
		if err = s.dao.DeleteUserFlowCache(c, key); err != nil {
			log.Error("load unicom s.dao.DeleteUserFlowCache error(%v)", err)
			continue
		}
	}
	if err = s.dao.AddUserFlowListCache(context.TODO(), list); err != nil {
		log.Error("load unicom s.dao.AddUserFlowListCache error(%v)", err)
		return
	}
	log.Info("load unicom flow success")
}

// returnPoints retutn user integral and flow
func (s *Service) returnPoints(c context.Context, u *unicom.UnicomUserFlow) (err error) {
	var (
		userbind *unicom.UserBind
		result   int64
	)
	if userbind, err = s.unicomBindInfo(c, u.Mid); err != nil {
		return
	}
	ub := &unicom.UserBind{}
	*ub = *userbind
	ub.Flow = ub.Flow + u.Flow
	ub.Integral = ub.Integral + u.Integral
	if err = s.dao.AddUserBindCache(c, ub.Mid, ub); err != nil {
		log.Error("unicom s.dao.AddUserBindCache error(%v)", err)
		return
	}
	if result, err = s.dao.UpUserIntegral(c, ub); err != nil || result == 0 {
		log.Error("unicom s.dao.UpUserIntegral error(%v) or result==0", err)
		return
	}
	var packInt int
	if u.Integral > 0 {
		packInt = u.Integral
	} else {
		packInt = u.Flow
	}
	log.Info("unicom_pack(%v) mid(%v)", u.Desc+",领取失败并返还", userbind.Mid)
	s.unicomPackInfoc(userbind.Usermob, u.Desc+",领取失败并返还", u.Orderid, userbind.Phone, packInt, 0, userbind.Mid, time.Now())
	return
}
