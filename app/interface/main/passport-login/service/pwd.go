package service

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	tmpRsaTimeOut = ecode.New(662)
	tmpPWDError   = ecode.New(629)
)

// checkUserPwd compare pwd
func (s *Service) checkUserPwd(rsaPwd, originPwd, salt string) (err error) {
	var (
		tsHash, pwd string
	)
	if tsHash, pwd, err = parseRSAPwd(rsaPwd); err != nil {
		return err
	}

	ts, err := Hash2TsSeconds(tsHash)
	if err != nil {
		return tmpRsaTimeOut
	}
	now := time.Now().Unix()
	if now-ts > _rsaTimeoutSeconds {
		return tmpRsaTimeOut
	}

	if pwd == "" || len(pwd) == 0 {
		return tmpPWDError
	}

	if !matching(pwd, salt, originPwd) {
		return tmpPWDError
	}
	return
}

func parseRSAPwd(rsaPwd string) (tsHash, pwd string, err error) {
	if len(rsaPwd) < 88 {
		return "", "", tmpPWDError
	}
	rsaPwd = addBase64Padding(rsaPwd)
	var tsHashPwd string
	if tsHashPwd, err = rsaDecrypt(rsaPwd); err != nil {
		return "", "", err
	}
	if len(tsHashPwd) < _tsHashLen {
		return "", "", tmpPWDError
	}
	tsHash = tsHashPwd[:_tsHashLen]
	pwd = tsHashPwd[_tsHashLen:]
	return
}

func rsaDecrypt(rsaPwd string) (res string, err error) {
	var (
		rs          []byte
		priv        *rsa.PrivateKey
		decryptByte []byte
	)
	if rs, err = base64.StdEncoding.DecodeString(rsaPwd); err != nil {
		log.Error("failed to base64 StdEncoding decode RSA pwd , error(%+v),rsaPwd = (%s)", err, rsaPwd)
		if rs, err = base64.URLEncoding.DecodeString(rsaPwd); err != nil {
			log.Error("failed to base64 URLEncoding decode RSA pwd , error(%+v),rsaPwd = (%s)", err, rsaPwd)
		}
		if err != nil {
			log.Error("base64 decode all fail pwd(len:%d)(%s) error(%+v)", len(rsaPwd), rsaPwd, err)
			return
		}
	}
	if priv, err = priKey([]byte(privateKey)); err != nil {
		log.Error("priKey err, error (%+v)", err)
	}
	if decryptByte, err = rsaDecryptPKCS8(priv, rs); err != nil {
		log.Error("failed to decrypt RSA pwd , error(%v)", err)
		return
	}
	res = string(decryptByte)
	return
}

func matching(plainPwd, salt, originPwd string) bool {
	if salt == "" {
		return md5Hex(plainPwd) == originPwd
	}
	return md5Hex(fmt.Sprintf("%s>>BiLiSaLt<<%s", md5Hex(plainPwd), salt)) == originPwd
}
