package service

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddBusiness .
func (s *Service) AddBusiness(c context.Context, b *model.Business) (err error) {
	if err = s.dao.AddBusiness(c, b); err != nil {
		log.Error("s.dao.AddBusiness error(%v)", err)
	}
	return
}

// UpdateBusiness .
func (s *Service) UpdateBusiness(c context.Context, b *model.Business) (err error) {
	flowState := []int64{}
	binStr := strconv.FormatInt(b.Flow, 2)
	for k, v := range []byte(binStr) {
		if string(v) == "1" {
			flowState = append(flowState, int64(len(binStr)-k))
		}
	}
	if err = s.dao.BatchUpdateChildState(c, b.ID, flowState); err != nil {
		log.Error("s.dao.BatchUpdateChildState error(%v)", err)
		return
	}
	if err = s.dao.UpdateBusiness(c, b); err != nil {
		log.Error("s.dao.UpdateBusiness error(%v)", err)
	}
	return
}

// AddRole .
func (s *Service) AddRole(c context.Context, br *model.BusinessRole) (err error) {
	var (
		maxRid    interface{}
		maxRidInt int64
	)
	if maxRid, err = s.dao.MaxRidByBid(c, br.BID); err != nil {
		return
	}
	if maxRid != nil {
		maxRidInt = maxRid.(int64)
	}
	maxRidInt++
	br.RID = maxRidInt
	if err = s.dao.AddRole(c, br); err != nil {
		log.Error("s.dao.AddRole error(%v)", err)
	}
	return
}

// UpdateRole .
func (s *Service) UpdateRole(c context.Context, br *model.BusinessRole) (err error) {
	if err = s.dao.UpdateRole(c, br); err != nil {
		log.Error("s.dao.UpdateRole error(%v)", err)
	}
	return
}

// AddUser .
func (s *Service) AddUser(c context.Context, bur *model.BusinessUserRole) (err error) {
	for _, uid := range bur.UIDs {
		if _, ok := s.userNames[uid]; !ok {
			err = ecode.ManagerUIDNOTExist
			return
		}
	}
	if err = s.dao.AddUser(c, bur); err != nil {
		log.Error("s.dao.AddUser error(%v)", err)
	}
	return
}

// UpdateUser .
func (s *Service) UpdateUser(c context.Context, bur *model.BusinessUserRole) (err error) {
	if err = s.dao.UpdateUser(c, bur); err != nil {
		log.Error("s.dao.Update error(%v)", err)
	}
	return
}

// UpdateState .
func (s *Service) UpdateState(c context.Context, su *model.StateUpdate) (err error) {
	var (
		childInfo  *model.BusinessList
		parentInfo *model.BusinessList
	)
	if su.Type == model.BusinessOpenType { // business
		if su.State == model.BusinessOpenState {
			if childInfo, err = s.dao.BusinessByID(c, su.ID); err != nil {
				log.Error("s.dao.BusinessByID param id (%d) error(%v)", su.ID, err)
				return
			}
			if childInfo.PID != 0 {
				if parentInfo, err = s.dao.BusinessByID(c, childInfo.PID); err != nil {
					log.Error("s.dao.BusinessByID param id (%d) error(%v)", childInfo, err)
					return
				}
				flowState := map[int64]int64{}
				binStr := strconv.FormatInt(parentInfo.Flow, 2)
				for k, v := range []byte(binStr) {
					if string(v) == "1" {
						flowState[int64(len(binStr)-k)] = int64(len(binStr) - k)
					}
				}
				if _, ok := flowState[childInfo.FlowState]; !ok {
					err = ecode.ManagerFlowForbidden
					return
				}
			}
		}
		if err = s.dao.UpdateBusinessState(c, su); err != nil {
			log.Error("s.dao.UpdateBusinessState error(%v)", err)
		}
		return
	}
	// role
	if err = s.dao.UpdateBusinessRoleState(c, su); err != nil {
		log.Error("s.dao.UpdateBusinessRoleState error(%v)", err)
	}
	return
}

