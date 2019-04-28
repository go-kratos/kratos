package cms

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// MixedFilter filters ugc and pgc data to get the allowed data
func (d *Dao) MixedFilter(ctx context.Context, sids []int64, aids []int64) (okSids map[int64]int, okAids map[int64]int) {
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		okAids = d.aidsFilter(context.Background(), aids)
		return
	})
	g.Go(func() (err error) {
		okSids = d.sidsFilter(context.Background(), sids)
		return
	})
	g.Wait()
	return
}

// filter canPlay Aids
func (d *Dao) aidsFilter(ctx context.Context, aids []int64) (okAids map[int64]int) {
	var (
		arcMetas map[int64]*model.ArcCMS
		err      error
	)
	okAids = make(map[int64]int)
	if arcMetas, err = d.LoadArcsMediaMap(ctx, aids); err != nil {
		log.Error("MixedFilter Aids %v, Err %v", aids, err)
		return
	}
	if len(arcMetas) == 0 {
		return
	}
	for aid, arcMeta := range arcMetas {
		if arcMeta.CanPlay() {
			okAids[aid] = 1
		}
	}
	return
}

// filter canPlay Sids
func (d *Dao) sidsFilter(ctx context.Context, sids []int64) (okSids map[int64]int) {
	var (
		snsAuth map[int64]*model.SnAuth
		err     error
	)
	okSids = make(map[int64]int)
	if snsAuth, err = d.LoadSnsAuthMap(ctx, sids); err != nil {
		log.Error("MixedFilter Sids %v, Err %v", sids, err)
	}
	if len(snsAuth) == 0 {
		return
	}
	for sid, snAuth := range snsAuth {
		if snAuth.CanPlay() {
			okSids[sid] = 1
		}
	}
	return
}
