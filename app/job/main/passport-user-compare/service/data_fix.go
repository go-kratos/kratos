package service

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"go-common/app/job/main/passport-user-compare/model"
	"go-common/library/log"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func (s *Service) fullFixed(msg chan *model.ErrorFix) {
	for {
		fix, ok := <-msg
		if !ok {
			log.Error("consumer full info closed")
			return
		}
		log.Info("full data fixed consumer msg,(%+v)", fix)
		if s.c.DataFixSwitch {
			errType := fix.ErrorType
			if notExistUserBase == errType || pwdErrorType == errType || statusErrorType == errType || statusErrorType == userIDErrorType {
				s.fixUserBase(fix.Mid, fix.Action, "full")
			}
			if notExistUserTel == errType || telErrorType == errType {
				s.fixUserTel(fix.Mid, fix.Action, "full")
			}
			if notExistUserMail == errType || mailErrorType == errType {
				s.fixUserMail(fix.Mid, fix.Action, "full")
			}
			if notExistUserSafeQuestion == errType || safeErrorType == errType {
				s.fixSafeQuestion(fix.Mid, fix.Action, "full")
			}
			if notExistUserThirdBind == errType || sinaErrorType == errType || qqErrorType == errType {
				s.fixUserSns(fix.Mid, errType, fix.Action, "full")
			}
		}
	}
}

func (s *Service) incFixed(msg chan *model.ErrorFix) {
	for {
		fix, ok := <-msg
		if !ok {
			log.Error("consumer inc info closed")
			return
		}
		log.Info("dynamic data fixed consumer msg,(%+v)", fix)
		if s.incrDataFixSwitch {
			errType := fix.ErrorType
			if notExistUserBase == errType || pwdErrorType == errType || statusErrorType == errType || statusErrorType == userIDErrorType {
				s.fixUserBase(fix.Mid, fix.Action, "incr")
			}
			if notExistUserTel == errType || telErrorType == errType {
				s.fixUserTel(fix.Mid, fix.Action, "incr")
			}
			if notExistUserMail == errType || mailErrorType == errType {
				s.fixUserMail(fix.Mid, fix.Action, "incr")
			}
			if notExistUserSafeQuestion == errType || safeErrorType == errType {
				s.fixSafeQuestion(fix.Mid, fix.Action, "incr")
			}
			if notExistUserThirdBind == errType || sinaErrorType == errType || qqErrorType == errType {
				s.fixUserSns(fix.Mid, errType, fix.Action, "incr")
			}
			if notExistUserRegOriginType == errType || userRegOriginErrorType == errType {
				s.fixUserRegOrigin(fix.Mid, errType, fix.Action, "incr")
			}
		}
	}
}

