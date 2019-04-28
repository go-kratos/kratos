package service

import (
	"context"
	"io"
	"time"

	"go-common/library/log"

	"go-common/app/admin/main/macross/conf"
	"go-common/app/admin/main/macross/dao"
	"go-common/app/admin/main/macross/dao/oss"
	model "go-common/app/admin/main/macross/model/manager"
)

// Service service struct info.
type Service struct {
	c   *conf.Config
	oss *oss.Dao
	dao *dao.Dao
	// manager cache
	user         map[string]map[string]*model.User // system => { username => managerInfo }
	role         map[string]map[int64]*model.Role  // system => { roleId => roleInfo }
	authRelation map[int64][]int64                 // role_id => [ auth_id ]
	auth         map[string]map[int64]*model.Auth  // system => { authId => authInfo }
	// ios cache
	modelNameCache map[string]map[string]int64
}

// New service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		oss: oss.New(c),
		// init manager cache
		user:         make(map[string]map[string]*model.User),
		role:         make(map[string]map[int64]*model.Role),
		authRelation: make(map[int64][]int64),
		auth:         make(map[string]map[int64]*model.Auth),
		// init ios cache
		modelNameCache: make(map[string]map[string]int64),
	}
	// manager cache
	if err := s.loadUserCache(); err != nil {
		panic(err)
	}
	if err := s.loadRoleCache(); err != nil {
		panic(err)
	}
	if err := s.loadAuthRelationCache(); err != nil {
		panic(err)
	}
	if err := s.loadAuthCache(); err != nil {
		panic(err)
	}
	go s.loadproc()
	return
}

// loadproc is a routine load to cache
func (s *Service) loadproc() {
	for {
		time.Sleep(time.Duration(conf.Conf.Reload))
		s.loadUserCache()
		s.loadRoleCache()
		s.loadAuthRelationCache()
		s.loadAuthCache()
	}
}

func (s *Service) loadUserCache() (err error) {
	var tmpUser map[string]map[string]*model.User
	if tmpUser, err = s.dao.Users(context.TODO()); err != nil {
		log.Error("s.dao.Users() error(%v)", err)
		return
	}
	s.user = tmpUser
	return
}

func (s *Service) loadRoleCache() (err error) {
	var tmpRole map[string]map[int64]*model.Role
	if tmpRole, err = s.dao.Roles(context.TODO()); err != nil {
		log.Error("s.dao.Roles() error(%v)", err)
		return
	}
	s.role = tmpRole
	return
}

func (s *Service) loadAuthRelationCache() (err error) {
	var tmpAuthRelation map[int64][]int64
	if tmpAuthRelation, err = s.dao.AuthRelation(context.TODO()); err != nil {
		log.Error("s.dao.AuthRelation() error(%v)", err)
		return
	}
	s.authRelation = tmpAuthRelation
	return
}

func (s *Service) loadAuthCache() (err error) {
	var tmpAuth map[string]map[int64]*model.Auth
	if tmpAuth, err = s.dao.Auths(context.TODO()); err != nil {
		log.Error("s.dao.Auths() error(%v)", err)
		return
	}
	s.auth = tmpAuth
	return
}

// DiffPutOss upload diff to oss
func (s *Service) DiffPutOss(c context.Context, f io.Reader, filename string) (uri string, err error) {
	if uri, err = s.oss.Put(c, f, filename); err != nil {
		log.Error("s.oss.Put(%s) error(%v)", filename, err)
	}
	return
}

// Ping dao
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("s.dao error(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.dao.Close()
}
