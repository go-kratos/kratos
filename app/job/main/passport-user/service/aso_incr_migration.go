package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/passport-user/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type asoAccountBMsg struct {
	Action    string
	Table     string
	New       *model.OriginAccount
	Timestamp int64
}

type asoAccountInfoBMsg struct {
	Action    string
	Table     string
	New       *model.OriginAccountInfo
	Timestamp int64
}

type asoAccountRegBMsg struct {
	Action    string
	Table     string
	New       *model.OriginAccountReg
	Timestamp int64
}

type asoAccountSnsBMsg struct {
	Action    string
	Table     string
	New       *model.OriginAccountSns
	Timestamp int64
}

func (s *Service) consumeproc() {
	s.group.New = func(msg *databus.Message) (res interface{}, err error) {
		bmsg := new(model.BMsg)
		if err = json.Unmarshal(msg.Value, bmsg); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", string(msg.Value), err)
			return
		}
		log.Info("receive msg action(%s) table(%s) key(%s) partition(%d) offset(%d) timestamp(%d) New(%s) Old(%s)",
			bmsg.Action, bmsg.Table, msg.Key, msg.Partition, msg.Offset, msg.Timestamp, string(bmsg.New), string(bmsg.Old))
		if bmsg.Table == _asoAccountTable {
			asoAccountBMsg := &asoAccountBMsg{
				Action:    bmsg.Action,
				Table:     bmsg.Table,
				Timestamp: msg.Timestamp,
			}
			newAccount := new(model.OriginAccount)
			if err = json.Unmarshal(bmsg.New, newAccount); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", string(bmsg.New), err)
				return
			}
			asoAccountBMsg.New = newAccount
			return asoAccountBMsg, nil
		} else if bmsg.Table == _asoAccountSnsTable {
			asoAccountSnsBMsg := &asoAccountSnsBMsg{
				Action:    bmsg.Action,
				Table:     bmsg.Table,
				Timestamp: msg.Timestamp,
			}
			newAccountSns := new(model.OriginAccountSns)
			if err = json.Unmarshal(bmsg.New, newAccountSns); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", string(bmsg.New), err)
				return
			}
			asoAccountSnsBMsg.New = newAccountSns
			return asoAccountSnsBMsg, nil
		} else if strings.HasPrefix(bmsg.Table, _asoAccountInfoTable) {
			asoAccountInfoBMsg := &asoAccountInfoBMsg{
				Action:    bmsg.Action,
				Table:     bmsg.Table,
				Timestamp: msg.Timestamp,
			}
			newAccountInfo := new(model.OriginAccountInfo)
			if err = json.Unmarshal(bmsg.New, newAccountInfo); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", string(bmsg.New), err)
				return
			}
			asoAccountInfoBMsg.New = newAccountInfo
			return asoAccountInfoBMsg, nil
		} else if strings.HasPrefix(bmsg.Table, _asoAccountRegOriginTable) {
			asoAccountRegBMsg := &asoAccountRegBMsg{
				Action:    bmsg.Action,
				Table:     bmsg.Table,
				Timestamp: msg.Timestamp,
			}
			newAccountReg := new(model.OriginAccountReg)
			if err = json.Unmarshal(bmsg.New, newAccountReg); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", string(bmsg.New), err)
				return
			}
			asoAccountRegBMsg.New = newAccountReg
			return asoAccountRegBMsg, nil
		}
		return
	}
	s.group.Split = func(msg *databus.Message, data interface{}) int {
		if t, ok := data.(*asoAccountBMsg); ok {
			return int(t.New.Mid)
		} else if t, ok := data.(*asoAccountSnsBMsg); ok {
			return int(t.New.Mid)
		} else if t, ok := data.(*asoAccountInfoBMsg); ok {
			return int(t.New.Mid)
		} else if t, ok := data.(*asoAccountRegBMsg); ok {
			return int(t.New.Mid)
		}
		return 0
	}
	s.group.Do = func(msgs []interface{}) {
		for _, m := range msgs {
			if msg, ok := m.(*asoAccountBMsg); ok {
				for {
					if err := s.handleAsoAccount(msg); err != nil {
						log.Error("fail to handleAsoAccount msg(%+v) new(%+v) error(%+v)", msg, msg.New, err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			} else if msg, ok := m.(*asoAccountSnsBMsg); ok {
				for {
					if err := s.handleAsoAccountSns(msg); err != nil {
						log.Error("fail to handleAsoAccountSns msg(%+v) new(%+v) error(%+v)", msg, msg.New, err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			} else if msg, ok := m.(*asoAccountInfoBMsg); ok {
				for {
					if err := s.handleAsoAccountInfo(msg); err != nil {
						log.Error("fail to handleAsoAccountInfo msg(%+v) new(%+v) error(%+v)", msg, msg.New, err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			} else if msg, ok := m.(*asoAccountRegBMsg); ok {
				for {
					if err := s.handleAsoAccountReg(msg); err != nil {
						log.Error("fail to handleAsoAccountReg msg(%+v) new(%+v) error(%+v)", msg, msg.New, err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			}
		}
	}
	// start the group
	s.group.Start()
	log.Info("s.group.Start()")
}

func (s *Service) handleAsoAccount(msg *asoAccountBMsg) (err error) {
	switch msg.Action {
	case _insertAction:
		var (
			a                *model.OriginAccount
			userBase         *model.UserBase
			userEmail        *model.UserEmail
			userTel          *model.UserTel
			currentUserBase  *model.UserBase
			currentUserTel   *model.UserTel
			currentUserEmail *model.UserEmail
		)
		if a, err = s.d.GetAsoAccountByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get asoAccount by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		if userBase, err = s.convertAccountToUserBase(a); err != nil {
			log.Error("fail to convert AsoAccount(%+v) to UserBase error(%+v)", a, err)
			return
		}
		if currentUserBase, err = s.d.GetUserBaseByMid(context.Background(), a.Mid); err != nil {
			return
		}
		if currentUserBase == nil {
			if _, err = s.d.AddUserBase(context.Background(), userBase); err != nil {
				return
			}
		}
		if a.Tel != "" {
			if currentUserTel, err = s.d.GetUserTelByMid(context.Background(), a.Mid); err != nil {
				return
			}
			if currentUserTel == nil {
				userTel = s.convertAccountToUserTel(a)
				if err = s.addUserTel(userTel, msg.Timestamp); err != nil {
					return
				}
			}
		}
		if a.Email != "" {
			if currentUserEmail, err = s.d.GetUserEmailByMid(context.Background(), a.Mid); err != nil {
				return
			}
			if currentUserEmail == nil {
				userEmail = s.convertAccountToUserEmail(a)
				if err = s.addUserEmail(userEmail, msg.Timestamp); err != nil {
					return
				}
			}
			if err = s.updateUserEmailVerified(msg.New.Mid); err != nil {
				return
			}
		}
		s.addCache(func() {
			ub, e := s.d.GetUserBaseByMid(context.Background(), msg.New.Mid)
			if e == nil && ub != nil {
				s.d.SetUserBaseCache(context.Background(), ub)
			}
		})
		s.addCache(func() {
			ut, e := s.d.GetUserTelByMid(context.Background(), msg.New.Mid)
			if e == nil && ut != nil {
				s.d.SetUserTelCache(context.Background(), ut)
			}
		})
		s.addCache(func() {
			ue, e := s.d.GetUserEmailByMid(context.Background(), msg.New.Mid)
			if e == nil && ue != nil {
				s.d.SetUserEmailCache(context.Background(), ue)
			}
		})
	case _deleteAction:
		var tx *sql.Tx
		tx, err = s.d.BeginTran(context.Background())
		if err != nil {
			log.Error("s.dao.Begin error(%+v)", err)
			return
		}
		if _, err = s.d.TxDelUserBase(tx, msg.New.Mid); err != nil {
			tx.Rollback()
			return
		}
		if _, err = s.d.TxDelUserTel(tx, msg.New.Mid); err != nil {
			tx.Rollback()
			return
		}
		if _, err = s.d.TxDelUserEmail(tx, msg.New.Mid); err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit(), error(%v)", err)
			return
		}
		s.addCache(func() {
			s.d.DelUserBaseCache(context.Background(), msg.New.Mid)
		})
		s.addCache(func() {
			s.d.DelUserTelCache(context.Background(), msg.New.Mid)
		})
		s.addCache(func() {
			s.d.DelUserEmailCache(context.Background(), msg.New.Mid)
		})
	case _updateAction:
		var (
			a                *model.OriginAccount
			userBase         *model.UserBase
			userEmail        *model.UserEmail
			userTel          *model.UserTel
			currentUserBase  *model.UserBase
			currentUserEmail *model.UserEmail
			currentUserTel   *model.UserTel
		)
		if a, err = s.d.GetAsoAccountByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get asoAccount by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		if userBase, err = s.convertAccountToUserBase(a); err != nil {
			log.Error("fail to convert AsoAccount(%+v) to UserBase error(%+v)", a, err)
			return
		}
		userTel = s.convertAccountToUserTel(a)
		userEmail = s.convertAccountToUserEmail(a)
		if currentUserEmail, err = s.d.GetUserEmailByMid(context.Background(), a.Mid); err != nil {
			return
		}
		if currentUserTel, err = s.d.GetUserTelByMid(context.Background(), a.Mid); err != nil {
			return
		}
		if currentUserBase, err = s.d.GetUserBaseByMid(context.Background(), a.Mid); err != nil {
			return
		}
		if currentUserBase == nil {
			if _, err = s.d.AddUserBase(context.Background(), userBase); err != nil {
				return
			}
		} else {
			if _, err = s.d.UpdateUserBase(context.Background(), userBase); err != nil {
				return
			}
		}
		if currentUserTel != nil {
			if currentUserTel.TelBindTime == 0 {
				if err = s.updateUserTelAndBindTime(userTel, msg.Timestamp); err != nil {
					return
				}
			} else {
				if err = s.updateUserTel(userTel, msg.Timestamp); err != nil {
					return
				}
			}
		} else {
			if a.Tel != "" {
				if err = s.addUserTel(userTel, msg.Timestamp); err != nil {
					return
				}
			}
		}
		if currentUserEmail != nil {
			if currentUserEmail.EmailBindTime == 0 {
				if err = s.updateUserEmailAndBindTime(userEmail, msg.Timestamp); err != nil {
					return
				}
			} else {
				if err = s.updateUserEmail(userEmail, msg.Timestamp); err != nil {
					return
				}
			}
			if err = s.updateUserEmailVerified(msg.New.Mid); err != nil {
				return
			}
		} else {
			if a.Email != "" {
				if err = s.addUserEmail(userEmail, msg.Timestamp); err != nil {
					return
				}
				if err = s.updateUserEmailVerified(msg.New.Mid); err != nil {
					return
				}
			}
		}
		s.addCache(func() {
			ub, e := s.d.GetUserBaseByMid(context.Background(), msg.New.Mid)
			if e == nil && ub != nil {
				s.d.SetUserBaseCache(context.Background(), ub)
			}
		})
		s.addCache(func() {
			ut, e := s.d.GetUserTelByMid(context.Background(), msg.New.Mid)
			if e == nil && ut != nil {
				s.d.SetUserTelCache(context.Background(), ut)
			}
		})
		s.addCache(func() {
			ue, e := s.d.GetUserEmailByMid(context.Background(), msg.New.Mid)
			if e == nil && ue != nil {
				s.d.SetUserEmailCache(context.Background(), ue)
			}
		})
	}
	return
}

func (s *Service) handleAsoAccountSns(msg *asoAccountSnsBMsg) (err error) {
	var (
		a           *model.OriginAccountSns
		currentSina *model.UserThirdBind
		currentQQ   *model.UserThirdBind
	)
	switch msg.Action {
	case _insertAction:
		if a, err = s.d.GetAsoAccountSnsByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get asoAccountSns by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		if err = s.syncAsoAccountSns(a); err != nil {
			log.Error("fail to sync asoAccountSns(%+v) error(%+v)", a, err)
			return
		}
	case _updateAction:
		if a, err = s.d.GetAsoAccountSnsByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get asoAccountSns by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		qq := s.convertAccountSnsToUserThirdBindQQ(a)
		sina := s.convertAccountSnsToUserThirdBindSina(a)
		if currentSina, err = s.d.GetUserThirdBindByMidAndPlatform(context.Background(), a.Mid, _platformSINA); err != nil {
			return
		}
		if currentQQ, err = s.d.GetUserThirdBindByMidAndPlatform(context.Background(), a.Mid, _platformQQ); err != nil {
			return
		}
		if currentSina != nil {
			if _, err = s.d.UpdateUserThirdBind(context.Background(), sina); err != nil {
				return
			}
		} else {
			if a.SinaUID != 0 {
				if _, err = s.d.AddUserThirdBind(context.Background(), sina); err != nil {
					return
				}
			}
		}
		if currentQQ != nil {
			if _, err = s.d.UpdateUserThirdBind(context.Background(), qq); err != nil {
				return
			}
		} else {
			if a.QQOpenid != "" {
				if _, err = s.d.AddUserThirdBind(context.Background(), qq); err != nil {
					return
				}
			}
		}
	case _deleteAction:
		if _, err = s.d.DelUserThirdBind(context.Background(), msg.New.Mid); err != nil {
			return
		}
		s.addCache(func() {
			s.d.DelUserThirdBindQQCache(context.Background(), msg.New.Mid)
		})
		s.addCache(func() {
			s.d.DelUserThirdBindSinaCache(context.Background(), msg.New.Mid)
		})
		return
	}
	s.addCache(func() {
		resQQ, e := s.d.GetUserThirdBindByMidAndPlatform(context.Background(), msg.New.Mid, _platformQQ)
		if e == nil && resQQ != nil {
			s.d.SetUserThirdBindQQCache(context.Background(), resQQ)
		}
	})
	s.addCache(func() {
		resSina, e := s.d.GetUserThirdBindByMidAndPlatform(context.Background(), msg.New.Mid, _platformSINA)
		if e == nil && resSina != nil {
			s.d.SetUserThirdBindSinaCache(context.Background(), resSina)
		}
	})
	return
}

func (s *Service) handleAsoAccountInfo(msg *asoAccountInfoBMsg) (err error) {
	switch msg.Action {
	case _insertAction:
		var (
			a             *model.OriginAccountInfo
			tx            *sql.Tx
			userRegOrigin *model.UserRegOrigin
		)
		if a, err = s.d.GetAsoAccountInfoByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get asoAccountInfo by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		userRegOrigin = s.convertAccountInfoToUserRegOrigin(a)
		tx, err = s.d.BeginTran(context.Background())
		if err != nil {
			log.Error("s.dao.Begin error(%v)", err)
			return
		}
		if _, err = s.d.TxInsertUpdateUserRegOrigin(tx, userRegOrigin); err != nil {
			log.Error("fail to insert update userRegOrigin(%+v) error(%+v)", userRegOrigin, err)
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit(), error(%v)", err)
			return
		}
		s.addCache(func() {
			uro, e := s.d.GetUserRegOriginByMid(context.Background(), msg.New.Mid)
			if e == nil && uro != nil {
				s.d.SetUserRegOriginCache(context.Background(), uro)
			}
		})
	case _updateAction:
		var (
			a                       *model.OriginAccountInfo
			tx                      *sql.Tx
			userSafeQuestion        *model.UserSafeQuestion
			userEmail               *model.UserEmail
			currentUserSafeQuestion *model.UserSafeQuestion
		)
		if a, err = s.d.GetAsoAccountInfoByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get asoAccountInfo by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		userEmail = s.convertAccountInfoToUserEmail(a)
		userSafeQuestion = s.convertAccountInfoToUserSafeQuestion(a)
		if currentUserSafeQuestion, err = s.d.GetUserSafeQuestionByMid(context.Background(), msg.New.Mid); err != nil {
			log.Error("fail to get userSafeQuestion by mid(%d) error(%+v)", msg.New.Mid, err)
			return
		}
		tx, err = s.d.BeginTran(context.Background())
		if err != nil {
			log.Error("s.dao.Begin error(%v)", err)
			return
		}
		if currentUserSafeQuestion != nil {
			if _, err = s.d.TxUpdateUserSafeQuesion(tx, userSafeQuestion); err != nil {
				log.Error("fail to update user safe question userSafeQuestion(%+v) error(%+v)", userSafeQuestion, err)
				tx.Rollback()
				return
			}
		} else {
			if a.SafeQuestion != 0 || a.SafeAnswer != "" {
				if _, err = s.d.TxAddUserSafeQuestion(tx, userSafeQuestion); err != nil {
					log.Error("fail to add user safe question userSafeQuestion(%+v) error(%+v)", userSafeQuestion, err)
					tx.Rollback()
					return
				}
			}
		}
		if _, err = s.d.TxUpdateUserEmailVerified(tx, userEmail); err != nil {
			log.Error("fail to update user email userEmail(%+v) error(%+v)", userEmail, err)
			tx.Rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit(), error(%v)", err)
			return
		}
	}
	s.addCache(func() {
		usq, e := s.d.GetUserSafeQuestionByMid(context.Background(), msg.New.Mid)
		if e == nil && usq != nil {
			s.d.SetUserSafeQuestionCache(context.Background(), usq)
		}
	})
	s.addCache(func() {
		ue, e := s.d.GetUserEmailByMid(context.Background(), msg.New.Mid)
		if e == nil && ue != nil {
			s.d.SetUserEmailCache(context.Background(), ue)
		}
	})
	return
}

func (s *Service) handleAsoAccountReg(msg *asoAccountRegBMsg) (err error) {
	var a *model.OriginAccountReg
	if a, err = s.d.GetAsoAccountRegByMid(context.Background(), msg.New.Mid); err != nil {
		log.Error("fail to get asoAccountReg by mid(%d) error(%+v)", msg.New.Mid, err)
		return
	}
	switch msg.Action {
	case _insertAction:
		if err = s.syncAsoAccountReg(a); err != nil {
			log.Error("fail to sync asoAccountReg(%+v) error(%+v)", a, err)
			return
		}
		s.addCache(func() {
			uro, e := s.d.GetUserRegOriginByMid(context.Background(), msg.New.Mid)
			if e == nil && uro != nil {
				s.d.SetUserRegOriginCache(context.Background(), uro)
			}
		})
	}
	return
}

func (s *Service) handlerUpdateEmailDuplicate(userEmail *model.UserEmail, timestamp int64) (err error) {
	var (
		duplicateMid       int64
		asoAccount         *model.OriginAccount
		duplicateUserEmail *model.UserEmail
	)
	// 记录冲突日志
	userEmailDuplicate := &model.UserEmailDuplicate{
		Mid:           userEmail.Mid,
		Email:         userEmail.Email,
		Verified:      userEmail.Verified,
		EmailBindTime: userEmail.EmailBindTime,
		Timestamp:     timestamp,
	}
	if _, err = s.d.AddUserEmailDuplicate(context.Background(), userEmailDuplicate); err != nil {
		log.Error("fail to add user email duplicate userEmailDuplicate(%+v) error(%+v)", userEmailDuplicate, err)
		return
	}
	//1. 获取冲突mid
	if duplicateMid, err = s.d.GetMidByEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to get mid by email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	//2. 根据冲突mid查询老库最新数据
	if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserEmail = s.convertAccountToUserEmail(asoAccount)
	// 3. 将冲突的Email设置为NULL
	duplicateUserMailToNil := &model.UserEmail{
		Mid: duplicateMid,
	}
	log.Info("handle update email duplicate, mid(%d) duplicateMid(%d)", userEmail.Mid, duplicateMid)
	if _, err = s.d.UpdateUserEmail(context.Background(), duplicateUserMailToNil); err != nil {
		log.Error("fail to update user email duplicateUserMailToNil(%+v) error(%+v)", duplicateUserMailToNil, err)
		return
	}
	// 4. 更新email
	if _, err = s.d.UpdateUserEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to update user email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	// 5. 更新冲突的email
	if _, err = s.d.UpdateUserEmail(context.Background(), duplicateUserEmail); err != nil {
		log.Error("fail to update user email duplicateUserEmail(%+v) error(%+v)", duplicateUserEmail, err)
		return
	}
	return
}

func (s *Service) handlerUpdateEmailAndBindTimeDuplicate(userEmail *model.UserEmail, timestamp int64) (err error) {
	var (
		duplicateMid       int64
		asoAccount         *model.OriginAccount
		duplicateUserEmail *model.UserEmail
	)
	// 记录冲突日志
	userEmailDuplicate := &model.UserEmailDuplicate{
		Mid:           userEmail.Mid,
		Email:         userEmail.Email,
		Verified:      userEmail.Verified,
		EmailBindTime: userEmail.EmailBindTime,
		Timestamp:     timestamp,
	}
	if _, err = s.d.AddUserEmailDuplicate(context.Background(), userEmailDuplicate); err != nil {
		log.Error("fail to add user email duplicate userEmailDuplicate(%+v) error(%+v)", userEmailDuplicate, err)
		return
	}
	//1. 获取冲突mid
	if duplicateMid, err = s.d.GetMidByEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to get mid by email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	//2. 根据冲突mid查询老库最新数据
	if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserEmail = s.convertAccountToUserEmail(asoAccount)
	// 3. 将冲突的Email设置为NULL
	duplicateUserMailToNil := &model.UserEmail{
		Mid: duplicateMid,
	}
	log.Info("handle update email and bind time duplicate, mid(%d) duplicateMid(%d)", userEmail.Mid, duplicateMid)
	if _, err = s.d.UpdateUserEmail(context.Background(), duplicateUserMailToNil); err != nil {
		log.Error("fail to update user email duplicateUserMailToNil(%+v) error(%+v)", duplicateUserMailToNil, err)
		return
	}
	// 4. 更新email and bind time
	if _, err = s.d.UpdateUserEmailAndBindTime(context.Background(), userEmail); err != nil {
		log.Error("fail to update user email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	// 5. 更新冲突的email
	if _, err = s.d.UpdateUserEmail(context.Background(), duplicateUserEmail); err != nil {
		log.Error("fail to update user email duplicateUserEmail(%+v) error(%+v)", duplicateUserEmail, err)
		return
	}
	return
}

func (s *Service) handlerInsertEmailDuplicate(userEmail *model.UserEmail, timestamp int64) (err error) {
	var (
		duplicateMid       int64
		asoAccount         *model.OriginAccount
		duplicateUserEmail *model.UserEmail
	)
	// 记录冲突日志
	userEmailDuplicate := &model.UserEmailDuplicate{
		Mid:           userEmail.Mid,
		Email:         userEmail.Email,
		Verified:      userEmail.Verified,
		EmailBindTime: userEmail.EmailBindTime,
		Timestamp:     timestamp,
	}
	if _, err = s.d.AddUserEmailDuplicate(context.Background(), userEmailDuplicate); err != nil {
		log.Error("fail to add user email duplicate userEmailDuplicate(%+v) error(%+v)", userEmailDuplicate, err)
		return
	}
	//1. 获取冲突mid
	if duplicateMid, err = s.d.GetMidByEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to get mid by email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	//2. 根据冲突mid查询老库最新数据
	if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserEmail = s.convertAccountToUserEmail(asoAccount)
	// 3. 将冲突的Email设置为NULL
	duplicateUserMailToNil := &model.UserEmail{
		Mid: duplicateMid,
	}
	log.Info("handle insert email duplicate, mid(%d) duplicateMid(%d)", userEmail.Mid, duplicateMid)
	if _, err = s.d.UpdateUserEmail(context.Background(), duplicateUserMailToNil); err != nil {
		log.Error("fail to update user email duplicateUserMailToNil(%+v) error(%+v)", duplicateUserMailToNil, err)
		return
	}
	// 4. 新增email
	if _, err = s.d.AddUserEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to add user email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	// 5. 更新冲突的email
	if _, err = s.d.UpdateUserEmail(context.Background(), duplicateUserEmail); err != nil {
		log.Error("fail to update user email duplicateUserEmail(%+v) error(%+v)", duplicateUserEmail, err)
		return
	}
	return
}

func (s *Service) handlerUpdateTelDuplicate(userTel *model.UserTel, timestamp int64) (err error) {
	var (
		duplicateMid     int64
		asoAccount       *model.OriginAccount
		duplicateUserTel *model.UserTel
	)
	// 记录冲突日志
	userTelDuplicate := &model.UserTelDuplicate{
		Mid:         userTel.Mid,
		Tel:         userTel.Tel,
		Cid:         userTel.Cid,
		TelBindTime: userTel.TelBindTime,
		Timestamp:   timestamp,
	}
	if _, err = s.d.AddUserTelDuplicate(context.Background(), userTelDuplicate); err != nil {
		log.Error("fail to add user tel duplicate userTelDuplicate(%+v) error(%+v)", userTelDuplicate, err)
		return
	}
	//1. 获取冲突mid
	if duplicateMid, err = s.d.GetMidByTel(context.Background(), userTel); err != nil {
		log.Error("fail to get mid by tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	//2. 根据冲突mid查询老库最新数据
	if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserTel = s.convertAccountToUserTel(asoAccount)
	// 3. 将冲突的Tel设置为NULL
	duplicateUserTelToNil := &model.UserTel{
		Mid: duplicateMid,
	}
	log.Info("handle update tel duplicate, mid(%d) duplicateMid(%d)", userTel.Mid, duplicateMid)
	if _, err = s.d.UpdateUserTel(context.Background(), duplicateUserTelToNil); err != nil {
		log.Error("fail to update user tel duplicateUserTelToNil(%+v) error(%+v)", duplicateUserTelToNil, err)
		return
	}
	// 4. 更新tel
	if _, err = s.d.UpdateUserTel(context.Background(), userTel); err != nil {
		log.Error("fail to update user tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	// 5. 更新冲突tel
	if _, err = s.d.UpdateUserTel(context.Background(), duplicateUserTel); err != nil {
		log.Error("fail to update user tel duplicateUserTel(%+v) error(%+v)", duplicateUserTel, err)
		return
	}
	return
}

func (s *Service) handlerUpdateTelAndBindTimeDuplicate(userTel *model.UserTel, timestamp int64) (err error) {
	var (
		duplicateMid     int64
		asoAccount       *model.OriginAccount
		duplicateUserTel *model.UserTel
	)
	// 记录冲突日志
	userTelDuplicate := &model.UserTelDuplicate{
		Mid:         userTel.Mid,
		Tel:         userTel.Tel,
		Cid:         userTel.Cid,
		TelBindTime: userTel.TelBindTime,
		Timestamp:   timestamp,
	}
	if _, err = s.d.AddUserTelDuplicate(context.Background(), userTelDuplicate); err != nil {
		log.Error("fail to add user tel duplicate userTelDuplicate(%+v) error(%+v)", userTelDuplicate, err)
		return
	}
	//1. 获取冲突mid
	if duplicateMid, err = s.d.GetMidByTel(context.Background(), userTel); err != nil {
		log.Error("fail to get mid by tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	//2. 根据冲突mid查询老库最新数据
	if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserTel = s.convertAccountToUserTel(asoAccount)
	// 3. 将冲突的Tel设置为NULL
	duplicateUserTelToNil := &model.UserTel{
		Mid: duplicateMid,
	}
	log.Info("handle update tel and bind time duplicate, mid(%d) duplicateMid(%d)", userTel.Mid, duplicateMid)
	if _, err = s.d.UpdateUserTel(context.Background(), duplicateUserTelToNil); err != nil {
		log.Error("fail to update user tel duplicateUserTelToNil(%+v) error(%+v)", duplicateUserTelToNil, err)
		return
	}
	// 4. 更新tel
	if _, err = s.d.UpdateUserTelAndBindTime(context.Background(), userTel); err != nil {
		log.Error("fail to update user tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	// 5. 更新冲突tel
	if _, err = s.d.UpdateUserTel(context.Background(), duplicateUserTel); err != nil {
		log.Error("fail to update user tel duplicateUserTel(%+v) error(%+v)", duplicateUserTel, err)
		return
	}
	return
}

func (s *Service) handlerInsertTelDuplicate(userTel *model.UserTel, timestamp int64) (err error) {
	var (
		duplicateMid     int64
		asoAccount       *model.OriginAccount
		duplicateUserTel *model.UserTel
	)
	// 记录冲突日志
	userTelDuplicate := &model.UserTelDuplicate{
		Mid:         userTel.Mid,
		Tel:         userTel.Tel,
		Cid:         userTel.Cid,
		TelBindTime: userTel.TelBindTime,
		Timestamp:   timestamp,
	}
	if _, err = s.d.AddUserTelDuplicate(context.Background(), userTelDuplicate); err != nil {
		log.Error("fail to add user tel duplicate userTelDuplicate(%+v) error(%+v)", userTelDuplicate, err)
		return
	}
	//1. 获取冲突mid
	if duplicateMid, err = s.d.GetMidByTel(context.Background(), userTel); err != nil {
		log.Error("fail to get mid by tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	//2. 根据冲突mid查询老库最新数据
	if asoAccount, err = s.d.GetAsoAccountByMid(context.Background(), duplicateMid); err != nil {
		log.Error("fail to get asoAccount by mid(%d) error(%+v)", duplicateMid, err)
		return
	}
	duplicateUserTel = s.convertAccountToUserTel(asoAccount)
	// 3. 将冲突的Tel设置为NULL
	duplicateUserTelToNil := &model.UserTel{
		Mid: duplicateMid,
	}
	log.Info("handle insert tel duplicate, mid(%d) duplicateMid(%d)", userTel.Mid, duplicateMid)
	if _, err = s.d.UpdateUserTel(context.Background(), duplicateUserTelToNil); err != nil {
		log.Error("fail to update user tel duplicateUserTelToNil(%+v) error(%+v)", duplicateUserTelToNil, err)
		return
	}
	// 4. 更新tel
	if _, err = s.d.AddUserTel(context.Background(), userTel); err != nil {
		log.Error("fail to update user tel userTel(%+v) error(%+v)", userTel, err)
		return
	}
	// 5. 更新冲突tel
	if _, err = s.d.UpdateUserTel(context.Background(), duplicateUserTel); err != nil {
		log.Error("fail to update user tel duplicateUserTel(%+v) error(%+v)", duplicateUserTel, err)
		return
	}
	return
}

func (s *Service) updateUserEmailVerified(mid int64) (err error) {
	var (
		info      *model.OriginAccountInfo
		userEmail *model.UserEmail
	)
	if info, err = s.d.GetAsoAccountInfoByMid(context.Background(), mid); err != nil {
		log.Error("fail to get asoAccountInfo by mid(%d) error(%+v)", mid, err)
		return
	}
	userEmail = s.convertAccountInfoToUserEmail(info)
	if _, err = s.d.UpdateUserEmailVerified(context.Background(), userEmail); err != nil {
		log.Error("fail to update user email userEmail(%+v) error(%+v)", userEmail, err)
		return
	}
	return
}

func (s *Service) addUserEmail(userEmail *model.UserEmail, timestamp int64) (err error) {
	if _, err = s.d.AddUserEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to add user email userEmail(%+v) error(%+v)", userEmail, err)
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				if err = s.handlerInsertEmailDuplicate(userEmail, timestamp); err != nil {
					log.Error("fail to handlerInsertEmailDuplicate userEmail(%+v) error(%+v)", userEmail, err)
					return
				}
				err = nil
				return
			}
		}
		return
	}
	return
}

func (s *Service) updateUserEmail(userEmail *model.UserEmail, timestamp int64) (err error) {
	if _, err = s.d.UpdateUserEmail(context.Background(), userEmail); err != nil {
		log.Error("fail to update user email userEmail(%+v) error(%+v)", userEmail, err)
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				if err = s.handlerUpdateEmailDuplicate(userEmail, timestamp); err != nil {
					log.Error("fail to handlerUpdateEmailDuplicate userEmail(%+v) error(%+v)", userEmail, err)
					return
				}
				err = nil
				return
			}
		}
		return
	}
	return
}

func (s *Service) updateUserEmailAndBindTime(userEmail *model.UserEmail, timestamp int64) (err error) {
	if _, err = s.d.UpdateUserEmailAndBindTime(context.Background(), userEmail); err != nil {
		log.Error("fail to update user email and bind time userEmail(%+v) error(%+v)", userEmail, err)
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				if err = s.handlerUpdateEmailAndBindTimeDuplicate(userEmail, timestamp); err != nil {
					log.Error("fail to handlerUpdateEmailAndBindTimeDuplicate userEmail(%+v) error(%+v)", userEmail, err)
					return
				}
				err = nil
				return
			}
		}
		return
	}
	return
}

func (s *Service) addUserTel(userTel *model.UserTel, timestamp int64) (err error) {
	if _, err = s.d.AddUserTel(context.Background(), userTel); err != nil {
		log.Error("fail to add user tel userTel(%+v) error(%+v)", userTel, err)
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				if err = s.handlerInsertTelDuplicate(userTel, timestamp); err != nil {
					log.Error("fail to handlerTelDuplicate userTel(%+v) error(%+v)", userTel, err)
					return
				}
				err = nil
				return
			}
		}
		return
	}
	return
}

func (s *Service) updateUserTel(userTel *model.UserTel, timestamp int64) (err error) {
	if _, err = s.d.UpdateUserTel(context.Background(), userTel); err != nil {
		log.Error("fail to update user tel userTel(%+v) error(%+v)", userTel, err)
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				if err = s.handlerUpdateTelDuplicate(userTel, timestamp); err != nil {
					log.Error("fail to handlerTelDuplicate userTel(%+v) error(%+v)", userTel, err)
					return
				}
				err = nil
				return
			}
		}
		return
	}
	return
}

func (s *Service) updateUserTelAndBindTime(userTel *model.UserTel, timestamp int64) (err error) {
	if _, err = s.d.UpdateUserTelAndBindTime(context.Background(), userTel); err != nil {
		log.Error("fail to update user tel and bind time userTel(%+v) error(%+v)", userTel, err)
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				if err = s.handlerUpdateTelAndBindTimeDuplicate(userTel, timestamp); err != nil {
					log.Error("fail to handlerTelDuplicate userTel(%+v) error(%+v)", userTel, err)
					return
				}
				err = nil
				return
			}
		}
		return
	}
	return
}
