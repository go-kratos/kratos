package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// AddOrUpCallback will add or update a callback
func (s *Service) AddOrUpCallback(c context.Context, cbp *param.AddCallbackParam) (cbID int32, err error) {
	cb := &model.Callback{}

	if err = s.dao.ORM.Model(cb).
		Where(&model.Callback{Business: cbp.Business}).
		Assign(&model.Callback{
			Business:    cbp.Business,
			URL:         cbp.URL,
			IsSobot:     cbp.IsSobot,
			State:       cbp.State,
			ExternalAPI: cbp.ExternalAPI,
			SourceAPI:   cbp.SourceAPI,
		}).
		FirstOrCreate(cb).Error; err != nil {
		log.Error("Failed to create a callback(%+v): %v", cb, err)
		return
	}
	cbID = cb.CbID

	// load callbacks from datebase in cache
	go s.loadCallbacks()

	return
}

// ListCallback will list all enabled callbacks
func (s *Service) ListCallback(c context.Context) (cbList model.CallbackSlice, err error) {
	for _, cb := range s.callbackCache {
		cbList = append(cbList, cb)
	}

	sort.Slice(cbList, func(i, j int) bool {
		return cbList[i].Business > cbList[j].Business
	})
	return
}

// SendCallbackRetry will try to send callback with specified attempts
func (s *Service) SendCallbackRetry(c context.Context, cb *model.Callback, payload *model.Payload) (err error) {
	attempts := 3
	for i := 0; i < attempts; i++ {
		err = s.dao.SendCallback(c, cb, payload)
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(time.Second * 1)
	}
	return errors.Wrapf(err, "after %d attempts", attempts)
}

// SetExtAPI .
func (s *Service) SetExtAPI(ctx context.Context, ea *param.BusAttrExtAPI) (err error) {
	if err = s.dao.ORM.Table("workflow_callback").Where("business = ?", ea.Bid).Update("external_api", ea.ExternalAPI).Error; err != nil {
		log.Error("Failed to set external_api field where bid = %d : %v", ea.Bid, err)
	}
	return
}
