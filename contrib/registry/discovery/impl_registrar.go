package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/registry"
)

func (d *Discovery) Register(ctx context.Context, service *registry.ServiceInstance) (err error) {
	ins := fromServerInstance(service, d.config)

	d.mutex.Lock()
	if _, ok := d.registry[ins.AppID]; ok {
		err = errors.Wrap(ErrDuplication, ins.AppID)
	} else {
		d.registry[ins.AppID] = struct{}{}
	}
	d.mutex.Unlock()
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(d.ctx)
	if err = d.register(ctx, ins); err != nil {
		d.mutex.Lock()
		delete(d.registry, ins.AppID)
		d.mutex.Unlock()
		cancel()
		return
	}

	ch := make(chan struct{}, 1)
	d.cancelFunc = func() {
		cancel()
		<-ch
	}

	// renew the current register_service
	go func() {
		defer log.Warn("Discovery:register_service goroutine quit")
		ticker := time.NewTicker(_registerGap)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = d.renew(ctx, ins)
			case <-ctx.Done():
				_ = d.cancel(ins)
				ch <- struct{}{}
				return
			}
		}
	}()

	return
}

// register an instance with Discovery
func (d *Discovery) register(ctx context.Context, ins *discoveryInstance) (err error) {
	d.mutex.RLock()
	c := d.config
	d.mutex.RUnlock()

	var metadata []byte
	if ins.Metadata != nil {
		if metadata, err = json.Marshal(ins.Metadata); err != nil {
			log.Error("Discovery: register instance marshal metadata failed", "metadata", ins.Metadata, "error", err)
		}
	}
	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	})
	uri := fmt.Sprintf(_registerURL, d.pickNode())

	// params
	p := newParams(d.config)
	p.Set(_paramKeyAppID, ins.AppID)
	for _, addr := range ins.Addrs {
		p.Add(_paramKeyAddrs, addr)
	}
	p.Set(_paramKeyVersion, ins.Version)
	if ins.Status == 0 {
		p.Set(_paramKeyStatus, _statusUP)
	} else {
		p.Set(_paramKeyStatus, strconv.FormatInt(ins.Status, 10))
	}
	p.Set(_paramKeyMetadata, string(metadata))

	// send request to Discovery server.
	if _, err = d.httpClient.R().
		SetContext(ctx).
		SetQueryParamsFromValues(p).
		SetResult(&res).
		Post(uri); err != nil {
		d.switchNode()
		log.Error("Discovery: register client.Get failed",
			"uri", uri+"?"+p.Encode(), "zone", c.Zone, "env", c.Env, "appid", ins.AppID, "addrs", ins.Addrs, "error", err)
		return
	}

	if res.Code != 0 {
		err = fmt.Errorf("ErrorCode: %d", res.Code)
		log.Error("Discovery: register client.Get returned code",
			"uri", uri, "env", c.Env, "appid", ins.AppID, "addrs", ins.Addrs, "code", res.Code)
	}

	log.Info("Discovery: register client.Get succeeded", "uri", uri, "env", c.Env, "appid", ins.AppID, "addrs", ins.Addrs)

	return
}

func (d *Discovery) Deregister(_ context.Context, service *registry.ServiceInstance) error {
	ins := fromServerInstance(service, d.config)
	return d.cancel(ins)
}
