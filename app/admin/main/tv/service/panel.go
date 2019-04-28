package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_configTableName = "tv_price_config"
	_valid           = 0
	_invalid         = 1
	_noPid           = 0
	_orderField      = "ctime"
	_vip             = 10
	_online          = 0
	_delete          = 2
)

// PanelInfo select panel info by id
func (s *Service) PanelInfo(id int64) (panelInfo *model.TvPriceConfigResp, err error) {
	if panelInfo, err = s.dao.GetById(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("PanelInfo (%v) error(%v)", panelInfo, err)
		return
	}

	if panelInfo != nil {
		panelInfo.OriginPrice = panelInfo.Price
		hasDiscount, price, discountInfos := s.hasDiscount(nil, panelInfo.ID)
		if hasDiscount {
			panelInfo.Price = price
		}
		panelInfo.Items = discountInfos
	}
	return
}

// PanelStatus change panel status by id
func (s *Service) PanelStatus(id, status int64) (err error) {
	var (
		flag      bool
		panelInfo *model.TvPriceConfigResp
	)
	if panelInfo, err = s.dao.GetById(id); err != nil {
		log.Error("PanelInfo (%v) error(%v)", panelInfo, err)
	}
	if status == _online && panelInfo.SuitType == _vip {
		if flag, err = s.hasOnlineUpgradeVipProduct(); err != nil {
			log.Error("GetValidUpgradeVipProduct Err %v", err)
			return
		}
		if flag {
			err = ecode.TVVipSuitTypeConflict
			return err
		}
	}

	if panelInfo.PID == 0 && status == _delete {
		if err = s.dao.DB.Table(_configTableName).Where("pid = ?", id).Update("status", _delete).Error; err != nil {
			log.Error("PanelStatus discount (%v) error(%v)", id, err)
		}
	}

	if err = s.dao.PanelStatus(id, status); err != nil {
		log.Error("PanelStatus (%v) error(%v)", id, err)
	}

	return
}

// SavePanel add or update panel info
func (s *Service) SavePanel(c context.Context, panel *model.TvPriceConfig) (err error) {
	opType := s.c.YSTParam.Update
	//start tx
	tx := s.DB.Begin()

	if panel.PID != 0 {
		_, _, discounts := s.hasDiscount(nil, panel.PID)
		if flag := checkDisCountTime(discounts, panel); !flag {
			err = ecode.TvPriceTimeConflict
			return err
		}
		s.copyParentExtraField(panel)
	}

	if panel.ID != 0 && panel.PID == 0 {
		if err = tx.Table(_configTableName).Where("pid = ?", panel.ID).
			Update(map[string]interface{}{
				"suit_type":   panel.SuitType,
				"sub_type":    panel.SubType,
				"selected":    panel.Selected,
				"superscript": panel.Superscript,
				"month":       panel.Month,
			}).Error; err != nil {
			log.Error("Update discount failed while update panel, Err %v", err)
			tx.Rollback()
			return
		}
		_, _, discounts := s.hasDiscount(tx, panel.ID)
		for _, discount := range discounts {
			if err = s.syncPanels(c, opType, &discount); err != nil {
				err = ecode.TvVipProdSyncErr
				tx.Rollback()
				return err
			}
		}
	}

	if panel.ID == 0 {
		opType = s.c.YSTParam.Insert

		if panel.PID == 0 {
			panel.Status = _invalid
		} else {
			s.copyParentExtraField(panel)
		}
		if flag := s.dao.ExistProduct(panel.ProductID); flag {
			err = ecode.TvVipProductExit
			return err
		}
	}

	if err = tx.Save(panel).Error; err != nil {
		log.Error("SavePanel %s, Err %v", panel, err)
		return err
	}
	if err = s.syncPanels(c, opType, panel); err != nil {
		err = ecode.TvVipProdSyncErr
		tx.Rollback()
		return err
	}
	tx.Commit()

	return
}

// PanelList get panle list
func (s *Service) PanelList(platform, month, subType, suitType int64) (panels []*model.TvPriceConfigListResp, err error) {
	var (
		db = s.dao.DB.Model(&model.TvPriceConfigListResp{}).Where("pid = ? and status in (?, ?)", _noPid, _valid, _invalid)
	)

	if platform != 0 {
		db = db.Where("platform = ?", platform)
	}
	if month != 0 {
		db = db.Where("month = ?", month)
	}
	if subType != -1 {
		db = db.Where("sub_type = ?", subType)
	}
	if suitType != -1 {
		db = db.Where("suit_type = ?", suitType)
	}

	if err = db.Order("suit_type, sub_type desc, month desc").Find(&panels).Error; err != nil {
		log.Error("OrderList %v, Err %v", panels, err)
		return panels, err
	}

	for _, panel := range panels {
		hasDiscount, price, _ := s.hasDiscount(nil, panel.ID)
		panel.OriginPrice = panel.Price

		if hasDiscount {
			panel.Price = price
		}
	}

	return
}

