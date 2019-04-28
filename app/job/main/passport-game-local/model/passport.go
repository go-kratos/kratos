package model

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	_cloudSalt = "bi_clould_tencent_01"
)

// AsoAccount aso account.
type AsoAccount struct {
	Mid            int64  `json:"mid"`
	Userid         string `json:"userid"`
	Uname          string `json:"uname"`
	Pwd            string `json:"pwd"`
	Salt           string `json:"salt"`
	Email          string `json:"email"`
	Tel            string `json:"tel"`
	CountryID      int64  `json:"country_id"`
	MobileVerified int8   `json:"mobile_verified"`
	Isleak         int8   `json:"isleak"`
	Mtime          string `json:"mtime"`
}

// OriginAsoAccount origin aso account.
type OriginAsoAccount struct {
	Mid            int64  `json:"mid"`
	Userid         string `json:"userid"`
	Uname          string `json:"uname"`
	Pwd            string `json:"pwd"`
	Salt           string `json:"salt"`
	Email          string `json:"email"`
	Tel            string `json:"tel"`
	CountryID      int64  `json:"country_id"`
	MobileVerified int8   `json:"mobile_verified"`
	Isleak         int8   `json:"isleak"`
	Mtime          string `json:"modify_time"`
}

// Default doHash aso account, including the followings fields: userid, uname, pwd, email, tel.
func Default(a *OriginAsoAccount) *AsoAccount {
	return &AsoAccount{
		Mid:            a.Mid,
		Userid:         a.Userid,
		Uname:          a.Uname,
		Pwd:            doHash(a.Pwd, _cloudSalt),
		Salt:           a.Salt,
		Email:          doHash(a.Email, _cloudSalt),
		Tel:            doHash(a.Tel, _cloudSalt),
		CountryID:      a.CountryID,
		MobileVerified: a.MobileVerified,
		Isleak:         a.Isleak,
		Mtime:          a.Mtime,
	}
}

// DefaultHash hash a plain text using default salt.
func DefaultHash(plaintext string) string {
	return doHash(plaintext, _cloudSalt)
}

func doHash(plaintext, salt string) string {
	if plaintext == "" {
		return ""
	}
	hash := md5.New()
	hash.Write([]byte(plaintext))
	hash.Write([]byte(salt))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
