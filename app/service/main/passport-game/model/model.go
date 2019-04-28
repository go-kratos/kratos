package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"

	xtime "go-common/library/time"
)

const (
	// EmptyFace empty face
	EmptyFace = "http://static.hdslb.com/images/member/noface.gif"

	_cloudSalt = "bi_clould_tencent_01"
	_leak      = 1
)

// FullFace account full face.
func (m *Info) FullFace() string {
	if m.Face == "" {
		return EmptyFace
	}
	return fmt.Sprintf("http://i%d.hdslb.com%s", m.Mid%3, m.Face)
}

// Token token resp.
type Token struct {
	Mid         int64  `json:"mid"`
	AppID       int32  `json:"appid"`
	AccessToken string `json:"access_key"`
	CreateAt    int64  `json:"create_at"`
	UserID      string `json:"userid"`
	Uname       string `json:"uname"`
	Expires     int64  `json:"expires"`
	Permission  string `json:"permission"`
}

// LoginToken login token.
type LoginToken struct {
	Mid       int64  `json:"mid"`
	AccessKey string `json:"access_key"`
	Expires   int64  `json:"expires"`
}

// RSAKey rsa key.
type RSAKey struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

// RenewToken renew token resp.
type RenewToken struct {
	Expires int64 `json:"expires"`
}

// AsoAccount aso account.
type AsoAccount struct {
	Mid            int64      `json:"mid"`
	UserID         string     `json:"userid"`
	Uname          string     `json:"uname"`
	Pwd            string     `json:"pwd"`
	Salt           string     `json:"salt"`
	Email          *string    `json:"email"`
	Tel            *string    `json:"tel"`
	CountryID      int64      `json:"country_id"`
	MobileVerified int8       `json:"mobile_verified"`
	Isleak         int8       `json:"isleak"`
	Ctime          xtime.Time `json:"-"`
	Mtime          xtime.Time `json:"-"`
}

// ResRegV3 ResRegV3
type ResRegV3 struct {
	Mid int `json:"mid"`
}

// ResRegV2 ResRegV2
type ResRegV2 struct {
	Mid       int    `json:"mid"`
	AccessKey string `json:"accessKey"`
}

// ResReg ResReg
type ResReg struct {
	Mid int `json:"mid"`
}

// ResByTel ResByTel
type ResByTel struct {
	Mid       int    `json:"mid"`
	AccessKey string `json:"accessKey"`
}

// ResCaptcha ResCaptcha
type ResCaptcha struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    CaptchaData `json:"data"`
}

// CaptchaData CaptchaData
type CaptchaData struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// ArgRegV3 ArgRegV3
type ArgRegV3 struct {
	User    string `form:"user" validate:"required"`
	Pwd     string `form:"pwd" validate:"required"`
	Captcha string `form:"captcha"`
	Ctoken  string `form:"ctoken"`
}

// ArgRegV2 ArgRegV2
type ArgRegV2 struct {
	Captcha string `form:"captcha"`
	Ctoken  string `form:"ctoken"`
}

// ArgReg ArgReg
type ArgReg struct {
	Email   string `form:"email" validate:"required"`
	Userpwd string `form:"pwd" validate:"required"`
	User    string `form:"user" validate:"required"`
}

// ArgByTel ArgByTel
type ArgByTel struct {
	Tel       string `form:"tel" validate:"required"`
	Uname     string `form:"user" validate:"required"`
	Userpwd   string `form:"pwd" validate:"required"`
	CountryID string `form:"country_id" validate:"required"`
	Captcha   string `form:"captcha" validate:"required"`
}

// SendSms SendSms
type SendSms struct {
	Tel       string `form:"tel" validate:"required"`
	CountryID string `form:"country_id" validate:"required"`
	Captcha   string `form:"captcha" `
	Ctoken    string `form:"ctoken" `
	ResetPwd  bool   `form:"reset_pwd" `
}

// TdoRegV3 TdoRegV3
type TdoRegV3 struct {
	Arg    ArgRegV3
	Cookie string
	IP     string
}

// TdoRegV2 TdoRegV2
type TdoRegV2 struct {
	Arg    ArgRegV2
	Cookie string
	IP     string
}

// TdoReg TdoReg
type TdoReg struct {
	Arg    ArgReg
	Cookie string
	IP     string
}

// TdoByTel TdoByTel
type TdoByTel struct {
	Arg    ArgByTel
	Cookie string
	IP     string
}

// TdoSendSms TdoSendSms
type TdoSendSms struct {
	Arg    SendSms
	Cookie string
	IP     string
}

// Leak leak.
func (a *AsoAccount) Leak() bool {
	return a.Isleak == _leak
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

// DefaultUserID get default user id.
func DefaultUserID(mid int64) string {
	return "用户" + strconv.FormatInt(mid, 10)
}

// DefaultUname get default uname.
func DefaultUname(mid int64) string {
	return "用户" + strconv.FormatInt(mid, 10)
}
