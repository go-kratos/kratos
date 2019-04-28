package service

import (
	"context"
	"encoding/hex"
	"strings"
	"time"

	"go-common/app/job/main/passport-user-compare/model"
	"go-common/library/log"
)

// Dynamic inc fixed data
func (s *Service) incCompareAndFix() {
	startTime, err := time.ParseInLocation(timeFormat, s.c.IncTask.StartTime, loc)
	if err != nil {
		log.Error("failed to parse end time, time.ParseInLocation(%s, %s, %v), error(%v)", timeFormat, s.c.IncTask.StartTime, loc, err)
		return
	}
	log.Info("start inc compare and fixed")
	for {
		stepDuration := time.Duration(s.c.IncTask.StepDuration)
		endTime := startTime.Add(stepDuration)
		if time.Now().Before(endTime) {
			log.Info("break dynamic ")
			break
		}
		s.compareAndFixedUserBaseInc(context.Background(), startTime, endTime)
		s.compareAndFixedSafeQuestionInc(context.Background(), startTime, endTime)
		s.compareAndFixedUserRegOriginInc(context.Background(), startTime, endTime)
		startTime = endTime
		endTime = startTime.Add(stepDuration)
	}
	if len(dynamicAccountStat) != 0 {
		if err = s.d.SendWechat(dynamicAccountStat); err != nil {
			log.Error("s.d.SendWeChat account stat,error is (%+v)", err)
			return
		}
	}
	if len(dynamicAccountInfoStat) != 0 {
		if err := s.d.SendWechat(dynamicAccountInfoStat); err != nil {
			log.Error("s.d.SendWeChat account info stat ,error is (%+v)", err)
			return
		}
	}
	if len(dynamicAccountRegStat) != 0 {
		if err := s.d.SendWechat(dynamicAccountRegStat); err != nil {
			log.Error("s.d.SendWeChat account reg stat ,error is (%+v)", err)
			return
		}
	}
	log.Info("dynamic compare and fix ,chan size is %d", len(s.incFixChan))
}