// BusinessList .
func (s *Service) BusinessList(c context.Context, bp *model.BusinessListParams) (res []*model.BusinessList, err error) {
	var (
		pids        []int64
		sKeys       []int
		cKeys       []int
		pidChildRes map[int64][]*model.BusinessList
		fKeys       = make(map[int64][]int)
		pidsMap     = make(map[int64]int64)
		tempRes     = make(map[int64]*model.BusinessList)
		childRes    = make(map[int64]*model.BusinessList)
	)
	flowBusinessMap := make(map[int64]map[int64][]*model.BusinessList)
	if tempRes, err = s.dao.ParentBusiness(c, bp.State); err != nil {
		log.Error("s.dao.ParentBusiness error(%v)", err)
		return
	}
	if len(tempRes) <= 0 {
		res = []*model.BusinessList{}
		return
	}
	if childRes, err = s.dao.ChildBusiness(c, bp); err != nil {
		log.Error("s.dao.ChildBusiness error(%v)", err)
		return
	}
	for _, cr := range childRes {
		cKeys = append(cKeys, int(cr.ID))
	}
	// auth verify
	if bp.Auth > 0 {
		up := &model.UserListParams{
			UID:  bp.UID,
			Role: model.UserRoleDefaultVal,
		}
		var userBusiness []*model.BusinessUserRoleList
		if userBusiness, err = s.dao.UserList(c, up); err != nil {
			log.Error("s.dao.UserList error(%v)", err)
			return
		}
		if len(userBusiness) <= 0 {
			res = []*model.BusinessList{}
			return
		}
		bidsMap := make(map[int64]struct{})
		for _, u := range userBusiness {
			bidsMap[u.BID] = struct{}{}
		}
		cKeys = []int{}
		for _, cr := range childRes {
			if _, ok := bidsMap[cr.BID]; ok {
				cKeys = append(cKeys, int(cr.ID))
			}
		}
	}
	sort.Ints(cKeys)
	for _, c := range cKeys {
		cidR := childRes[int64(c)]
		if _, ok := pidsMap[cidR.PID]; !ok {
			pidsMap[cidR.PID] = cidR.PID
			pids = append(pids, cidR.PID)
		}
		cidR.FlowChild = []*model.FlowBusiness{}
		if _, ok := flowBusinessMap[cidR.PID]; !ok {
			flowBusinessMap[cidR.PID] = make(map[int64][]*model.BusinessList)
		}
		fKeys[cidR.PID] = append(fKeys[cidR.PID], int(cidR.FlowState))
		flowBusinessMap[cidR.PID][cidR.FlowState] = append(flowBusinessMap[cidR.PID][cidR.FlowState], cidR)
	}
	if bp.Check != 0 {
		if pidChildRes, err = s.dao.ChildBusinessByPIDs(c, pids); err != nil {
			log.Error("s.dao.ChildBusinessByPIDs error(%v)", err)
			return
		}
	}
	for k, r := range tempRes {
		r.FlowChild = []*model.FlowBusiness{}
		if _, ok := pidChildRes[k]; !ok && bp.Check != 0 {
			delete(tempRes, k)
			continue
		}
		if _, ok := flowBusinessMap[k]; !ok && bp.Check != 0 {
			continue
		}
		flowsMap := make(map[int]int)
		flowsSlice := []int{}
		if flows, ok := fKeys[r.ID]; ok {
			for _, f := range flows {
				if _, ok := flowsMap[f]; ok {
					continue
				}
				flowsMap[f] = f
				flowsSlice = append(flowsSlice, f)
			}
			sort.Ints(flowsSlice)
			for _, f := range flowsSlice {
				r.FlowChild = append(r.FlowChild, &model.FlowBusiness{
					FlowState: int64(f),
					Child:     flowBusinessMap[r.ID][int64(f)],
				})
			}
		}
		sKeys = append(sKeys, int(k))
	}
	sort.Ints(sKeys)
	for _, k := range sKeys {
		if _, ok := tempRes[int64(k)]; !ok {
			continue
		}
		res = append(res, tempRes[int64(k)])
	}
	return
}

// FlowList .
func (s *Service) FlowList(c context.Context, bp *model.BusinessListParams) (res []*model.BusinessList, err error) {
	var (
		tempRes = []*model.BusinessList{}
		userRes = []*model.BusinessUserRoleList{}
	)
	res = []*model.BusinessList{}
	if tempRes, err = s.dao.FlowList(c, bp); err != nil {
		log.Error("s.dao.FlowList error(%v)", err)
		return
	}
	if bp.Auth <= 0 {
		res = tempRes
		return
	}
	up := &model.UserListParams{
		UID:  bp.UID,
		Role: model.UserRoleDefaultVal,
	}

	if userRes, err = s.dao.UserList(c, up); err != nil {
		log.Error("s.dao.UserList error(%v)", err)
		return
	}
	if len(userRes) <= 0 {
		return
	}
	tMap := make(map[int64]*model.BusinessList)
	for _, t := range tempRes {
		tMap[t.BID] = t
	}
	for _, u := range userRes {
		if u.Role == "" {
			continue
		}
		if t, ok := tMap[u.BID]; ok {
			res = append(res, t)
		}
	}
	return
}

