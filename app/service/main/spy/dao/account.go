package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/service/main/spy/model"
	"go-common/library/log"
)

// tel level.
const (
	TelRiskLevelLow     = 1
	TelRiskLevelMedium  = 2
	TelRiskLevelHigh    = 3
	TelRiskLevelUnknown = 4
)

// TelRiskLevel tel risk level.
func (d *Dao) TelRiskLevel(c context.Context, mid int64) (riskLevel int8, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("type", "123")
	var resp struct {
		TS   int64 `json:"ts"`
		Code int64 `json:"code"`
		Data struct {
			Level int8  `json:"level"`
			Mid   int64 `json:"mid"`
			Score int64 `json:"score"`
		} `json:"data"`
	}
	if err = d.httpClient.Get(c, d.c.Property.TelValidateURL, "", params, &resp); err != nil {
		log.Error("d.httpClient.Do() error(%v) , riskLevel = TelRiskLevelUnknown", err)
		riskLevel = TelRiskLevelUnknown
		err = nil
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("GET TelRiskLevel url resp(%v)", resp)
		return
	}
	log.Info("GET TelValidateURL suc url(%s) resp(%v)", d.c.Property.TelValidateURL+"?"+params.Encode(), resp)
	riskLevel = resp.Data.Level
	return
}

// BlockAccount block account.
func (d *Dao) BlockAccount(c context.Context, mid int64) (err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	params.Set("admin_reason", "spy")
	params.Set("blockType", "3")
	params.Set("operator", "spy")
	params.Set("type", "json")
	var resp struct {
		Code int64 `json:"code"`
	}
	if err = d.httpClient.Get(c, d.c.Property.BlockAccountURL, "", params, &resp); err != nil {
		log.Error("httpClient.Do() error(%v)", err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("GET block account url resp(%v)", resp)
	}
	return
}

// SecurityLogin security login
func (d *Dao) SecurityLogin(c context.Context, mid int64, reason string) (err error) {
	params := url.Values{}
	params.Set("mids", fmt.Sprintf("%d", mid))
	params.Set("operator", "spy")
	params.Set("desc", reason)
	params.Set("type", "json")
	var resp struct {
		Code int64 `json:"code"`
	}
	if err = d.httpClient.Post(c, d.c.Property.SecurityLoginURL, "", params, &resp); err != nil {
		log.Error("message url(%s) error(%v)", d.c.Property.SecurityLoginURL+"?"+params.Encode(), err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("POST SecurityLogin url resp(%v)", resp)
		return
	}
	log.Info("POST SecurityLogin suc url(%s) resp(%v)", resp)
	return
}

// TelInfo tel info.
func (d *Dao) TelInfo(c context.Context, mid int64) (tel *model.TelInfo, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	var resp struct {
		Code int64         `json:"code"`
		Data model.TelInfo `json:"data"`
	}
	if err = d.httpClient.Get(c, d.c.Property.TelInfoByMidURL, "", params, &resp); err != nil {
		log.Error("d.httpClient.Do() error(%v) ,TelInfo", err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("GET TelRiskLevel url resp(%v)", resp)
		return
	}
	tel = &resp.Data
	log.Info("GET TelInfoByMidURL suc url(%s) resp(%v)", d.c.Property.TelInfoByMidURL+"?"+params.Encode(), resp)
	return
}

// ProfileInfo get user profile info from account service.
func (d *Dao) ProfileInfo(c context.Context, mid int64, remoteIP string) (profile *model.ProfileInfo, err error) {
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	var resp struct {
		Code int64             `json:"code"`
		Data model.ProfileInfo `json:"data"`
	}
	if err = d.httpClient.Get(c, d.c.Property.ProfileInfoByMidURL, remoteIP, params, &resp); err != nil {
		log.Error("d.httpClient.Do() error(%v), ProfileInfo", err)
		return
	}
	if resp.Code != 0 {
		err = fmt.Errorf("GET ProfileInfo url resp(%v)", resp)
		return
	}
	profile = &resp.Data
	log.Info("GET ProfileInfoByMidURL suc url(%s) resp(%v)", d.c.Property.ProfileInfoByMidURL+"?"+params.Encode(), resp)
	return
}
