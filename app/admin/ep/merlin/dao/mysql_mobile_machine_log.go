package dao

import (
	"database/sql"
	"fmt"
	"strconv"

	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

const (
	_mobileRightJoinLogSQL = "SELECT m.id,m.serial,ml.username,ml.operation_type,ml.operation_result,ml.ctime FROM mobile_machines AS m INNER JOIN mobile_machine_logs AS ml ON m.id = ml.machine_id"
	_mobileRightCountSQL   = "SELECT count(m.id) FROM mobile_machines AS m INNER JOIN mobile_machine_logs AS ml ON m.id = ml.machine_id"
)

// InsertMobileMachineLog Insert Mobile Machine Log.
func (d *Dao) InsertMobileMachineLog(mobileMachineLog *model.MobileMachineLog) (err error) {
	return pkgerr.WithStack(d.db.Create(mobileMachineLog).Error)
}

// FindMobileMachineLogs Find Mobile Machine Logs.
func (d *Dao) FindMobileMachineLogs(queryRequest *model.QueryMobileMachineLogRequest) (total int64, mobileMachineLogs []*model.AboundMobileMachineLog, err error) {
	var (
		qSQL = _mobileRightJoinLogSQL
		cSQL = _mobileRightCountSQL
		rows *sql.Rows
	)

	if queryRequest.MachineID > 0 || queryRequest.OperateType != "" || queryRequest.OperateUser != "" || queryRequest.Serial != "" {
		var (
			strSQL      = ""
			logicalWord = _where
		)

		if queryRequest.MachineID > 0 {
			strSQL = fmt.Sprintf("%s %s  ml.machine_id = %s", strSQL, logicalWord, strconv.FormatInt(queryRequest.MachineID, 10))
			logicalWord = _and
		}

		if queryRequest.OperateType != "" {
			strSQL = fmt.Sprintf("%s %s  ml.operation_type like '%s'", strSQL, logicalWord, _wildcards+queryRequest.OperateType+_wildcards)
			logicalWord = _and
		}

		if queryRequest.Serial != "" {
			strSQL = fmt.Sprintf("%s %s  m.serial like '%s'", strSQL, logicalWord, _wildcards+queryRequest.Serial+_wildcards)
			logicalWord = _and
		}

		if queryRequest.OperateUser != "" {
			strSQL = fmt.Sprintf("%s %s  ml.username like '%s'", strSQL, logicalWord, _wildcards+queryRequest.OperateUser+_wildcards)
			logicalWord = _and
		}

		qSQL = _mobileRightJoinLogSQL + " " + strSQL
		cSQL = _mobileRightCountSQL + " " + strSQL

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
		ml := &model.AboundMobileMachineLog{}
		if err = rows.Scan(&ml.MachineID, &ml.Serial, &ml.Username, &ml.OperateType, &ml.OperateResult, &ml.OperateTime); err != nil {
			return
		}
		mobileMachineLogs = append(mobileMachineLogs, ml)
	}

	return
}
