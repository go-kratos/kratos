package service

import (
	"context"
	"sort"

	model "go-common/app/admin/main/macross/model/manager"
	"go-common/library/log"
)

// GetAuths get auths.
func (s *Service) GetAuths(c context.Context, name string) (res map[string]map[string]*model.Auth, err error) {
	res = make(map[string]map[string]*model.Auth)
	for system, users := range s.user {
		var (
			user    *model.User
			authIDs []int64
			auths   map[int64]*model.Auth
			resTmp  map[string]*model.Auth
			ok      bool
		)
		if user, ok = users[name]; !ok {
			continue
		}
		if authIDs, ok = s.authRelation[user.RoleID]; !ok {
			continue
		}
		if auths, ok = s.auth[system]; !ok {
			continue
		}
		for _, authID := range authIDs {
			if auth, ok := auths[authID]; ok {
				if resTmp, ok = res[system]; !ok {
					resTmp = make(map[string]*model.Auth)
					res[system] = resTmp
				}
				resTmp[auth.AuthFlag] = auth
			}
		}
	}
	return
}

// Auth get auth.
func (s *Service) Auth(c context.Context, system string) (res []*model.Auth) {
	for _, auth := range s.auth[system] {
		res = append(res, auth)
	}
	sort.Sort(model.Auths(res))
	return
}

// SaveAuth save auth.
func (s *Service) SaveAuth(c context.Context, authID int64, system, authName, authFlag string) (err error) {
	var rows int64
	if authID == 0 {
		if rows, err = s.dao.AddAuth(c, system, authName, authFlag); err != nil {
			log.Error("s.dao.AddAuth(%s, %s, %s) error(%v)", system, authName, authFlag, err)
			return
		}
	} else {
		if rows, err = s.dao.UpAuth(c, authName, authID); err != nil {
			log.Error("s.dao.UpAuth(%s, %d) error(%v)", authName, authID, err)
			return
		}
	}
	if rows != 0 {
		// update cache
		s.loadAuthCache()
	}
	return
}

// DelAuth del auth.
func (s *Service) DelAuth(c context.Context, authID int64) (err error) {
	var rows int64
	if rows, err = s.dao.DelAuth(c, authID); err != nil {
		log.Error("s.dao.DelAuth(%d) error(%s)", authID, err)
		return
	} else if rows != 0 {
		// update cache
		s.loadAuthCache()
		if rows, err = s.dao.CleanAuthRelationByAuth(c, authID); err != nil {
			log.Error("s.dao.CleanAuthRelationByAuth(%d) error(%v)", authID, err)
			return
		} else if rows != 0 {
			s.loadAuthRelationCache()
		}
	}
	return
}

// AuthRelation get auth relation.
func (s *Service) AuthRelation(c context.Context, roleID, authID int64, state int) (err error) {
	var rows int64
	if state == 0 {
		if rows, err = s.dao.DelAuthRelation(c, roleID, authID); err != nil {
			log.Error("s.dao.DelAuthRelation(%d, %d) error(%v)", roleID, authID, err)
			return
		}
	} else {
		if rows, err = s.dao.AddAuthRelation(c, roleID, authID); err != nil {
			log.Error("s.dao.AddAuthRelation(%d, %d) error(%v)", roleID, authID, err)
			return
		}
	}
	if rows != 0 {
		s.loadAuthRelationCache()
	}
	return
}
