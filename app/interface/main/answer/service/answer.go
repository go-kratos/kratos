package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"go-common/app/interface/main/answer/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/text/translate/chinese"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	_hashSalt = "bilirqeust"
	// _ansURI     = "/answer/img?qs_id=%v&ans1_hash=%s&ans2_hash=%s&ans3_hash=%s&ans4_hash=%s"
	_baseTypeID = 36 // 官方基础题库
	_rankBtn    = 7 * 24 * time.Hour
	_minType    = 3
	_maxType    = 10
)

var (
	_typeIdsMapping = map[int][]int{
		100001: {15, 16, 17},
		100002: {29, 30},
		100003: {12, 13},
		// 22: 21,
		// 24: 23,
		// 35: 31, 36: 31,
		// 32: 30, 33: 30, 34: 30,
		// 29: 28, 37: 28, 7: 28, 8: 28,
		// 41: 5, 42: 5, 43: 5, 44: 5, 45: 5, 46: 5, 47: 5, 48: 5, 49: 5, 50: 5, 51: 5, 52: 5,
	}
	// 对外展示方式1
	_typeMap1 = []*model.TypeInfo{
		{Name: "游戏", Subs: []*model.SubType{
			{ID: 8, Name: "动作射击"},
			{ID: 9, Name: "冒险格斗"},
			{ID: 100003, Name: "策略模拟"},
			// {ID: 13, Name: "策略模拟"},
			{ID: 14, Name: "音乐体育"},
		}},
		{Name: "影视", Subs: []*model.SubType{
			{ID: 15, Name: "纪录片"},
			{ID: 16, Name: "电影"},
			{ID: 17, Name: "电视剧"},
		}},
		{Name: "科技", Subs: []*model.SubType{
			{ID: 18, Name: "军事"},
			{ID: 19, Name: "地理"},
			{ID: 20, Name: "历史"},
			{ID: 21, Name: "文学"},
			{ID: 22, Name: "数学"},
			{ID: 23, Name: "物理"},
			{ID: 24, Name: "化学"},
			{ID: 25, Name: "生物"},
			{ID: 26, Name: "数码科技"},
		}},
		{Name: "动画", Subs: []*model.SubType{
			{ID: 27, Name: "国创"},
			{ID: 28, Name: "番剧"},
		}},
		{Name: "艺术", Subs: []*model.SubType{
			{ID: 100002, Name: "音乐"},
			// {ID: 30, Name: "音乐"},
			{ID: 31, Name: "绘画"},
		}},
		{Name: "流行前线", Subs: []*model.SubType{
			{ID: 32, Name: "娱乐"},
			{ID: 33, Name: "时尚"},
			{ID: 34, Name: "运动"},
		}},
		{Name: "鬼畜", Subs: []*model.SubType{
			{ID: 35, Name: "鬼畜"},
		}},
	}
	// 推荐分区映射
	_recTypeIDMap = map[int]map[string][]int{
		124: {"main_tid": []int{1, 167}, "sub_tid": []int{3, 129}},
		127: {"main_tid": []int{1, 167}, "sub_tid": []int{3, 4}},
		126: {"main_tid": []int{3, 119}, "sub_tid": []int{1}},
		123: {"main_tid": []int{3}, "sub_tid": []int{129, 119}},
		121: {"main_tid": []int{36, 177}, "sub_tid": []int{160}},
		125: {"main_tid": []int{4}, "sub_tid": []int{129}},
		129: {"main_tid": []int{36, 177}, "sub_tid": []int{160}},
		130: {"main_tid": []int{23, 11}, "sub_tid": []int{160}},
		128: {"main_tid": []int{119}, "sub_tid": []int{160}},
	}
)

// BaseQ base question.
func (s *Service) BaseQ(c context.Context, mid int64, lang string, mobile bool) (res *model.AnsQueDetailList, err error) {
	var aqs *model.AnsQuesList
	if aqs, err = s.BaseQs(c, mid, lang, mobile); err != nil {
		err = errors.Wrapf(err, "s.ansRPC.BaseQs(%d,%t)", mid, mobile)
		return
	}
	res = s.convertModel(aqs)
	return
}

