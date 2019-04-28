package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

// QueryAccount is verify account is exist
func (s *Service) QueryAccount(c context.Context, req *model.QueryInfoReq) (res *model.QueryInfoResp, err error) {
	if req.CToken == "" || req.Code == "" {
		err = ecode.RequestErr
	}
	if req.QType == "" {
		err = ecode.RequestErr
		return
	}
	if _, err = s.d.Verify(c, req.Code, req.CToken); err != nil {
		log.Error("Verify code error(%v)", err)
		err = ecode.CaptchaErr
		return
	}
	res = new(model.QueryInfoResp)
	var v *model.MIDInfo
	if v, err = s.d.GetMidInfo(c, req.QType, req.QValue); err != nil {
		return
	}
	log.Info("QueryAccount:GetMidInfo req: %+v, res: %+v", req, v)
	if v == nil || v.Count == 0 {
		// 用户不存在
		res.Status = 2
		return
	}
	if v != nil {
		// 用户存在多个
		if v.Count > 1 {
			res.Status = 3
			return
		}
		if v.Count == 1 {
			res.Status = 1
		}
	}
	if res.UID, err = strconv.ParseInt(v.Mids, 10, 64); err != nil {
		log.Error("QueryAccount strconv ParseInt err(%v)", err)
		return
	}
	var count int64
	//查询该账号是否在申诉中
	if count, err = s.d.GetNoDeal(c, res.UID); err != nil {
		return
	} else if count != 0 {
		res.Status = 4
	} else if count == 0 {
		res.Status = 1
	}
	return
}

// CommitInfo is commit appeal info
func (s *Service) CommitInfo(c context.Context, uinfo *model.UserInfoReq) (err error) {
	if uinfo.Captcha == "" {
		err = ecode.CaptchaErr
		return
	}
	var code string
	retry := 0
	for {
		if retry > 3 {
			break
		}
		if code, err = s.d.GetEMailCode(c, uinfo.Mid, uinfo.LinkMail); err != nil {
			retry = retry + 1
			continue
		}
		break
	}
	if err != nil || code == "" || code != uinfo.Captcha {
		err = ecode.CaptchaErr
		return
	}
	//验证通过,删除验证码
	s.d.DelEMailCode(c, uinfo.Mid, uinfo.LinkMail)
	// params check
	regType := uinfo.RegType
	if uinfo.RegTime == 0 {
		err = ecode.RequestErr
		if err != nil {
			return
		}
	}
	if uinfo.LoginAddrs == "" {
		err = ecode.RequestErr
		if err != nil {
			return
		}
	}
	if !(regType == 1 || regType == 2 || regType == 3) {
		err = ecode.RequestErr
		if err != nil {
			return
		}
	}
	if uinfo.RegAddr == "" {
		err = ecode.RequestErr
		if err != nil {
			return
		}
	}
	//字符传切割成数组后的长度校验
	if err = checkCountStr(uinfo.Unames, 3); err != nil {
		return
	}
	if err = checkCountStr(uinfo.LoginAddrs, 3); err != nil {
		return
	}
	if err = checkCountStr(uinfo.Emails, 3); err != nil {
		return
	}
	if uinfo.Phones != "" {
		phoneAddr := strings.Split(uinfo.Phones, ";")
		if err = checkCountStr(string(phoneAddr[0]), 3); err != nil {
			return
		}
		if len(phoneAddr) > 1 {
			if err = checkCountStr(string(phoneAddr[1]), 3); err != nil {
				return
			}
		}
	}

	if uinfo.Emails != "" {
		split := strings.Split(uinfo.Emails, ",")
		for _, email := range split {
			if !strings.Contains(email, "@") {
				err = ecode.RequestErr
				return
			}
		}
	}
	if !strings.Contains(uinfo.LinkMail, "@") {
		err = ecode.RequestErr
		return
	}
	if checkCountStr(uinfo.Pwds, 3) != nil {
		err = ecode.RequestErr
		if err != nil {
			return
		}
	}
	//要传递，则两个都必须有值 默认不传设置为99
	if uinfo.CardID != "" {
		if uinfo.CardType == 99 {
			err = ecode.RequestErr
			if err != nil {
				return
			}
		}
	} else {
		if uinfo.CardType != 99 {
			err = ecode.RequestErr
			if err != nil {
				return
			}
		}
	}
	if uinfo.SafeAnswer != "" {
		if uinfo.SafeQuestion == 99 {
			err = ecode.RequestErr
			if err != nil {
				return
			}
		}
	} else {
		if uinfo.SafeQuestion != 99 {
			err = ecode.RequestErr
			if err != nil {
				return
			}
		}
	}
	log.Info("CommitInfo uinfo: %+v", uinfo)
	//检查mid是否存在
	_, err = s.d.Info3(c, uinfo.Mid)
	if err != nil {
		log.Error("s.d.Info3 err(%v)", err)
		return
	}

	// 查询mid对应的上次成功找回的案件的信息：成功次数，提交时间
	sucCount, lastSucCtime, err := s.getlastSuc(c, uinfo.Mid)
	if err != nil {
		return
	}
	uinfo.LastSucCount = sucCount
	uinfo.LastSucCTime = lastSucCtime

	rid, err := s.d.InsertRecoveryInfo(c, uinfo)
	if err != nil {
		log.Error("InsertRecoveryInfo err(%v)", err)
		return
	}

	if err = s.InsertRecoveryAddit(c, rid, uinfo.Files, uinfo.BusinessMap); err != nil {
		return
	}
	//整个提交通过后,验证通过,删除验证码
	s.d.DelEMailCode(c, uinfo.Mid, uinfo.LinkMail)
	return
}

