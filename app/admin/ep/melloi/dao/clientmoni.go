package dao

import "go-common/app/admin/ep/melloi/model"

//AddClientMoni add ClientMoni
func (d *Dao) AddClientMoni(clientMoni *model.ClientMoni) (clientMoniID int, err error) {
	err = d.DB.Create(clientMoni).Error
	clientMoniID = clientMoni.ID
	return
}

//QueryClientMoni query ClientMoni
func (d *Dao) QueryClientMoni(clientMoni *model.ClientMoni) (clm []*model.ClientMoni, err error) {
	err = d.DB.Table(model.ClientMoni{}.TableName()).Where(clientMoni).Find(&clm).Error
	return
}
