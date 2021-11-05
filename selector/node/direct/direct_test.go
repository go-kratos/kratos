package direct

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"

	"github.com/stretchr/testify/assert"
)

func TestDirect(t *testing.T) {
	b := &Builder{}
	wn := b.Build(selector.NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	done := wn.Pick()
	assert.NotNil(t, done)
	time.Sleep(time.Millisecond * 10)
	done(context.Background(), selector.DoneInfo{})
	assert.Equal(t, float64(10), wn.Weight())
	assert.Greater(t, time.Millisecond*15, wn.PickElapsed())
	assert.Less(t, time.Millisecond*5, wn.PickElapsed())
}

func TestDirectDefaultWeight(t *testing.T) {
	b := &Builder{}
	wn := b.Build(selector.NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
		}))

	done := wn.Pick()
	assert.NotNil(t, done)
	time.Sleep(time.Millisecond * 10)
	done(context.Background(), selector.DoneInfo{})
	assert.Equal(t, float64(100), wn.Weight())
	assert.Greater(t, time.Millisecond*20, wn.PickElapsed())
	assert.Less(t, time.Millisecond*5, wn.PickElapsed())
}
