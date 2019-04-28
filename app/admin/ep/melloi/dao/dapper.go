package dao

import (
	"context"
	"net/http"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_dapperDependUri = "/x/internal/dapper/service-depend?service_name="
	_dapperRspCode   = 0
)

// QueryServiceDepend query service depend
func (d *Dao) QueryServiceDepend(c context.Context, serviceName string) (ret []string, err error) {
	var (
		req       *http.Request
		DependURL = d.c.Dapper.Host + _dapperDependUri + serviceName
		res       struct {
			Code    int                   `json:"code"`
			Data    *model.DependResponse `json:"data"`
			Message string                `json:"message"`
		}
	)

	if req, err = d.newRequest(http.MethodGet, DependURL, nil); err != nil {
		return
	}

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.Depend url(%s) res($s) error(%v)", DependURL, res, err)
		return
	}
	if res.Code != _dapperRspCode {
		err = ecode.MelloiTreeRequestErr
		log.Error("d.Tree.Response url(%s) resCode(%s) error(%v)", DependURL, res.Code, err)
		return
	}

	for _, item := range res.Data.Items {
		ret = append(ret, item.ServiceName)
	}

	return
}
