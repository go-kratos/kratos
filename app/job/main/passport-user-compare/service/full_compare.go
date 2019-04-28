package service

import (
	"context"
	"encoding/hex"
	"strconv"
	"strings"

	"go-common/app/job/main/passport-user-compare/model"
	"go-common/library/log"
)

// full compare data and fix data
func (s *Service) fullCompareAndFix() {
	go s.compareAndFixedUserBase(context.Background())
	go s.compareAndFixedSafeQuestion(context.Background())
	go s.compareAndFixedSns(context.Background())
}

// full div compare and fixed use base info .
func (s *Service) compareAndFixedUserBase(c context.Context) {
	diff := s.c.FullTask.AccountEnd / fullDivSegment
	var i int64
	for i = 0; i < fullDivSegment; i++ {
		go s.compareAndFixedUserBaseDiv(context.Background(), i, diff*i, diff*(i+1))
	}
}

// full div compare safe question .
func (s *Service) compareAndFixedSafeQuestion(c context.Context) {
	for i := 0; i < 30; i++ {
		go s.compareAndFixedSafeQuestionDiv(context.Background(), i)
	}
}

// full div compare and fixed sns.
func (s *Service) compareAndFixedSns(c context.Context) {
	diff := s.c.FullTask.AccountSnsEnd / fullDivSegment
	var i int64
	for i = 0; i < fullDivSegment; i++ {
		go s.compareAndFixSnsDiv(context.Background(), i, diff*i, diff*(i+1))
	}
}

