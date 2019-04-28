package dao

import (
	"go-common/app/admin/ep/melloi/model"

	pkgerr "github.com/pkg/errors"
)

// QueryGRPC Query grpc
func (d *Dao) QueryGRPC(grpc *model.GRPC, pn, ps int32, treeNodes []string) (qgr *model.QueryGRPCResponse, err error) {
	qgr = &model.QueryGRPCResponse{}
	grpc.Active = 1

	gDB := d.DB.Table(model.GRPC{}.TableName()).Where("app in (?)", treeNodes)
	if grpc.Department != "" && grpc.Project != "" && grpc.APP != "" {
		gDB = gDB.Where("department = ? and project = ? and app = ?", grpc.Department, grpc.Project, grpc.APP)
	} else if grpc.Department != "" && grpc.Project != "" {
		gDB = gDB.Where("department = ? and project = ?", grpc.Department, grpc.Project)
	} else if grpc.Department != "" {
		gDB = gDB.Where("department = ?", grpc.Department)
	}

	err = gDB.Where(grpc).Count(&qgr.TotalSize).Offset((pn - 1) * ps).
		Limit(ps).Order("id desc").Find(&qgr.GRPCS).Error
	qgr.PageNum = pn
	qgr.PageSize = ps
	return
}

// QueryGRPCByWhiteName Query grpc By WhiteName
func (d *Dao) QueryGRPCByWhiteName(grpc *model.GRPC, pn, ps int32) (qgr *model.QueryGRPCResponse, err error) {
	qgr = &model.QueryGRPCResponse{}
	grpc.Active = 1
	err = d.DB.Table(model.GRPC{}.TableName()).Where(grpc).Count(&qgr.TotalSize).Offset((pn - 1) * ps).
		Limit(ps).Order("id desc").Find(&qgr.GRPCS).Error
	qgr.PageNum = pn
	qgr.PageSize = ps
	return
}

// QueryGRPCByID Query GRPC By ID
func (d *Dao) QueryGRPCByID(id int) (grpc *model.GRPC, err error) {
	grpc = &model.GRPC{}
	err = pkgerr.WithStack(d.DB.Where("id = ?", id).First(grpc).Error)
	return
}

// UpdateGRPC Update grpc
func (d *Dao) UpdateGRPC(grpc *model.GRPC) error {
	return d.DB.Table(model.GRPC{}.TableName()).Where("id=?", grpc.ID).Update(grpc).Error
}

//CreateGRPC  new grpc
func (d *Dao) CreateGRPC(grpc *model.GRPC) (g *model.GRPC, err error) {
	grpc.Active = 1
	return grpc, pkgerr.WithStack(d.DB.Create(grpc).Error)
}

// DeleteGRPC Delete grpc
func (d *Dao) DeleteGRPC(id int) error {
	return d.DB.Table(model.GRPC{}.TableName()).Where("id=?", id).Update("active", -1).Error
}
