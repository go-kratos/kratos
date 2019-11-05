package main

var _noneTemplate = `
// NAME {{or .Comment "get data from cache if miss will call source method, then add to cache."}} 
func (d *{{.StructName}}) NAME(c context.Context) (res VALUE, err error) {
	addCache := true
	res, err = CACHEFUNC(c)
	if err != nil {
		{{if .CacheErrContinue}}
		addCache = false
		err = nil
		{{else}}
		return
		{{end}}
	}
	{{if .EnableNullCache}}
	defer func() {
		{{if .SimpleValue}} if res == {{.NullCache}} { {{else}} if {{.CheckNullCode}} { {{end}}
			res = {{.ZeroValue}}
		}
	}()
	{{end}}
	{{if .GoValue}}
	if len(res) != 0 {
	{{else}}
	if res != {{.ZeroValue}} {
	{{end}}
	cache.MetricHits.Inc("bts:NAME")
		return
	}
	{{if .EnableSingleFlight}}
		var rr interface{}
		sf := d.cacheSFNAME()
		rr, err, _ = cacheSingleFlights[SFNUM].Do(sf, func() (r interface{}, e error) {
			cache.MetricMisses.Inc("bts:NAME")
			r, e = RAWFUNC(c)
			return
		})
		res = rr.(VALUE)
	{{else}}
		cache.MetricMisses.Inc("bts:NAME")
		res, err = RAWFUNC(c)
	{{end}}
	if err != nil {
		return
	}
	var miss = res
	{{if .EnableNullCache}}
	{{if .GoValue}}
	if len(miss) == 0 {
	{{else}}
	if miss == {{.ZeroValue}} {
	{{end}}
		miss = {{.NullCache}}
	}
	{{end}}
	if !addCache {
		return
	}
	{{if .Sync}}
		ADDCACHEFUNC(c, miss)
	{{else}}
	d.cache.Do(c, func(c context.Context) {
		ADDCACHEFUNC(c, miss)
	})
	{{end}}
	return
}
`
