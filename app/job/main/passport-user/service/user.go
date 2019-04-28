package service

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"go-common/app/job/main/passport-user/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_platformSINA = 1
	_platformQQ   = 2
)

func (s *Service) convertAccountToUserBase(a *model.OriginAccount) (res *model.UserBase, err error) {
	pwdBytes, err := hex.DecodeString(a.Pwd)
	if err != nil {
		log.Error("fail to hex decode pwd(%+v) error(%+v)", a.Pwd, err)
		return
	}
	res = &model.UserBase{
		Mid:     a.Mid,
		UserID:  a.UserID,
		Pwd:     pwdBytes,
		Salt:    a.Salt,
		Status:  a.Isleak,
		Deleted: 0,
		MTime:   a.MTime,
	}
	return
}

func (s *Service) convertAccountToUserEmail(a *model.OriginAccount) (res *model.UserEmail) {
	if a.Email == "" {
		res = &model.UserEmail{
			Mid:           a.Mid,
			EmailBindTime: time.Now().Unix(),
			MTime:         a.MTime,
		}
		return
	}
	res = &model.UserEmail{
		Mid:           a.Mid,
		Email:         s.doEncrypt(a.Email),
		EmailBindTime: time.Now().Unix(),
		MTime:         a.MTime,
	}
	return
}

func (s *Service) convertAccountToUserTel(a *model.OriginAccount) (res *model.UserTel) {
	if a.Tel == "" {
		res = &model.UserTel{
			Mid:         a.Mid,
			TelBindTime: time.Now().Unix(),
			MTime:       a.MTime,
		}
		return
	}
	res = &model.UserTel{
		Mid:         a.Mid,
		Tel:         s.doEncrypt(a.Tel),
		Cid:         s.countryMap[a.CountryID],
		TelBindTime: time.Now().Unix(),
		MTime:       a.MTime,
	}
	return
}

func (s *Service) convertAccountInfoToUserSafeQuestion(a *model.OriginAccountInfo) (res *model.UserSafeQuestion) {
	res = &model.UserSafeQuestion{
		Mid:          a.Mid,
		SafeQuestion: a.SafeQuestion,
		SafeAnswer:   s.doHash(a.SafeAnswer),
		SafeBindTime: time.Now().Unix(),
	}
	return
}

func (s *Service) convertAccountInfoToUserRegOrigin(a *model.OriginAccountInfo) (res *model.UserRegOrigin) {
	res = &model.UserRegOrigin{
		Mid:      a.Mid,
		JoinIP:   InetAtoN(a.JoinIP),
		JoinIPV6: a.JoinIPV6,
		JoinTime: a.JoinTime,
	}
	return
}

func (s *Service) convertAccountInfoToUserEmail(a *model.OriginAccountInfo) (res *model.UserEmail) {
	var (
		verified   int32
		activeTime int64
	)
	if a.Spacesta >= 0 {
		verified = 1
		activeTime = a.ActiveTime
	}
	res = &model.UserEmail{
		Mid:           a.Mid,
		Verified:      verified,
		EmailBindTime: activeTime,
	}
	return
}

func (s *Service) convertAccountRegToUserRegOrigin(a *model.OriginAccountReg) (res *model.UserRegOrigin) {
	res = &model.UserRegOrigin{
		Mid:     a.Mid,
		Origin:  a.OriginType,
		RegType: a.RegType,
		AppID:   a.AppID,
		CTime:   a.MTime,
		MTime:   a.CTime,
	}
	return
}

func (s *Service) convertAccountSnsToUserThirdBindQQ(a *model.OriginAccountSns) (res *model.UserThirdBind) {
	res = &model.UserThirdBind{
		Mid:      a.Mid,
		OpenID:   a.QQOpenid,
		PlatForm: _platformQQ,
		Token:    a.QQAccessToken,
		Expires:  a.QQAccessExpires,
	}
	return
}

func (s *Service) convertAccountSnsToUserThirdBindSina(a *model.OriginAccountSns) (res *model.UserThirdBind) {
	res = &model.UserThirdBind{
		Mid:      a.Mid,
		OpenID:   strconv.FormatInt(a.SinaUID, 10),
		PlatForm: _platformSINA,
		Token:    a.SinaAccessToken,
		Expires:  a.SinaAccessExpires,
	}
	return
}

