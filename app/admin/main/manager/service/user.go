package service

import (
	"context"
	"time"

	"go-common/app/admin/main/manager/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Auth return user's auth infomation.
func (s *Service) Auth(c context.Context, username string) (res *model.Auth, err error) {
	var (
		user model.User
		resp *model.RespPerm
	)
	if err = s.dao.DB().Where("username = ?", username).First(&user).Error; err != nil {
		if err == ecode.NothingFound {
			err = ecode.Int(10001)
			return
		}
		err = errors.Wrapf(err, "s.dao.DB().user.fitst(%s)", username)
		return
	}
	res = &model.Auth{
		UID:      user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}
	if resp, err = s.permsByMid(c, user.ID); err != nil {
		err = errors.Wrapf(err, "s.permsByMid(%d)", user.ID)
		return
	}
	res.Perms = resp.Res
	res.Admin = resp.Admin
	var aa []*model.AuthAssign
	if result := s.dao.DB().Joins("left join auth_item on auth_item.id=auth_assignment.item_id").Where("auth_assignment.user_id=? and auth_item.type=?", user.ID, model.TypeGroup).Find(&aa); !result.RecordNotFound() {
		res.Assignable = true
	}
	return
}

// Permissions return user's permissions.
func (s *Service) Permissions(c context.Context, username string) (res *model.Permissions, err error) {
	res = new(model.Permissions)
	var (
		user model.User
		resp *model.RespPerm
	)
	if err = s.dao.DB().Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ecode.NothingFound
		}
		log.Error("find user username(%s) error(%v)", username, err)
		return
	}
	res.UID = user.ID
	if resp, err = s.permsByMid(c, user.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ecode.NothingFound
		}
		log.Error("s.permsByMid(%d) error(%v)", user.ID, err)
		return
	}
	res.Perms = resp.Res
	res.Admin = resp.Admin
	res.Orgs = resp.Groups
	res.Roles = resp.Roles
	return
}

func (s *Service) permsByMid(c context.Context, mid int64) (resp *model.RespPerm, err error) {
	var (
		pointIDs []int64
		as       []*model.AuthAssign
	)
	resp = &model.RespPerm{}
	resp.Res = make([]string, 0)
	resp.Groups = make([]*model.AuthOrg, 0)
	resp.Roles = make([]*model.AuthOrg, 0)
	if s.admins[mid] {
		resp.Admin = true
		for _, p := range s.pointList {
			resp.Res = append(resp.Res, p.Data)
		}
		return
	}
	if err = s.dao.DB().Where("user_id = ?", mid).Find(&as).Error; err != nil {
		err = errors.Wrapf(err, "s.dao.DB().user.find(%d)", mid)
		return
	}
	for _, a := range as {
		if info, found := s.orgAuth[a.ItemID]; found { // group & role information
			if info.Type == model.TypeRole {
				resp.Roles = append(resp.Roles, info)
			}
			if info.Type == model.TypeGroup {
				resp.Groups = append(resp.Groups, info)
			}
		}
		if vs, ok := s.roleAuth[a.ItemID]; ok {
			pointIDs = append(pointIDs, vs...)
			continue
		}
		if vs, ok := s.groupAuth[a.ItemID]; ok {
			pointIDs = append(pointIDs, vs...)
		}
	}
	repeat := map[string]byte{}
	for _, id := range pointIDs {
		if assign, ok := s.points[id]; ok {
			l := len(repeat)
			repeat[assign.Data] = 0
			if len(repeat) != l {
				resp.Res = append(resp.Res, assign.Data)
			}
		}
	}
	return
}

// Users get user list.
func (s *Service) Users(c context.Context, pn, ps int) (res *model.UserPager, err error) {
	res = &model.UserPager{
		Pn: pn,
		Ps: ps,
	}
	var items []*model.User
	if err = s.dao.DB().Where("state = ?", model.UserStateOn).Offset((pn - 1) * ps).Limit(ps).Find(&items).Error; err != nil {
		if err != ecode.NothingFound {
			err = errors.Wrapf(err, "s.dao.DB().users.find(%d,%d)", pn, ps)
			return
		}
		err = nil
	}
	res.Items = items
	return
}

