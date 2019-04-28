package service

import (
	"context"
	"math"
	"time"

	"go-common/app/service/main/ugcpay/conf"
	"go-common/app/service/main/ugcpay/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AssetRegister register a content to asset
func (s *Service) AssetRegister(ctx context.Context, mid int64, oid int64, otype string, currency string, price int64) (err error) {
	var (
		asset = &model.Asset{
			ID:       1,
			MID:      mid,
			OID:      oid,
			OType:    otype,
			Currency: currency,
			Price:    price,
			State:    "valid",
			CTime:    time.Now(),
			MTime:    time.Now(),
		}
	)
	if err = s.dao.UpsertAsset(ctx, asset); err != nil {
		return
	}
	if err = s.dao.DelCacheAsset(ctx, oid, otype, currency); err != nil {
		return
	}
	return
}

// AssetQuery get asset by content
func (s *Service) AssetQuery(ctx context.Context, oid int64, otype string, currency string) (as *model.Asset, pp map[string]int64, err error) {
	log.Info("asset query oid : %d, otype : %s, currency : %s", oid, otype, currency)
	if as, err = s.dao.Asset(ctx, oid, otype, currency); err != nil {
		return
	}
	if as == nil {
		err = ecode.UGCPayAssetInvalid
		return
	}
	pp = s.calcPlatformPrice(ctx, as.Price)
	return
}

func (s *Service) calcPlatformPrice(ctx context.Context, price int64) (pp map[string]int64) {
	pp = make(map[string]int64)
	for platform, rate := range conf.Conf.Biz.Price.PlatformTax {
		if rate > 1.0 {
			pp[platform] = int64(math.Ceil(float64(price) * rate))
		}
	}
	return
}

// func (s *Service) calcRawPrice(ctx context.Context, platform string, fee int64) (rawFee int64) {
// 	tax, ok := conf.Conf.Biz.Price.PlatformTax[platform]
// 	if !ok {
// 		rawFee = fee
// 		return
// 	}
// 	rawFee = int64(float64(fee) / tax)
// 	return
// }

// AssetRelation get relation
func (s *Service) AssetRelation(ctx context.Context, mid int64, oid int64, otype string) (state string, err error) {
	if mid <= 0 {
		return "none", nil
	}
	var (
		assetRelation *model.AssetRelation
	)
	if state, err = s.dao.CacheAssetRelationState(ctx, oid, otype, mid); err != nil {
		return
	}
	if state != "miss" {
		return
	}
	if assetRelation, err = s.dao.RawAssetRelation(ctx, mid, oid, otype); err != nil {
		return
	}
	if assetRelation == nil {
		state = "none"
	} else {
		state = assetRelation.State
	}
	s.cache.Save(func() {
		if theErr := s.dao.AddCacheAssetRelationState(context.Background(), oid, otype, mid, state); theErr != nil {
			log.Error("s.dao.AddCacheAssetRelationState error : %+v", theErr)
			return
		}
	})
	return
}
