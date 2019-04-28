package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// const .
const (
	_linkTypeH5  = 7
	_linkTypeApp = 10
)

// SavePushData save push data
func (s *Service) SavePushData(c context.Context, arg *model.VipPushData) (err error) {
	if err = s.checkPushData(arg); err != nil {
		err = errors.WithStack(err)
		return
	}

	if arg.ID == 0 {
		arg.ProgressStatus = model.NotStart
		arg.Status = model.Normal
		if _, err = s.dao.AddPushData(c, arg); err != nil {
			err = errors.WithStack(err)
		}
		return
	}

	var pushData *model.VipPushData
	if pushData, err = s.dao.GetPushData(c, arg.ID); err != nil {
		err = errors.WithStack(err)
		return
	}
	if pushData == nil {
		err = ecode.VipPushDataNotExitErr
		return
	}

	if pushData.Status == model.Fail {
		err = ecode.VipPushDataUpdateErr
		return
	}

	if pushData.ProgressStatus == model.Started {
		err = ecode.VipPushDataUpdateErr
		return
	}

	if pushData.GroupName != arg.GroupName || !pushData.EffectStartDate.Time().Equal(arg.EffectStartDate.Time()) {
		if !(pushData.Status == model.Normal && pushData.ProgressStatus == model.NotStart) {
			err = ecode.VipPushDataUpdateErr
			return
		}
	}

	arg.ProgressStatus = pushData.ProgressStatus
	if pushData.PushedCount == arg.PushTotalCount {
		arg.ProgressStatus = model.Started
	}

	if _, err = s.dao.UpdatePushData(c, arg); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (s *Service) checkPushData(arg *model.VipPushData) (err error) {
	if len(arg.GroupName) > 10 {
		err = ecode.VipPushGroupLenErr
		return
	}

	if len(arg.Title) > 30 {
		err = ecode.VipPushTitleLenErr
		return
	}

	if len(arg.Content) > 200 {
		err = ecode.VipPushContentLenErr
		return
	}

	if arg.LinkType != _linkTypeH5 && arg.LinkType != _linkTypeApp {
		err = ecode.VipPushLinkTypeErr
		return
	}

	if arg.EffectEndDate.Time().Before(arg.EffectStartDate.Time()) {
		err = ecode.VipPushEffectTimeErr
		return
	}

	duration := arg.EffectEndDate.Time().Sub(arg.EffectStartDate.Time())
	day := duration.Hours() / 24
	arg.PushTotalCount = int32(day) + 1

	if _, err = time.Parse("15:04:05", arg.PushStartTime); err != nil {
		err = ecode.VipPushFmtTimeErr
		return
	}
	if _, err = time.Parse("15:04:05", arg.PushEndTime); err != nil {
		err = ecode.VipPushFmtTimeErr
		return
	}

	platformMap := make(map[string]*model.PushDataPlatform)
	platformArr := make([]*model.PushDataPlatform, 0)
	var (
		data      []byte
		key       string
		condition string
		ok        bool
	)
	if err = json.Unmarshal([]byte(arg.Platform), &platformArr); err != nil {
		log.Error("error(%+v)", err)
		err = ecode.VipPushPlatformErr
		return
	}

	for _, v := range platformArr {

		if key, ok = model.PushPlatformNameMap[v.Name]; !ok {
			err = ecode.VipPushPlatformErr
			return
		}
		if condition, ok = model.ConditionNameMap[v.Condition]; !ok {
			err = ecode.VipPushPlatformErr
			return
		}

		v.Condition = condition
		platformMap[key] = v

	}

	if data, err = json.Marshal(platformMap); err != nil {
		err = errors.WithStack(err)
		return
	}

	arg.Platform = string(data)
	return
}

// GetPushData get push data
func (s *Service) GetPushData(c context.Context, id int64) (res *model.VipPushData, err error) {
	if res, err = s.dao.GetPushData(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		return
	}
	res.PushProgress = fmt.Sprintf("%v/%v", res.PushedCount, res.PushTotalCount)
	if err = s.fmtPushDataPlatform(res); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DisablePushData .
func (s *Service) DisablePushData(c context.Context, id int64) (err error) {
	var (
		res *model.VipPushData
		now = time.Now()
	)
	if res, err = s.dao.GetPushData(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		err = ecode.VipPushDataNotExitErr
		return
	}

	if !(res.ProgressStatus == model.Starting && res.Status == model.Normal && res.DisableType == model.UnDisable) {
		err = ecode.VipPushDataDisableErr
		return
	}
	duration := now.Sub(res.EffectStartDate.Time())
	day := duration.Hours() / 24
	res.PushTotalCount = int32(day) + 1
	res.EffectEndDate = xtime.Time(now.Unix())
	if res.PushTotalCount > res.PushedCount {
		res.PushTotalCount--
	}
	if res.PushTotalCount == res.PushedCount {
		res.ProgressStatus = model.Started
	}
	if err = s.dao.DisablePushData(c, res); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DelPushData .
func (s *Service) DelPushData(c context.Context, id int64) (err error) {
	var res *model.VipPushData
	if res, err = s.dao.GetPushData(c, id); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		err = ecode.VipPushDataNotExitErr
		return
	}
	if !(res.ProgressStatus == model.NotStart && res.Status == model.Normal) {
		err = ecode.VipPushDataDelErr
		return
	}

	if err = s.dao.DelPushData(c, id); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// PushDatas get push datas
func (s *Service) PushDatas(c context.Context, arg *model.ArgPushData) (res []*model.VipPushData, count int64, err error) {
	if count, err = s.dao.PushDataCount(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}

	if res, err = s.dao.PushDatas(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, v := range res {
		v.PushProgress = fmt.Sprintf("%v/%v", v.PushedCount, v.PushTotalCount)
	}
	return
}

func (s *Service) fmtPushDataPlatform(res *model.VipPushData) (err error) {
	var platformArr []*model.PushDataPlatform

	platform := make(map[string]*model.PushDataPlatform)
	if err = json.Unmarshal([]byte(res.Platform), &platform); err != nil {
		err = errors.WithStack(err)
		return
	}

	for k, v := range platform {
		r := new(model.PushDataPlatform)
		r.Name = model.PushPlatformMap[k]
		r.Build = v.Build
		r.Condition = model.ConditionMap[v.Condition]
		platformArr = append(platformArr, r)
	}
	res.PlatformArr = platformArr
	return
}
