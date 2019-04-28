package rank

import (
	"context"

	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-intl/conf"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

const (
	_allRank = "/data/rank/recent_all-app.json"
)

// Dao is rank dao.
type Dao struct {
	// http client
	clientAsyn *httpx.Client
	// all rank
	allRank string
}

// New new a rank dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		// http client
		clientAsyn: httpx.NewClient(c.HTTPClientAsyn),
		allRank:    c.Host.Rank + _allRank,
	}
	return d
}

// AllRank is.
func (d *Dao) AllRank(c context.Context) (ranks []*rank.Rank, err error) {
	var res struct {
		Code int          `json:"code"`
		List []*rank.Rank `json:"list"`
	}
	if err = d.clientAsyn.Get(c, d.allRank, "", nil, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.Wrap(err, d.allRank)
		return
	}
	ranks = res.List
	return
}
