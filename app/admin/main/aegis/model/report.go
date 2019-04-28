package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/aegis/model/common"
	modtask "go-common/app/admin/main/aegis/model/task"
	"go-common/library/log"
	xtime "go-common/library/time"
)

//hash fields
const (
	Dispatch = "ds"
	Delay    = "dy"
	Submit   = "st_%d_%d" // 参数1:提交状态(任务提交，资源提交，任务关闭) 参数2:提交前属于谁
	Release  = "rl"
	RscState = "rs_%d"
	UseTime  = "ut"
	SetKey   = "report_set"

	In  = "in"
	Out = "out"

	//type
	TypeMeta          = int8(0)
	TypeTotal         = int8(1)
	TypeCreate        = int8(2) //每日任务增量
	TypeDaySubmit     = int8(3) //每日任务提交信息
	TypeDayUserSubmit = int8(4) //每日审核人员提交信息
)

//TaskReport 任务报表
type TaskReport struct {
	ID         int64     `json:"id"`
	BusinessID int64     `json:"business_id"`
	FlowID     int64     `json:"flow_id"`
	Type       int8      `json:"type"`
	UID        int64     `json:"uid" gorm:"column:uid"`
	StatDate   time.Time `json:"stat_date"`
	Content    string    `json:"content"`
}

func (t *TaskReport) TableName() string {
	return "task_report"
}

//DayPager 按天分页
type DayPager struct {
	Pn      string `json:"pn"`
	IsFirst bool   `json:"is_first"` //第一页
	IsLast  bool   `json:"is_last"`  //最后一页
}

//ReportSubmitRes .
type ReportSubmitRes struct {
	Order  []string              `json:"order"`  //表格列顺序
	Header map[string]string     `json:"header"` //表格头
	Rows   [][]map[string]string `json:"rows"`   //表格内容体
}

//ReportSubmitItem .
type ReportSubmitItem struct {
	StatDate  string                        `json:"stat_date"`
	CreateCnt int64                         `json:"create_cnt"`
	UID       int64                         `json:"uid"`
	Username  string                        `json:"username"`
	AvgDur    string                        `json:"avg_dur"`
	AvgUtime  string                        `json:"avg_utime"`
	FlowID    int64                         `json:"flow_id"`
	SubmitCnt int64                         `json:"submit_cnt"`
	StateCnt  map[string]*ReportSubmitState `json:"state_cnt"`
}

//ReportSubmitState .
type ReportSubmitState struct {
	Cnt   int64  `json:"cnt"`
	Ratio string `json:"ratio"`
}

type ReportSubmitContent struct {
	RscStates  map[string]int64 `json:"rsc_states"`
	AvgDur     int64            `json:"avg_dur"`
	AvgUtime   int64            `json:"avg_utime"`
	TotalDur   int64            `json:"omitempty,total_dur"`
	TotalUtime int64            `json:"omitempty,total_utime"`
}

//OptReport 参数
type OptReport struct {
	BizID  int64  `form:"business_id" validate:"required"`
	FlowID int64  `form:"flow_id"`
	UName  string `form:"username"`
	Bt     string `form:"bt"`
	Btime  time.Time
	Et     string `form:"et"`
	Etime  time.Time
	Type   int8 `form:"type"`
}

//OptReportSubmit 参数
type OptReportSubmit struct {
	BizID  int64  `form:"business_id" validate:"required"`
	FlowID int64  `form:"flow_id"`
	Bt     string `form:"bt"`
	Et     string `form:"et"`
	Pn     string `form:"pn"`
}

//ReportFlowItem 流量
type ReportFlowItem struct {
	Hour string `json:"hour"`
	In   int64  `json:"in"`
	Out  int64  `json:"out"`
}

//MetaItem .
type MetaItem struct {
	DS int64 // 领取
	RL int64 // 释放
	RP int64 // 通过
	ST int64 // 提交
	UT int64 // 时长
}

