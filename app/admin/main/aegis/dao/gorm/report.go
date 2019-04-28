package gorm

import (
	"context"
	"database/sql"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/log"
)

//ReportTaskMetas 任务数据统计
func (d *Dao) ReportTaskMetas(c context.Context, bt string, et string, bizid, flowid int64, uids []int64, mnames map[int64]string, tp int8) (metas []*model.ReportMeta, missuid []int64, err error) {
	db := d.orm.Table("task_report").Select("mtime,uid,type,content").Where("business_id=? AND type=?", bizid, tp)
	if flowid != 0 {
		db.Where("flow_id=?", flowid)
	}
	db = db.Where("mtime>=? AND mtime<?", bt, et)

	if len(uids) > 0 {
		db = db.Where("uid IN (?)", uids)
	}

	var rows *sql.Rows
	if rows, err = db.Order("mtime asc").Rows(); err != nil {
		log.Error("ReportTaskFlow error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		meta := &model.ReportMeta{}
		if err = rows.Scan(&meta.Mtime, &meta.UID, &meta.Type, &meta.Content); err != nil {
			return
		}
		if uname, ok := mnames[meta.UID]; ok {
			meta.Uname = uname
		} else {
			missuid = append(missuid, meta.UID)
		}
		metas = append(metas, meta)
	}
	return
}

//TaskReports 任务报表记录
func (d *Dao) TaskReports(c context.Context, biz int64, flowID int64, tp []int8, statdateFrom string, statdateTo string) (res []*model.TaskReport, err error) {
	res = []*model.TaskReport{}
	db := d.orm
	if statdateFrom != "" {
		db = db.Where("stat_date>=?", statdateFrom)
	}
	if statdateTo != "" {
		db = db.Where("stat_date <=?", statdateTo)
	}
	db = db.Where("business_id=?", biz)
	if flowID > 0 {
		db = db.Where("flow_id=?", flowID)
	}
	if len(tp) >= 0 {
		db = db.Where("type in (?)", tp)
	}
	db = db.Order("stat_date desc, business_id desc, flow_id desc")
	if err = db.Find(&res).Error; err != nil {
		log.Error("TaskReports error(%v)", err)
	}
	return
}
