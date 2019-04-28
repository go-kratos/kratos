package upcrmdao

import (
	"github.com/jinzhu/gorm"
	"go-common/app/service/main/upcredit/conf"

	"context"
	"fmt"
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"go-common/library/log"
)

//Dao upcrm dao
type Dao struct {
	conf  *conf.Config
	crmdb *gorm.DB
}

//New create
func New(c *conf.Config) *Dao {
	var d = &Dao{
		conf: c,
	}
	crmdb, err := gorm.Open("mysql", c.DB.Upcrm.DSN)
	if crmdb == nil {
		log.Error("connect to db fail, err=%v", err)
		return nil
	}
	d.crmdb = crmdb
	crmdb.SingularTable(true)
	d.crmdb.LogMode(c.IsTest)
	return d
}

//Close close
func (d *Dao) Close() {
	if d.crmdb != nil {
		d.crmdb.Close()
	}
}

//AddLog add log
func (d *Dao) AddLog(arg *upcrmmodel.ArgCreditLogAdd) error {
	var creditLog = &upcrmmodel.CreditLog{}
	creditLog.CopyFrom(arg)
	return d.crmdb.Create(creditLog).Error
}

//AddCreditScore add score
func (d *Dao) AddCreditScore(creditScore *upcrmmodel.UpScoreHistory) error {
	return d.crmdb.Create(creditScore).Error
}

//AddOrUpdateCreditScore update score
func (d *Dao) AddOrUpdateCreditScore(creditScore *upcrmmodel.UpScoreHistory) (err error) {
	var tablename = creditScore.TableName()
	var insertSQL = fmt.Sprintf("insert into %s (mid, score_type, score, generate_date, ctime) values (?,?,?,?,?) "+
		"on duplicate key update score=?", tablename)
	err = d.crmdb.Exec(
		insertSQL,
		creditScore.Mid, creditScore.ScoreType, creditScore.Score, creditScore.GenerateDate, creditScore.CTime,
		creditScore.Score).Error
	if err != nil {
		log.Error("add credit score fail, mid=%d, err=%+v", creditScore.Mid, err)
	}
	return
}

//GetCreditScore get score
func (d *Dao) GetCreditScore(c context.Context, arg *upcrmmodel.GetScoreParam) (results []*upcrmmodel.UpScoreHistory, err error) {
	var mod = upcrmmodel.UpScoreHistory{
		Mid: arg.Mid,
	}
	err = d.crmdb.Table(mod.TableName()).Select("score, generate_date").
		Where("mid=? AND generate_date>=? AND generate_date<=? AND score_type=?", arg.Mid, arg.FromDate, arg.ToDate, arg.ScoreType).
		Find(&results).Error
	if err != nil {
		log.Error("get score history fail, arg=%+v, err=%+v", arg, err)
	}
	return
}

//GetCreditLog get log
func (d *Dao) GetCreditLog(c context.Context, arg *upcrmmodel.ArgGetLogHistory) (results []*upcrmmodel.SimpleCreditLogWithContent, err error) {
	var mod = upcrmmodel.SimpleCreditLogWithContent{}
	mod.Mid = arg.Mid

	err = d.crmdb.Table(mod.TableName()).Select("type, op_type, reason, business_type, ctime, content, oid").
		//Limit(arg.Limit).
		//Offset(arg.Offset).
		Where("mid=? AND ctime>=? AND ctime<=?", arg.Mid, arg.FromDate, arg.ToDate).
		Order("ctime").
		Find(&results).Error
	if err != nil {
		log.Error("get log history fail, arg=%+v, err=%+v", arg, err)
	}
	return
}
