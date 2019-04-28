package service

import (
	"context"
	"sort"

	"go-common/library/log"

	model "go-common/app/admin/main/macross/model/manager"
)

// Role get role.
func (s *Service) Role(c context.Context, system string) (res []*model.Role) {
	for _, role := range s.role[system] {
		if role == nil {
			continue
		}
		role.Auths = make(map[string]*model.Auth)
		if authIDs, ok := s.authRelation[role.RoleID]; ok {
			if auths, ok := s.auth[system]; ok {
				for _, authID := range authIDs {
					if auth, ok := auths[authID]; ok {
						role.Auths[auth.AuthFlag] = auth
					}
				}
			}
		}
		res = append(res, role)
	}
	sort.Sort(model.Roles(res))
	return
}

// SaveRole save role.
func (s *Service) SaveRole(c context.Context, roleID int64, system, roleName string) (err error) {
	var rows int64
	if roleID == 0 {
		if rows, err = s.dao.AddRole(c, system, roleName); err != nil {
			log.Error("s.dao.AddRole(%s, %s) error(%v)", system, roleName, err)
			return
		}
	} else {
		if rows, err = s.dao.UpRole(c, roleName, roleID); err != nil {
			log.Error("s.dao.UpRole(%s, %d) error(%v)", roleName, roleID, err)
			return
		}
	}
	if rows != 0 {
		// update cache
		s.loadRoleCache()
	}
	return
}

// DelRole del role.
func (s *Service) DelRole(c context.Context, roleID int64) (err error) {
	var rows int64
	if rows, err = s.dao.DelRole(c, roleID); err != nil {
		log.Error("s.dao.DelRole(%d) error(%s)", roleID, err)
		return
	} else if rows != 0 {
		// update cache
		s.loadRoleCache()
		if rows, err = s.dao.CleanAuthRelationByRole(c, roleID); err != nil {
			log.Error("s.dao.CleanAuthRelationByRole(%d) error(%s)", roleID, err)
			return
		} else if rows != 0 {
			s.loadAuthRelationCache()
		}
	}
	return
}
