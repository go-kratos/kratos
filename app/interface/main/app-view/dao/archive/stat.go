package archive

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// Stat get a archive stat.
func (d *Dao) Stat(c context.Context, aid int64) (st *api.Stat, err error) {
	if st, err = d.statCache(c, aid); err != nil {
		log.Error("%+v", err)
	} else if st != nil {
		return
	}
	arg := &archive.ArgAid2{Aid: aid}
	if st, err = d.arcRPC.Stat3(c, arg); err != nil {
		log.Error("d.arcRPC.Stat3(%v) error(%v)", arg, err)
	}
	return
}