func (s *Service) compareAndFixedUserBaseDiv(c context.Context, index, start, end int64) {
	var (
		originAccounts []*model.OriginAccount
		err            error
	)
	for {
		log.Info("start full compare basic_account,index %d,start is %d ,end is %d ,step is  %d ", index, start, end, s.c.FullTask.Step)
		if originAccounts, err = s.d.BatchQueryAccount(c, start, s.c.FullTask.Step); err != nil {
			log.Error("query batch account error,error is  (%+v)", err)
			continue
		}
		for _, originAccount := range originAccounts {
			mid := originAccount.Mid
			// 密码、状态对比
			var userBase *model.UserBase
			if userBase, err = s.d.QueryUserBase(c, mid); err != nil {
				log.Error("full query user base error,mid is %d,err is(+v)", mid, err)
				continue
			}
			if userBase == nil {
				log.Info("full userBase not exist,mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserBase, Action: insertAction}
				s.fullFixChan <- errorFix
				s.mu.Lock()
				accountStat["notExistUserBase"] = accountStat["notExistUserBase"] + 1
				s.mu.Unlock()
				continue
			}
			if originAccount.Pwd != hex.EncodeToString(userBase.Pwd) || originAccount.Salt != userBase.Salt {
				log.Info("full pwd compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccount.Pwd, hex.EncodeToString(userBase.Pwd))
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: pwdErrorType, Action: updateAction}
				s.fullFixChan <- errorFix
				s.mu.Lock()
				accountStat["pwd"] = accountStat["pwd"] + 1
				s.mu.Unlock()
			}
			if originAccount.Isleak != userBase.Status {
				log.Info("full status compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccount.Isleak, userBase.Status)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: statusErrorType, Action: updateAction}
				s.fullFixChan <- errorFix
				s.mu.Lock()
				accountStat["status"] = accountStat["status"] + 1
				s.mu.Unlock()
			}
			// 对比手机号
			var userTel *model.UserTel
			if userTel, err = s.d.QueryUserTel(c, mid); err != nil {
				log.Error("full query user tel error,mid is %d,err is(+v)", mid, err)
				continue
			}
			originTel := originAccount.Tel
			if originTel != "" && userTel == nil {
				log.Info("dynamic tel not exist,mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserTel, Action: insertAction}
				s.fullFixChan <- errorFix
				s.mu.Lock()
				accountStat["notExistUserTel"] = accountStat["notExistUserTel"] + 1
				s.mu.Unlock()
				continue
			}
			if userTel != nil {
				var tel string
				if tel, err = s.doDecrypt(userTel.Tel); err != nil {
					log.Error("full doDecrypt tel error,mid is %d,tel is (%+v), err is(+v)", mid, userTel.Tel, err)
					continue
				}
				ot := strings.Trim(strings.ToLower(originTel), "")
				if ot != tel {
					log.Info("full tel compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, ot, tel)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: telErrorType, Action: updateAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					dynamicAccountStat["tel"] = dynamicAccountStat["tel"] + 1
					s.mu.Unlock()
				}
			}
			// 对比邮箱
			var userEmail *model.UserEmail
			if userEmail, err = s.d.QueryUserMail(c, mid); err != nil {
				log.Error("full query user mail error,error is  (%+v)", err)
				continue
			}
			originMail := originAccount.Email
			if originMail != "" && userEmail == nil {
				log.Info("full mail not exist,mid is %d", mid)
				errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserMail, Action: insertAction}
				s.fullFixChan <- errorFix
				s.mu.Lock()
				accountStat["notExistUserMail"] = accountStat["notExistUserMail"] + 1
				s.mu.Unlock()
				continue
			}
			if userEmail != nil {
				var mail string
				if mail, err = s.doDecrypt(userEmail.Email); err != nil {
					log.Error("full doDecrypt email error,mid is %d,email is (%+v), err is(+v)", mid, userEmail.Email, err)
					continue
				}
				om := strings.Trim(strings.ToLower(originMail), "")
				if om != mail {
					log.Info("full mail compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originMail, mail)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: mailErrorType, Action: updateAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountStat["mail"] = accountStat["mail"] + 1
					s.mu.Unlock()
				}
			}
		}
		if start > end || len(originAccounts) == 0 {
			break
		}
		start = originAccounts[len(originAccounts)-1].Mid
	}
	log.Info("End full compare basic_account,index %d ", index)
	s.mu.Lock()
	accountStat["user_base_index"] = index
	s.mu.Unlock()
	if err = s.d.SendWechat(accountStat); err != nil {
		log.Error("s.d.SendWechat account stat,error is (%+v)", err)
		return
	}
}

func (s *Service) compareAndFixedSafeQuestionDiv(c context.Context, tableIndex int) {
	var (
		start int64
		err   error
	)
	for {
		log.Info("start full compare aso_account_info,table is %d,start is %d,step is %d", tableIndex, start, s.c.FullTask.Step)
		var originAccountInfos []*model.OriginAccountInfo
		if originAccountInfos, err = s.d.BatchQueryAccountInfo(c, start, s.c.FullTask.Step, tableIndex); err != nil {
			log.Error("full query batch account info error,error is  (%+v)", err)
			continue
		}
		for _, originAccountInfo := range originAccountInfos {
			mid := originAccountInfo.Mid
			if len(originAccountInfo.SafeAnswer) != 0 {
				var userSafeQuestion *model.UserSafeQuestion
				if userSafeQuestion, err = s.d.QueryUserSafeQuestion(c, mid); err != nil {
					log.Error("full query user safe question err, mid is %d,err is(+v)", mid, err)
					continue
				}
				if userSafeQuestion == nil {
					log.Info("full userSafeQuestion not exist,mid is %d", mid)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserSafeQuestion, Action: insertAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountInfoStat["notExistUserSafeQuestion"] = accountInfoStat["notExistUserSafeQuestion"] + 1
					s.mu.Unlock()
					continue
				}
				originSafeAnswerBytes := s.doHash(originAccountInfo.SafeAnswer)
				if hex.EncodeToString(originSafeAnswerBytes) != hex.EncodeToString(userSafeQuestion.SafeAnswer) {
					log.Info("full safe question compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, originAccountInfo.SafeAnswer, hex.EncodeToString(userSafeQuestion.SafeAnswer))
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: safeErrorType, Action: updateAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountInfoStat["safe"] = accountInfoStat["safe"] + 1
					s.mu.Unlock()
				}
			}
		}
		if start > s.c.FullTask.AccountInfoEnd || len(originAccountInfos) == 0 {
			break
		}
		start = originAccountInfos[len(originAccountInfos)-1].ID
	}
	log.Info("end full compare and fix account info table is %d", tableIndex)
	s.mu.Lock()
	accountInfoStat["user_safe_question_index"] = int64(tableIndex)
	s.mu.Unlock()
	if err := s.d.SendWechat(accountInfoStat); err != nil {
		log.Error("s.d.SendWeChat account info stat ,error is (%+v)", err)
		return
	}
}

func (s *Service) compareAndFixSnsDiv(c context.Context, index, start, end int64) {
	var err error
	for {
		log.Info("start full compare aso_account_sns,index is %d,start is %d,step is %d", index, start, s.c.FullTask.Step)
		var originAccountSnses []*model.OriginAccountSns
		if originAccountSnses, err = s.d.BatchQueryAccountSns(c, start, s.c.FullTask.Step); err != nil {
			log.Error("full query batch account sns error,error is  (%+v)", err)
			continue
		}
		for _, sns := range originAccountSnses {
			mid := sns.Mid
			if sns.QQOpenid != "" {
				var userThirdBind *model.UserThirdBind
				if userThirdBind, err = s.d.QueryUserThirdBind(c, mid, platformQQ); err != nil {
					log.Error("full query user bind error ,mid is %d, error is %(+v)", mid, err)
					continue
				}
				if userThirdBind == nil {
					log.Info("full sns not exist,mid is %d", mid)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserThirdBind, Action: insertAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountSnsStat["notExistUserThirdBind"] = accountSnsStat["notExistUserThirdBind"] + 1
					s.mu.Unlock()
					continue
				}
				if userThirdBind.PlatForm == 2 && userThirdBind.OpenID != sns.QQOpenid {
					log.Info("full sns qq compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, string(sns.QQOpenid), userThirdBind.OpenID)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: qqErrorType, Action: updateAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountSnsStat["sns"] = accountSnsStat["sns"] + 1
					s.mu.Unlock()
				}
			}
			if sns.SinaUID != 0 {
				var userThirdBind *model.UserThirdBind
				if userThirdBind, err = s.d.QueryUserThirdBind(c, mid, platformSina); err != nil {
					log.Error("full query user bind error ,mid is %d, error is %(+v)", mid, err)
					continue
				}
				if userThirdBind == nil {
					log.Info("full sns not exist,mid is %d", mid)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: notExistUserThirdBind, Action: insertAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountSnsStat["notExistUserThirdBind"] = accountSnsStat["notExistUserThirdBind"] + 1
					s.mu.Unlock()
					continue
				}
				var openID int64
				if openID, err = strconv.ParseInt(userThirdBind.OpenID, 10, 64); err != nil {
					log.Error("parse error.")
					return
				}
				if userThirdBind.PlatForm == 1 && openID != sns.SinaUID {
					log.Info("full sns sina compare not match,mid is %d, origin is (%+v),new is (%+v)", mid, sns.SinaUID, userThirdBind.OpenID)
					errorFix := &model.ErrorFix{Mid: mid, ErrorType: sinaErrorType, Action: updateAction}
					s.fullFixChan <- errorFix
					s.mu.Lock()
					accountSnsStat["sns"] = accountSnsStat["sns"] + 1
					s.mu.Unlock()
				}
			}
		}
		if start > end || len(originAccountSnses) == 0 {
			break
		}
		start = originAccountSnses[len(originAccountSnses)-1].Mid
	}
	log.Info("end full compare and fix account sns index %d", index)
	s.mu.Lock()
	accountSnsStat["user_sns_index"] = index
	s.mu.Unlock()
	if err = s.d.SendWechat(accountSnsStat); err != nil {
		log.Error("s.d.SendWeChat sns stat,error is (%+v)", err)
		return
	}
}
