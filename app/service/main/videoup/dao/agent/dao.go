package agent

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/videoup/conf"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
)

// Dao is redis dao.
type Dao struct {
	c            *conf.Config
	proxyClient  *xhttp.Client
	pgcSubmitURI string
}

// New pub agent
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:            c,
		proxyClient:  xhttp.NewClient(c.HTTPClient.Read),
		pgcSubmitURI: c.PubAgent.PGCSubmit,
	}
	// add proxy
	proxyURL, _ := url.Parse(c.PubAgent.Proxy)
	dialer := &net.Dialer{
		Timeout:   time.Duration(c.HTTPClient.Write.Timeout),
		KeepAlive: time.Duration(c.HTTPClient.Write.KeepAlive),
	}
	transport := &http.Transport{
		DialContext:     dialer.DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyURL),
	}
	d.proxyClient.SetTransport(transport)
	return d
}

// AgentMsg send message to upper.
func (d *Dao) AgentMsg(c context.Context, route, filename string) (err error) {
	query := "content_type=json&r=pgc_submit/msg_pgc_submit&filename=" + filename
	msg := &archive.PubAgentParam{
		Route:       route,
		Timestamp:   strconv.FormatInt(time.Now().Unix(), 10),
		Filename:    filename,
		Xcode:       0,
		VideoDesign: "[]",
		Submit:      1,
	}
	// new request
	bs, err := json.Marshal(msg)
	if err != nil {
		log.Error("json.Marshal pubagent msg error (%v) | msg(%v)", err, msg)
		return
	}
	req, err := http.NewRequest("POST", d.pgcSubmitURI+"?"+query, bytes.NewReader(bs))
	if err != nil {
		log.Error("http.NewRequest error(%v) | uri(%s) msg(%v) req(%v)", err, d.pgcSubmitURI+"?"+query, msg, req)
		return
	}
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = d.proxyClient.Do(c, req, &res); err != nil {
		log.Error("AgentMsg d.httpW.Do error(%v) | uri(%s) msg(%v) req(%v)", err, d.pgcSubmitURI+"?"+query, msg, req)
		return
	}
	if res.Code != 0 {
		log.Error("AgentMsg res(%v) | uri(%s) msg(%v) req(%v)", res, d.pgcSubmitURI+"?"+query, msg, req)
		return
	}
	log.Info("AgentMsg success msg(%v) uri(%s)", msg, d.pgcSubmitURI+"?"+query)
	return
}
