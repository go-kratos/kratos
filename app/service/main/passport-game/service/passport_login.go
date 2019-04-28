package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/satori/go.uuid"
)

const (
	_tsHashLen         = 16
	_rsaTimeoutSeconds = 20

	_asoPasswordError = -627
)

// Login login via cloud, if err occurred, via origin.
func (s *Service) Login(c context.Context, app *model.App, subid int32, userid, rsaPwd string) (res *model.LoginToken, err error) {
	if userid == "" || rsaPwd == "" {
		err = ecode.UsernameOrPasswordErr
		return
	}
	cache := true
	a, pwd, tsHash, err := s.checkUserData(c, userid, rsaPwd)
	if err != nil {
		if err == ecode.PasswordHashExpires || err == ecode.PasswordTooLeak {
			return
		}
		err = nil
	} else {
		var t *model.Perm
		if t, cache, err = s.saveToken(c, app.AppID, subid, a.Mid); err != nil {
			err = nil
		} else {
			res = &model.LoginToken{
				Mid:       t.Mid,
				AccessKey: t.AccessToken,
				Expires:   t.Expires,
			}
			return
		}
	}
	if res, err = s.loginOrigin(c, userid, pwd, tsHash); err != nil {
		return
	}
	if cache && res != nil {
		s.d.SetTokenCache(c, &model.Perm{
			Mid:         res.Mid,
			AppID:       app.AppID,
			AppSubID:    subid,
			AccessToken: res.AccessKey,
			CreateAt:    res.Expires - _expireSeconds,
			Expires:     res.Expires,
		})
	}
	return
}

func (s *Service) checkUserData(c context.Context, userid, rsaPwd string) (a *model.AsoAccount, pwd, tsHash string, err error) {
	if tsHash, pwd, err = s.parseRSAPwd(rsaPwd); err != nil {
		return
	}
	if timeout(time.Now().Unix(), tsHash) {
		err = ecode.PasswordHashExpires
		return
	}
	if pwd == "" {
		err = ecode.UsernameOrPasswordErr
		return
	}
	if a, err = s.asoAccountInfo(c, userid); err != nil {
		return
	}
	// TODO 2017.12.27 when leak business changes, update here
	if a.Leak() {
		err = ecode.PasswordTooLeak
		return
	}
	if pwd == "" {
		err = ecode.UsernameOrPasswordErr
		return
	}
	if !pwdMatches(pwd, a.Salt, a.Pwd) {
		err = ecode.UsernameOrPasswordErr
	}
	return
}

func (s *Service) parseRSAPwd(rsaPwd string) (tsHash, pwd string, err error) {
	if len(rsaPwd) < 88 {
		err = ecode.UsernameOrPasswordErr
		return
	}
	var tsHashPwd string
	if tsHashPwd, err = s.rsaDecrypt(rsaPwd); err != nil {
		return
	}
	if len(tsHashPwd) < _tsHashLen {
		err = ecode.UsernameOrPasswordErr
		return
	}
	tsHash = tsHashPwd[:_tsHashLen]
	pwd = tsHashPwd[_tsHashLen:]
	return
}

// asoAccountInfo get user info by userid or hash.
func (s *Service) asoAccountInfo(c context.Context, userid string) (info *model.AsoAccount, err error) {
	var acs []*model.AsoAccount
	if acs, err = s.d.AsoAccount(c, userid, model.DefaultHash(userid)); err != nil {
		return
	}
	if len(acs) == 0 {
		err = ecode.UserNotExist
		return
	}
	if len(acs) > 1 {
		err = ecode.UserDuplicate
		return
	}
	info = acs[0]
	return
}

func (s *Service) rsaDecrypt(rsaPwd string) (res string, err error) {
	rs, err := base64.StdEncoding.DecodeString(rsaPwd)
	if err != nil {
		log.Error("failed to base64 decode RSA pwd for cloud, error(%v)", err)
		return
	}
	d, err := rsaDecryptPKCS8(s.cloudRSAKey.priv, rs)
	if err != nil {
		log.Error("failed to decrypt RSA pwd for cloud, error(%v)", err)
		return
	}
	res = string(d)
	return
}

func (s *Service) originRSAEncrypt(pwd string, timeHash string) (res string, err error) {
	ed, err := rsaEncryptPKCS8(s.originPubKey, []byte(timeHash+pwd))
	if err != nil {
		log.Error("failed to RSA encrypt pwd for origin, error(%v)", err)
		return
	}
	res = base64.StdEncoding.EncodeToString(ed)
	return
}

func (s *Service) loginOrigin(c context.Context, userid, pwd, tsHash string) (res *model.LoginToken, err error) {
	rsaPwd, err := s.originRSAEncrypt(pwd, tsHash)
	if err != nil {
		return
	}
	res, err = s.d.LoginOrigin(c, userid, rsaPwd)
	if ecode.Cause(err).Code() == _asoPasswordError {
		err = ecode.UsernameOrPasswordErr
	}
	return
}

func timeout(now int64, rsaTimeHash string) bool {
	ts, err := Hash2TsSeconds(rsaTimeHash)
	if err != nil {
		return false
	}
	return now-ts > _rsaTimeoutSeconds
}

// pwdMatches check if password matches.
// NOTE: since passport did not use salt in matching password in early period,
// those accounts which registered that time and never changed
// after passport start to use salt have empty salt,
// the schema for generating password hash in passport origin is:
// if the salt is empty, take the result of md5Hex(plain) as password hash,
// else take the result of fmt.Sprintf("%s>>BiLiSaLt<<%s", md5Hex(plainPwd), salt) as password hash.
func pwdMatches(plainPwd, salt, cloudPwdHash string) bool {
	if salt == "" {
		return model.DefaultHash(md5Hex(plainPwd)) == cloudPwdHash
	}
	return model.DefaultHash(md5Hex(fmt.Sprintf("%s>>BiLiSaLt<<%s", md5Hex(plainPwd), salt))) == cloudPwdHash
}

func (s *Service) saveToken(c context.Context, appid, subid int32, mid int64) (res *model.Perm, cache bool, err error) {
	cache = true
	now := time.Now().Unix()
	token := &model.Perm{
		Mid:         mid,
		AppID:       appid,
		AppSubID:    subid,
		AccessToken: md5Hex(fmt.Sprintf("%d,%s", mid, uuid.NewV4().String())) + _segmentation + s.currentRegion,
		CreateAt:    now,
		Expires:     now + _expireSeconds,
	}
	if _, err = s.d.AddToken(c, token); err != nil {
		return
	}
	res = token
	if err = s.d.SetTokenCache(c, token); err != nil {
		err = nil
		cache = false
	}
	return
}

// LoginOrigin login via passport api.
func (s *Service) LoginOrigin(c context.Context, query, cookie string) (res *model.LoginToken, err error) {
	return s.d.Login(c, query, cookie)
}
