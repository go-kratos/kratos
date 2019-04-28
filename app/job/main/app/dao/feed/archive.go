package feed

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Archives multi get archives.
func (d *Dao) Archives(c context.Context, aids []int64, ip string) (as map[int64]*api.Arc, err error) {
	if len(aids) == 0 {
		return
	}
	var (
		missed []int64
		missAv map[int64]*api.Arc
	)
	// get data from memecached
	if as, missed, err = d.arcsCache(c, aids); err != nil {
		log.Error("%+v", err)
	}
	if len(as) == 0 {
		as = make(map[int64]*api.Arc, len(aids))
		missed = aids
	}
	if len(missed) != 0 {
		arg := &archive.ArgAids2{Aids: missed, RealIP: ip}
		if missAv, err = d.arcRPC.Archives3(c, arg); err != nil {
			err = errors.Wrapf(err, "%v", arg)
			return
		}
		for aid, a := range missAv {
			as[aid] = a
		}
	}
	return
}
