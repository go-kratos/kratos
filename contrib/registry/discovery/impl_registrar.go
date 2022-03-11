package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/SeeMusic/kratos/v2/registry"
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
		defer d.Logger().Warn("Discovery:register_service goroutine quit")
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

//  register an instance with Discovery
func (d *Discovery) register(ctx context.Context, ins *discoveryInstance) (err error) {
	d.mutex.RLock()
	c := d.config
	d.mutex.RUnlock()

	var metadata []byte
	if ins.Metadata != nil {
		if metadata, err = json.Marshal(ins.Metadata); err != nil {
			d.Logger().Errorf(
				"Discovery:register instance Marshal metadata(%v) failed!error(%v)", ins.Metadata, err,
			)
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
		d.Logger().Errorf("Discovery: register client.Get(%s)  zone(%s) env(%s) appid(%s) addrs(%v) error(%v)",
			uri+"?"+p.Encode(), c.Zone, c.Env, ins.AppID, ins.Addrs, err)
		return
	}

	if res.Code != 0 {
		err = fmt.Errorf("ErrorCode: %d", res.Code)
		d.Logger().Errorf("Discovery: register client.Get(%v)  env(%s) appid(%s) addrs(%v) code(%v)",
			uri, c.Env, ins.AppID, ins.Addrs, res.Code)
	}

	d.Logger().Infof(
		"Discovery: register client.Get(%v) env(%s) appid(%s) addrs(%s) success\n",
		uri, c.Env, ins.AppID, ins.Addrs,
	)

	return
}

func (d *Discovery) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	ins := fromServerInstance(service, d.config)
	return d.cancel(ins)
}
