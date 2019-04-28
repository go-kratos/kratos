package dao

import (
	"context"

	"go-common/app/service/main/share/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Shares get shares
func (d *Dao) Shares(ctx context.Context, oids []int64, tp int) (shares map[int64]int64, err error) {
	shares, err = d.SharesCache(ctx, oids, tp)
	if err != nil {
		log.Error("d.SharesCache(%v) tp(%d) error(%v)", oids, tp, err)
		err = nil
		shares = make(map[int64]int64, len(oids))
	}
	var missed []int64
	for _, oid := range oids {
		if _, ok := shares[oid]; !ok {
			missed = append(missed, oid)
		}
	}
	if len(missed) == 0 {
		return
	}
	// 最大30个id，并且分了100张表，用in的优化空间也不大，暂时循环单个查
	for _, oid := range missed {
		cnt, err := d.ShareCount(ctx, oid, tp)
		if err != nil {
			continue
		}
		shares[oid] = cnt
	}
	return
}

// ShareCount get share from cache/db
func (d *Dao) ShareCount(ctx context.Context, oid int64, tp int) (count int64, err error) {
	count, err = d.ShareCache(ctx, oid, tp)
	if count != -1 && err == nil {
		return
	}
	var share *model.Share
	if share, err = d.Share(ctx, oid, tp); err != nil {
		err = errors.WithStack(err)
		return
	}
	count = 0
	if share != nil {
		count = share.Count
	}
	d.asyncCache.Save(func() {
		if err = d.SetShareCache(context.Background(), oid, tp, count); err != nil {
			log.Error("%+v", err)
			return
		}
	})
	return
}

// Add add share
func (d *Dao) Add(ctx context.Context, p *model.ShareParams) (shared int64, err error) {
	var ok bool
	if ok, err = d.AddShareMember(ctx, p); err != nil {
		return
	}
	if !ok {
		err = ecode.ShareAlreadyAdd
		return
	}
	if err = d.AddShare(ctx, p.OID, p.TP); err != nil {
		err = errors.WithStack(err)
		return
	}
	var share *model.Share
	if share, err = d.Share(ctx, p.OID, p.TP); err != nil {
		err = errors.WithStack(err)
		return
	}
	shared = share.Count
	d.asyncCache.Save(func() {
		if err = d.SetShareCache(context.Background(), p.OID, p.TP, shared); err != nil {
			log.Error("%+v", err)
			return
		}
	})
	return
}
