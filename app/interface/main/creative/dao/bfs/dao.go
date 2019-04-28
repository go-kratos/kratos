package bfs

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	nurl "net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_bucket = "archive"
	_url    = "http://bfs.bilibili.co/bfs/archive/"
	_method = "PUT"
	_key    = "8d4e593ba7555502"
	_secret = "0bdbd4c7caeeddf587c3c4daec0475"
)

var (
	errUpload   = errors.New("Upload failed")
	errDownload = errors.New("Download out image link failed")
)

// Dao is bfs dao.
type Dao struct {
	c          *conf.Config
	client     *http.Client
	captureCli *http.Client
}

// New new a bfs dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
		client: &http.Client{
			Timeout: time.Duration(c.BFS.Timeout),
		},
		captureCli: &http.Client{
			Timeout: time.Duration(time.Second),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return ecode.RequestErr
			},
		},
	}
	return d
}

// Upload upload bfs.
func (d *Dao) Upload(c context.Context, fileType string, bs []byte) (location string, err error) {
	req, err := http.NewRequest(d.c.BFS.Method, d.c.BFS.URL, bytes.NewBuffer(bs))
	if err != nil {
		log.Error("http.NewRequest error (%v) | fileType(%s)", err, fileType)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(d.c.BFS.Key, d.c.BFS.Secret, d.c.BFS.Method, d.c.BFS.Bucket, expire)
	req.Header.Set("Host", d.c.BFS.URL)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	ctx, cancel := context.WithTimeout(c, time.Duration(d.c.BFS.Timeout))
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("d.Client.Do error(%v) | url(%s)", err, d.c.BFS.URL)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, d.c.BFS.URL)
		err = errUpload
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("strconv.Itoa err, code(%s) | url(%s)", code, d.c.BFS.URL)
		err = errUpload
		return
	}
	location = header.Get("Location")
	return
}

// UploadArc upload bfs to archive bucket.
func (d *Dao) UploadArc(c context.Context, fileType string, body io.Reader) (location string, err error) {
	req, err := http.NewRequest(_method, _url, body)
	if err != nil {
		log.Error("http.NewRequest error (%v) | fileType(%s)", err, fileType)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(_key, _secret, _method, _bucket, expire)
	req.Header.Set("Host", _url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	c, cancel := context.WithTimeout(c, time.Duration(d.c.BFS.Timeout))
	req = req.WithContext(c)
	defer cancel()
	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("d.Client.Do error(%v) | url(%s)", err, _url)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, _url)
		err = errUpload
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("strconv.Itoa err, code(%s) | url(%s)", code, _url)
		err = errUpload
		return
	}
	location = header.Get("Location")
	return
}

// authorize returns authorization for upload file to bfs
func authorize(key, secret, method, bucket string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf("%s\n%s\n\n%d\n", method, bucket, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}

// UploadByFile upload local img file.
func (d *Dao) UploadByFile(c context.Context, imgpath string) (location string, err error) {
	data, err := ioutil.ReadFile(imgpath)
	if err != nil {
		log.Error("UploadByFile ioutil.ReadFile error (%v) | imgpath(%s)", err, imgpath)
		return
	}
	fileType := http.DetectContentType(data)
	if fileType != "image/jpeg" && fileType != "image/png" {
		log.Error("file type not allow file type(%s)", fileType)
		err = ecode.CreativeArticleImageTypeErr
	}
	body := new(bytes.Buffer)
	_, err = body.Write(data)
	if err != nil {
		log.Error("body.Write error (%v)", err)
		return
	}
	req, err := http.NewRequest(_method, _url, body)
	if err != nil {
		log.Error("http.NewRequest error (%v) | fileType(%s)", err, fileType)
		return
	}
	expire := time.Now().Unix()
	authorization := authorize(_key, _secret, _method, _bucket, expire)
	req.Header.Set("Host", _url)
	req.Header.Add("Date", fmt.Sprint(expire))
	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", fileType)
	// timeout
	c, cancel := context.WithTimeout(c, time.Duration(d.c.BFS.Timeout))
	req = req.WithContext(c)
	defer cancel()
	resp, err := d.client.Do(req)
	if err != nil {
		log.Error("d.Client.Do error(%v) | url(%s)", err, _url)
		err = ecode.BfsUploadServiceUnavailable
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("Upload http.StatusCode nq http.StatusOK (%d) | url(%s)", resp.StatusCode, _url)
		err = errUpload
		return
	}
	header := resp.Header
	code := header.Get("Code")
	if code != strconv.Itoa(http.StatusOK) {
		log.Error("strconv.Itoa err, code(%s) | url(%s)", code, _url)
		err = errUpload
		return
	}
	location = header.Get("Location")
	return
}

//Capture performs a HTTP Get request for the image url and upload bfs.
func (d *Dao) Capture(c context.Context, url string) (loc string, size int, err error) {
	if err = checkURL(url); err != nil {
		return
	}
	bs, ct, err := d.download(c, url)
	if err != nil {
		return
	}
	size = len(bs)
	if size == 0 {
		log.Error("capture image size(%d)|url(%s)", size, url)
		return
	}
	if ct != "image/jpeg" && ct != "image/jpg" && ct != "image/png" && ct != "image/gif" {
		log.Error("capture not allow image file type(%s)", ct)
		err = ecode.CreativeArticleImageTypeErr
		return
	}
	loc, err = d.Upload(c, ct, bs)
	return loc, size, err
}

func (d *Dao) download(c context.Context, url string) (bs []byte, ct string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("capture http.NewRequest error(%v)|url (%s)", err, url)
		return
	}
	// timeout
	ctx, cancel := context.WithTimeout(c, 800*time.Millisecond)
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.captureCli.Do(req)
	if err != nil {
		log.Error("capture d.client.Do error(%v)|url(%s)", err, url)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("capture http.StatusCode nq http.StatusOK(%d)|url(%s)", resp.StatusCode, url)
		err = errDownload
		return
	}
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("capture ioutil.ReadAll error(%v)", err)
		err = errDownload
		return
	}
	ct = http.DetectContentType(bs)
	return
}

func checkURL(url string) (err error) {
	// http || https
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		log.Error("capture url invalid(%s)", url)
		err = ecode.RequestErr
		return
	}
	u, err := nurl.Parse(url)
	if err != nil {
		log.Error("capture url.Parse error(%v)", err)
		err = ecode.RequestErr
		return
	}
	// make sure ip is public. avoid ssrf
	ips, err := net.LookupIP(u.Host) // take from 1st argument
	if err != nil {
		log.Error("capture url(%s) LookupIP failed", url)
		err = ecode.RequestErr
		return
	}
	if len(ips) == 0 {
		log.Error("capture url(%s) LookupIP length 0", url)
		err = ecode.RequestErr
		return
	}
	for _, v := range ips {
		if !isPublicIP(v) {
			log.Error("capture url(%s) is not public ip(%v)", url, v)
			err = ecode.RequestErr
			return
		}
	}
	return
}

func isPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
