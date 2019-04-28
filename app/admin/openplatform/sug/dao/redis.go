package dao

import (
	"context"
	"fmt"
	//
	"go-common/app/admin/openplatform/sug/model"
	"go-common/library/cache/redis"
)

const (
	_sugSeason = "SUG:SEASON:%d"
	_sugItem   = "SUGITEM:%d"
	_expire    = 93600
)

// SetSug set redis sug
func (d *Dao) SetSug(c context.Context, seasonID int64, ItemID int64, score int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	key := fmt.Sprintf(_sugSeason, seasonID)
	if err = conn.Send("ZADD", key, score, ItemID); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", key, _expire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	return
}

// DelSug del redis sug
func (d *Dao) DelSug(c context.Context, seasonID, itemsID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if _, err = conn.Do("ZREM", fmt.Sprintf(_sugSeason, seasonID), itemsID); err != nil {
		return
	}
	return
}

// SetItem set redis item
func (d *Dao) SetItem(c context.Context, item *model.Item, location string) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ItemsID), "items_id", item.ItemsID); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ItemsID), "name", item.Name); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ItemsID), "brief", item.Brief); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ItemsID), "head", item.Img); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ItemsID), "pic", location); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", fmt.Sprintf(_sugItem, item.ItemsID), _expire); err != nil {
		return
	}
	if err = conn.Flush(); err != nil {
		return
	}
	if b, err = redis.Bool(conn.Receive()); err != nil {
		return
	}
	conn.Receive()
	return
}
