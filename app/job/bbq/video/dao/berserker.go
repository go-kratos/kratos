package dao

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go-common/app/job/bbq/video/conf"
	"go-common/app/job/bbq/video/model"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	xhttp "net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"io/ioutil"

	pkgerr "github.com/pkg/errors"
)

const (
	_jobStatusSuccess = 1
	_jobStatusFailed  = 2
	_jobStatusDoing   = 3
	_jobStatusWaiting = 4
	//_httpHeaderUser     = "x1-bilispy-user"
	//_httpHeaderColor    = "x1-bilispy-color"
	//_httpHeaderTimeout  = "x1-bilispy-timeout"
	_httpHeaderRemoteIP = "x-backend-bili-real-ip"
	_userAgent          = "User-Agent"
	_noKickUserAgent    = "yangyucheng@bilibili.com"
	_queryJSON          = `{"select":[],"where":{"log_date":{"in":["%s"]}},"page":{"limit":1000},"sort":{"play":-1}}`
	_queryJSONOper      = `{"select":[],"where":{"log_date":{"in":["%s"]},"cid":{"gt":%d}},"page":{"limit":5000},"sort":{"cid":1}}`
	_hscUserAgent       = "huangshancheng@bilibili.com"
	_lzqUserAgent       = "liuzhiquan@bilibili.com"
	_chmUserAgent       = "caiheming@bilibili.com"
	_ljUserAgent        = "liujin@bilibili.com"
	//_userDmgQueryJSON   = `{"select":[],"where":{"log_date":{"in":["%s"]},"mid":{"gt":"%s"}},"sort":{"mid":1},"page":{"limit":200}}`
	_upUserDmgQueryJSON = `{"select":[],"where":{"mid":{"gt":%d}},"sort":{"mid":1},"page":{"limit":200}}`
	_userDmgQueryHive   = `select mid, gender, age, geo, content_tag, viewed_video, content_zone, content_count, follow_ups from sycpb.hbase_dmp_tag where last_active_date >= %s and length(viewed_video) > 0`
	_upMidQueryHive     = `select mid from ods.ods_member_relation_stat where log_date = %s  and follower>= 10000 limit 100`
	//_upMidQueryHive = `{"select":["name":"mid"],"where":{"log_date":{"in":["%s"]},"follower":{"gte":10000}, "pages":{"limit":10}}`
	_basePathUserProfile      = "/tmp/"
	_basePathUserProfileBuvid = "/data/"
)

var (
	signParams = []string{"appKey", "timestamp", "version"}
)

// QueryPlayDaily get video play rank list from berserker
func (d *Dao) QueryPlayDaily(c context.Context, date string) (vlist []*model.VideoHiveInfo, err error) {
	v := make(url.Values, 8)
	query := fmt.Sprintf(_queryJSON, date)
	v.Set("query", query)
	var res struct {
		Code   int                   `json:"code"`
		Result []model.VideoHiveInfo `json:"result"`
	}

	if err = d.doHTTPGet(c, d.c.Berserker.API.Rankdaily, "", v, d.c.Berserker.Key.YYC, _noKickUserAgent, &res); err != nil {
		log.Error("d.doHTTPGet err[%v]", err)
		return
	}
	if res.Code != 200 || len(res.Result) == 0 {
		err = ecode.NothingFound
		log.Warn("Berserker return err, url:%s;res:%d", d.c.Berserker.API.Rankdaily+"?"+v.Encode(), res.Code)
		return
	}
	for _, info := range res.Result {
		i := info
		vlist = append(vlist, &i)
	}
	return
}

