package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/interface/main/dm2/model"
	"go-common/library/log"
)

const (
	_waveFormCallBackSuccess = 1

	_waveFormPrefix = "http://i0.hdslb.com/bfs"

	_expire = 20 // 2*10
)

func (s *Service) waveForm(c context.Context, oid int64, tp int32) (waveForm *model.WaveForm, err error) {
	var (
		cacheError bool
	)
	if waveForm, err = s.dao.WaveFormCache(c, oid, tp); err != nil {
		cacheError = true
		err = nil
	}
	if waveForm != nil && !waveForm.Empty {
		return
	}
	if waveForm, err = s.dao.GetWaveForm(c, oid, tp); err != nil {
		log.Error("params(oid:%v,tp:%v),error(%v)", oid, tp, err)
		return
	}
	if waveForm == nil {
		waveForm = &model.WaveForm{
			Oid:   oid,
			Type:  tp,
			Empty: true,
		}
	}
	if !cacheError {
		temp := waveForm
		s.cache.Do(c, func(ctx context.Context) {
			s.dao.SetWaveFormCache(ctx, temp)
		})
	}
	return
}

// WaveForm .
func (s *Service) WaveForm(c context.Context, aid, oid int64, tp int32, mid int64) (waveFormResp *model.WaveFormResp, err error) {
	var (
		uposErr     error
		waveFromURL string
		waveForm    *model.WaveForm
	)
	if err = s.SubtitlePermission(c, aid, oid, tp, mid); err != nil {
		return
	}
	if waveForm, err = s.waveForm(c, oid, tp); err != nil {
		log.Error("params(oid:%v,tp:%v),error(%v)", oid, tp, err)
		return
	}
	if !waveForm.Empty {
		waveFormResp = &model.WaveFormResp{
			State:       waveForm.State,
			WaveFromURL: waveForm.WaveFromURL,
		}
		switch waveForm.State {
		case model.WaveFormStatusFailed, model.WaveFormStatusSuccess:
			return
		case model.WaveFormStatusWaitting:
			if time.Since(time.Unix(waveForm.Mtime, 0)) < _expire {
				return
			}
		}
	}
	waveForm.State = model.WaveFormStatusWaitting
	if waveFromURL, uposErr = s.dao.Upos(c, oid); uposErr != nil {
		log.Error("postUpos(oid:%v),error(%v)", oid, err)
		waveForm.State = model.WaveFormStatusError
	}
	waveForm.WaveFromURL = fmt.Sprintf("%s/%s", _waveFormPrefix, waveFromURL)
	if err = s.dao.UpsertWaveFrom(c, waveForm); err != nil {
		log.Error("params(waveForm:%+v),error(%v)", waveForm, err)
		return
	}
	if err = s.dao.DelWaveFormCache(c, oid, tp); err != nil {
		log.Error("DelWaveFormCache.params(oid:%v,tp:%v),error(%v)", oid, tp, err)
		return
	}
	waveFormResp = &model.WaveFormResp{
		State:       waveForm.State,
		WaveFromURL: waveForm.WaveFromURL,
	}
	return
}

// WaveFormCallBack .
func (s *Service) WaveFormCallBack(c context.Context, oid int64, tp int32, code int32, info string) (err error) {
	var (
		waveForm *model.WaveForm
	)
	if waveForm, err = s.waveForm(c, oid, tp); err != nil {
		log.Error("params(oid:%v,tp:%v),error(%v)", oid, tp, err)
		return
	}
	if code == _waveFormCallBackSuccess {
		waveForm.State = model.WaveFormStatusSuccess
	} else {
		waveForm.State = model.WaveFormStatusFailed
		log.Error("WaveFormCallBack.params(oid:%v,tp:%v).errorInfo(%s)", oid, tp, info)
	}
	if err = s.dao.UpsertWaveFrom(c, waveForm); err != nil {
		log.Error("params(waveForm:%+v),error(%v)", waveForm, err)
		return
	}
	if err = s.dao.DelWaveFormCache(c, oid, tp); err != nil {
		log.Error("DelWaveFormCache.params(oid:%v,tp:%v),error(%v)", oid, tp, err)
		return
	}
	return
}
