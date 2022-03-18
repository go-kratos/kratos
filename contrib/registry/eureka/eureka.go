package eureka

import (
	"context"
	"strings"
	"sync"
	"time"
)

type subscriber struct {
	appID    string
	callBack func()
}

type API struct {
	cli             *Client
	allInstances    map[string][]Instance
	subscribers     map[string]*subscriber
	refreshInterval time.Duration
	lock            sync.Mutex
}

func NewAPI(ctx context.Context, client *Client, refreshInterval time.Duration) *API {
	e := &API{
		cli:             client,
		allInstances:    make(map[string][]Instance),
		subscribers:     make(map[string]*subscriber),
		refreshInterval: refreshInterval,
	}

	// 首次广播一次
	e.broadcast()

	go e.refresh(ctx)

	return e
}

func (e *API) refresh(ctx context.Context) {
	ticker := time.NewTicker(e.refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.broadcast()
		}
	}
}

func (e *API) broadcast() {
	instances := e.cacheAllInstances()
	if instances == nil {
		return
	}

	for _, subscriber := range e.subscribers {
		go subscriber.callBack()
	}
	defer e.lock.Unlock()
	e.lock.Lock()
	e.allInstances = instances
}

func (e *API) cacheAllInstances() map[string][]Instance {
	items := make(map[string][]Instance)
	instances := e.cli.FetchAllUpInstances(context.Background())
	for _, instance := range instances {
		items[e.ToAppID(instance.App)] = append(items[instance.App], instance)
	}

	return items
}

func (e *API) Register(ctx context.Context, serviceName string, endpoints ...Endpoint) error {
	appID := e.ToAppID(serviceName)
	upInstances := make(map[string]struct{})

	for _, ins := range e.GetService(ctx, appID) {
		upInstances[ins.InstanceID] = struct{}{}
	}

	for _, ep := range endpoints {
		if _, ok := upInstances[ep.InstanceID]; !ok {
			if err := e.cli.Register(ctx, ep); err != nil {
				return err
			}
			go e.cli.Heartbeat(ep)
		}
	}

	return nil
}

// Deregister 中的ctx 和 register ctx 是同一个
func (e *API) Deregister(ctx context.Context, endpoints []Endpoint) error {
	for _, ep := range endpoints {
		if err := e.cli.Deregister(ctx, ep.AppID, ep.InstanceID); err != nil {
			return err
		}
	}

	return nil
}

func (e *API) Subscribe(serverName string, fn func()) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	appID := e.ToAppID(serverName)
	e.subscribers[appID] = &subscriber{
		appID:    appID,
		callBack: fn,
	}
	return nil
}

func (e *API) GetService(ctx context.Context, serverName string) []Instance {
	appID := e.ToAppID(serverName)
	if ins, ok := e.allInstances[appID]; ok {
		return ins
	}

	// 如果不再allinstances 中可以尝试再单独获取一次
	return e.cli.FetchAppUpInstances(ctx, appID)
}

func (e *API) Unsubscribe(serverName string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.subscribers, e.ToAppID(serverName))
}

func (e *API) ToAppID(serverName string) string {
	return strings.ToUpper(serverName)
}
