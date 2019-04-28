package v1

import (
	"go-common/app/service/live/xuser/model/exp"
	expm "go-common/app/service/live/xuser/model/exp"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	gtime "go-common/library/time"
	"strconv"
	"time"
)

// consts
const (
	_liveUserExpID = 104 // http://info.bilibili.co/pages/viewpage.action?pageId=8731603
)

func (s *UserExpService) addUserExpLog(uid int64, uexp int64, level map[int64]*expm.LevelInfo, logDesc string) (err error) {
	if len(level) > 0 {
		for _, v := range level {
			logParams := &exp.ModelExpLog{
				Mid:        uid,
				Uexp:       uexp,
				Rexp:       0,
				Ts:         time.Now().Unix(),
				ReqBizDesc: logDesc,
				Content: map[string]string{
					"经验变更类型":    "增加用户维度经验",
					"变更数量":      strconv.FormatInt(uexp, 10),
					"当前用户维度总经验": strconv.FormatInt(v.UserLevel.UserExp, 10),
					"当前主播维度总经验": strconv.FormatInt(v.AnchorLevel.UserExp, 10),
					"当前用户维度等级":  strconv.FormatInt(v.UserLevel.Level, 10),
					"当前主播维度等级":  strconv.FormatInt(v.AnchorLevel.Level, 10),
				},
			}
			s.addLog(_liveUserExpID, "exp_change", logParams, "增加用户经验")
		}
	}
	return
}

func (s *UserExpService) addAnchorExpLog(uid int64, uexp int64, level map[int64]*expm.LevelInfo, logDesc string) (err error) {
	if len(level) > 0 {
		for _, v := range level {
			logParams := &exp.ModelExpLog{
				Mid:        uid,
				Uexp:       uexp,
				Rexp:       0,
				Ts:         time.Now().Unix(),
				ReqBizDesc: logDesc,
				Content: map[string]string{
					"经验变更类型":    "增加主播维度经验",
					"变更数量":      strconv.FormatInt(uexp, 10),
					"当前用户维度总经验": strconv.FormatInt(v.UserLevel.UserExp, 10),
					"当前主播维度总经验": strconv.FormatInt(v.AnchorLevel.UserExp, 10),
					"当前用户维度等级":  strconv.FormatInt(v.UserLevel.Level, 10),
					"当前主播维度等级":  strconv.FormatInt(v.AnchorLevel.Level, 10),
				},
			}
			s.addLog(_liveUserExpID, "exp_change", logParams, "增加主播经验")
		}
	}
	return
}

func (s *UserExpService) addLog(business int, action string, expInfo *exp.ModelExpLog, desc string) (err error) {
	t := gtime.Time(expInfo.Ts)
	content := make(map[string]interface{}, len(expInfo.Content))
	for k, v := range expInfo.Content {
		content[k] = v
	}
	ui := &report.UserInfo{
		Mid:      expInfo.Mid,
		Platform: desc,
		Build:    0,
		Buvid:    expInfo.Buvid,
		Business: business,
		Type:     0,
		Action:   action,
		Ctime:    t.Time(),
		IP:       expInfo.ReqBizDesc,
		// extra
		Index:   []interface{}{int64(expInfo.Mid), 0, "", "", ""},
		Content: content,
	}
	report.User(ui)
	log.Info("add log to report: userexplog: %+v userinfo: %+v,error(%v)", expInfo, ui)
	return
}

// RecordTimeCostLog  ...
// 记录时间
func (s *UserExpService) RecordTimeCostLog(nowTime int64, desc string) {
	log.Info(desc+"|%d", nowTime)
}

// RecordTimeCost ...
// 记录时间
func (s *UserExpService) RecordTimeCost() (NowTime int64) {
	return time.Now().UnixNano() / 1000000 // 用毫秒
}
