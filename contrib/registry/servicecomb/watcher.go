package servicecomb

import (
	"github.com/go-chassis/sc-client"

	"github.com/go-kratos/kratos/v2/registry"
	"golang.org/x/net/context"
)

var _ registry.Watcher = (*Watcher)(nil)

type Watcher struct {
	cli RegistryClient
	ch  chan *registry.ServiceInstance
}

func newWatcher(_ context.Context, cli RegistryClient, serviceName string) (*Watcher, error) {
	// 构建当前服务与目标服务之间的依赖关系，完成discovery
	_, err := cli.FindMicroServiceInstances(curServiceID, appID, serviceName, "")
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		cli: cli,
		ch:  make(chan *registry.ServiceInstance),
	}
	go func() {
		watchErr := w.cli.WatchMicroService(curServiceID, func(event *sc.MicroServiceInstanceChangedEvent) {
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
			w.Put(svcIns)
		})
		if watchErr != nil {
			return
		}
	}()
	return w, nil
}

// Put only for UT
func (w *Watcher) Put(svcIns *registry.ServiceInstance) {
	w.ch <- svcIns
}

func (w *Watcher) Next() ([]*registry.ServiceInstance, error) {
	var svcInstances []*registry.ServiceInstance
	svcIns := <-w.ch
	svcInstances = append(svcInstances, svcIns)
	return svcInstances, nil
}

func (w *Watcher) Stop() error {
	close(w.ch)
	return nil
}
