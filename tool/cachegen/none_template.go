package main

var _noneTemplate = `
// NAME {{or .Comment "get data from cache if miss will call source method, then add to cache."}} 
func (d *Dao) NAME(c context.Context) (res VALUE, err error) {
	addCache := true
	res, err = CACHEFUNC(c)
	if err != nil {
		addCache = false
		err = nil
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
	prom.CacheHit.Incr("NAME")
		return
	}
	{{if .EnableSingleFlight}}
		var rr interface{}
		sf := d.cacheSFNAME()
		rr, err, _ = cacheSingleFlights[SFNUM].Do(sf, func() (r interface{}, e error) {
			prom.CacheMiss.Incr("NAME")
			r, e = RAWFUNC(c)
			return
		})
		res = rr.(VALUE)
	{{else}}
		prom.CacheMiss.Incr("NAME")
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
