package dao

import (
	"context"
	"fmt"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
)

const (
	_replyStatFormat = "rs_%d"
)

// PingMc ping
func (d *Dao) PingMc(ctx context.Context) (err error) {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	item := memcache.Item{Key: "ping", Value: []byte("pong"), Expiration: 0}
	return conn.Set(&item)
}

func keyReplyStat(rpID int64) string {
	return fmt.Sprintf(_replyStatFormat, rpID)
}

// RemReplyStatMc ...
func (d *Dao) RemReplyStatMc(ctx context.Context, rpID int64) (err error) {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	return conn.Delete(keyReplyStat(rpID))
}

// SetReplyStatMc set reply stat into mc.
func (d *Dao) SetReplyStatMc(ctx context.Context, rs *model.ReplyStat) (err error) {
	conn := d.mc.Get(ctx)
	defer conn.Close()
	key := keyReplyStat(rs.RpID)
	item := &memcache.Item{
		Key:        key,
		Object:     rs,
		Expiration: d.mcExpire,
		Flags:      memcache.FlagJSON,
	}
	if err = conn.Set(item); err != nil {
		log.Error("memcache Set(%s, %v), error(%v)", key, item, err)
	}
	return
}

// ReplyStatsMc get multi repies stat from memcache.
func (d *Dao) ReplyStatsMc(ctx context.Context, rpIDs []int64) (rsMap map[int64]*model.ReplyStat, missIDs []int64, err error) {
	rsMap = make(map[int64]*model.ReplyStat)
	keys := make([]string, len(rpIDs))
	mapping := make(map[string]int64)
	for i, rpID := range rpIDs {
		key := keyReplyStat(rpID)
		keys[i] = key
		mapping[key] = rpID
	}
	for _, chunkedKeys := range splitString(keys, 2000) {
		var (
			conn  = d.mc.Get(ctx)
			items map[string]*memcache.Item
		)
		if items, err = conn.GetMulti(chunkedKeys); err != nil {
			if err == memcache.ErrNotFound {
				missIDs = rpIDs
				err = nil
				conn.Close()
				return
			}
			conn.Close()
			log.Error("memcache GetMulti error(%v)", err)
			return
		}
		for _, item := range items {
			stat := new(model.ReplyStat)
			if err = conn.Scan(item, stat); err != nil {
				log.Error("memcache Scan(%v) error(%v)", item.Value, err)
				continue
			}
			rsMap[mapping[item.Key]] = stat
			delete(mapping, item.Key)
		}
		conn.Close()
	}
	for _, rpID := range mapping {
		missIDs = append(missIDs, rpID)
	}
	return
}
