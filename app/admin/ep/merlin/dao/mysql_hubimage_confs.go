package dao

import (
	"go-common/app/admin/ep/merlin/model"
	"go-common/library/ecode"

	pkgerr "github.com/pkg/errors"
)

// InsertHubImageConf Insert Hub Image conf.
func (d *Dao) InsertHubImageConf(hubImageConf *model.HubImageConf) (err error) {
	return pkgerr.WithStack(d.db.Create(hubImageConf).Error)
}

// UpdateHubImageConf Update Hub Image Conf.
func (d *Dao) UpdateHubImageConf(hubImageConf *model.HubImageConf) (err error) {
	return pkgerr.WithStack(d.db.Model(&model.HubImageConf{}).Where("image_name=?", hubImageConf.ImageName).Updates(map[string]interface{}{
		"update_by": hubImageConf.UpdateBy, "command": hubImageConf.Command, "environments": hubImageConf.Envs, "hosts": hubImageConf.Hosts}).Error)
}

// FindHubImageConfByImageName Find Hub Image Conf By Image Name.
func (d *Dao) FindHubImageConfByImageName(imageName string) (hubImageConf *model.HubImageConf, err error) {
	hubImageConf = &model.HubImageConf{}
	if err = d.db.Where("image_name = ?", imageName).First(hubImageConf).Error; err == ecode.NothingFound {
		err = nil
	}
	return
}
