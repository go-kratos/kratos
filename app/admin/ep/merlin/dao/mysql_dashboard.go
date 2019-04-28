package dao

import (
	"database/sql"
	"fmt"

	"go-common/app/admin/ep/merlin/model"
)

// MachineLifeCycle Statistics Machine Life Cycle.
func (d *Dao) MachineLifeCycle() (machineLifeCycles []*model.MachineLifeCycle, err error) {
	var rows *sql.Rows

	SQL := "select CASE " +
		"when TIMESTAMPDIFF(DAY,ctime,end_time)>0 and TIMESTAMPDIFF(DAY,ctime,end_time) <=100 THEN'0-100' " +
		"when TIMESTAMPDIFF(DAY,ctime,end_time)>100 and TIMESTAMPDIFF(DAY,ctime,end_time) <=200 THEN'100-200' " +
		"when TIMESTAMPDIFF(DAY,ctime,end_time)>200 THEN'200+' ELSE '0-100' " +
		"End as lifecycle,count(-1)from machines GROUP BY lifecycle"

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mlc := &model.MachineLifeCycle{}
		if err = rows.Scan(&mlc.Duration, &mlc.Count); err != nil {
			return
		}
		machineLifeCycles = append(machineLifeCycles, mlc)
	}
	return

}

// MachineCountGroupByBusiness Machine Count Group By Business
func (d *Dao) MachineCountGroupByBusiness() (machinesCount []*model.MachineCountGroupByBusiness, err error) {
	var rows *sql.Rows

	SQL := "SELECT business_unit,count(-1) as count FROM machines GROUP BY business_unit ORDER BY count DESC"

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MachineCountGroupByBusiness{}
		if err = rows.Scan(&mc.BusinessUnit, &mc.Count); err != nil {
			return
		}
		machinesCount = append(machinesCount, mc)
	}
	return
}

// MachineCountGroupByBusinessInRunning Machine Count Group By Business In Running
func (d *Dao) MachineCountGroupByBusinessInRunning() (machinesCount []*model.MachineCountGroupByBusiness, err error) {
	var rows *sql.Rows

	SQL := "SELECT business_unit,count(-1) as count FROM machines where status>=100 GROUP BY business_unit ORDER BY count DESC"

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MachineCountGroupByBusiness{}
		if err = rows.Scan(&mc.BusinessUnit, &mc.Count); err != nil {
			return
		}
		machinesCount = append(machinesCount, mc)
	}
	return
}

// MachineWillBeExpired Statistics Machine Will Be Expired
func (d *Dao) MachineWillBeExpired() (machines []*model.MachineCreatedAndEndTime, err error) {
	var rows *sql.Rows

	SQL := "SELECT id,name,username,app,ctime,end_time from machines where `status`>=100 order by end_time limit 10"

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MachineCreatedAndEndTime{}
		if err = rows.Scan(&mc.ID, &mc.MachineName, &mc.Username, &mc.App, &mc.CreateTime, &mc.EndTime); err != nil {
			return
		}
		machines = append(machines, mc)
	}
	return
}

// MachineLatestCreated Machine Latest Created
func (d *Dao) MachineLatestCreated() (machines []*model.MachineCreatedAndEndTime, err error) {
	var rows *sql.Rows

	SQL := "SELECT id,name,username,app,ctime,end_time from machines where `status`>=100 order by ctime desc limit 10"

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MachineCreatedAndEndTime{}
		if err = rows.Scan(&mc.ID, &mc.MachineName, &mc.Username, &mc.App, &mc.CreateTime, &mc.EndTime); err != nil {
			return
		}
		machines = append(machines, mc)
	}
	return
}