func (s *Service) getlastSuc(c context.Context, mid int64) (int64, xtime.Time, error) {
	// 查询成功找回次数
	sucCount, err := s.d.GetSuccessCount(c, mid)
	if err != nil {
		log.Error("GetSuccessCount err(%+v)", err)
		return 0, 0, err
	}
	if sucCount == 0 {
		return 0, 0, nil
	}
	// 查询上次成功找回的提交时间，无记录则默认0即可
	lastSuc, err := s.d.GetLastSuccess(c, mid)
	if err != nil {
		log.Error("GetLastSuccess err(%+v)", err)
		return 0, 0, err
	}
	return sucCount, lastSuc.LastApplyTime, nil
}

// QueryCon is Multi conditional combinatorial query
func (s *Service) QueryCon(c context.Context, aq *model.QueryRecoveryInfoReq) (res *model.MultiQueryRes, err error) {
	var (
		infos []*model.AccountRecoveryInfo
		total int64
	)
	res = new(model.MultiQueryRes)

	infos, total, err = s.d.GetAllByCon(c, aq)
	rids := make([]int64, 0, len(infos))
	mids := make([]int64, 0, len(infos))
	m := make(map[int64]*model.AccountRecoveryInfo)
	successMap := make(map[int64]*model.RecoverySuccess, len(infos))
	lastSuccessMap := make(map[int64]*model.LastSuccessData, len(infos))
	for _, i := range infos {
		m[i.Rid] = i
		rids = append(rids, i.Rid)
		mids = append(mids, i.Mid)
	}
	if len(infos) > 0 {
		if successMap, err = s.d.BatchGetRecoverySuccess(c, mids); err != nil {
			return
		}
		if lastSuccessMap, err = s.d.BatchGetLastSuccess(c, mids); err != nil {
			return
		}
	}
	if len(rids) == 0 {
		return
	}

	bizData, err := s.QueryRecoveryAddit(c, rids)
	if err != nil {
		log.Error("QueryRecoveryAddit err(%+v)", err)
		return
	}

	res.Info = make([]*model.RecoveryResInfo, 0, len(infos))
	for _, rid := range rids {
		addit, ok := bizData[rid]
		if !ok {
			addit = &model.RecoveryAddit{
				Files: []string{},
				Extra: map[string]interface{}{},
			}
		}
		res.Info = append(res.Info, &model.RecoveryResInfo{
			AccountRecoveryInfo: *m[rid],
			RecoveryAddit:       *addit,
		})
	}
	hideQueryResInfo(res.Info, aq.IsAdvanced)
	for _, info := range res.Info {
		suc, ok := successMap[info.Mid]
		if !ok {
			suc = new(model.RecoverySuccess)
		}
		info.RecoverySuccess = *suc

		lastSuc, ok := lastSuccessMap[info.Mid]
		if !ok {
			lastSuc = new(model.LastSuccessData)
		}
		info.LastSuccessData = *lastSuc
	}
	res.Page = &model.Page{
		Size:  aq.Size,
		Total: total,
	}
	return
}