// UsersTotal get user total
func (s *Service) UsersTotal(c context.Context) (total int64, err error) {
	var item *model.User
	if err = s.dao.DB().Model(&item).Where("state = ?", model.UserStateOn).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "s.dao.DB().users.count")
	}
	return
}

// Heartbeat user activity record
func (s *Service) Heartbeat(c context.Context, username string) (err error) {
	var user model.User
	if err = s.dao.DB().Where("username = ?", username).First(&user).Error; err != nil {
		if err == ecode.NothingFound {
			err = ecode.Int(10001)
			return
		}
		err = errors.Wrapf(err, "s.dao.DB().user.fitst(%s)", username)
		return
	}
	now := time.Now()
	if err = s.dao.DB().Exec("insert into user_heartbeat (uid,mtime) values (?,?) on duplicate key update mtime=?", user.ID, now, now).Error; err != nil {
		err = errors.Wrapf(err, "s.dao.DB().heartbeat(%d)", user.ID)
	}
	return
}

// loadUnames loads the relation of uid & unames in two maps
func (s *Service) loadUnames() {
	var (
		err        error
		items      []*model.User
		unames     = make(map[int64]string)
		unicknames = make(map[int64]string)
		uids       = make(map[string]int64)
	)
	if err = s.dao.DB().Where("state = ?", model.UserStateOn).Find(&items).Error; err != nil {
		if err != ecode.NothingFound {
			log.Error("LoadUnames Error (%v)", err)
			return
		}
		if len(items) == 0 {
			log.Info("LoadUnames No Active User was found")
			return
		}
	}
	for _, v := range items {
		unames[v.ID] = v.Username
		unicknames[v.ID] = v.Nickname
		uids[v.Username] = v.ID
	}
	if length := len(unames); length != 0 {
		s.userNames = unames
		s.userNicknames = unicknames
		log.Info("LoadUnames Refresh Success! Lines:%d", length)
	}
	if length := len(uids); length != 0 {
		s.userIds = uids
		log.Info("LoadUIds Refresh Success! Lines:%d", length)
	}
}

// regularly load unames
func (s *Service) loadUnamesproc() {
	var duration time.Duration
	if duration = time.Duration(s.c.UnameTicker); duration == 0 {
		duration = time.Duration(5 * time.Minute) // default value
	}
	ticker := time.NewTicker(duration)
	for range ticker.C {
		s.loadUnames()
		s.loadUdepts()
	}
	ticker.Stop()
}

// Unames treats the param "uids" and give back the unames data
func (s *Service) Unames(c context.Context, uids []int64) (res map[int64]string) {
	res = make(map[int64]string)
	for _, v := range uids {
		if uname, ok := s.userNames[v]; ok {
			res[v] = uname
		}
	}
	return
}

// UIds treats the param "unames" and give back the uids data
func (s *Service) UIds(c context.Context, unames []string) (res map[string]int64) {
	res = make(map[string]int64)
	for _, v := range unames {
		if uid, ok := s.userIds[v]; ok {
			res[v] = uid
		}
	}
	return
}

// loadUdeps loads the relation of uid & department in one map
func (s *Service) loadUdepts() {
	var (
		items  []*model.UserDept
		udepts = make(map[int64]string)
	)
	/**
	select `user`.id,user_department.`name` from `user` LEFT JOIN user_department on `user`.department_id = user_department.id where user_department.`status` = 1 and `user`.state = 0
	*/
	err := s.dao.DB().Table("user").
		Select("user.id, user_department.`name` AS department").
		Joins("LEFT JOIN user_department ON `user`.department_id = user_department.id").
		Where("user_department.`status` = ?", model.UserDepOn).
		Where("state = ?", model.UserStateOn).Find(&items).Error
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("loadUdepts Error (%v)", err)
			return
		}
		if len(items) == 0 {
			log.Info("loadUdepts No Active User was found")
			return
		}
	}
	for _, v := range items {
		udepts[v.ID] = v.Department
	}
	if length := len(udepts); length != 0 {
		s.userDeps = udepts
		log.Info("loadUdepts Refresh Success! Lines:%d", length)
	}
}

// Udepts treats the param "uids" and give back the users' department data
func (s *Service) Udepts(c context.Context, uids []int64) (res map[int64]string) {
	res = make(map[int64]string)
	for _, v := range uids {
		if udept, ok := s.userDeps[v]; ok {
			res[v] = udept
		}
	}
	return
}
