package dao

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

const (
	_wildcards = "%"
	_where     = "WHERE"
	_and       = "AND"

	_pathAndPodNamesMappingSQL = "SELECT CONCAT(business_unit, '.', project, '.', app) AS \"path\" , GROUP_CONCAT(pod_name) AS \"pod_names\" FROM machines WHERE status >=0 AND status < 100 AND ctime < %s GROUP BY business_unit, project, app;"
	_nodeSQL                   = "SELECT DISTINCT(machine_id) FROM machine_nodes"
)

// FindExpiredMachineByDay find expired machines.
func (d *Dao) FindExpiredMachineByDay(day int) (ms []*model.Machine, err error) {
	err = pkgerr.WithStack(d.db.Where("status >= 0 AND DATEDIFF(end_time,NOW())>=0 AND DATEDIFF(end_time,NOW())<= ?", day).Find(&ms).Error)
	return
}

// FindExpiredMachine find a day-old machines.
func (d *Dao) FindExpiredMachine() (ms []*model.Machine, err error) {
	err = pkgerr.WithStack(d.db.Where("status >= 0 AND DATEDIFF(NOW(),end_time) = ?", -1).Find(&ms).Error)
	return
}

// QueryMachine query machine.
func (d *Dao) QueryMachine(id int64) (machine *model.Machine, err error) {
	machine = &model.Machine{}
	err = pkgerr.WithStack(d.db.Where("id = ?", id).First(machine).Error)
	return
}

// QueryMachineByName Query Machine By Name.
func (d *Dao) QueryMachineByName(machineName string) (machine *model.Machine, err error) {
	machine = &model.Machine{}
	err = pkgerr.WithStack(d.db.Where("name = ?", machineName).First(machine).Error)
	return
}

// QueryOnlineMachineByName Query Online Machine By Name.
func (d *Dao) QueryOnlineMachineByName(machineName string) (machine *model.Machine, err error) {
	machine = &model.Machine{}
	err = pkgerr.WithStack(d.db.Where("status >=100 and name = ?", machineName).First(machine).Error)
	return
}

// DelMachine delete machine.
func (d *Dao) DelMachine(id int64, username string) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("id = ?", id).Update("status", model.RemovedMachineInMerlin).Update("update_by", username).Error)
}

// HasMachine verify machine name.
func (d *Dao) HasMachine(name string) (b bool, err error) {
	var size int64
	if err = pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("name LIKE ?", name+"-"+_wildcards).Where("status <> ?", model.RemovedMachineInMerlin).Count(&size).Error); err != nil {
		return
	}
	b = size > 0
	return
}

// UpdateMachineStatus update machine status by giving value.
func (d *Dao) UpdateMachineStatus(id int64, status int) error {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("id = ?", id).Update("status", status).Error)
}

// QueryMachines query machines and update machine in names status to boot.
func (d *Dao) QueryMachines(names []string, qmr *model.QueryMachineRequest) (total int64, machines []*model.Machine, err error) {
	if err = pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("status >= 0 AND status < 100 AND pod_name IN (?)", names).Update("status", model.BootMachineInMerlin).Error); err != nil {
		return
	}
	gDB := d.db.Model(&model.Machine{}).Where("status >= 100 AND status < 300 AND pod_name IN (?) OR (status > -300 AND status <= -200 OR  status >= 0 AND status < 100) AND username = ?", names, qmr.Requester)
	if qmr.Username != "" {
		gDB = gDB.Where("username LIKE ?", qmr.Username+_wildcards)
	}
	if qmr.MachineName != "" {
		gDB = gDB.Where("BINARY name LIKE ?", _wildcards+qmr.MachineName+_wildcards)
	}
	if qmr.BusinessUnit != "" || qmr.Project != "" || qmr.App != "" {
		var (
			nSQL        = _nodeSQL
			logicalWord = _where
		)
		if qmr.BusinessUnit != "" {
			nSQL = fmt.Sprintf("%s %s BINARY business_unit = '%s'", nSQL, logicalWord, qmr.BusinessUnit)
			logicalWord = _and
		}
		if qmr.Project != "" {
			nSQL = fmt.Sprintf("%s %s BINARY project = '%s'", nSQL, logicalWord, qmr.Project)
			logicalWord = _and
		}
		if qmr.App != "" {
			nSQL = fmt.Sprintf("%s %s BINARY app = '%s'", nSQL, logicalWord, qmr.App)
		}
		gDB = gDB.Where("ID IN ( " + nSQL + " )")
	}
	if err = pkgerr.WithStack(gDB.Count(&total).Error); err != nil {
		return
	}
	err = pkgerr.WithStack(gDB.Order("ctime DESC").Offset((qmr.PageNum - 1) * qmr.PageSize).Limit(qmr.PageSize).Find(&machines).Error)
	return
}

// UpdateMachineEndTime update machine end time.
func (d *Dao) UpdateMachineEndTime(id int64, delayStatus int, endTime time.Time) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("id = ?", id).Update("delay_status", delayStatus).Update("end_time", endTime).Error)
}

// UpdateMachineDelayStatus update machine delay status.
func (d *Dao) UpdateMachineDelayStatus(id int64, delayStatus int) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("id = ?", id).Update("delay_status", delayStatus).Error)
}

// QueryPathAndPodNamesMapping query path and pod names map.
func (d *Dao) QueryPathAndPodNamesMapping() (pathAndPodNames map[string][]string, err error) {
	var (
		rows     *sql.Rows
		path     string
		podNames string
	)
	if rows, err = d.db.Raw(fmt.Sprintf(_pathAndPodNamesMappingSQL, d.c.Paas.MachineTimeout)).Rows(); err != nil {
		err = pkgerr.WithStack(err)
		return
	}
	defer rows.Close()
	pathAndPodNames = make(map[string][]string)
	for rows.Next() {
		if err = rows.Scan(&path, &podNames); err != nil {
			return
		}
		pathAndPodNames[path] = strings.Split(podNames, ",")
	}
	return
}

// UpdateMachineStatusByPodNames update machine status by pod names.
func (d *Dao) UpdateMachineStatusByPodNames(podNames []string, target int) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("pod_name IN (?)", podNames).Update("status", target).Error)
}

// UpdateMachineUser Update Machine User By Id.
func (d *Dao) UpdateMachineUser(machineID int64, username string) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("id = ?", machineID).Update("username", username).Error)
}

// QueryMachinesByPodNames query machines by giving pod name slices.
func (d *Dao) QueryMachinesByPodNames(podNames []string) (machines []*model.Machine, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("pod_name IN (?)", podNames).Find(&machines).Error)
	return
}

// InsertMachines insert machines.
func (d *Dao) InsertMachines(u string, gmr *model.GenMachinesRequest, ins []*model.CreateInstance) (err error) {
	var m *model.Machine
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}
	for _, in := range ins {
		m = in.ToMachine(u, gmr)
		if err = tx.Create(m).Error; err != nil {
			tx.Rollback()
			return
		}
		for _, n := range gmr.ToMachineNode(m.ID) {
			if err = tx.Create(n).Error; err != nil {
				tx.Rollback()
				return
			}
		}
		if err = tx.Create(m.ToMachineLog()).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// QueryMachineInRunning Query Machine In Running
func (d *Dao) QueryMachineInRunning() (ms []*model.Machine, err error) {
	err = d.db.Where("status >= 100").Find(&ms).Error
	return
}

// QueryMachineCount Query Machine Count
func (d *Dao) QueryMachineCount() (total int, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.Machine{}).Count(&total).Error)
	return
}