// BaseQs get base question
func (s *Service) BaseQs(c context.Context, mid int64, lang string, mobile bool) (rqs *model.AnsQuesList, err error) {
	var (
		ids []int64
		now = time.Now()
	)
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	h, err := s.history(c, mid)
	if err == nil && h != nil {
		if h.StartTime.Add(s.answerDuration()).After(now) && h.Score == 0 {
			if h.StepExtraCompleteTime != 0 {
				err = ecode.AnswerProNoPass // extra question pass
				return
			}
			err = ecode.AnswerExtraNoPass
			return
		}
		if h.Score > 0 && h.IsPassCaptcha == 0 {
			err = ecode.AnswerCaptchaNoPassed
			return
		}
	}
	ids, err = s.answerDao.IdsCache(c, mid, model.Q)
	if err != nil || len(ids) != s.c.Answer.BaseNum {
		ids, err = s.answerDao.QidByType(c, _baseTypeID, uint8(s.c.Answer.BaseNum))
		if err != nil {
			log.Error("s.answerDao.QidByType(%d,%d)  error(%v)", _baseTypeID, s.c.Answer.BaseNum, err)
			return
		}
		if len(ids) == 0 {
			err = ecode.AnswerQsNumErr
			log.Error("qidByType ids len(%d) is 0", len(ids))
			return
		}
	}
	rqs, err = s.concatData(c, mid, ids, lang, mobile, s.c.Answer.BaseNum)
	if err != nil {
		log.Error("BaseQs s.concatData(%d, %d, %d) error(%v)", c, mid, ids, err)
		return
	}
	at := &model.AnswerTime{
		Stime:  now,
		Etimes: 0,
	}
	err = s.answerDao.SetExpireCache(c, mid, at)
	rqs.CurrentTime = at.Stime
	rqs.EndTime = at.Stime.Add(s.answerDuration())
	s.answerDao.DelIdsCache(c, mid, model.BaseExtraPassQ)
	s.answerDao.DelIdsCache(c, mid, model.BaseExtraNoPassQ)
	return
}

// ConvertExtraQs extra question.
func (s *Service) ConvertExtraQs(c context.Context, mid int64, lang string, mobile bool) (res *model.AnsQueDetailList, err error) {
	var ans *model.AnsQuesList
	if ans, err = s.ExtraQs(c, mid, lang, mobile); err != nil {
		err = errors.Wrapf(err, "s.ansRPC.ExtraQues(%d,%t)", mid, mobile)
		return
	}
	res = s.convertExtraModel(ans)
	return
}

