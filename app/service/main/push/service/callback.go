package service

import (
	"context"

	"go-common/app/service/main/push/dao/huawei"
	"go-common/app/service/main/push/model"
)

// AddCallback adds callback.
func (s *Service) AddCallback(c context.Context, cb *model.Callback) (err error) {
	// db 中能取到优先从 db 取（客户端有时候在送达和点击时刻取 mid 或 buvid 时会有问题）
	if r, e := s.dao.Report(c, cb.Token); e == nil && r != nil {
		cb.Mid = r.Mid
		cb.Buvid = r.Buvid
		cb.APP = r.APPID
		cb.Platform = r.PlatformID
		cb.Brand = model.DeviceBrand(r.DeviceBrand)
	}
	if cb.Platform == model.PlatformHuawei && cb.Extra != nil &&
		(cb.Extra.Status == huawei.CallbackTokenUninstalled || cb.Extra.Status == huawei.CallbackTokenNotApply || cb.Extra.Status == huawei.CallbackTokenInactive) {
		s.reportCache.Save(func() {
			if cb.Mid > 0 {
				s.DelReport(context.TODO(), cb.APP, cb.Mid, cb.Token)
			} else {
				s.dao.DelReport(context.TODO(), cb.Token)
			}
		})
	}
	err = s.dao.AddCallback(c, cb)
	return
}
