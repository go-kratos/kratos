package web

import (
	"context"
	"encoding/json"
	xhttp "net/http"
	"net/url"

	"go-common/app/interface/main/web-goblin/model/web"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Recruit .
func (d *Dao) Recruit(ctx context.Context, param url.Values, route *web.Params) (res json.RawMessage, err error) {
	var (
		req     *xhttp.Request
		rs      json.RawMessage
		mokaURI = d.c.Recruit.MokaURI + "/" + route.Route + "/" + d.c.Recruit.Orgid
	)
	if route.JobID != "" {
		mokaURI = mokaURI + "/" + route.JobID
	}
	param.Del("route")
	param.Del("job_id")
	if req, err = xhttp.NewRequest("GET", mokaURI+"?"+param.Encode(), nil); err != nil {
		log.Error("Recruit xhttp.NewRequest url(%s) error(%v)", mokaURI, err)
		return
	}
	if err = d.httpR.Do(ctx, req, &rs); err != nil {
		log.Error("Recruit d.httpR.Do url(%s) error(%v)", mokaURI, err)
		err = ecode.NothingFound
		return
	}
	res = rs
	return
}