func (s *Service) fixUserBase(mid int64, action, tpe string) {
	var (
		origin *model.OriginAccount
		err    error
	)
	log.Info("data fix user base,mid is %d,action %s,type %s", mid, action, tpe)
	if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
		log.Error("data fix query account by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	var pwdByte []byte
	if pwdByte, err = hex.DecodeString(origin.Pwd); err != nil {
		log.Error("data fix hex.DecodeString(origin.Pwd) error,mid is %d,err is (%+v),origin is(%+v)", mid, err, origin)
		return
	}
	if insertAction == action {
		a := &model.UserBase{
			Mid:     origin.Mid,
			UserID:  origin.UserID,
			Pwd:     pwdByte,
			Salt:    origin.Salt,
			Status:  origin.Isleak,
			Deleted: 0,
			MTime:   origin.MTime,
		}
		if _, err = s.d.InsertUserBase(context.Background(), a); err != nil {
			log.Error("data fix s.d.InsertUserBase by mid error,mid is %d,err is (%+v)", mid, err)
			return
		}
	}
	if updateAction == action {
		a := &model.UserBase{
			Mid:    origin.Mid,
			UserID: origin.UserID,
			Pwd:    pwdByte,
			Salt:   origin.Salt,
			Status: origin.Isleak,
		}
		if _, err = s.d.UpdateUserBase(context.Background(), a); err != nil {
			log.Error("data fix s.d.UpdateUserBase by mid error,mid is %d,err is (%+v)", mid, err)
			return
		}
	}
}

func (s *Service) fixUserTel(mid int64, action, tpe string) {
	var (
		origin *model.OriginAccount
		err    error
	)
	log.Info("data fix user tel,mid is %d,action %s,type %s", mid, action, tpe)
	if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
		log.Error("data fix query account by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	ot := strings.Trim(strings.ToLower(origin.Tel), "")
	var telByte []byte
	if telByte, err = s.doEncrypt(ot); err != nil {
		log.Error("data fix  doEncrypt tel by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	var telBindTime int64
	if insertAction == action && len(ot) != 0 {
		ut := &model.UserTel{
			Mid:   origin.Mid,
			Tel:   telByte,
			Cid:   s.countryMap[origin.CountryID],
			MTime: origin.MTime,
		}
		if telBindTime, err = s.d.QueryTelBindLog(context.Background(), mid); err != nil {
			log.Error("user not exist tel.mid %d", mid)
		}
		if telBindTime > int64(filterStart) && telBindTime < int64(filterEnd) {
			telBindTime = 0
		}
		ut.TelBindTime = telBindTime
		if _, err = s.d.InsertUserTel(context.Background(), ut); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == mySQLErrCodeDuplicateEntry {
					if err = s.handlerInsertTelDuplicate(ut); err != nil {
						log.Error("fail to handlerInsertTelDuplicate userTel(%+v) error(%+v)", ut, err)
						return
					}
					err = nil
					return
				}
			}
			log.Error("fail to add user tel userTel(%+v) error(%+v)", ut, err)
			return
		}
	}
	if updateAction == action {
		ut := &model.UserTel{
			Mid: origin.Mid,
			Tel: telByte,
			Cid: s.countryMap[origin.CountryID],
		}
		if _, err = s.d.UpdateUserTel(context.Background(), ut); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == mySQLErrCodeDuplicateEntry {
					if err = s.handlerUpdateTelDuplicate(ut); err != nil {
						log.Error("fail to handlerInsertTelDuplicate userTel(%+v) error(%+v)", ut, err)
						return
					}
					err = nil
					return
				}
			}
			log.Error("fail to update user tel userTel(%+v) error(%+v)", ut, err)
			return
		}
	}
}

func (s *Service) fixUserMail(mid int64, action, tpe string) {
	var (
		origin     *model.OriginAccount
		originInfo *model.OriginAccountInfo
		err        error
	)
	if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
		log.Error("data fix query account by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	log.Info("data fix  user mail,mid is %d,action %s,type %s, origin(%+v)", mid, action, tpe, origin)
	om := strings.Trim(strings.ToLower(origin.Email), "")
	var emailByte []byte
	if emailByte, err = s.doEncrypt(om); err != nil {
		log.Error("data fix doEncrypt mail by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	if insertAction == action && len(om) != 0 {
		userMail := &model.UserEmail{
			Mid:   origin.Mid,
			Email: emailByte,
			MTime: origin.MTime,
		}
		if originInfo, err = s.d.QueryAccountInfoByMid(context.Background(), mid); err != nil {
			log.Error("fail to QueryAccountInfoByMid mid is (%+v) error(%+v)", mid, err)
			return
		}
		timestamp := originInfo.ActiveTime
		if originInfo.Spacesta >= 0 {
			userMail.Verified = 1
			userMail.EmailBindTime = timestamp
		}
		if _, err = s.d.InsertUserEmail(context.Background(), userMail); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == mySQLErrCodeDuplicateEntry {
					if err = s.handlerEmailInsertDuplicate(userMail); err != nil {
						log.Error("fail to handlerEmailInsertDuplicate userEmail(%+v) error(%+v)", userMail, err)
						return
					}
					err = nil
					return
				}
			}
			log.Error("fail to add user email userEmail(%+v) error(%+v)", userMail, err)
			return
		}
	}
	if updateAction == action {
		userMail := &model.UserEmail{
			Mid:   origin.Mid,
			Email: emailByte,
		}
		if _, err = s.d.UpdateUserMail(context.Background(), userMail); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == mySQLErrCodeDuplicateEntry {
					if err = s.handlerEmailUpdateDuplicate(userMail); err != nil {
						log.Error("fail to handlerEmailDuplicate userEmail(%+v) error(%+v)", userMail, err)
						return
					}
					err = nil
					return
				}
			}
			log.Error("fail to update user email userEmail(%+v) error(%+v)", userMail, err)
			return
		}
	}
}

