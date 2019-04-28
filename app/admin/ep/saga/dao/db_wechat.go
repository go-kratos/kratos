package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"

	"github.com/pkg/errors"
	pkgerr "github.com/pkg/errors"
)

const (
	_where                  = "WHERE"
	_and                    = "AND"
	_wildcards              = "%"
	_contactRightJoinLogSQL = "SELECT m.id,m.name,ml.username,ml.operation_type,ml.operation_result,ml.ctime FROM machines AS m INNER JOIN machine_logs AS ml ON m.id = ml.machine_id"
	_contactRightCountSQL   = "SELECT count(m.id) FROM machines AS m INNER JOIN machine_logs AS ml ON m.id = ml.machine_id"
)

//var (
//	regUserID = regexp.MustCompile(`^\d+$`)
//)

// QueryUserByUserName query user by user name
func (d *Dao) QueryUserByUserName(userName string) (contactInfo *model.ContactInfo, err error) {
	contactInfo = &model.ContactInfo{}
	err = pkgerr.WithStack(d.db.Where(&model.ContactInfo{UserName: userName}).First(contactInfo).Error)
	return
}

// QueryUserByID query user by user ID
func (d *Dao) QueryUserByID(userID string) (contactInfo *model.ContactInfo, err error) {
	contactInfo = &model.ContactInfo{}
	err = pkgerr.WithStack(d.db.Where(&model.ContactInfo{UserID: userID}).First(contactInfo).Error)
	return
}

// UserIds query user ids for the user names 从数据库表ContactInfo查询员工编号
func (d *Dao) UserIds(userNames []string) (userIds string, err error) {
	var (
		userName    string
		ids         []string
		contactInfo *model.ContactInfo
	)
	if len(userNames) == 0 {
		err = errors.Errorf("UserIds: userNames is empty!")
		return
	}

	for _, userName = range userNames {
		if contactInfo, err = d.QueryUserByUserName(userName); err != nil {
			err = errors.Wrapf(err, "UserIds: no such user (%s) in db, err (%s)", userName, err.Error())
			return
		}

		log.Info("UserIds: username (%s), userid (%s)", userName, contactInfo.UserID)
		if contactInfo.UserID != "" {
			ids = append(ids, contactInfo.UserID)
		}
	}

	if len(ids) <= 0 {
		err = errors.Wrapf(err, "UserIds: failed to find all the users in db, what a pity!")
		return
	}

	userIds = strings.Join(ids, "|")

	return
}

// ContactInfos query all the records in contact_infos
func (d *Dao) ContactInfos() (contactInfos []*model.ContactInfo, err error) {
	err = pkgerr.WithStack(d.db.Find(&contactInfos).Error)
	return
}

// FindContacts find application record by auditor.
func (d *Dao) FindContacts(pn, ps int) (total int64, ars []*model.ContactInfo, err error) {
	cdb := d.db.Model(&model.ContactInfo{}).Order("ID desc").Offset((pn - 1) * ps).Limit(ps)
	if err = pkgerr.WithStack(cdb.Find(&ars).Error); err != nil {
		return
	}
	if err = pkgerr.WithStack(d.db.Model(&model.ContactInfo{}).Count(&total).Error); err != nil {
		return
	}
	return
}

// CreateContact create contact info record
func (d *Dao) CreateContact(contact *model.ContactInfo) (err error) {
	return pkgerr.WithStack(d.db.Create(contact).Error)
}

// DelContact delete the contact info with the specified UserID
func (d *Dao) DelContact(contact *model.ContactInfo) (err error) {
	return pkgerr.WithStack(d.db.Delete(contact).Error)
}

// UptContact update the contact information
func (d *Dao) UptContact(contact *model.ContactInfo) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.ContactInfo{}).Where(&model.ContactInfo{UserID: contact.UserID}).Updates(*contact).Error)
}

// InsertContactLog insert machine log.
func (d *Dao) InsertContactLog(contactlog *model.ContactLog) (err error) {
	return pkgerr.WithStack(d.db.Create(contactlog).Error)
}

// FindMachineLogs Find Machine Logs.
func (d *Dao) FindMachineLogs(queryRequest *model.QueryContactLogRequest) (total int64, machineLogs []*model.AboundContactLog, err error) {
	var (
		qSQL = _contactRightJoinLogSQL
		cSQL = _contactRightCountSQL
		rows *sql.Rows
	)

	if queryRequest.UserID > 0 || queryRequest.UserName != "" || queryRequest.OperateType != "" || queryRequest.OperateUser != "" {
		var (
			strSQL      = ""
			logicalWord = _where
		)

		if queryRequest.UserID > 0 {
			strSQL = fmt.Sprintf("%s %s  ml.machine_id = %s", strSQL, logicalWord, strconv.FormatInt(queryRequest.UserID, 10))
			logicalWord = _and
		}

		if queryRequest.UserName != "" {
			strSQL = fmt.Sprintf("%s %s  m.name like '%s'", strSQL, logicalWord, _wildcards+queryRequest.UserName+_wildcards)
			logicalWord = _and
		}

		if queryRequest.OperateType != "" {
			strSQL = fmt.Sprintf("%s %s  ml.operation_type like '%s'", strSQL, logicalWord, _wildcards+queryRequest.OperateType+_wildcards)
			logicalWord = _and
		}

		if queryRequest.OperateUser != "" {
			strSQL = fmt.Sprintf("%s %s  ml.username like '%s'", strSQL, logicalWord, _wildcards+queryRequest.OperateUser+_wildcards)
			logicalWord = _and
		}

		qSQL = _contactRightJoinLogSQL + " " + strSQL
		cSQL = _contactRightCountSQL + " " + strSQL

	}

	cDB := d.db.Raw(cSQL)
	if err = pkgerr.WithStack(cDB.Count(&total).Error); err != nil {
		return
	}
	gDB := d.db.Raw(qSQL)
	if rows, err = gDB.Order("ml.ctime DESC").Offset((queryRequest.PageNum - 1) * queryRequest.PageSize).Limit(queryRequest.PageSize).Rows(); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ml := &model.AboundContactLog{}
		if err = rows.Scan(&ml.MachineID, &ml.Name, &ml.Username, &ml.OperateType, &ml.OperateResult, &ml.OperateTime); err != nil {
			return
		}
		machineLogs = append(machineLogs, ml)
	}

	return
}

// AddWechatCreateLog ...
func (d *Dao) AddWechatCreateLog(req *model.WechatCreateLog) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

// QueryWechatCreateLog ...
func (d *Dao) QueryWechatCreateLog(ifpage bool, req *model.Pagination, wechatCreateInfo *model.WechatCreateLog) (wechatCreateInfos []*model.WechatCreateLog, total int, err error) {

	gDB := d.db.Table("wechat_create_logs").Where(wechatCreateInfo).Model(&model.WechatCreateLog{})
	if err = errors.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}
	if ifpage {
		gDB = gDB.Order("id DESC").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&wechatCreateInfos)
	} else {
		gDB = gDB.Find(&wechatCreateInfos)
	}
	if gDB.Error != nil {
		if gDB.RecordNotFound() {
			err = nil
		} else {
			err = errors.WithStack(gDB.Error)
		}
	}
	return
}

// CreateChatLog ...
func (d *Dao) CreateChatLog(req *model.WechatChatLog) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}

// CreateMessageLog ...
func (d *Dao) CreateMessageLog(req *model.WechatMessageLog) (err error) {
	return pkgerr.WithStack(d.db.Create(req).Error)
}