// QueryRecoveryAddit query recovery addit
func (s *Service) QueryRecoveryAddit(c context.Context, rids []int64) (map[int64]*model.RecoveryAddit, error) {
	log.Info("QueryRecoveryAddit rids=%v", rids)
	addits, err := s.d.BatchGetRecoveryAddit(c, rids)
	if err != nil {
		return nil, err
	}
	recoveryAddits := make(map[int64]*model.RecoveryAddit, len(rids))
	for _, addit := range addits {
		recoveryAddit := addit.AsRecoveryAddit()
		recoveryAddits[addit.Rid] = recoveryAddit
	}
	return recoveryAddits, nil
}

// Judge is reject or agree one operation
func (s *Service) Judge(c context.Context, req *model.JudgeReq) (err error) {
	if req.Status == 1 || req.Status == 2 { // 1通过 2驳回
		if err = s.d.UpdateStatus(c, req.Status, req.Rid, req.Operator, req.OptTime, req.Remark); err != nil {
			log.Error("UpdateStatus status=%d (1 agree,2 reject) error(%v)", req.Status, err)
			return
		}
	}
	return
}

// BatchJudge is reject or agree more operation.
func (s *Service) BatchJudge(c context.Context, req *model.BatchJudgeReq) (err error) {
	if req.Status == 1 || req.Status == 2 { // 1通过 2驳回
		if err = s.batchJudge(c, req.Status, req.RidsAry, req.Operator, req.OptTime, req.Remark); err != nil {
			log.Error("batchJudge update status fail")
			return
		}
	}
	return
}

// GetCaptchaMail send captcha mail.
func (s *Service) GetCaptchaMail(c context.Context, req *model.CaptchaMailReq) (state int64, err error) {
	//邮件次数是否到达最大值
	state, _ = s.d.SetLinkMailCount(c, req.LinkMail)
	if state == 10 {
		return
	}
	go func() {
		vcode := randCode()
		if err1 := s.SendMailM(c, model.VerifyMail, req.LinkMail, vcode); err1 == nil {
			s.d.SetCaptcha(c, vcode, req.Mid, req.LinkMail)
			s.SendMailLog(c, req.Mid, model.VerifyMail, req.LinkMail, vcode)
		} else {
			log.Error("SendMailM VerifyMail fail")
		}
	}()
	state = 1
	return
}

// SendMail reject more.
func (s *Service) SendMail(c context.Context, req *model.SendMailReq) (err error) {
	log.Info("SendMail rid=%d,status=%d", req.RID, req.Status)
	mailStatus, err := s.d.GetMailStatus(c, req.RID)
	if err != nil {
		log.Error("GetMailStatus no record rid=%d,status=%d", req.RID, req.Status)
		return
	}
	if mailStatus == 0 {
		log.Info("SendMail rid=%d,status=%d", req.RID, req.Status)
		switch req.Status {
		case model.DOAgree:
			s.AddMailch(func() {
				s.agree(context.Background(), req.RID, "账号找回服务")
			})
		case model.DOReject:
			s.AddMailch(func() {
				s.reject(context.Background(), req.RID)
			})
		}
		err = s.d.UpdateMailStatus(c, req.RID)
		if err != nil {
			log.Error("UpdateMailStatus no record rid=%d,status=%d,err: %+v", req.RID, req.Status, err)
			return
		}
	}

	return
}

