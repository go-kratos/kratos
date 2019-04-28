package service

import (
	"context"

	"go-common/app/interface/main/passport-login/conf"
	"go-common/app/interface/main/passport-login/dao"
	"go-common/library/log"
)

const (
	_publicKeyType     int8  = 0
	_privateKeyType    int8  = 1
	_aesKeyType        int8  = 2
	_securitySaltType  int8  = 3
	_rsaTimeoutSeconds int64 = 20
	_tsHashLen               = 16
)

var (
	privateKey   = ""
	publicKey    = ""
	aesKey       = ""
	securitySalt = ""
)

// Service struct
type Service struct {
	c      *conf.Config
	dao    *dao.Dao
	missch chan func()
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:      c,
		dao:    dao.New(c),
		missch: make(chan func(), 10240),
	}
	s.initSecret()
	go s.cacheproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) initSecret() {
	res, err := s.dao.LoadSecret(context.Background())
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		if _publicKeyType == r.KeyType {
			publicKey = r.Key
		}
		if _privateKeyType == r.KeyType {
			privateKey = r.Key
		}
		if _aesKeyType == r.KeyType {
			aesKey = r.Key
		}
		if _securitySaltType == r.KeyType {
			securitySalt = r.Key
		}
	}
	if len(publicKey) == 0 || len(privateKey) == 0 || len(aesKey) == 0 || len(securitySalt) == 0 {
		panic("load secret error")
	}
}

func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for executing closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
