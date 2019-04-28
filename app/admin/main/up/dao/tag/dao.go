package tag

import (
	"go-common/app/admin/main/up/conf"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_tagList   = "/x/internal/tag/minfo"
	_tagCheck  = "/x/internal/tag/check"
	_appealTag = "/videoup/archive/reason/tag"
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	tagList   string
	tagCheck  string
	appealTag string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:         c,
		client:    httpx.NewClient(c.HTTPClient.Normal),
		tagList:   c.Host.API + _tagList,
		tagCheck:  c.Host.API + _tagCheck,
		appealTag: c.Host.Videoup + _appealTag,
	}
	return
}
