package service

import (
	"context"
	"image/color"
	"sync"
	"time"

	"go-common/app/interface/main/captcha/conf"
	"go-common/app/interface/main/captcha/dao"
	"go-common/library/cache"
	"go-common/library/ecode"

	"github.com/golang/freetype/truetype"
)

// Captcha captcha.
type Captcha struct {
	frontColors  []color.Color
	bkgColors    []color.Color
	disturbLevel int
	fonts        []*truetype.Font
}

// Service captcha service.
type Service struct {
	conf    *conf.Config
	dao     *dao.Dao
	captcha *Captcha
	cacheCh *cache.Cache
	// captcha mem.
	init   bool
	lock   sync.RWMutex
	mCode  map[string][]string
	mImage map[string]map[string][]byte
}

// New new a service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		conf:    c,
		dao:     dao.New(c),
		captcha: newCaptcha(c.Captcha),
		cacheCh: cache.New(1, 1024),
		mCode:   make(map[string][]string, len(c.Business)),
		mImage:  make(map[string]map[string][]byte, len(c.Business)),
	}
	go s.generaterProc()
	return s
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}

// Ping ping dao.
func (s *Service) Ping(c context.Context) error {
	if !s.init {
		return ecode.CaptchaNotCreate
	}
	return s.dao.Ping(c)
}

func (s *Service) generaterProc() {
	waiter := &sync.WaitGroup{}
	for _, b := range s.conf.Business {
		waiter.Add(1)
		go s.initGenerater(waiter, b.BusinessID, b.LenStart, b.LenEnd, b.Width, b.Length)
	}
	waiter.Wait()
	s.init = true
	for {
		for _, b := range s.conf.Business {
			go s.generater(b.BusinessID, b.LenStart, b.LenEnd, b.Width, b.Length)
		}
		time.Sleep(time.Duration(s.conf.Captcha.Interval))
	}
}
