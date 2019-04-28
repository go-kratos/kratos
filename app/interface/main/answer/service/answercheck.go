package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/interface/main/answer/conf"
	"go-common/app/interface/main/answer/model"
	accoutCli "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/net/metadata"
	"go-common/library/text/translate/chinese"
)

var (
	// concat type type_id : _pendantIDNameMap id
	_rankIDPendantMap = map[int]int{
		// 11: 6,
		// 12: 7,
		// 14: 8,
		// 15: 9,
		// 17: 11,
		// 27: 10,
		// 28: 12,
		// 41: 13,
		// 9:  14,
		27: 124,
		28: 127,
		31: 126,
		29: 123,
		18: 121,
		8:  125,
		19: 129,
		15: 130,
		7:  128,
	}

	// concat pendant
	_pendantIDNameMap = map[int]string{
		5:   "哔哩王",
		6:   "声控",
		7:   "追番党",
		8:   "调教师",
		9:   "动感DJ",
		10:  "局座",
		11:  "攻略组",
		12:  "学霸",
		13:  "迷影者",
		14:  "全明星",
		122: "哔哩王",
		124: "声控",
		127: "追番党",
		126: "调教师",
		123: "动感DJ",
		121: "局座",
		125: "攻略组",
		129: "学霸",
		130: "迷影者",
		128: "全明星",
	}

	// 老挂件id对应新挂件id
	_oldPIDToNewMap = map[int]int{
		5:   122,
		6:   124,
		7:   127,
		8:   126,
		9:   123,
		10:  121,
		11:  125,
		12:  129,
		13:  130,
		14:  128,
		122: 122,
		124: 124,
		127: 127,
		126: 126,
		123: 123,
		121: 121,
		125: 125,
		129: 129,
		130: 130,
		128: 128,
	}

	_pendantIDImgMap = map[int]string{
		122: "/bfs/face/67ed957ae789852bcc59b1c1e3097ea23179f793.png",
		124: "/bfs/face/ff61b405cdcf8f7860c67293218340aeaed6e233.png",
		127: "/bfs/face/369098093a07af821b767eac44b51f97ee8501c5.png",
		126: "/bfs/face/9e775c3ebe224a774d4b2f99fd5be342eb6f51ec.png",
		123: "/bfs/face/939fa982d8b1c1fd653de5c7890db03d62e87226.png",
		121: "/bfs/face/7f6b5cb11ea7abd2e05b04f65f190dfb10456554.png",
		125: "/bfs/face/90cc47168e40326dc934fad7b9abb82aa748d6ac.png",
		129: "/bfs/face/42869dad53926c75e3010150c15b16a8925fb268.png",
		130: "/bfs/face/3d5ee491c125bf452b2dbec082dbb8209b645316.png",
		128: "/bfs/face/b53937110e8009a720e2426ea69c449483718b3c.png",
	}

	// 125: "攻略组",--> 题库(8,9,12,13,14)
	// 130: "迷影者",--> 题库(15,16,17)
	// 121: "局座",--> 题库(18)
	// 129: "学霸",--> 题库(19,20,21,22,23,24,25,26)
	// 124: "声控",--> 题库(27)
	// 127: "追番党",--> 题库(28)
	// 126: "调教师",--> 题库(31)
	// 123: "动感DJ",--> 题库(30,29)
	// 128: "全明星",--> 题库(35,34,33,32)
	// 122: "哔哩王",

	// 分区合并归类
	_typeIDMap = map[int][]int{
		// 11: {12, 13},                                         // 动漫作品+动漫内容
		// 15: {15, 16},                                         // ACG+三次元音乐
		// 17: {17, 18, 19, 20, 21, 22, 23, 24, 25},             // 各类游戏
		// 28: {28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 7, 8},   // 科学技术+音频+视频技术
		// 41: {41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52}, // 各类影视剧
		8:  {8, 9, 12, 13, 14},               // 游戏
		19: {19, 20, 21, 22, 23, 24, 25, 26}, // 科技
		15: {15, 16, 17},                     // 影视
		29: {29, 30},                         // 音乐
		7:  {7, 35, 34, 33, 32},              // 鬼畜+流行前线
	}

	// 兼容账号rank错误
	_rank0 = int32(0)
)