// HasDiscount Judge whether there is a discount
func (s *Service) hasDiscount(tx *gorm.DB, id int64) (hasDiscount bool, price int64, panels []model.TvPriceConfig) {
	if tx == nil {
		tx = s.dao.DB
	}
	if err := tx.Table(_configTableName).Where("pid = ? and status = ?", id, _valid).Find(&panels).Error; err != nil {
		log.Error("HasDiscount %v, Err %v", id, err)
		return
	}
	nowTime := time.Now().Unix()

	for _, panel := range panels {
		if int64(panel.Stime) < nowTime && nowTime < int64(panel.Etime) {
			hasDiscount = true
			price = panel.Price
			return
		}
	}

	return
}

// checkDisCountTime check discount time conflict
func checkDisCountTime(discounts []model.TvPriceConfig, panel *model.TvPriceConfig) (flag bool) {
	var (
		startTime = panel.Stime
		endTime   = panel.Etime
	)

	for _, discount := range discounts {
		// do not compare with self
		if panel.ID != discount.ID {
			if discount.Stime < startTime && startTime < discount.Etime {
				return false
			}

			if discount.Stime < endTime && endTime < discount.Etime {
				return false
			}

			if discount.Stime > startTime && endTime > discount.Etime {
				return false
			}
		}
	}

	return true
}

// checkRemotePanel check YST panel
func (s *Service) checkRemotePanel(c context.Context) {
	var (
		panels []*model.TvPriceConfigListResp
	)
	if err := s.dao.DB.Table(_configTableName).Order(_orderField).Find(&panels).Error; err != nil {
		log.Error("CheckRemotePanel Err %v", err)
		return
	}
	res, _ := s.getRemotePanels(c)
	remotePanels := res.Product

	panelMap := make(map[string]*model.TvPriceConfigListResp, len(panels))
	remotePaneMap := make(map[string]model.Product, len(remotePanels))
	for i := 0; i < len(panels); i++ {
		panelMap[panels[i].ProductID] = panels[i]
	}
	for i := 0; i < len(remotePanels); i++ {
		remotePaneMap[remotePanels[i].ID] = remotePanels[i]
	}

	for i := 0; i < len(panels); i++ {
		rp, exists := remotePaneMap[panels[i].ProductID]
		if exists {
			s.compareFiled(rp, panels[i])
		} else {
			log.Error("Our panel not exists in YST, panel id is (%v)", panels[i].ProductID)
		}
	}

	for i := 0; i < len(remotePanels); i++ {
		p, exists := panelMap[remotePanels[i].ID]
		if exists {
			s.compareFiled(remotePanels[i], p)
		} else {
			log.Error("YST panel not exists in our db, panel id is (%v)", remotePanels[i].ID)
		}
	}
}

