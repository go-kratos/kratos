package dao

import (
	"database/sql"
	"fmt"
	"strconv"

	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

const (
	_machineRightJoinLogSQL = "SELECT m.id,m.name,ml.username,ml.operation_type,ml.operation_result,ml.ctime FROM machines AS m INNER JOIN machine_logs AS ml ON m.id = ml.machine_id"
	_machineRightCountSQL   = "SELECT count(m.id) FROM machines AS m INNER JOIN machine_logs AS ml ON m.id = ml.machine_id"
)

// InsertMachineLog insert machine log.
func (d *Dao) InsertMachineLog(machineLog *model.MachineLog) (err error) {
	return pkgerr.WithStack(d.db.Create(machineLog).Error)
}

// FindMachineLogsByMachineID find machine log by machineId.
func (d *Dao) FindMachineLogsByMachineID(machineID int64, pn, ps int) (total int64, machineLogs []*model.AboundMachineLog, err error) {
	var (
		qSQL = _machineRightJoinLogSQL
		cSQL = _machineRightCountSQL
		rows *sql.Rows
	)
	if machineID > 0 {
		qSQL = _machineRightJoinLogSQL + " WHERE machine_id = " + strconv.FormatInt(machineID, 10)
		cSQL = _machineRightCountSQL + " WHERE machine_id = " + strconv.FormatInt(machineID, 10)
	}

	cDB := d.db.Raw(cSQL)
	if err = pkgerr.WithStack(cDB.Count(&total).Error); err != nil {
		return
	}
	gDB := d.db.Raw(qSQL)
	if rows, err = gDB.Order("ml.ctime DESC").Offset((pn - 1) * ps).Limit(ps).Rows(); err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		ml := &model.AboundMachineLog{}
		if err = rows.Scan(&ml.MachineID, &ml.Name, &ml.Username, &ml.OperateType, &ml.OperateResult, &ml.OperateTime); err != nil {
			return
		}
		machineLogs = append(machineLogs, ml)
	}
	return
}

// FindMachineLogs Find Machine Logs.
func (d *Dao) FindMachineLogs(queryRequest *model.QueryMachineLogRequest) (total int64, machineLogs []*model.AboundMachineLog, err error) {
	var (
		qSQL = _machineRightJoinLogSQL
		cSQL = _machineRightCountSQL
		rows *sql.Rows
	)

	if queryRequest.MachineID > 0 || queryRequest.OperateType != "" || queryRequest.OperateUser != "" || queryRequest.MachineName != "" {
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

		if queryRequest.MachineName != "" {
			strSQL = fmt.Sprintf("%s %s  m.name like '%s'", strSQL, logicalWord, _wildcards+queryRequest.MachineName+_wildcards)
			logicalWord = _and
		}

		if queryRequest.OperateUser != "" {
			strSQL = fmt.Sprintf("%s %s  ml.username like '%s'", strSQL, logicalWord, _wildcards+queryRequest.OperateUser+_wildcards)
			logicalWord = _and
		}

		qSQL = _machineRightJoinLogSQL + " " + strSQL
		cSQL = _machineRightCountSQL + " " + strSQL

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
		ml := &model.AboundMachineLog{}
		if err = rows.Scan(&ml.MachineID, &ml.Name, &ml.Username, &ml.OperateType, &ml.OperateResult, &ml.OperateTime); err != nil {
			return
		}
		machineLogs = append(machineLogs, ml)
	}

	return
}
