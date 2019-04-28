package staff

import (
	"go-common/app/admin/main/videoup/conf"
	bm "go-common/library/net/http/blademaster"
)

const (
	_staffListURL   = "/videoup/staff"
	_staffSubmitURL = "/videoup/staff/apply/batch"
)

// Dao is search dao
type Dao struct {
	c          *bm.ClientConfig
	httpClient *bm.Client
	staffURI   string
	submitUrl  string
}

var (
	d *Dao
)

// New new staff dao
func New(c *conf.Config) *Dao {
	return &Dao{
		c:          c.HTTPClient.Read,
		httpClient: bm.NewClient(c.HTTPClient.Read),
		staffURI:   c.Host.Archive + _staffListURL,
		submitUrl:  c.Host.Archive + _staffSubmitURL,
	}
}
