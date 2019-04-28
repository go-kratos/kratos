package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/app/admin/main/aegis/model/resource"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/xstr"
)

// send to log service
func (s *Service) sendAuditLog(c context.Context, action string, opt *model.SubmitOptions, flowres interface{}, logtype int) (err error) {
	// send
	logData := &report.ManagerInfo{
		Uname:    opt.Uname,
		UID:      opt.UID,
		Business: model.LogBusinessAudit,
		Type:     logtype,
		Oid:      opt.RID,
		Action:   action,
		Ctime:    time.Now(),
		Index:    []interface{}{opt.BusinessID, opt.NewFlowID, opt.TaskID, strconv.Itoa(opt.Result.State)},
		Content: map[string]interface{}{
			"opt":  opt,
			"flow": flowres,
		},
	}
	if err = report.Manager(logData); err != nil {
		log.Error("report.Manager(%+v) error(%v)", logData, err)
	}

	return
}

// send to log service
func (s *Service) sendTaskConsumerLog(c context.Context, action string, opt *common.BaseOptions) (err error) {
	logData := &report.ManagerInfo{
		Uname:    opt.Uname,
		UID:      opt.UID,
		Business: model.LogBusinessTask,
		Type:     model.LogTypeTaskConsumer,
		Oid:      0,
		Action:   action,
		Ctime:    time.Now(),
		Index:    []interface{}{opt.BusinessID, opt.FlowID, opt.Role},
	}
	if err = report.Manager(logData); err != nil {
		log.Error("report.Manager(%+v) error(%v)", logData, err)
	}
	return
}

// send to log service
func (s *Service) sendRscLog(c context.Context, acction string, opt *model.AddOption, res *net.TriggerResult, update interface{}, err error) {
	var rid, flowid int64
	if res != nil {
		flowid = res.NewFlowID
		rid = res.RID
	}

	// send
	logData := &report.ManagerInfo{
		Uname:    "business",
		UID:      399,
		Business: model.LogBusinessResource,
		Type:     model.LogTypeFromAdd,
		Oid:      rid,
		Action:   acction,
		Ctime:    time.Now(),
		Index:    []interface{}{opt.BusinessID, flowid, opt.OID},
		Content: map[string]interface{}{
			"opt":    opt,
			"res":    res,
			"update": update,
			"err":    err,
		},
	}
	if err1 := report.Manager(logData); err1 != nil {
		log.Error("report.Manager(%+v) error(%v)", logData, err1)
	}
}

func (s *Service) sendRscCancleLog(c context.Context, BusinessID int64, oids []string, uid int64, username string, err error) {
	logData := &report.ManagerInfo{
		Uname:    username,
		UID:      uid,
		Business: model.LogBusinessResource,
		Type:     model.LogTypeFromCancle,
		Oid:      0,
		Action:   "cancle",
		Ctime:    time.Now(),
		Index:    []interface{}{BusinessID},
		Content: map[string]interface{}{
			"oids": oids,
			"err":  err,
		},
	}
	if err1 := report.Manager(logData); err1 != nil {
		log.Error("report.Manager(%+v) error(%v)", logData, err1)
	}
}

func (s *Service) sendRscSubmitLog(c context.Context, action string, opt *model.SubmitOptions, res interface{}) {
	logData := &report.ManagerInfo{
		Uname:    opt.Uname,
		UID:      opt.UID,
		Business: model.LogBusinessResource,
		Type:     model.LogTypeFormAuditor,
		Oid:      opt.RID,
		Action:   action,
		Ctime:    time.Now(),
		Index:    []interface{}{opt.BusinessID, opt.FlowID, opt.OID},
		Content: map[string]interface{}{
			"opt": opt,
			"res": res,
		},
	}
	if err := report.Manager(logData); err != nil {
		log.Error("report.Manager(%+v) error(%v)", logData, err)
	}
}

/**
 * 记录流程流转日志
 * oid=rid,action=new_flow_id, index=[net_id, old_flow_id, transition_id, from], content=submit+result
 * from=单个提交/批量提交/跳流程/启动/取消
 */
