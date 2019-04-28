package service

import (
	"context"

	"go-common/app/admin/main/aegis/model/business"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/task"
	"go-common/library/log"
)

// GetRole .
func (s *Service) GetRole(c context.Context, opt *common.BaseOptions) (uname string, role int8, err error) {
	role, uname, err = s.mc.GetRole(c, opt)
	if err == nil && role > 0 {
		return
	}

	var (
		cfgs  map[int64][]int64 //bid:flowids
		mngid int64
	)

	if cfgs = s.getTaskRole(c, opt.BusinessID); err != nil {
		return
	}
	for bid, flows := range cfgs {
		exist := false
		for _, flowID := range flows {
			if flowID == opt.FlowID {
				exist = true
				break
			}
		}
		if exist {
			mngid = bid
			break
		}
	}
	log.Info("GetRole %d %d", mngid, opt.UID)

	if mngid <= 0 {
		return
	}

	roles, err := s.http.GetRole(c, mngid, opt.UID)
	if err != nil {
		return
	}

	for _, item := range roles {
		log.Info("GetRole %d %d %+v", mngid, opt.UID, item)
		if int8(item.RID) == task.TaskRoleLeader {
			role = task.TaskRoleLeader
			break
		}
		if int8(item.RID) == task.TaskRoleMember {
			role = task.TaskRoleMember
		}
	}
	log.Info("GetRole %d %d %d", mngid, opt.UID, role)

	if opt.Uname == "" {
		unames, _ := s.http.GetUnames(c, []int64{opt.UID})
		opt.Uname = unames[opt.UID]
	}

	uname = opt.Uname
	s.mc.SetRole(c, opt, role)
	return
}

//GetTaskBizFlows uid能查看哪些任务节点
func (s *Service) GetTaskBizFlows(c context.Context, uid int64) (businessID []int64, flows []int64, err error) {
	var (
		uroles []*task.Role
	)
	businessID = []int64{}
	flows = []int64{}
	//用户角色
	if uroles, err = s.http.GetUserRoles(c, uid); err != nil {
		log.Error("GetTaskBizFlows s.http.GetUserRoles(%d) error(%v)", uid, err)
		return
	}

	//业务与用户角色的映射关系
	bizMap := map[int64]int{}
	bizs := []int64{}
	for _, item := range uroles {
		if item == nil || item.BID <= 0 {
			continue
		}

		bizMap[item.BID] = 1
		bizs = append(bizs, item.BID)
	}

	businessID, flows = s.getTaskBiz(c, bizs)
	log.Info("checkAccessTask uid(%d) can see business(%+v)", uid, businessID)
	return
}

//GetRoleBiz uid能查看哪些业务
func (s *Service) GetRoleBiz(c context.Context, uid int64, role string, noAdmin bool) (businessID []int64, err error) {
	var (
		uroles []*task.Role
	)
	businessID = []int64{}
	//用户角色
	if uroles, err = s.http.GetUserRoles(c, uid); err != nil {
		log.Error("GetRoleBiz s.http.GetUserRoles(%d) error(%v)", uid, err)
		return
	}

	//业务与用户角色的映射关系
	bizMap := map[int64]int{}
	for _, item := range uroles {
		if item == nil || item.BID <= 0 || item.RID <= 0 {
			continue
		}

		roles, bizID := s.getBizRole(c, item.BID)
		log.Info("GetRoleBiz s.getBizRole roles(%+v) bizID(%d)", roles, bizID)
		if len(roles) == 0 || bizID <= 0 {
			continue
		}

		//没有role, 走bid, 有role，走role过滤; 然后在走noadmin过滤
		if (role != "" && roles[role] != item.RID) || (noAdmin && roles[business.BizBIDAdmin] == item.RID) || bizMap[bizID] > 0 {
			continue
		}

		businessID = append(businessID, bizID)
		bizMap[bizID] = 1
	}
	log.Info("uid(%d) can see biz(%v) as role(%s) noadmin(%v)", uid, businessID, role, noAdmin)
	return
}

//FlushRole .
func (s *Service) FlushRole(c context.Context, BizID, flowID int64, uids []int64) (err error) {
	return s.mc.DelRole(c, BizID, flowID, uids)
}
