package errgroup

import (
	"context"
	"sync"
)

func fakeRunTask(ctx context.Context) error {
	return nil
}

func ExampleGroup_group() {
	g := Group{}
	g.Go(func() error {
		return fakeRunTask(context.Background())
	})
	g.Go(func() error {
		return fakeRunTask(context.Background())
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_ctx() {
	g, ctx := WithContext(context.Background())
	g.Go(func() error {
		return fakeRunTask(ctx)
	})
	g.Go(func() error {
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
	g.Go(func() error {
		return fakeRunTask(context.Background())
	})
	g.Go(func() error {
		return fakeRunTask(context.Background())
	})
	if err := g.Wait(); err != nil {
		// handle err
	}
}

func ExampleGroup_waitgroup() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		// do something
		wg.Done()
	}()
	go func() {
		// do something
		wg.Done()
	}()
	wg.Wait()
}