// ExtraQs extra question.
func (s *Service) ExtraQs(c context.Context, mid int64, lang string, mobile bool) (rqs *model.AnsQuesList, err error) {
	var (
		ids      []int64
		passids  []int64
		npassids []int64
		now      = time.Now()
	)
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	h, err := s.checkExtraState(c, mid, now)
	if err != nil {
		return
	}
	// keep on answer
	passids, _ = s.answerDao.IdsCache(c, mid, model.BaseExtraPassQ)
	npassids, err = s.answerDao.IdsCache(c, mid, model.BaseExtraNoPassQ)
	if err != nil || len(passids) != s.c.Answer.BaseExtraPassNum || len(npassids) != s.c.Answer.BaseExtraNoPassNum {
		var (
			ok bool
		)
		ok, passids, npassids = s.extraQueByBigData(c, mid, "")
		if !ok {
			// if bigdata get extra mid faild
			passids, err = s.answerDao.ExtraQidByType(c, model.BaseExtraPassQ, uint8(s.c.Answer.BaseExtraPassNum))
			if err != nil {
				log.Error("s.answerDao.ExtraQidByType(%d, %d, %d) error(%v)", model.BaseExtraPassQ, s.c.Answer.BaseExtraPassNum, len(passids), err)
				return
			}
			if len(passids) != s.c.Answer.BaseExtraPassNum {
				err = ecode.AnswerQsNumErr
				log.Warn("passids lenth(%d) neq BaseExtraPassNum(%d)", len(passids), s.c.Answer.BaseExtraPassNum)
				return
			}
			npassids, err = s.answerDao.ExtraQidByType(c, model.BaseExtraNoPassQ, uint8(s.c.Answer.BaseExtraNoPassNum))
			if err != nil {
				log.Error("s.answerDao.ExtraQidByType(%d, %d, %d) error(%v)", model.BaseExtraNoPassQ, s.c.Answer.BaseExtraNoPassNum, len(npassids), err)
				return
			}
			if len(npassids) != s.c.Answer.BaseExtraNoPassNum {
				err = ecode.AnswerQsNumErr
				log.Warn("npassids lenth(%d) neq BaseExtraNoPassNum(%d)", len(npassids), s.c.Answer.BaseExtraNoPassNum)
				return
			}
		}
	}
	ids = append(passids, npassids...)
	rqs, err = s.concatExtraData(c, mid, ids, passids, npassids, lang, mobile, s.c.Answer.BaseExtraPassNum+s.c.Answer.BaseExtraNoPassNum)
	if err != nil {
		log.Error("BaseExtraQs s.concatExtraData(%d, %d, %d) error(%v)", c, mid, ids, err)
		return
	}
	rqs.CurrentTime = now
	rqs.EndTime = h.StartTime.Add(s.answerDuration())
	if _, err = s.answerDao.UpdateExtraStartTime(c, h.ID, mid, now); err != nil {
		log.Error("s.answerDao.UpdateExtraStartTime( %d, %d) error(%v)", h.ID, mid, err)
		return
	}
	h.StepExtraStartTime = now
	h.Mtime = now
	s.userActionLog(mid, model.ExtraStartTime, h)
	s.answerDao.DelHistoryCache(c, mid)
	return
}

func (s *Service) checkExtraState(c context.Context, mid int64, now time.Time) (h *model.AnswerHistory, err error) {
	h, err = s.history(c, mid)
	if err != nil {
		log.Error("s.history(%v) is nil error(%v)", h, err)
		err = ecode.AnswerBaseNotPassed
		return
	}
	if h != nil {
		if h.Score > 0 && h.IsPassCaptcha == 0 {
			err = ecode.AnswerCaptchaNoPassed
			return
		}
		// if base pass
		if h.StartTime.Add(s.answerDuration()).After(now) && h.Score == 0 {
			if h.StepExtraCompleteTime != 0 {
				err = ecode.AnswerProNoPass
			}
			return
		}
		err = ecode.AnswerBaseNotPassed
		return
	}
	err = ecode.AnswerBaseNotPassed
	return
}

// ProTypes get promotion types.
func (s *Service) proTypes(c context.Context, mid int64) (res *model.ProTypes, err error) {
	var (
		repro bool
		ah    *model.AnswerHistory
		now   = time.Now()
	)
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	if ah, err = s.checkBase(c, mid, now); err != nil {
		return
	}
	qsidsMc, err := s.answerDao.IdsCache(c, mid, model.Q)
	if err == nil && len(qsidsMc) == s.c.Answer.ProNum {
		repro = true
	}
	res = &model.ProTypes{List: _typeMap1, EndTime: ah.StartTime.Add(s.answerDuration()), CurrentTime: now, Repro: repro}
	return
}

