package usersuit

import (
	"go-common/app/interface/main/account/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao struct.
type Dao struct {
	c             *conf.Config
	http          *bm.Client
	orderURL      string
	orderHistory  string
	groupURL      string
	entryGroupURL string
	vipGroupURL   string
	pendantURL    string
	packageURL    string
}

// New new a Dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		http:          bm.NewClient(c.HTTPClient.Normal),
		groupURL:      c.Host.API + _groupInfo,
		orderURL:      c.Host.API + _orderInfo,
		orderHistory:  c.Host.API + _orderHistory,
		entryGroupURL: c.Host.API + _entryGroup,
		vipGroupURL:   c.Host.API + _vipGroup,
		pendantURL:    c.Host.API + _pendantInfo,
		packageURL:    c.Host.API + _packageInfo,
	}
	return
}
