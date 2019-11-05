package main

var _singleTemplate = `
// NAME {{or .Comment "get data from cache if miss will call source method, then add to cache."}} 
func (d *{{.StructName}}) NAME(c context.Context, {{.IDName}} KEY{{.ExtraArgsType}}) (res VALUE, err error) {
	addCache := true
	res, err = CACHEFUNC(c, {{.IDName}} {{.ExtraCacheArgs}})
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
	{{if .EnablePaging}}
	var miss VALUE
	{{end}}
	{{if .EnableSingleFlight}}
		var rr interface{}
		sf := d.cacheSFNAME({{.IDName}} {{.ExtraArgs}})
		rr, err, _ = cacheSingleFlights[SFNUM].Do(sf, func() (r interface{}, e error) {
			cache.MetricMisses.Inc("bts:NAME")
			{{if .EnablePaging}}
				var rrs [2]interface{}
				rrs[0], rrs[1], e = RAWFUNC(c, {{.IDName}} {{.ExtraRawArgs}})
				r = rrs
			{{else}}
				r, e = RAWFUNC(c, {{.IDName}} {{.ExtraRawArgs}})
			{{end}}
			return
		})
		{{if .EnablePaging}}
			res = rr.([2]interface{})[0].(VALUE)
			miss = rr.([2]interface{})[1].(VALUE)
		{{else}}
			res = rr.(VALUE)
		{{end}}
	{{else}}
		cache.MetricMisses.Inc("bts:NAME")
		{{if .EnablePaging}}
		res, miss, err = RAWFUNC(c, {{.IDName}} {{.ExtraRawArgs}})
		{{else}}
		res, err = RAWFUNC(c, {{.IDName}} {{.ExtraRawArgs}})
		{{end}}
	{{end}}
	if err != nil {
		return
	}
	{{if .EnablePaging}}
	{{else}}
		miss := res
	{{end}}
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
		ADDCACHEFUNC(c, {{.IDName}}, miss {{.ExtraAddCacheArgs}})
	{{else}}
	d.cache.Do(c, func(c context.Context) {
		ADDCACHEFUNC(c, {{.IDName}}, miss {{.ExtraAddCacheArgs}})
	})
	{{end}}
	return
}
`
