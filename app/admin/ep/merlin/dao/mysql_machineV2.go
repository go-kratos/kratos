package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// InsertMachinesV2 insert machines for v2.
func (d *Dao) InsertMachinesV2(u string, gmr *model.GenMachinesRequest, pgmr *model.PaasGenMachineRequest) (ms []*model.Machine, err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}
	for _, pm := range pgmr.Machines {
		m := pm.ToMachine(u, gmr)
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
		ms = append(ms, m)
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// UpdateStatusForMachines update status for machines.
func (d *Dao) UpdateStatusForMachines(status int, ids []int64) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("id IN (?)", ids).Update("status", status).Error)
}

// UpdateMachineStatusByName update status by name.
func (d *Dao) UpdateMachineStatusByName(status int, n string) error {
	return pkgerr.WithStack(d.db.Model(&model.Machine{}).Where("name = ?", n).Update("status", status).Error)
}
