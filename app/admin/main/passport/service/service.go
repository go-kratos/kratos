package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"

	"go-common/app/admin/main/passport/conf"
	"go-common/app/admin/main/passport/dao"
)

// Service struct
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	AESBlock cipher.Block
	hashSalt []byte
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		hashSalt: []byte(c.Encode.Salt),
	}
	s.AESBlock, _ = aes.NewCipher([]byte(c.Encode.AesKey))
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