// ProType type.
func (s *Service) ProType(c context.Context, mid int64, lang string) (res *model.AnsProType, err error) {
	var (
		repro string
		list  = []*model.AnsTypeList{}
	)
	rpcRes, err := s.proTypes(c, mid)
	if err != nil {
		log.Error("s.proTypes(%+d) error (%v)", mid, err)
		return
	}
	for _, vt := range rpcRes.List {
		var sub = []*model.AnsType{}
		for _, vst := range vt.Subs {
			ansType := &model.AnsType{ID: vst.ID, Name: vst.Name}
			if lang == model.LangZhTW {
				ansType.Name = chinese.Convert(c, ansType.Name)
			}
			sub = append(sub, ansType)
		}
		ansTypeList := &model.AnsTypeList{Name: vt.Name, Fields: sub}
		if lang == model.LangZhTW {
			ansTypeList.Name = chinese.Convert(c, ansTypeList.Name)
		}
		list = append(list, ansTypeList)
	}
	if rpcRes.Repro {
		repro = "yes"
	} else {
		repro = "no"
	}
	res = &model.AnsProType{List: list, CurrentTime: rpcRes.CurrentTime.Unix(), EndTime: rpcRes.EndTime.Unix(), Repro: repro}
	return
}

// ConvertProQues pro question.
func (s *Service) ConvertProQues(c context.Context, mid int64, tIds string, lang string, mobile bool) (res []*model.AnsQueDetail, err error) {
	var (
		ans   *model.AnsQuesList
		ansdl *model.AnsQueDetailList
	)
	if ans, err = s.ProQues(c, mid, tIds, lang, mobile); err != nil {
		err = errors.Wrapf(err, "s.ProQues(%d,%v,%t)", mid, tIds, mobile)
		return
	}
	ansdl = s.convertModel(ans)
	res = ansdl.QuesList
	return
}

// ProQues question info.
func (s *Service) ProQues(c context.Context, mid int64, qtsStr string, lang string, mobile bool) (rqs *model.AnsQuesList, err error) {
	var (
		ah             *model.AnswerHistory
		now            = time.Now()
		allQids        []int64
		tIds, realTIDs []int
	)
	if s.checkAnswerBlock(c, mid) {
		err = ecode.AnswerBlock
		return
	}
	if ah, err = s.checkBase(c, mid, now); err != nil {
		return
	}
	allQids, err = s.answerDao.IdsCache(c, mid, model.Q)
	if err != nil || len(allQids) != s.c.Answer.ProNum {
		tIDStrArr := strings.Split(qtsStr, ",")
		if len(tIDStrArr) < _minType || len(tIDStrArr) > _maxType {
			err = ecode.AnswerTypeIDsErr
			return
		}
		if tIds, err = sliceAtoi(tIDStrArr); err != nil {
			err = ecode.AnswerTypeIDsErr
			return
		}
		for _, qt := range tIds {
			if qt <= 0 {
				err = ecode.RequestErr
				return
			}
			if mapIDS, ok := _typeIdsMapping[qt]; ok {
				realTIDs = append(realTIDs, mapIDS...)
				continue
			}
			realTIDs = append(realTIDs, qt)
		}
		num := math.Ceil(float64(s.c.Answer.ProNum) / float64(len(realTIDs)))
		log.Warn("realTIDs:%v", realTIDs)
		for _, qt := range realTIDs {
			var t []int64
			t, err = s.answerDao.QidByType(c, qt, uint8(num))
			if err != nil {
				log.Error("s.answerDao.QidByType(%d, %f, %d) error(%+v)", qt, num, len(t), err)
				return
			}
			if len(t) == 0 {
				log.Error("mid:%d the QidByType(%d, %f, %d) of len is 0", mid, qt, num, len(t))
				err = ecode.AnswerMidDBQueErr
				return
			}
			allQids = append(allQids, t...)
		}
		if len(allQids) == 0 || len(allQids) < s.c.Answer.ProNum {
			log.Error("ProQues allQids len is 0 or allQids len less(%d, %d, %f, %v, %d)", len(allQids), s.c.Answer.ProNum, num, realTIDs, mid)
			err = ecode.NothingFound
			return
		}
	}
	if rqs, err = s.concatData(c, mid, allQids, lang, mobile, s.c.Answer.ProNum); err != nil {
		log.Error("ProQues s.concatData(%d, %d, %d) error(%v)", c, mid, allQids, err)
		return
	}
	if _, err = s.answerDao.UpdateStepTwoTime(c, ah.ID, mid, now); err != nil {
		return
	}
	ah.StepTwoStartTime = now
	ah.Mtime = now
	s.userActionLog(mid, model.ProQues, ah)
	s.answerDao.DelHistoryCache(c, mid)
	return
}

