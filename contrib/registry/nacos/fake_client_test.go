package nacos

import (
	"errors"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// fakeNamingClient is a lightweight in-memory fake implementation of
// naming_client.INamingClient used by tests to avoid needing a running
// Nacos server.
type fakeNamingClient struct {
	mu          sync.RWMutex
	services    map[string][]model.Instance // key: group@@serviceName
	subscribers map[string][]*vo.SubscribeParam
}

// NewFakeNamingClient creates a new fake naming client.
func NewFakeNamingClient() naming_client.INamingClient {
	return &fakeNamingClient{
		services:    make(map[string][]model.Instance),
		subscribers: make(map[string][]*vo.SubscribeParam),
	}
}

func (f *fakeNamingClient) notify(serviceKey string) {
	f.mu.RLock()
	subs := f.subscribers[serviceKey]
	f.mu.RUnlock()
	for _, sp := range subs {
		if sp.SubscribeCallback != nil {
			// make a local copy of hosts
			f.mu.RLock()
			hosts := append([]model.Instance(nil), f.services[serviceKey]...)
			f.mu.RUnlock()
			sp.SubscribeCallback(hosts, nil)
		}
	}
}

func (f *fakeNamingClient) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	if param.ServiceName == "" {
		return false, errors.New("service name empty")
	}
	key := param.GroupName + "@@" + param.ServiceName
	cluster := param.ClusterName
	if cluster == "" {
		cluster = "DEFAULT"
	}
	inst := model.Instance{
		Ip:          param.Ip,
		Port:        param.Port,
		ServiceName: key, // store with group prefix to mimic server GetService behavior
		ClusterName: cluster,
		Metadata:    param.Metadata,
	}
	// leave InstanceId empty to exercise fallback behavior
	f.mu.Lock()
	f.services[key] = append(f.services[key], inst)
	f.mu.Unlock()
	// notify subscribers
	go f.notify(key)
	return true, nil
}

func (f *fakeNamingClient) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	key := param.GroupName + "@@" + param.ServiceName
	f.mu.Lock()
	defer f.mu.Unlock()
	list := f.services[key]
	for i := 0; i < len(list); i++ {
		if list[i].Ip == param.Ip && list[i].Port == param.Port {
			// remove
			list = append(list[:i], list[i+1:]...)
			i--
		}
	}
	f.services[key] = list
	go f.notify(key)
	return true, nil
}

func (f *fakeNamingClient) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	key := param.GroupName + "@@" + param.ServiceName
	f.mu.RLock()
	defer f.mu.RUnlock()
	list := f.services[key]
	if len(list) == 0 {
		return nil, errors.New("service not found")
	}
	// return a copy
	out := make([]model.Instance, len(list))
	copy(out, list)
	return out, nil
}

func (f *fakeNamingClient) GetService(param vo.GetServiceParam) (model.Service, error) {
	key := param.GroupName + "@@" + param.ServiceName
	f.mu.RLock()
	defer f.mu.RUnlock()
	list := f.services[key]
	hosts := make([]model.Instance, len(list))
	copy(hosts, list)
	// Return empty Service with nil error when no hosts — watcher expects
	// an initial empty response instead of an error.
	return model.Service{Hosts: hosts}, nil
}

func (f *fakeNamingClient) Subscribe(param *vo.SubscribeParam) error {
	key := param.GroupName + "@@" + param.ServiceName
	f.mu.Lock()
	f.subscribers[key] = append(f.subscribers[key], param)
	f.mu.Unlock()
	// call once to prime the watcher
	go func() {
		f.mu.RLock()
		hosts := append([]model.Instance(nil), f.services[key]...)
		f.mu.RUnlock()
		if param.SubscribeCallback != nil {
			param.SubscribeCallback(hosts, nil)
		}
	}()
	return nil
}

func (f *fakeNamingClient) Unsubscribe(param *vo.SubscribeParam) error {
	key := param.GroupName + "@@" + param.ServiceName
	f.mu.Lock()
	defer f.mu.Unlock()
	subs := f.subscribers[key]
	for i, s := range subs {
		if s == param {
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}
	f.subscribers[key] = subs
	return nil
}

// The rest of the INamingClient methods are not used by tests but
// must be provided to satisfy the interface. Return sensible defaults.
func (f *fakeNamingClient) BatchRegisterInstance(_ vo.BatchRegisterInstanceParam) (bool, error) {
	return true, nil
}
func (f *fakeNamingClient) CloseClient() {}
func (f *fakeNamingClient) GetAllServicesInfo(_ vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return model.ServiceList{}, nil
}

func (f *fakeNamingClient) SelectAllInstances(_ vo.SelectAllInstancesParam) ([]model.Instance, error) {
	return nil, nil
}

func (f *fakeNamingClient) SelectOneHealthyInstance(arg0 vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	// pick the first match if any
	key := arg0.GroupName + "@@" + arg0.ServiceName
	f.mu.RLock()
	defer f.mu.RUnlock()
	list := f.services[key]
	if len(list) == 0 {
		return nil, errors.New("no instance")
	}
	return &list[0], nil
}
func (f *fakeNamingClient) ServerHealthy() bool { return true }
func (f *fakeNamingClient) UpdateInstance(_ vo.UpdateInstanceParam) (bool, error) {
	return true, nil
}
