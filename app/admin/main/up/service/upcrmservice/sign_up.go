package upcrmservice

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-common/app/admin/main/up/dao/global"
	"go-common/app/admin/main/up/model"
	"go-common/app/admin/main/up/model/signmodel"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/app/admin/main/up/service/cache"
	"go-common/app/admin/main/up/util"
	accgrpc "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

// SignUpAuditLogs .
func (s *Service) SignUpAuditLogs(c context.Context, arg *signmodel.SignOpSearchArg) (res *signmodel.SignAuditListReply, err error) {
	return s.mng.SignUpAuditLogs(c, arg)
}

// SignAdd add sign info
func (s *Service) SignAdd(context context.Context, arg *signmodel.SignUpArg) (result signmodel.CommonResponse, err error) {
	if arg == nil {
		log.Error("sign add arg is nil")
		return
	}
	// 处理合同信息
	var contractInfo []*signmodel.SignContractInfoArg
	for _, v := range arg.ContractInfo {
		if v == nil || strings.Trim(v.Filename, " ") == "" {
			continue
		}
		if strings.Trim(v.Filelink, " \n\r") == "" {
			err = model.ErrNoFileLink
			log.Error("no file link for contract, please upload file, arg=%v", arg)
			return
		}
		contractInfo = append(contractInfo, v)
	}
	arg.ContractInfo = contractInfo
	// 从context 里拿后台登录信息
	var bmContext, ok = context.(*blademaster.Context)
	if ok {
		uid, ok := util.GetContextValueInt64(bmContext, "uid")
		if ok {
			arg.AdminID = int(uid)
		}
		name, ok := util.GetContextValueString(bmContext, "username")
		if ok {
			arg.AdminName = name
		}
	}
	log.Info("add sign up, req=%+v, admin id=%d, admin name=%s", arg, arg.AdminID, arg.AdminName)
	// 事物
	tx := s.crmdb.BeginTran(context)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// 获取up主基础信息中的分区id
	upInfo, _ := s.crmdb.QueryUpBaseInfo(arg.Mid, "active_tid")
	// 1.先去sign up里插入一条记录，然后获得对应的id
	var dbSignUp signmodel.SignUp
	// 把请求信息制作签约up的db内容
	arg.CopyTo(&dbSignUp)
	if upInfo.ActiveTid != 0 {
		// 把up主的分区id赋值到签约up主信息
		dbSignUp.ActiveTid = int16(upInfo.ActiveTid)
	}
	_, e := s.crmdb.InsertSignUp(tx, &dbSignUp)
	err = e
	if err != nil {
		log.Error("fail to add into sign up db, req=%+v, err=%+v", arg, err)
		return
	}

	log.Info("add sign up ok, new id=%d, next add other info", dbSignUp.ID)
	// 2.把id写入pay/task/contract的sign_id字段中，然后分别将这三种信息插入到数据库
	for _, v := range arg.PayInfo {
		v.SignID = dbSignUp.ID
		v.Mid = dbSignUp.Mid
		if _, err = s.addPayInfo(tx, v); err != nil {
			log.Error("insert payinfo db fail, err=%+v", err)
			break
		}
	}

	for _, v := range arg.TaskInfo {
		v.SignID = dbSignUp.ID
		v.Mid = dbSignUp.Mid
		if _, err = s.addTaskInfo(tx, v); err != nil {
			log.Error("insert payinfo db fail, err=%+v", err)
			break
		}
	}

	for _, v := range arg.ContractInfo {
		v.SignID = dbSignUp.ID
		v.Mid = dbSignUp.Mid
		if _, err = s.addContractInfo(tx, v); err != nil {
			log.Error("insert payinfo db fail, err=%+v", err)
			break
		}
	}
	log.Info("add sign up, new id=%d, all info finish", dbSignUp.ID)
	index := []interface{}{}
	content := map[string]interface{}{
		"new":         arg,
		"old":         nil,
		"change_type": new([]int8),
	}
	// 上报添加的签约日志
	s.AddAuditLog(signmodel.SignUpLogBizID, signmodel.SignUpMidAdd, "新增", int64(arg.AdminID), arg.AdminName, []int64{int64(arg.Mid)}, index, content)
	return
}