func (s *Service) syncAsoAccount(a *model.OriginAccount) (err error) {
	var (
		tx       *sql.Tx
		userBase *model.UserBase
	)
	userBase, err = s.convertAccountToUserBase(a)
	if err != nil {
		log.Error("fail to convert AsoAccount(%+v) to UserBase error(%+v)", a, err)
		return
	}
	tx, err = s.d.BeginTran(context.Background())
	if err != nil {
		log.Error("s.dao.Begin error(%+v)", err)
		return
	}
	if _, err = s.d.TxAddUserBase(tx, userBase); err != nil {
		log.Error("fail to add UserBase(%+v) error(%+v)", userBase, err)
		tx.Rollback()
		return
	}
	if a.Email != "" {
		userEmail := s.convertAccountToUserEmail(a)
		userEmail.EmailBindTime = 0
		if _, err = s.d.TxAddUserEmail(tx, userEmail); err != nil {
			log.Error("fail to add userEmail(%+v) error(%+v)", userEmail, err)
			tx.Rollback()
			return
		}
	}
	if a.Tel != "" {
		userTel := s.convertAccountToUserTel(a)
		userTel.TelBindTime = 0
		if _, err = s.d.TxAddUserTel(tx, userTel); err != nil {
			log.Error("fail to add userTel(%+v) error(%+v)", userTel, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(), error(%v)", err)
		return
	}
	return
}

func (s *Service) syncAsoAccountInfo(a *model.OriginAccountInfo) (err error) {
	var (
		tx               *sql.Tx
		userSafeQuestion *model.UserSafeQuestion
		userRegOrigin    *model.UserRegOrigin
		userEmail        *model.UserEmail
	)
	userRegOrigin = s.convertAccountInfoToUserRegOrigin(a)
	userEmail = s.convertAccountInfoToUserEmail(a)
	tx, err = s.d.BeginTran(context.Background())
	if err != nil {
		log.Error("s.dao.Begin error(%+v)", err)
		return
	}
	if a.SafeQuestion != 0 || a.SafeAnswer != "" {
		userSafeQuestion = s.convertAccountInfoToUserSafeQuestion(a)
		userSafeQuestion.SafeBindTime = 0
		if _, err = s.d.TxInsertIgnoreUserSafeQuestion(tx, userSafeQuestion); err != nil {
			tx.Rollback()
			return
		}
	}
	if _, err = s.d.TxInsertIgnoreUserRegOrigin(tx, userRegOrigin); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(), error(%v)", err)
		return
	}
	if _, err = s.d.UpdateUserEmailBindTime(context.Background(), userEmail); err != nil {
		tx.Rollback()
		return
	}
	return
}

func (s *Service) syncAsoAccountReg(a *model.OriginAccountReg) (err error) {
	userRegOrigin := s.convertAccountRegToUserRegOrigin(a)
	if _, err = s.d.InsertUpdateUserRegOriginType(context.Background(), userRegOrigin); err != nil {
		log.Error("fail to insert update reg origin type userRegOrigin(%+v) error(%+v)", userRegOrigin, err)
	}
	return
}

func (s *Service) syncAsoAccountSns(a *model.OriginAccountSns) (err error) {
	var (
		tx   *sql.Tx
		qq   *model.UserThirdBind
		sina *model.UserThirdBind
		utb  *model.UserThirdBind
	)
	sina = s.convertAccountSnsToUserThirdBindSina(a)
	qq = s.convertAccountSnsToUserThirdBindQQ(a)
	tx, err = s.d.BeginTran(context.Background())
	if err != nil {
		log.Error("s.dao.Begin error(%+v)", err)
		return
	}
	if a.QQOpenid != "" {
		if utb, err = s.d.GetUserThirdBindByMidAndPlatform(context.Background(), a.Mid, _platformQQ); err != nil {
			return
		}
		if utb == nil {
			if _, err = s.d.TxAddUserThirdBind(tx, qq); err != nil {
				log.Error("fail to add third bind qq userThirdBind(%+v) error(%+v)", qq, err)
				tx.Rollback()
				return
			}
		} else {
			log.Error("third bind qq is exist, userThirdBind(%+v)", utb)
		}
	}
	if a.SinaUID != 0 {
		if utb, err = s.d.GetUserThirdBindByMidAndPlatform(context.Background(), a.Mid, _platformSINA); err != nil {
			return
		}
		if utb == nil {
			if _, err = s.d.TxAddUserThirdBind(tx, sina); err != nil {
				log.Error("fail to add third bind sina userThirdBind(%+v) error(%+v)", sina, err)
				tx.Rollback()
				return
			}
		} else {
			log.Error("third bind weibo is exist, userThirdBind(%+v)", utb)
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(), error(%v)", err)
		return
	}
	return
}

func (s *Service) syncAsoCountryCode() (err error) {
	res, err := s.d.AsoCountryCode(context.Background())
	if err != nil {
		log.Error("fail to get country code error(%+v)", err)
		return
	}
	for _, r := range res {
		s.d.AddCountryCode(context.Background(), r)
	}
	return
}
