package report

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	xhttp "net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	mdlrep "go-common/app/job/main/tv/model/report"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_httpOK     = "0"
	_pid        = "73" // platform: 73 is tv
	_plat       = 5    // platform: 5  is tv
	_tv         = "android_tv_yst"
	_actDur     = `select request_uri,time_iso,ip,version,buvid,fts,proid,chid,pid,brand,deviceid,model,osver,ctime,mid,ver,net,oid,eid,start_time,end_time,duration,openudid,idfa,mac,is_coldstart,session_id,buvid_ext from ods.app_active_duration where unix_timestamp(ctime, 'yyyyMMddHHmmss')>=%d and unix_timestamp(ctime, 'yyyyMMddHHmmss')<%d and log_date=%s and pid=%s`
	_playDur    = `select stime,build,buvid,mobi_app,platform,session,mid,aid,cid,sid,epid,type,sub_type,quality,total_time,paused_time,played_time,video_duration,play_type,network_type,last_play_progress_time,max_play_progress_time,device,epid_status,play_status,user_status,actual_played_time,auto_play,detail_play_time,list_play_time from ods.app_play_duration where stime>=%d and stime<%d and log_date=%s and mobi_app=%s`
	_visitEvent = `select request_uri,time_iso,ip,version,buvid,fts,proid,chid,pid,brand,deviceid,model,osver,mid,ctime,ver,net,oid,page_name,page_arg,ua,h5_chid,unix_timestamp(ctime, 'yyyyMMddHHmmss') from ods.app_visit_event where unix_timestamp(ctime, 'yyyyMMddHHmmss')>=%d and unix_timestamp(ctime, 'yyyyMMddHHmmss')<%d and log_date=%s and pid=%s`
	_arcClick   = `select r_type,avid,cid,part,mid,stime,did,ip,ua,buvid,cookie_sid,refer,type,sub_type,sid,epid,platform,device from ods.ods_archive_click where stime>=%d and stime<%d and log_date=%s and plat=%d`
)

var (
	signParams = []string{"appKey", "timestamp", "version"}
	_userAgent = "User-Agent"
)

// Report .
func (d *Dao) Report(ctx context.Context, table string) (info string, err error) {
	var (
		v          = url.Values{}
		ip         = metadata.String(ctx, metadata.RemoteIP)
		logdata    = d.queryfmt()
		start, end = d.dealDate()
		query      string
		res        struct {
			Code      int    `json:"code"`
			Msg       string `json:"msg"`
			StatusURL string `json:"jobStatusUrl"`
		}
	)
	logdata = "'" + logdata + "'"
	switch table {
	case mdlrep.ArchiveClick:
		query = fmt.Sprintf(_arcClick, start, end, logdata, _plat)
	case mdlrep.ActiveDuration:
		query = fmt.Sprintf(_actDur, start, end, logdata, `"`+_pid+`"`)
	case mdlrep.PlayDuration:
		query = fmt.Sprintf(_playDur, start, end, logdata, "'"+_tv+"'")
	case mdlrep.VisitEvent:
		query = fmt.Sprintf(_visitEvent, start, end, logdata, `"`+_pid+`"`)
	default:
		err = errors.New("table is nill")
		return
	}
	v.Set("appKey", d.conf.DpClient.Key)
	v.Set("signMethod", "md5")
	v.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	v.Set("version", "1.0")
	v.Set("query", query)
	if err = d.newRequest(ctx, d.conf.Report.ReportURI, ip, v, &res); err != nil {
		log.Error("newRequest url(%v), err(%v)", d.conf.Report.ReportURI+"?"+v.Encode(), err)
		return
	}
	if res.Code == 200 && res.StatusURL != "" {
		info = res.StatusURL
	}
	return
}