// SignUpdate .
func (s *Service) SignUpdate(context context.Context, arg *signmodel.SignUpArg) (result signmodel.CommonResponse, err error) {
	if arg == nil || arg.ID == 0 {
		log.Error("sign up arg is nil")
		return
	}
	// 处理合同信息
	var contractInfo []*signmodel.SignContractInfoArg
	for _, v := range arg.ContractInfo {
		if v == nil || strings.Trim(v.Filename, " ") == "" {
			continue
		}
		if strings.Trim(v.Filelink, " \n\r") == "" {
			err = model.ErrNoFileLink
			log.Error("no file link for contract, please upload file, arg=%v", arg)
			return
		}
		contractInfo = append(contractInfo, v)
	}
	arg.ContractInfo = contractInfo
	// 从context 里拿后台登录信息
	var bmContext, ok = context.(*blademaster.Context)
	if ok {
		uid, ok := util.GetContextValueInt64(bmContext, "uid")
		if ok {
			arg.AdminID = int(uid)
		}
		name, ok := util.GetContextValueString(bmContext, "username")
		if ok {
			arg.AdminName = name
		}
	}
	log.Info("add sign up, req=%+v, admin id=%d, admin name=%s", arg, arg.AdminID, arg.AdminName)
	// 预处理声明参数
	var (
		oriSignUp                                         *signmodel.SignUp
		oriSignPayMap                                     map[int64]*signmodel.SignPay
		oriSignTaskMap                                    map[int64]*signmodel.SignTask
		oriSignContractMap                                map[int64]*signmodel.SignContract
		signPayIDMap, signTaskIDMap, signContractIDMap    map[int64]struct{}
		delSignPayIDs, delSignTaskIDs, delSignContractIDs []int64
		changeType                                        []int8
		signPays                                          []*signmodel.SignPay
		signTasks                                         []*signmodel.SignTask
		signContracts                                     []*signmodel.SignContract
		fields                                            = make(map[int8]struct{})
		signUp                                            = &signmodel.SignUp{}
		oriSignUpArg                                      = new(signmodel.SignUpArg)
	)
	// 从db获取签约信息 、付款信息、任务信息、合同信息
	if oriSignUp, oriSignPayMap, oriSignTaskMap, oriSignContractMap, err = s.crmdb.SignUpID(arg.ID); err != nil {
		log.Error("s.crmdb.SignUpID(%+d) error(%+v)", arg.ID, err)
		return
	}
	if oriSignUp.State == 1 || oriSignUp.State == 100 {
		err = fmt.Errorf("up签约已过期或者已被删除")
		return
	}
	// 把请求参数制作成签约up的db
	arg.CopyTo(signUp)
	// 请求的签约up信息和db的签约up信息做diff, fields包含了修改哪些信息
	signUp.Diff(oriSignUp, fields)
	// 把db内的签约信息制作成请求参数类型数据, 放入上报日志需要
	oriSignUpArg.SignUpBaseInfo.CopyFrom(oriSignUp)
	// 需要更新的付款信息的db id
	signPayIDMap = map[int64]struct{}{}
	// 对比请求的付款信息,diff 出变更的数据
	for _, v := range arg.PayInfo {
		var sp = &signmodel.SignPay{}
		v.CopyTo(sp)
		sp.Diff(oriSignPayMap, fields)
		sp.SignID = arg.ID
		sp.Mid = arg.Mid
		v.SignID = arg.ID
		v.Mid = arg.Mid
		signPays = append(signPays, sp)
		signPayIDMap[v.ID] = struct{}{}
	}
	// 把db内付款信息制作成上报日志需要结构
	for _, v := range oriSignPayMap {
		var pi = &signmodel.SignPayInfoArg{}
		pi.CopyFrom(v)
		pi.SignID = arg.ID
		pi.Mid = arg.Mid
		oriSignUpArg.PayInfo = append(oriSignUpArg.PayInfo, pi)
		if _, ok := signPayIDMap[v.ID]; !ok && v.ID != 0 {
			delSignPayIDs = append(delSignPayIDs, int64(v.ID))
		}
	}
	// 付款信息是否存在删除
	if len(delSignPayIDs) > 0 {
		fields[signmodel.ChangeSignPayHistory] = struct{}{}
	}
	// 需要更新的任务的db id
	signTaskIDMap = map[int64]struct{}{}
	// 对比请求的任务信息,diff 出变更的数据
	for _, v := range arg.TaskInfo {
		var st = &signmodel.SignTask{SignID: arg.ID, Mid: arg.Mid}
		v.CopyTo(st)
		st.Diff(oriSignTaskMap, fields)
		st.SignID = arg.ID
		st.Mid = arg.Mid
		v.SignID = arg.ID
		v.Mid = arg.Mid
		signTasks = append(signTasks, st)
		signTaskIDMap[v.ID] = struct{}{}
	}
	// 把db内任务信息制作成上报日志需要结构
	for _, v := range oriSignTaskMap {
		var ti = &signmodel.SignTaskInfoArg{}
		ti.CopyFrom(v)
		ti.SignID = arg.ID
		ti.Mid = arg.Mid
		oriSignUpArg.TaskInfo = append(oriSignUpArg.TaskInfo, ti)
		if _, ok := signTaskIDMap[v.ID]; !ok && v.ID != 0 {
			delSignTaskIDs = append(delSignTaskIDs, int64(v.ID))
		}
	}
	// 任务是否存在删除
	if len(delSignTaskIDs) > 0 {
		fields[signmodel.ChangeSignTaskHistory] = struct{}{}
	}
	// 需要更新的合同的db id
	signContractIDMap = map[int64]struct{}{}
	// 对比请求的合同信息,diff 出变更的数据
	for _, v := range arg.ContractInfo {
		var sc = &signmodel.SignContract{SignID: arg.ID, Mid: arg.Mid}
		v.CopyTo(sc)
		sc.Diff(oriSignContractMap, fields)
		sc.SignID = arg.ID
		sc.Mid = arg.Mid
		v.SignID = arg.ID
		v.Mid = arg.Mid
		signContracts = append(signContracts, sc)
		signContractIDMap[v.ID] = struct{}{}
	}
	// 把db内合同信息制作成上报日志需要结构
	for _, v := range oriSignContractMap {
		var si = &signmodel.SignContractInfoArg{}
		si.CopyFrom(v)
		si.SignID = arg.ID
		si.Mid = arg.Mid
		oriSignUpArg.ContractInfo = append(oriSignUpArg.ContractInfo, si)
		if _, ok := signContractIDMap[v.ID]; !ok && v.ID != 0 {
			delSignContractIDs = append(delSignContractIDs, int64(v.ID))
		}
	}
	// 合同是否存在删除
	if len(delSignContractIDs) > 0 {
		fields[signmodel.ChangeSignContractHistory] = struct{}{}
	}
	for k := range fields {
		changeType = append(changeType, k)
	}
	if len(changeType) == 0 {
		err = fmt.Errorf("up签约信息暂无修改")
		return
	}
	tx := s.crmdb.BeginTran(context)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if _, err = s.crmdb.InsertSignUp(tx, signUp); err != nil {
		log.Error("fail to add into sign up db, req=%+v, err=%+v", arg, err)
		return
	}
	for _, v := range signPays {
		if _, err = s.crmdb.InsertPayInfo(tx, v); err != nil {
			log.Error("insert pay info db fail, err=%+v", err)
			break
		}
	}
	for _, v := range signTasks {
		if _, err = s.crmdb.InsertTaskInfo(tx, v); err != nil {
			log.Error("insert task info db fail, err=%+v", err)
			break
		}
		var (
			init bool
			sth  *signmodel.SignTaskHistory
		)
		if sth, init, err = s.crmdb.GetOrCreateTaskHistory(tx, v); err != nil {
			log.Error("s.crmdb.GetOrCreateTaskHistory, err=%+v", err)
			break
		}
		if !init {
			sth.Attribute = v.Attribute
			sth.TaskCondition = v.TaskCondition
			sth.TaskType = v.TaskType
			if err = s.crmdb.UpSignTaskHistory(tx, sth); err != nil {
				log.Error("s.crmdb.UpSignTaskHistory, err=%+v", err)
				break
			}
		}
	}
	for _, v := range signContracts {
		if _, err = s.crmdb.InsertContractInfo(tx, v); err != nil {
			log.Error("insert contract info db fail, err=%+v", err)
			break
		}
	}
	if _, err = s.crmdb.DelPayInfo(tx, delSignPayIDs); err != nil {
		log.Error("delete task info db fail, err=%+v", err)
		return
	}
	if _, err = s.crmdb.DelTaskInfo(tx, delSignTaskIDs); err != nil {
		log.Error("delete task info db fail, err=%+v", err)
		return
	}
	if _, err = s.crmdb.DelSignContract(tx, delSignContractIDs); err != nil {
		log.Error("delete task info db fail, err=%+v", err)
		return
	}
	index := []interface{}{int64(arg.ID)}
	content := map[string]interface{}{
		"new":         arg,
		"old":         oriSignUpArg,
		"change_type": changeType,
	}
	// 上报变更信息
	s.AddAuditLog(signmodel.SignUpLogBizID, signmodel.SignUpMidUpdate, "修改", int64(arg.AdminID), arg.AdminName, []int64{int64(arg.Mid)}, index, content)
	return
}

