package service

import (
	"context"
	"encoding/hex"
	"strings"

	"go-common/app/job/main/passport-user-compare/model"
	"go-common/library/log"
)

func (s *Service) checkTelDuplicateJob() {
	log.Info("check tel duplicate job start...")
	res, err := s.d.UserTelDuplicate(context.Background())
	if err != nil {
		log.Error("fail to get UserTelDuplicate, error(%+v)", err)
		return
	}
	var (
		asoAccount *model.OriginAccount
		userTel    *model.UserTel
	)
	for _, r := range res {
		if asoAccount, err = s.d.QueryAccountByMid(context.Background(), r.Mid); err != nil {
			log.Error("fail to check tel duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		if userTel, err = s.d.QueryUserTel(context.Background(), r.Mid); err != nil {
			log.Error("fail to check tel duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		om := strings.Trim(strings.ToLower(asoAccount.Tel), "")
		var telByte []byte
		if telByte, err = s.doEncrypt(om); err != nil {
			log.Error("checkTelDuplicateJob doEncrypt mail by mid error,mid is %d,err is (%+v)", r.Mid, err)
			return
		}
		originHex := hex.EncodeToString(telByte)
		newHex := hex.EncodeToString(userTel.Tel)
		if originHex == newHex {
			log.Info("check user tel duplicate success, userTelDuplicate(%+v)", r)
			if _, err = s.d.UpdateUserTelDuplicateStatus(context.Background(), r.ID); err != nil {
				log.Error("fail to update user tel duplicate status, id(%d) error(%+v)", r.ID, err)
			}
		} else {
			log.Info("fail to check user tel duplicate, new(%s) origin(%s) userTelDuplicate(%+v)", newHex, originHex, r)
		}
	}
	log.Info("update tel duplicate job end...")
}

func (s *Service) checkEmailDuplicateJob() {
	log.Info("check email duplicate job start...")
	res, err := s.d.UserEmailDuplicate(context.Background())
	if err != nil {
		log.Error("fail to get UserEmailDuplicate, error(%+v)", err)
		return
	}
	var (
		asoAccount *model.OriginAccount
		userEmail  *model.UserEmail
	)
	for _, r := range res {
		if asoAccount, err = s.d.QueryAccountByMid(context.Background(), r.Mid); err != nil {
			log.Error("fail to check email duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		if userEmail, err = s.d.QueryUserMail(context.Background(), r.Mid); err != nil {
			log.Error("fail to check email duplicate mid(%d) error(%+v)", r.Mid, err)
			return
		}
		om := strings.Trim(strings.ToLower(asoAccount.Email), "")
		var emailByte []byte
		if emailByte, err = s.doEncrypt(om); err != nil {
			log.Error("checkEmailDuplicateJob doEncrypt mail by mid error,mid is %d,err is (%+v)", r.Mid, err)
			return
		}
		originHex := hex.EncodeToString(emailByte)
		newHex := hex.EncodeToString(userEmail.Email)
		if originHex == newHex {
			log.Info("check user email duplicate success, userEmailDuplicate(%+v)", r)
			if _, err = s.d.UpdateUserEmailDuplicateStatus(context.Background(), r.ID); err != nil {
				log.Error("fail to update user email duplicate status, id(%d) error(%+v)", r.ID, err)
			}
		} else {
			log.Info("fail to check user email duplicate, new(%s) origin(%s) userEmailDuplicate(%+v)", newHex, originHex, r)
		}
	}
	log.Info("update email duplicate job end...")
}