// getRemotePanels get YST panel
func (s *Service) getRemotePanels(c context.Context) (res *model.RemotePanel, err error) {
	var (
		req *http.Request
	)
	res = &model.RemotePanel{}
	params := map[string]string{
		"vod_type": s.c.YSTParam.QueryPanelType,
		"source":   s.c.YSTParam.Source,
	}
	reqBody, _ := json.Marshal(params)
	getRemotePanelUrl := s.c.URLConf.GetRemotePanelUrl
	if req, err = http.NewRequest(http.MethodPost, getRemotePanelUrl, bytes.NewReader(reqBody)); err != nil {
		log.Error("MerakNotify NewRequest Err %v, Url %v", err, getRemotePanelUrl)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err = s.client.Do(c, req, &res); err != nil {
		log.Error("MergeUpInfo http req failed ,err:%v", err)
		return
	}
	return
}

// compareFiled compare our panel field with YST
func (s *Service) compareFiled(remotePanel model.Product, panel *model.TvPriceConfigListResp) {
	rpid := remotePanel.ID
	prodiuctId := panel.ProductID
	if rpid != prodiuctId {
		log.Error("PanelInfo id different, Remote panel prodiuctId is (%v), Our panel prodiuctId is (%v)", rpid, prodiuctId)
	}
	if remotePanel.Price != panel.Price {
		log.Error("PanelInfo price different, Remote panel prodiuctId is (%v), Our panel prodiuctId is (%v)", rpid, prodiuctId)
	}
	if remotePanel.Contract != strconv.Itoa(int(panel.SubType)) {
		log.Error("PanelInfo subType different, Remote panel prodiuctId is (%v), Our panel prodiuctId is (%v)", rpid, prodiuctId)
	}
	if remotePanel.SuitType != panel.SuitType {
		log.Error("PanelInfo suitType different, Remote panel id is (%v), Our panel id is (%v)", rpid, prodiuctId)
	}
	if remotePanel.ProductDuration != strconv.FormatInt(panel.Month*31, 10) {
		log.Error("PanelInfo months different, Remote panel prodiuctId is (%v), Our panel prodiuctId is (%v)", rpid, prodiuctId)
	}
	if remotePanel.Title != panel.ProductName {
		log.Error("PanelInfo productName different, Remote panel prodiuctId is (%v), Our panel id prodiuctId (%v)", rpid, prodiuctId)
	}
	if panel.PID != 0 {
		parentPanel, err := s.dao.GetById(panel.PID)
		if err != nil {
			log.Error("PanelInfo (%v) error(%v)", parentPanel, err)
		}
		if remotePanel.ComboPkgID != parentPanel.ProductID {
			log.Error("PanelInfo pid different, Remote panel prodiuctId is (%v), Our panel prodiuctId is (%v)", rpid, prodiuctId)
		}
	}
}

// syncPanels send our panel to YST
func (s *Service) syncPanels(c context.Context, opType string, panel *model.TvPriceConfig) (err error) {
	var (
		req *http.Request
		res struct {
			Result  string `json:"result"`
			Message string `json:"message"`
		}
		params struct {
			OpType  string        `json:"optype"`
			Source  string        `json:"source"`
			Product model.Product `json:"product"`
		}
	)
	params.OpType = opType
	params.Source = s.c.YSTParam.Source
	params.Product.ID = panel.ProductID
	params.Product.VodType = s.c.YSTParam.InsertPanelType
	params.Product.Title = panel.ProductName
	params.Product.ProductDuration = strconv.FormatInt(31*panel.Month, 10)
	params.Product.Description = panel.Remark
	params.Product.Contract = strconv.Itoa(int(panel.SubType))
	if panel.PID != 0 {
		parentPanel, _ := s.dao.GetById(panel.PID)
		if parentPanel != nil {
			params.Product.ComboPkgID = parentPanel.ProductID
			params.Product.ComboDes = parentPanel.Remark
		} else {
			log.Error("Prarent PanelInfo not found, pid =(%v)", panel.PID)
			err = ecode.NothingFound
			return err
		}

	}
	params.Product.Price = panel.Price
	params.Product.SuitType = panel.SuitType

	reqBody, _ := json.Marshal(params)

	SyncPanelUrl := s.c.URLConf.SyncPanelUrl
	if req, err = http.NewRequest(http.MethodPost, SyncPanelUrl, bytes.NewReader(reqBody)); err != nil {
		log.Error("MerakNotify NewRequest Err %v, Url %v", err, SyncPanelUrl)
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if err = s.client.Do(c, req, &res); err != nil {
		log.Error("MergeUpInfo http req failed ,err:%v", err)
		return err
	}

	if res.Result != "SUCCESS" {
		err = ecode.TvVipProdSyncErr
		log.Info("Sync panel To YST Fail,err:%v", res.Message)
		return err
	}

	return
}

func (s *Service) hasOnlineUpgradeVipProduct() (flag bool, err error) {
	var (
		panels []*model.TvPriceConfig
	)
	if err = s.dao.DB.Table(_configTableName).Where("suit_type = ? and status = ? and pid = ?", _vip, _valid, _noPid).Find(&panels).Error; err != nil {
		log.Error("GetValidUpgradeVipProduct Err %v", err)
		return
	}
	flag = !(len(panels) == 0)

	return
}

func (s *Service) copyParentExtraField(panel *model.TvPriceConfig) {
	parentPanel, _ := s.dao.GetById(panel.PID)
	if parentPanel != nil {
		panel.SuitType = parentPanel.SuitType
		panel.SubType = parentPanel.SubType
		panel.Selected = parentPanel.Selected
		panel.Superscript = parentPanel.Superscript
		panel.Month = parentPanel.Month
	}
}