// ViolationAdd .
func (s *Service) ViolationAdd(context context.Context, arg *signmodel.ViolationArg) (result signmodel.CommonResponse, err error) {
	if arg == nil || arg.SignID == 0 {
		log.Error("violation add arg is nil")
		return
	}
	su := &signmodel.SignUp{}
	if err = s.crmdb.GetDb().Table(signmodel.TableSignUp).Where("id = ? AND state = 0", arg.SignID).Find(su).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("签约ID(%d)不存在", arg.SignID)
		return
	}
	var bmContext, ok = context.(*blademaster.Context)
	if ok {
		uid, ok := util.GetContextValueInt64(bmContext, "uid")
		if ok {
			arg.AdminID = uid
		}
		name, ok := util.GetContextValueString(bmContext, "username")
		if ok {
			arg.AdminName = name
		}
	}
	log.Info("add sign up, req=%+v, admin id=%d, admin name=%s", arg, arg.AdminID, arg.AdminName)
	tx := s.crmdb.BeginTran(context)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// 新增违约历史
	vh := &signmodel.SignViolationHistory{}
	arg.CopyTo(vh)
	if err = tx.Save(vh).Error; err != nil {
		log.Error("sign violation up fail, err=%+v", err)
		return
	}
	// 新增违约次数
	if err = tx.Table(signmodel.TableSignUp).Where("id = ?", arg.SignID).UpdateColumns(
		map[string]interface{}{
			"violation_times": gorm.Expr("violation_times + ?", 1),
			"admin_id":        arg.AdminID,
			"admin_name":      arg.AdminName,
		}).Error; err != nil {
		log.Error("sign up add violation times fail, err=%+v", err)
	}
	return
}

// ViolationRetract .
func (s *Service) ViolationRetract(context context.Context, arg *signmodel.IDArg) (result signmodel.CommonResponse, err error) {
	var bmContext, ok = context.(*blademaster.Context)
	if ok {
		uid, ok := util.GetContextValueInt64(bmContext, "uid")
		if ok {
			arg.AdminID = uid
		}
		name, ok := util.GetContextValueString(bmContext, "username")
		if ok {
			arg.AdminName = name
		}
	}
	log.Info("add sign up, req=%+v, admin id=%d, admin name=%s", arg, arg.AdminID, arg.AdminName)
	tx := s.crmdb.BeginTran(context)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// 删除违约历史
	if err = tx.Table(signmodel.TableSignViolationHistory).Where("id = ?", arg.ID).UpdateColumns(
		map[string]interface{}{
			"state":      100,
			"admin_id":   arg.AdminID,
			"admin_name": arg.AdminName,
		}).Error; err != nil {
		log.Error("sign violation Retract  fail, err=%+v", err)
	}
	// 减少违约次数
	if err = tx.Table(signmodel.TableSignUp).Where("id = ?", arg.SignID).UpdateColumns(
		map[string]interface{}{
			"violation_times": gorm.Expr("violation_times - ?", 1),
			"admin_id":        arg.AdminID,
			"admin_name":      arg.AdminName,
		}).Error; err != nil {
		log.Error("sign up dec violation times fail, err=%+v", err)
	}
	return
}