func (s *Service) checkBase(c context.Context, mid int64, now time.Time) (ah *model.AnswerHistory, err error) {
	ah, err = s.history(c, mid)
	if err != nil || ah == nil || ah.StartTime.Add(s.answerDuration()).Before(now) || ah.Score != 0 || ah.StepOneCompleteTime == 0 {
		err = ecode.AnswerBaseNotPassed
		log.Error("checkBase(%d, %v) AnswerExpire error(%v)", mid, now, err)
		return
	}
	if ah.StepExtraCompleteTime == 0 {
		err = ecode.AnswerExtraNoPass
		return
	}
	if ah.Score > 0 && ah.IsPassCaptcha == 0 {
		err = ecode.AnswerCaptchaNoPassed
	}
	return
}

func (s *Service) checkTime(c context.Context, mid int64, now time.Time) (at *model.AnswerTime, rs bool) {
	var err error
	if at, err = s.answerDao.ExpireCache(c, mid); err != nil {
		return
	}
	if at == nil || at.Stime.Add(s.answerDuration()).Before(now) {
		return
	}
	rs = true
	return
}

func (s *Service) concatData(c context.Context, mid int64, ids []int64, lang string, mobile bool, qs int) (rqs *model.AnsQuesList, err error) {
	var (
		list []*model.AnsQue
		qm   map[int64]*model.Question
	)
	if qm, err = s.answerDao.ByIds(c, ids); err != nil {
		log.Error("s.answerDao.ByIds(%v) error(%v)", ids, err)
		err = ecode.NothingFound
		return
	}
	for _, d := range ids {
		i := qm[d]
		rq := s.imgPosition(c, i, mid, lang, mobile)
		list = append(list, rq)
	}
	if len(list) > qs {
		list = list[:qs]
	}
	rqs = &model.AnsQuesList{QuesList: list}
	if err := s.answerDao.SetIdsCache(c, mid, ids, model.Q); err != nil {
		log.Error("s.answerDao.SetIdsCache(%d, %d) error(%v)", mid, ids, err)
	}
	log.Info("s.concatData load que success(%d, %v, %v, %d)", mid, ids, mobile, qs)
	return
}

func (s *Service) concatExtraData(c context.Context, mid int64, ids []int64, passids []int64, nopassids []int64, lang string, mobile bool, qs int) (rqs *model.AnsQuesList, err error) {
	var (
		list []*model.AnsQue
		qm   map[int64]*model.ExtraQst
	)
	if qm, err = s.answerDao.ExtraByIds(c, ids); err != nil || len(qm) < qs {
		log.Error("s.answerDao.ExtraByIds(%v) error(%+v)", ids, err)
		return
	}
	for _, d := range ids {
		i := qm[d]
		rq := s.imgExtraPosition(c, i, mid, lang, mobile)
		list = append(list, rq)
	}
	if len(list) > qs {
		list = list[:qs]
	}
	rqs = &model.AnsQuesList{QuesList: list}
	if err = s.answerDao.SetIdsCache(c, mid, passids, model.BaseExtraPassQ); err != nil {
		log.Error("s.answerDao.SetIdsCache(%d, %d) error(%v)", mid, passids, err)
		return
	}
	if err = s.answerDao.SetIdsCache(c, mid, nopassids, model.BaseExtraNoPassQ); err != nil {
		log.Error("s.answerDao.SetIdsCache(%d, %d) error(%v)", mid, nopassids, err)
		return
	}
	log.Info("s.concatData extra load que success(%d, %v, %v, %d)", mid, ids, mobile, qs)
	return
}