func (s *Service) sendNetTriggerLog(c context.Context, pm *net.TriggerResult) (err error) {
	var (
		submitValue, resValue []byte
		tran                  string
		content               = map[string]interface{}{}
	)
	if len(pm.TransitionID) > 0 {
		tran = xstr.JoinInts(pm.TransitionID)
	}
	if pm.SubmitToken != nil {
		if submitValue, err = json.Marshal(pm.SubmitToken); err != nil {
			log.Error("sendNetTriggerLog json.Marshal error(%v) submit(%v)", err, pm.SubmitToken)
			return
		}
		content["submit"] = string(submitValue)
	}
	if pm.ResultToken != nil {
		if resValue, err = json.Marshal(pm.ResultToken); err != nil {
			log.Error("sendNetTriggerLog json.Marshal error(%v) result(%v)", err, pm.ResultToken)
			return
		}
		content["result"] = string(resValue)
	}

	data := &report.ManagerInfo{
		Business: model.LogBusinessNet,
		Type:     model.LogTypeNetTrigger,
		Oid:      pm.RID,
		Action:   strconv.FormatInt(pm.NewFlowID, 10),
		Ctime:    time.Now(),
		Index:    []interface{}{pm.NetID, tran, pm.From, pm.OldFlowID},
		Content:  content,
	}
	log.Info("sendNetTriggerLog start send log(%+v)", data)
	report.Manager(data)
	return
}

/**
 * oid: 各元素id, type=level, action=禁用/创建/更新/启用, ctime=time.now, index=[net_id, ch_name, flow_id, tran_id], content=diff
 * level in (net/token/token_bind_flow/token_bind_transition/flow/transition/direction)
 * diff如下：
 * token: obj=name+compare+value(type)
 * token_bind: obj=flow_chname/tran_chname:从token_obj变成token_obj, ch_name从xx变成xx
 * flow: ch_name从xx变成xx, name从xx变成xx
 * tran: ch_name从xx变成xx, name从xx变成xx,trigger从xx变成xx,limit从xx变成xx
 * dir:direction从xx变成xx,order从xx变成xx,guard从xx变成xx,output从xx变成yy
 *
 */
func (s *Service) sendNetConfLog(c context.Context, tp int, oper *model.NetConfOper) (err error) {
	data := &report.ManagerInfo{
		UID:      oper.UID,
		Uname:    "",
		Business: model.LogBusinessNetConf,
		Type:     tp,
		Oid:      oper.OID,
		Action:   oper.Action,
		Ctime:    time.Now(),
		Index:    []interface{}{oper.NetID, oper.FlowID, oper.TranID, oper.ChName},
		Content: map[string]interface{}{
			"diff": strings.Join(oper.Diff, "\r\n"),
		},
	}

	log.Info("sendNetConfLog data(%+v)", data)
	report.Manager(data)
	return
}

//SearchAuditLogCSV 操作日志结果csv
func (s *Service) SearchAuditLogCSV(c context.Context, pm *model.SearchAuditLogParam) (csv [][]string, err error) {
	var (
		res []*model.SearchAuditLog
	)
	if res, _, err = s.SearchAuditLog(c, pm); err != nil {
		return
	}

	csv = make([][]string, len(res)+1)
	csv[0] = []string{"rid", "oid", "task id", "状态", "操作时间", "操作人", "其他信息"}
	for i, item := range res {
		csv[i+1] = []string{
			strconv.FormatInt(item.RID, 10),
			item.OID,
			strconv.FormatInt(item.TaskID, 10),
			item.State,
			item.Stime,
			fmt.Sprintf("%s(%s)", item.Uname, item.Department),
			item.Extra,
		}
	}

	return
}