// ViolationList .
func (s *Service) ViolationList(context context.Context, arg *signmodel.PageArg) (result *signmodel.ViolationResult, err error) {
	if arg == nil {
		log.Error("arg is nil")
		return
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	if arg.Size <= 0 || arg.Size >= 50 {
		arg.Size = 20
	}
	result = new(signmodel.ViolationResult)
	result.Result = []*signmodel.ViolationArg{}
	result.Page = arg.Page
	result.Size = arg.Size
	var (
		count  int
		offset = (arg.Page - 1) * arg.Size
		vhs    = []*signmodel.SignViolationHistory{}
	)
	if err = s.crmdb.GetDb().Table(signmodel.TableSignViolationHistory).Where("sign_id = ?", arg.SignID).Count(&count).Error; err != nil {
		log.Error("violation count fail, err=%+v", err)
		return
	}
	if count <= 0 {
		return
	}
	result.TotalCount = count
	if err = s.crmdb.GetDb().Table(signmodel.TableSignViolationHistory).Where("sign_id = ?", arg.SignID).Order(fmt.Sprintf("%s %s", "mtime", "DESC")).
		Offset(offset).
		Limit(arg.Size).
		Find(&vhs).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("violationList fail, err=%+v", err)
	}
	for _, v := range vhs {
		var re = &signmodel.ViolationArg{}
		re.CopyFrom(v)
		result.Result = append(result.Result, re)
	}
	return
}

// AbsenceAdd .
func (s *Service) AbsenceAdd(context context.Context, arg *signmodel.AbsenceArg) (result signmodel.CommonResponse, err error) {
	if arg == nil || arg.SignID == 0 {
		log.Error("violation add arg is nil")
		return
	}
	su := &signmodel.SignUp{}
	if err = s.crmdb.GetDb().Table(signmodel.TableSignUp).Where("id = ? AND state = 0", arg.SignID).Find(su).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("签约ID(%d)不存在", arg.SignID)
		return
	}
	var Result struct {
		LeaveTimes int
	}
	if err = s.crmdb.GetDb().Table(signmodel.TableSignTaskAbsence).Select("SUM(absence_count) as leave_times").Where("sign_id = ? AND state = 0", arg.SignID).Scan(&Result).Error; err != nil {
		log.Error("sign task absence sum fail, err=%+v", err)
		return
	}
	var bmContext, ok = context.(*blademaster.Context)
	if ok {
		uid, ok := util.GetContextValueInt64(bmContext, "uid")
		if ok {
			arg.AdminID = uid
		}
		name, ok := util.GetContextValueString(bmContext, "username")
		if ok {
			arg.AdminName = name
		}
	}
	log.Info("add sign up, req=%+v, admin id=%d, admin name=%s", arg, arg.AdminID, arg.AdminName)
	tx := s.crmdb.BeginTran(context)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var sth *signmodel.SignTaskHistory
	// 查询签约周期任务历史数据，没有新增
	if sth, err = s.getOrCreateTaskHistory(tx, arg.SignID); err != nil {
		log.Error("s.getOrCreateTaskHistory(%d), err=%+v", arg.SignID, err)
		return
	}
	// 增加请假历史
	ta := &signmodel.SignTaskAbsence{}
	arg.CopyTo(ta)
	ta.TaskHistoryID = sth.ID
	if err = tx.Save(ta).Error; err != nil {
		log.Error("sign task absence  fail, err=%+v", err)
		return
	}
	// 新增请假历史
	if err = tx.Table(signmodel.TableSignUp).Where("id = ?", arg.SignID).UpdateColumns(
		map[string]interface{}{
			"leave_times": Result.LeaveTimes + arg.AbsenceCount,
			"admin_id":    arg.AdminID,
			"admin_name":  arg.AdminName,
		}).Error; err != nil {
		log.Error("sign up add leave time fail, err=%+v", err)
	}
	return
}

// get task history, if not exist, then will create it
func (s *Service) getOrCreateTaskHistory(tx *gorm.DB, signID int64) (res *signmodel.SignTaskHistory, err error) {
	st := new(signmodel.SignTask)
	if err = s.crmdb.GetDb().Select("*").Where("sign_id = ?", signID).Find(&st).Error; err != nil {
		return
	}
	if res, _, err = s.crmdb.GetOrCreateTaskHistory(tx, st); err != nil {
		log.Error("s.crmdb.GetOrCreateTaskHistory, err=%+v", err)
	}
	return
}

// AbsenceRetract .
func (s *Service) AbsenceRetract(context context.Context, arg *signmodel.IDArg) (result signmodel.CommonResponse, err error) {
	var bmContext, ok = context.(*blademaster.Context)
	if ok {
		uid, ok := util.GetContextValueInt64(bmContext, "uid")
		if ok {
			arg.AdminID = uid
		}
		name, ok := util.GetContextValueString(bmContext, "username")
		if ok {
			arg.AdminName = name
		}
	}
	log.Info("add sign up, req=%+v, admin id=%d, admin name=%s", arg, arg.AdminID, arg.AdminName)
	tx := s.crmdb.BeginTran(context)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	signTaskAbsences := &signmodel.SignTaskAbsence{}
	if err = s.crmdb.GetDb().Table(signmodel.TableSignTaskAbsence).Select("absence_count").Where("id = ? AND state = 0", arg.ID).Find(&signTaskAbsences).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign task absence_count fail, err=%+v", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("签约ID(%d)的请假ID(%d)不存在", arg.SignID, arg.ID)
		return
	}
	// 删除请假历史
	if err = tx.Table(signmodel.TableSignTaskAbsence).Where("id = ?", arg.ID).UpdateColumns(
		map[string]interface{}{
			"state":      100,
			"admin_id":   arg.AdminID,
			"admin_name": arg.AdminName,
		}).Error; err != nil {
		log.Error("task absence Retract  fail, err=%+v", err)
	}
	// 减少请假次数
	if err = tx.Table(signmodel.TableSignUp).Where("id = ?", arg.SignID).UpdateColumns(
		map[string]interface{}{
			"leave_times": gorm.Expr("leave_times - ?", signTaskAbsences.AbsenceCount),
			"admin_id":    arg.AdminID,
			"admin_name":  arg.AdminName,
		}).Error; err != nil {
		log.Error("sign up dec leave times fail, err=%+v", err)
	}
	return
}