//QueryOperaVideo query operation video once
func (d *Dao) QueryOperaVideo(c context.Context, date string, ch chan<- *model.VideoHiveInfo) (err error) {
	i := int64(0)
	var mid int64
	for {
		v := make(url.Values, 8)
		var res struct {
			Code   int                   `json:"code"`
			Result []model.VideoHiveInfo `json:"result"`
		}
		query := fmt.Sprintf(_queryJSONOper, date, i)
		v.Set("query", query)
		if err = d.doHTTPGet(c, d.c.Berserker.API.Operaonce, "", v, d.c.Berserker.Key.LZQ, _lzqUserAgent, &res); err != nil {
			log.Error("d.doHTTPGet err[%v]", err)
			return
		}
		if res.Code == 200 && len(res.Result) == 0 {
			return
		}
		if res.Code != 200 {
			err = ecode.NothingFound
			log.Warn("Berserker return err, url:%s;res:%d", d.c.Berserker.API.Operaonce+"?"+v.Encode(), res.Code)
			return
		}

		for _, info := range res.Result {
			ch <- &info
			mid = info.CID
		}
		i = mid
	}
}

//QueryUserBasic ...
func (d *Dao) QueryUserBasic(c context.Context) (jobURL string, err error) {

	v := make(url.Values, 8)
	var res struct {
		Code   int      `json:"code"`
		Msg    string   `json:"msg"`
		Result []string `json:"result"`
	}
	query := "{}"
	v.Set("query", query)
	if err = d.doHTTPGet(c, d.c.Berserker.API.Userbasic, "", v, d.c.Berserker.Key.LZQ, _lzqUserAgent, &res); err != nil {
		log.Error("d.doHTTPGet err[%v]", err)
		return
	}
	for i, file := range res.Result {
		query = fmt.Sprintf("{\"fileSuffix\": \"%s\"}", file)
		v.Set("query", query)
		bs, err := d.doHTTPGetRaw(c, d.c.Berserker.API.Userbasic, "", v, d.c.Berserker.Key.LZQ, _lzqUserAgent, &res)
		if err != nil {
			log.Error("d.doHTTPGet err[%v]", err)
		} else {
			fileName := fmt.Sprintf("/data/basic_profile/part_%d", i)
			if ioutil.WriteFile(fileName, bs, 0644) == nil {
				log.Info("write file success")
			} else {
				log.Error("write file error(%v)", err)
			}
		}

	}

	return
}

//UserProfileGet ...
func (d *Dao) UserProfileGet(c context.Context) (jobURL []string, err error) {
	//
	v := make(url.Values, 8)
	var res struct {
		Code   int      `json:"code"`
		Msg    string   `json:"msg"`
		Result []string `json:"result"`
	}
	query := "{}"
	v.Set("query", query)
	if err = d.doHTTPGet(c, d.c.Berserker.API.UserProfile, "", v, d.c.Berserker.Key.HM, _chmUserAgent, &res); err != nil {
		log.Error("d.doHTTPGet err[%v]", err)
		return
	}

	for i, file := range res.Result {
		query = fmt.Sprintf("{\"fileSuffix\": \"/%s\"}", file)
		//fmt.Printf("query: %v\n", query)
		v.Set("query", query)

		time.Sleep(3 * time.Second)
		var bs []byte
		bs, err = d.doHTTPGetRaw(c, d.c.Berserker.API.UserProfile, "", v, d.c.Berserker.Key.HM, _chmUserAgent, &res)

		if err != nil {
			log.Error("d.doHTTPGet err[%v]", err)
		} else {
			fileName := fmt.Sprintf(_basePathUserProfile+"part_%d", i)
			if ioutil.WriteFile(fileName, bs, 0644) == nil {
				log.Info("write file success")
			} else {
				log.Error("write file error(%v)", err)
			}
			d.ReadLine(fmt.Sprintf(_basePathUserProfile+"part_%d", i), d.HandlerUserBbqDmg)
			os.RemoveAll(fmt.Sprintf(_basePathUserProfile+"part_%d", i))
		}

	}

	time.Sleep(3 * time.Second)
	v2 := make(url.Values, 8)
	var res2 struct {
		Code   int      `json:"code"`
		Msg    string   `json:"msg"`
		Result []string `json:"result"`
	}

	query2 := "{}"
	v2.Set("query2", query2)
	if err = d.doHTTPGet(c, d.c.Berserker.API.UserProfileBuvid, "", v2, d.c.Berserker.Key.HM, _chmUserAgent, &res2); err != nil {
		log.Error("d.doHTTPGet err[%v]", err)
		return
	}

	for i, file := range res2.Result {
		query2 = fmt.Sprintf("{\"fileSuffix\": \"/%s\"}", file)
		//fmt.Printf("query: %v\n", query)
		v2.Set("query", query2)

		time.Sleep(3 * time.Second)
		bs, err := d.doHTTPGetRaw(c, d.c.Berserker.API.UserProfileBuvid, "", v2, d.c.Berserker.Key.HM, _chmUserAgent, &res2)

		if err != nil {
			log.Error("d.doHTTPGet err[%v]", err)
		} else {
			fileName := fmt.Sprintf(_basePathUserProfileBuvid+"part_%d", i)
			if ioutil.WriteFile(fileName, bs, 0644) == nil {
				log.Info("write file success")
			} else {
				log.Error("write file error(%v)", err)
			}
			d.ReadLine(fmt.Sprintf(_basePathUserProfileBuvid+"part_%d", i), d.HandlerUserBbqDmgBuvid)
			os.RemoveAll(fmt.Sprintf(_basePathUserProfileBuvid+"part_%d", i))
		}

	}
	return
}

