package creative

import (
	"context"

	"go-common/app/interface/main/videoup/conf"
	bm "go-common/library/net/http/blademaster"
)

// Dao is elec dao.
type Dao struct {
	c                 *conf.Config
	httpW             *bm.Client
	setWatermarkURL   string
	uploadMaterialURL string
}

// New new a elec dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                 c,
		httpW:             bm.NewClient(c.HTTPClient.Write),
		setWatermarkURL:   c.Host.APICo + _setWatermark,
		uploadMaterialURL: c.Host.APICo + _uploadMaterial,
	}
	return d
}

// Ping ping success.
func (d *Dao) Ping(c context.Context) (err error) {
	return
}