//ReportMeta .
type ReportMeta struct {
	Mtime   xtime.Time
	Content string
	UID     int64
	Type    int8
	Uname   string
}

//Genform 元数据的加工函数,进出审核表格
func Genform(res map[string][]*ReportFlowItem) (form map[string][24]*ReportFlowItem) {
	form = make(map[string][24]*ReportFlowItem)

	date := func(hour string) (d string, h string) {
		d, h = hour[:10], hour[11:]
		return
	}

	emptyflow := func() (res [24]*ReportFlowItem) {
		res = [24]*ReportFlowItem{}
		for i := 0; i < 24; i++ {
			res[i] = &ReportFlowItem{Hour: fmt.Sprintf("%.2d", i)}
		}
		return res
	}

	mapcheck := make(map[string]int64) //完全没数据的就不展示了
	for uname, arr := range res {
		for _, item := range arr {
			d, h := date(item.Hour)
			hi, _ := strconv.Atoi(h)
			key := fmt.Sprintf("%s %s", uname, d)
			if arr, ok := form[key]; ok {
				arr[hi].In = item.In
				arr[hi].Out = item.Out
				form[key] = arr
			} else {
				data := emptyflow()
				data[hi].In = item.In
				data[hi].Out = item.Out
				form[key] = data
			}
			mapcheck[key] += (item.In + item.Out)
		}
	}

	for k, v := range mapcheck {
		if v == 0 {
			delete(form, k)
		}
	}

	return
}

//Gentempres 元数据的加工函数,进出审核柱状图1
func Gentempres(opt *OptReport, mnames map[int64]string, metas []*ReportMeta) (tempres map[string]map[int64]*ReportFlowItem) {
	var funcin, funcout caculate
	if opt.Type == TypeMeta {
		funcin = metaIn
		funcout = metaOut
	} else {
		funcin = totalIn
		funcout = totalOut
	}

	tempres = make(map[string]map[int64]*ReportFlowItem)
	for _, meta := range metas {
		metaM := make(map[string]int64)
		if err := json.Unmarshal([]byte(meta.Content), &metaM); err != nil {
			log.Error("json.Unmarshal (%s) error(%v)", meta.Content, err)
			continue
		}

		var in, out int64
		in += funcin(metaM, meta)
		out += funcout(metaM, meta)
		hour := meta.Mtime.Time().Format("2006-01-02 15")

		item := &ReportFlowItem{
			Hour: hour,
			In:   in,
			Out:  out,
		}
		if val, ok := tempres[hour]; ok {
			if v, ok := val[meta.UID]; ok {
				v.In += in
				v.Out += out
				tempres[hour][meta.UID] = v
			} else {
				tempres[hour][meta.UID] = item
			}
		} else {
			tempres[hour] = map[int64]*ReportFlowItem{
				meta.UID: item,
			}
		}
	}
	return
}

type caculate func(map[string]int64, *ReportMeta) int64

func totalIn(metaM map[string]int64, meta *ReportMeta) (in int64) {
	return metaM[In]
}
func totalOut(metaM map[string]int64, meta *ReportMeta) (in int64) {
	return metaM[Out]
}

func metaIn(metaM map[string]int64, meta *ReportMeta) (in int64) {
	if ds, ok := metaM[Dispatch]; ok { //领取到的
		in += ds
	}
	if rl, ok := metaM[Release]; ok { //释放的
		in -= rl
	}
	return
}

func metaOut(metaM map[string]int64, meta *ReportMeta) (out int64) {
	if st, ok := metaM[fmt.Sprintf(Submit, modtask.TaskStateSubmit, meta.UID)]; ok { //本人提交的
		out = st
	}
	return
}

