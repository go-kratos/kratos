package service

import (
	"context"

	"go-common/app/service/main/ugcpay-rank/internal/service/rank"
	"go-common/library/log"
	"go-common/library/stat/prom"

	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

var (
	sfRank = [1]*singleflight.Group{{}}
	_ctx   = context.Background()
)

func (s *Service) rank(ctx context.Context, sfKey string, r rank.Rank, pr rank.PrepRank) (data interface{}, err error) {
	var (
		theRank     interface{}
		thePrepRank interface{}
	)
	// load rank from cache
	if theRank, err = r.Load(ctx); err != nil {
		log.Error("%s, err: %+v", r, err)
		err = nil
	}
	// if rank exist return
	if theRank != nil {
		prom.CacheHit.Incr("Rank")
		data = theRank
		return
	}
	// sf 闭包
	fn := func() (res interface{}, ferr error) {
		prom.CacheMiss.Incr("Rank")
		if thePrepRank, ferr = s.prepRank(ctx, pr); ferr != nil {
			log.Error("loadRank get prep_rank failed, err: %+v", ferr)
			return
		}
		if thePrepRank == nil {
			return
		}
		// rebuild rank
		if theRank, ferr = r.Rebuild(ctx, thePrepRank); ferr != nil {
			ferr = errors.WithMessage(ferr, "loadRank rebuild rank failed")
			return
		}
		if theRank == nil {
			err = errors.Errorf("loadRank rebuild rank failed, nil rank returned, r: %s, pr: %s", r, pr)
			return
		}
		// save rank to cache
		s.Asyncer.Do(_ctx, func(ctx context.Context) {
			if theErr := r.Save(ctx, theRank); theErr != nil {
				log.Error("loadRank save rank failed, err: %+v", theErr)
				return
			}
		})
		res = theRank
		return
	}
	data, err, _ = sfRank[0].Do(sfKey, fn)
	return
}

// UpdateElecPrepRankFromOrder .
func (s *Service) UpdateElecPrepRankFromOrder(ctx context.Context, pr rank.PrepRank, payMID int64, fee int64) (err error) {
	var (
		thePrepRank     interface{}
		updatedPrepRank interface{}
	)
	if thePrepRank, err = s.prepRank(ctx, pr); err != nil {
		err = errors.WithMessage(err, "updatePrepRank load prep_rank failed")
		return
	}
	if thePrepRank == nil {
		return
	}
	if updatedPrepRank, err = pr.UpdateOrder(ctx, thePrepRank, payMID, fee); err != nil {
		err = errors.WithMessage(err, "updatePrepRank update prep_rank failed")
		return
	}
	if err = pr.Save(ctx, updatedPrepRank); err != nil {
		err = errors.WithMessage(err, "updatePrepRank save prep_rank failed")
		return
	}
	return
}

// UpdateElecPrepRankFromMessage .
func (s *Service) UpdateElecPrepRankFromMessage(ctx context.Context, pr rank.PrepRank, payMID int64, message string, hidden bool) (err error) {
	var (
		thePrepRank     interface{}
		updatedPrepRank interface{}
	)
	if thePrepRank, err = s.prepRank(ctx, pr); err != nil {
		err = errors.WithMessage(err, "updatePrepRank load prep_rank failed")
		return
	}
	if thePrepRank == nil {
		return
	}
	if updatedPrepRank, err = pr.UpdateMessage(ctx, thePrepRank, payMID, message, hidden); err != nil {
		err = errors.WithMessage(err, "updatePrepRank update prep_rank failed")
		return
	}
	if err = pr.Save(ctx, updatedPrepRank); err != nil {
		err = errors.WithMessage(err, "updatePrepRank save prep_rank failed")
		return
	}
	return
}

// prepRank 一定返回非nil rank, 或非nil err
func (s *Service) prepRank(ctx context.Context, pr rank.PrepRank) (rank interface{}, err error) {
	var (
		thePrepRank interface{}
	)
	// load prep_rank from cache
	if thePrepRank, err = pr.Load(ctx); err != nil {
		log.Error("%s, err: %+v", pr, err)
		err = nil
	}
	if thePrepRank != nil {
		prom.CacheHit.Incr("PrepRank")
		rank = thePrepRank
		return
	}
	// if prep_rank = nil, rebuild prep_rank from db
	prom.CacheMiss.Incr("PrepRank")
	if thePrepRank, err = pr.Rebuild(ctx); err != nil {
		err = errors.WithMessage(err, "loadRank rebuild prep_rank failed")
		return
	}
	// if still prep_rank = nil, return err
	if thePrepRank == nil {
		err = errors.Errorf("loadRank rebuild prep_rank failed, nil prep_rank returned, pr: %s", pr)
		return
	}
	rank = thePrepRank
	// save prep_rank to cache
	s.Asyncer.Do(_ctx, func(ctx context.Context) {
		if theErr := pr.Save(ctx, thePrepRank); theErr != nil {
			log.Error("loadRank save prep_rank failed, err: %+v", theErr)
			return
		}
	})
	return
}
