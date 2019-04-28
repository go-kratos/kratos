package main

var _singleTemplate = `
// NAME {{or .Comment "get data from cache if miss will call source method, then add to cache."}} 
func (d *Dao) NAME(c context.Context, id KEY{{.ExtraArgsType}}) (res VALUE, err error) {
	addCache := true
	res, err = CACHEFUNC(c, id {{.ExtraCacheArgs}})
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
	{{if .EnablePaging}}
	var miss VALUE
	{{end}}
	{{if .EnableSingleFlight}}
		var rr interface{}
		sf := d.cacheSFNAME(id {{.ExtraArgs}})
		rr, err, _ = cacheSingleFlights[SFNUM].Do(sf, func() (r interface{}, e error) {
			prom.CacheMiss.Incr("NAME")
			{{if .EnablePaging}}
				var rrs [2]interface{}
				rrs[0], rrs[1], e = RAWFUNC(c, id {{.ExtraRawArgs}})
				r = rrs
			{{else}}
				r, e = RAWFUNC(c, id {{.ExtraRawArgs}})
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
		prom.CacheMiss.Incr("NAME")
		{{if .EnablePaging}}
		res, miss, err = RAWFUNC(c, id {{.ExtraRawArgs}})
		{{else}}
		res, err = RAWFUNC(c, id {{.ExtraRawArgs}})
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
		ADDCACHEFUNC(c, id, miss {{.ExtraAddCacheArgs}})
	{{else}}
	d.cache.Do(c, func(c context.Context) {
		ADDCACHEFUNC(c, id, miss {{.ExtraAddCacheArgs}})
	})
	{{end}}
	return
}
`