// doHttpRequest make a http request for data platform api
func (d *Dao) doHTTPGet(c context.Context, uri, realIP string, params url.Values, key *conf.BerSerkerKey, userAgent string, res interface{}) (err error) {

	enc, err := d.berserkeSign(params, key)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s,params:%v", uri, params)
		return
	}
	if enc != "" {
		uri = uri + "?" + enc
	}
	req, err := xhttp.NewRequest(xhttp.MethodGet, uri, nil)

	fmt.Printf("Req: %s ", req.URL)
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", xhttp.MethodGet, uri)
		return
	}

	req.Header.Set(_userAgent, userAgent+" "+env.AppID)
	if err != nil {
		return
	}

	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}

	return d.HTTPClient.Do(c, req, res)
}

// doHTTPGetRaw make a http request for data platform api
func (d *Dao) doHTTPGetRaw(c context.Context, uri, realIP string, params url.Values, key *conf.BerSerkerKey, userAgent string, res interface{}) (bs []byte, err error) {
	enc, err := d.berserkeSign(params, key)
	if err != nil {
		err = pkgerr.Wrapf(err, "uri:%s,params:%v", uri, params)
		return
	}
	if enc != "" {
		uri = uri + "?" + enc
	}
	req, err := xhttp.NewRequest(xhttp.MethodGet, uri, nil)
	if err != nil {
		err = pkgerr.Wrapf(err, "method:%s,uri:%s", xhttp.MethodGet, uri)
		return
	}

	req.Header.Set(_userAgent, userAgent+" "+env.AppID)
	if err != nil {
		return
	}

	if realIP != "" {
		req.Header.Set(_httpHeaderRemoteIP, realIP)
	}

	return d.HTTPClient.Raw(c, req)
}

