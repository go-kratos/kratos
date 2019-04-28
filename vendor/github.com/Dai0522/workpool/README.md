# workpool
基于ringbuffer的无锁golang workpool  

# usage
```
type TestTask struct {
	name string
}

func (t *TestTask) Run() *[]byte {
	fmt.Println(t.name)
	res := []byte(t.name)
	time.Sleep(time.Duration(1 * time.Second))
	return &res
}

func createPool() *workpool.Pool {
	conf := &workpool.PoolConfig{
		MaxWorkers:     1024,
		MaxIdleWorkers: 512,
		MinIdleWorkers: 128,
		KeepAlive:      time.Duration(30 * time.Second),
	}
	p, err := workpool.NewWorkerPool(1024, conf)
	if err != nil {
		panic(err)
	}
	p.Start()
	return p
}

wp := createPool()
ft := workpool.NewFutureTask(&TestTask{
    name: "daiwei",
})
wp.Submit(ft)
res, _ := ft.Wait(time.Duration(3 * time.Second))
```