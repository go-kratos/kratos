package dao

import (
	"go-common/app/admin/ep/melloi/model"

	pkgerr "github.com/pkg/errors"
)

// QueryGRPCSnapByID  query grpcsnap by id
func (d *Dao) QueryGRPCSnapByID(id int) (grpcSnap *model.GRPCSnap, err error) {
	grpcSnap = &model.GRPCSnap{}
	err = pkgerr.WithStack(d.DB.Table(model.GRPCSnap{}.TableName()).Where("id = ?", id).First(grpcSnap).Error)
	return
}

// UpdateGRPCSnap Update grpc
func (d *Dao) UpdateGRPCSnap(grpcSnap *model.GRPCSnap) error {
	return d.DB.Table(model.GRPCSnap{}.TableName()).Where("id=?", grpcSnap.ID).Update(grpcSnap).Error
}

// CreateGRPCSnap CreateGRPC  new grpc
func (d *Dao) CreateGRPCSnap(grpcSnap *model.GRPCSnap) (err error) {
	grpcSnap.Active = 1
	return pkgerr.WithStack(d.DB.Table(model.GRPCSnap{}.TableName()).Create(grpcSnap).Error)
}

// DeleteGRPCSnap Delete grpc snap
func (d *Dao) DeleteGRPCSnap(id int) error {
	return pkgerr.WithStack(d.DB.Table(model.GRPCSnap{}.TableName()).Where("id=?", id).Update("active", -1).Error)
}
