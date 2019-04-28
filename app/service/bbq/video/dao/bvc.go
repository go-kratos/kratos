package dao

import (
	"context"
	"go-common/app/service/bbq/video/model"
	"regexp"
)

const (
	_insertRecord = "insert bvc_flow_record (`bvcid`,`svid`,`type`) values (?,?,?)"
	_updateRecord = "update bvc_flow_record set `svid` = ?,`type` = ? where `bvcid` = ?"
)

// AddOrUpdateFlowRecord 添加bvc flow记录
func (d *Dao) AddOrUpdateFlowRecord(c context.Context, r *model.BVCRecord) error {
	_, err := d.db.Exec(c, _insertRecord, r.FLowID, r.SVID, r.Type)
	if err != nil {
		if matched, _ := regexp.MatchString("Duplicate entry", err.Error()); matched {
			_, err = d.db.Exec(c, _updateRecord, r.SVID, r.Type, r.FLowID)
		}
	}
	return err
}
