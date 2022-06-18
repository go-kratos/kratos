package servicecomb

import (
	"context"
	"github.com/go-chassis/sc-client"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestRegistry
func TestRegistry(t *testing.T) {
	c, err := sc.NewClient(sc.Options{
		Endpoints: []string{"127.0.0.1:30100"},
	})
	assert.NoError(t, err)
	r := NewRegistry(c)
	instanceId, err := uuid.NewV4()
	assert.NoError(t, err)
	svc := &registry.ServiceInstance{
		Name:      "KratosServicecomb",
		Version:   "0.0.1",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		ID:        instanceId.String(),
	}
	ctx := context.TODO()
	t.Run("Register test, expected: success.", func(t *testing.T) {
		err = r.Register(ctx, svc)
		assert.NoError(t, err)
	})
	t.Run("GetService test, expected: success.", func(t *testing.T) {
		insts, err := r.GetService(ctx, svc.Name)
		assert.NoError(t, err)
		assert.Greater(t, len(insts), 0)
	})
	t.Run("Deregister test, expected: success.", func(t *testing.T) {
		svc.ID = instanceId.String()
		err = r.Deregister(ctx, svc)
		assert.NoError(t, err)
	})
}

func TestWatcher(t *testing.T) {
	c, err := sc.NewClient(sc.Options{
		Endpoints: []string{"127.0.0.1:30100"},
	})
	assert.NoError(t, err)
	r := NewRegistry(c)
	ctx := context.TODO()
	instanceId1, err := uuid.NewV4()
	assert.NoError(t, err)
	svc1 := &registry.ServiceInstance{
		Name:      "WatcherTest",
		Version:   "0.0.1",
		Metadata:  map[string]string{"app": "kratos"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		ID:        instanceId1.String(),
	}
	err = r.Register(ctx, svc1)
	assert.NoError(t, err)
	w, err := r.Watch(ctx, "WatcherTest")
	assert.NoError(t, err)
	assert.NotEmpty(t, w)
	t.Run("Watch register event, expected: success", func(t *testing.T) {
		instances, err := w.Next()
		assert.NoError(t, err)
		assert.NotEmpty(t, instances)
		assert.Equal(t, instanceId1.String(), instances[0].ID)
		err = w.Stop()
		assert.NoError(t, err)
	})
	t.Run("Watch deregister event, expected: success", func(t *testing.T) {
		//Deregister instance1 after 5 seconds.
		_, err := w.Next()
		assert.NoError(t, err)
		err = r.Deregister(ctx, svc1)
		assert.NoError(t, err)
		instances, err := w.Next()
		assert.NoError(t, err)
		assert.NotEmpty(t, instances)
		assert.Equal(t, instanceId1.String(), instances[0].ID)
		err = w.Stop()
		assert.NoError(t, err)
	})
	t.Run("Stop test, expected: success", func(t *testing.T) {
		err = w.Stop()
		assert.NoError(t, err)
	})

}