// Sign calc appkey and appsecret sign.
func (d *Dao) berserkeSign(params url.Values, key *conf.BerSerkerKey) (query string, err error) {
	params.Set("appKey", key.Appkey)
	params.Set("signMethod", "md5")
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	tmp := params.Encode()
	signTmp := d.encode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(key.Secret)
	b.WriteString(signTmp)
	b.WriteString(key.Secret)
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(strings.ToUpper(hex.EncodeToString(mh[:])))
	query = qb.String()
	return
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
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

// QueryUserDmg .
func (d *Dao) QueryUserDmg(c context.Context) (jobURL string, err error) {
	logDay := time.Now().AddDate(0, 0, -1).Format("20060102")
	params := url.Values{}
	params.Set("query", fmt.Sprintf(_userDmgQueryHive, logDay))
	var res struct {
		Code         int    `json:"code"`
		Msg          string `json:"msg"`
		JobStatusURL string `json:"jobStatusUrl"`
	}
	if err = d.doHTTPGet(c, d.c.Berserker.API.Userdmg, "", params, d.c.Berserker.Key.HSC, _hscUserAgent, &res); err != nil {
		return
	}
	if res.Code != 200 {
		log.Error("Berserker user_dmg err(%v)", err)
		return
	}
	jobURL = res.JobStatusURL
	return
}

// QueryJobStatus 查询hive脚本执行结果
func (d *Dao) QueryJobStatus(c context.Context, jobURL string) (urls []string, err error) {
	var res struct {
		Code      int      `json:"code"`
		Msg       string   `json:"msg"`
		StatusID  int      `json:"statusId"`
		StatusMsg string   `json:"statusMsg"`
		HdfsPath  []string `json:"hdfsPath"`
	}
	req, err := xhttp.NewRequest(xhttp.MethodGet, jobURL, nil)
	if err != nil {
		log.Error("QueryJobStatus NewRequest, err(%v)", err)
		return
	}
	for {
		if err = d.HTTPClient.Do(c, req, &res); err != nil {
			log.Error("QueryJobStatus do get failed, joburl(%v), err(%v)", jobURL, err)
			return
		}
		if res.Code != 200 {
			log.Error("QueryJobStatus http code error, joburl(%v), err(%v)", jobURL, err)
			return
		}
		if res.StatusID == _jobStatusDoing || res.StatusID == _jobStatusWaiting {
			//等待1min
			log.Info("QueryJobStatus got job status %v, joburl(%v)", res.StatusID, jobURL)
			time.Sleep(60 * time.Second)
			continue
		}
		if res.StatusID == _jobStatusFailed {
			log.Error("QueryJobStatus got job status failed joburl(%v), err(%v)", jobURL, err)
			return
		}
		if res.StatusID == _jobStatusSuccess {
			log.Info("QueryJobStatus got job status success joburl(%v), err(%v)", jobURL, err)
			urls = res.HdfsPath
			return
		}
		if res.StatusID != _jobStatusSuccess && res.StatusID != _jobStatusFailed && res.StatusID != _jobStatusDoing && res.StatusID != _jobStatusWaiting {
			log.Error("QueryJobStatus got wrong job status status(%v), joburl(%v)", res.StatusID, jobURL)
			return
		}
	}
}

//QueryUpUserDmg .
func (d *Dao) QueryUpUserDmg(c context.Context, mid int64) (upUserDmg []*model.UpUserDmg, err error) {
	params := url.Values{}
	params.Set("query", fmt.Sprintf(_upUserDmgQueryJSON, mid))
	var res struct {
		Code   int                `json:"code"`
		Result []*model.UpUserDmg `json:"result"`
	}

	if err = d.doHTTPGet(c, d.c.Berserker.API.Upuserdmg, "", params, d.c.Berserker.Key.HSC, _hscUserAgent, &res); err != nil {
		return
	}
	if res.Code != 200 {
		log.Error("Berserker up_user_dmg err(%v)", err)
		return
	}
	upUserDmg = res.Result
	return
}

//QueryUpMid .发起hive查询，取粉丝数大于1万的up mid
func (d *Dao) QueryUpMid(c context.Context, date string) (jobURL string, err error) {
	params := url.Values{}
	params.Set("query", fmt.Sprintf(_upMidQueryHive, date))
	var res struct {
		Code         int    `json:"code"`
		Msg          string `json:"msg"`
		JobStatusURL string `json:"jobStatusUrl"`
	}
	if err = d.doHTTPGet(c, d.c.Berserker.API.Upmid, "", params, d.c.Berserker.Key.LJ, _ljUserAgent, &res); err != nil {
		log.Error("hive QueryUpMid failed, err(%v)", err)
		return
	}
	if res.Code != 200 {
		fmt.Println(res.Code)
		log.Error("hive QueryUpMid failed, err(%v), httpcode(%v)", err, res.Code)
		return
	}
	jobURL = res.JobStatusURL
	fmt.Println(jobURL)
	return
}
