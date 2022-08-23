package servicecomb

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-chassis/cari/discovery"
	pb "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/cari/pkg/errsvc"
	"github.com/go-chassis/sc-client"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/gofrs/uuid"
)

func init() {
	appID = os.Getenv(appIDVar)
	if appID == "" {
		appID = "default"
	}
	env = os.Getenv(envVar)
}

var (
	_ registry.Registrar = (*Registry)(nil)
	_ registry.Discovery = (*Registry)(nil)
)

var (
	curServiceID string
	appID        string
	env          string
)

const (
	appIDKey         = "appId"
	envKey           = "environment"
	envVar           = "CAS_ENVIRONMENT_ID"
	appIDVar         = "CAS_APPLICATION_NAME"
	frameWorkName    = "kratos"
	frameWorkVersion = "v2"
)

type RegistryClient interface {
	GetMicroServiceID(appID, microServiceName, version, env string, opts ...sc.CallOption) (string, error)
	FindMicroServiceInstances(consumerID, appID, microServiceName, versionRule string, opts ...sc.CallOption) ([]*discovery.MicroServiceInstance, error)
	RegisterService(microService *discovery.MicroService) (string, error)
	RegisterMicroServiceInstance(microServiceInstance *discovery.MicroServiceInstance) (string, error)
	Heartbeat(microServiceID, microServiceInstanceID string) (bool, error)
	UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) (bool, error)
	WatchMicroService(microServiceID string, callback func(*sc.MicroServiceInstanceChangedEvent)) error
}

// Registry is servicecomb registry.
type Registry struct {
	cli RegistryClient
}

func NewRegistry(client RegistryClient) *Registry {
	r := &Registry{
		cli: client,
	}
	return r
}

func (r *Registry) GetService(_ context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	instances, err := r.cli.FindMicroServiceInstances("", appID, serviceName, "")
	if err != nil {
		return nil, err
	}
	svcInstances := make([]*registry.ServiceInstance, 0, len(instances))
	for _, instance := range instances {
		svcInstances = append(svcInstances, &registry.ServiceInstance{
			ID:        instance.InstanceId,
			Name:      serviceName,
			Metadata:  instance.Properties,
			Endpoints: instance.Endpoints,
			Version:   instance.ServiceId,
		})
	}
	return svcInstances, nil
}

func (r *Registry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newWatcher(ctx, r.cli, serviceName)
}

func (r *Registry) Register(_ context.Context, svcIns *registry.ServiceInstance) error {
	fw := &discovery.FrameWork{
		Name:    frameWorkName,
		Version: frameWorkVersion,
	}
	ms := &discovery.MicroService{
		ServiceName: svcIns.Name,
		AppId:       appID,
		Version:     svcIns.Version,
		Environment: env,
		Framework:   fw,
	}
	// 先尝试创建微服务
	sid, err := r.cli.RegisterService(ms)
	// 若失败，说明服务可能已注册
	if err != nil {
		registryException, ok := err.(*sc.RegistryException)
		if !ok {
			return err
		}
		var svcErr errsvc.Error
		parseErr := json.Unmarshal([]byte(registryException.Message), &svcErr)
		if parseErr != nil {
			return parseErr
		}
		// 若错误码显示服务未注册，直接返回
		if svcErr.Code != pb.ErrServiceAlreadyExists {
			return err
		}
		sid, err = r.cli.GetMicroServiceID(appID, ms.ServiceName, ms.Version, ms.Environment)
		if err != nil {
			return err
		}
	} else {
		// 保存当前版本微服务对应的sid
		curServiceID = sid
	}
	props := make(map[string]string)
	props[appIDKey] = appID
	props[envKey] = env
	if svcIns.ID == "" {
		var id uuid.UUID
		id, err = uuid.NewV4()
		if err != nil {
			return err
		}
		svcIns.ID = id.String()
	}
	_, err = r.cli.RegisterMicroServiceInstance(&discovery.MicroServiceInstance{
		InstanceId: svcIns.ID,
		ServiceId:  sid,
		Endpoints:  svcIns.Endpoints,
		HostName:   svcIns.ID,
		Properties: props,
		Version:    svcIns.Version,
	})
	if err != nil {
		return err
	}
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			<-ticker.C
			_, err = r.cli.Heartbeat(sid, svcIns.ID)
			if err != nil {
				log.Errorf("failed to send heartbeat: %v", err)
				continue
			}
		}
	}()
	return nil
}

func (r *Registry) Deregister(_ context.Context, svcIns *registry.ServiceInstance) error {
	sid, err := r.cli.GetMicroServiceID(appID, svcIns.Name, svcIns.Version, env)
	if err != nil {
		return err
	}
	_, err = r.cli.UnregisterMicroServiceInstance(sid, svcIns.ID)
	if err != nil {
		return err
	}
	return nil
}
