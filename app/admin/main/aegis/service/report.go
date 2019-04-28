package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

//ReportTaskflow .
func (s *Service) ReportTaskflow(c context.Context, opt *model.OptReport) (res map[string][]*model.ReportFlowItem, form map[string][24]*model.ReportFlowItem, err error) {
	// TODO 从redis 补充实时增量
	var (
		mnames map[int64]string
		metas  []*model.ReportMeta
	)

	// 取出统计元数据
	mnames, metas, err = s.reportMeta(c, opt)
	if err != nil {
		return
	}

	tempres := model.Gentempres(opt, mnames, metas)
	res = model.Genres(opt, tempres, mnames)
	form = model.Genform(res)
	return
}

//MemberStats . 统计个人24h 处理量 处理率 通过率 平均耗时
func (s *Service) MemberStats(c context.Context, bizid, flowid int64, uids []int64) (res map[int64][]interface{}, err error) {
	var (
		metas []*model.ReportMeta
		opt   = &model.OptReport{
			BizID:  bizid,
			FlowID: flowid,
			Type:   model.TypeMeta,
		}
		passState = s.passState(bizid, flowid)
	)

	if _, metas, err = s.reportMeta(c, opt); err != nil {
		return
	}

	return model.GenMemberStat(metas, passState)
}

//MemberStats . 统计个人24h 处理量 处理率 通过率 平均耗时

func (s *Service) parseOption(c context.Context, opt *model.OptReport) (uids []int64, mnames map[int64]string, err error) {
	var muids map[string]int64
	mnames = make(map[int64]string)

	bTime, _ := time.ParseInLocation("2006-01-02", opt.Bt, time.Local)
	eTime, _ := time.ParseInLocation("2006-01-02", opt.Et, time.Local)
	switch {
	case bTime.IsZero() && eTime.IsZero():
		bTime = time.Now().Add(-24 * time.Hour)
		eTime = time.Now()
	case bTime.IsZero() && !eTime.IsZero():
		bTime = eTime.Add(-24 * time.Hour)
	case !bTime.IsZero() && eTime.IsZero():
		eTime = bTime.Add(24 * time.Hour)
	}

	eTime = eTime.Add(+1 * time.Hour)

	opt.Btime, opt.Etime = bTime, eTime
	opt.Bt, opt.Et = bTime.Format("2006-01-02 15"), eTime.Format("2006-01-02 15")

	if bTime.After(eTime) {
		err = ecode.RequestErr
		return
	}

	if len(opt.UName) > 0 {
		if muids, err = s.http.GetUIDs(c, opt.UName); err != nil {
			return
		}
		for name, uid := range muids {
			uids = append(uids, uid)
			mnames[uid] = name
		}
	}

	return
}

func (s *Service) reportMeta(c context.Context, opt *model.OptReport) (mnames map[int64]string, metas []*model.ReportMeta, err error) {
	// TODO 从redis 补充实时增量
	var (
		uids []int64
	)

	uids, mnames, err = s.parseOption(c, opt)
	if err != nil {
		return
	}

	// 取出统计元数据
	metas, missuids, err := s.gorm.ReportTaskMetas(c, opt.Bt, opt.Et, opt.BizID, opt.FlowID, uids, mnames, opt.Type)
	if err != nil {
		return
	}
	if len(missuids) > 0 && opt.Type == model.TypeMeta {
		var tunames map[int64]string
		if tunames, err = s.http.GetUnames(c, missuids); err != nil {
			return
		}
		for uid, uname := range tunames {
			mnames[uid] = uname
		}
	}
	return
}

//TODO 根据业务配置,获取通过状态号
func (s *Service) passState(bizid, flowid int64) int64 {
	return 1
}

