package main

var _multiTemplate = `
// NAME {{or .Comment "get data from cache if miss will call source method, then add to cache."}} 
func (d *{{.StructName}}) NAME(c context.Context, {{.IDName}} []KEY{{.ExtraArgsType}}) (res map[KEY]VALUE, err error) {
	if len({{.IDName}}) == 0 {
		return
	}
	addCache := true
	if res, err = CACHEFUNC(c, {{.IDName}} {{.ExtraCacheArgs}});err != nil {
		{{if .CacheErrContinue}}
		addCache = false
		res = nil
		err = nil
		{{else}}
		return
		{{end}}
	}
	var miss []KEY
	for _, key := range {{.IDName}} {
	{{if .GoValue}}
	if (res == nil) || (len(res[key]) == 0) {
	{{else}}
		{{if .NumberValue}}
		if _, ok := res[key]; !ok {
		{{else}}
		if (res == nil) || (res[key] == {{.ZeroValue}}) {
		{{end}}
	{{end}}
			miss = append(miss, key)
		}
	}
	cache.MetricHits.Add(float64(len({{.IDName}}) - len(miss)), "bts:NAME")
	{{if .EnableNullCache}}
	for k, v := range res {
		{{if .SimpleValue}} if v == {{.NullCache}} { {{else}} if {{.CheckNullCode}} { {{end}}
			delete(res, k)
		}
	}
	{{end}}
	missLen := len(miss)
	if missLen == 0 {
		return 
	}
	{{if .EnableBatch}}
	missData := make(map[KEY]VALUE, missLen)
	{{else}}
	var missData map[KEY]VALUE
	{{end}}
	{{if .EnableSingleFlight}}
		var rr interface{}
		sf := d.cacheSFNAME({{.IDName}} {{.ExtraArgs}})
		rr, err, _ = cacheSingleFlights[SFNUM].Do(sf, func() (r interface{}, e error) {
			cache.MetricMisses.Add(float64(len(miss)), "bts:NAME")
			r, e = RAWFUNC(c, miss {{.ExtraRawArgs}})
			return
		})
		missData = rr.(map[KEY]VALUE)
	{{else}}
		{{if .EnableBatch}}
			cache.MetricMisses.Add(float64(missLen), "bts:NAME")
			var mutex sync.Mutex
			{{if .BatchErrBreak}}
				group := errgroup.WithCancel(c)
			{{else}}
				group := errgroup.WithContext(c)
			{{end}}
			if missLen > MAXGROUP {
			group.GOMAXPROCS(MAXGROUP)
			}
			var run = func(ms []KEY) {
				group.Go(func(ctx context.Context) (err error) {
					data, err := RAWFUNC(ctx, ms {{.ExtraRawArgs}})
					mutex.Lock()
					for k, v := range data {
						missData[k] = v
					}
					mutex.Unlock()
					return
				})
			}
			var (
				i int
				n = missLen/GROUPSIZE
			)
			for i=0; i< n; i++{
				run(miss[i*GROUPSIZE:(i+1)*GROUPSIZE])
			}
			if len(miss[i*GROUPSIZE:]) > 0 {
				run(miss[i*GROUPSIZE:])
			}
			err = group.Wait()
		{{else}}
			cache.MetricMisses.Add(float64(len(miss)), "bts:NAME")
			missData, err = RAWFUNC(c, miss {{.ExtraRawArgs}})
		{{end}}
	{{end}}
	if res == nil {
		res = make(map[KEY]VALUE, len({{.IDName}}))
	}
	for k, v := range missData {
		res[k] = v
	}
	if err != nil {
		return
	}
	{{if .EnableNullCache}}
		for _, key := range miss {
			{{if .GoValue}}
			if len(res[key]) == 0 { 
			{{else}}
			if res[key] == {{.ZeroValue}} { 
			{{end}}
				missData[key] = {{.NullCache}}
			}
		}
	{{end}}
	if !addCache {
		return
	}
	{{if .Sync}}
	ADDCACHEFUNC(c, missData {{.ExtraAddCacheArgs}})
	{{else}}
	d.cache.Do(c, func(c context.Context) {
		ADDCACHEFUNC(c, missData {{.ExtraAddCacheArgs}})
	})
	{{end}}
	return
}
`