func (s *Service) batchJudge(c context.Context, status int64, rids []int64, operator string, optTime xtime.Time, remark string) (err error) {
	var tx *sql.Tx
	//开启事务
	if tx, err = s.d.BeginTran(c); err != nil {
		return
	}
	for _, rid := range rids {
		if err = s.d.UpdateStatus(c, status, rid, operator, optTime, remark); err != nil {
			tx.Rollback()
			return
		}
	}
	//提交
	err = tx.Commit()
	return
}

// deal deal user appeal
func (s *Service) deal(c context.Context, rid int64) (err error) {
	mid, linkMail, ctime, err := s.d.GetUinfoByRid(c, rid)
	if mid <= 0 || err != nil {
		log.Error("deal mid error,no record in account_recovery_info; rid=%d, err(%v)", rid, err)
		return
	}
	uid := strconv.FormatInt(mid, 10)
	rid1 := strconv.FormatInt(rid, 10)
	err = s.SendMailM(c, model.CommitMail, linkMail, hideUID(uid), ctime, rid1)
	if err != nil {
		log.Error("deal SendMailM  rid=%d, err(%v)", rid, err)
		return
	}
	s.SendMailLog(c, mid, model.CommitMail, linkMail, hideUID(uid), ctime, rid1)
	return
}

// agree use appeal is agree
func (s *Service) agree(c context.Context, rid int64, operator string) (err error) {
	mid, linkMail, ctime, err := s.d.GetUinfoByRid(c, rid)
	if mid <= 0 || err != nil {
		log.Error("agree mid error,no record in account_recovery_info; rid=%d, err(%v)", rid, err)
		return
	}
	if err = s.d.UpdateSuccessCount(c, mid); err != nil {
		log.Error("UpdateSuccessCount error(%v)", err)
		return
	}
	user, err := s.d.UpdatePwd(c, mid, operator)
	if err != nil {
		return
	}
	uid := strconv.FormatInt(mid, 10)
	rid1 := strconv.FormatInt(rid, 10)
	userID := user.UserID
	pwd := user.Pwd
	err = s.SendMailM(c, model.AgreeMail, linkMail, hideUID(uid), ctime, rid1, userID, pwd)
	if err != nil {
		log.Error("agree SendMailM  rid=%d, err(%v)", rid, err)
		return
	}
	s.SendMailLog(c, mid, model.AgreeMail, linkMail, hideUID(uid), ctime, rid1, userID, pwd)
	return
}

// reject use appeal is reject
func (s *Service) reject(c context.Context, rid int64) (err error) {
	mid, linkMail, ctime, err := s.d.GetUinfoByRid(c, rid)
	if mid <= 0 || err != nil {
		log.Error("reject mid error,no record in account_recovery_info; rid=%d, err(%v)", rid, err)
		return
	}
	uid := strconv.FormatInt(mid, 10)
	rid1 := strconv.FormatInt(rid, 10)
	err = s.SendMailM(c, model.RejectMail, linkMail, hideUID(uid), ctime, rid1)
	if err != nil {
		log.Error("reject SendMailM  rid=%d, err(%v)", rid, err)
		return
	}
	s.SendMailLog(c, mid, model.RejectMail, linkMail, hideUID(uid), ctime, rid1)
	return
}

// agreeMore agree more. todo 暂时未使用
func (s *Service) agreeMore(c context.Context, ridsStr string, operator string) (err error) {
	batchRes, err := s.d.GetUinfoByRidMore(c, ridsStr)
	if err != nil {
		log.Error("GetUinfoByRidMore err(%v)", err)
	}
	var mids string
	for _, res := range batchRes {
		mids += "," + res.Mid
	}
	if err = s.d.BatchUpdateSuccessCount(c, mids[1:]); err != nil {
		log.Error("BatchUpdateSuccessCount error(%v)", err)
		return
	}
	res, err := s.d.UpdateBatchPwd(c, mids[1:], operator)
	if err != nil {
		return
	}
	s.SendMailMany(c, model.AgreeMail, batchRes, res)
	return
}

