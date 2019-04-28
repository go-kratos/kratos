package http

import (
	"go-common/app/service/main/ugcpay/model"
)

// ArgAssetRegister .
type ArgAssetRegister struct {
	MID      int64  `form:"mid" validate:"required"`
	OID      int64  `form:"oid" validate:"required"`
	OType    string `form:"otype" validate:"required"`
	Price    int64  `form:"price" validate:"required"`
	Currency string `form:"currency" validate:"required"`
}

// ArgAssetQuery .
type ArgAssetQuery struct {
	OID      int64  `form:"oid" validate:"required"`
	OType    string `form:"otype" validate:"required"`
	Currency string `form:"currency" validate:"required"`
}

// RespAssetQuery .
type RespAssetQuery struct {
	Price         int64            `json:"price"`
	PlatformPrice map[string]int64 `json:"platform_price"`
}

// Parse .
func (r *RespAssetQuery) Parse(as *model.Asset, pp map[string]int64) {
	if as == nil {
		return
	}
	r.Price = as.Price
	r.PlatformPrice = pp
}

// ArgAssetRelation .
type ArgAssetRelation struct {
	MID   int64  `form:"mid" validate:"required"`
	OID   int64  `form:"oid" validate:"required"`
	OType string `form:"otype" validate:"required"`
}

// RespAssetRelation .
type RespAssetRelation struct {
	State string `json:"state"`
}

// ArgAssetRelationDetail .
type ArgAssetRelationDetail struct {
	ArgAssetRelation
	Currency string `form:"currency" validate:"required"`
}

// RespAssetRelationDetail .
type RespAssetRelationDetail struct {
	RelationState      string           `json:"relation_state"`
	AssetPrice         int64            `json:"asset_price"`
	AssetPlatformPrice map[string]int64 `json:"asset_platform_price"`
}
