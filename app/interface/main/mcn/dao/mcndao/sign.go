package mcndao

import (
	"time"

	adminmodel "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/app/interface/main/mcn/tool/validate"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var (
	signNotInStates = []model.MCNSignState{model.MCNSignStateOnDelete, model.MCNSignStateOnPreOpen, model.MCNSignStateOnExpire, model.MCNSignStateOnClear}
	// 从高到低，先找到的返回
	// 状态优先级：
	// 0： 封禁
	// 1： 签约
	// 2：
	statePriority = []model.MCNSignState{
		model.MCNSignStateOnBlock,
		model.MCNSignStateOnSign,
		model.MCNSignStateOnReject,
		model.MCNSignStateOnReview,
		model.MCNSignStateNoApply,
	}

	// mcnSign = mcnmodel.McnSign{}
	// mcnUp   = mcnmodel.McnUp{}

	// UpPermissionApplyCannotApplyStates 在这此状态下，不能再申请改变Up主
	UpPermissionApplyCannotApplyStates = []adminmodel.MCNUPPermissionState{
		adminmodel.MCNUPPermissionStateNoAuthorize, // 待Up主同意
		adminmodel.MCNUPPermissionStateReview,      // 待审核
	}

	// UpSignedStates up signed state
	UpSignedStates = []model.MCNUPState{
		model.MCNUPStateOnSign,
		model.MCNUPStateOnPreOpen,
	}
)

// GetMcnSignState .
// mcnList, it's all mcn sign found with the state
// state, the mcn's state of qualified, if multiple state found, will be return in priority
func (d *Dao) GetMcnSignState(fields string, mcnMid int64) (mcn *mcnmodel.McnSign, state model.MCNSignState, err error) {
	var mcnList []*mcnmodel.McnSign
	if err = d.mcndb.Select(fields).Where("mcn_mid=? and state not in(?)", mcnMid, signNotInStates).Find(&mcnList).Error; err != nil {
		err = errors.WithStack(err)
		return
	}

	if len(mcnList) == 0 {
		log.Warn("mcn not exist, mcn id=%d", mcnMid)
		err = ecode.NothingFound
		return
	}

	var stateMap = make(map[model.MCNSignState]*mcnmodel.McnSign)
	for _, v := range mcnList {
		stateMap[model.MCNSignState(v.State)] = v
	}

	for _, v := range statePriority {
		if mcnValue, ok := stateMap[v]; ok {
			state = v
			mcn = mcnValue
			break
		}
	}
	return
}

// GetUpBind .
func (d *Dao) GetUpBind(query interface{}, args ...interface{}) (upList []*mcnmodel.McnUp, err error) {
	if err = d.mcndb.Select("*").Where(query, args...).Find(&upList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			log.Error("fail to get bind up from db, err=%s", err)
			return
		}
	}
	return
}

// BindUp .
func (d *Dao) BindUp(up *mcnmodel.McnUp, sign *mcnmodel.McnSign, arg *mcnmodel.McnBindUpApplyReq) (result *mcnmodel.McnUp, affectedRow int64, err error) {
	if arg == nil || sign == nil {
		return nil, 0, ecode.ServerErr
	}
	var db *gorm.DB
	if up == nil {
		up = &mcnmodel.McnUp{
			SignID: sign.ID,
		}
	}
	arg.CopyTo(up)

	// 如果绑定自己，那么直接接受
	if sign.McnMid == arg.UpMid {
		up.State = model.MCNUPStateOnSign
		// 签约周期为MCN的签约周期
		if up.BeginDate == 0 {
			up.BeginDate = sign.BeginDate
		}
		if up.EndDate == 0 {
			up.EndDate = sign.EndDate
		}
		var (
			now  = time.Now()
			date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		)
		if up.BeginDate.Time().After(date) {
			up.State = model.MCNUPStateOnPreOpen
		}
	} else {
		// 清除状态为未授权
		up.State = model.MCNUPStateNoAuthorize
		if !validate.RegHTTPCheck.MatchString(arg.ContractLink) || !validate.RegHTTPCheck.MatchString(arg.UpAuthLink) {
			log.Error("contract link or up auth link is not http, arg=%v", arg)
			err = ecode.RequestErr
			return
		}
	}
	// 判断开始时间与结束时间
	if up.BeginDate == 0 || up.EndDate < up.BeginDate {
		log.Error("begin date is after end date, arg=%v", arg)
		err = ecode.MCNUpBindUpDateError
		return
	}
	db = d.mcndb.Save(up)
	affectedRow, err = db.RowsAffected, db.Error
	if err != nil {
		log.Error("save bind up info fail, err=%s, sign=%v", err, sign)
		err = ecode.ServerErr
	}
	result = up
	return
}