// rejectMore reject more. todo 暂时未使用
func (s *Service) rejectMore(c context.Context, ridsStr string) (err error) {
	batchRes, err := s.d.GetUinfoByRidMore(c, ridsStr)
	if err != nil {
		log.Error("s.d.GetUinfoByRidMore err(%v)", err)
	}
	s.SendMailMany(c, model.RejectMail, batchRes, nil)
	return
}

// WebToken get captcha token
func (s *Service) WebToken(c context.Context) (token *model.Token, err error) {
	var (
		tokenReq *model.TokenResq
	)
	if tokenReq, err = s.d.GetToken(c, s.c.CaptchaConf.TokenBID); err != nil {
		log.Error("GetToken error(%v)", err)
		return
	}
	token = tokenReq.Data
	return
}

// Verify verify captcha
func (s *Service) Verify(c context.Context, token, code string) (err error) {
	if _, err = s.d.Verify(c, code, token); err != nil {
		log.Error("Verify error(%v)", err)
		return
	}
	return
}

// hideUID
func hideUID(mid string) (uid string) {
	if len(mid) >= 3 {
		uid = mid[:1] + "****" + mid[len(mid)-1:]
	} else {
		uid = mid[:1] + "****"
	}
	return
}

// randCode
func randCode() (vcode string) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode = fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return
}

// checkCountStr 校验常用登陆地址,密码，手机，邮箱，昵称的个数
func checkCountStr(str string, count int) (err error) {
	if str != "" {
		split := strings.Split(str, ",")
		if len(split) > count {
			err = ecode.RequestErr
		}
	}
	return
}

// hideQueryResInfo hide response info
func hideQueryResInfo(res []*model.RecoveryResInfo, IsAdvanced bool) {
	for _, resInfo := range res {
		//其他1 手机注册2 邮箱注册3
		if resInfo.RegType == 1 {
			resInfo.RegTypeStr = "其他注册"
		} else if resInfo.RegType == 2 {
			resInfo.RegTypeStr = "手机注册"
		} else if resInfo.RegType == 3 {
			resInfo.RegTypeStr = "邮箱注册"
		}

		cardMap := make(map[int64]string)
		cardMap[99] = "没有绑定证件类型"
		cardMap[0] = "身份证"
		cardMap[1] = "护照(境外签发)"
		cardMap[2] = "港澳居民来往内地通行证"
		cardMap[3] = "台湾居民来往大陆通行证"
		cardMap[4] = "护照(中国签发)"
		cardMap[5] = "外国人永久居留证"
		cardMap[6] = "其他国家或地区身份证"

		safeMap := make(map[int64]string)
		safeMap[99] = "没安全提示问题"
		safeMap[0] = "没安全提示问题"
		safeMap[1] = "你最喜欢的格言什么?"
		safeMap[2] = "你家乡的名称是什么?"
		safeMap[3] = "你读的小学叫什么?"
		safeMap[4] = "你的父亲叫什么名字?"
		safeMap[5] = "你的母亲叫什么名字?"
		safeMap[6] = "你最喜欢的偶像是谁?"
		safeMap[7] = "你最喜欢的歌曲是什么?"

		resInfo.CardTypeStr = cardMap[resInfo.CardType]

		if IsAdvanced {
			resInfo.SafeQuestionStr = safeMap[resInfo.SafeQuestion]
			resInfo.RegTimeStr = resInfo.RegTime.Time().Format("2006年01月")
		} else {
			resInfo.SafeAnswer = model.HIDEALL
			resInfo.SafeQuestionStr = model.HIDEALL
			resInfo.CardID = hideCardID(resInfo.CardID)
			resInfo.Emails = hideEmails(resInfo.Emails)
			resInfo.RegTypeStr = "**注册"
			resInfo.RegAddr = "**_**_**"
			resInfo.RegTimeStr = "**年**月"
		}

		//电话地址单独处理
		resInfo.Phones = hidePhones(resInfo.Phones, IsAdvanced)
		//数据缺失(没有设置密保，没有绑定实名认证)则不打码 99表示默认不传
		if resInfo.SysSafe == "null" || resInfo.SafeQuestion == 99 {
			resInfo.SafeAnswer = "没安全提示问题"
			resInfo.SafeQuestionStr = "无"
		}
		if resInfo.CardType == 99 {
			resInfo.CardTypeStr = "没有绑定证件类型"
			resInfo.CardID = "无"
		}
		if resInfo.Phones == "" {
			resInfo.Phones = "无"
		}
		if resInfo.Emails == "" {
			resInfo.Emails = "无"
		}
		if resInfo.UNames == "" {
			resInfo.UNames = "无"
		}
		//密码所有人均不可见
		resInfo.Pwd = hidePwds(resInfo.Pwd)
	}
}

