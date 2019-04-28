package dao

import (
	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
)

const (
	_configTableName = "tv_price_config"
	_valid           = 0
	_invalid         = 1
)

// GetById select panel info by id
func (d *Dao) GetById(id int64) (panelInfo *model.TvPriceConfigResp, err error) {
	panelInfo = &model.TvPriceConfigResp{}

	if err = d.DB.Table(_configTableName).Where("id = ? and status in (?, ?)", id, _valid, _invalid).First(panelInfo).Error; err != nil {
		log.Error("GetById (%v) error(%v)", panelInfo, err)
	}

	return
}

// PanelStatus update panel status info by id
func (d *Dao) PanelStatus(id, status int64) (err error) {
	if err = d.DB.Table(_configTableName).Where("id = ? and status in (?, ?)", id, _valid, _invalid).Update("status", status).Error; err != nil {
		log.Error("PanelStatus (%v) error(%v)", id, err)
	}

	return
}

// SavePanel update or add panel info
func (d *Dao) SavePanel(panel *model.TvPriceConfig) (err error) {
	if err = d.DB.Save(panel).Error; err != nil {
		log.Error("SavePanel (%v) error(%v)", panel, err)
		return err
	}

	return
}

// ExistProduct check duplicated productId, productName
func (d *Dao) ExistProduct(productID string) (flag bool) {
	panel := []model.TvPriceConfig{}
	if err := d.DB.Table(_configTableName).Where("status in (?, ?) and product_id = ?", _valid, _invalid, productID).Find(&panel).Error; err != nil {
		log.Error("HasDiscount %v, Err %v", productID, err)
		return false
	}
	return !(len(panel) == 0)
}
