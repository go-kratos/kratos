package search

import (
	"go-common/app/interface/main/reply/conf"
	es "go-common/library/database/elastic"
	bm "go-common/library/net/http/blademaster"
)

const (
	_searchURL    = "/x/internal/search/reply"
	_searchLogURL = "/api/reply/external/search"
)

// Dao search dao.
type Dao struct {
	logURL    string
	searchURL string
	httpCli   *bm.Client
	es        *es.Elastic
}

// New new a dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		logURL:    c.Host.Search + _searchLogURL,
		searchURL: c.Host.API + _searchURL,
		httpCli:   bm.NewClient(c.HTTPClient),
		es:        es.NewElastic(c.Es),
	}
	return d
}
