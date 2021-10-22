package ewma

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

	assert.Equal(t, float64(100), wn.Weight())
	done := wn.Pick()
	assert.NotNil(t, done)
	done2 := wn.Pick()
	assert.NotNil(t, done2)

	time.Sleep(time.Millisecond * 10)
	done(context.Background(), selector.DoneInfo{})
	assert.Less(t, float64(30000), wn.Weight())
	assert.Greater(t, float64(60000), wn.Weight())

	assert.Greater(t, time.Millisecond*15, wn.PickElapsed())
	assert.Less(t, time.Millisecond*5, wn.PickElapsed())
}

func TestDirectError(t *testing.T) {
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

	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = context.DeadlineExceeded
		}
		done := wn.Pick()
		assert.NotNil(t, done)
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}

	assert.Less(t, float64(30000), wn.Weight())
	assert.Greater(t, float64(60000), wn.Weight())
}

func TestDirectErrorHandler(t *testing.T) {
	b := &Builder{
		ErrHandler: func(err error) bool {
			return err != nil
		},
	}
	wn := b.Build(selector.NewNode(
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = context.DeadlineExceeded
		}
		done := wn.Pick()
		assert.NotNil(t, done)
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}

	assert.Less(t, float64(30000), wn.Weight())
	assert.Greater(t, float64(60000), wn.Weight())
}
