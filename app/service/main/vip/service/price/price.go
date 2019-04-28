package price

import (
	"context"
	"sync"

	"go-common/app/service/main/vip/model"
)

// Price .
type Price struct {
	priceConfig map[int64]map[int8][]*model.VipPriceConfig
	sync.RWMutex
}

// New create a price map.
func New() (p *Price) {
	return &Price{
		priceConfig: make(map[int64]map[int8][]*model.VipPriceConfig),
	}
}

// GetPriceConfig .
func (p *Price) GetPriceConfig(m int64) (res map[int8][]*model.VipPriceConfig, ok bool) {
	p.RLock()
	defer p.RUnlock()
	res, ok = p.priceConfig[m]
	return
}

// SetPriceConfig .
func (p *Price) SetPriceConfig(c context.Context, priceConfig map[int64]map[int8][]*model.VipPriceConfig) (err error) {
	p.Lock()
	p.priceConfig = priceConfig
	p.Unlock()
	return
}
