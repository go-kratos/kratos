package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// DelMachineNodeByMachineID delete node by machine id.
func (d *Dao) DelMachineNodeByMachineID(mID int64) error {
	return pkgerr.WithStack(d.db.Where("machine_id = ?", mID).Delete(&model.MachineNode{}).Error)
}

// GenMachineNodes generate some machine nodes.
func (d *Dao) GenMachineNodes(nodes []*model.MachineNode) (err error) {
	tx := d.db.Begin()
	if err = tx.Error; err != nil {
		return
	}
	for _, n := range nodes {
		if err = tx.Create(n).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return
}

// QueryMachineNodes query machine nodes by machine id.
func (d *Dao) QueryMachineNodes(mID int64) (nodes []*model.MachineNode, err error) {
	err = pkgerr.WithStack(d.db.Model(&model.MachineNode{}).Where("machine_id = ?", mID).Find(&nodes).Error)
	return
}
