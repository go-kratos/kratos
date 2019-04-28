package dao

import (
	"context"
	"net/http"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_rmTokenURI     = "/api/v1/auth"
	_clusterNodeURI = "/api/v1/clusters?env=uat"
)

//RmToken get PaaS token
func (d *Dao) RmToken(c context.Context) (token string, err error) {
	var (
		req         *http.Request
		tokenURL    = d.c.ServiceCluster.TestHost + _rmTokenURI
		res         = &model.ClusterRmTokenResponse{}
		rmTokenPost = &model.RmTokenPost{
			APIToken:   conf.Conf.Paas.APIToken,
			PlatformID: conf.Conf.Paas.PlatformID,
		}
	)
	if req, err = d.newRequest(http.MethodPost, tokenURL, rmTokenPost); err != nil {
		log.Error("query paas machine url(%s) err(%v)", tokenURL, err)
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("RmToken error :(%v)", err)
		err = ecode.MelloiPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	token = res.Data.Token
	return
}

//NetInfo get melloi server use
func (d *Dao) NetInfo(c context.Context, token string) (cluster []*model.ClusterResponseItemsSon, err error) {
	var (
		url        = d.c.ServiceCluster.UatHost + _clusterNodeURI
		req        *http.Request
		resCluster = &model.ClusterResponse{}
	)

	if req, err = d.newRequest(http.MethodGet, url, nil); err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Authorization-Token", token)

	if err = d.httpClient.Do(c, req, &resCluster); err != nil {
		log.Error("d.Token url(%s) res($s) err(%v)", url, err)
		err = ecode.MelloiPaasRequestErr
		return
	}
	if err = resCluster.CheckStatus(); err != nil {
		return
	}
	cluster = resCluster.Data.Items
	return
}
