package dao

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/cache/memcache"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_genPaasMachines           = "/api/merlin/machine/create"
	_delPaasMachine            = "/api/merlin/machine/free"
	_queryPaasMachineStatus    = "/api/merlin/machine/status"
	_queryPaasMachine          = "/api/merlin/machine/detail"
	_updatePaasMachineNode     = "/api/merlin/machine/update"
	_updatePaasMachineSnapshot = "/api/merlin/machine/snapshot"
	_queryPaasClusters         = "/api/merlin/clusters"
	_queryPaasClusterByNetwork = "/api/merlin/clusters/network/"
	_auth                      = "/api/v1/auth"
	_authHeader                = "X-Authorization-Token"
)

// GenPaasMachines create machine in paas.
func (d *Dao) GenPaasMachines(c context.Context, mc *model.PaasGenMachineRequest) (instances []*model.CreateInstance, err error) {
	var (
		req *http.Request
		res = &model.PaasGenMachineResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodPost, _genPaasMachines, mc); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("gen paas machine url(%s) err(%v)", _genPaasMachines, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	instances = res.Data
	return
}

// DelPaasMachine delete machine in paas.
func (d *Dao) DelPaasMachine(c context.Context, pqadmr *model.PaasQueryAndDelMachineRequest) (instance *model.ReleaseInstance, err error) {
	var (
		req *http.Request
		res = &model.PaasDelMachineResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodPost, _delPaasMachine, pqadmr); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("delete paas machine url(%s) err(%v)", _delPaasMachine, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	instance = &res.Data
	return
}

// QueryPaasMachineStatus query status of machine in paas.
func (d *Dao) QueryPaasMachineStatus(c context.Context, pqadmr *model.PaasQueryAndDelMachineRequest) (machineStatus *model.MachineStatus, err error) {
	var (
		req *http.Request
		res = &model.PaasQueryMachineStatusResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodPost, _queryPaasMachineStatus, pqadmr); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("query paas machine status url(%s) err(%v)", _queryPaasMachineStatus, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	machineStatus = &res.Data
	return
}

// SnapshotPaasMachineStatus Snapshot Paas Machine Status.
func (d *Dao) SnapshotPaasMachineStatus(c context.Context, pqadmr *model.PaasQueryAndDelMachineRequest) (status int, err error) {
	var (
		req *http.Request
		res = &model.PaasSnapshotMachineResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodPost, _updatePaasMachineSnapshot, pqadmr); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("snapshot machine status url(%s) err(%v)", _updatePaasMachineSnapshot, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}

	status = res.Status
	return
}

// QueryPaasMachine query detail information of machine in paas.
func (d *Dao) QueryPaasMachine(c context.Context, pqadmr *model.PaasQueryAndDelMachineRequest) (md *model.PaasMachineDetail, err error) {
	var (
		req *http.Request
		res = &model.PaasQueryMachineResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodPost, _queryPaasMachine, pqadmr); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("query paas machine url(%s) err(%v)", _queryPaasMachine, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	md = &res.Data
	return
}

// QueryClusters query cluster information in paas.
func (d *Dao) QueryClusters(c context.Context) (clusters []*model.Cluster, err error) {
	var (
		req *http.Request
		res = &model.PaasQueryClustersResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodGet, _queryPaasClusters, nil); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.httpSearch url(%s) error(%v)", d.c.Paas.Host+"?"+_queryPaasClusters, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	clusters = res.Data.Items
	return
}

// QueryCluster query cluster information in paas by giving network.
func (d *Dao) QueryCluster(c context.Context, netWordID int64) (cluster *model.Cluster, err error) {
	var (
		req *http.Request
		res = &model.PaasQueryClusterResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodGet, _queryPaasClusterByNetwork+strconv.FormatInt(netWordID, 10), nil); err != nil {
		log.Error("http new request err(%v)", err)
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.httpSearch url(%s) error(%v)", d.c.Paas.Host+"?"+_queryPaasClusters, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	cluster = res.Data
	return
}

// UpdatePaasMachineNode update paas machine node.
func (d *Dao) UpdatePaasMachineNode(c context.Context, pumnr *model.PaasUpdateMachineNodeRequest) (data string, err error) {
	var (
		req *http.Request
		res = &model.PaasUpdateMachineNodeResponse{}
	)
	if req, err = d.newPaasRequest(c, http.MethodPost, _updatePaasMachineNode, pumnr); err != nil {
		log.Error("http new request err(%v)", err)
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.httpSearch url(%s) error(%v)", d.c.Paas.Host+_updatePaasMachineNode, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	data = res.Data
	return
}

func (d *Dao) authPaas(c context.Context) (token string, err error) {
	var (
		req         *http.Request
		res         = &model.PaasAuthResponse{}
		authRequest = model.PaasAuthRequest{
			APIToken:   d.c.Paas.Token,
			PlatformID: "merlin",
		}
	)
	if req, err = d.newRequest(http.MethodPost, d.c.Paas.Host+_auth, authRequest); err != nil {
		return
	}
	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("query paas machine url(%s) err(%v)", _auth, err)
		err = ecode.MerlinPaasRequestErr
		return
	}
	if err = res.CheckStatus(); err != nil {
		return
	}
	token = res.Data.Token
	return
}

// paasToken TODO:当前放在dao层有点不规范，放在service层，封装上又不如这样更好，后续再考虑一下.
func (d *Dao) paasToken(c context.Context) (authToken string, err error) {
	var (
		item *memcache.Item
		conn = d.mc.Get(c)
	)
	defer conn.Close()
	if item, err = conn.Get(d.c.Paas.Token); err == nil {
		if err = json.Unmarshal(item.Value, &authToken); err != nil {
			log.Error("Json unmarshal err(%v)", err)
		}
		return
	}
	if authToken, err = d.authPaas(c); err != nil {
		return
	}
	item = &memcache.Item{Key: d.c.Paas.Token, Object: authToken, Flags: memcache.FlagJSON, Expiration: d.expire}
	d.tokenCacheSave(c, item)
	return
}

func (d *Dao) newPaasRequest(c context.Context, method, uri string, v interface{}) (req *http.Request, err error) {
	var authToken string
	if authToken, err = d.paasToken(c); err != nil {
		return
	}
	if req, err = d.newRequest(method, d.c.Paas.Host+uri, v); err != nil {
		return
	}
	req.Header.Set(_authHeader, authToken)
	return
}