func (s *Service) compareAndFixedUserBaseInc(c context.Context, start, end time.Time) {
	var (
		originAccounts []*model.OriginAccount
		origin         *model.OriginAccount
		err            error
	)
	log.Info("dynamic compare and fix user base,tel,mail ,time start is(%+v), end time is (%+v)", start, end)
	if originAccounts, err = s.d.BatchQueryAccountByTime(c, start, end); err != nil {
		log.Error("dynamic query batch account inc error,error is  (%+v)", err)
		return
	}
	for _, originAccount := range originAccounts {
		for {
			mid := originAccount.Mid
			var userBase *model.UserBase
			if userBase, err = s.d.QueryUserBase(c, mid); err != nil {
				log.Error("dynamic query user base error,mid is %d,err is(+v)", mid, err)
				continue
			}
			// 对比密码和盐
			if userBase == nil {
				log.Info("dynamic user base not exist,mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserBase, Action: insertAction}
				s.incFixChan <- errorFix
				s.mu.Lock()
				dynamicAccountStat["notExistUserBase"] = dynamicAccountStat["notExistUserBase"] + 1
				s.mu.Unlock()
				continue
			}
			if originAccount.UserID != userBase.UserID {
				if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
					continue
				}
				if origin.UserID != userBase.UserID {
					log.Info("dynamic pwd compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccount.Pwd, hex.EncodeToString(userBase.Pwd))
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: userIDErrorType, Action: updateAction}
					s.incFixChan <- errorFix
					s.mu.Lock()
					dynamicAccountStat["userid"] = dynamicAccountStat["userid"] + 1
					s.mu.Unlock()
				}
			}
			if originAccount.Pwd != hex.EncodeToString(userBase.Pwd) || originAccount.Salt != userBase.Salt {
				if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
					continue
				}
				if origin.Pwd != hex.EncodeToString(userBase.Pwd) || origin.Salt != userBase.Salt {
					log.Info("dynamic pwd compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccount.Pwd, hex.EncodeToString(userBase.Pwd))
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: pwdErrorType, Action: updateAction}
					s.incFixChan <- errorFix
					s.mu.Lock()
					dynamicAccountStat["pwd"] = dynamicAccountStat["pwd"] + 1
					s.mu.Unlock()
				}
			}

			if originAccount.Isleak != userBase.Status {
				if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
					continue
				}
				if origin.Isleak != userBase.Status {
					log.Info("dynamic status compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccount.Isleak, userBase.Status)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: statusErrorType, Action: updateAction}
					s.incFixChan <- errorFix
					s.mu.Lock()
					dynamicAccountStat["status"] = dynamicAccountStat["status"] + 1
					s.mu.Unlock()
				}
			}

			// 对比手机号
			var userTel *model.UserTel
			if userTel, err = s.d.QueryUserTel(c, mid); err != nil {
				log.Error("dynamic query user tel error,mid is %d,err is(+v)", mid, err)
				continue
			}
			originTel := originAccount.Tel
			if originTel != "" && userTel == nil {
				log.Info("dynamic tel not exist,mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserTel, Action: insertAction}
				s.incFixChan <- errorFix
				s.mu.Lock()
				dynamicAccountStat["notExistUserTel"] = dynamicAccountStat["notExistUserTel"] + 1
				s.mu.Unlock()
				continue
			}
			if userTel != nil {
				var tel string
				if tel, err = s.doDecrypt(userTel.Tel); err != nil {
					log.Error("dynamic doDecrypt tel error,mid is %d,tel is (%+v), err is(+v)", mid, userTel.Tel, err)
					continue
				}
				ot := strings.Trim(strings.ToLower(originTel), "")
				if ot != tel {
					if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
						continue
					}
					ot = strings.Trim(strings.ToLower(origin.Tel), "")
					if ot != tel {
						log.Info("dynamic tel compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, ot, tel)
						errorFix := &model.ErrorFix{Mid: mid, ErrorType: telErrorType, Action: updateAction}
						s.incFixChan <- errorFix
						s.mu.Lock()
						dynamicAccountStat["tel"] = dynamicAccountStat["tel"] + 1
						s.mu.Unlock()
					}
				}
			}
			// 对比邮箱
			var userEmail *model.UserEmail
			if userEmail, err = s.d.QueryUserMail(c, mid); err != nil {
				log.Error("dynamic query user mail error,error is  (%+v)", err)
				continue
			}
			originMail := originAccount.Email
			if originMail != "" && userEmail == nil {
				log.Info("dynamic mail not exist,mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserMail, Action: insertAction}
				s.incFixChan <- errorFix
				s.mu.Lock()
				dynamicAccountStat["notExistUserMail"] = dynamicAccountStat["notExistUserMail"] + 1
				s.mu.Unlock()
				continue
			}
			if userEmail != nil {
				var mail string
				if mail, err = s.doDecrypt(userEmail.Email); err != nil {
					log.Error("dynamic doDecrypt email error,mid is %d,email is (%+v), err is(+v)", mid, userEmail.Email, err)
					continue
				}
				om := strings.Trim(strings.ToLower(originMail), "")
				if om != mail {
					if origin, err = s.d.QueryAccountByMid(context.Background(), mid); err != nil {
						continue
					}
					om = strings.Trim(strings.ToLower(origin.Email), "")
					if om != mail {
						log.Info("dynamic mail compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, om, mail)
						errorFix := &model.ErrorFix{Mid: mid, ErrorType: mailErrorType, Action: updateAction}
						s.incFixChan <- errorFix
						s.mu.Lock()
						dynamicAccountStat["mail"] = dynamicAccountStat["mail"] + 1
						s.mu.Unlock()
					}
				}
			}
			break
		}
	}
	s.mu.Lock()
	dynamicAccountStat["dynamicUserBase"] = 0
	s.mu.Unlock()
}

func (s *Service) compareAndFixedSafeQuestionInc(c context.Context, start, end time.Time) {

	for i := 0; i < 30; i++ {
		var (
			err                error
			originAccountInfos []*model.OriginAccountInfo
			originInfo         *model.OriginAccountInfo
		)
		log.Info("dynamic compare and fix safe ,time start is(%+v), end time is (%+v)", start, end)
		if originAccountInfos, err = s.d.BatchQueryAccountInfoByTime(c, start, end, i); err != nil {
			log.Error("dynamic query batch account info inc error,error is  (%+v)", err)
			return
		}
		for _, originAccountInfo := range originAccountInfos {
			for {
				mid := originAccountInfo.Mid
				// fixed userSafeQuestion
				var userSafeQuestion *model.UserSafeQuestion
				if userSafeQuestion, err = s.d.QueryUserSafeQuestion(c, mid); err != nil {
					log.Error("dynamic query user safe question err, mid is %d,err is(+v)", mid, err)
					continue
				}
				if len(originAccountInfo.SafeAnswer) == 0 && userSafeQuestion == nil {
					continue
				}
				if len(originAccountInfo.SafeAnswer) != 0 && userSafeQuestion == nil {
					log.Info("dynamic safe question not exist, mid is %d", mid)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserSafeQuestion, Action: insertAction}
					s.incFixChan <- errorFix
					s.mu.Lock()
					dynamicAccountInfoStat["notExistUserSafeQuestion"] = dynamicAccountInfoStat["notExistUserSafeQuestion"] + 1
					s.mu.Unlock()
					continue
				}
				originSafeQuestion := originAccountInfo.SafeQuestion
				newSafeQuestion := userSafeQuestion.SafeQuestion
				if originSafeQuestion != newSafeQuestion {
					if originInfo, err = s.d.QueryAccountInfoByMid(context.Background(), mid); err != nil {
						continue
					}
					if originInfo.SafeQuestion != newSafeQuestion {
						log.Info("dynamic safe question index compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originSafeQuestion, newSafeQuestion)
						errorFix := &model.ErrorFix{Mid: mid, ErrorType: safeErrorType, Action: updateAction}
						s.incFixChan <- errorFix
						s.mu.Lock()
						dynamicAccountInfoStat["safe"] = dynamicAccountInfoStat["safe"] + 1
						s.mu.Unlock()
						continue
					}
				}
				originSafeAnswerBytes := s.doHash(originAccountInfo.SafeAnswer)
				if hex.EncodeToString(originSafeAnswerBytes) != hex.EncodeToString(userSafeQuestion.SafeAnswer) {
					if originInfo, err = s.d.QueryAccountInfoByMid(context.Background(), mid); err != nil {
						continue
					}
					if hex.EncodeToString(s.doHash(originAccountInfo.SafeAnswer)) != hex.EncodeToString(userSafeQuestion.SafeAnswer) {
						log.Info("dynamic safe question answer compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccountInfo.SafeAnswer, hex.EncodeToString(userSafeQuestion.SafeAnswer))
						errorFix := &model.ErrorFix{Mid: mid, ErrorType: safeErrorType, Action: updateAction}
						s.incFixChan <- errorFix
						s.mu.Lock()
						dynamicAccountInfoStat["safe"] = dynamicAccountInfoStat["safe"] + 1
						s.mu.Unlock()
						continue
					}
				}
				var uro *model.UserRegOrigin
				if uro, err = s.d.GetUserRegOriginByMid(c, mid); err != nil {
					log.Error("dynamic query user reg origin err, mid is %d,err is(+v)", mid, err)
					continue
				}
				if uro == nil {
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserRegOriginType, Action: insertAction}
					s.incFixChan <- errorFix
					s.mu.Lock()
					dynamicAccountRegStat["notExistUserRegOrigin"] = dynamicAccountRegStat["notExistUserRegOrigin"] + 1
					s.mu.Unlock()
					continue
				}
				if uro.JoinIP != InetAtoN(originAccountInfo.JoinIP) || uro.JoinTime != originAccountInfo.JoinTime {
					if originInfo, err = s.d.QueryAccountInfoByMid(context.Background(), mid); err != nil {
						continue
					}
					if uro.JoinIP != InetAtoN(originInfo.JoinIP) || uro.JoinTime != originInfo.JoinTime {
						errorFix := &model.ErrorFix{Mid: mid, ErrorType: userRegOriginErrorType, Action: updateAction}
						s.incFixChan <- errorFix
						s.mu.Lock()
						dynamicAccountRegStat["userRegOrigin"] = dynamicAccountRegStat["userRegOrigin"] + 1
						s.mu.Unlock()
						continue
					}
				}
				break
			}
		}
	}
	s.mu.Lock()
	dynamicAccountInfoStat["dynamicUserSafeQuestion"] = 0
	s.mu.Unlock()
}

func (s *Service) compareAndFixedUserRegOriginInc(c context.Context, start, end time.Time) {

	for i := 0; i < 20; i++ {
		var (
			err               error
			originAccountRegs []*model.OriginAccountReg
		)
		log.Info("dynamic compare and fix account origin reg ,time start is(%+v), end time is (%+v)", start, end)
		if originAccountRegs, err = s.d.BatchQueryAccountRegByTime(c, start, end, i); err != nil {
			log.Error("dynamic query batch account info inc error,error is  (%+v)", err)
			return
		}
		for _, originAccountReg := range originAccountRegs {
			mid := originAccountReg.Mid
			// fixed userSafeQuestion
			var uro *model.UserRegOrigin
			if uro, err = s.d.GetUserRegOriginByMid(c, mid); err != nil {
				log.Error("dynamic query user reg origin err, mid is %d,err is(+v)", mid, err)
				continue
			}
			if uro == nil {
				log.Info("dynamic user reg origin not exist, mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserRegOriginType, Action: insertAction}
				s.incFixChan <- errorFix
				s.mu.Lock()
				dynamicAccountRegStat["notExistUserRegOrigin"] = dynamicAccountRegStat["notExistUserRegOrigin"] + 1
				s.mu.Unlock()
				continue
			}
			if mid <= 250531100 {
				continue
			}
			if uro.RegType != originAccountReg.RegType || uro.Origin != originAccountReg.OriginType {
				log.Info("dynamic user reg origin update, mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: userRegOriginErrorType, Action: updateAction}
				s.incFixChan <- errorFix
				s.mu.Lock()
				dynamicAccountRegStat["userRegOrigin"] = dynamicAccountRegStat["userRegOrigin"] + 1
				s.mu.Unlock()
				continue
			}
		}
	}
	s.mu.Lock()
	dynamicAccountRegStat["dynamicUserRegOrigin"] = 0
	s.mu.Unlock()
}
