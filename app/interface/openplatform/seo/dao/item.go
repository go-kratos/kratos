package dao

import (
	"context"
)

// GetItem get item page
func (d *Dao) GetItem(c context.Context, id int, bot bool) (res []byte, err error) {
	key := getKey(id, _item, bot)
	res, err = d.GetCache(c, key)
	if err == nil && res != nil {
		return
	}

	url := getUrl(id, _item, bot)
	res, err = d.GetUrl(c, url)
	if err == nil {
		d.AddCache(c, key, res)
	}
	return
}
