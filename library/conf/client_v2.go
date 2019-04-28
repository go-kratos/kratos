package conf

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// api
	_apiGet1    = "http://%s/config/v2/get?%s"
	_apiCheck1  = "http://%s/config/v2/check?%s"
	_apiCreate  = "http://%s/config/v2/create"
	_apiUpdate  = "http://%s/config/v2/update"
	_apiConfIng = "http://%s/config/v2/config/ing?%s"
)

type version1 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *ver   `json:"data"`
}

type ver struct {
	Version int64   `json:"version"`
	Diffs   []int64 `json:"diffs"`
}

type confIng struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *Value `json:"data"`
}

type res struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//Value value.
type Value struct {
	CID    int64  `json:"cid"`
	Name   string `json:"name"`
	Config string `json:"config"`
}

// Toml2 return config value.
func (c *Client) Toml2() (cf string, ok bool) {
	var (
		m   map[string]*Value
		val *Value
	)
	if m, ok = c.data.Load().(map[string]*Value); !ok {
		return
	}
	if val, ok = m[commonKey]; !ok {
		return
	}
	cf = val.Config
	return
}

// Value2 return config value.
func (c *Client) Value2(key string) (cf string, ok bool) {
	var (
		m   map[string]*Value
		val *Value
	)
	if m, ok = c.data.Load().(map[string]*Value); !ok {
		return
	}
	if val, ok = m[key]; !ok {
		return
	}
	cf = val.Config
	return
}

// init check local config is ok
func (c *Client) init2() (err error) {
	var v *ver
	c.data.Store(make(map[string]*Value))
	if v, err = c.checkVersion2(&ver{Version: _unknownVersion}); err != nil {
		fmt.Printf("get remote version error(%v)\n", err)
		return
	}
	for i := 0; i < 3; i++ {
		if v.Version == _unknownVersion {
			fmt.Println("get null version")
			return
		}
		if err = c.download2(v, true); err == nil {
			return
		}
		fmt.Printf("retry times: %d, c.download() error(%v)\n", i, err)
		time.Sleep(_retryInterval)
	}
	return
}

func (c *Client) updateproc2() (err error) {
	var ver *ver
	for {
		time.Sleep(_retryInterval)
		if ver, err = c.checkVersion2(c.diff); err != nil {
			log.Error("c.checkVersion(%d) error(%v)", c.ver, err)
			continue
		} else if ver.Version == c.diff.Version {
			continue
		}
		if err = c.download2(ver, false); err != nil {
			log.Error("c.download() error(%s)", err)
			continue
		}
	}
}

// download download config from config service
func (c *Client) download2(ver *ver, isFirst bool) (err error) {
	var (
		d             *data
		tmp           []*Value
		oConfs, confs map[string]*Value
		buf           = new(bytes.Buffer)
		ok            bool
	)
	if d, err = c.getConfig2(ver); err != nil {
		return
	}
	bs := []byte(d.Content)
	// md5 file
	if mh := md5.Sum(bs); hex.EncodeToString(mh[:]) != d.Md5 {
		err = fmt.Errorf("md5 mismatch, local:%s, remote:%s", hex.EncodeToString(mh[:]), d.Md5)
		return
	}

	// write conf
	if err = json.Unmarshal(bs, &tmp); err != nil {
		return
	}
	confs = make(map[string]*Value)
	if oConfs, ok = c.data.Load().(map[string]*Value); ok {
		for k, v := range oConfs {
			confs[k] = v
		}
	}
	for _, v := range tmp {
		if err = ioutil.WriteFile(path.Join(conf.Path, v.Name), []byte(v.Config), 0644); err != nil {
			return
		}
		confs[v.Name] = v
	}
	for _, v := range confs {
		if strings.Contains(v.Name, ".toml") {
			buf.WriteString(v.Config)
			buf.WriteString("\n")
		}
	}
	confs[commonKey] = &Value{Config: buf.String()}
	// update current version
	c.diff = ver
	c.data.Store(confs)
	if isFirst {
		return
	}
	for _, v := range tmp {
		if c.watchAll {
			c.event <- v.Name
			continue
		}
		if c.watchFile == nil {
			continue
		}
		if _, ok := c.watchFile[v.Name]; ok {
			c.event <- v.Name
		}
	}
	return
}

// poll config server
func (c *Client) checkVersion2(reqVer *ver) (ver *ver, err error) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		rb   []byte
	)
	if url, err = c.makeURL2(_apiCheck1, reqVer); err != nil {
		err = fmt.Errorf("checkVersion() c.makeUrl() error url empty")
		return
	}
	// http
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	if resp, err = c.httpCli.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("checkVersion() http error url(%s) status: %d", url, resp.StatusCode)
		return
	}
	// ok
	if rb, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	v := &version1{}
	if err = json.Unmarshal(rb, v); err != nil {
		return
	}
	switch v.Code {
	case _codeOk:
		if v.Data == nil {
			err = fmt.Errorf("checkVersion() response error result: %v", v)
			return
		}
		ver = v.Data
	case _codeNotModified:
		ver = reqVer
	default:
		err = fmt.Errorf("checkVersion() response error result: %v", v)
	}
	return
}

