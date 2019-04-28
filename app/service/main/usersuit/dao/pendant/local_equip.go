package pendant

import (
	"context"
	"time"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"github.com/bluele/gcache"
	"github.com/pkg/errors"
)

func (d *Dao) loadEquip(c context.Context, mid int64) (*model.PendantEquip, error) {
	equip, err := d.equipCache(c, mid)
	if err != nil {
		return nil, err
	}
	d.storeEquip(mid, equip)
	return equip, nil
}

func (d *Dao) storeEquip(mid int64, equip *model.PendantEquip) {
	if equip == nil {
		return
	}
	d.equipStore.SetWithExpire(mid, equip, time.Duration(d.c.EquipCache.Expire))
}

func (d *Dao) localEquip(mid int64) (*model.PendantEquip, error) {
	prom.CacheHit.Incr("local_equip_cache")
	item, err := d.equipStore.Get(mid)
	if err != nil {
		prom.CacheMiss.Incr("local_equip_cache")
		return nil, err
	}
	equip, ok := item.(*model.PendantEquip)
	if !ok {
		prom.CacheMiss.Incr("local_equip_cache")
		return nil, errors.New("Not a equip")
	}
	return equip, nil
}

// EquipCache get equip cache.
func (d *Dao) EquipCache(c context.Context, mid int64) (*model.PendantEquip, error) {
	equip, err := d.localEquip(mid)
	if err != nil {
		if err != gcache.KeyNotFoundError {
			log.Error("Failed to get equip from local: mid: %d: %+v", mid, err)
		}
		return d.loadEquip(c, mid)
	}
	return equip, nil
}

// EquipsCache get multi equip cache.
func (d *Dao) EquipsCache(c context.Context, mids []int64) (map[int64]*model.PendantEquip, []int64, error) {
	equips := make(map[int64]*model.PendantEquip, len(mids))
	lcMissed := make([]int64, 0, len(mids))
	for _, mid := range mids {
		equip, err := d.localEquip(mid)
		if err != nil {
			if err != gcache.KeyNotFoundError {
				log.Error("Failed to get equip from local: mid: %d: %+v", mid, err)
			}
			lcMissed = append(lcMissed, mid)
			continue
		}
		equips[equip.Mid] = equip
	}
	if len(lcMissed) == 0 {
		return equips, nil, nil
	}
	rdsEquips, rdsMissed, err := d.equipsCache(c, lcMissed)
	if err != nil {
		return nil, nil, err
	}
	for _, equip := range rdsEquips {
		d.storeEquip(equip.Mid, equip)
		equips[equip.Mid] = equip
	}
	return equips, rdsMissed, nil
}