const (
	_unBindTel = 0
)

// ProCheck check second step questions
func (s *Service) ProCheck(c context.Context, mid int64, ids []int64, ansHash map[int64]string, lang string) (hid int64, err error) {
	var now = time.Now()
	if len(ids) != s.c.Answer.ProNum {
		err = ecode.AnswerQsNumErr
		return
	}
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	ah, err := s.history(c, mid)
	if err != nil || ah == nil || ah.StartTime.Add(s.answerDuration()).Before(now) || ah.Score != 0 {
		err = ecode.AnswerBaseNotPassed
		return
	}
	if (now.Unix() - ah.StepTwoStartTime.Unix()) < s.c.Answer.BlockedTimestamp {
		s.answerDao.SetBlockCache(c, mid)
		log.Error("member user answer block, time space(%v)", now.Unix()-ah.StepTwoStartTime.Unix())
		err = ecode.AnswerBlock
		return
	}
	qsidsMc, err := s.answerDao.IdsCache(c, mid, model.Q)
	if err != nil {
		err = ecode.AnswerMidCacheQidsErr
		log.Error("s.answerDao.IdsCache(%d) err(%v) ", mid, err)
		return
	}
	ok, err := s.checkQsIDs(c, ids, mid, qsidsMc, s.c.Answer.ProNum)
	if !ok {
		return
	}
	errIds, rc, err := s.checkAns(c, mid, ids, ansHash, lang, s.c.Answer.ProNum)
	if err != nil {
		return
	}
	rcJSON, err := json.Marshal(rc)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", rc, err)
		return
	}
	total := s.c.Answer.BaseNum + s.c.Answer.ProNum + int(ah.StepExtraScore)
	score := total - len(errIds)
	ahDB := &model.AnswerHistory{
		ID:             ah.ID,
		Hid:            ah.Hid,
		CompleteResult: string(rcJSON),
		CompleteTime:   now,
		Score:          int8(score),
		IsFirstPass:    0,
	}
	log.Info("user: %d, score:%d, his: %v", mid, score, ahDB)
	member, err := s.accInfo(c, mid)
	if err == nil && member != nil && score >= model.Score60 && member.Rank == model.UserInfoRank {
		ahDB.IsFirstPass = 1
	}
	ahDB.RankID = s.pendant(c, ahDB, mid, metadata.String(c, metadata.RemoteIP), rc)
	r, err := s.answerDao.SetHistory(c, mid, ahDB)
	if err != nil || r != 1 {
		return
	}
	ah.CompleteResult = ahDB.CompleteResult
	ah.CompleteTime = ahDB.CompleteTime
	ah.Score = ahDB.Score
	ah.IsFirstPass = ahDB.IsFirstPass
	ah.RankID = ahDB.RankID
	ah.Mtime = now
	s.userActionLog(mid, model.ProCheck, ah)
	if ahDB.Score >= model.Score60 && ahDB.RankID > 0 {
		if hid, _, err = s.answerDao.PendantHistory(c, mid); err != nil {
			return
		}
		if hid <= 0 {
			s.answerDao.AddPendantHistory(c, mid, ah.Hid)
		}
	}
	hid = ah.Hid
	s.missch.Do(c, func(ctx context.Context) {
		s.answerDao.DelHistoryCache(ctx, mid)
		s.answerDao.DelIdsCache(ctx, mid, model.Q)
	})
	return
}

