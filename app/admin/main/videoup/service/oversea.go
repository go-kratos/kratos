package service

import (
	"go-common/app/admin/main/videoup/model/oversea"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

// PolicyGroups get policy group
func (s *Service) PolicyGroups(c *bm.Context, uid, groupID int64, gType, state int8, count, page int64, order, sort string) (data *oversea.PolicyGroupData, err error) {
	groups, total, err := s.oversea.PolicyGroups(c, uid, groupID, gType, state, count, page, order, sort)
	if len(groups) != 0 {
		s.mulIDtoName(c, groups, s.mng.GetNameByUID, "UID", "UserName")
	}
	data = &oversea.PolicyGroupData{
		Items: groups,
		Pager: &oversea.Pager{
			Num:   page,
			Size:  count,
			Total: total,
		},
	}
	return
}

// ArchiveGroups get policy group by archive id
func (s *Service) ArchiveGroups(c *bm.Context, aid int64) (groups []*oversea.PolicyGroup, err error) {
	groups, err = s.oversea.ArchiveGroups(c, aid)
	return
}

// AddPolicyGroup add policy group
func (s *Service) AddPolicyGroup(c *bm.Context, group *oversea.PolicyGroup) (err error) {
	if err = s.oversea.AddPolicyGroup(c, group); err != nil {
		return
	}
	s.sendPolicyLog(c, &oversea.PolicyGroup{}, group)
	return
}

// UpdatePolicyGroup update policy group
func (s *Service) UpdatePolicyGroup(c *bm.Context, id int64, attrs map[string]interface{}) (err error) {
	var (
		oldG = &oversea.PolicyGroup{}
		newG = &oversea.PolicyGroup{}
	)
	if oldG, err = s.oversea.PolicyGroup(c, id); err != nil {
		return
	}
	if err = s.oversea.UpdatePolicyGroup(c, id, attrs); err != nil {
		return
	}
	if newG, err = s.oversea.PolicyGroup(c, id); err != nil {
		log.Error("s.oversea.PolicyGroup(%d) err(%v)", id, err)
		err = nil
	} else {
		s.sendPolicyLog(c, oldG, newG)
	}
	return
}

// UpdatePolicyGroups multi update policy group
func (s *Service) UpdatePolicyGroups(c *bm.Context, ids []int64, attrs map[string]interface{}) (err error) {
	var (
		oldGs  []*oversea.PolicyGroup
		newGs  []*oversea.PolicyGroup
		newMap = make(map[int64]*oversea.PolicyGroup)
	)
	if oldGs, err = s.oversea.PolicyGroupsByIds(c, ids); err != nil {
		return
	}
	if err = s.oversea.UpdatePolicyGroups(c, ids, attrs); err != nil {
		return
	}
	if newGs, err = s.oversea.PolicyGroupsByIds(c, ids); err != nil {
		log.Error("s.oversea.PolicyGroupsByIds(%d) err(%v)", ids, err)
		err = nil
	} else {
		for _, v := range newGs {
			newMap[v.ID] = v
		}
		for _, oldG := range oldGs {
			newG := &oversea.PolicyGroup{}
			if _, ok := newMap[oldG.ID]; ok {
				newG = newMap[oldG.ID]
			}
			s.sendPolicyLog(c, oldG, newG)
		}
	}
	return
}

// PolicyItems get polices by group id
func (s *Service) PolicyItems(c *bm.Context, gid int64) (items []*oversea.PolicyItem, err error) {
	return s.oversea.PolicyItems(c, gid)
}

// AddPolicies add policies
func (s *Service) AddPolicies(c *bm.Context, uid, gid int64, items []*oversea.PolicyParams) (err error) {
	var (
		zids []int64
		oldG = &oversea.PolicyGroup{}
		newG = &oversea.PolicyGroup{}
	)
	if oldG, err = s.oversea.PolicyGroup(c, gid); err != nil {
		return
	}
	policies := make([]oversea.PolicyItem, len(items))
	for i, v := range items {
		policies[i].ID = v.ID
		policies[i].GroupID = gid
		policies[i].PlayAuth = v.PlayAuth
		policies[i].DownAuth = v.DownAuth
		policies[i].AreaID = xstr.JoinInts(v.AreaIds)
		policies[i].State = oversea.StateOK
		zids, _ = s.oversea.ZoneIDs(c, v.AreaIds)
		policies[i].ZoneID = xstr.JoinInts(zids)
	}
	if err = s.oversea.AddPolicies(c, policies); err != nil {
		return
	}
	if err = s.oversea.UpdatePolicyGroup(c, gid, map[string]interface{}{"uid": uid}); err != nil {
		log.Error("s.oversea.UpdatePolicyGroup(%d) err(%v)", gid, err)
		err = nil
	}
	if newG, err = s.oversea.PolicyGroup(c, gid); err != nil {
		log.Error("s.oversea.PolicyGroup(%d) err(%v)", gid, err)
		err = nil
	} else {
		s.sendPolicyLog(c, oldG, newG)
	}
	return
}

// DelPolices soft delete policies
func (s *Service) DelPolices(c *bm.Context, uid, gid int64, ids []int64) (err error) {
	var (
		oldG = &oversea.PolicyGroup{}
		newG = &oversea.PolicyGroup{}
	)
	if oldG, err = s.oversea.PolicyGroup(c, gid); err != nil {
		return
	}
	if err = s.oversea.DelPolices(c, gid, ids); err != nil {
		return
	}
	if err = s.oversea.UpdatePolicyGroup(c, gid, map[string]interface{}{"uid": uid}); err != nil {
		log.Error("s.oversea.UpdatePolicyGroup(%d) err(%v)", gid, err)
		err = nil
	}
	if newG, err = s.oversea.PolicyGroup(c, gid); err != nil {
		log.Error("s.oversea.PolicyGroup(%d) err(%v)", gid, err)
		err = nil
	} else {
		s.sendPolicyLog(c, oldG, newG)
	}
	return
}
