package reply

import (
	"go-common/app/job/main/archive/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	client       *bm.Client
	changeSubMid string
}

// New is new redis dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: bm.NewClient(c.HTTPClient),
		// path
		changeSubMid: c.Host.APICo + _changeSubjectMid,
	}
	return d
}