// UpdateBindUp .
func (d *Dao) UpdateBindUp(values map[string]interface{}, query interface{}, args ...interface{}) (affectedRow int64, err error) {
	var db = d.mcndb.Table(mcnmodel.TableNameMcnUp).Where(query, args...).Updates(values)
	affectedRow, err = db.RowsAffected, db.Error
	if err != nil {
		log.Error("fail to update bind up, err=%s", err)
	}
	return
}

//UpConfirm up confrim
func (d *Dao) UpConfirm(arg *mcnmodel.McnUpConfirmReq, state model.MCNUPState) (err error) {
	var tx = d.mcndb.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	var changeMap = map[string]interface{}{
		"state":             state,
		"state_change_time": time.Now(),
	}
	if arg.Choice {
		changeMap["confirm_time"] = time.Now()
	}
	err = tx.Table(mcnmodel.TableNameMcnUp).
		Where("id=? and up_mid=? and state=?", arg.BindID, arg.UpMid, model.MCNUPStateNoAuthorize).
		Updates(changeMap).Error
	if err != nil {
		log.Error("fail to  update db, err=%s", err)
		return
	}
	// 表示同意
	if arg.Choice {
		// 驳回其他的绑定请求
		err = tx.Table(mcnmodel.TableNameMcnUp).
			Where("id !=? and up_mid=? and state=?", arg.BindID, arg.UpMid, model.MCNUPStateNoAuthorize).
			Updates(map[string]interface{}{
				"state":             model.MCNUPStateOnRefuse,
				"state_change_time": time.Now(),
			}).Error
		if err != nil {
			log.Error("fail to  update db, err=%s", err)
			return
		}
	}
	return tx.Commit().Error
}

// GetBindInfo .
func (d *Dao) GetBindInfo(arg *mcnmodel.McnUpGetBindReq) (res *mcnmodel.McnGetBindReply, err error) {
	var result mcnmodel.McnGetBindReply
	err = d.mcndb.Raw(`select s.company_name, s.mcn_mid, u.up_auth_link, u.id as bind_id, u.permission as new_permission
		from mcn_up as u inner join mcn_sign as s 
		on s.id = u.sign_id 
		where u.id = ? and u.up_mid=? and u.state = 0;`, arg.BindID, arg.UpMid).Find(&result).Error
	res = &result
	return
}

//GetMcnOldInfo 获取冷却中的信息
func (d *Dao) GetMcnOldInfo(mcnMid int64) (res *mcnmodel.McnSign, err error) {
	res = new(mcnmodel.McnSign)
	err = d.mcndb.Where("mcn_mid=? and state=?", mcnMid, model.MCNSignStateOnCooling).Order("id desc").Limit(1).Find(res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("fail to get db, err=%s", err)
		return
	}
	return
}

//GetUpPermissionApply 从permission apply表中读取数据
func (d *Dao) GetUpPermissionApply(fields string, query interface{}, args ...interface{}) (res []*mcnmodel.McnUpPermissionApply, err error) {
	err = d.mcndb.Select(fields).Where(query, args...).Find(&res).Error
	if err != nil {
		log.Error("fail to get db, err=%v", err)
		return
	}
	return
}

// GetUpPermissionBindInfo .
func (d *Dao) GetUpPermissionBindInfo(arg *mcnmodel.McnUpGetBindReq) (res *mcnmodel.McnGetBindReply, err error) {
	var result mcnmodel.McnGetBindReply
	err = d.mcndb.Raw(`select s.company_name, u.mcn_mid, u.up_auth_link, u.id as bind_id, u.old_permission, u.new_permission
		from mcn_up_permission_apply as u inner join mcn_sign as s 
		on s.id = u.sign_id 
		where u.id = ? and u.up_mid=? and u.state = 0;`, arg.BindID, arg.UpMid).Find(&result).Error
	res = &result
	return
}

//UpPermissionConfirm up confrim
func (d *Dao) UpPermissionConfirm(arg *mcnmodel.McnUpConfirmReq, state adminmodel.MCNUPPermissionState) (err error) {
	var db = d.mcndb
	var changeMap = map[string]interface{}{
		"state": state,
	}

	err = db.Table(mcnmodel.TableMcnUpPermissionApply).
		Where("id=? and up_mid=? and state=?", arg.BindID, arg.UpMid, adminmodel.MCNUPPermissionStateNoAuthorize).
		Updates(changeMap).Error
	if err != nil {
		log.Error("fail to  update db, err=%s", err)
		return
	}

	return
}