// updateVersion update config version
func (c *Client) getConfig2(ver *ver) (data *data, err error) {
	var (
		url  string
		req  *http.Request
		resp *http.Response
		rb   []byte
		res  = &result{}
	)
	if url, err = c.makeURL2(_apiGet1, ver); err != nil {
		err = fmt.Errorf("getConfig() c.makeUrl() error url empty")
		return
	}
	// http
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	if resp, err = c.httpCli.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	// ok
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("getConfig() http error url(%s) status: %d", url, resp.StatusCode)
		return
	}
	if rb, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(rb, res); err != nil {
		return
	}
	switch res.Code {
	case _codeOk:
		// has new config
		if res.Data == nil {
			err = fmt.Errorf("getConfig() response error result: %v", res)
			return
		}
		data = res.Data
	default:
		err = fmt.Errorf("getConfig() response error result: %v", res)
	}
	return
}

// makeUrl signed url
func (c *Client) makeURL2(api string, ver *ver) (query string, err error) {
	var ids []byte
	params := url.Values{}
	// service
	params.Set("service", service())
	params.Set("hostname", conf.Host)
	params.Set("build", conf.Ver)
	params.Set("version", fmt.Sprint(ver.Version))
	if ids, err = json.Marshal(ver.Diffs); err != nil {
		return
	}
	params.Set("ids", string(ids))
	params.Set("ip", localIP())
	params.Set("token", conf.Token)
	params.Set("appoint", conf.Appoint)
	params.Set("customize", c.customize)
	// api
	query = fmt.Sprintf(api, conf.Addr, params.Encode())
	return
}

//Create create.
func (c *Client) Create(name, content, operator, mark string) (err error) {
	var (
		resp *http.Response
		rb   []byte
		res  = &res{}
	)
	params := url.Values{}
	params.Set("service", service())
	params.Set("name", name)
	params.Set("content", content)
	params.Set("operator", operator)
	params.Set("mark", mark)
	params.Set("token", conf.Token)
	if resp, err = c.httpCli.PostForm(fmt.Sprintf(_apiCreate, conf.Addr), params); err != nil {
		return
	}
	defer resp.Body.Close()
	// ok
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Create() http error url(%s) status: %d", fmt.Sprintf(_apiCreate, conf.Addr), resp.StatusCode)
		return
	}
	if rb, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(rb, res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
	}
	return
}

//Update update.
func (c *Client) Update(ID int64, content, operator, mark string) (err error) {
	var (
		resp *http.Response
		rb   []byte
		res  = &result{}
	)
	params := url.Values{}
	params.Set("conf_id", fmt.Sprintf("%d", ID))
	params.Set("content", content)
	params.Set("operator", operator)
	params.Set("mark", mark)
	params.Set("service", service())
	params.Set("token", conf.Token)
	if resp, err = c.httpCli.PostForm(fmt.Sprintf(_apiUpdate, conf.Addr), params); err != nil {
		return
	}
	defer resp.Body.Close()
	// ok
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Update() http error url(%s) status: %d", fmt.Sprintf(_apiUpdate, conf.Addr), resp.StatusCode)
		return
	}
	if rb, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(rb, res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
	}
	return
}

//ConfIng confIng.
func (c *Client) ConfIng(name string) (v *Value, err error) {
	var (
		req  *http.Request
		resp *http.Response
		rb   []byte
		res  = &confIng{}
	)
	params := url.Values{}
	params.Set("name", name)
	params.Set("service", service())
	params.Set("token", conf.Token)
	// http
	if req, err = http.NewRequest("GET", fmt.Sprintf(_apiConfIng, conf.Addr, params.Encode()), nil); err != nil {
		return
	}
	if resp, err = c.httpCli.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	// ok
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("ConfIng() http error url(%s) status: %d", _apiCreate, resp.StatusCode)
		return
	}
	if rb, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(rb, res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = ecode.Int(res.Code)
		return
	}
	v = res.Data
	return
}

//Configs configs.
func (c *Client) Configs() (confs []*Value, ok bool) {
	var (
		m map[string]*Value
	)
	if m, ok = c.data.Load().(map[string]*Value); !ok {
		return
	}
	for _, v := range m {
		if v.CID == 0 {
			continue
		}
		confs = append(confs, v)
	}
	return
}

func service() string {
	return fmt.Sprintf("%s_%s_%s", conf.TreeID, conf.DeployEnv, conf.Zone)
}