// AbsenceList .
func (s *Service) AbsenceList(context context.Context, arg *signmodel.PageArg) (result *signmodel.AbsenceResult, err error) {
	if arg == nil {
		log.Error("arg is nil")
		return
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	if arg.Size <= 0 || arg.Size >= 50 {
		arg.Size = 20
	}
	result = new(signmodel.AbsenceResult)
	result.Result = []*signmodel.AbsenceArg{}
	result.Page = arg.Page
	result.Size = arg.Size
	var (
		count      int
		taskHisIDs []int64
		mst        = make(map[int64]*signmodel.SignTaskHistory)
		offset     = (arg.Page - 1) * arg.Size
		tas        = []*signmodel.SignTaskAbsence{}
		sts        = []*signmodel.SignTaskHistory{}
		su         = &signmodel.SignUp{}
	)
	if err = s.crmdb.GetDb().Table(signmodel.TableSignUp).Where("id = ? AND state IN (0,1)", arg.SignID).Find(su).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("up签约不存在")
		return
	}
	if err = s.crmdb.GetDb().Table(signmodel.TableSignTaskAbsence).Where("sign_id = ?", arg.SignID).Count(&count).Error; err != nil {
		log.Error("sign task absence count fail, err=%+v", err)
		return
	}
	if count <= 0 {
		return
	}
	result.TotalCount = count
	if err = s.crmdb.GetDb().Table(signmodel.TableSignTaskAbsence).Where("sign_id = ?", arg.SignID).Order(fmt.Sprintf("%s %s", "mtime", "DESC")).
		Offset(offset).
		Limit(arg.Size).
		Find(&tas).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("absenceList fail, err=%+v", err)
	}
	for _, v := range tas {
		taskHisIDs = append(taskHisIDs, v.TaskHistoryID)
	}
	if err = s.crmdb.GetDb().Table(signmodel.TableSignTaskHistory).Where("id IN (?) AND sign_id = ?", taskHisIDs, arg.SignID).Find(&sts).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign task history fail, err=%+v", err)
		return
	}
	for _, v := range sts {
		mst[v.ID] = v
	}
	for _, v := range tas {
		var (
			ok           bool
			sDate, eDate time.Time
			st           *signmodel.SignTaskHistory
			re           = &signmodel.AbsenceArg{}
		)
		// 从任务历史里面获取签约周期
		if st, ok = mst[int64(v.TaskHistoryID)]; ok {
			if st.TaskType != signmodel.TaskTypeAccumulate {
				sDate, eDate = signmodel.GetTaskDuration(st.GenerateDate.Time(), st.TaskType)
				re.TaskBegin = xtime.Time(sDate.Unix())
				re.TaskEnd = xtime.Time(eDate.Unix())
			} else {
				re.TaskBegin = su.BeginDate
				re.TaskEnd = su.EndDate
			}
		}
		re.CopyFrom(v)
		result.Result = append(result.Result, re)
	}
	return
}

// ViewCheck .
func (s *Service) ViewCheck(context context.Context, arg *signmodel.PowerCheckArg) (res *signmodel.PowerCheckReply, err error) {
	res = &signmodel.PowerCheckReply{}
	if arg == nil || arg.Mid == 0 {
		log.Error("view arg is nil")
		return
	}
	var count int64
	if err = s.crmdb.GetDb().Table(signmodel.TableSignUp).Where("mid = ?", arg.Mid).Count(&count).Error; err != nil {
		log.Error("db fail, err=%+v", err)
		return
	}
	if count > 0 {
		res.IsSign = true
	}
	var baseInfo upcrmmodel.UpBaseInfo
	if baseInfo, err = s.crmdb.QueryUpBaseInfo(arg.Mid, "active_tid"); err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.crmdb.QueryUpBaseInfo(%d), err=%+v", arg.Mid, err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	for _, tid := range arg.TIDs {
		if tid == int16(baseInfo.ActiveTid) {
			res.IsPower = true
			return
		}
	}
	return
}

/*
	id, 最后插入的id
*/
func (s *Service) addPayInfo(tx *gorm.DB, arg *signmodel.SignPayInfoArg) (id int64, err error) {
	if arg == nil {
		err = errors.New("add pay info nil pointer")
		return
	}

	var dbInfo signmodel.SignPay
	arg.CopyTo(&dbInfo)
	var _, e = s.crmdb.InsertPayInfo(tx, &dbInfo)
	err = e
	id = dbInfo.ID
	return
}

func (s *Service) addTaskInfo(tx *gorm.DB, arg *signmodel.SignTaskInfoArg) (id int64, err error) {
	if arg == nil {
		err = errors.New("add task info nil pointer")
		return
	}

	var dbInfo signmodel.SignTask
	arg.CopyTo(&dbInfo)
	if _, err = s.crmdb.InsertTaskInfo(tx, &dbInfo); err != nil {
		log.Error("s.crmdb.InsertTaskInfo(%+v) error(%+v)", &dbInfo, err)
		return
	}
	id = dbInfo.ID
	if _, _, err = s.crmdb.GetOrCreateTaskHistory(tx, &dbInfo); err != nil {
		log.Error("s.crmdb.GetOrCreateTaskHistory(%+v) error(%+v)", &dbInfo, err)
	}
	return
}

