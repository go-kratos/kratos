package errgroup

import (
	"context"
)

func fakeRunTask(ctx context.Context) error {
	return nil
}

func ExampleGroup_group() {
	g := Group{}
	g.Go(func(context.Context) error {
		return fakeRunTask(context.Background())
	})
	g.Go(func(context.Context) error {
		return fakeRunTask(context.Background())
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_ctx() {
	g := WithContext(context.Background())
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_cancel() {
	g := WithCancel(context.Background())
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(ctx)
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_maxproc() {
	g := Group{}
	// set max concurrency
	g.GOMAXPROCS(2)
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(context.Background())
	})
	g.Go(func(ctx context.Context) error {
		return fakeRunTask(context.Background())
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}
