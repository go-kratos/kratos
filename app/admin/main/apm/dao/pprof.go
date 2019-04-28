package dao

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"go-common/app/admin/main/apm/model/pprof"
	"go-common/library/conf/env"
	"go-common/library/log"
)

// Instances get instances
func (d *Dao) Instances(c context.Context, appName string) (ins *pprof.Ins, err error) {
	var (
		req    *http.Request
		uri    = d.c.Host.SVENCo + "/x/admin/apm/noauth/discovery/fetch"
		ret    = &pprof.Response{}
		params = url.Values{}
	)
	params.Add("zone", env.Zone)
	params.Add("env", env.DeployEnv)
	params.Add("region", env.Region)
	params.Add("appid", appName)
	params.Add("status", "3")
	if req, err = d.client.NewRequest(http.MethodGet, uri, "", params); err != nil {
		log.Error("d.Instances http.NewRequest error(%v)", err)
		return
	}
	if err = d.client.Do(c, req, ret); err != nil {
		log.Error("d.Instances client Do error(%v)", err)
		return
	}
	if ret.Code != 0 {
		err = fmt.Errorf("%s params(%s) response return_code(%d)", uri, params.Encode(), ret.Code)
		log.Error("d.Instance error(%v)", err)
		return
	}
	ins = ret.Data
	return
}
