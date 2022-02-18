package context

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	type ctxKey1 struct{}
	type ctxKey2 struct{}
	ctx1 := context.WithValue(context.Background(), ctxKey1{}, "https://github.com/go-kratos/")
	ctx2 := context.WithValue(context.Background(), ctxKey2{}, "https://go-kratos.dev/")

	ctx, cancel := Merge(ctx1, ctx2)
	defer cancel()

	got := ctx.Value(ctxKey1{})
	value1, ok := got.(string)
	if !ok {
		t.Errorf("expect %v, got %v", true, ok)
	}
	if !reflect.DeepEqual(value1, "https://github.com/go-kratos/") {
		t.Errorf("expect %v, got %v", "https://github.com/go-kratos/", value1)
	}

	got2 := ctx.Value(ctxKey2{})
	value2, ok := got2.(string)
	if !ok {
		t.Errorf("expect %v, got %v", true, ok)
	}
	if !reflect.DeepEqual("https://go-kratos.dev/", value2) {
		t.Errorf("expect %v, got %v", "https://go-kratos.dev/", value2)
	}

	t.Log(value1)
	t.Log(value2)
}

func TestMerge(t *testing.T) {
	type ctxKey1 struct{}
	type ctxKey2 struct{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ctx1 := context.WithValue(context.Background(), ctxKey1{}, "https://github.com/go-kratos/")
	ctx2 := context.WithValue(ctx, ctxKey2{}, "https://go-kratos.dev/")

	ctx, cancel = Merge(ctx1, ctx2)
	defer cancel()

	got := ctx.Value(ctxKey1{})
	value1, ok := got.(string)
	if !ok {
		t.Errorf("expect %v, got %v", true, ok)
	}
	if !reflect.DeepEqual(value1, "https://github.com/go-kratos/") {
		t.Errorf("expect %v, got %v", "https://github.com/go-kratos/", value1)
	}

	got2 := ctx.Value(ctxKey2{})
	value2, ok := got2.(string)
	if !ok {
		t.Errorf("expect %v, got %v", true, ok)
	}
	if !reflect.DeepEqual(value2, "https://go-kratos.dev/") {
		t.Errorf("expect %v, got %v", " https://go-kratos.dev/", value2)
	}

	t.Log(ctx)
}

func TestErr(t *testing.T) {
	ctx1, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()
	time.Sleep(time.Millisecond)

	ctx, cancel := Merge(ctx1, context.Background())
	defer cancel()
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Errorf("expect %v, got %v", context.DeadlineExceeded, ctx.Err())
	}
}

func TestDone(t *testing.T) {
	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, cancel := Merge(ctx1, context.Background())
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel()
	}()

	if <-ctx.Done() != struct{}{} {
		t.Errorf("expect %v, got %v", struct{}{}, <-ctx.Done())
	}
}

func TestFinish(t *testing.T) {
	mc := &mergeCtx{
		parent1:  context.Background(),
		parent2:  context.Background(),
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	err := mc.finish(context.DeadlineExceeded)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expect %v, got %v", context.DeadlineExceeded, err)
	}
	if !reflect.DeepEqual(mc.doneMark, uint32(1)) {
		t.Errorf("expect %v, got %v", 1, mc.doneMark)
	}
	if <-mc.done != struct{}{} {
		t.Errorf("expect %v, got %v", struct{}{}, <-mc.done)
	}
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
	if !errors.Is(mc.doneErr, context.Canceled) {
		t.Errorf("expect %v, got %v", context.Canceled, mc.doneErr)
	}

	ctx2, cancel2 := context.WithCancel(context.Background())

	mc = &mergeCtx{
		parent1:  ctx2,
		parent2:  context.Background(),
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	go func() {
		time.Sleep(time.Millisecond * 50)
		cancel2()
	}()

	mc.wait()
	t.Log(mc.doneErr)
	if !errors.Is(mc.doneErr, context.Canceled) {
		t.Errorf("expect %v, got %v", context.Canceled, mc.doneErr)
	}
}

func TestCancel(t *testing.T) {
	mc := &mergeCtx{
		parent1:  context.Background(),
		parent2:  context.Background(),
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	mc.cancel()
	if <-mc.cancelCh != struct{}{} {
		t.Errorf("expect %v, got %v", struct{}{}, <-mc.cancelCh)
	}
}

func Test_mergeCtx_Deadline(t *testing.T) {
	type fields struct {
		parent1Timeout time.Time
		parent2Timeout time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want1  bool
	}{
		{
			name:   "parent1 not deadline",
			fields: fields{time.Time{}, time.Now().Add(time.Second * 100)},
			want1:  true,
		},
		{
			name:   "parent2 not deadline",
			fields: fields{time.Now().Add(time.Second * 100), time.Time{}},
			want1:  true,
		},
		{
			name:   " parent1 parent2 not deadline",
			fields: fields{time.Time{}, time.Time{}},
			want1:  false,
		},
		{
			name:   " parent1 < parent2",
			fields: fields{time.Now().Add(time.Second * 100), time.Now().Add(time.Second * 200)},
			want1:  true,
		},
		{
			name:   " parent1 > parent2",
			fields: fields{time.Now().Add(time.Second * 100), time.Now().Add(time.Second * 50)},
			want1:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parent1, parent2 context.Context
			var cancel1, cancel2 context.CancelFunc
			if reflect.DeepEqual(tt.fields.parent1Timeout, time.Time{}) {
				parent1 = context.Background()
			} else {
				parent1, cancel1 = context.WithDeadline(context.Background(), tt.fields.parent1Timeout)
				defer cancel1()
			}
			if reflect.DeepEqual(tt.fields.parent2Timeout, time.Time{}) {
				parent2 = context.Background()
			} else {
				parent2, cancel2 = context.WithDeadline(context.Background(), tt.fields.parent2Timeout)
				defer cancel2()
			}

			mc := &mergeCtx{
				parent1: parent1,
				parent2: parent2,
			}
			got, got1 := mc.Deadline()
			t.Log(got)
			if got1 != tt.want1 {
				t.Errorf("Deadline() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_Err2(t *testing.T) {
	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()
	time.Sleep(time.Millisecond)

	ctx, cancel := Merge(ctx1, context.Background())
	defer cancel()

	if ctx.Err() != nil {
		t.Errorf("expect %v, got %v", nil, ctx.Err())
	}

	ctx1, cancel1 := context.WithCancel(context.Background())
	time.Sleep(time.Millisecond)

	ctx, cancel = Merge(ctx1, context.Background())
	defer cancel()

	cancel1()

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Errorf("expect %v, got %v", context.Canceled, ctx.Err())
	}

	ctx1, cancel1 = context.WithCancel(context.Background())
	time.Sleep(time.Millisecond)

	ctx, cancel = Merge(context.Background(), ctx1)
	defer cancel()

	cancel1()

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Errorf("expect %v, got %v", context.Canceled, ctx.Err())
	}

	ctx, cancel = Merge(context.Background(), context.Background())
	cancel()
	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Errorf("expect %v, got %v", context.Canceled, ctx.Err())
	}
}