//SearchAuditLog 查询审核日志
func (s *Service) SearchAuditLog(c context.Context, pm *model.SearchAuditLogParam) (res []*model.SearchAuditLog, p common.Pager, err error) {
	var (
		logs                *model.SearchLogResult
		ridoid, udepartment map[int64]string
		oidrid              map[string]int64
	)

	p = common.Pager{
		Ps: pm.Ps,
		Pn: pm.Pn,
	}
	//oid转换成rid查询
	if len(pm.OID) > 0 {
		if oidrid, err = s.gorm.ResIDByOID(c, pm.BusinessID, pm.OID); err != nil {
			log.Error("SearchAuditLog s.gorm.ResIDByOID error(%+v) pm(%+v)", err, pm)
			return
		}
		if len(oidrid) == 0 {
			return
		}
		ridoid = map[int64]string{}
		for oid, rid := range oidrid {
			pm.RID = append(pm.RID, rid)
			ridoid[rid] = oid
		}
	}

	if logs, err = s.searchAuditLog(c, pm); err != nil {
		err = ecode.AegisSearchErr
		return
	}
	p.Total = logs.Page.Total
	if len(logs.Result) == 0 {
		return
	}

	uids := []int64{}
	rids := []int64{}
	unameuids := []int64{}
	uidunameexist := map[int64]string{}
	res = make([]*model.SearchAuditLog, len(logs.Result))
	for i, item := range logs.Result {
		if item.UID > 0 {
			uids = append(uids, item.UID)
		}
		oid, exist := ridoid[item.OID]
		if !exist {
			rids = append(rids, item.OID)
		}
		if item.Uname != "" {
			uidunameexist[item.UID] = item.Uname
		} else {
			item.Uname = uidunameexist[item.UID]
		}

		if item.Uname == "" && item.UID > 0 {
			unameuids = append(unameuids, item.UID)
		}

		change := &model.Change{}
		if err = json.Unmarshal([]byte(item.Extra), &change); err != nil {
			log.Error("searchAuditLog json.Unmarshal error(%v) extra(%s) pm(%+v)", err, item.Extra, pm)
			return
		}

		flowaction, submitopt := change.GetSubmitOper()

		res[i] = &model.SearchAuditLog{
			RID:        item.OID,
			OID:        oid,
			TaskID:     item.Int2,
			State:      item.Str0,
			Stime:      item.Ctime,
			UID:        item.UID,
			Uname:      item.Uname,
			Department: "",
			Extra:      fmt.Sprintf("操作详情：[%s]%s %s", item.Action, flowaction, submitopt),
		}
	}

	//由搜索结果提供了rid
	if len(rids) > 0 {
		if ridoid, err = s.gorm.ResOIDByID(c, rids); err != nil {
			return
		}
	}

	unames, _ := s.http.GetUnames(c, unameuids)
	udepartment, _ = s.http.GetUdepartment(c, uids)
	for _, item := range res {
		if item.OID == "" {
			item.OID = ridoid[item.RID]
		}
		if item.Uname == "" && item.UID > 0 {
			item.Uname = unames[item.UID]
		}
		item.Department = udepartment[item.UID]
	}
	return
}

func (s *Service) trackAuditLog(c context.Context, pm *model.SearchAuditLogParam) (res []*model.TrackAudit, err error) {
	var (
		logs   *model.SearchLogResult
		flowch map[int64]string
	)
	res = []*model.TrackAudit{}
	if logs, err = s.searchAuditLog(c, pm); err != nil {
		log.Error("trackAuditLog s.searchAuditLog error(%v) pm(%+v)", err, pm)
		return
	}

	res = make([]*model.TrackAudit, len(logs.Result))
	flows := []int64{}
	for i, item := range logs.Result {
		change := &model.Change{}
		if err = json.Unmarshal([]byte(item.Extra), change); err != nil {
			log.Error("trackAuditLog json.Unmarshal error(%v) pm(%+v)", err, pm)
			return
		}

		res[i] = &model.TrackAudit{
			Ctime:  item.Ctime,
			FlowID: []int64{},
			State:  "",
			Uname:  item.Uname,
		}
		if int(item.Type) == model.LogTypeAuditCancel {
			res[i].State = "删除"
		}
		if change.Flow == nil {
			continue
		}
		if int(item.Type) == model.LogTypeAuditCancel {
			var one []int64
			one, err = xstr.SplitInts(change.Flow.OldFlowID.String())
			if err != nil {
				log.Error("trackAuditLog xstr.SplitInts(%s) error(%v)", change.Flow.OldFlowID.String(), err)
				err = nil
				continue
			}
			if len(one) == 0 {
				continue
			}

			res[i].FlowID = one
			flows = append(flows, one...)
			continue
		}

		flows = append(flows, change.Flow.NewFlowID)
		res[i].FlowID = []int64{change.Flow.NewFlowID}
		if change.Flow.ResultToken != nil {
			res[i].State = change.Flow.ResultToken.ChName
		}
	}

	//get flows names
	if len(flows) == 0 {
		return
	}
	if flowch, err = s.gorm.ColumnMapString(c, net.TableFlow, "ch_name", flows, ""); err != nil {
		log.Error("trackAuditLog s.gorm.ColumnMapString error(%v) pm(%+v)", err, pm)
		return
	}
	for _, item := range res {
		fnames := make([]string, len(item.FlowID))
		for i, fid := range item.FlowID {
			fnames[i] = flowch[fid]
		}
		item.FlowName = strings.Join(fnames, ",")
	}
	return
}

