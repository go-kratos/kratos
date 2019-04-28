package dao

import (
	"context"
	"fmt"

	amodel "go-common/app/admin/openplatform/sug/model"
	"go-common/app/job/openplatform/open-sug/model"
	"go-common/library/cache/redis"
	"strconv"
)

const (
	_sugSeason = "SUG:SEASON:%s"
	_sugItem   = "SUGITEM:%d"
	_expire    = 93600
)

// SetItem set item info to redis
func (d *Dao) SetItem(c context.Context, item *model.Item) (b bool, err error) {
	var location string
	picItem := amodel.Item{ItemsID: item.ID, Name: item.Name, Brief: item.Brief, Img: item.HeadImg}
	if location, _ = d.CreateItemPNG(picItem); location == "" {
		return
	}
	item.SugImg = location
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ID), "items_id", item.ID); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ID), "name", item.Name); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ID), "brief", item.Brief); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ID), "head", item.HeadImg); err != nil {
		return
	}
	if err = conn.Send("HSET", fmt.Sprintf(_sugItem, item.ID), "pic", item.SugImg); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", fmt.Sprintf(_sugItem, item.ID), _expire); err != nil {
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

// SetSug ...
func (d *Dao) SetSug(c context.Context, seasonID string, ItemID int64, score float64) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZADD", fmt.Sprintf(_sugSeason, seasonID), score, ItemID); err != nil {
		return
	}
	if err = conn.Send("EXPIRE", fmt.Sprintf(_sugSeason, seasonID), _expire); err != nil {
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

// DelSug ...
func (d *Dao) DelSug(c context.Context, seasonID string) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("DEL", fmt.Sprintf(_sugSeason, seasonID)); err != nil {
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

// DelSugItem ...
func (d *Dao) DelSugItem(c context.Context, seasonID, itemsID int64) (b bool, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	if err = conn.Send("ZREM", fmt.Sprintf(_sugSeason, strconv.FormatInt(seasonID, 10)), itemsID); err != nil {
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
