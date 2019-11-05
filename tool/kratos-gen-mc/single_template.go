package main

import (
	"strings"
)

var _singleGetTemplate = `
// NAME {{or .Comment "get data from mc"}} 
func (d *{{.StructName}}) NAME(c context.Context, id KEY {{.ExtraArgsType}}) (res VALUE, err error) {
	key := {{.KeyMethod}}(id{{.ExtraArgs}})
	{{if .GetSimpleValue}}
		var v string
		err = d.mc.Get(c, key).Scan(&v)
	{{else}}
		{{if .GetDirectValue}}
			err = d.mc.Get(c, key).Scan(&res)
		{{else}}
			{{if .PointType}}
				res = &{{.OriginValueType}}{}
				if err = d.mc.Get(c, key).Scan(res); err != nil {
					res = nil
					if err == memcache.ErrNotFound {
						err = nil
					}
				}
			{{else}}
				err = d.mc.Get(c, key).Scan(&res)
			{{end}}
		{{end}}
	{{end}}
	if err != nil {
		{{if .PointType}}
		{{else}}
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		{{end}}
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	{{if .GetSimpleValue}}
		r, err := {{.ConvertBytes2Value}}
		if err != nil {
			log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return
		}
		res = {{.ValueType}}(r)
	{{end}}
	return
}
`

var _singleSetTemplate = `
// NAME {{or .Comment "Set data to mc"}} 
func (d *{{.StructName}}) NAME(c context.Context, id KEY, val VALUE {{.ExtraArgsType}}) (err error) {
	{{if .PointType}}
      if val == nil {
        return 
      }
	{{end}}
	{{if .LenType}}
      if len(val) == 0 {
        return 
      }
	{{end}}
	key := {{.KeyMethod}}(id{{.ExtraArgs}})
	{{if .SimpleValue}}
		bs := {{.ConvertValue2Bytes}}
		item := &memcache.Item{Key: key, Value: bs, Expiration: {{.ExpireCode}}, Flags: {{.Encode}}}
	{{else}}
		item := &memcache.Item{Key: key, Object: val, Expiration: {{.ExpireCode}}, Flags: {{.Encode}}}
	{{end}}
	{{if .EnableNullCode}}
		if {{.CheckNullCode}} {
			item.Expiration = {{.ExpireNullCode}}
		}
	{{end}}
	if err = d.mc.Set(c, item); err != nil {
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
`
var _singleAddTemplate = strings.Replace(_singleSetTemplate, "Set", "Add", -1)
var _singleReplaceTemplate = strings.Replace(_singleSetTemplate, "Set", "Replace", -1)

var _singleDelTemplate = `
// NAME {{or .Comment "delete data from mc"}} 
func (d *{{.StructName}}) NAME(c context.Context, id KEY {{.ExtraArgsType}}) (err error) {
	key := {{.KeyMethod}}(id{{.ExtraArgs}})
	if err = d.mc.Delete(c, key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
`
