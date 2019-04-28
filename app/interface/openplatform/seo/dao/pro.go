package dao

import (
	"context"
)

// GetPro get project page
func (d *Dao) GetPro(c context.Context, id int, bot bool) (res []byte, err error) {
	key := getKey(id, _pro, bot)
	res, err = d.GetCache(c, key)
	if err == nil && res != nil {
		return
	}

	url := getUrl(id, _pro, bot)
	res, err = d.GetUrl(c, url)
	if err == nil {
		d.AddCache(c, key, res)
	}
	return
}