// CheckBase check base question all
func (s *Service) CheckBase(c context.Context, mid int64, ids []int64, ansHas map[int64]string, lang string) (res *model.AnsCheck, err error) {
	var (
		now          = time.Now()
		errIds       []int64
		profileReply *accoutCli.ProfileReply
	)
	// 检查手机绑定
	if profileReply, err = s.accountSvc.Profile3(c, &accoutCli.MidReq{Mid: mid}); err != nil || profileReply == nil || profileReply.Profile == nil {
		log.Error("s.accRPC.Profile3(%d) err(%+v)", mid, err)
		err = ecode.AnswerAccCallErr
		return
	}
	if profileReply.Profile.TelStatus == _unBindTel {
		err = ecode.AnswerNeedBindTel
		return
	}
	if len(ids) < s.c.Answer.BaseNum {
		err = ecode.RequestErr
		return
	}
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	res = &model.AnsCheck{}
	at, ok := s.checkTime(c, mid, now)
	if !ok {
		err = ecode.AnswerTimeExpire
		return
	}
	if len(ids) != s.c.Answer.BaseNum {
		err = ecode.AnswerQsNumErr
		return
	}
	qsIdsMc, err := s.answerDao.IdsCache(c, mid, model.Q)
	if err != nil {
		log.Error("s.answerDao.IdsCache(%d) err(%v) ", mid, err)
		err = ecode.AnswerMidCacheQidsErr
		return
	}
	ok, err = s.checkQsIDs(c, ids, mid, qsIdsMc, s.c.Answer.BaseNum)
	if err != nil || !ok {
		return
	}
	errIds, _, err = s.checkAns(c, mid, ids, ansHas, lang, s.c.Answer.BaseNum)
	res.QidList = errIds
	if err != nil {
		return
	}
	if len(errIds) > 0 {
		return
	}
	s.basePass(c, mid, at, now)
	res.Pass = true
	return
}

// Captcha get question captcha
func (s *Service) Captcha(c context.Context, mid int64, clientType string, newCaptcha int) (res *model.ProcessRes, err error) {
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	ah, err := s.history(c, mid)
	if err != nil || ah == nil || ah.Score == model.Score0 {
		log.Info("answer Captcha(%d) answer history is null or score is zero err(%v) ", mid, err)
		if ah != nil {
			if ah.StepOneCompleteTime == 0 {
				err = ecode.AnswerBaseNotPassed
				return
			}
			if ah.StepExtraCompleteTime == 0 {
				err = ecode.AnswerExtraNoPass
				return
			}
		}
		err = ecode.AnswerProNoPass
		return
	}
	if ah.IsPassCaptcha == model.CaptchaPass {
		err = ecode.AnswerCaptchaPassed
		return
	}
	if !conf.Conf.Answer.Captcha {
		if res, err = s.preProcess(c, mid, metadata.String(c, metadata.RemoteIP), clientType, newCaptcha); err == nil {
			return
		}
		log.Error("s.preProcess(%d,%s,%d) err:%+v", mid, clientType, newCaptcha, err)
	}
	var token, url string
	if token, url, err = s.answerDao.Captcha(c); err != nil {
		return
	}
	res = &model.ProcessRes{
		Token:       token,
		URL:         url,
		CaptchaType: model.BiliCaptcha,
	}
	return
}

