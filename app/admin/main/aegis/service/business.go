package service

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/middleware"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddBusiness .
func (s *Service) AddBusiness(c context.Context, b *business.Business) (id int64, err error) {
	if id, err = s.gorm.AddBusiness(c, b); err != nil {
		log.Error("s.gorm.AddBusiness error(%v)", err)
	}
	return
}

// UpdateBusiness .
func (s *Service) UpdateBusiness(c context.Context, b *business.Business) (err error) {
	if err = s.gorm.UpdateBusiness(c, b); err != nil {
		log.Error("s.gorm.UpdateBusiness error(%v)", err)
	}
	return
}

// SetBusinessState .
func (s *Service) SetBusinessState(c context.Context, b *business.Business) (err error) {
	if int8(b.State) == business.StateEnable {
		if err = s.gorm.EnableBusiness(c, b.ID); err != nil {
			log.Error("s.gorm.EnableBusiness error(%v)", err)
		}
	} else {
		if err = s.gorm.DisableBusiness(c, b.ID); err != nil {
			log.Error("s.gorm.EnableBusiness error(%v)", err)
		}
	}
	return
}

// Business .
func (s *Service) Business(c context.Context, b *business.Business) (res *business.Business, err error) {
	if res, err = s.gorm.Business(c, b.ID); err != nil {
		log.Error("s.gorm.Business error(%v)", err)
	}
	return
}

// BusinessList .
func (s *Service) BusinessList(c context.Context, businessID []int64, onlyEnable bool) (res []*business.Business, err error) {
	if res, err = s.gorm.BusinessList(c, 0, businessID, onlyEnable); err != nil {
		log.Error("s.gorm.BusinessList error(%v) ids(%+v)", err, businessID)
	}

	s.mulIDtoName(c, res, s.transUnames, "UID", "UserName")
	return
}

func (s *Service) transUnames(c context.Context, uids []int64) (unames map[int64][]interface{}, err error) {
	unames = make(map[int64][]interface{})
	uns, err := s.http.GetUnames(c, uids)
	if err != nil {
		return
	}
	for uid, uname := range uns {
		unames[uid] = []interface{}{uname}
	}
	return
}

// AddBizCFG .
func (s *Service) AddBizCFG(c context.Context, b *business.BizCFG) (lastid int64, err error) {
	return s.gorm.AddBizConfig(c, b)
}

// UpdateBizCFG .
func (s *Service) UpdateBizCFG(c context.Context, b *business.BizCFG) (err error) {
	if err = s.gorm.EditBizConfig(c, b); err != nil {
		log.Error("s.gorm.UpdateBizCFG error(%v)", err)
		return
	}

	delete(s.bizCfgCache, fmt.Sprintf("%d_%d", b.BusinessID, b.TP))
	delete(s.bizRoleCache, b.BusinessID)
	delete(s.taskRoleCache, b.BusinessID)
	delete(s.bizMiddlewareCache, b.BusinessID)
	return
}

// ListBizCFGs .
func (s *Service) ListBizCFGs(c context.Context, bizid int64) (res []*business.BizCFG, err error) {
	if res, err = s.gorm.GetConfigs(c, bizid); err != nil {
		log.Error("s.gorm.ListBizCFGs error(%v)", err)
	}
	return
}

// ReserveCFG 保留字配置
func (s *Service) ReserveCFG(c context.Context, bizid int64) (rsvcfg map[string]interface{}, err error) {
	var rsv string
	if rsv = s.getConfig(c, bizid, business.TypeReverse); len(rsv) == 0 {
		return
	}

	rsvcfg = make(map[string]interface{})

	if err = json.Unmarshal([]byte(rsv), &rsvcfg); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		err = ecode.AegisReservedCfgErr
	}
	return
}

// AttributeCFG  属性位配置
func (s *Service) AttributeCFG(c context.Context, bizid int64) (attrcfg map[string]uint, err error) {
	var attr string
	if attr = s.getConfig(c, bizid, business.TypeAttribute); len(attr) == 0 {
		return
	}

	attrcfg = make(map[string]uint)

	if err = json.Unmarshal([]byte(attr), &attrcfg); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
	}
	return
}

func (s *Service) syncBizCache(c context.Context) (err error) {
	var (
		bizcfg []*business.BizCFG
	)
	if bizcfg, err = s.gorm.ActiveConfigs(c); err != nil {
		return
	}

	bizTpMap := map[string]string{}
	bizRoleMap := map[int64]map[string]int64{}
	taskRoleMap := map[int64]map[int64][]int64{}
	bizMiddleware := map[int64][]*middleware.Aggregate{}
	for _, item := range bizcfg {
		bizTpMap[fmt.Sprintf("%d_%d", item.BusinessID, item.TP)] = item.Config

		biz, taskRole, err := item.FormatMngBID()
		if err == nil && len(taskRole) > 0 {
			taskRoleMap[biz] = taskRole
		}

		bizID, bizRoles, err := item.FormatBizBID()
		if err == nil && len(bizRoles) > 0 {
			bizRoleMap[bizID] = bizRoles
		}

		mws, err := item.FormatAggregate()
		if err == nil && len(mws) > 0 {
			bizMiddleware[item.BusinessID] = mws
		}
	}
	s.bizCfgCache = bizTpMap
	s.bizRoleCache = bizRoleMap
	s.taskRoleCache = taskRoleMap
	s.bizMiddlewareCache = bizMiddleware
	return
}

func (s *Service) getConfig(c context.Context, bizid int64, tp int8) (cfg string) {
	var ok bool
	if cfg, ok = s.bizCfgCache[fmt.Sprintf("%d_%d", bizid, tp)]; ok {
		return cfg
	}
	s.syncBizCache(c)
	cfg = s.bizCfgCache[fmt.Sprintf("%d_%d", bizid, tp)]
	return
}

func (s *Service) getTaskRole(c context.Context, bizid int64) (cfgs map[int64][]int64) {
	var exist bool
	if cfgs, exist = s.taskRoleCache[bizid]; exist {
		return cfgs
	}

	s.syncBizCache(c)
	cfgs = s.taskRoleCache[bizid]
	return
}

func (s *Service) getTaskBiz(c context.Context, bids []int64) (bizs []int64, flows []int64) {
	bizs = []int64{}
	flows = []int64{}

	s.syncBizCache(c)
	for _, bid := range bids {
		for biz, item := range s.taskRoleCache {
			if v, exist := item[bid]; exist && len(v) > 0 {
				bizs = append(bizs, biz)
				flows = append(flows, v...)
			}
		}
	}
	return
}

func (s *Service) getBizRole(c context.Context, mngID int64) (bizRole map[string]int64, bizID int64) {
	for biz, item := range s.bizRoleCache {
		if item[business.BizBIDMngID] == mngID {
			bizID = biz
			bizRole = item
			return
		}
	}

	s.syncBizCache(c)
	for biz, item := range s.bizRoleCache {
		if item[business.BizBIDMngID] == mngID {
			bizID = biz
			bizRole = item
			break
		}
	}
	return
}
