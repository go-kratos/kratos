package v2

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"go-common/app/infra/config/model"
	"go-common/library/log"
)

// SetFile set config file.
func (d *Dao) SetFile(name string, conf *model.Content) (err error) {
	b, err := json.Marshal(conf)
	if err != nil {
		log.Error("json.Marshal(%v) error(%v)", conf, err)
		return
	}
	p := path.Join(d.pathCache, name)
	if err = ioutil.WriteFile(p, b, 0644); err != nil {
		log.Error("ioutil.WriteFile(%s) error(%v)", p, err)
	}
	return
}

// File return config file.
func (d *Dao) File(name string) (res *model.Content, err error) {
	p := path.Join(d.pathCache, name)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error("ioutil.ReadFile(%s) error(%v)", p, err)
		return
	}
	res = &model.Content{}
	if err = json.Unmarshal(b, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", b, err)
	}
	return
}

// DelFile delete file cache.
func (d *Dao) DelFile(name string) (err error) {
	p := path.Join(d.pathCache, name)
	if err = os.Remove(p); err != nil {
		log.Error("os.Remove(%s) error(%v)", p, err)
	}
	return
}

// SetFileStr save string file.
func (d *Dao) SetFileStr(name string, val string) (err error) {
	p := path.Join(d.pathCache, name)
	if err = ioutil.WriteFile(p, []byte(val), 0644); err != nil {
		log.Error("ioutil.WriteFile(%s) error(%v)", p, err)
	}
	return
}

// FileStr get string file.
func (d *Dao) FileStr(name string) (file string, err error) {
	p := path.Join(d.pathCache, name)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Error("ioutil.ReadFile(%s) error(%v)", p, err)
		return
	}
	file = string(b)
	return
}
