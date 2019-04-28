package anchorReward

import (
	"context"
	"fmt"

	"go-common/library/cache/memcache"
	"go-common/library/log"

	"go-common/app/service/live/xrewardcenter/model/anchorTask"

	"github.com/pkg/errors"
)

const (
	//NewRewardCache have new reward.
	NewRewardCache = "new_reward:v1:%d"
	//HasRewardCache have a reward
	HasRewardCache = "has_reward:v1:%d"
)

// HasNewReward .
type HasNewReward struct {
	HasNewReward int64 `json:"has_new_reward"`
}

// pingMC .
func (d *Dao) pingMC(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err = conn.Set(&memcache.Item{Key: "ping", Value: []byte{1}, Expiration: 0}); err != nil {
		err = errors.Wrap(err, "conn.Store(set,ping,1)")
	}
	return
}

func newRewardKey(uid int64) string {
	return fmt.Sprintf(NewRewardCache, uid)
}

// SetNewReward get MRRecord from mc
func (d *Dao) SetNewReward(c context.Context, uid int64, hasNew int64) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = newRewardKey(uid)
	)

	defer conn.Close()

	log.Info("SetNewReward.mc.set(%v) (%v)", key, hasNew)

	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     hasNew,
		Flags:      memcache.FlagJSON,
		Expiration: 3600,
	}); err != nil {
		log.Error("SetNewReward.mc.set(%v %v) error(%v)", key, hasNew, err)
		return
	}
	return
}

// GetNewReward set MRRecord to mc
func (d *Dao) GetNewReward(c context.Context, uid int64) (r int64, err error) {
	r = 0

	var (
		key  = newRewardKey(uid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()

	item, err := conn.Get(key)

	if err != nil {
		if err == memcache.ErrNotFound {
			return r, err
		}
		err = errors.Wrapf(err, "conn.get(%v,%v)", key, item)
		return r, err
	}

	if err = conn.Scan(item, &r); err != nil {
		log.Error("mc.Scan(%s) error(%v)", string(item.Value), err)
	}

	return
}

// GetHasReward .
func (d *Dao) GetHasReward(c context.Context, uid int64) (r int64, err error) {
	var (
		key  = hasRewardKey(uid)
		conn = d.mc.Get(c)
	)
	defer conn.Close()

	r = 0
	item, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			return r, err
		}
		err = errors.Wrapf(err, "mc.get(%s,%v)", key, item)
		return r, err
	}

	if err = conn.Scan(item, &r); err != nil {
		log.Error("mc.Scan(%s) error(%v)", string(item.Value), err)
	}
	//spew.Dump(r, err)

	return
}

func hasRewardKey(uid int64) string {
	return fmt.Sprintf(HasRewardCache, uid)
}

// SetHasReward save whether one user have reward.
func (d *Dao) SetHasReward(c context.Context, uid int64, isHave int64) (err error) {
	var (
		conn = d.mc.Get(c)
		key  = hasRewardKey(uid)
	)

	defer conn.Close()

	log.Info("SetHasReward.mc.set(%v %v)", key, isHave)

	if err = conn.Set(&memcache.Item{
		Key:        key,
		Object:     isHave,
		Flags:      memcache.FlagJSON,
		Expiration: 3600 * 24,
	}); err != nil {
		log.Error("SetHasReward.mc.set(%v %v) error(%v)", key, isHave, err)
		return
	}
	return
}

// GetOrder get order by id from cache.
func (d *Dao) GetOrder(c context.Context, id string) (exists int64, err error) {

	conn := d.mc.Get(c)
	defer conn.Close()

	_, err = conn.Get(id)

	if err != nil {
		exists = 0
		return exists, err
	}
	exists = anchorTask.RewardExists
	return exists, err
}

// SaveOrder save orderid in cache.
func (d *Dao) SaveOrder(c context.Context, id string) (err error) {

	conn := d.mc.Get(c)
	defer conn.Close()

	if err = conn.Set(&memcache.Item{
		Key:        id,
		Value:      []byte{1},
		Expiration: 3600,
	}); err != nil {
		log.Error("(%v) error(%v)", id, err)
		return
	}

	return
}
