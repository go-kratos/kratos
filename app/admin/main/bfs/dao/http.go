package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_addVolume     = "/pitchfork/group/add_volume"
	_addFreeVolume = "/pitchfork/group/add_free_volume"
	_compact       = "/pitchfork/group/compact"
	_groupStatus   = "/pitchfork/group/status"
)

func (d *Dao) addVolumeURI() string {
	return d.c.Host.Pitchfork + _addVolume
}

func (d *Dao) addFreeVolumeURI() string {
	return d.c.Host.Pitchfork + _addFreeVolume
}

func (d *Dao) compactURI() string {
	return d.c.Host.Pitchfork + _compact
}

func (d *Dao) groupStatusURI() string {
	return d.c.Host.Pitchfork + _groupStatus
}

// AddVolume add volumes.
func (d *Dao) AddVolume(c context.Context, group string, num int64) (err error) {
	var (
		params = url.Values{}
		uri    = d.addVolumeURI()
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("group", group)
	params.Set("num", fmt.Sprint(num))
	if err = d.httpCli.Post(c, uri, ip, params, nil); err != nil {
		log.Error("add volume error(%v)", err)
	}
	return
}

// AddFreeVolume add free volume.
func (d *Dao) AddFreeVolume(c context.Context, group, dir string, num int64) (err error) {
	var (
		params = url.Values{}
		uri    = d.addFreeVolumeURI()
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("group", group)
	params.Set("idir", dir)
	params.Set("bdir", dir)
	params.Set("num", fmt.Sprint(num))
	if err = d.httpCli.Post(c, uri, ip, params, nil); err != nil {
		log.Error("add free volume error(%v)", err)
	}
	return
}

// Compact compact store disk by group.
func (d *Dao) Compact(c context.Context, group string, vid int64) (err error) {
	var (
		params = url.Values{}
		uri    = d.compactURI()
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("group", group)
	if vid > 0 {
		params.Set("vid", fmt.Sprint(vid))
	}
	if err = d.httpCli.Post(c, uri, ip, params, nil); err != nil {
		log.Error("compact group error(%v)", err)
	}
	return
}

// SetGroupStatus set store status by group id.
func (d *Dao) SetGroupStatus(c context.Context, group, status string) (err error) {
	var (
		params = url.Values{}
		uri    = d.groupStatusURI()
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	params.Set("group", group)
	params.Set("status", status)
	if err = d.httpCli.Post(c, uri, ip, params, nil); err != nil {
		log.Error("set group status error(%v)", err)
	}
	return
}
