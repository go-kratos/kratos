package notify

import (
	"context"
	"errors"
	"go-common/app/infra/notify/conf"
	"go-common/app/infra/notify/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"net/url"
)

var (
	errUnknownSchema = errors.New("callback url with unknown schema")
)

// Clients notify clients
type Clients struct {
	httpClient     *bm.Client
	liverpcClients *LiverpcClients
}

// NewClients New NotifyClients
func NewClients(c *conf.Config, w *model.Watcher) *Clients {
	nc := &Clients{
		httpClient:     bm.NewClient(c.HTTPClient),
		liverpcClients: newLiverpcClients(w),
	}
	log.Info("Notify.NewClients topic(%s), group(%s), callback len(%d), liverpc clients(%d)",
		w.Topic, w.Group, len(w.Callbacks), len(nc.liverpcClients.clients))
	return nc
}

// Post do callback with different client vary schemas
func (nc *Clients) Post(c context.Context, notifyURL *model.NotifyURL, msg string) (err error) {
	switch notifyURL.Schema {
	case model.LiverpcSchema:
		err = nc.liverpcClients.Post(c, notifyURL, msg)
	case model.HTTPSchema:
		params := url.Values{}
		params.Set("msg", msg)
		client := nc.httpClient
		err = client.Post(c, notifyURL.RawURL, "", params, nil)
	default:
		err = errUnknownSchema
	}
	return
}
