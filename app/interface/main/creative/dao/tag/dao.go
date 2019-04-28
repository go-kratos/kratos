package tag

import (
	"go-common/app/interface/main/creative/conf"
	httpx "go-common/library/net/http/blademaster"
)

const (
	_tagList    = "/x/internal/tag/minfo"
	_tagCheck   = "/x/internal/tag/check"
	_appealTag  = "/videoup/archive/reason/tag"
	_mngTagList = "/x/admin/manager/internal/tag/list"
)

var (
	// StaffTagBid var
	StaffTagBid int64 = 15
)

// Dao  define
type Dao struct {
	c *conf.Config
	// http
	client *httpx.Client
	// uri
	tagList       string
	tagCheck      string
	appealTag     string
	mngTagListURI string
}

// New init dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:             c,
		client:        httpx.NewClient(c.HTTPClient.Normal),
		tagList:       c.Host.API + _tagList,
		tagCheck:      c.Host.API + _tagCheck,
		appealTag:     c.Host.Videoup + _appealTag,
		mngTagListURI: c.Host.MainSearch + _mngTagList,
	}
	return
}
