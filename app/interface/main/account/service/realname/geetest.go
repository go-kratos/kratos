package realname

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CaptchaGTRegister register geetest
func (s *Service) CaptchaGTRegister(c context.Context, mid int64, ip, gtType string) (urlstr string, remote int, err error) {
	var (
		p         = url.Values{}
		challenge string
	)
	p.Set("ct", "1")
	p.Set("gt", conf.Conf.Realname.Geetest.CaptchaID)
	if challenge, err = s.realnameDao.RealnameCaptchaGTRegister(c, mid, ip, gtType, 1); err != nil || challenge == "" {
		p.Set("success", "0")
		remote = 0
		err = nil
		randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		p.Set("challenge", hex.EncodeToString(randOne[:])+hex.EncodeToString(randTwo[:])[0:2])
	} else {
		p.Set("success", "1")
		remote = 1
		array := md5.Sum([]byte(challenge + conf.Conf.Realname.Geetest.PrivateKey))
		p.Set("challenge", hex.EncodeToString(array[:]))
	}
	p.Set("hash", fmt.Sprintf("%x", md5.Sum([]byte(p.Encode()))))
	urlstr = "http://passport.bilibili.com/register/verification.html?" + p.Encode()
	return
}

// CaptchaGTRefresh refresh geetest
func (s *Service) CaptchaGTRefresh(c context.Context, mid int64, ip, gtType string, hash string) (challenge string, gt string, success int, err error) {
	log.Info("CaptchaGTRefresh got hash : %s", hash)
	gt = conf.Conf.Realname.Geetest.CaptchaID
	if challenge, err = s.realnameDao.RealnameCaptchaGTRegister(c, mid, ip, gtType, 1); err != nil || challenge == "" {
		success = 0
		randOne := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		randTwo := md5.Sum([]byte(strconv.Itoa(rand.Intn(100))))
		challenge = hex.EncodeToString(randOne[:]) + hex.EncodeToString(randTwo[:])[0:2]
	} else {
		success = 1
		array := md5.Sum([]byte(challenge + conf.Conf.Realname.Geetest.PrivateKey))
		challenge = hex.EncodeToString(array[:])
	}
	return
}

// CaptchaGTValidate validate geetest
func (s *Service) CaptchaGTValidate(c context.Context, mid int64, ip, clientType string, param *model.ParamRealnameCaptchaGTCheck) (res *model.RealnameCaptchaGTValidate, err error) {
	res = &model.RealnameCaptchaGTValidate{
		State: model.RealnameFalse,
	}
	switch param.Remote {
	// 极验远程校验
	case model.RealnameTrue:
		md5Array := md5.Sum([]byte(conf.Conf.Realname.Geetest.PrivateKey + "geetest" + param.Challenge))
		if hex.EncodeToString(md5Array[:]) != param.Validate {
			return
		}
		var remoteSeccode string
		if remoteSeccode, err = s.realnameDao.RealnameCaptchaGTRegisterValidate(c, param.Challenge, param.Seccode, clientType, ip, conf.Conf.Realname.Geetest.CaptchaID, mid); err != nil {
			return
		}
		log.Info("CaptchaGTValidate remoteSec : %s", remoteSeccode)
		md5Array = md5.Sum([]byte(remoteSeccode))
		res.State = model.RealnameTrue
		// if hex.EncodeToString(md5Array[:]) == param.Seccode {
		// 	res.State = model.RealnameTrue
		// }
	// 极验本地校验
	case model.RealnameFalse:
		// md5Array := md5.Sum([]byte(param.Challenge))
		// if hex.EncodeToString(md5Array[:]) == param.Validate {
		// 	res.State = model.RealnameTrue
		// }
		res.State = model.RealnameTrue
	default:
		err = ecode.RequestErr
	}
	s.addmiss(func() {
		if missErr := s.setAlipayAntispamPassFlag(context.Background(), mid, true); missErr != nil {
			log.Error("%+v", err)
		}
	})
	return
}
