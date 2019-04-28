package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"go-common/app/admin/main/open/model"
	"go-common/library/log"
)

// AddApp represents add an app.
func (s *Service) AddApp(c context.Context, appname string) (err error) {
	// Generate the secret and key
	timestamp := strconv.Itoa((int)(time.Now().Unix()) + rand.Intn(1))
	hashsec := md5.New()
	hashsec.Write([]byte(timestamp))
	appsecret := hex.EncodeToString(hashsec.Sum(nil))
	hashkey := md5.New()
	timestamp += "biliappkey"
	hashkey.Write([]byte(timestamp))
	appkey := hex.EncodeToString(hashkey.Sum(nil))
	g := &model.App{
		AppName:   appname,
		AppKey:    appkey,
		AppSecret: appsecret,
		Enabled:   1,
	}
	if err = s.dao.AddApp(c, g); err != nil {
		log.Error("s.dao.AddApp(%+v) error(%v)", g, err)
	}
	return
}

// DelApp represents delete an app.
func (s *Service) DelApp(c context.Context, appid int64) (err error) {
	if err = s.dao.DelApp(c, appid); err != nil {
		log.Error("s.DelApp(%d) error(%v)", appid, err)
	}
	return
}

// UpdateApp represents update an app.
func (s *Service) UpdateApp(c context.Context, arg *model.AppParams) (err error) {
	if err = s.dao.UpdateApp(c, arg); err != nil {
		log.Error("s.UpdateApps id (%d) appname(%s) error", arg.AppID, arg.AppName)
	}
	return
}

// ListApp represents search an app.
func (s *Service) ListApp(c context.Context, t *model.AppListParams) (res []*model.App, total int64, err error) {
	if res, err = s.dao.ListApp(c, t); err != nil {
		log.Error("s.dao.ListApp error (%v)", err)
		return
	}
	// output the data
	total = int64(len(res))
	start := (t.PN - 1) * t.PS
	if start >= total {
		res = []*model.App{}
		return
	}
	end := start + t.PS
	if end > total {
		end = total
	}
	res = res[start:end]
	return
}