// RoleList .
func (s *Service) RoleList(c context.Context, br *model.BusinessRole) (res []*model.BusinessRole, err error) {
	var (
		rids       = []int{}
		tempRes    = []*model.BusinessRole{}
		userRes    = []*model.BusinessUserRoleList{}
		tempResMap = make(map[int64]*model.BusinessRole)
	)
	res = []*model.BusinessRole{}
	if tempRes, err = s.dao.RoleListByBID(c, br); err != nil {
		log.Error("s.dao.RoleList error(%v)", err)
		return
	}
	for _, tr := range tempRes {
		rids = append(rids, int(tr.RID))
		tempResMap[tr.RID] = tr
	}
	if br.Auth <= 0 {
		res = tempRes
		return
	}
	up := &model.UserListParams{
		BID:  br.BID,
		UID:  br.UID,
		Role: model.UserRoleDefaultVal,
	}
	if userRes, err = s.dao.UserList(c, up); err != nil {
		log.Error("s.dao.UserList error(%v)", err)
		return
	}
	if len(userRes) <= 0 {
		return
	}
	roleStr := userRes[0].Role
	roleSlice := strings.Split(roleStr, ",")
	sort.Ints(rids)
	for _, r := range rids {
		for _, rs := range roleSlice {
			rsInt64, _ := strconv.ParseInt(rs, 10, 64)
			if int64(r) == rsInt64 {
				res = append(res, tempResMap[int64(r)])
			}
		}
	}
	return
}

// UserList .
func (s *Service) UserList(c context.Context, u *model.UserListParams) (res []*model.BusinessUserRoleList, total int64, err error) {
	if u.UName != "" {
		if uid, ok := s.userIds[u.UName]; ok {
			u.UID = uid
		}
	}
	if res, err = s.dao.UserList(c, u); err != nil {
		log.Error("s.dao.UserList error(%v)", err)
		return
	}
	rids := []int64{}
	for _, r := range res {
		if uname, ok := s.userNames[r.UID]; ok {
			r.UName = uname
		}
		if unickname, ok := s.userNicknames[r.UID]; ok {
			r.UNickname = unickname
		}
		if cuname, ok := s.userNames[r.CUID]; ok {
			r.CUName = cuname
		}
		roleSlice := strings.Split(r.Role, ",")
		for _, ridStr := range roleSlice {
			ridInt64, _ := strconv.ParseInt(ridStr, 10, 64)
			rids = append(rids, ridInt64)
		}
	}
	var roleRes []*model.BusinessRole
	if roleRes, err = s.dao.RoleListByRIDs(c, u.BID, rids); err != nil {
		log.Error("s.dao.RoleListByIDs error(%v)", err)
		return
	}
	for _, r := range res {
		r.RoleName = []string{}
		rids := strings.Split(r.Role, ",")
		for _, rid := range rids {
			ridInt64, _ := strconv.ParseInt(rid, 10, 64)
			for _, rr := range roleRes {
				if ridInt64 == rr.RID && r.BID == rr.BID {
					r.RoleName = append(r.RoleName, rr.Name)
				}
			}
		}
	}
	total = int64(len(res))
	start := (u.PN - 1) * u.PS
	if start >= total {
		res = []*model.BusinessUserRoleList{}
		return
	}
	end := start + u.PS
	if end > total {
		end = total
	}
	res = res[start:end]
	return
}

// DeleteUser .
func (s *Service) DeleteUser(c context.Context, bur *model.BusinessUserRole) (err error) {
	if err = s.dao.DeleteUser(c, bur); err != nil {
		log.Error("s.dao.DeleteUser error(%v)", err)
	}
	return
}

