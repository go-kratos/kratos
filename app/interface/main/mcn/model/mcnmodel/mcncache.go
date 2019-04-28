package mcnmodel

import "time"

// PublicationPriceCache .
type PublicationPriceCache struct {
	ModifyTime time.Time `json:"mtime"`
}

// UpPermissionCache .
type UpPermissionCache struct {
	IsNull     bool   `json:"is_null"` // 缓存用来标记
	Permission uint32 `json:"permission"`
}
