package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// OrderPayResult order pay result.
func (s *Service) OrderPayResult(c context.Context, orderNo string, mid, appID int64, platform, device, mobiApp string, build int64, panelType string) (res *model.OrderResult, err error) {
	var (
		dialog *model.ConfDialog
		order  *model.OrderInfo
		vip    *model.VipInfo
	)
	res = &model.OrderResult{OrderNo: orderNo}
	// check order status
	if order, err = s.OrderInfo(c, orderNo); err != nil {
		return
	}
	if order == nil {
		err = ecode.VipOrderNoErr
		return
	}
	if order.Mid != mid {
		err = ecode.AccessDenied
		return
	}
	res.Status = order.Status
	if order.Status != model.SUCCESS {
		return
	}
	// get available dialog
	if dialog = s.getDialog(c, appID, platform, device, mobiApp, build, panelType); dialog == nil {
		return
	}
	if vip, err = s.VipInfo(c, mid); err != nil {
		return
	}
	// format and return
	if strings.Contains(dialog.Content, "[buy_months]") {
		dialog.Content = strings.Replace(dialog.Content, "[buy_months]", strconv.FormatInt(int64(order.BuyMonths), 10), -1)
	}
	if strings.Contains(dialog.Content, "[overdue_time]") {
		dialog.Content = strings.Replace(dialog.Content, "[overdue_time]", vip.VipOverdueTime.Time().Format("2006-01-02"), -1)
	}
	res.Dialog = dialog
	return
}

func (s *Service) getDialog(c context.Context, appID int64, platform, device, mobiApp string, build int64, panelType string) (res *model.ConfDialog) {
	var platAppMap map[int64]*model.ConfDialog
	vplat := s.GetPlatID(c, platform, panelType, mobiApp, device, build)
	log.Info("s.GetPlatID(%s, %s, %s, %s) plat(%d)", platform, panelType, mobiApp, device, vplat)
	if platAppMapI, ok := s.dialogMap.Load(vplat); ok {
		platAppMap = platAppMapI.(interface{}).(map[int64]*model.ConfDialog)
		if res, ok = platAppMap[appID]; !ok {
			res = platAppMap[0]
		}
	}
	if res != nil {
		return
	}
	var (
		appMap  map[int64]*model.ConfDialog
		allPlat int64
	)
	//PRD: 限定平台+限定app_id > 限定平台+未限定app_id > 未限定平台+限定app_id > 未限定平台+未限定app_id
	if appMapI, ok := s.dialogMap.Load(allPlat); ok {
		appMap = appMapI.(interface{}).(map[int64]*model.ConfDialog)
		if res, ok = appMap[appID]; !ok {
			res = appMap[0]
		}
	}
	return
}

func (s *Service) loaddialogproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("service.loaddialogproc panic(%v)", x)
			go s.loaddialogproc()
			log.Info("service.loaddialogproc recover")
		}
	}()
	for {
		time.Sleep(time.Second)
		if err := s.loadDialog(); err != nil {
			log.Error("s.loadDialog() err(%+v)", err)
		}
	}
}

func (s *Service) loadDialog() (err error) {
	var (
		res  map[int64]map[int64]*model.ConfDialog
		ds   []*model.ConfDialog
		c    = context.Background()
		curr = time.Now()
	)
	if ds, err = s.dao.DialogAll(c); err != nil {
		return
	}
	log.Info("s.dao.DialogAll() len(%d)", len(ds))
	// [plat][appID]*model.ConfDialog
	res = make(map[int64]map[int64]*model.ConfDialog)
	for _, dialog := range ds {
		if dialog.StartTime.Time().After(curr) {
			log.Warn("loadDialog (%+v) startTime not available.", dialog)
			continue
		}
		appMap, ok := res[dialog.Platform]
		if !ok {
			appMap = make(map[int64]*model.ConfDialog)
			res[dialog.Platform] = appMap
		}
		td, ok := appMap[dialog.AppID]
		if ok {
			if td.StartTime > dialog.StartTime {
				continue
			}
		}
		appMap[dialog.AppID] = dialog
	}
	if err != nil {
		log.Error("s.loadDialog err(%+v)", err)
		return
	}
	log.Info("loadDialog success %+v", res)
	for k, v := range res {
		s.dialogMap.Store(k, v)
	}
	return
}
