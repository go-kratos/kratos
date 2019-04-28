package upcrm

import (
	"time"

	"go-common/app/admin/main/up/util/now"
	"go-common/app/job/main/up/model/signmodel"
	"go-common/app/job/main/up/model/upcrmmodel"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

const (
	//PayStateUnpay not pay
	PayStateUnpay = 0
	//PayStatePayed payed
	PayStatePayed = 1
)

// var days
var (
	Day   = time.Hour * 24
	Day3  = Day * 3
	Week  = Day * 7
	Week3 = Week * 3
)

// GetTaskDuration this will return task duration, [startDate, endDate)
func GetTaskDuration(date time.Time, taskType int8) (startDate, endDate time.Time) {
	var ndate = now.New(date)
	now.WeekStartDay = time.Monday
	switch taskType {
	case signmodel.TaskTypeDay:
		var begin = ndate.BeginningOfDay()
		return begin, begin.AddDate(0, 0, 1)
	case signmodel.TaskTypeWeek:
		var begin = ndate.BeginningOfWeek()
		return begin, begin.AddDate(0, 0, 7)
	case signmodel.TaskTypeMonth:
		var begin = ndate.BeginningOfMonth()
		return begin, begin.AddDate(0, 1, 0)
	case signmodel.TaskTypeQuarter:
		var begin = ndate.BeginningOfQuarter()
		return begin, begin.AddDate(0, 3, 0)
	}
	return
}

// GetTaskExpireLimit this will return generate date range that needs to send email
func GetTaskExpireLimit(today time.Time, taskType int8) (minDate, maxDate time.Time) {
	// 0<=（endDate-today)<= limit
	//  0<=(generateDate + duration - today) <= limit
	//  today-duration <= generateDate <= today - duration + limit
	// duration = endDate - startDate
	switch taskType {
	case signmodel.TaskTypeDay:
		var tmp = today.AddDate(0, 0, -1)
		return tmp, tmp.Add(Day)
	case signmodel.TaskTypeWeek:
		var tmp = today.AddDate(0, 0, -7)
		return tmp, tmp.Add(Day3)
	case signmodel.TaskTypeMonth:
		var tmp = today.AddDate(0, -1, 0)
		return tmp, tmp.Add(Week)
	case signmodel.TaskTypeQuarter:
		var tmp = today.AddDate(0, -3, 0)
		return tmp, tmp.Add(Week3)
	}
	return
}

//InsertSignUp insert sign up
// up : sign up
func (d *Dao) InsertSignUp(up *signmodel.SignUp) (affectedRow int64, err error) {
	var db = d.crmdb.Save(up)
	err = db.Error
	affectedRow = db.RowsAffected
	return
}

//InsertPayInfo inert pay
func (d *Dao) InsertPayInfo(info *signmodel.SignPay) (affectedRow int64, err error) {
	var db = d.crmdb.Save(info)
	err = db.Error
	affectedRow = db.RowsAffected
	return
}

//InsertTaskInfo insert task
func (d *Dao) InsertTaskInfo(info *signmodel.SignTask) (affectedRow int64, err error) {
	var db = d.crmdb.Save(info)
	err = db.Error
	affectedRow = db.RowsAffected
	return
}

//InsertContractInfo insert contract
func (d *Dao) InsertContractInfo(info interface{}) (affectedRow int64, err error) {
	var db = d.crmdb.Save(info)
	err = db.Error
	affectedRow = db.RowsAffected
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
	d.crmdb.Table(signmodel.TableNameSignUp).Where(query, args).Count(&count)
	return
}

//GetTask get task by sign id and state
func (d *Dao) GetTask(signID []uint32, state ...int) (result []signmodel.SignTask, err error) {
	err = d.crmdb.Table(signmodel.TableNameSignTask).Where("sign_id in (?) and state = ?", signID, state).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//GetPay get get sign id
func (d *Dao) GetPay(signID []uint32) (result []signmodel.SignPay, err error) {
	err = d.crmdb.Table(signmodel.TableNameSignPay).Where("sign_id in (?)", signID).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//GetContract get get sign id
func (d *Dao) GetContract(signID []uint32) (result []signmodel.SignContract, err error) {
	err = d.crmdb.Table(signmodel.TableNameSignContract).Where("sign_id in (?)", signID).Find(&result).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

//PayComplete finish pay by pay id
func (d *Dao) PayComplete(ids []uint32) (affectedRow int64, err error) {
	var db = d.crmdb.Table(signmodel.TableNameSignPay).Where("id in (?)", ids).Updates(map[string]interface{}{"state": PayStatePayed})
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

//GetEndDate used for template
func (s *SignWithName) GetEndDate() string {
	return s.EndDate.Time().Format(upcrmmodel.TimeFmtDate)
}

//GetDueSignUp check due
// expireAfterDays : how many days to expire
func (d *Dao) GetDueSignUp(now time.Time, expireAfterDays int) (result []*SignWithName, err error) {
	var dueDate = now.AddDate(0, 0, expireAfterDays)
	var db = d.crmdb.Table(signmodel.TableNameSignUp).Select("id, mid, end_date, admin_id, admin_name")
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

//GetEndDate used for template
func (s *PayWithAdmin) GetEndDate() string {
	return s.DueDate.Time().Format(upcrmmodel.TimeFmtDate)
}

//GetPayValue for template
func (s *PayWithAdmin) GetPayValue() float64 {
	return float64(s.PayValue) / 1000.0
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

//TaskWithAdmin task with admin
type TaskWithAdmin struct {
	signmodel.SignTaskHistory
	Name      string
	AdminID   int
	AdminName string
}

//EndDate used for template
func (s *TaskWithAdmin) EndDate() xtime.Time {
	var _, end = GetTaskDuration(s.GenerateDate.Time(), s.TaskType)
	return xtime.Time(end.Unix())
}

//TypeDesc descrption
func (s *TaskWithAdmin) TypeDesc() string {
	return signmodel.TaskTypeStr(int(s.TaskType))
}

var (
	//NeedCheckTaskType task type
	NeedCheckTaskType = []int8{
		signmodel.TaskTypeWeek,
		signmodel.TaskTypeMonth,
		signmodel.TaskTypeQuarter,
	}
)

//GetDueTask get due tasks
func (d *Dao) GetDueTask(now time.Time) (result []*TaskWithAdmin, err error) {
	// 到期的任务，
	// 从sign_task_history表中查询任务
	for _, t := range NeedCheckTaskType {
		var min, max = GetTaskExpireLimit(now, t)
		// 从sign_task_history中查询min<=generate_date <=max and task_type = t and state != 2
		var tasks []*TaskWithAdmin
		err = d.crmdb.Raw(`
select t.*, u.admin_id, u.admin_name 
from sign_task_history as t 
left join sign_up as u 
on t.sign_id=u.id 
where t.generate_date=? and t.task_type=? and t.state!=2`, // t.generate_date=? 现在只在时间刚好达到时发送邮件，这样可以只发送一次，但是有可能会漏掉那些早已经过期的。
			max.Format(upcrmmodel.TimeFmtDate), t).
			Find(&tasks).Error
		if err != nil {
			log.Error("fail to get from db, task type=%d, min=%s, max=%s", t, min, max)
			return
		}

		log.Info("get from db, task type=%d, min=%s, max=%s, len=%d", t, min, max, len(tasks))

		result = append(result, tasks...)
	}
	return
}

//UpdateEmailState update email send state
// state : @
func (d *Dao) UpdateEmailState(table string, ids []uint32, state int8) (affectedRow int64, err error) {
	if len(ids) == 0 {
		log.Warn("no ids to update email state, state=%d", state)
		return
	}
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
	err = d.crmdb.Table(signmodel.TableNameSignUp).Select("id").Where("mid=? and end_date>=?", mid, date.Format(upcrmmodel.TimeFmtDate)).Limit(1).Find(&ids).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("check exist from db fail, err=%+v", err)
		return
	}
	exist = len(ids) > 0
	return
}
