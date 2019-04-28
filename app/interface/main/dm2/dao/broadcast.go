package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_goimChatURI = "/x/internal/chat/push/room"

	_broadcastURI = "/x/internal/broadcast/push/room"

	_broadcastDmOperation = 1000
	_broadcastDmRoomFmt   = "video://%d/%d" //video://{aid}/{cid}
)

func (d *Dao) goimChatURI() string {
	return d.conf.Host.API + _goimChatURI
}

func (d *Dao) broadcastURI() string {
	return d.conf.Host.API + _broadcastURI
}

// BroadcastInGoim send dm msg in realtime.
func (d *Dao) BroadcastInGoim(c context.Context, cid, aid int64, info json.RawMessage) (err error) {
	var (
		res struct {
			Code int64 `json:"code"`
		}
		data = map[string]interface{}{
			"cmd":  "DM",
			"info": info,
		}
	)
	v, err := json.Marshal(data)
	if err != nil {
		log.Error("json.Marshal(%s) error(%v)", info, err)
		return
	}
	url := fmt.Sprintf("%s?cids=%d&aid=%d", d.goimChatURI(), cid, aid)
	req, err := http.NewRequest("POST", url, bytes.NewReader(v))
	if err != nil {
		log.Error("broadcast http.NewRequest() error(%v)", err)
		return
	}
	if err = d.httpCli.Do(c, req, &res); err != nil {
		log.Error("httpCli.Do(%s) error(%v)", url, err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("broadcast api failed(%d)", res.Code)
		log.Error("broadcast(%s) res code(%d)", url, res.Code)
	}
	return
}

// Broadcast send dm msg in realtime.
func (d *Dao) Broadcast(c context.Context, cid, aid int64, msg string) (err error) {
	var (
		res struct {
			Code int64 `json:"code"`
		}
	)
	params := url.Values{}
	params.Set("operation", fmt.Sprint(_broadcastDmOperation))
	params.Set("room", fmt.Sprintf(_broadcastDmRoomFmt, aid, cid))
	params.Set("message", msg)
	if err = d.httpCli.Post(c, d.broadcastURI(), metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("httpCli.Do(%s) error(%v)", d.broadcastURI(), err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("broadcast api failed(%d)", res.Code)
		log.Error("broadcast(%s) res code(%d)", d.broadcastURI(), res.Code)
	}
	return
}