// MobileMachineUserUsageCount Mobile Machine User Usage Count
func (d *Dao) MobileMachineUserUsageCount() (mobileMachinesUsageCount []*model.MobileMachineUserUsageCount, err error) {
	var rows *sql.Rows

	SQL := fmt.Sprintf("select username,count(username) as count from mobile_machine_logs "+
		"where operation_type='%s' and operation_result ='%s' "+
		"group by username order by count desc", model.MBBindLog, model.OperationSuccessForMachineLog)

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		mc := &model.MobileMachineUserUsageCount{}
		if err = rows.Scan(&mc.Username, &mc.Count); err != nil {
			return
		}
		mobileMachinesUsageCount = append(mobileMachinesUsageCount, mc)
	}

	return
}

// MobileMachineUserLendCount Mobile Machine Lend Count
func (d *Dao) MobileMachineUserLendCount() (mobileMachinesUsageCount []*model.MobileMachineUserUsageCount, err error) {
	var rows *sql.Rows

	SQL := fmt.Sprintf("select username,count(username) as count from mobile_machine_logs "+
		"where operation_type='%s' and operation_result ='%s' "+
		"group by username order by count desc", model.MBLendOutLog, model.OperationSuccessForMachineLog)

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MobileMachineUserUsageCount{}
		if err = rows.Scan(&mc.Username, &mc.Count); err != nil {
			return
		}
		mobileMachinesUsageCount = append(mobileMachinesUsageCount, mc)
	}
	return
}

// MobileMachineUsageCount Mobile Machine Usage Count
func (d *Dao) MobileMachineUsageCount() (mobileMachinesUsageCount []*model.MobileMachineUsageCount, err error) {
	var rows *sql.Rows

	SQL := fmt.Sprintf("select b.id,b.name,count(b.name) as count from mobile_machine_logs as a "+
		"left join mobile_machines as b on a.machine_id = b.id "+
		"where a.operation_type='%s' and a.operation_result = '%s' "+
		"group by b.id,b.`name` order by count desc", model.MBBindLog, model.OperationSuccessForMachineLog)

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MobileMachineUsageCount{}
		if err = rows.Scan(&mc.MobileMachineID, &mc.MobileMachineName, &mc.Count); err != nil {
			return
		}
		mobileMachinesUsageCount = append(mobileMachinesUsageCount, mc)
	}
	return
}

// MobileMachineLendCount Mobile Machine Lend Count
func (d *Dao) MobileMachineLendCount() (mobileMachinesUsageCount []*model.MobileMachineUsageCount, err error) {
	var rows *sql.Rows

	SQL := fmt.Sprintf("select b.id,b.name,count(b.name) as count from mobile_machine_logs as a "+
		"left join mobile_machines as b on a.machine_id = b.id "+
		"where a.operation_type='%s' and a.operation_result = '%s' "+
		"group by b.id,b.`name` order by count desc", model.MBLendOutLog, model.OperationSuccessForMachineLog)

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MobileMachineUsageCount{}
		if err = rows.Scan(&mc.MobileMachineID, &mc.MobileMachineName, &mc.Count); err != nil {
			return
		}
		mobileMachinesUsageCount = append(mobileMachinesUsageCount, mc)
	}
	return
}

// MobileMachineModeCount Mobile Machine Mode Count.
func (d *Dao) MobileMachineModeCount() (mobileMachinesTypeCount []*model.MobileMachineTypeCount, err error) {
	var rows *sql.Rows

	SQL := "select mode,count(mode) as count from mobile_machines where action > 0 group by mode order by count desc"

	if rows, err = d.db.Raw(SQL).Rows(); err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		mc := &model.MobileMachineTypeCount{}
		if err = rows.Scan(&mc.ModeName, &mc.Count); err != nil {
			return
		}
		mobileMachinesTypeCount = append(mobileMachinesTypeCount, mc)
	}
	return
}

// MobileMachineUseRecord Mobile Machine Use Record
func (d *Dao) MobileMachineUseRecord() (mobileMachineLogs []*model.MobileMachineLog, err error) {
	err = d.db.Where("operation_type= ? or operation_type = ?", model.MBBindLog, model.MBReleaseLog).Order("machine_id,id").Find(&mobileMachineLogs).Error
	return
}
