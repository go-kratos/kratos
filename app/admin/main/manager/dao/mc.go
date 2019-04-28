package dao

import (
	"context"

	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/net/http/blademaster/middleware/permit"
)

// Session .
func (d *Dao) Session(ctx context.Context, sid string) (res *permit.Session, err error) {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	r, err := conn.Get(sid)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Error("conn.Get(%s) error(%v)", sid, err)
		return
	}
	res = &permit.Session{}
	if err = conn.Scan(r, res); err != nil {
		log.Error("conn.Scan(%s) error(%v)", string(r.Value), err)
	}
	return
}

// SetSession .
func (d *Dao) SetSession(ctx context.Context, p *permit.Session) (err error) {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	item := &memcache.Item{
		Key:        p.Sid,
		Object:     p,
		Flags:      memcache.FlagJSON,
		Expiration: int32(_sessionLife),
	}
	err = conn.Set(item)
	return
}
