package dao

import (
	"context"
	"errors"

	"go-common/app/interface/openplatform/seo/conf"
)

// Sitemap 生成站点地图
func (d *Dao) Sitemap(c context.Context, host string) (res []byte, err error) {
	s := conf.GetSitemap(host)
	if s == nil || s.Url == "" {
		err = errors.New(host + " sitemap config not exist")
		return
	}

	res, err = d.GetUrl(c, s.Url)
	return
}
