package account

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/up"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

//Service struct.
type Service struct {
	c                           *conf.Config
	acc                         *account.Dao
	article                     *article.Dao
	archive                     *archive.Dao
	up                          *up.Dao
	exemptIDCheckUps            map[int64]int64
	exemptZeroLevelAndAnswerUps map[int64]int64
	uploadTopSizeUps            map[int64]int64
	missch                      chan func()
	pCacheHit                   *prom.Prom
	pCacheMiss                  *prom.Prom
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:          c,
		acc:        rpcdaos.Acc,
		article:    rpcdaos.Art,
		archive:    rpcdaos.Arc,
		up:         rpcdaos.Up,
		missch:     make(chan func(), 1024),
		pCacheHit:  prom.CacheHit,
		pCacheMiss: prom.CacheMiss,
	}
	s.loadExemptIDCheckUps()
	s.loadExemptZeroLevelAndAnswer()
	s.loadUpSpecialUploadTopSize()
	go s.loadproc()
	go s.cacheproc()
	return s
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.loadExemptIDCheckUps()
		s.loadExemptZeroLevelAndAnswer()
		s.loadUpSpecialUploadTopSize()
	}
}

// loadExemptIDCheckUps 白名单人员，免实名认证 type=8.
func (s *Service) loadExemptIDCheckUps() {
	ups, err := s.up.UpSpecial(context.TODO(), 8)
	if err != nil {
		return
	}
	s.exemptIDCheckUps = ups
	log.Info("exemptIDCheckUps ups: (%v)", len(ups))
}

// loadExemptIDCheckUps  白名单人员，免账号激活，免账号升级到1级 type=12
func (s *Service) loadExemptZeroLevelAndAnswer() {
	ups, err := s.up.UpSpecial(context.TODO(), 12)
	if err != nil {
		return
	}
	s.exemptZeroLevelAndAnswerUps = ups
}

// loadUpSpecialUploadTopSize 投稿名单限制升级 上传视频可以超过4G,但是在8G以下 type=16
func (s *Service) loadUpSpecialUploadTopSize() {
	ups, err := s.up.UpSpecial(context.TODO(), 16)
	if err != nil {
		return
	}
	s.uploadTopSizeUps = ups
}

// AddCache add to chan for cache
func (s *Service) addCache(f func()) {
	select {
	case s.missch <- f:
	default:
		log.Warn("cacheproc chan full")
	}
}

// cacheproc is a routine for execute closure.
func (s *Service) cacheproc() {
	for {
		f := <-s.missch
		f()
	}
}
