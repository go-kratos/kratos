package mcndao

import (
	"time"

	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	dateFmt    = "2006-01-02"
	tidSummary = 65535 // 特殊的tid表示所有分区数据之和
)

//GetMcnDataSummary .
func (d *Dao) GetMcnDataSummary(selec string, query interface{}, args ...interface{}) (res *mcnmodel.McnDataSummary, err error) {
	res = new(mcnmodel.McnDataSummary)
	err = d.mcndb.Select(selec).Where(query, args...).Limit(1).Find(res).Error
	if err != nil {
		log.Error("query db fail, err=%s", err)
		return
	}
	return
}

//GetMcnDataSummaryWithDiff get data with datacenter diff
func (d *Dao) GetMcnDataSummaryWithDiff(signID int64, dataTYpe mcnmodel.McnDataType, generateDate time.Time) (res *mcnmodel.McnGetDataSummaryReply, err error) {
	dataDay0, err := d.GetMcnDataSummary("up_count, fans_count_accumulate, archive_count_accumulate, play_count_accumulate, generate_date", "sign_id=? and data_type=? and generate_date=? and active_tid=?", signID, dataTYpe, generateDate.Format(dateFmt), tidSummary)
	if int64(dataDay0.GenerateDate) == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			log.Warn("not found generate date=%s, sign id=%d", generateDate.Format(dateFmt), signID)
			return
		}
		log.Error("fail to get data db, err=%s, sign id=%d", err, signID)
		return
	}

	dataDay1, err := d.GetMcnDataSummary("up_count, fans_count_accumulate, archive_count_accumulate, play_count_accumulate, generate_date", "sign_id=? and data_type=? and generate_date=? and active_tid=?", signID, dataTYpe, generateDate.AddDate(0, 0, -1).Format(dateFmt), tidSummary)
	if int64(dataDay1.GenerateDate) == 0 {
		err = gorm.ErrRecordNotFound
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			log.Warn("not found generate date=%s, sign id=%d", generateDate.Format(dateFmt), signID)
			dataDay1 = new(mcnmodel.McnDataSummary)
		} else {
			log.Error("fail to get data db, err=%s, sign id=%d", err, signID)
			return
		}
	}

	res = new(mcnmodel.McnGetDataSummaryReply)
	res.CopyFrom(dataDay0)
	res.CalcDiff(dataDay1)
	return
}

//GetDataUpLatestDate .
func (d *Dao) GetDataUpLatestDate(dataType mcnmodel.DataType, signID int64) (generateDate time.Time, err error) {
	var model = mcnmodel.McnDataUp{}
	err = d.mcndb.Select("generate_date").Where("data_type=? and sign_id=?", dataType, signID).Order("generate_date desc").Limit(1).Find(&model).Error
	if err != nil {
		log.Error("get latest date from mcn_data_up fail, err=%s", err)
		return
	}

	generateDate = model.GenerateDate.Time()
	return
}

//GetAllUpData .
func (d *Dao) GetAllUpData(signID int64, upmid int64, generateDate time.Time) (res []*mcnmodel.McnUpDataInfo, err error) {
	var sqlstr = `select 
		u.begin_date,
		u.end_date,
		u.state,
		u.up_mid,
		u.publication_price,
		u.permission,
		d.fans_increase_accumulate,
		d.archive_count,
		d.play_count,
		d.fans_increase_month,
		d.fans_count,
		d.fans_count_active,
		generate_date
		from
		mcn_up as u left join mcn_data_up as d
		on (d.sign_id = u.sign_id and d.up_mid = u.up_mid and d.data_type=1 and generate_date=? )
		where u.sign_id=? and u.state not in (?)`
	var values = []interface{}{generateDate, signID, []model.MCNUPState{model.MCNUPStateOnDelete}}
	if upmid != 0 {
		sqlstr += " and u.up_mid=?"
		values = append(values, upmid)
	}
	err = d.mcndb.Raw(sqlstr, values...).Find(&res).Error
	if err != nil {
		log.Error("query db fail, err=%s", err)
		return
	}
	return
}

//GetAllUpDataTemp .
func (d *Dao) GetAllUpDataTemp(signID int64, upmid int64, generateDate time.Time) (res []*mcnmodel.McnUpDataInfo, err error) {
	/*
		前台这边一期的数据
		1.首页 - 绑定up主总数
		2.up主列表 - 总粉数
		3.up主列表 - 投稿数
		4.up主列表 - up分区
		5.up主列表 - 签约及到期时间
	*/
	var sqlstr = `select 
		u.begin_date,
		u.end_date,
		u.state,
		u.up_mid,
		d.article_count_accumulate as archive_count,
		d.fans_count,
		d.active_fans as fans_count_active
		from
		mcn_up as u left join up_base_info as d
		on (d.mid = u.up_mid and d.business_type = 1)
		where u.sign_id=? and u.state not in (?)`
	var values = []interface{}{signID, []model.MCNUPState{model.MCNUPStateOnDelete, model.MCNUPStateOnExpire, model.MCNUPStateOnClear}}
	if upmid != 0 {
		sqlstr += " and u.up_mid=?"
		values = append(values, upmid)
	}
	err = d.mcndb.Raw(sqlstr, values...).Find(&res).Error
	if err != nil {
		log.Error("query db fail, err=%s", err)
		return
	}
	return
}
