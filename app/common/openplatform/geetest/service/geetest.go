package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"go-common/app/common/openplatform/geetest/dao"
	"go-common/app/common/openplatform/geetest/model"
	"math/rand"
	"strconv"
)

// PreProcess preprocessing the geetest and get to challenge
func PreProcess(c context.Context, mid int64, newCaptcha int, ip, clientType, captchaID, privateKey string) (res *model.ProcessRes, err error) {
	var pre string
	res = &model.ProcessRes{}
	res.CaptchaID = captchaID
	res.NewCaptcha = newCaptcha
	if pre, err = dao.PreProcess(c, mid, ip, clientType, newCaptcha, captchaID); err != nil || pre == "" {
		randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		challenge := hex.EncodeToString(randOne[:]) + hex.EncodeToString(randTwo[:])[0:2]
		res.Challenge = challenge
		return
	}
	res.Success = 1
	slice := md5.Sum([]byte(pre + privateKey))
	res.Challenge = hex.EncodeToString(slice[:])
	return
}

// Validate recheck the challenge code and get to seccode
func Validate(c context.Context, challenge, validate, seccode, clientType, ip, captchaID, privateKey string, success int, mid int64) (stat bool) {
	if len(validate) != 32 {
		return
	}
	if success != 1 {
		slice := md5.Sum([]byte(challenge))
		stat = hex.EncodeToString(slice[:]) == validate
		return
	}
	slice := md5.Sum([]byte(privateKey + "geetest" + challenge))
	if hex.EncodeToString(slice[:]) != validate {
		return
	}
	res, err := dao.Validate(c, challenge, seccode, clientType, ip, captchaID, mid)
	if err != nil {
		return
	}
	slice = md5.Sum([]byte(seccode))
	stat = hex.EncodeToString(slice[:]) == res.Seccode
	return
}