// ansHash get answer hash.
func (s *Service) ansHash(mid int64, ans string) (ansHash string) {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%d%s", ans, mid, _hashSalt)))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Service) imgPosition(c context.Context, qs *model.Question, mid int64, lang string, mobile bool) (rq *model.AnsQue) {
	var (
		y                 float64
		qsLineLength      float64 = 36
		questionFontSize  float64 = 10
		questionTitleSize float64 = 12
		ans                       = make([]*model.AnsPosition, 4)
		imgStr                    = "v3_%s_A-%s_B-%s_C-%s_D-%s_%s"
		p                         = map[bool]string{true: "H5", false: "PC"}
		bfsHost                   = "https://i0.hdslb.com/bfs/member/"
		as                [4]string
	)
	rq = &model.AnsQue{ID: qs.ID}
	if mobile {
		qsLineLength = 11
		questionFontSize = 12
		questionTitleSize = 16
	}

	qsLength := utf8.RuneCountInString(qs.Question)
	if float64(qsLength) > qsLineLength {
		line := math.Ceil(float64(qsLength) / qsLineLength)
		rq.Height = 2 * line * questionTitleSize
		rq.PositionY = y
	} else {
		rq.Height = 2 * questionTitleSize
	}
	y = rq.Height

	if lang == model.LangZhTW {
		qs.Question = chinese.Convert(c, qs.Question)
		qs.Ans[0] = chinese.Convert(c, qs.Ans[0])
		qs.Ans[1] = chinese.Convert(c, qs.Ans[1])
		qs.Ans[2] = chinese.Convert(c, qs.Ans[2])
		qs.Ans[3] = chinese.Convert(c, qs.Ans[3])
	}
	idx := rand.Perm(4)
	for i := range qs.Ans {
		ans[i] = &model.AnsPosition{
			AnsHash:   s.ansHash(mid, qs.Ans[idx[i]]),
			Height:    2 * questionFontSize,
			PositionY: y,
		}
		y += 2 * questionFontSize
		as[i] = qs.Ans[idx[i]]
	}
	m := md5.New()
	m.Write([]byte(fmt.Sprintf(imgStr, strconv.FormatInt(qs.ID, 10), as[0], as[1], as[2], as[3], p[mobile])))
	fname := hex.EncodeToString(m.Sum(nil)) + ".jpg"
	if s.c.Answer.Debug {
		fname = fmt.Sprintf("debug_%s", fname)
	}
	rq.Img = bfsHost + fname
	rq.Ans = ans
	return
}

func (s *Service) imgExtraPosition(c context.Context, qs *model.ExtraQst, mid int64, lang string, mobile bool) (rq *model.AnsQue) {
	var (
		y                 float64
		qsLineLength      float64 = 36
		questionFontSize  float64 = 10
		questionTitleSize float64 = 12
		ans                       = make([]*model.AnsPosition, 2)
		imgStr                    = "%s_A-%s_B-%s_%s"
		p                         = map[bool]string{true: "H5", false: "PC"}
		bfsHost                   = "https://i0.hdslb.com/bfs/member/"
		as                [2]string
	)
	rq = &model.AnsQue{ID: qs.ID}
	if mobile {
		qsLineLength = 11
		questionFontSize = 12
		questionTitleSize = 16
	}

	qsLength := utf8.RuneCountInString(qs.Question)
	if float64(qsLength) > qsLineLength {
		line := math.Ceil(float64(qsLength) / qsLineLength)
		rq.Height = 2 * line * questionTitleSize
		rq.PositionY = y
	} else {
		rq.Height = 2 * questionTitleSize
	}
	y = rq.Height
	if lang == model.LangZhTW {
		as = [2]string{chinese.Convert(c, model.ExtraAnsA), chinese.Convert(c, model.ExtraAnsB)}
	} else {
		as = [2]string{model.ExtraAnsA, model.ExtraAnsB}
	}

	for k, v := range as {
		ans[k] = &model.AnsPosition{
			AnsHash:   s.ansHash(mid, v),
			Height:    2 * questionFontSize,
			PositionY: y,
		}
		y += 2 * questionFontSize
	}

	m := md5.New()
	m.Write([]byte(fmt.Sprintf(imgStr, strconv.FormatInt(qs.OriginID, 10), as[0], as[1], p[mobile])))
	fname := hex.EncodeToString(m.Sum(nil)) + ".jpg"
	if s.c.Answer.Debug {
		fname = fmt.Sprintf("debug_%s", fname)
	}
	rq.Img = bfsHost + fname
	rq.Ans = ans
	return
}

