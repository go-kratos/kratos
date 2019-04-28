package dao

import (
	"context"
	"net/url"
	"strings"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_uidsURL   = "/x/admin/manager/users/uids"
	_unamesURL = "/x/admin/manager/users/unames"
)

// Unames get unames by uid
func (d *Dao) Unames(c context.Context, uids []int64) (res map[int64]string, err error) {
	var (
		param    = url.Values{}
		uidStr   = xstr.JoinInts(uids)
		unameURI = d.c.Host.Manager + _unamesURL
	)
	param.Set("uids", uidStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[int64]string `json:"data"`
		Message string           `json:"message"`
	}

	err = d.hclient.Get(c, unameURI, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", unameURI+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", unameURI+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

// Uids get uids by unames
func (d *Dao) Uids(c context.Context, names []string) (res map[string]int64, err error) {
	var (
		param    = url.Values{}
		namesStr = strings.Join(names, ",")
		uidURI   = d.c.Host.Manager + _uidsURL
	)
	param.Set("unames", namesStr)

	var httpRes struct {
		Code    int              `json:"code"`
		Data    map[string]int64 `json:"data"`
		Message string           `json:"message"`
	}

	err = d.hclient.Get(c, uidURI, "", param, &httpRes)
	if err != nil {
		log.Error("d.client.Get(%s) error(%v)", uidURI+"?"+param.Encode(), err)
		return
	}
	if httpRes.Code != 0 {
		log.Error("url(%s) error(%v), code(%d), message(%s)", uidURI+"?"+param.Encode(), err, httpRes.Code, httpRes.Message)
	}
	res = httpRes.Data
	return
}

// GetUIDByName 获取uid
func (d *Dao) GetUIDByName(c context.Context, name string) (uid int64, err error) {
	var res map[string]int64
	if res, err = d.Uids(c, []string{name}); err != nil {
		return
	}
	if uid, ok := res[name]; ok {
		return uid, nil
	}
	return
}

// GetNameByUID 获取用户名
func (d *Dao) GetNameByUID(c context.Context, uids []int64) (mcases map[int64][]interface{}, err error) {
	var res map[int64]string

	if res, err = d.Unames(c, uids); err != nil {
		return
	}
	mcases = make(map[int64][]interface{})
	for uid, uname := range res {
		mcases[uid] = []interface{}{uname}
	}
	return
}
