package dao

import (
	"context"
	"go-common/app/admin/main/tv/model"
	"go-common/library/cache/memcache"
)

//ModulePublish is used for module status MC key
const ModulePublish = "ModulePublish"

//SetModPub is used for set module publish status to MC
func (d *Dao) SetModPub(c context.Context, pageID string, p model.ModPub) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	itemJSON := &memcache.Item{
		Key:        ModulePublish + pageID,
		Object:     p,
		Flags:      memcache.FlagJSON,
		Expiration: 0,
	}
	if err = conn.Set(itemJSON); err != nil {
		return
	}
	return
}

//GetModPub is used for getting module publish status
func (d *Dao) GetModPub(c context.Context, pageID string) (p model.ModPub, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = d.mc.Get(c)
	defer conn.Close()
	k := ModulePublish + pageID
	if item, err = conn.Get(k); err != nil {
		return
	}
	if err = conn.Scan(item, &p); err != nil {
		return
	}
	return
}