// UserRole .
func (s *Service) UserRole(c context.Context, brl *model.BusinessUserRoleList) (res []*model.UserRole, err error) {
	var (
		ridsInt64 []int64
		userRole  []*model.BusinessUserRoleList
		roleList  []*model.BusinessRole
		u         = &model.UserListParams{
			BID:  brl.BID,
			UID:  brl.UID,
			Role: -1,
		}
	)
	if userRole, err = s.dao.UserList(c, u); err != nil {
		log.Error("s.dao.UserList error(%v)", err)
		return
	}
	if len(userRole) <= 0 {
		res = []*model.UserRole{}
		return
	}
	rids := strings.Split(userRole[0].Role, ",")
	for _, rid := range rids {
		ridInt64, _ := strconv.ParseInt(rid, 10, 64)
		ridsInt64 = append(ridsInt64, ridInt64)
	}
	if roleList, err = s.dao.RoleListByRIDs(c, brl.BID, ridsInt64); err != nil {
		log.Error("s.dao.RoleListByRIDs error(%v)", err)
		return
	}
	for _, rid := range ridsInt64 {
		for _, rl := range roleList {
			if rid == rl.RID {
				r := &model.UserRole{
					ID:   rl.ID,
					BID:  rl.BID,
					RID:  rl.RID,
					Type: rl.Type,
					Name: rl.Name,
				}
				res = append(res, r)
			}
		}
	}
	return
}

// StateUp .
func (s *Service) StateUp(c context.Context, p *model.UserStateUp) (err error) {
	err = s.dao.DB().Table("manager_business_user_role").Where("bid = ? AND uid = ?", p.BID, p.AdminID).Update("state = ?", p.State).Error
	return
}

// UserRoles .
func (s *Service) UserRoles(c context.Context, uid int64) (res []*model.UserRole, err error) {
	var roleRes []*model.BusinessUserRoleList
	if roleRes, err = s.dao.UserRoles(c, uid); err != nil {
		log.Error("s.dao.UserRoles error(%v)", err)
		return
	}
	if len(roleRes) <= 0 {
		res = []*model.UserRole{}
		return
	}
	var roleList []*model.BusinessRole
	for _, r := range roleRes {
		bid := r.BID
		rids := []int64{}
		ridsStr := strings.Split(r.Role, ",")
		for _, r := range ridsStr {
			rid, _ := strconv.ParseInt(r, 10, 64)
			rids = append(rids, rid)
		}
		var tmpRoleList []*model.BusinessRole
		if tmpRoleList, err = s.dao.RoleListByRIDs(c, bid, rids); err != nil {
			log.Error("s.dao.RoleListByRIDs error(%v)", err)
			return
		}
		roleList = append(roleList, tmpRoleList...)
	}
	for _, rl := range roleList {
		r := &model.UserRole{
			ID:   rl.ID,
			BID:  rl.BID,
			RID:  rl.RID,
			Type: rl.Type,
			Name: rl.Name,
		}
		res = append(res, r)
	}
	return
}

// AllRoles .
func (s *Service) AllRoles(c context.Context, pid, uid int64) (res []*model.UserRole, err error) {
	var childs []*model.Business
	if childs, err = s.dao.BusinessChilds(c, pid); err != nil {
		log.Error("s.dao.BusinessChilds error(%v)", err)
		return
	}
	var childIDs []int64
	for _, c := range childs {
		childIDs = append(childIDs, c.BID)
	}
	var roleRes []*model.BusinessUserRoleList
	if roleRes, err = s.dao.UserRoleByBIDs(c, uid, childIDs); err != nil {
		log.Error("s.dao.UserRoleByBIDs error(%v)", err)
		return
	}
	if len(roleRes) <= 0 {
		res = []*model.UserRole{}
		return
	}
	var roleList []*model.BusinessRole
	for _, r := range roleRes {
		bid := r.BID
		rids := []int64{}
		ridsStr := strings.Split(r.Role, ",")
		for _, r := range ridsStr {
			rid, _ := strconv.ParseInt(r, 10, 64)
			rids = append(rids, rid)
		}
		var tmpRoleList []*model.BusinessRole
		if tmpRoleList, err = s.dao.RoleListByRIDs(c, bid, rids); err != nil {
			log.Error("s.dao.RoleListByRIDs error(%v)", err)
			return
		}
		roleList = append(roleList, tmpRoleList...)
	}
	for _, rl := range roleList {
		r := &model.UserRole{
			ID:   rl.ID,
			BID:  rl.BID,
			RID:  rl.RID,
			Type: rl.Type,
			Name: rl.Name,
		}
		res = append(res, r)
	}
	return
}

// IsAdmin .
func (s *Service) IsAdmin(c context.Context, bid, uid int64) bool {
	var (
		err error
		res []*model.UserRole
	)
	p := &model.BusinessUserRoleList{
		BID: bid,
		UID: uid,
	}
	if res, err = s.UserRole(c, p); err != nil {
		log.Error("s.UserRole error(%v)", err)
		return false
	}
	for _, r := range res {
		if r.Type == 1 {
			return true
		}
	}
	return false
}
