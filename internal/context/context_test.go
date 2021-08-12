package context

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx1 := context.WithValue(context.Background(), "go-kratos", "https://github.com/go-kratos/")
	ctx2 := context.WithValue(context.Background(), "kratos", "https://go-kratos.dev/")

	ctx, cancel := Merge(ctx1, ctx2)
	defer cancel()

	got := ctx.Value("go-kratos")
	value1, ok := got.(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, value1, "https://github.com/go-kratos/")
	//
	got2 := ctx.Value("kratos")
	value2, ok := got2.(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, value2, "https://go-kratos.dev/")

	t.Log(value1)
	t.Log(value2)
}

func TestErr(t *testing.T) {
	ctx1, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	ctx, cancel := Merge(ctx1, context.Background())
	defer cancel()

	assert.Equal(t, ctx.Err(), context.DeadlineExceeded)
}

func TestDone(t *testing.T) {
	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, cancel := Merge(ctx1, context.Background())
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	assert.Equal(t, <-ctx.Done(), struct{}{})
}

func TestFinish(t *testing.T) {
	mc := &mergeCtx{
		parent1:  context.Background(),
		parent2:  context.Background(),
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	err := mc.finish(context.DeadlineExceeded)
	assert.Equal(t, err, context.DeadlineExceeded)
	assert.Equal(t, mc.doneMark, uint32(1))
	assert.Equal(t, <-mc.done, struct{}{})
}

func TestWait(t *testing.T) {
	ctx1, cancel := context.WithCancel(context.Background())

	mc := &mergeCtx{
		parent1:  ctx1,
		parent2:  context.Background(),
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	mc.wait()
	t.Log(mc.doneErr)
	assert.Equal(t, mc.doneErr, context.Canceled)
}

func TestCancel(t *testing.T) {
	mc := &mergeCtx{
		parent1:  context.Background(),
		parent2:  context.Background(),
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	mc.cancel()

	assert.Equal(t, <-mc.cancelCh, struct{}{})
}