func (s *Service) loadQidsCache() {
	qs, err := s.answerDao.QidsByState(context.Background(), model.PassCheck)
	if len(qs) == 0 || err != nil {
		log.Error("s.answerDao.loadQidsCache(%d) size is zero error(%v)", model.PassCheck, err)
	}
	qmap := map[int8][]int64{}
	for _, q := range qs {
		qmap[q.TypeID] = append(qmap[q.TypeID], q.ID)
	}
	for k, v := range qmap {
		s.answerDao.DelQidsCache(context.Background(), int(k))
		s.answerDao.SetQids(context.Background(), v, int(k))
	}
	log.Info("s.answerDao.loadQidsCache suc(%v)", qmap)
}

func (s *Service) loadExtraQidsCache() {
	qs, err := s.answerDao.QidsExtraByState(context.Background(), model.MaxLoadQueSize)
	if len(qs) == 0 || err != nil {
		log.Error("s.answerDao.QidsExtraByState(%d) size is zero error(%v)", model.MaxLoadQueSize, err)
		return
	}
	qmap := map[int8][]int64{}
	for _, q := range qs {
		qmap[q.Ans] = append(qmap[q.Ans], q.ID)
	}
	for k, v := range qmap {
		s.answerDao.DelExtraQidsCache(context.Background(), k)
		s.answerDao.SetExtraQids(context.Background(), v, k)
	}
	log.Info("s.answerDao.loadExtraQidsCache suc(%v)", qmap)
}

// Cool .
func (s *Service) Cool(c context.Context, hid, mid int64) (cool *model.AnsCool, err error) {
	var (
		his   *model.AnswerHistory
		types []*model.TypeInfo
		li    = []*model.CoolPower{
			{Name: "动画", Num: 0},
			{Name: "艺术", Num: 0},
			{Name: "游戏", Num: 0},
			{Name: "科技", Num: 0},
			{Name: "影视", Num: 0},
			{Name: "鬼畜", Num: 0},
		}
		completeResult = make(map[int8]int64)
	)
	his, err = s.historyByHid(c, hid)
	if err != nil {
		return
	}
	cool = &model.AnsCool{
		Score:       his.Score,
		IsSameUser:  his.Mid == mid,
		IsFirstPass: his.IsFirstPass,
		Level:       his.PassedLevel,
		Share:       &model.CoolShare{},
		VideoInfo:   &model.CoolVideo{},
		Rank:        &model.CoolRank{},
	}
	cool.CanShowRankBtn = his.Score >= 85 && his.Mtime.Before(time.Now().Add(_rankBtn))
	r := _pendantIDNameMap[int(his.RankID)]
	if r != "" {
		if rid, ok := _oldPIDToNewMap[his.RankID]; ok {
			his.RankID = rid
		}
		rs := _rankShire[his.RankID]
		idx := rand.Perm(len(rs.VideoArr))
		cool.ViewMore = rs.ViewMore
		cool.Share = rs.Share
		cool.VideoInfo = rs.VideoArr[idx[0]]
		cool.Rank = &model.CoolRank{
			ID:   int(his.RankID),
			Name: r,
			Img:  "https://i0.hdslb.com" + _pendantIDImgMap[his.RankID],
		}
	}
	us, err := s.accInfo(c, his.Mid)
	if err != nil || us == nil {
		log.Error("CheckQueCaptcha accInfo(%d) info is null error(%v)", mid, err)
		return
	}
	cool.Name = us.Name
	cool.Face = us.Face
	if err = json.Unmarshal([]byte(his.CompleteResult), &completeResult); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", his.CompleteResult, err)
		err = nil
	}
	log.Info("hid(%d), completeResult: %v+", hid, completeResult)
	for k := range completeResult {
		for _, t := range s.questionTypeCache {
			if len(t.Subs) != 0 {
				for _, s := range t.Subs {
					if int64(k) == s.ID {
						types = append(types, &model.TypeInfo{ID: s.ID, Name: s.Name, LabelName: s.LabelName})
					}
				}
			}
			if int64(k) == t.ID {
				types = append(types, t)
			}
		}
	}
	log.Info("hid(%d), cool types: %v+", hid, types)
	for _, t := range types {
		for _, p := range li {
			if p.Name == t.LabelName {
				p.Num += completeResult[int8(t.ID)]
			}
		}
	}
	log.Info("hid(%d), power types: %v+", hid, li)
	cool.Powers = append(cool.Powers, li...)
	if _, ok := _recTypeIDMap[int(his.RankID)]; ok {
		cool.MainTids = _recTypeIDMap[int(his.RankID)]["main_tid"]
		cool.SubTids = _recTypeIDMap[int(his.RankID)]["sub_tid"]
	}
	return
}

