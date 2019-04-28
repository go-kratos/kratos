package bvc

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/videoup/conf"
	"go-common/library/log"
	"go-common/library/xstr"
)

// Dao is message dao.
type Dao struct {
	c          *conf.Config
	httpCli    *http.Client
	capableURL string
}

// New new a activity dao.
func New(c *conf.Config) (d *Dao) {
	// http://bvc-playurl.bilibili.co/video_capable?cid=18669677&capable=0&ts=1497426598&sign=e822b4ce02fdcf91eb29fd47dcbbc3c8
	d = &Dao{
		c: c,
		httpCli: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				DisableKeepAlives: true,
			},
			Timeout: 10 * time.Second,
		},
		capableURL: c.Host.Bvc.Bvc + "/video_capable",
	}
	return
}

// VideoCapable sync cid audit result to bvc.
func (d *Dao) VideoCapable(c context.Context, aid int64, allCids []int64, capable int) (err error) {
	var times int
	if len(allCids)%50 == 0 {
		times = len(allCids) / 50
	} else {
		times = len(allCids)/50 + 1
	}
	var cids []int64
	for i := 0; i < times; i++ {
		if i == times-1 {
			cids = allCids[i*50:]
		} else {
			cids = allCids[i*50 : (i+1)*50]
		}
		if err = d.call(aid, cids, capable); err != nil {
			return
		}
	}
	return
}

func (d *Dao) call(aid int64, cids []int64, capable int) (err error) {
	var (
		// http params
		cidStr     = xstr.JoinInts(cids)
		capableStr = strconv.Itoa(capable)
		ts         = strconv.FormatInt(time.Now().Unix(), 10)
		signs      = []string{cidStr, ts, capableStr, d.c.Host.Bvc.AppendKey}
	)
	log.Info("VideoCapable  aid(%d) cid(%s) capable(%d) process start", aid, cidStr, capable)
	params := url.Values{}
	params.Set("ts", ts)
	params.Set("cid", cidStr)
	params.Set("capable", capableStr)
	// sign = md5(cid + gap_key + ts + gap_key + capable + gap_key + append_key)
	// gap_key(d.c.Host.Bvc.GapKey) = "--->"
	signStr := strings.Join(signs, d.c.Host.Bvc.GapKey)
	mh := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(mh[:])
	params.Set("sign", sign)
	url := d.capableURL + "?" + params.Encode()
	// http
	var (
		req  *http.Request
		resp *http.Response
	)
	if req, err = http.NewRequest("POST", url, nil); err != nil {
		log.Error("VideoCapable cid(%s) http.NewRequest(POST,%s) error(%v)", cidStr, url, err)
		return
	}
	if resp, err = d.httpCli.Do(req); err != nil {
		log.Error("VideoCapable cid(%s) d.httpCli.Do() error(%v)", cidStr, err)
		return
	}
	defer resp.Body.Close()
	// Update successfully only when response is HTTP/200.(bvc said this)
	if resp.StatusCode != http.StatusOK {
		log.Error("VideoCapable cid(%s) sync cid result faild. StatusCode(%d)", cidStr, resp.StatusCode)
		return
	}
	var rb []byte
	if rb, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("VideoCapable cid(%s) ioutil.ReadAll(resp.Body) error(%v)", cidStr, err)
		err = nil
		return
	}
	var res struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
	}
	if err = json.Unmarshal(rb, &res); err != nil {
		log.Error("VideoCapable cid(%s) json.Unmarshal(%s) error(%v)", cidStr, string(rb), err)
		err = nil
		return
	}
	log.Info("VideoCapable cid(%s) capable(%d) res.error(%d) or res.message(%v) process end", cidStr, capable, res.Error, res.Message)
	return
}
