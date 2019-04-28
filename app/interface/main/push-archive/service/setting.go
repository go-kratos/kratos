package service

import (
	"context"
	"time"

	"go-common/app/interface/main/push-archive/model"
	pb "go-common/app/service/main/push/api/grpc/v1"
	"go-common/library/log"
)

// Setting gets user's archive-result setting.
func (s *Service) Setting(c context.Context, mid int64) (st *model.Setting, err error) {
	st, err = s.dao.Setting(c, mid)
	if err != nil {
		return
	}
	if st == nil {
		st = &model.Setting{Type: model.PushTypeSpecial} // 如果用户还未上传配置，则默认为特殊关注
	}
	return
}

// SetSetting saves user's archive-result setting.
func (s *Service) SetSetting(c context.Context, mid int64, st *model.Setting) (err error) {
	err = s.dao.SetSetting(c, mid, st)
	if err == nil {
		set := &pb.SetSettingRequest{
			Mid:   mid,
			Type:  1,
			Value: int32(st.Type),
		}
		s.settingCh <- set
	}
	return
}

func (s *Service) setSettingProc() (err error) {
	defer func() {
		if msg := recover(); msg != nil {
			log.Error("setSettingProc got panic(%+v)", msg)
		}
	}()

	for {
		time.Sleep(time.Millisecond * 200)
		set, open := <-s.settingCh
		if !open {
			log.Error("setSettingProc settingCh is closed")
			return
		}

		// before send rpc, check db value with new value, if diff, then rpc is later than another update
		dbSet, err := s.dao.Setting(context.TODO(), set.Mid)
		if err != nil {
			log.Error("setSettingProc s.dao.Setting error(%v) set(%+v)", err, set)
			s.settingCh <- set
			continue
		}
		if dbSet == nil || dbSet.Type != int(set.Value) {
			log.Info("setSettingProc push setting value diff, db(%+v) rpc(%+v)", dbSet, set)
			continue
		}

		// rpc中0-关闭，1-开启
		var tp int32
		if set.Value == model.PushTypeSpecial || set.Value == model.PushTypeAttention {
			tp = 1
		}
		set.Value = tp

		if _, err := s.pushRPC.SetSetting(context.TODO(), set); err != nil {
			log.Error("s.pushRPC.SetSetting error(%v) set(%+v)", err, set)
			s.settingCh <- set
		}
	}
}
