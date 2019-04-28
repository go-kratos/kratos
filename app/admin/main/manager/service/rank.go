package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/manager/model"
	"go-common/library/log"
)

// RankGroups gets archive groups.
func (s *Service) RankGroups(c context.Context, pn, ps int) (res []*model.RankGroup, total int, err error) {
	start := (pn - 1) * ps
	if err = s.dao.DB().Where("isdel=0").Offset(start).Limit(ps).Find(&res).Error; err != nil {
		log.Error("s.RankGroups() get groups error(%v)", err)
		return
	}
	if err = s.dao.DB().Model(&model.RankGroup{}).Where("isdel=0").Count(&total).Error; err != nil {
		log.Error("s.RankGroups() get total error(%v)", err)
		return
	}
	// todo 优化，合并查询
	for _, g := range res {
		g.Auths, _ = s.rankAuths(c, g.ID)
	}
	return
}

// RankGroup gets rank group info.
func (s *Service) RankGroup(c context.Context, id int64) (res *model.RankGroup, err error) {
	res = new(model.RankGroup)
	if err = s.dao.DB().Where("id=? and isdel=0", id).First(res).Error; err != nil {
		log.Error("s.RankGroup(%d) error(%v)", id, err)
		return
	}
	res.Auths, _ = s.rankAuths(c, id)
	return
}

// AddRankGroup adds rank group.
func (s *Service) AddRankGroup(c context.Context, g *model.RankGroup, auths []int64) (id int64, err error) {
	if err = s.dao.DB().Create(g).Error; err != nil {
		log.Error("s.AddRankGroup(%+v) error(%v)", g, err)
		return
	}
	id = g.ID
	for _, a := range auths {
		if err = s.dao.DB().Exec("insert into rank_auths (group_id,auth_id,isdel) values (?,?,0)", id, a).Error; err != nil {
			log.Error("s.AddRankGroup(%d,%d) add auth error(%v)", id, a, err)
			return
		}
	}
	return
}

// UpdateRankGroup udpates rank group.
func (s *Service) UpdateRankGroup(c context.Context, g *model.RankGroup, auths []int64) (err error) {
	if err = s.dao.DB().Model(g).Update(g).Error; err != nil {
		log.Error("s.UpdateRankGroup(%+v) error(%v)", g, err)
		return
	}
	// todo 优化，找出差异auth，多了删，少了加
	if err = s.delRankAuthAll(c, g.ID); err != nil {
		log.Error("s.delRankAuthAll(%d) error(%d)", g.ID, err)
		return
	}
	for _, a := range auths {
		if err = s.dao.DB().Exec("insert into rank_auths (group_id,auth_id,isdel) values (?,?,0) on duplicate key update isdel=0", g.ID, a).Error; err != nil {
			log.Error("s.UpdateRankGroup(%d,%d) add auth error(%v)", g.ID, a, err)
			return
		}
	}
	return
}

// DelRankGroup deletes rank group...
func (s *Service) DelRankGroup(c context.Context, id int64) (err error) {
	if err = s.delRankAuthAll(c, id); err != nil {
		log.Error("s.delRankAuthAll(%d) error(%v)", id, err)
		return
	}
	if err = s.dao.DB().Model(&model.RankGroup{ID: id}).Update("isdel", 1).Error; err != nil {
		log.Error("s.DelRankGroup(%d) error(%v)", id, err)
	}
	return
}

// AddRankUser adds rank user.
func (s *Service) AddRankUser(c context.Context, uid int64) (err error) {
	if err = s.dao.DB().Create(&model.RankUser{UID: uid}).Error; err != nil {
		log.Error("s.AddRankUser(%d) error(%v)", uid, err)
	}
	return
}

