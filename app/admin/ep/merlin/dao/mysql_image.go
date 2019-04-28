package dao

import (
	"go-common/app/admin/ep/merlin/model"

	pkgerr "github.com/pkg/errors"
)

// Images Search all images in db.
func (d *Dao) Images() (images []*model.Image, err error) {
	err = pkgerr.WithStack(d.db.Where("status = ?", model.AliveImageStatus).Find(&images).Error)
	return
}

// AddImage Create new image in db.
func (d *Dao) AddImage(image *model.Image) error {
	return pkgerr.WithStack(d.db.Create(image).Error)
}

// UpdateImage Update image in db.
func (d *Dao) UpdateImage(image *model.Image) error {
	return pkgerr.WithStack(d.db.Model(&model.Image{}).Updates(image).Error)
}

// DelImage Delete image in db.
func (d *Dao) DelImage(iID int64) error {
	return pkgerr.WithStack(d.db.Model(&model.Image{}).Where("id = ?", iID).Update("status", model.DeletedImageStatus).Error)
}