// CheckJob .
func (d *Dao) CheckJob(ctx context.Context, urls string) (res *mdlrep.DpCheckJobResult, err error) {
	var (
		v  = url.Values{}
		ip = metadata.String(ctx, metadata.RemoteIP)
	)
	res = &mdlrep.DpCheckJobResult{}
	v.Set("appKey", d.conf.DpClient.Key)
	v.Set("signMethod", "md5")
	v.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	v.Set("version", "1.0")
	if err = d.newRequest(ctx, urls, ip, v, &res); err != nil {
		log.Error("d.newRequest error(%v)", err)
		return
	}
	if res.Code != xhttp.StatusOK {
		log.Error("d.CheckJob newRequest error code:%d ; url(%s) ", res.Code, urls+"?"+v.Encode())
		err = fmt.Errorf("code(%d) msg(%s) statusID(%d) statusID(%s)", res.Code, res.Msg, res.StatusID, res.StatusMsg)
	}
	return
}

// PostRequest .
func (d *Dao) PostRequest(ctx context.Context, body string) (err error) {
	var (
		req *xhttp.Request
		res struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
	)
	if req, err = xhttp.NewRequest(xhttp.MethodPost, d.conf.Report.UpDataURI, strings.NewReader(body)); err != nil {
		log.Error("xhttp.NewRequest url(%s) error (%v)", d.conf.Report.UpDataURI, err)
		return
	}
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	if err = d.httpR.Do(ctx, req, &res); err != nil {
		log.Error("d.httpReq.Do error(%v) url(%s)", err, d.conf.Report.UpDataURI)
		return
	}
	if res.Code != _httpOK {
		log.Error("PostRequest error code:%s ; url(%s) ", res.Code, d.conf.Report.UpDataURI)
		err = fmt.Errorf("code(%s)", res.Code)
	}
	return
}

// NewRequest new http request with method, url, ip, values and headers .
func (d *Dao) newRequest(c context.Context, url, realIP string, params url.Values, res interface{}) (err error) {
	enc, err := d.sign(params)
	if err != nil {
		log.Error("url:%s,params:%v", url, params)
		return
	}
	if enc != "" {
		url = url + "?" + enc
	}
	req, err := xhttp.NewRequest(xhttp.MethodGet, url, nil)
	if err != nil {
		log.Error("xhttp.NewRequest method:%s,url:%s", xhttp.MethodGet, url)
		return
	}
	req.Header.Set(_userAgent, "haoguanwei@bilibili.com  "+env.AppID)
	if err != nil {
		return
	}
	return d.httpR.Do(c, req, res)
}

// Sign calc appkey and appsecret sign .
func (d *Dao) sign(params url.Values) (query string, err error) {
	tmp := params.Encode()
	signTmp := d.encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(d.conf.DpClient.Secret)
	b.WriteString(signTmp)
	b.WriteString(d.conf.DpClient.Secret)
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(strings.ToUpper(hex.EncodeToString(mh[:])))
	query = qb.String()
	return
}

// Encode encodes the values into ``URL encoded'' form .
// ("bar=baz&foo=quux") sorted by key .
func (d *Dao) encode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		found := false
		for _, p := range signParams {
			if p == k {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		vs := v[k]
		prefix := k
		for _, v := range vs {
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}

// delDate create start and end time .
func (d *Dao) dealDate() (start, end int64) {
	var (
		timeDelay = "-" + d.conf.Report.TimeDelay
		timeSpan  = "-" + time.Duration(d.conf.Report.SeTimeSpan).String()
	)
	endTime := time.Now()
	et, _ := time.ParseDuration(timeDelay)
	endTmp := endTime.Add(et)
	end = endTime.Add(et).Unix()
	st, _ := time.ParseDuration(timeSpan)
	start = endTmp.Add(st).Unix()
	return
}

func (d *Dao) queryfmt() (logdata string) {
	var (
		dtime    = time.Now()
		timeData = "-" + d.conf.Report.TimeDelay
	)
	et, _ := time.ParseDuration(timeData)
	logdata = dtime.Add(et).Format("20060102")
	return
}
