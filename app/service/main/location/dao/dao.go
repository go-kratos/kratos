package dao

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"go-common/app/service/main/location/conf"
	"go-common/app/service/main/location/model"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

const (
	_anonym  = "/app/geoip_download"
	_checkip = "/ipip/checkipipnetversion"
)

// Dao dao.
type Dao struct {
	// mysql
	c       *conf.Config
	db      *sql.DB
	client  *http.Client
	client2 *httpx.Client
	// redis
	redis  *redis.Pool
	expire int32
	// host
	anonym  string
	checkip string
}

// New new a dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:      c,
		db:     sql.NewMySQL(c.DB.Zlimit),
		redis:  redis.NewPool(c.Redis.Zlimit.Config),
		expire: int32(time.Duration(c.Redis.Zlimit.Expire) / time.Second),
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		anonym:  c.Host.Maxmind + _anonym,
		checkip: c.Host.Bvcip + _checkip,
		client2: httpx.NewClient(c.HTTPClient),
	}
	return
}

// Ping ping a dao.
func (d *Dao) Ping(c context.Context) (err error) {
	conn := d.redis.Get(c)
	_, err = conn.Do("SET", "PING", "PONG")
	conn.Close()
	return
}

// Close close a dao.
func (d *Dao) Close() (err error) {
	if d.redis != nil {
		if err1 := d.redis.Close(); err1 != nil {
			err = err1
		}
	}
	if d.db != nil {
		if err1 := d.db.Close(); err1 != nil {
			err = err1
		}
	}
	return
}

// DownloadAnonym download anonym file.
func (d *Dao) DownloadAnonym() (err error) {
	// get file
	var (
		req  *http.Request
		resp *http.Response
	)
	params := url.Values{}
	params.Set("edition_id", "GeoIP2-Anonymous-IP")
	params.Set("date", "")
	params.Set("license_key", d.c.AnonymKey)
	params.Set("suffix", "tar.gz")
	enc := params.Encode()
	if req, err = http.NewRequest(http.MethodGet, d.anonym+"?"+enc, nil); err != nil {
		log.Error("http.NewRequest(%v) error(%v)", d.anonym+"?"+enc, err)
		return
	}
	if resp, err = d.client.Do(req); err != nil {
		log.Error("d.client.Do error(%v)", d.anonym+"?"+enc, err)
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(fmt.Sprintf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, d.anonym, enc))
		log.Error("%v", err)
		return
	}
	defer resp.Body.Close()
	// get md5
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, resp.Body); err != nil {
		log.Error("io.Copy error(%v)", err)
		return
	}
	var md5Str string
	md5Bs := md5.Sum(buf.Bytes())
	if md5Str, err = d.downloadAnonymMd5(); err != nil {
		log.Error("downloadAnonymMd5 error(%v)", err)
		return
	}
	if md5Str != hex.EncodeToString(md5Bs[:]) {
		err = errors.New("md5 not matched")
		return
	}
	var gr *gzip.Reader
	if gr, err = gzip.NewReader(buf); err != nil {
		log.Error("gzip.NewReader error(%v)", err)
		return
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	var hdr *tar.Header
	for {
		if hdr, err = tr.Next(); err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				log.Error("DownloadAnonym error(%v)", err)
				return
			}
		}
		if strings.Index(hdr.Name, d.c.AnonymFileName) > -1 {
			var f *os.File
			downfile := path.Join(d.c.FilePath, d.c.AnonymFileName)
			if f, err = os.OpenFile(downfile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
				log.Error(" os.OpenFile(%v) error(%v)", downfile, err)
				return
			}
			defer f.Close()
			if _, err = io.Copy(f, tr); err != nil {
				log.Error("io.Copy error(%v)", err)
				return
			}
			break
		}
	}
	return
}

// downloadAnonymMd5 get anonym file md5.
func (d *Dao) downloadAnonymMd5() (md5 string, err error) {
	var (
		req  *http.Request
		resp *http.Response
		body []byte
	)
	params := url.Values{}
	params.Set("edition_id", "GeoIP2-Anonymous-IP")
	params.Set("date", "")
	params.Set("license_key", d.c.AnonymKey)
	params.Set("suffix", "tar.gz.md5")
	enc := params.Encode()
	if req, err = http.NewRequest(http.MethodGet, d.anonym+"?"+enc, nil); err != nil {
		log.Error("DownloadAnonym http.NewRequest(%v) error(%v)", d.anonym+"?"+enc, err)
		return
	}
	if resp, err = d.client.Do(req); err != nil {
		log.Error("DownloadAnonym d.client.Do error(%v)", d.anonym+"?"+enc, err)
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(fmt.Sprintf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, d.anonym, enc))
		log.Error("%v", err)
		return
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("DownloadAnonymMd5 ioutil.ReadAll() error(%v)", err)
		return
	}
	md5 = string(body)
	return
}

func (d *Dao) CheckVersion(c context.Context) (version *model.Version, err error) {
	var (
		req *http.Request
		ip  = metadata.String(c, metadata.RemoteIP)
	)
	if req, err = d.client2.NewRequest("GET", d.checkip, ip, nil); err != nil {
		log.Error("%v", err)
		return
	}
	var res struct {
		Code int            `json:"code"`
		Msg  string         `json:"msg"`
		Data *model.Version `json:"data"`
	}
	if err = d.client2.Do(c, req, &res); err != nil {
		return
	}
	if res.Code != ecode.OK.Code() {
		err = errors.New(fmt.Sprintf("checkVersion falid err_code(%v)", res.Code))
		return
	}
	version = res.Data
	return
}

func (d *Dao) DownIPLibrary(c context.Context, version, file string) (err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	if req, err = http.NewRequest(http.MethodGet, version, nil); err != nil {
		log.Error("http.NewRequest(%v) error(%v)", version, err)
		return
	}
	if resp, err = d.client.Do(req); err != nil {
		log.Error("d.client.Do error(%v)", version, err)
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(fmt.Sprintf("incorrect http status:%d url:%s", resp.StatusCode, version))
		log.Error("%v", err)
		return
	}
	defer resp.Body.Close()
	var f *os.File
	if f, err = os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err != nil {
		log.Error(" os.OpenFile(%v) error(%v)", file, err)
		return
	}
	defer f.Close()
	if _, err = io.Copy(f, resp.Body); err != nil {
		log.Error("io.Copy error(%v)", err)
	}
	return
}