func (s *Service) addContractInfo(tx *gorm.DB, arg *signmodel.SignContractInfoArg) (id int64, err error) {
	if arg == nil {
		err = errors.New("add contract info nil pointer")
		return
	}

	var dbInfo signmodel.SignContract
	arg.CopyTo(&dbInfo)
	var _, e = s.crmdb.InsertContractInfo(tx, &dbInfo)
	err = e
	id = dbInfo.ID
	return
}

// type sortPayFunc func(p1, p2 *signmodel.SignPayInfoArg) bool

// type paySorter struct {
// 	datas []*signmodel.SignPayInfoArg
// 	by    sortPayFunc // Closure used in the Less method.
// }

// Len is part of sort.Interface.
// func (s *paySorter) Len() int {
// 	return len(s.datas)
// }

// // Swap is part of sort.Interface.
// func (s *paySorter) Swap(i, j int) {
// 	s.datas[i], s.datas[j] = s.datas[j], s.datas[i]
// }

// // Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
// func (s *paySorter) Less(i, j int) bool {
// 	return s.by(s.datas[i], s.datas[j])
// }

// func sortPayInfo(planets []*signmodel.SignPayInfoArg, sortfunc sortPayFunc) {
// 	ps := &paySorter{
// 		datas: planets,
// 		by:    sortfunc, // The Sort method's receiver is the function (closure) that defines the sort order.
// 	}
// 	sort.Sort(ps)
// }

// func sortByDueAsc(p1, p2 *signmodel.SignPayInfoArg) bool {
// 	var v1, _ = time.Parse(upcrmmodel.TimeFmtDate, p1.DueDate)
// 	var v2, _ = time.Parse(upcrmmodel.TimeFmtDate, p2.DueDate)
// 	return v1.Before(v2)
// }

