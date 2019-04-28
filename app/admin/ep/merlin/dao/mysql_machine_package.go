package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// FindAllMachinePackages find all machine packages.
func (d *Dao) FindAllMachinePackages() (machinePackages []*model.MachinePackage, err error) {
	err = pkgerr.WithStack(d.db.Find(&machinePackages).Error)
	return
}