// Validate check question captcha
func (s *Service) Validate(c context.Context, challenge, validate, seccode, clientType string, success int, mid int64,
	cookie, captchaType string, comargs map[string]string) (res *model.AnsCheck, err error) {
	var now = time.Now()
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	res = &model.AnsCheck{}
	ah, err := s.history(c, mid)
	log.Info(" Validate ah (%d) res(%v) ", mid, ah)
	if err != nil || ah == nil || ah.Score == model.Score0 {
		log.Info("answer Validate(%d) answer history is null or score is zero err(%v) ", mid, err)
		if ah != nil {
			if ah.StepOneCompleteTime == 0 {
				err = ecode.AnswerBaseNotPassed
				return
			}
			if ah.StepExtraCompleteTime == 0 {
				err = ecode.AnswerExtraNoPass
				return
			}
		}
		err = ecode.AnswerProNoPass
		return
	}
	// passed go to next page
	if ah.IsPassCaptcha == model.CaptchaPass {
		res.Pass = true
		res.HistoryID = ah.Hid
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	switch captchaType {
	case model.BiliCaptcha:
		if err = s.answerDao.Verify(c, validate, seccode, ip); err != nil {
			log.Error("answerDao.Verify(%s,%s,%s) error:%+v", validate, seccode, ip, err)
			return
		}
		res.Pass = true
	default:
		if ok := s.validate(c, challenge, validate, seccode, clientType, ip, success, mid); !ok {
			log.Error("Validate validate(%v,%v,%v,%v,%v,%d) error(%v)", challenge, validate, seccode, clientType, success, mid, err)
			err = ecode.AnswerGeetestVaErr
			return
		}
		res.Pass = true
	}
	member, err := s.accInfo(c, mid)
	if err != nil || member == nil {
		log.Error("Validate accInfo(%d) info is null error(%v)", mid, err)
		return
	}
	if _, err = s.answerDao.UpdateCaptcha(c, ah.ID, ah.Mid, model.CaptchaPass); err != nil {
		log.Error("s.answerDao.UpdateCaptcha error (%v) ", err)
		err = ecode.ServerErr
		return
	}
	ah.IsPassCaptcha = model.CaptchaPass
	ah.Mtime = now
	s.userActionLog(mid, model.Captcha, ah)
	s.answerDao.DelHistoryCache(c, mid)
	log.Info(" Validate member (%v) rank(%d) ", ah, member.Rank)
	if ah.Score >= model.Score60 && (member.Rank == model.UserInfoRank || member.Rank == _rank0) {
		log.Info(" beFormal in (%d) ", mid)
		s.sendData(c, comargs, ah, ip)
		if err = s.accountDao.BeFormal(c, mid, ip); err != nil {
			log.Error(" beFormal fail(%d) err(%v)", mid, err)
			s.addRetryBeFormal(&model.Formal{Mid: mid, IP: ip})
			err = ecode.AnswerFormalFailed
			return
		}
		s.answerDao.UpdateLevel(c, ah.ID, ah.Mid, 1, 1)
		ah.IsFirstPass = 1
		ah.PassedLevel = 1
		ah.Mtime = now
		s.userActionLog(mid, model.Level, ah)
		s.answerDao.DelHistoryCache(c, mid)
		s.PendantRec(c, &model.ReqPendant{HID: ah.Hid, MID: mid})
	}
	res.HistoryID = ah.Hid
	return
}

// checkQsIDs check question id param.
func (s *Service) checkQsIDs(c context.Context, ids []int64, mid int64, qsIdsMc []int64, qs int) (ok bool, err error) {
	if qsIdsMc == nil {
		log.Error("CheckBase.qsIdsMc is nil (%d,%v) )", mid, qsIdsMc)
		err = ecode.AnswerMidCacheQidsErr
		return
	}
	if len(ids) != qs {
		err = ecode.AnswerQsNumErr
		return
	}
	qidMap := map[int64]bool{}
	for _, v := range qsIdsMc {
		qidMap[v] = true
	}
	i := 0
	for _, v := range ids {
		if qidMap[v] {
			i++
		}
	}
	if i == qs {
		ok = true
	} else {
		err = ecode.AnswerQidDiffRequestErr
	}
	return
}

// checkAns check question ans.
func (s *Service) checkAns(c context.Context, mid int64, ids []int64, ansHash map[int64]string, lang string, count int) (errIds []int64, rc map[int8]int, err error) {
	qs, err := s.answerDao.ByIds(c, ids)
	if err != nil || qs == nil || len(qs) != count {
		log.Error("checkAns.qs is nil (%v,%v) error(%v)", ids, qs, err)
		err = ecode.AnswerMidDBQueErr
		return
	}
	errIds = []int64{}
	rc = make(map[int8]int)
	for _, q := range qs {
		if lang == model.LangZhTW {
			q.Ans[0] = chinese.Convert(c, q.Ans[0])
		}
		if h := s.ansHash(mid, q.Ans[0]); h != ansHash[q.ID] {
			errIds = append(errIds, q.ID)
		} else {
			rc[q.TypeID]++
		}
	}
	return
}

// basePass base question pass.
func (s *Service) basePass(c context.Context, mid int64, at *model.AnswerTime, now time.Time) {
	h := &model.AnswerHistory{
		Mid:                 mid,
		StartTime:           at.Stime,
		StepOneErrTimes:     at.Etimes,
		StepOneCompleteTime: now.Unix() - at.Stime.Unix(),
		Ctime:               now,
		Mtime:               now,
	}
	r, hid, err := s.answerDao.AddHistory(c, mid, h)
	if err != nil || r != 1 {
		log.Error("answerDao.AddHistory r !=1 (%d,%v) error(%v)", mid, h, err)
		return
	}
	h.Hid, _ = strconv.ParseInt(hid, 10, 64)
	s.userActionLog(mid, model.BasePass, h)
	s.answerDao.DelHistoryCache(c, mid)
	s.answerDao.DelExpireCache(c, mid)
	s.answerDao.DelIdsCache(c, mid, model.Q)
}

// setPendant set pendant.
func (s *Service) pendant(c context.Context, ah *model.AnswerHistory, mid int64, ip string, rc map[int8]int) (rankID int) {
	var (
		ok          bool
		ht          int
		typeIDScore = map[int8]int{ // key:_typeIDMap`key,value:Score
			8:  0,
			19: 0,
			15: 0,
			29: 0,
			7:  0,
		}
	)
	if ah.Score == model.FullScore {
		return model.RankTop // 122: "哔哩王",
	}
	for k, v := range rc {
		switch k {
		case 8, 9, 12, 13, 14: // 游戏 125: "攻略组",--> 题库(8,9,12,13,14)
			typeIDScore[8] += v
		case 15, 16, 17: // 影视 130: "迷影者",--> 题库(15,16,17)
			typeIDScore[15] += v
		case 19, 20, 21, 22, 23, 24, 25, 26: // 科技 129: "学霸",--> 题库(19,20,21,22,23,24,25,26)
			typeIDScore[19] += v
		case 29, 30: // 音乐 123: "动感DJ",--> 题库(30,29)
			typeIDScore[29] += v
		case 7, 35, 34, 33, 32: // 鬼畜+流行前线  128: "全明星",--> 题库(35,34,33,32)
			typeIDScore[7] += v
		default:
			// 121: "局座",--> 题库(18)
			// 124: "声控",--> 题库(27)
			// 127: "追番党",--> 题库(28)
			// 126: "调教师",--> 题库(31)
			typeIDScore[k] += v
		}
	}
	score := 0
	for k, v := range typeIDScore {
		if score < v {
			score = v
			ht = int(k)
		}
	}
	rankID, ok = _rankIDPendantMap[ht]
	if !ok {
		log.Warn("user(%d),pendant() rankId(%d) result:%+v ", mid, _rankIDPendantMap[ht], rc)
	}
	return
}

func (s *Service) checkAnswerBlock(c context.Context, mid int64) (block bool) {
	block, _ = s.answerDao.CheckBlockCache(c, mid)
	return
}

func (s *Service) sendData(c context.Context, comargs map[string]string, ah *model.AnswerHistory, ip string) {
	s.promBeFormal.Incr("count")
	// add report bigdata log
	ans := []interface{}{
		strconv.FormatInt(ah.StepOneCompleteTime, 10),
		ah.CompleteResult,
		strconv.FormatInt(ah.CompleteTime.Unix()-ah.StepTwoStartTime.Unix(), 10),
		fmt.Sprintf("%d", ah.Score),
		strconv.FormatInt(time.Now().Unix(), 10),
	}
	s.missch.Do(c, func(ctx context.Context) {
		ac := map[string]string{
			"itemType": infoc.ItemTypeLV,
			"action":   infoc.ActionAnswer,
			"ip":       ip,
			"mid":      strconv.FormatInt(ah.Mid, 10),
			"sid":      comargs["sid"],
			"ua":       comargs["ua"],
			"buvid":    comargs["buvid"],
			"refer":    comargs["refer"],
			"url":      comargs["url"],
		}
		log.Info("s.infoc2.ServiceAntiCheatBus(%v,%v)", ac, ans)
		s.infoc2.ServiceAntiCheatBus(ac, ans)
	})
}

// ExtraCheck extra check.
func (s *Service) ExtraCheck(c context.Context, mid int64, ids []int64, ansHash map[int64]string, ua string, lang string, refer string, buvid string) (err error) {
	var now = time.Now()
	if len(ids) < s.c.Answer.ExtraNum {
		err = ecode.RequestErr
		return
	}
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	ah, err := s.history(c, mid)
	if err != nil || ah == nil || ah.StartTime.Add(s.answerDuration()).Before(now) || ah.Score != 0 {
		err = ecode.AnswerBaseNotPassed
		return
	}
	if len(ids) != (s.c.Answer.BaseExtraPassNum + s.c.Answer.BaseExtraNoPassNum) {
		return ecode.AnswerQsNumErr
	}
	passids, err := s.answerDao.IdsCache(c, mid, model.BaseExtraPassQ)
	if err != nil {
		log.Error("s.answerDao.IdsCache(%d) extra pass err(%v) ", mid, err)
		return ecode.AnswerMidCacheQidsErr
	}
	nopassids, err := s.answerDao.IdsCache(c, mid, model.BaseExtraNoPassQ)
	if err != nil {
		log.Error("s.answerDao.IdsCache(%d) extra nopass err(%v) ", mid, err)
		return ecode.AnswerMidCacheQidsErr
	}
	idsmc := append(passids, nopassids...)
	ok, err := s.checkQsIDs(c, ids, mid, idsmc, s.c.Answer.BaseExtraPassNum+s.c.Answer.BaseExtraNoPassNum)
	if err != nil || !ok {
		return
	}
	ret, qs, _ := s.checkExtraPassAns(c, mid, passids, ansHash, lang, s.c.Answer.BaseExtraPassNum)
	ah.StepExtraScore = int64(ret * s.c.Answer.BaseExtraScore)
	ah.StepExtraCompleteTime = now.Unix() - ah.StartTime.Unix()
	if _, err = s.answerDao.UpdateExtraRet(c, ah.ID, mid, ah.StepExtraCompleteTime, ah.StepExtraScore); err != nil {
		log.Error("s.answerDao.UpdateExtraRet(%d) err(%v) ", mid, err)
		return
	}
	ah.Mtime = now
	s.userActionLog(mid, model.ExtraCheck, ah)
	s.answerDao.DelHistoryCache(c, mid)
	s.answerDao.DelIdsCache(c, mid, model.BaseExtraPassQ)
	s.answerDao.DelIdsCache(c, mid, model.BaseExtraNoPassQ)
	// send answer ret to bigdata
	rs, err := s.sendExtraRetMsg(c, mid, qs, nopassids, ansHash, s.c.Answer.BaseExtraNoPassNum)
	if err != nil {
		log.Error("s.sendExtraRetMsg(%d,%v,%v,%v) err(%v) ", mid, qs, nopassids, ansHash, err)
		return
	}
	s.answerDao.PubExtraRet(c, mid, &model.DataBusResult{
		Mid:   mid,
		Buvid: buvid,
		IP:    metadata.String(c, metadata.RemoteIP),
		Ua:    ua,
		Refer: refer,
		Score: int8(ah.StepExtraScore),
		Rs:    rs,
		Hid:   ah.Hid,
	})
	return
}

// checkExtraPassAns check extra question ans.
func (s *Service) checkExtraPassAns(c context.Context, mid int64, ids []int64, ansHash map[int64]string, lang string, count int) (ret int, qs map[int64]*model.ExtraQst, err error) {
	qs, err = s.answerDao.ExtraByIds(c, ids)
	if err != nil || qs == nil || len(qs) != count {
		log.Error("checkAns extra qs is nil (%v,%v) error(%v)", ids, qs, err)
		err = ecode.AnswerMidDBQueErr
		return
	}
	for _, q := range qs {
		var ans string
		switch q.Ans {
		case model.NormalQ:
			if lang == model.LangZhTW {
				ans = s.ansHash(mid, chinese.Convert(c, model.ExtraAnsA))
			} else {
				ans = s.ansHash(mid, model.ExtraAnsA)
			}
		case model.ViolationQ:
			if lang == model.LangZhTW {
				ans = s.ansHash(mid, chinese.Convert(c, model.ExtraAnsB))
			} else {
				ans = s.ansHash(mid, model.ExtraAnsB)
			}
		}
		if ansHash[q.ID] == ans {
			ret++
		}
	}
	return
}

func (s *Service) sendExtraRetMsg(c context.Context, mid int64, passqs map[int64]*model.ExtraQst, nopassids []int64,
	ansHash map[int64]string, count int) (rs []*model.Rs, err error) {
	var (
		qs map[int64]*model.ExtraQst
	)
	qs, err = s.answerDao.ExtraByIds(c, nopassids)
	if err != nil || qs == nil || len(qs) != count {
		log.Error("checkAns extra nopassqs is nil (%v) error(%v)", qs, err)
		err = ecode.AnswerMidDBQueErr
		return
	}
	for k, v := range passqs {
		qs[k] = v
	}
	for _, q := range qs {
		var (
			userAns int8
		)
		ansA := s.ansHash(mid, model.ExtraAnsA)
		ansB := s.ansHash(mid, model.ExtraAnsB)
		switch ansHash[q.ID] {
		case ansA:
			userAns = model.NormalQ
		case ansB:
			userAns = model.ViolationQ
		default:
			userAns = model.UnKownQ
		}
		rs = append(rs, &model.Rs{
			ID:       q.OriginID,
			Question: q.Question,
			Ans:      userAns,
			TrueAns:  q.Ans,
			AvID:     q.AvID,
			Status:   q.Status,
			Source:   q.Source,
			Ctime:    q.Ctime,
			Mtime:    q.Mtime,
		})
	}
	return
}

// PendantRec .
func (s *Service) PendantRec(c context.Context, arg *model.ReqPendant) (err error) {
	var (
		ok       bool
		status   int8
		hid, ret int64
		his      *model.AnswerHistory
	)
	if hid, status, err = s.answerDao.PendantHistory(c, arg.MID); err != nil {
		return
	}
	if hid != arg.HID {
		log.Warn("mid(%d) arg.hid(%d) db.hid(%d) is invald!", arg.MID, arg.HID, hid)
		return
	}
	if status != model.PendantNotGet {
		log.Warn("mid(%d) hid(%d) not first get!", arg.MID, arg.HID)
		return
	}
	his, err = s.historyByHid(c, arg.HID)
	if err != nil {
		return
	}
	if his.Score < model.Score60 || his.IsFirstPass != 1 {
		log.Warn("mid(%d) hid(%d) score(%d) isFirstPass(%d) not pass or first answer !", arg.MID, arg.HID, his.Score, his.IsFirstPass)
		return
	}
	if _, ok = _pendantIDNameMap[int(his.RankID)]; !ok {
		log.Warn("mid(%d) get illegal pid(%d) by answer first!", arg.MID, int(his.RankID))
		return
	}
	if ret, err = s.answerDao.UpPendantHistory(c, arg.MID, arg.HID); err != nil {
		return
	}
	if ret <= 0 {
		log.Warn("mid(%d) hid(%d) pid(%d) history answer not get!", arg.MID, arg.HID, int(his.RankID))
		return
	}
	s.missch.Do(c, func(ctx context.Context) {
		if pendantErr := s.accountDao.GivePendant(ctx, arg.MID, int64(his.RankID), model.PenDantDays, metadata.String(c, metadata.RemoteIP)); pendantErr != nil {
			log.Error("s.accountDao.GivePendant(%d,%d) error(%+v)", arg.MID, int64(his.RankID), pendantErr)
		}
	})
	return
}
