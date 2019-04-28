package dao

import "go-common/app/admin/ep/tapd/model"

// AddEventLog Add Event Log.
func (d *Dao) AddEventLog(eventLog *model.EventLog) error {
	return d.db.Create(eventLog).Error
}

// UpdateEventLog Update Event Log.
func (d *Dao) UpdateEventLog(eventLog *model.EventLog) error {
	return d.db.Model(&model.EventLog{}).Where("id=?", eventLog.ID).Update(eventLog).Error
}

//FindEventLogs Find Event Logs.
func (d *Dao) FindEventLogs(req *model.QueryEventLogReq) (total int64, eventLogs []*model.EventLog, err error) {
	gDB := d.db.Model(&model.EventLog{})

	if req.WorkspaceID > 0 {
		gDB = gDB.Where("workspace_id=?", req.WorkspaceID)
	}
	if req.EventID > 0 {
		gDB = gDB.Where("event_id=?", req.EventID)
	}

	if string(req.Event) != "" {
		gDB = gDB.Where("event like ?", string(req.Event)+_wildcards)
	}

	if err = gDB.Count(&total).Error; err != nil {
		return
	}

	err = gDB.Order("ctime desc").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&eventLogs).Error

	return
}
