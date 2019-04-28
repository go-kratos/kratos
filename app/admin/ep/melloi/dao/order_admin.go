package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

// QueryOrderAdmin get administrator for Order of this project
func (d *Dao) QueryOrderAdmin(userName string) (admin *model.OrderAdmin, err error) {
	admin = &model.OrderAdmin{}
	err = d.DB.Table(model.OrderAdmin{}.TableName()).Where("user_name = ?", userName).First(admin).Error
	return
}

// AddOrderAdmin add order of admin
func (d *Dao) AddOrderAdmin(admin *model.OrderAdmin) error {
	return d.DB.Create(admin).Error
}
