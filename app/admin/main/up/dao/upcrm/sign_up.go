package upcrm

import (
	"time"

	"go-common/app/admin/main/up/model/signmodel"
	"go-common/app/admin/main/up/model/upcrmmodel"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

const (
	// PayStateUnpay not pay
	PayStateUnpay = 0
	// PayStatePayed payed
	PayStatePayed = 1
)

// InsertSignUp insert sign up
func (d *Dao) InsertSignUp(db *gorm.DB, up *signmodel.SignUp) (affectedRow int64, err error) {
	var handle = db.Save(up)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// InsertPayInfo inert pay
func (d *Dao) InsertPayInfo(db *gorm.DB, info *signmodel.SignPay) (affectedRow int64, err error) {
	var handle = db.Save(info)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// InsertTaskInfo insert task
func (d *Dao) InsertTaskInfo(db *gorm.DB, info *signmodel.SignTask) (affectedRow int64, err error) {
	var handle = db.Save(info)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// InsertContractInfo insert contract
func (d *Dao) InsertContractInfo(db *gorm.DB, info interface{}) (affectedRow int64, err error) {
	var handle = db.Save(info)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// DelPayInfo update payinfo
func (d *Dao) DelPayInfo(db *gorm.DB, ids []int64) (affectedRow int64, err error) {
	var handle = db.Model(&signmodel.SignPay{}).Where("id IN (?)", ids).Update("state", 100)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// DelTaskInfo update taskInfo
func (d *Dao) DelTaskInfo(db *gorm.DB, ids []int64) (affectedRow int64, err error) {
	var handle = db.Model(&signmodel.SignTask{}).Where("id IN (?)", ids).Update("state", 100)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// DelSignContract update signcontract
func (d *Dao) DelSignContract(db *gorm.DB, ids []int64) (affectedRow int64, err error) {
	var handle = db.Model(&signmodel.SignContract{}).Where("id IN (?)", ids).Update("state", 100)
	err = handle.Error
	affectedRow = handle.RowsAffected
	return
}

// SignUpID .
func (d *Dao) SignUpID(sigID int64) (su *signmodel.SignUp, msp map[int64]*signmodel.SignPay, mst map[int64]*signmodel.SignTask, msc map[int64]*signmodel.SignContract, err error) {
	su = &signmodel.SignUp{}
	if err = d.crmdb.Table(signmodel.TableSignUp).Where("id = ? AND state IN (0,1)", sigID).Find(su).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
		su, msp, mst, msc = nil, nil, nil, nil
		return
	}
	var (
		sps []*signmodel.SignPay
		sts []*signmodel.SignTask
		scs []*signmodel.SignContract
	)
	msp = make(map[int64]*signmodel.SignPay)
	if err = d.crmdb.Table(signmodel.TableSignPay).Where("sign_id = ? AND state IN (0,1)", sigID).Find(&sps).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	for _, v := range sps {
		msp[v.ID] = v
	}
	mst = make(map[int64]*signmodel.SignTask)
	if err = d.crmdb.Table(signmodel.TableSignTask).Where("sign_id = ? AND state = 0", sigID).Find(&sts).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	for _, v := range sts {
		mst[v.ID] = v
	}
	msc = make(map[int64]*signmodel.SignContract)
	if err = d.crmdb.Table(signmodel.TableSignContract).Where("sign_id = ? AND state = 0", sigID).Find(&scs).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db fail, err=%+v", err)
		return
	}
	for _, v := range scs {
		msc[v.ID] = v
	}
	return
}

// GetSignIDByCondition get sign id
// arg query args
func (d *Dao) GetSignIDByCondition(arg *signmodel.SignQueryArg) (signIDs []uint32, err error) {
	var signIDMap = map[uint32]struct{}{}

	switch {
	default:

		// 如果是mid，则不进行其他的查询
		if arg.Mid != 0 {
			var result []signmodel.SignUpOnlyID
			var db = d.crmdb.Table("sign_up").Where("mid=?", arg.Mid)
			err = db.Select("id").Find(&result).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error("db fail, err=%+v", err)
				return
			}
			for _, v := range result {
				signIDMap[v.ID] = struct{}{}
			}
			break
		}
		var now = time.Now()
		// 1.增加查询条件
		if arg.DuePay == 1 {
			var duedate = now.AddDate(0, 0, 7)
			var result []signmodel.SignUpOnlySignID
			// due_date and state = PayStateUnpay
			err = d.crmdb.Table("sign_pay").Select("sign_id").Where("due_date <= ? and state = 0", duedate.Format(upcrmmodel.TimeFmtDate)).Where("state != 100").
				Find(&result).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error("db fail, err=%+v", err)
			} else {
				for _, v := range result {
					signIDMap[v.SignID] = struct{}{}
				}
			}
		}

		if arg.DueSign == 1 || arg.ExpireSign == 1 {
			var duedate = now.AddDate(0, 0, 30)
			var result []signmodel.SignUpOnlyID
			var db = d.crmdb.Table("sign_up").Select("id")
			if arg.DueSign == 1 && arg.ExpireSign == 1 {
				db = db.Where("end_date < ?", duedate.Format(upcrmmodel.TimeFmtDate)).Where("state != 100") // 去掉已删除的
			} else {
				if arg.DueSign == 1 {
					db = db.Where("end_date >= ? and end_date <= ?", now.Format(upcrmmodel.TimeFmtDate), duedate.Format(upcrmmodel.TimeFmtDate))
				} else if arg.ExpireSign == 1 {
					db = db.Where("end_date < ?", now.Format(upcrmmodel.TimeFmtDate))
				}
			}
			err = db.
				Find(&result).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				log.Error("db fail, err=%+v", err)
			} else {
				for _, v := range result {
					signIDMap[v.ID] = struct{}{}
				}
			}
		}
	}

	for k := range signIDMap {
		signIDs = append(signIDs, k)
	}

	return
}

//GetSignUpByID signid 可以是nil，如果是nil，则会取所有的信息
// query, args, 额外的查询条件
func (d *Dao) GetSignUpByID(signID []uint32, order string, offset int, limit int, query interface{}, args ...interface{}) (result []signmodel.SignUp, err error) {
	var db = d.crmdb.Table("sign_up")
	if signID != nil {
		db = db.Where("id in (?)", signID)
	}
	if query != nil {
		db = db.Where(query, args)
	}
	if order != "" {
		db = db.Order(order)
	}
	err = db.
		Offset(offset).
		Limit(limit).
		Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//GetSignUpCount get sign up's count
func (d *Dao) GetSignUpCount(query string, args ...interface{}) (count int) {
	d.crmdb.Table(signmodel.TableSignUp).Where(query, args).Count(&count)
	return
}

//GetTask get task by sign id and state
func (d *Dao) GetTask(signID []uint32, state ...int) (result []signmodel.SignTask, err error) {
	err = d.crmdb.Table(signmodel.TableSignTask).Where("sign_id in (?) and state = ?", signID, state).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//GetPay get get sign id
func (d *Dao) GetPay(signID []uint32) (result []signmodel.SignPay, err error) {
	err = d.crmdb.Table(signmodel.TableSignPay).Where("sign_id in (?)", signID).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//GetContract get get sign id
func (d *Dao) GetContract(signID []uint32) (result []signmodel.SignContract, err error) {
	err = d.crmdb.Table(signmodel.TableSignContract).Where("sign_id in (?)", signID).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//PayComplete finish pay by pay id
func (d *Dao) PayComplete(ids []int64) (affectedRow int64, err error) {
	var db = d.crmdb.Table(signmodel.TableSignPay).Where("id in (?)", ids).Updates(map[string]interface{}{"state": PayStatePayed})
	err = db.Error
	if err == nil {
		affectedRow = db.RowsAffected
	}
	return
}

//SignWithName sign with name, used to send mail
type SignWithName struct {
	signmodel.SignUp
	Name string
}

//GetDueSignUp check due
// expireAfterDays : how many days to expire
func (d *Dao) GetDueSignUp(now time.Time, expireAfterDays int) (result []*SignWithName, err error) {
	var dueDate = now.AddDate(0, 0, expireAfterDays)
	var db = d.crmdb.Table(signmodel.TableSignUp).Select("id, mid, end_date, admin_id, admin_name")
	// email_state = 0 是未发过邮件的意思
	db = db.Where("end_date <= ? and end_date >= ? and email_state = 0", dueDate.Format(upcrmmodel.TimeFmtDate), now.Format(upcrmmodel.TimeFmtDate))
	err = db.Find(&result).Error
	return
}

//PayWithAdmin  pay with name, used to send mail
type PayWithAdmin struct {
	signmodel.SignPay
	Name      string
	AdminID   int
	AdminName string
}

//GetDuePay check due
func (d *Dao) GetDuePay(now time.Time, expireAfterDays int) (result []*PayWithAdmin, err error) {
	var dueDate = now.AddDate(0, 0, expireAfterDays)
	err = d.crmdb.Raw("select p.id, "+
		"p.sign_id, "+
		"p.mid, "+
		"p.due_date, "+
		"p.pay_value,"+
		"s.admin_id, "+
		"s.admin_name"+
		" from sign_pay as p left join sign_up as s on p.sign_id = s.id "+
		" where p.due_date <= ? and p.email_state = 0 and p.state = 0", dueDate.Format(upcrmmodel.TimeFmtDate)).
		Scan(&result).Error
	return
}

//UpdateEmailState update email send state
// state : @
func (d *Dao) UpdateEmailState(table string, ids []int64, state int8) (affectedRow int64, err error) {
	var db = d.crmdb.Table(table).Where("id in (?)", ids).Update("email_state", state)
	err = db.Error
	if err == nil {
		affectedRow = db.RowsAffected
	} else {
		log.Error("update email state fail, err=%+v", err)
	}
	return
}

//CheckUpHasValidContract check if has valid contract
func (d *Dao) CheckUpHasValidContract(mid int64, date time.Time) (exist bool, err error) {
	var ids []struct {
		ID int
	}
	err = d.crmdb.Table(signmodel.TableSignUp).Select("id").Where("mid=? and end_date>=?", mid, date.Format(upcrmmodel.TimeFmtDate)).Limit(1).Find(&ids).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("check exist from db fail, err=%+v", err)
		return
	}
	exist = len(ids) > 0
	return
}

// GetOrCreateTaskHistory .
func (d *Dao) GetOrCreateTaskHistory(db *gorm.DB, st *signmodel.SignTask) (sth *signmodel.SignTaskHistory, init bool, err error) {
	sDate, _ := signmodel.GetTaskDuration(time.Now(), st.TaskType)
	sth = new(signmodel.SignTaskHistory)
	err = db.Select("*").Where("task_template_id=? and generate_date=?", st.ID, sDate).Find(&sth).Error
	// 创建一条，如果没找到的话
	if err == gorm.ErrRecordNotFound {
		sth = &signmodel.SignTaskHistory{
			Mid:            st.Mid,
			SignID:         int64(st.SignID),
			TaskTemplateID: int(st.ID),
			TaskType:       st.TaskType,
			TaskCondition:  int(st.TaskCondition),
			GenerateDate:   xtime.Time(sDate.Unix()),
			Attribute:      st.Attribute,
			State:          signmodel.SignTaskStateRunning,
		}
		if err = db.Save(&sth).Error; err != nil {
			log.Error("create task history fail, err=%v, task=%v", err, st)
		}
		init = true
	}
	return
}

// UpSignTaskHistory .
func (d *Dao) UpSignTaskHistory(db *gorm.DB, sth *signmodel.SignTaskHistory) (err error) {
	if err = db.Table(signmodel.TableSignTaskHistory).Where("id = ?", sth.ID).UpdateColumns(
		map[string]interface{}{
			"task_type":      sth.TaskType,
			"task_condition": sth.TaskCondition,
			"attribute":      sth.Attribute,
			"mtime":          time.Now(),
		}).Error; err != nil {
		log.Error("dao.UpSignTaskHistory(%+v) , err=%+v", sth, err)
	}
	return
}
