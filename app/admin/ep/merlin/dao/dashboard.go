package dao

import (
	"context"

	"go-common/app/admin/ep/merlin/model"
	"go-common/library/cache/memcache"
)

//QueryMachineUsageSummaryFromCache Machine Usage Summary In Cache.
func (d *Dao) QueryMachineUsageSummaryFromCache(c context.Context, pqadmrs []*model.PaasQueryAndDelMachineRequest) (pmds []*model.PaasMachineDetail, err error) {
	var (
		conn = d.mc.Get(c)
		item *memcache.Item
	)

	defer conn.Close()

	for _, pqadmr := range pqadmrs {

		pmd := &model.PaasMachineDetail{}

		if item, err = conn.Get(pqadmr.Name); err == nil {
			if err = conn.Scan(item, &pmd); err == nil {
				pmds = append(pmds, pmd)
				continue
			}
		}

		if pmd, err = d.QueryPaasMachine(c, pqadmr); err != nil {
			continue
		}
		pmds = append(pmds, pmd)

		item = &memcache.Item{Key: pqadmr.Name, Object: pmd, Flags: memcache.FlagJSON, Expiration: d.expire}

		d.tokenCacheSave(c, item)

	}
	return
}
