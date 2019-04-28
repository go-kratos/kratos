package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"image/png"
	"math/rand"
	"strings"
	"sync"
	"time"

	"go-common/library/ecode"

	uuid "github.com/satori/go.uuid"
)

const (
	_captchaURL = "%s/x/v1/captcha/get?bid=%s&token=%s"
)

// Token use bid, get a token.
func (s *Service) Token(c context.Context, bid string) (url string, token string, err error) {
	token = hex.EncodeToString(uuid.NewV4().Bytes())
	business := s.LookUp(bid)
	if err = s.dao.AddTokenCache(c, token, int32(time.Duration(business.TTL)/time.Second)); err != nil {
		return
	}
	url = fmt.Sprintf(_captchaURL, s.conf.Captcha.OuterHost, bid, token)
	return
}

// CaptchaImg get a captcha by token,bid.
func (s *Service) CaptchaImg(c context.Context, token, bid string) (img []byte, err error) {
	code, img, ttl := s.randomCaptcha(bid)
	realCode, _, err := s.dao.CaptchaCache(c, token)
	if err != nil {
		return
	}
	if realCode == "" {
		err = ecode.CaptchaTokenExpired
		return
	}
	err = s.dao.UpdateTokenCache(c, token, code, ttl)
	return
}

// VerifyCaptcha verify captcha by token and code.
func (s *Service) VerifyCaptcha(c context.Context, token, code string) (err error) {
	var (
		realCode string
		isInit   bool
	)
	if realCode, isInit, err = s.dao.CaptchaCache(c, token); err != nil {
		return
	}
	if realCode == "" {
		err = ecode.CaptchaCodeNotFound
		return
	}
	if isInit {
		err = ecode.CaptchaNotCreate
		return
	}
	if ok := strings.ToLower(realCode) == strings.ToLower(code); ok {
		s.cacheCh.Save(func() {
			s.dao.DelCaptchaCache(context.Background(), token)
		})
	} else {
		err = ecode.CaptchaErr
	}
	return
}

func (s *Service) initGenerater(waiter *sync.WaitGroup, bid string, lenStart, lenEnd, width, length int) {
	s.generater(bid, lenStart, lenEnd, width, length)
	waiter.Done()
}

func (s *Service) generater(bid string, lenStart, lenEnd, width, length int) {
	images := make(map[string][]byte, s.conf.Captcha.Capacity)
	codes := make([]string, 0, s.conf.Captcha.Capacity)
	for i := 0; i < s.conf.Captcha.Capacity; i++ {
		img, code := s.captcha.createImage(lenStart, lenEnd, width, length, TypeALL)
		var b bytes.Buffer
		switch s.conf.Captcha.Ext {
		case "png":
			png.Encode(&b, img)
		case "jpeg":
			jpeg.Encode(&b, img, &jpeg.Options{Quality: 100})
		default:
			jpeg.Encode(&b, img, &jpeg.Options{Quality: 100})
		}
		images[code] = b.Bytes()
		codes = append(codes, code)
	}
	s.lock.Lock()
	s.mImage[bid] = images
	s.mCode[bid] = codes
	s.lock.Unlock()
}

func (s *Service) randomCaptcha(bid string) (code string, img []byte, ttl int32) {
	business := s.LookUp(bid)
	ttl = int32(time.Duration(business.TTL) / time.Second)
	rnd := rand.Intn(s.conf.Captcha.Capacity)
	code = s.mCode[business.BusinessID][rnd]
	img = s.mImage[business.BusinessID][code]
	return
}