//TrackResource 资源信息追踪, 获取资源add/update日志，并分页，以此为基准，获取对应时间端内的资源audit日志；若add/update日志只有不超过1页，则获取全部audit日志；超过1页，最后一页会返回剩余的全部audit日志
func (s *Service) TrackResource(c context.Context, pm *model.TrackParam) (res *model.TrackInfo, p common.Pager, err error) {
	var (
		obj        *resource.Resource
		rsc        []*model.TrackRsc
		audit      []*model.TrackAudit
		rela       [][]int
		LogMinTime = "2018-11-01 10:00:00"
	)

	if obj, err = s.gorm.ResourceByOID(c, pm.OID, pm.BusinessID); err != nil || obj == nil {
		log.Error("TrackResource s.gorm.ResourceByOID error(%v)/not found, pm(%+v)", err, pm)
		return
	}
	if rsc, p, err = s.searchResourceLog(c, obj.ID, pm.Pn, pm.Ps); err != nil {
		err = ecode.AegisSearchErr
		return
	}

	//超过部分不需要查询audit
	topn := int(math.Ceil(float64(p.Total) / float64(p.Ps)))
	if (topn > 0 && topn < p.Pn) || (topn <= 0 && p.Pn > 1) {
		return
	}

	//没有资源日志，则不查询审核日志--资源日志添加失败，还是需要展示审核日志啊，审核日志分页有规律
	ap := &model.SearchAuditLogParam{
		BusinessID: pm.BusinessID,
		RID:        []int64{obj.ID},
		CtimeFrom:  LogMinTime,
		CtimeTo:    "",
		Ps:         1000, //一次性拿出来所有的日志
	}
	//对于audit日志，当add日志各种情况下会返回如下:no data(p.Total <= 0)---全量, 1页(p.Total <= p.Ps)---全量, 2或多页(p1=最新->p1.lasttime, p2=p1.lasttime-p2.lasttime,...pn=pn-1.lasttime-mintime)
	if p.Total > p.Ps { //有多页
		llen := len(rsc)
		if llen > 0 && topn > p.Pn {
			ap.CtimeFrom = rsc[llen-1].Ctime
		}
		if p.Pn > 1 {
			ap.CtimeTo = pm.LastPageTime
		}
	}

	if audit, err = s.trackAuditLog(c, ap); err != nil {
		err = ecode.AegisSearchErr
		return
	}

	//根据ctime聚合，以资源日志为基准
	llen := len(rsc) + 2
	rscctime := make([]string, llen)
	rscctime[0] = time.Now().Format("2006-01-02 15:04:05") //max
	for i, item := range rsc {
		rscctime[i+1] = item.Ctime
	}
	rscctime[llen-1] = time.Time{}.Format("2006-01-02 15:04:05") //min
	index := 0
	for i := 1; i < llen; i++ {
		rel := []int{}
		for ; index < len(audit); index++ {
			t := audit[index].Ctime
			if t >= rscctime[i] && t < rscctime[i-1] {
				rel = append(rel, index)
				continue
			}
			break
		}

		if i == llen-1 && len(rel) == 0 {
			continue
		}
		rela = append(rela, rel)
	}

	res = &model.TrackInfo{
		Add:      rsc,
		Audit:    audit,
		Relation: rela,
	}
	return
}

func (s *Service) searchAuditLog(c context.Context, pm *model.SearchAuditLogParam) (resp *model.SearchLogResult, err error) {
	args := &model.ParamsQueryLog{
		Business:  model.LogBusinessAudit,
		Oid:       pm.RID,
		CtimeFrom: pm.CtimeFrom,
		CtimeTo:   pm.CtimeTo,
		Int2:      pm.TaskID,
		Uname:     pm.Username,
	}
	if pm.State != "" {
		args.Str0 = []string{pm.State}
	}
	if pm.BusinessID > 0 {
		args.Int0 = []int64{pm.BusinessID}
	}

	escm := model.EsCommon{
		Ps:    pm.Ps,
		Pn:    pm.Pn,
		Order: "ctime",
		Sort:  "desc",
	}

	return s.http.QueryLogSearch(c, args, escm)
}

