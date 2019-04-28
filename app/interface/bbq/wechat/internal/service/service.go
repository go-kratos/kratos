package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"go-common/app/interface/bbq/wechat/dao"
	"go-common/app/interface/bbq/wechat/internal/conf"
	"go-common/app/interface/bbq/wechat/internal/model"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

// Ping Service
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// TokenGet .
func (s *Service) TokenGet(ctx context.Context, arg *model.TokenReq, url string) (sig string, err error) {
	token := ""

	if token, err = s.dao.TokenGetLast(ctx); err != nil {
		log.Warnv(ctx, log.KV("log", "weixin token get last fail"))
		return
	}

	if token == "" {
		if token, err = s.dao.TokenGet(ctx); err != nil {
			log.Warnv(ctx, log.KV("log", "weixin token get fail"))
			return
		}
	}

	sig = s.TokenCalc(ctx, token, arg.Noncestr, arg.Timestamp, url)
	return
}

// TokenCalc 计算签名
func (s *Service) TokenCalc(ctx context.Context, ticket string, noncestr string, timestamp string, url string) (sig string) {
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, noncestr, timestamp, url)

	hash := sha1.New()
	hash.Write([]byte(str))
	bs := hash.Sum(nil)

	sig = fmt.Sprintf("%x", bs)
	return
}
