package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

const (
	actionConnect    = 1
	actionDisconnect = 2

	_apiReport    = "http://dataflow.biliapi.com/log/system"
	_reportFormat = "001659%d%d|%d|%d|%d|%d|%d|%s|%s|%d"
)

var (
	httpCli     = &http.Client{Timeout: 1 * time.Second}
	_reportSucc = []byte("succeed")
	_reportOK   = []byte("OK")
)

// Report is report params.
type Report struct {
	From int64
	Aid  int64
	Cid  int64
	Mid  int64
	Key  string
	IP   string
}

func reportCh(action int, ch *Channel) {
	if ch.Room == nil {
		return
	}
	u, err := url.Parse(ch.Room.ID)
	if err != nil {
		return
	}
	if u.Scheme != "video" {
		return
	}
	paths := strings.Split(u.Path, "/")
	if len(paths) < 2 {
		return
	}
	r := &Report{Key: ch.Key, Mid: ch.Mid, IP: ch.IP}
	r.Aid, _ = strconv.ParseInt(u.Host, 10, 64)
	r.Cid, _ = strconv.ParseInt(paths[1], 10, 64)
	switch ch.Platform {
	case "ios":
		r.From = 3
	case "android":
		r.From = 2
	default:
		r.From = 1
	}
	report(action, r, ch.Room.OnlineNum())
}

func report(action int, r *Report, online int32) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	lines := fmt.Sprintf(_reportFormat, timestamp, timestamp, r.From, r.Aid, r.Cid, r.Mid, online, r.Key, r.IP, action)
	req, err := http.NewRequest("POST", _apiReport, strings.NewReader(lines))
	if err != nil {
		return
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("report: httpCli.POST(%s) error(%d)", lines, resp.StatusCode)
		return
	}
	if !bytes.Equal(b, _reportSucc) && !bytes.Equal(b, _reportOK) {
		log.Error("report error(%s)", b)
		return
	}
	if r.Mid == 19158909 {
		log.Info("report: line(%s)", lines)
	}
}