func (s *Service) fixSafeQuestion(mid int64, action, tpe string) {
	var (
		accountInfo *model.OriginAccountInfo
		err         error
	)
	log.Info("data fix safe question,mid is %d,action %s,type %s", mid, action, tpe)
	if accountInfo, err = s.d.QueryAccountInfoByMid(context.Background(), mid); err != nil {
		log.Error("data fix query account info by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	if insertAction == action && len(accountInfo.SafeAnswer) != 0 {
		usq := &model.UserSafeQuestion{
			Mid:          accountInfo.Mid,
			SafeQuestion: accountInfo.SafeQuestion,
			SafeAnswer:   s.doHash(accountInfo.SafeAnswer),
			SafeBindTime: time.Now().Unix(),
		}
		if _, err = s.d.InsertUserSafeQuestion(context.Background(), usq); err != nil {
			log.Error("data fix s.d.InsertUserSafeQuestion error,mid is %d,err is (%+v)", mid, err)
			return
		}
	}
	if updateAction == action {
		usq := &model.UserSafeQuestion{
			Mid:          accountInfo.Mid,
			SafeQuestion: accountInfo.SafeQuestion,
			SafeAnswer:   s.doHash(accountInfo.SafeAnswer),
		}
		if _, err = s.d.UpdateUserSafeQuestion(context.Background(), usq); err != nil {
			log.Error("data fix s.d.UpdateUserSafeQuestion error,mid is %d,err is (%+v)", mid, err)
			return
		}
	}
}

func (s *Service) fixUserSns(mid, errType int64, action, tpe string) {
	var (
		accountSns *model.OriginAccountSns
		err        error
	)
	log.Info("data fix third bind ,mid is %d,action %s,type %s", mid, action, tpe)
	if accountSns, err = s.d.QueryAccountSnsByMid(context.Background(), mid); err != nil {
		log.Error("data fix query account sns by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}

	if insertAction == action {
		if len(accountSns.SinaAccessToken) != 0 {
			sina := &model.UserThirdBind{
				Mid:      accountSns.Mid,
				PlatForm: platformSina,
				OpenID:   string(accountSns.SinaUID),
				Token:    accountSns.SinaAccessToken,
				Expires:  accountSns.SinaAccessExpires,
			}
			if _, err = s.d.InsertUserThirdBind(context.Background(), sina); err != nil {
				log.Error("data fix s.d.InsertUserThirdBind by mid error,mid is %d,err is (%+v)", mid, err)
				return
			}
		}
		if len(accountSns.QQAccessToken) != 0 {
			qq := &model.UserThirdBind{
				Mid:      accountSns.Mid,
				PlatForm: platformQQ,
				OpenID:   accountSns.QQOpenid,
				Token:    accountSns.QQAccessToken,
				Expires:  accountSns.QQAccessExpires,
			}
			if _, err = s.d.InsertUserThirdBind(context.Background(), qq); err != nil {
				log.Error("data fix s.d.UpdateUserThirdBind by mid error,mid is %d,err is (%+v)", mid, err)
				return
			}
		}
	}

	if updateAction == action {
		if sinaErrorType == errType {
			sns := &model.UserThirdBind{
				Mid:      accountSns.Mid,
				PlatForm: platformSina,
				OpenID:   string(accountSns.SinaUID),
				Token:    accountSns.SinaAccessToken,
			}
			if _, err = s.d.UpdateUserThirdBind(context.Background(), sns); err != nil {
				log.Error("data fix s.d.UpdateUserThirdBind by mid error,mid is %d,err is (%+v)", mid, err)
				return
			}
		}
		if qqErrorType == errType {
			sns := &model.UserThirdBind{
				Mid:      accountSns.Mid,
				PlatForm: platformQQ,
				OpenID:   accountSns.QQOpenid,
				Token:    accountSns.QQAccessToken,
			}
			if _, err = s.d.UpdateUserThirdBind(context.Background(), sns); err != nil {
				log.Error("data fix s.d.UpdateUserThirdBind by mid error,mid is %d,err is (%+v)", mid, err)
				return
			}
		}
	}
}

func (s *Service) fixUserRegOrigin(mid, errType int64, action, tpe string) {
	var (
		accountInfo *model.OriginAccountInfo
		accountReg  *model.OriginAccountReg
		affected    int64
		err         error
	)
	log.Info("data fix user reg ,mid is %d,action %s,type %s", mid, action, tpe)
	if accountInfo, err = s.d.QueryAccountInfoByMid(context.Background(), mid); err != nil {
		log.Error("data fix query account info by mid error,mid is %d,err is (%+v)", mid, err)
		return
	}
	uro := &model.UserRegOrigin{
		Mid:      accountInfo.Mid,
		JoinTime: accountInfo.JoinTime,
		JoinIP:   InetAtoN(accountInfo.JoinIP),
		MTime:    accountInfo.MTime,
		CTime:    accountInfo.MTime,
	}
	if mid >= 250531100 {
		if accountReg, err = s.d.QueryAccountRegByMid(context.Background(), mid); err != nil {
			log.Error("data fix query account reg by mid error,mid is %d,err is (%+v)", mid, err)
			return
		}
		if accountReg != nil {
			uro.RegType = accountReg.RegType
			uro.Origin = accountReg.OriginType
			uro.MTime = accountReg.MTime
			uro.CTime = accountReg.CTime
			uro.AppID = accountReg.AppID
		}
	}
	if insertAction == action || updateAction == action {
		if affected, err = s.d.InsertUpdateUserRegOriginType(context.Background(), uro); err != nil {
			log.Error("data fix  InsertUpdateUserRegOrigin by mid error,mid is %d,err is (%+v)", mid, err)
			return
		}
		if affected == 0 {
			log.Error("data fix  InsertUpdateUserRegOrigin opt error,not affected ", mid, err)
			return
		}
	}
}

func (s *Service) handlerEmailInsertDuplicate(userEmail *model.UserEmail) (err error) {
	var (
		duplicateMid int64
		asoAccount   *model.OriginAccount
		affected     int64
	)
	if duplicateMid, err = s.d.GetMidByEmail(context.Background(), userEmail); err != nil {
		log.Error("handlerEmailInsertDuplicate fail to get mid by email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	if asoAccount, err = s.d.QueryAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("handlerEmailInsertDuplicate fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	// 3. 将冲突的Email设置为NULL
	dunplicateUserMailToNil := &model.UserEmail{
		Mid: duplicateMid,
	}
	if affected, err = s.d.UpdateUserMail(context.Background(), dunplicateUserMailToNil); err != nil || affected == 0 {
		log.Error("handlerEmailInsertDuplicate  s.d.UpdateUserMail to nil error.")
	}
	// 4. 插入新的email
	if affected, err = s.d.InsertUserEmail(context.Background(), userEmail); err != nil {
		log.Error("handlerEmailInsertDuplicate s.d.InsertUserEmail. userMail is (%+v),affect is %d", userEmail, affected)
	}
	// 5. 更新的tel
	om := strings.Trim(strings.ToLower(asoAccount.Email), "")
	var emailByte []byte
	if emailByte, err = s.doEncrypt(om); err != nil {
		log.Error("handlerEmailInsertDuplicate data fix doEncrypt mail by mid error,mid is %d,err is (%+v)", asoAccount.Mid, err)
		return
	}
	dunplicateUserEmail := &model.UserEmail{
		Mid:   asoAccount.Mid,
		Email: emailByte,
	}
	if affected, err = s.d.UpdateUserMail(context.Background(), dunplicateUserEmail); err != nil {
		log.Error("handlerEmailInsertDuplicate s.d.UpdateUserMail to right. userMail is (%+v),affected is %d", dunplicateUserEmail, affected)
	}
	return
}

func (s *Service) handlerEmailUpdateDuplicate(userEmail *model.UserEmail) (err error) {
	var (
		duplicateMid int64
		asoAccount   *model.OriginAccount
		affected     int64
	)
	if duplicateMid, err = s.d.GetMidByEmail(context.Background(), userEmail); err != nil {
		log.Error("handlerEmailInsertDuplicate fail to get mid by email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	if asoAccount, err = s.d.QueryAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("handlerEmailInsertDuplicate fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	// 3. 将冲突的Email设置为NULL
	duplicateUserMailToNil := &model.UserEmail{
		Mid: duplicateMid,
	}
	if affected, err = s.d.UpdateUserMail(context.Background(), duplicateUserMailToNil); err != nil || affected == 0 {
		log.Error("handlerEmailInsertDuplicate  s.d.UpdateUserMail to nil error.")
	}
	// 4. 插入新的email
	if affected, err = s.d.UpdateUserMail(context.Background(), userEmail); err != nil {
		log.Error("handlerEmailInsertDuplicate s.d.InsertUserEmail. userMail is (%+v),affected %d", userEmail, affected)
	}
	// 5. 更新的tel
	om := strings.Trim(strings.ToLower(asoAccount.Email), "")
	var emailByte []byte
	if emailByte, err = s.doEncrypt(om); err != nil {
		log.Error("handlerEmailInsertDuplicate data fix doEncrypt mail by mid error,mid is %d,err is (%+v)", asoAccount.Mid, err)
		return
	}
	duplicateUserEmail := &model.UserEmail{
		Mid:   asoAccount.Mid,
		Email: emailByte,
	}
	if affected, err = s.d.UpdateUserMail(context.Background(), duplicateUserEmail); err != nil {
		log.Error("handlerEmailInsertDuplicate s.d.UpdateUserMail to right. userMail is (%+v),affected %d", duplicateUserEmail, affected)
	}
	return
}

func (s *Service) handlerInsertTelDuplicate(userTel *model.UserTel) (err error) {
	var (
		duplicateMid        int64
		duplicateAsoAccount *model.OriginAccount
		affected            int64
	)
	// 1. 查询duplicateMid
	if duplicateMid, err = s.d.GetMidByTel(context.Background(), userTel); err != nil {
		log.Error("handlerInsertTelDuplicate to get mid by tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	// 2. 查询冲突的mid
	if duplicateAsoAccount, err = s.d.QueryAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("handlerInsertTelDuplicate to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	// 3. 将冲突的tel设置为NULL
	duplicateUserTelToNil := &model.UserTel{
		Mid: duplicateMid,
	}
	if affected, err = s.d.UpdateUserTel(context.Background(), duplicateUserTelToNil); err != nil || affected == 0 {
		log.Error("handlerInsertTelDuplicate  s.d.UpdateUserTel to nil error.")
	}
	// 4. 插入新的tel
	if affected, err = s.d.InsertUserTel(context.Background(), userTel); err != nil {
		log.Error("handlerInsertTelDuplicate s.d.InsertUserTel.userTel is (%+v),affected %d", userTel, affected)
	}
	// 5. 更新的tel
	ot := strings.Trim(strings.ToLower(duplicateAsoAccount.Tel), "")
	var telByte []byte
	if telByte, err = s.doEncrypt(ot); err != nil {
		log.Error("data fix  doEncrypt tel by mid error,mid is %d,err is (%+v)", duplicateAsoAccount.Mid, err)
		return
	}
	duplicateUserTel := &model.UserTel{
		Mid: duplicateAsoAccount.Mid,
		Cid: s.countryMap[duplicateAsoAccount.CountryID],
		Tel: telByte,
	}
	if affected, err = s.d.UpdateUserTel(context.Background(), duplicateUserTel); err != nil {
		log.Error("handlerInsertTelDuplicate s.d.UpdateUserTel to right. userTel is (%+v),affected %d", userTel, affected)
	}
	return
}

func (s *Service) handlerUpdateTelDuplicate(userTel *model.UserTel) (err error) {
	var (
		duplicateMid        int64
		duplicateAsoAccount *model.OriginAccount
		affected            int64
	)
	// 1. 查询duplicateMid
	if duplicateMid, err = s.d.GetMidByTel(context.Background(), userTel); err != nil {
		log.Error("handlerUpdateTelDuplicate to get mid by tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	// 2. 查询冲突的mid
	if duplicateAsoAccount, err = s.d.QueryAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("handlerUpdateTelDuplicate to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserTelToNil := &model.UserTel{
		Mid: duplicateMid,
	}
	// 3. 冲突的Tel设置为NULL
	if affected, err = s.d.UpdateUserTel(context.Background(), duplicateUserTelToNil); err != nil || affected == 0 {
		log.Error("handlerUpdateTelDuplicate  s.d.UpdateUserTel to nil error.")
	}
	// 4. update tel
	if affected, err = s.d.UpdateUserTel(context.Background(), userTel); err != nil {
		log.Error("handlerUpdateTelDuplicate s.d.UpdateUserTel.userTel is (%+v),affected %d", userTel, affected)
	}
	// 5. 设置冲突的tel
	ot := strings.Trim(strings.ToLower(duplicateAsoAccount.Tel), "")
	var telByte []byte
	if telByte, err = s.doEncrypt(ot); err != nil {
		log.Error("handlerUpdateTelDuplicate data fix  doEncrypt tel by mid error,mid is %d,err is (%+v)", duplicateAsoAccount.Mid, err)
		return
	}
	duplicateUserTel := &model.UserTel{
		Mid: duplicateAsoAccount.Mid,
		Cid: s.countryMap[duplicateAsoAccount.CountryID],
		Tel: telByte,
	}
	if affected, err = s.d.UpdateUserTel(context.Background(), duplicateUserTel); err != nil {
		log.Error("handlerUpdateTelDuplicate s.d.UpdateUserTel to right. userTel is (%+v),affected %d", userTel, affected)
	}
	return
}

func (s *Service) fixEmailVerified() (err error) {
	var (
		res        []*model.UserEmail
		originInfo *model.OriginAccountInfo
		start      = int64(0)
	)
	for {
		log.Info("GetUnverifiedEmail, start %d", start)
		if res, err = s.d.GetUnverifiedEmail(context.Background(), start); err != nil {
			log.Error("fail to get UserTel error(%+v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if len(res) == 0 {
			log.Info("fix email verified finished!")
			break
		}
		for _, a := range res {
			for {
				if originInfo, err = s.d.QueryAccountInfoByMid(context.Background(), a.Mid); err != nil {
					log.Error("fail to QueryAccountInfoByMid mid is (%+v) error(%+v)", a.Mid, err)
					continue
				}
				break
			}
			if originInfo.Spacesta >= 0 {
				a.Verified = 1
				_, err = s.d.UpdateUserMailVerified(context.Background(), a)
			}
		}
		start = res[len(res)-1].Mid
	}
	return
}