// ExtraScore .
func (s *Service) ExtraScore(c context.Context, mid int64) (res *model.ExtraScoreReply, err error) {
	res = &model.ExtraScoreReply{}
	h, err := s.history(c, mid)
	if err != nil {
		return
	}
	res.Score = h.StepExtraScore + int64(s.c.Answer.BaseNum)
	return
}

func (s *Service) history(c context.Context, mid int64) (ah *model.AnswerHistory, err error) {
	cok := true
	if ah, err = s.answerDao.HistoryCache(c, mid); err != nil {
		cok = false
		return
	}
	if ah != nil {
		return
	}
	ah, err = s.answerDao.History(c, mid)
	if err != nil {
		return
	}
	if ah != nil && cok {
		s.answerDao.SetHistoryCache(c, mid, ah)
	}
	return
}

func (s *Service) answerDuration() (d time.Duration) {
	return time.Duration(s.c.Answer.Duration) * time.Minute
}

func sliceAtoi(sa []string) ([]int, error) {
	si := make([]int, 0, len(sa))
	for _, a := range sa {
		i, err := strconv.Atoi(a)
		if err != nil {
			return si, err
		}
		si = append(si, i)
	}
	return si, nil
}

func (s *Service) extraQueByBigData(c context.Context, mid int64, ip string) (ok bool, passids []int64, npassids []int64) {
	passids, npassids, err := s.accountDao.ExtraIds(c, mid, ip)
	if err != nil || len(passids) != s.c.Answer.BaseExtraPassNum || len(npassids) != s.c.Answer.BaseExtraNoPassNum {
		return
	}
	ids := append(passids, npassids...)
	if qm, err := s.answerDao.ExtraByIds(c, ids); err != nil || len(qm) != (s.c.Answer.BaseExtraPassNum+s.c.Answer.BaseExtraNoPassNum) {
		log.Error("s.answerDao.ExtraByIds(%v) error(%v)", ids, err)
		return
	}
	ok = true
	return
}

func (s *Service) loadtypes() (t map[int64]*model.TypeInfo) {
	tys, err := s.answerDao.Types(context.Background())
	if err != nil {
		log.Error("s.questionDao.Types error(%v)", err)
		return
	}
	tmp := map[int64]*model.TypeInfo{}
	for _, v := range tys {
		if v.Parentid == 0 && tmp[v.ID] == nil {
			tmp[v.ID] = &model.TypeInfo{ID: v.ID, Name: v.Name, Subs: []*model.SubType{}}
		} else if tmp[v.Parentid] != nil {
			tmp[v.Parentid].Subs = append(tmp[v.Parentid].Subs, &model.SubType{ID: v.ID, Name: v.Name, LabelName: v.LabelName})
		}
	}
	s.questionTypeCache = tmp
	t = tmp
	log.Info("load question type cacheproc success,%v", t)
	return
}
