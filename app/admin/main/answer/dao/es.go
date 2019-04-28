package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"go-common/app/admin/main/answer/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// Page search page info
type Page struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// SearchSubjectLog get subject logs
type SearchSubjectLog struct {
	Page   *Page
	Result []*struct {
		MID       int64  `json:"mid"`
		Business  string `json:"business"`
		Action    string `json:"action"`
		ExtraData string `json:"extra_data"`
		Ctime     string `json:"ctime"`
	}
}

// const .
const (
	answerLogID  = 15
	answerUpdate = "answer_update"

	basePass       = "basePass"
	extraStartTime = "extraStartTime"
	extraCheck     = "extraCheck"
	proQues        = "proQues"
	proCheck       = "proCheck"
	captcha        = "captchaPass"
	level          = "level"
)

// HistoryES search archives by es.
func (d *Dao) HistoryES(c context.Context, arg *model.ArgHistory) (reply []*model.AnswerHistoryDB, err error) {
	var index []string
	for year := time.Now().Year(); year >= 2018; year-- { // 2018.12 开始接入日志
		index = append(index, fmt.Sprintf("log_user_action_%d_%d", answerLogID, year))
	}
	r := d.es.NewRequest("log_user_action").Index(index...)
	r.Fields("mid", "action", "ctime", "extra_data")
	r.WhereEq("mid", arg.Mid)
	r.Pn(arg.Pn).Ps(arg.Ps)
	res := &SearchSubjectLog{}
	if err = r.Scan(c, &res); err != nil || res == nil {
		log.Error("HistoryES:Scan params(%s) error(%v)", r.Params(), err)
		return
	}
	log.Error("HistoryES:%+v", res)
	mmh := make(map[int64]map[string]*model.AnswerHistory)
	for _, v := range res.Result {
		mh := make(map[string]*model.AnswerHistory)
		if err = json.Unmarshal([]byte(v.ExtraData), &mh); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", v.ExtraData, err)
			return
		}
		for k, v := range mh {
			if mmh[v.Hid] == nil {
				mmh[v.Hid] = make(map[string]*model.AnswerHistory)
			}
			mmh[v.Hid][k] = v
		}
	}
	data := make([]*model.AnswerHistory, 0)
	for _, v := range mmh {
		if v == nil {
			continue
		}
		h := new(model.AnswerHistory)
		if eh, ok := v[basePass]; ok && eh != nil {
			h = eh
		}
		if eh, ok := v[extraStartTime]; ok && eh != nil {
			h.StepExtraStartTime = eh.StepExtraStartTime
			h.Mtime = eh.Mtime
		}
		if eh, ok := v[extraCheck]; ok && eh != nil {
			h.StepExtraCompleteTime = eh.StepExtraCompleteTime
			h.Score = eh.Score
			h.StepExtraScore = eh.StepExtraScore
			h.Mtime = eh.Mtime
		}
		if eh, ok := v[proQues]; ok && eh != nil {
			h.StepTwoStartTime = eh.StepTwoStartTime
			h.Mtime = eh.Mtime
		}
		if eh, ok := v[proCheck]; ok && eh != nil {
			h.CompleteResult = eh.CompleteResult
			h.CompleteTime = eh.CompleteTime
			h.Score = eh.Score
			h.IsFirstPass = eh.IsFirstPass
			h.RankID = eh.RankID
			h.Mtime = eh.Mtime
		}
		if eh, ok := v[captcha]; ok && eh != nil {
			h.IsPassCaptcha = eh.IsPassCaptcha
			h.Mtime = eh.Mtime
		}
		if eh, ok := v[level]; ok && eh != nil {
			h.IsFirstPass = eh.IsFirstPass
			h.PassedLevel = eh.PassedLevel
			h.Mtime = eh.Mtime
		}
		data = append(data, h)
	}
	sort.Sort(model.Histories(data))
	return convent(data), nil
}

func convent(src []*model.AnswerHistory) []*model.AnswerHistoryDB {
	res := make([]*model.AnswerHistoryDB, 0, len(src))
	for _, s := range src {
		r := &model.AnswerHistoryDB{
			ID:                    s.ID,
			Hid:                   s.Hid,
			Mid:                   s.Mid,
			StartTime:             xtime.Time(s.StartTime.Unix()),
			StepOneErrTimes:       s.StepOneErrTimes,
			StepOneCompleteTime:   s.StepOneCompleteTime,
			StepExtraStartTime:    xtime.Time(s.StepExtraStartTime.Unix()),
			StepExtraCompleteTime: s.StepExtraCompleteTime,
			StepExtraScore:        s.StepExtraScore,
			StepTwoStartTime:      xtime.Time(s.StepTwoStartTime.Unix()),
			CompleteTime:          xtime.Time(s.CompleteTime.Unix()),
			CompleteResult:        s.CompleteResult,
			Score:                 s.Score,
			IsFirstPass:           s.IsFirstPass,
			IsPassCaptcha:         s.IsPassCaptcha,
			PassedLevel:           s.PassedLevel,
			RankID:                s.RankID,
			Ctime:                 xtime.Time(s.Ctime.Unix()),
			Mtime:                 xtime.Time(s.Mtime.Unix()),
		}
		res = append(res, r)
	}
	return res
}