//ReportTaskSubmit 业务节点下的任务提交统计信息, 按天分页
func (s *Service) ReportTaskSubmit(c context.Context, pm *model.OptReportSubmit) (res *model.ReportSubmitRes, err error) {
	var (
		datas            []*model.TaskReport
		states           map[string]string
		unames           map[int64]string
		uids             = []int64{}
		dayUID           = int64(-1)
		dayUname         = "ALL"
		groups           = [][]*model.ReportSubmitItem{} //{stat_date:[]obj}
		creates          = map[string]string{}           //{bizid_flowid_statdate:cnt}
		stateOrderUnique = map[string]string{}
	)

	res = &model.ReportSubmitRes{}
	//取出元数据: 状态名 + 操作提交数据
	if datas, err = s.gorm.TaskReports(c, pm.BizID, pm.FlowID,
		[]int8{model.TypeCreate, model.TypeDaySubmit, model.TypeDayUserSubmit},
		pm.Bt, pm.Et); err != nil || len(datas) == 0 {
		return
	}
	if states, err = s.TokenByName(c, pm.BizID, "state"); err != nil {
		return
	}

	prevDate := ""
	groupItems := []*model.ReportSubmitItem{}
	for _, item := range datas {
		statdate := item.StatDate.Format("2006-01-02")
		//每日任务增量
		if item.Type == model.TypeCreate {
			creates[fmt.Sprintf("%d_%s", item.FlowID, statdate)] = item.Content
			continue
		}

		//每日uid维度的提交数据
		if item.Type == model.TypeDaySubmit {
			item.UID = dayUID
		} else {
			uids = append(uids, item.UID)
		}

		content := &model.ReportSubmitContent{}
		if err = json.Unmarshal([]byte(item.Content), content); err != nil {
			log.Error("ReportTaskSubmit json.Unmarshal(%s) error(%v), for id(%d)", item.Content, err, item.ID)
			continue
		}
		submitCnt := int64(0)
		statecnt := map[string]*model.ReportSubmitState{}
		for state, cnt := range content.RscStates {
			submitCnt = submitCnt + cnt
			statecnt[state] = &model.ReportSubmitState{
				Cnt: cnt,
			}
			if _, exist := stateOrderUnique[state]; !exist {
				stateOrderUnique[state] = states[state]
			}
		}
		//各状态的提交量所占比率
		for _, one := range statecnt {
			one.Ratio = fmt.Sprintf("%.2f", float64(one.Cnt)/float64(submitCnt))
		}
		if prevDate == "" {
			prevDate = statdate
		}
		if prevDate != statdate {
			groups = append(groups, groupItems)
			groupItems = []*model.ReportSubmitItem{}
			prevDate = statdate
		}
		groupItems = append(groupItems, &model.ReportSubmitItem{
			StatDate:  statdate,
			UID:       item.UID,
			AvgDur:    time.Unix(content.AvgDur, 0).In(time.UTC).Format("15:04:05"),
			AvgUtime:  time.Unix(content.AvgUtime, 0).In(time.UTC).Format("15:04:05"),
			FlowID:    item.FlowID,
			SubmitCnt: submitCnt,
			StateCnt:  statecnt,
		})
	}
	if len(groupItems) > 0 {
		groups = append(groups, groupItems)
	}

	//get usernames
	unames, _ = s.http.GetUnames(c, uids)

	res.Order = []string{
		"state_date",
		"create_cnt",
		"username",
		"avg_dur",
		"avg_utime",
		"submit_cnt",
	}
	res.Header = map[string]string{
		"state_date": "日期",
		"create_cnt": "新增量",
		"username":   "操作人",
		"avg_dur":    "平均耗时",
		"avg_utime":  "平均停留时间",
		"submit_cnt": "操作总量",
	}
	res.Rows = make([][]map[string]string, len(groups))
	for state, name := range stateOrderUnique {
		cnt := fmt.Sprintf("state_%s_cnt", state)
		ratio := fmt.Sprintf("state_%s_ratio", state)
		res.Header[cnt] = fmt.Sprintf("%s量", name)
		res.Header[ratio] = fmt.Sprintf("%s率", name)
		res.Order = append(res.Order, cnt, ratio)
	}

	for i, list := range groups {
		res.Rows[i] = make([]map[string]string, len(list))
		for j, item := range list {
			uname := strconv.FormatInt(item.UID, 10)
			create := "0"
			if item.UID == dayUID {
				create = creates[fmt.Sprintf("%d_%s", item.FlowID, item.StatDate)]
				uname = dayUname
			} else if v, exist := unames[item.UID]; exist {
				uname = v
			}

			one := map[string]string{
				"state_date": item.StatDate,
				"create_cnt": create,
				"username":   uname,
				"avg_dur":    item.AvgDur,
				"avg_utime":  item.AvgUtime,
				"submit_cnt": strconv.FormatInt(item.SubmitCnt, 10),
			}
			for state := range stateOrderUnique {
				cnt := "0"
				ration := "0.00"
				if v, exist := item.StateCnt[state]; exist {
					cnt = strconv.FormatInt(v.Cnt, 10)
					ration = v.Ratio
				}
				one[fmt.Sprintf("state_%s_cnt", state)] = cnt
				one[fmt.Sprintf("state_%s_ratio", state)] = ration
			}
			res.Rows[i][j] = one
		}
	}
	return
}
