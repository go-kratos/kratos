package servicecomb

import (
	"context"
	"testing"

	pb "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/sc-client"
	"github.com/gofrs/uuid"

	"github.com/go-kratos/kratos/v2/registry"
)

var r *Registry

func init() {
	r = NewRegistry(&mockClient{})
}

type mockClient struct{}

func (receiver *mockClient) WatchMicroService(_ string, _ func(*sc.MicroServiceInstanceChangedEvent)) error {
	return nil
}

//nolint
func (receiver *mockClient) FindMicroServiceInstances(_,
	_, microServiceName, _ string, _ ...sc.CallOption,
) ([]*pb.MicroServiceInstance, error) {
	if microServiceName == "KratosServicecomb" {
		return []*pb.MicroServiceInstance{{}}, nil
	}
	return nil, nil
}

func (receiver *mockClient) RegisterService(_ *pb.MicroService) (string, error) {
	return "", nil
}

func (receiver *mockClient) RegisterMicroServiceInstance(_ *pb.MicroServiceInstance) (string, error) {
	return "", nil
}

func (receiver *mockClient) Heartbeat(_, _ string) (bool, error) {
	return true, nil
}

func (receiver *mockClient) UnregisterMicroServiceInstance(_, _ string) (bool, error) {
	return true, nil
}

func (receiver *mockClient) GetMicroServiceID(_, _, _, _ string, _ ...sc.CallOption) (string, error) {
	return "", nil
}

func TestRegistry(t *testing.T) {
	instanceID, err := uuid.NewV4()
	if err != nil {
		t.Fatal(err)
	}
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
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("GetService test, expected: success.", func(t *testing.T) {
		var insts []*registry.ServiceInstance
		insts, err = r.GetService(ctx, svc.Name)
		if err != nil {
			t.Fatal(err)
		}
		if len(insts) <= 0 {
			t.Errorf("inst len less than 0")
		}
	})
	t.Run("Deregister test, expected: success.", func(t *testing.T) {
		svc.ID = instanceID.String()
		err = r.Deregister(ctx, svc)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestWatcher(t *testing.T) {
	instanceID1, err := uuid.NewV4()
	if err != nil {
		t.Fatal(err)
	}
	svc1 := &registry.ServiceInstance{
		Name:      "WatcherTest",
		Version:   "0.0.1",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		ID:        instanceID1.String(),
	}
	ctx := context.TODO()
	err = r.Register(ctx, svc1)
	if err != nil {
		t.Fatal(err)
	}
	w, err := r.Watch(ctx, "WatcherTest")
	if err != nil {
		t.Fatal(err)
	}
	if w == nil {
		t.Fatal("w is nil")
	}
	sbWatcher := w.(*Watcher)
	t.Run("Watch register event, expected: success", func(t *testing.T) {
		go sbWatcher.Put(svc1)
		var instances []*registry.ServiceInstance
		instances, err = w.Next()
		if err != nil {
			t.Fatal(err)
		}
		if len(instances) == 0 {
			t.Errorf("instances is empty")
		}
		if instanceID1.String() != instances[0].ID {
			t.Errorf("expected %v, got %v", instanceID1.String(), instances[0].ID)
		}
	})
	t.Run("Watch deregister event, expected: success", func(t *testing.T) {
		// Deregister instance1.
		err = r.Deregister(ctx, svc1)
		if err != nil {
			t.Fatal(err)
		}
		go sbWatcher.Put(svc1)
		var instances []*registry.ServiceInstance
		instances, err = w.Next()
		if err != nil {
			t.Fatal(err)
		}
		if len(instances) == 0 {
			t.Errorf("instances is empty")
		}
		if instanceID1.String() != instances[0].ID {
			t.Errorf("expected %v, got %v", instanceID1.String(), instances[0].ID)
		}
	})
	t.Run("Stop test, expected: success", func(t *testing.T) {
		err = w.Stop()
		if err != nil {
			t.Error(err)
		}
	})
}