//SignQuery sign query
func (s *Service) SignQuery(c context.Context, arg *signmodel.SignQueryArg) (res *signmodel.SignQueryResult, err error) {
	if arg == nil {
		log.Error("arg is nil")
		return
	}
	if arg.Page <= 0 {
		arg.Page = 1
	}
	if arg.IsDetail == signmodel.SignUpDetail {
		arg.Size = 5
	} else {
		if arg.Size <= 0 || arg.Size > 20 {
			arg.Size = 20
		}
	}
	var (
		count                 int
		mids                  []int64
		tids                  []int64
		signIDs               []int64
		signTaskIDs           []int64
		signTaskHistoryIDs    []int64
		tpNames               map[int64]string
		db                    = s.crmdb.GetDb()
		signUpHandle          = db.Table(signmodel.TableSignUp)
		signPayHandle         = db.Table(signmodel.TableSignPay)
		signContractHandle    = db.Table(signmodel.TableSignContract)
		signTaskHistoryHandle = db.Table(signmodel.TableSignTaskHistory)
		signAbsenceHandle     = db.Table(signmodel.TableSignTaskAbsence)
		signTaskHandle        = db.Table(signmodel.TableSignTask)
		signUps               = []*signmodel.SignUp{}
		signPayInfos          = []*signmodel.SignPay{}
		signContractInfos     = []*signmodel.SignContract{}
		signTaskHistorys      = []*signmodel.SignTaskHistory{}
		signTaskAbsences      = []*signmodel.SignTaskAbsence{}
		signTasks             = []*signmodel.SignTask{}
		signUpBaseInfos       = []*signmodel.SignUpBaseInfo{}
		signPayInfoMap        = make(map[int64][]*signmodel.SignPayInfoArg)
		signContractInfoMap   = make(map[int64][]*signmodel.SignContractInfoArg)
		signTaskHistoryMap    = make(map[int64][]*signmodel.SignTaskHistoryArg)
		offset                = (arg.Page - 1) * arg.Size
	)
	res = new(signmodel.SignQueryResult)
	res.Page = arg.Page
	res.Size = arg.Size
	res.Result = []*signmodel.SignUpsArg{}
	if len(arg.Tids) != 0 {
		signUpHandle = signUpHandle.Where("active_tid IN (?)", arg.Tids)
	}
	if arg.Mid != 0 {
		signUpHandle = signUpHandle.Where("mid = ?", arg.Mid)
	}
	if arg.DueSign != 0 {
		signUpHandle = signUpHandle.Where("due_warn = ?", 2)
	}
	if arg.DuePay != 0 {
		signUpHandle = signUpHandle.Where("pay_expire_state = ?", 2)
	}
	if arg.ExpireSign != 0 {
		signUpHandle = signUpHandle.Where("state = ?", 1)
	}
	if arg.Sex != -1 {
		signUpHandle = signUpHandle.Where("sex = ?", arg.Sex)
	}
	if len(arg.Country) != 0 {
		signUpHandle = signUpHandle.Where("country IN (?)", arg.Country)
	}
	if arg.ActiveTID != 0 {
		signUpHandle = signUpHandle.Where("active_tid = ?", arg.ActiveTID)
	}
	if arg.SignType != 0 {
		signUpHandle = signUpHandle.Where("sign_type = ?", arg.SignType)
	}
	if arg.TaskState != 0 {
		signUpHandle = signUpHandle.Where("task_state = ?", arg.TaskState)
	}
	if arg.SignBegin != 0 {
		signUpHandle = signUpHandle.Where("begin_date >= ?", arg.SignBegin)
	}
	if arg.SignEnd != 0 {
		signUpHandle = signUpHandle.Where("end_date <= ?", arg.SignEnd)
	}
	signUpHandle = signUpHandle.Where("state IN (0,1)")
	if err = signUpHandle.Count(&count).Error; err != nil {
		log.Error("signUps count fail, err=%+v", err)
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = count
	if arg.IsDetail == signmodel.SignUpDetail {
		signUpHandle = signUpHandle.Order(fmt.Sprintf("%s %s", "id", "DESC"))
	} else {
		signUpHandle = signUpHandle.Order(fmt.Sprintf("%s %s", "mtime", "DESC"))
	}
	if err = signUpHandle.Offset(offset).Limit(arg.Size).Find(&signUps).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("signUps fail, err=%+v", err)
		return
	}
	for _, v := range signUps {
		var signUpBaseInfo = &signmodel.SignUpBaseInfo{}
		signUpBaseInfo.CopyFrom(v)
		signUpBaseInfos = append(signUpBaseInfos, signUpBaseInfo)
		signIDs = append(signIDs, v.ID)
		mids = append(mids, v.Mid)
		tids = append(tids, int64(v.ActiveTid))
	}
	tpNames = cache.GetTidName(tids...)
	var infosReply *accgrpc.InfosReply
	if infosReply, err = global.GetAccClient().Infos3(c, &accgrpc.MidsReq{Mids: mids, RealIp: metadata.String(c, metadata.RemoteIP)}); err != nil {
		log.Error("global.GetAccClient().Infos3(%+v)  error(%+v)", mids, err)
		err = nil
	}
	if err = signPayHandle.Where("sign_id IN (?) AND state IN (0,1)", signIDs).Order(fmt.Sprintf("%s %s", "due_date", "ASC")).Find(&signPayInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign pays fail, err=%+v", err)
		return
	}
	for _, v := range signPayInfos {
		var signPayInfo = &signmodel.SignPayInfoArg{}
		signPayInfo.CopyFrom(v)
		signPayInfoMap[v.SignID] = append(signPayInfoMap[v.SignID], signPayInfo)
	}
	if err = signContractHandle.Where("sign_id IN (?) AND state = 0", signIDs).Order(fmt.Sprintf("%s %s", "id", "ASC")).Find(&signContractInfos).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign contract fail, err=%+v", err)
		return
	}
	for _, v := range signContractInfos {
		var signContractInfo = &signmodel.SignContractInfoArg{}
		signContractInfo.CopyFrom(v)
		signContractInfoMap[v.SignID] = append(signContractInfoMap[v.SignID], signContractInfo)
	}
	// task
	if err = signTaskHandle.Where("sign_id IN (?) AND state = 0", signIDs).Find(&signTasks).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign task fail, err=%+v", err)
		return
	}
	for _, v := range signTasks {
		signTaskIDs = append(signTaskIDs, v.ID)
	}
	// task history  NOTE: mtime 解决generate_date为累积时没用办法排序 有一定风险会错乱
	if arg.IsDetail == signmodel.SignUpList {
		var signTaskHistorySQL = `SELECT * FROM sign_task_history WHERE task_template_id IN (?) AND state IN (1,2) AND mtime = (SELECT MAX(mtime) FROM sign_task_history s 
		WHERE s.sign_id=sign_task_history.sign_id AND s.task_template_id IN (?) AND s.state IN (1,2))`
		if err = signTaskHistoryHandle.Raw(signTaskHistorySQL, signTaskIDs, signTaskIDs).Find(&signTaskHistorys).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("sign task history fail, err=%+v", err)
			return
		}
	} else {
		if err = signTaskHistoryHandle.Where("task_template_id IN (?) AND state IN (1,2)", signTaskIDs).Order(fmt.Sprintf("%s %s", "id", "DESC")).Find(&signTaskHistorys).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("sign task history fail, err=%+v", err)
			return
		}
	}
	for _, v := range signTaskHistorys {
		signTaskHistoryIDs = append(signTaskHistoryIDs, v.ID)
	}
	if err = signAbsenceHandle.Raw("select sum(absence_count)as absence_count, task_history_id from sign_task_absence where task_history_id IN (?) AND state = 0 group by task_history_id", signTaskHistoryIDs).Find(&signTaskAbsences).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign absence fail, err=%+v", err)
		return
	}
	var hidMap = make(map[int64]int)
	for _, v := range signTaskAbsences {
		hidMap[v.TaskHistoryID] = v.AbsenceCount
	}
	for _, v := range signTaskHistorys {
		var (
			absenceCounter     int
			signTaskHistoryArg = &signmodel.SignTaskHistoryArg{}
		)
		if count, ok := hidMap[v.ID]; ok {
			absenceCounter = count
		}
		signTaskHistoryArg.CopyFrom(v, absenceCounter)
		signTaskHistoryMap[v.SignID] = append(signTaskHistoryMap[v.SignID], signTaskHistoryArg)
	}
	for _, v := range signUpBaseInfos {
		if tpName, ok := tpNames[int64(v.ActiveTid)]; ok {
			v.TypeName = tpName
		}
		if infosReply != nil && infosReply.Infos != nil {
			if info, ok := infosReply.Infos[v.Mid]; ok {
				v.Name = info.Name
			}
		}
		if arg.IsDetail == signmodel.SignUpDetail {
			res.SignBaseInfo = v
		}
		var signUpsArg = &signmodel.SignUpsArg{}
		signUpsArg.SignUpBaseInfo = *v
		if signPayInfoArg, ok := signPayInfoMap[v.ID]; ok {
			signUpsArg.PayInfo = signPayInfoArg
		}
		if signContractInfoArg, ok := signContractInfoMap[v.ID]; ok {
			signUpsArg.ContractInfo = signContractInfoArg
		}
		if signTaskHistoryArg, ok := signTaskHistoryMap[v.ID]; ok {
			for _, sth := range signTaskHistoryArg {
				if sth.TaskType == signmodel.TaskTypeAccumulate {
					sth.TaskBegin = v.BeginDate
					sth.TaskEnd = v.EndDate
				}
			}
			signUpsArg.TaskHistoryInfo = signTaskHistoryArg
		}
		res.Result = append(res.Result, signUpsArg)
	}
	return
}

