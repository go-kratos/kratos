package polaris

import (
	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"
)

type Polaris struct {
	router    polaris.RouterAPI
	config    polaris.ConfigAPI
	limit     polaris.LimitAPI
	registry  polaris.ProviderAPI
	discovery polaris.ConsumerAPI
}

// New polaris Service governance.
func New(sdk api.SDKContext) Polaris {
	return Polaris{
		router:    polaris.NewRouterAPIByContext(sdk),
		config:    polaris.NewConfigAPIByContext(sdk),
		limit:     polaris.NewLimitAPIByContext(sdk),
		registry:  polaris.NewProviderAPIByContext(sdk),
		discovery: polaris.NewConsumerAPIByContext(sdk),
	}
}