// hidePhones
func hidePhones(phones string, IsAdvanced bool) (phoneStr string) {
	if phones != "" {
		phoneAddr := strings.Split(phones, ";")
		phoneInfo := strings.Split(string(phoneAddr[0]), ",")
		var addrInfo []string
		if len(phoneAddr) > 1 {
			addrInfo = strings.Split(string(phoneAddr[1]), ",")
		}
		for i, phone := range phoneInfo {
			if !IsAdvanced {
				phoneStr += "," + phone[:3] + "*****" + phone[len(phone)-3:]
			} else {
				phoneStr += "," + phone
			}
			if i < len(addrInfo) {
				phoneStr += "(" + addrInfo[i] + ")"
			}
		}
		phoneStr = strings.TrimLeft(phoneStr, ",")
		return
	}
	return phones
}

// hideEmails
func hideEmails(emails string) (emailStr string) {
	if emails != "" {
		split := strings.Split(emails, ",")
		for _, email := range split {
			index := strings.Index(email, "@")
			emailStr += "," + email[:3] + "*****" + email[index:]
		}
		emailStr = strings.TrimLeft(emailStr, ",")
		return
	}
	return emails
}

// hidePwds
func hidePwds(pwds string) (pwdStr string) {
	if pwds != "" {
		split := strings.Split(pwds, ",")
		m := len(split)
		for i := 0; i < m; i++ {
			pwdStr += ",*******"
		}
		pwdStr = strings.TrimLeft(pwdStr, ",")
		return
	}
	return "********"
}

// hideCardID
func hideCardID(cardID string) (cardIDStr string) {
	if cardID != "" && len(cardID) > 3 {
		cardIDStr = cardID[:2] + "********" + cardID[len(cardID)-2:]
		return
	}
	return cardID
}

// GameList GameList.
func (s *Service) GameList(c context.Context, mids string) (res []*model.GameListRes, err error) {
	res = make([]*model.GameListRes, 0)
	midArr := strings.Split(mids, ",")
	eg, _ := errgroup.WithContext(c)
	for _, midStr := range midArr {
		var mid int64
		mid, err = strconv.ParseInt(midStr, 10, 64)
		if err != nil {
			err = ecode.RequestErr
			return
		}
		eg.Go(func() (e error) {
			var items []*model.Game
			items, err = s.d.GetUserType(c, mid)
			if items == nil {
				items = make([]*model.Game, 0)
			}
			gameRes := &model.GameListRes{
				Mid:   mid,
				Items: items,
			}
			res = append(res, gameRes)
			return
		})
	}
	if err = eg.Wait(); err != nil {
		return
	}
	return
}

// InsertRecoveryAddit insert addit for commit info
func (s *Service) InsertRecoveryAddit(c context.Context, rid int64, Files []string, BusinessMap map[string]string) (err error) {
	log.Info("InsertRecoveryAddit rid=%v, BusinessMap=%v", rid, BusinessMap)
	bizFiles := strings.Join(Files, ",")
	bizExtra := BusinessMap
	extra, err := json.Marshal(bizExtra)
	if err != nil {
		log.Error("json.Unmarshal(%s) error(%+v)", string(extra), err)
		return
	}
	if err = s.d.InsertRecoveryAddit(c, rid, bizFiles, string(extra)); err != nil {
		log.Error("InsertRecoveryAddit rid=%v, extra=%v, error(%+v)", rid, string(extra), err)
		return
	}
	return nil
}