//Genres 元数据的加工函数,进出审核柱状图2
func Genres(opt *OptReport, tempres map[string]map[int64]*ReportFlowItem, mnames map[int64]string) (res map[string][]*ReportFlowItem) {
	res = make(map[string][]*ReportFlowItem)
	for i := opt.Btime; i.Format("2006-01-02 15") != opt.Etime.Format("2006-01-02 15"); i = i.Add(+1 * time.Hour) {
		hour := i.Format("2006-01-02 15")
		rt := &ReportFlowItem{
			Hour: hour,
		}
		if val, ok := tempres[hour]; ok {
			var tin, tout int64
			if opt.Type == TypeMeta {
				for uid, uname := range mnames {
					if item, ok := val[uid]; ok {
						rt = item
						tin += item.In
						tout += item.Out
					} else {
						rt.In = 0
						rt.Out = 0
					}
					addtores(res, uname, *rt)
				}
			} else {
				if item, ok := val[0]; ok {
					rt = item
					tin += item.In
					tout += item.Out
				}
			}

			rt.In = tin
			rt.Out = tout
			addtores(res, "total", *rt)
		} else {
			addtores(res, "total", *rt)
			for _, uname := range mnames {
				addtores(res, uname, *rt)
			}
		}
	}

	for k, v := range res {
		if opt.UName != "" && k != opt.UName {
			delete(res, k)
			continue
		}
		var (
			insum, outsum = int64(0), int64(0)
		)
		for _, item := range v {
			insum += item.In
			outsum += item.Out
		}
		if insum == 0 && outsum == 0 {
			delete(res, k)
		}
	}
	return
}

func addtores(res map[string][]*ReportFlowItem, key string, val ReportFlowItem) {
	if total, ok := res[key]; ok {
		res[key] = append(total, &val)
	} else {
		res[key] = []*ReportFlowItem{&val}
	}
}

//GenMemberStat 个人24小时统计报表
func GenMemberStat(metas []*ReportMeta, passState int64) (res map[int64][]interface{}, err error) {
	tempres := make(map[int64]*MetaItem)
	for _, meta := range metas {
		metaM := make(map[string]int64)
		if err = json.Unmarshal([]byte(meta.Content), &metaM); err != nil {
			log.Error("json.Unmarshal (%s) error(%v)", meta.Content, err)
			err = nil
			continue
		}

		tmt := &MetaItem{
			DS: metaM[Dispatch],
			RL: metaM[Release],
			RP: metaM[fmt.Sprintf(RscState, passState)],
			ST: metaM[fmt.Sprintf(Submit, modtask.TaskStateSubmit, meta.UID)],
			UT: metaM[UseTime],
		}

		if val, ok := tempres[meta.UID]; ok {
			val.DS += tmt.DS
			val.RL += tmt.RL
			val.RP += tmt.RP
			val.ST += tmt.ST
			val.UT += tmt.UT
		} else {
			tempres[meta.UID] = tmt
		}
	}

	for uid, v := range tempres {
		fmt.Printf("uid:%d v:%+v\n", uid, v)
	}
	res = make(map[int64][]interface{})
	for uid, meta := range tempres {
		var (
			completerate, passrate, avgtime = "", "", ""
		)
		total := meta.DS - meta.RL
		if total <= 0 {
			total = 0
			completerate, passrate, avgtime = "0%", "0%", "00:00:00"
		} else {
			cRate := float32(meta.ST) / float32(total)
			pRate := float32(meta.RP) / float32(total)
			at := int64(0)
			if meta.ST > 0 {
				at = meta.UT / meta.ST
			}
			if cRate < 0.0 || cRate > 1.0 {
				completerate = "100%"
			} else {
				completerate = fmt.Sprintf("%.2f%%", cRate*100)
			}
			if pRate < 0.0 || pRate > 1.0 {
				passrate = "100%"
			} else {
				passrate = fmt.Sprintf("%.2f%%", pRate*100)
			}
			if at < 0 {
				avgtime = "00:00:00"
			} else {
				avgtime = common.ParseWaitTime(at)
			}
		}
		res[uid] = []interface{}{total, completerate, passrate, avgtime}
	}
	return
}
