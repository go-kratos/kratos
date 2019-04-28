package archive

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// View3 view archive with pages pb.
func (d *Dao) View3(c context.Context, aid int64) (v *archive.View3, err error) {
	arg := &archive.ArgAid2{Aid: aid}
	if v, err = d.arcRPC.View3(c, arg); err != nil {
		log.Error("d.arcRPC.View3(%v) error(%+v)", arg, err)
		if ecode.Cause(err) == ecode.NothingFound {
			err = nil
			return
		}
	}
	return
}

// ViewCache get view static data from cache if cache missed from rpc.
func (d *Dao) ViewCache(c context.Context, aid int64) (vs *archive.View3, err error) {
	if aid == 0 {
		return
	}
	if vs, err = d.viewCache(c, aid); err != nil {
		return
	}
	if vs != nil && vs.Archive3 != nil && len(vs.Pages) != 0 {
		var st *api.Stat
		if st, err = d.Stat(c, aid); err != nil {
			log.Error("%+v", err)
			err = nil
			return
		}
		if st != nil {
			vs.Archive3.Stat = archive.Stat3{
				Aid:     st.Aid,
				View:    st.View,
				Danmaku: st.Danmaku,
				Reply:   st.Reply,
				Fav:     st.Fav,
				Coin:    st.Coin,
				Share:   st.Share,
				NowRank: st.NowRank,
				HisRank: st.HisRank,
				Like:    st.Like,
				DisLike: st.DisLike,
			}
		}
	}
	return
}

// Description get archive description by aid.
func (d *Dao) Description(c context.Context, aid int64) (desc string, err error) {
	arg := &archive.ArgAid{Aid: aid}
	if desc, err = d.arcRPC.Description2(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
