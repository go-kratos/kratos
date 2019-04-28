package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"go-common/app/interface/main/videoup/model/archive"
	gmdl "go-common/app/interface/main/videoup/model/geetest"
	"go-common/library/log"
)

// Validate fn
func (s *Service) Validate(c context.Context, geetest *archive.Geetest, upFormStr string, mid int64) (stat bool, err error) {
	if geetest == nil {
		log.Error("arc param geetest is nil, mid(%d)", mid)
		return
	}
	var res *gmdl.ValidateRes
	validate := geetest.Validate
	seccode := geetest.Seccode
	challenge := geetest.Challenge
	if len(validate) != 32 {
		log.Error("s.Validate(%s,%s,%s,%d) err(validate not eq 32byte)", challenge, validate, seccode, mid)
		return
	}
	if geetest.Success != 1 {
		slice := md5.Sum([]byte(challenge))
		stat = hex.EncodeToString(slice[:]) == validate
		return
	}
	slice := md5.Sum([]byte(s.c.Geetest.PrivateKEY + "geetest" + challenge))
	if hex.EncodeToString(slice[:]) != validate {
		log.Error("s.Validate(%s,%s,%s,%d) err(challenge not found)", challenge, validate, seccode, mid)
		return
	}
	res, err = s.gt.Validate(c, challenge, seccode, upFormStr, s.c.Geetest.CaptchaID, mid)
	if err != nil {
		log.Error("s.Validate(%s,%s,%s,%d) err(gtServer validate failed.)", challenge, validate, seccode, mid)
		return
	}
	slice = md5.Sum([]byte(seccode))
	stat = hex.EncodeToString(slice[:]) == res.Seccode
	return
}