// RankUsers gets rank user list.
func (s *Service) RankUsers(c context.Context, pn, ps int, username string) (res []*model.RankUserScores, total int, err error) {
	var (
		tmpUids, uids, ids []int64
		users              []*model.RankUser
	)
	if username != "" {
		var (
			us        []*model.User
			condition = fmt.Sprintf("%%%s%%", username)
		)
		if err = s.dao.DB().Where("username like ?", condition).Find(&us).Error; err != nil {
			log.Error("s.RankUsers search by username(%s) error(%v)", username, err)
			return
		}
		for _, u := range us {
			ids = append(ids, u.ID)
		}
	}
	if len(ids) > 0 {
		err = s.dao.DB().Where("uid in(?) and isdel=0", ids).Find(&users).Error
	} else {
		err = s.dao.DB().Where("isdel=0").Find(&users).Error
	}
	if err != nil {
		log.Error("s.RankUsers() get groups error(%v)", err)
		return
	}
	us := make(map[int64]map[int64]int)
	for _, u := range users {
		tmpUids = append(tmpUids, u.UID)
		if us[u.UID] == nil {
			us[u.UID] = make(map[int64]int)
		}
		if u.GroupID == 0 {
			continue
		}
		us[u.UID][u.GroupID] = u.Rank
	}
	dic := make(map[int64]bool)
	for _, i := range tmpUids {
		if !dic[i] {
			dic[i] = true
			uids = append(uids, i)
		}
	}
	names, err := s.usernames(uids)
	if err != nil {
		return
	}
	for _, uid := range uids {
		if names[uid] == nil {
			continue
		}
		u := &model.RankUserScores{
			UID:      uid,
			Username: names[uid].Username,
			Nickname: names[uid].Nickname,
			Ranks:    us[uid],
		}
		res = append(res, u)
	}
	total = len(res)
	start := (pn - 1) * ps
	if start >= total {
		res = []*model.RankUserScores{}
		return
	}
	end := start + ps
	if end > total {
		end = total
	}
	res = res[start:end]
	return
}

func (s *Service) usernames(uids []int64) (res map[int64]*model.User, err error) {
	var users []*model.User
	if err = s.dao.DB().Model(&model.User{}).Where("id in (?)", uids).Find(&users).Error; err != nil {
		log.Error("s.username(%v) error(%v)", uids, err)
		return
	}
	res = make(map[int64]*model.User, len(users))
	for _, u := range users {
		res[u.ID] = u
	}
	return
}

// SaveRankUser saves user group's rank.
func (s *Service) SaveRankUser(c context.Context, uid int64, ranks map[int64]int) (err error) {
	for gid, rank := range ranks {
		if err = s.dao.DB().Exec("insert into rank_users (group_id,uid,rank) values (?,?,?) on duplicate key update rank=?,isdel=0", gid, uid, rank, rank).Error; err != nil {
			log.Error("s.SaveRankUser(%d,%d,%d) save user rank error(%v)", gid, uid, rank, err)
			return
		}
	}
	return
}

// DelRankUser deletes user.
func (s *Service) DelRankUser(c context.Context, uid int64) (err error) {
	if err = s.dao.DB().Model(&model.RankUser{}).Where("uid=?", uid).Update("isdel", 1).Error; err != nil {
		log.Error("s.DelRankUser(%d) error(%v)", uid, err)
	}
	return
}

// delRankAuthAll deletes all rank auths by group.
func (s *Service) delRankAuthAll(c context.Context, gid int64) (err error) {
	if err = s.dao.DB().Model(&model.RankAuth{}).Where("group_id=?", gid).Update("isdel", 1).Error; err != nil {
		log.Error("s.DelRankAuthAll(%d) error(%v)", gid, err)
	}
	return
}

// rankAuths gets auths by group.
func (s *Service) rankAuths(c context.Context, gid int64) (res []*model.AuthItem, err error) {
	if err = s.dao.DB().Model(&model.RankAuth{}).Where("group_id=? and isdel=0", gid).Select("auth_item.id, auth_item.name, auth_item.data").Joins("left join auth_item on auth_item.id=rank_auths.auth_id").Scan(&res).Error; err != nil {
		log.Error("s.RankAuths(%d) error(%v)", gid, err)
	}
	return
}