func (s *Service) auditLogByRID(c context.Context, rid int64) (ls []string, err error) {
	resp, err := s.searchAuditLog(c, &model.SearchAuditLogParam{
		RID: []int64{rid},
		Ps:  1000,
		Pn:  1})
	if err != nil || resp == nil {
		return
	}

	for _, result := range resp.Result {
		change := &model.Change{}
		if err = json.Unmarshal([]byte(result.Extra), &change); err != nil {
			log.Error("json.Unmarshal error(%v)", err)
			return
		}

		flowaction, submitopt := change.GetSubmitOper()
		// 时间 + 操作人 + 操作/state + 操作内容
		l := fmt.Sprintf("%s %s[%s] %s %s", result.Ctime, result.Uname, result.Action, flowaction, submitopt)
		ls = append(ls, l)
	}
	return
}

func (s *Service) searchWeightLog(c context.Context, taskid int64, pn, ps int) (ls []*model.WeightLog, count int, err error) {
	args := &model.ParamsQueryLog{
		Business: model.LogBusinessTask,
		Type:     model.LogTYpeTaskWeight,
		Oid:      []int64{taskid},
		Action:   []string{"weight"},
	}
	escm := model.EsCommon{
		Pn:    pn,
		Ps:    ps,
		Order: "ctime",
		Sort:  "desc",
	}

	resp, err := s.http.QueryLogSearch(c, args, escm)
	if err != nil || resp == nil {
		return
	}

	count = resp.Page.Total
	for _, result := range resp.Result {
		logitem := make(map[string]*model.WeightLog)
		if err = json.Unmarshal([]byte(result.Extra), &logitem); err != nil {
			log.Error("json.Unmarshal error(%v)", err)
			return
		}
		ls = append(ls, logitem["weightlog"])
	}
	return
}

func (s *Service) searchConsumerLog(c context.Context, bizid, flowid int64, action []string, uids []int64, ps int) (at map[int64]string, err error) {
	args := &model.ParamsQueryLog{
		Business:  model.LogBusinessTask,
		Type:      model.LogTypeTaskConsumer,
		Action:    action,
		UID:       uids,
		CtimeFrom: time.Now().Add(-24 * time.Hour * 7).Format("2006-01-02 15:04:05"),
	}
	if bizid > 0 {
		args.Int0 = []int64{bizid}
	}
	if flowid > 0 {
		args.Int1 = []int64{flowid}
	}

	escm := model.EsCommon{
		Order: "ctime",
		Sort:  "desc",
		Pn:    1,
		Ps:    ps,
		Group: "uid",
	}
	resp, err := s.http.QueryLogSearch(c, args, escm)
	if err != nil || resp == nil {
		return
	}

	at = make(map[int64]string)
	for _, item := range resp.Result {
		if ct, ok := at[item.UID]; ok {
			if item.Ctime > ct {
				at[item.UID] = item.Ctime
			}
		} else {
			at[item.UID] = item.Ctime
		}
	}
	return
}

func (s *Service) searchResourceLog(c context.Context, rid int64, pn, ps int) (result []*model.TrackRsc, p common.Pager, err error) {
	//根据ctime降序排列
	args := &model.ParamsQueryLog{
		Business: model.LogBusinessResource,
		Type:     model.LogTypeFromAdd,
		Oid:      []int64{rid},
	}
	p = common.Pager{
		Pn: pn,
		Ps: ps,
	}
	escm := model.EsCommon{
		Order: "ctime",
		Sort:  "desc",
		Pn:    pn,
		Ps:    ps,
	}
	resp, err := s.http.QueryLogSearch(c, args, escm)
	if err != nil || resp == nil {
		return
	}

	p.Total = resp.Page.Total
	result = make([]*model.TrackRsc, len(resp.Result))
	for i, item := range resp.Result {
		extra := struct {
			Opt map[string]interface{} `json:"opt"`
		}{}

		if err = json.Unmarshal([]byte(item.Extra), &extra); err != nil {
			log.Error("ResourceLog json.Unmarshal error(%v) extra(%s)", err, item.Extra)
			return
		}

		result[i] = &model.TrackRsc{
			Ctime:   item.Ctime,
			Content: extra.Opt["content"].(string),
			Detail:  extra.Opt,
		}
	}

	//content变化，由于result是根据创建时间降序排列的，以result的最后一个为基础， 向result[0]判断
	content := ""
	for i := len(result) - 1; i >= 0; i-- {
		item := result[i]
		//固定content字段比较变化
		if item.Content != content {
			content = item.Content
			continue
		}
		item.Content = ""
	}
	return
}
