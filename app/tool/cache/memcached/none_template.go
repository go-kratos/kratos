package main

import (
	"strings"
)

var _noneGetTemplate = `
// NAME {{or .Comment "get data from mc"}} 
func (d *Dao) NAME(c context.Context) (res VALUE, err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := {{.KeyMethod}}()
	reply, err := conn.Get(key)
	if err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:NAME")
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	{{if .GetSimpleValue}}
		var v string
		err = conn.Scan(reply, &v)
	{{else}}
		{{if .GetDirectValue}}
			err = conn.Scan(reply, &res)
		{{else}}
			{{if .InitValue}}
				res = &{{.OriginValueType}}{}
				err = conn.Scan(reply, res)
			{{else}}
				res = {{.OriginValueType}}{}
				err = conn.Scan(reply, &res)
			{{end}}
		{{end}}
	{{end}}
	if err != nil {
		prom.BusinessErrCount.Incr("mc:NAME")
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	{{if .GetSimpleValue}}
		r, err := {{.ConvertBytes2Value}}
		if err != nil {
			prom.BusinessErrCount.Incr("mc:NAME")
			log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
			return
		}
		res = {{.ValueType}}(r)
	{{end}}
	return
}
`

var _noneSetTemplate = `
// NAME {{or .Comment "Set data to mc"}} 
func (d *Dao) NAME(c context.Context, val VALUE) (err error) {
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
	conn := d.mc.Get(c)
	defer conn.Close()
	key := {{.KeyMethod}}()
	{{if .SimpleValue}}
		bs := {{.ConvertValue2Bytes}}
		item := &memcache.Item{Key: key, Value: bs, Expiration: {{.ExpireCode}}, Flags: {{.Encode}}}
	{{else}}
		item := &memcache.Item{Key: key, Object: val, Expiration: {{.ExpireCode}}, Flags: {{.Encode}}}
	{{end}}
	if err = conn.Set(item); err != nil {
		prom.BusinessErrCount.Incr("mc:NAME")
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
`
var _noneAddTemplate = strings.Replace(_noneSetTemplate, "Set", "Add", -1)
var _noneReplaceTemplate = strings.Replace(_noneSetTemplate, "Set", "Replace", -1)

var _noneDelTemplate = `
// NAME {{or .Comment "delete data from mc"}} 
func (d *Dao) NAME(c context.Context) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	key := {{.KeyMethod}}()
	if err = conn.Delete(key); err != nil {
		if err == memcache.ErrNotFound {
			err = nil
			return
		}
		prom.BusinessErrCount.Incr("mc:NAME")
		log.Errorv(c, log.KV("NAME", fmt.Sprintf("%+v", err)), log.KV("key", key))
		return
	}
	return
}
`