// SignQueryID .
func (s *Service) SignQueryID(c context.Context, arg *signmodel.SignIDArg) (res *signmodel.SignUpArg, err error) {
	var (
		tpNames         map[int64]string
		signUp          *signmodel.SignUp
		signPayMap      map[int64]*signmodel.SignPay
		signTaskMap     map[int64]*signmodel.SignTask
		signContractMap map[int64]*signmodel.SignContract
	)
	res = new(signmodel.SignUpArg)
	if signUp, signPayMap, signTaskMap, signContractMap, err = s.crmdb.SignUpID(arg.ID); err != nil {
		log.Error("s.crmdb.SignUpID(%+d) error(%+v)", arg.ID, err)
		return
	}
	if signUp == nil {
		return
	}
	var infoReply *accgrpc.InfoReply
	if infoReply, err = global.GetAccClient().Info3(c, &accgrpc.MidReq{Mid: signUp.Mid, RealIp: metadata.String(c, metadata.RemoteIP)}); err != nil {
		log.Error("global.GetAccClient().Infos3(%d)  error(%+v)", signUp.Mid, err)
		err = nil
	}
	tpNames = cache.GetTidName(int64(signUp.ActiveTid))
	res.SignUpBaseInfo.CopyFrom(signUp)
	if infoReply != nil && infoReply.Info != nil {
		res.SignUpBaseInfo.Name = infoReply.Info.Name
	}
	if tpName, ok := tpNames[int64(signUp.ActiveTid)]; ok {
		res.SignUpBaseInfo.TypeName = tpName
	}
	for _, v := range signPayMap {
		var payInfo = &signmodel.SignPayInfoArg{}
		payInfo.CopyFrom(v)
		res.PayInfo = append(res.PayInfo, payInfo)
	}
	if signPayMap != nil {
		sort.Slice(res.PayInfo, func(i int, j int) bool {
			return res.PayInfo[i].DueDate < res.PayInfo[j].DueDate
		})
	}
	for _, v := range signTaskMap {
		var taskInfo = &signmodel.SignTaskInfoArg{}
		taskInfo.CopyFrom(v)
		res.TaskInfo = append(res.TaskInfo, taskInfo)
	}
	if signTaskMap != nil {
		sort.Slice(res.TaskInfo, func(i int, j int) bool {
			return res.TaskInfo[i].ID < res.TaskInfo[j].ID
		})
	}
	for _, v := range signContractMap {
		var contractInfo = &signmodel.SignContractInfoArg{}
		contractInfo.CopyFrom(v)
		res.ContractInfo = append(res.ContractInfo, contractInfo)
	}
	if signContractMap != nil {
		sort.Slice(res.ContractInfo, func(i int, j int) bool {
			return res.ContractInfo[i].ID < res.ContractInfo[j].ID
		})
	}
	return
}

// SignPayComplete complete sign pay
func (s *Service) SignPayComplete(con context.Context, arg *signmodel.SignPayCompleteArg) (result signmodel.SignPayCompleteResult, err error) {
	var affectedrow, e = s.crmdb.PayComplete(arg.IDs)
	if e != nil {
		err = e
		log.Error("fail to complete pay task, err=%+v", e)
		return
	}
	log.Info("complete pay, id=%+v, affected row=%d", arg.IDs, affectedrow)
	return
}

// SignCheckExist check sign up has an valid contract
func (s *Service) SignCheckExist(c context.Context, arg *signmodel.SignCheckExsitArg) (result signmodel.SignCheckExsitResult, err error) {
	result.Exist, err = s.crmdb.CheckUpHasValidContract(arg.Mid, time.Now())
	if err != nil {
		log.Error("check up has valid contract fail, err=%+v", err)
	}
	return
}

// Countrys .
func (s *Service) Countrys(c context.Context, arg *signmodel.CommonArg) (res *signmodel.SignCountrysReply, err error) {
	var (
		signUp       []*signmodel.SignUp
		db           = s.crmdb.GetDb()
		signUpHandle = db.Table(signmodel.TableSignUp)
	)
	res = new(signmodel.SignCountrysReply)
	if err = signUpHandle.Raw("select DISTINCT(country) from sign_up").Find(&signUp).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign up fail, err=%+v", err)
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	for _, v := range signUp {
		if v.Country == "" {
			continue
		}
		res.List = append(res.List, v.Country)
	}
	return
}

// Tids .
func (s *Service) Tids(c context.Context, arg *signmodel.CommonArg) (res *signmodel.SignTidsReply, err error) {
	var (
		tids         []int64
		signUp       []*signmodel.SignUp
		db           = s.crmdb.GetDb()
		signUpHandle = db.Table(signmodel.TableSignUp)
	)
	res = new(signmodel.SignTidsReply)
	if err = signUpHandle.Raw("select DISTINCT(active_tid) from sign_up").Find(&signUp).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("sign up fail, err=%+v", err)
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	for _, v := range signUp {
		if v.ActiveTid == 0 {
			continue
		}
		tids = append(tids, int64(v.ActiveTid))
	}
	res.List = cache.GetTidName(tids...)
	return
}
