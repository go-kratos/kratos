package servicecomb

import (
	"context"
	"testing"

	pb "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/sc-client"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

var r *Registry

func init() {
	c := &mockClient{}
	r = NewRegistry(c)
}

type mockClient struct{}

func (receiver *mockClient) WatchMicroService(microServiceID string, callback func(*sc.MicroServiceInstanceChangedEvent)) error {
	return nil
}

//nolint
func (receiver *mockClient) FindMicroServiceInstances(consumerID,
	appID, microServiceName, versionRule string, opts ...sc.CallOption,

) ([]*pb.MicroServiceInstance, error) {
	if microServiceName == "KratosServicecomb" {
		return []*pb.MicroServiceInstance{{}}, nil
	}
	return nil, nil
}

func (receiver *mockClient) RegisterService(microService *pb.MicroService) (string, error) {
	return "", nil
}

func (receiver *mockClient) RegisterMicroServiceInstance(microServiceInstance *pb.MicroServiceInstance) (string, error) {
	return "", nil
}

func (receiver *mockClient) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	return true, nil
}

func (receiver *mockClient) UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) (bool, error) {
	return true, nil
}

func (receiver *mockClient) GetMicroServiceID(appID, microServiceName, version, env string, opts ...sc.CallOption) (string, error) {
	return "", nil
}

// TestRegistry
func TestRegistry(t *testing.T) {
	instanceID, err := uuid.NewV4()
	assert.NoError(t, err)
	svc := &registry.ServiceInstance{
		Name:      "KratosServicecomb",
		Version:   "0.0.1",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		ID:        instanceID.String(),
	}
	ctx := context.TODO()
	t.Run("Register test, expected: success.", func(t *testing.T) {
		err = r.Register(ctx, svc)
		assert.NoError(t, err)
	})
	t.Run("GetService test, expected: success.", func(t *testing.T) {
		var insts []*registry.ServiceInstance
		insts, err = r.GetService(ctx, svc.Name)
		assert.NoError(t, err)
		assert.Greater(t, len(insts), 0)
	})
	t.Run("Deregister test, expected: success.", func(t *testing.T) {
		svc.ID = instanceID.String()
		err = r.Deregister(ctx, svc)
		assert.NoError(t, err)
	})
}

func TestWatcher(t *testing.T) {
	ctx := context.TODO()
	instanceID1, err := uuid.NewV4()
	assert.NoError(t, err)
	svc1 := &registry.ServiceInstance{
		Name:      "WatcherTest",
		Version:   "0.0.1",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		ID:        instanceID1.String(),
	}
	err = r.Register(ctx, svc1)
	assert.NoError(t, err)
	w, err := r.Watch(ctx, "WatcherTest")
	assert.NoError(t, err)
	assert.NotEmpty(t, w)
	sbWatcher := w.(*Watcher)
	t.Run("Watch register event, expected: success", func(t *testing.T) {
		go sbWatcher.Put(svc1)
		var instances []*registry.ServiceInstance
		instances, err = w.Next()
		assert.NoError(t, err)
		assert.NotEmpty(t, instances)
		assert.Equal(t, instanceID1.String(), instances[0].ID)
	})
	t.Run("Watch deregister event, expected: success", func(t *testing.T) {
		// Deregister instance1.
		err = r.Deregister(ctx, svc1)
		assert.NoError(t, err)
		go sbWatcher.Put(svc1)
		var instances []*registry.ServiceInstance
		instances, err = w.Next()
		assert.NoError(t, err)
		assert.NotEmpty(t, instances)
		assert.Equal(t, instanceID1.String(), instances[0].ID)
	})
	t.Run("Stop test, expected: success", func(t *testing.T) {
		err = w.Stop()
		assert.NoError(t, err)
	})
}
