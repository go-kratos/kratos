package videoup

import (
	"go-common/app/admin/main/mcn/conf"
	"go-common/app/service/main/videoup/model/archive"
	bm "go-common/library/net/http/blademaster"
)

const (
	_typeURL = "/videoup/types"
)

// Dao .
type Dao struct {
	c                *conf.Config
	client           *bm.Client
	videTypeURL      string
	videoUpTypeCache map[int]archive.Type
}

// New new a Dao and return.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		// http client
		client:           bm.NewClient(c.HTTPClient),
		videTypeURL:      c.Host.Videoup + _typeURL,
		videoUpTypeCache: make(map[int]archive.Type),
	}
	d.refreshUpType()
	go d.refreshUpTypeAsync()
	return
}
