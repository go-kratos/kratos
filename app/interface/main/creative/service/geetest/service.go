package geetest

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	d "go-common/app/interface/main/creative/dao/geetest"
	m "go-common/app/interface/main/creative/model/geetest"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

//Service struct.
type Service struct {
	c    *conf.Config
	gt   *d.Dao
	acc  *account.Dao
	pErr *prom.Prom
}

//New get service.
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:    c,
		gt:   d.New(c),
		acc:  rpcdaos.Acc,
		pErr: prom.BusinessErrCount,
	}
	return s
}

// PreProcessAdd fn
func (s *Service) PreProcessAdd(c context.Context, mid int64, ip, clientType string, newCaptcha int) (res *m.ProcessRes, err error) {
	var pre string
	res = &m.ProcessRes{
		Limit: map[string]bool{
			"add": false,
		},
	}
	if exist, _, _ := s.acc.HalfMin(c, mid); exist {
		log.Info("halfMin exist | mid(%d)", mid)
		res.Limit["add"] = true
	}
	if res.Limit["add"] {
		res.CaptchaID = s.c.Geetest.CaptchaID
		res.NewCaptcha = newCaptcha
		if pre, err = s.gt.PreProcess(c, mid, ip, clientType, newCaptcha); err != nil || pre == "" {
			randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
			randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
			challenge := hex.EncodeToString(randOne[:]) + hex.EncodeToString(randTwo[:])[0:2]
			res.Challenge = challenge
			log.Info("s.gt.PreProcess err != nil (%s,%d,%d) ", clientType, newCaptcha, mid)
			s.pErr.Incr("geetest_preprocess")
			res.Success = 0
			err = nil
			return
		}
		res.Success = 1
		slice := md5.Sum([]byte(pre + s.c.Geetest.PrivateKEY))
		res.Challenge = hex.EncodeToString(slice[:])
	}
	return
}

// PreProcess getGeetestChal
func (s *Service) PreProcess(c context.Context, mid int64, ip, clientType string, newCaptcha int) (res *m.ProcessRes, err error) {
	var pre string
	res = &m.ProcessRes{
		Limit: map[string]bool{
			"add": false,
		},
	}
	if exist, _, _ := s.acc.HalfMin(c, mid); exist {
		log.Info("halfMin exist | mid(%d)", mid)
		res.Limit["add"] = true
	}
	res.CaptchaID = s.c.Geetest.CaptchaID
	res.NewCaptcha = newCaptcha
	if pre, err = s.gt.PreProcess(c, mid, ip, clientType, newCaptcha); err != nil || pre == "" {
		randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		challenge := hex.EncodeToString(randOne[:]) + hex.EncodeToString(randTwo[:])[0:2]
		res.Challenge = challenge
		log.Info("s.gt.PreProcess err != nil (%s,%d,%d) ", clientType, newCaptcha, mid)
		s.pErr.Incr("geetest_preprocess")
		res.Success = 0
		err = nil
		return
	}
	res.Success = 1
	slice := md5.Sum([]byte(pre + s.c.Geetest.PrivateKEY))
	res.Challenge = hex.EncodeToString(slice[:])
	return
}

// Validate recheck the seccode
func (s *Service) Validate(c context.Context, challenge, validate, seccode, clientType, ip string, success int, mid int64) (stat bool) {
	if len(validate) != 32 {
		log.Error("s.Validate(%s,%s,%s,%d) err(validate not eq 32byte)", challenge, validate, seccode, mid)
		return
	}
	if success != 1 {
		slice := md5.Sum([]byte(challenge))
		stat = hex.EncodeToString(slice[:]) == validate
		return
	}
	slice := md5.Sum([]byte(s.c.Geetest.PrivateKEY + "geetest" + challenge))
	if hex.EncodeToString(slice[:]) != validate {
		log.Error("s.Validate(%s,%s,%s,%d) err(challenge not found)", challenge)
		return
	}
	res, err := s.gt.Validate(c, challenge, seccode, clientType, ip, s.c.Geetest.CaptchaID, mid)
	if err != nil {
		s.pErr.Incr("geetest_validate")
		log.Error("s.Validate(%s,%s,%s,%d) err(gtServer validate failed.)", challenge, validate, seccode, mid)
		return
	}
	slice = md5.Sum([]byte(seccode))
	stat = hex.EncodeToString(slice[:]) == res.Seccode
	return
}
