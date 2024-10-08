package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-kratos/kratos/v2/registry"
)

var _ registry.Watcher = (*watcher)(nil)

type watcher struct {
	serviceName    string
	clusters       []string
	groupName      string
	ctx            context.Context
	cancel         context.CancelFunc
	watchChan      chan struct{}
	cli            naming_client.INamingClient
	kind           string
	subscribeParam *vo.SubscribeParam
}

func newWatcher(ctx context.Context, cli naming_client.INamingClient, serviceName, groupName, kind string, clusters []string) (*watcher, error) {
	w := &watcher{
		serviceName: serviceName,
		clusters:    clusters,
		groupName:   groupName,
		cli:         cli,
		kind:        kind,
		watchChan:   make(chan struct{}, 1),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)

	w.subscribeParam = &vo.SubscribeParam{
		ServiceName: serviceName,
		Clusters:    clusters,
		GroupName:   groupName,
		SubscribeCallback: func([]model.SubscribeService, error) {
			select {
			case w.watchChan <- struct{}{}:
			default:
			}
		},
	}
	e := w.cli.Subscribe(w.subscribeParam)
	select {
	case w.watchChan <- struct{}{}:
	default:
	}
	return w, e
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
	}
	res, err := w.cli.GetService(vo.GetServiceParam{
		ServiceName: w.serviceName,
		GroupName:   w.groupName,
		Clusters:    w.clusters,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*registry.ServiceInstance, 0, len(res.Hosts))
	for _, in := range res.Hosts {
		kind := w.kind
		if k, ok := in.Metadata["kind"]; ok {
			kind = k
		}
		items = append(items, &registry.ServiceInstance{
			ID:        in.InstanceId,
			Name:      res.Name,
			Version:   in.Metadata["version"],
			Metadata:  in.Metadata,
			Endpoints: []string{fmt.Sprintf("%s://%s:%d", kind, in.Ip, in.Port)},
		})
	}
	return items, nil
}

func (w *watcher) Stop() error {
	err := w.cli.Unsubscribe(w.subscribeParam)
	w.cancel()
	return err
}
