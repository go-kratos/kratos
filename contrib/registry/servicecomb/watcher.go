package servicecomb

import (
	"github.com/go-chassis/sc-client"
	"github.com/go-kratos/kratos/v2/registry"
	"golang.org/x/net/context"
)

var _ registry.Watcher = (*watcher)(nil)

type watcher struct {
	cli *sc.Client
	ch  chan *registry.ServiceInstance
}

func newWatcher(_ context.Context, cli *sc.Client, serviceName string) (*watcher, error) {
	//构建当前服务与目标服务之间的依赖关系，完成discovery
	_, err := cli.FindMicroServiceInstances(curServiceId, appId, serviceName, "")
	if err != nil {
		return nil, err
	}
	w := &watcher{
		cli: cli,
		ch:  make(chan *registry.ServiceInstance),
	}
	go func() {
		watchErr := w.cli.WatchMicroService(curServiceId, func(event *sc.MicroServiceInstanceChangedEvent) {
			if event.Key.ServiceName != serviceName {
				return
			}
			svcIns := &registry.ServiceInstance{
				ID:        event.Instance.InstanceId,
				Name:      event.Key.ServiceName,
				Version:   event.Key.Version,
				Metadata:  event.Instance.Properties,
				Endpoints: event.Instance.Endpoints,
			}
			w.ch <- svcIns
		})
		if watchErr != nil {
			return
		}
	}()
	return w, nil
}

func (w watcher) Next() ([]*registry.ServiceInstance, error) {
	var svcInstances []*registry.ServiceInstance
	svcIns := <-w.ch
	svcInstances = append(svcInstances, svcIns)
	return svcInstances, nil
}

func (w watcher) Stop() error {
	close(w.ch)
	return nil
}
