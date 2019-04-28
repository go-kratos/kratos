package bfs

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/library/ecode"
	xhttp "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

const (
	_appKey = "appkey"
	_ts     = "ts"
)

var (
	_defaultHost             = "http://api.bilibili.co"
	_pathUpload              = "/x/internal/upload"
	_pathGenWatermark        = "/x/internal/image/gen"
	_defaultHTTPClientConfig = &xhttp.ClientConfig{
		App: &xhttp.App{
			Key:    "3c4e41f926e51656",
			Secret: "26a2095b60c24154521d24ae62b885bb",
		},
		Dial:    xtime.Duration(1 * time.Second),
		Timeout: xtime.Duration(1 * time.Second),
	}

	errBucket = errors.New("bucket is empty")
	errFile   = errors.New("file is empty")
)

// Config bfs upload config
type Config struct {
	Host       string
	HTTPClient *xhttp.ClientConfig
}

// Request bfs upload request
type Request struct {
	Bucket      string
	Dir         string
	ContentType string
	Filename    string
	File        []byte
	WMKey       string
	WMText      string
	WMPaddingX  uint32
	WMPaddingY  uint32
	WMScale     float64
}

// verify verify bfs request.
func (r *Request) verify() error {
	if r.Bucket == "" {
		return errBucket
	}
	if r.File == nil || len(r.File) == 0 {
		return errFile
	}
	return nil
}

// BFS bfs instance
type BFS struct {
	conf   *Config
	client *xhttp.Client
}

// New new a bfs client.
func New(c *Config) *BFS {
	if c == nil {
		c = &Config{
			Host:       _defaultHost,
			HTTPClient: _defaultHTTPClientConfig,
		}
	}
	return &BFS{
		conf:   c,
		client: xhttp.NewClient(c.HTTPClient),
	}
}

// Upload bfs internal upload.
func (b *BFS) Upload(ctx context.Context, req *Request) (location string, err error) {
	var (
		buf      = &bytes.Buffer{}
		bw       io.Writer
		request  *http.Request
		response *struct {
			Code int `json:"code"`
			Data struct {
				ETag     string `json:"etag"`
				Location string `json:"location"`
			} `json:"data"`
			Message string `json:"message"`
		}
		url = b.conf.Host + _pathUpload + "?" + b.sign()
	)
	if err = req.verify(); err != nil {
		return
	}
	w := multipart.NewWriter(buf)
	if bw, err = w.CreateFormFile("file", "bfs-upload"); err != nil {
		return
	}
	if _, err = bw.Write(req.File); err != nil {
		return
	}
	w.WriteField("bucket", req.Bucket)
	if req.Filename != "" {
		w.WriteField("file_name", req.Filename)
	}
	if req.Dir != "" {
		w.WriteField("dir", req.Dir)
	}
	if req.WMText != "" {
		w.WriteField("wm_text", req.WMKey)
	}
	if req.WMKey != "" {
		w.WriteField("wm_key", req.WMKey)
	}
	if req.WMPaddingX > 0 {
		w.WriteField("wm_padding_x", fmt.Sprint(req.WMPaddingX))
	}
	if req.WMPaddingY > 0 {
		w.WriteField("wm_padding_y", fmt.Sprint(req.WMPaddingY))
	}
	if req.WMScale > 0 {
		w.WriteField("wm_scale", strconv.FormatFloat(req.WMScale, 'f', 2, 64))
	}
	if req.ContentType != "" {
		w.WriteField("content_type", req.ContentType)
	}
	if err = w.Close(); err != nil {
		return
	}
	if request, err = http.NewRequest(http.MethodPost, url, buf); err != nil {
		return
	}
	request.Header.Set("Content-Type", w.FormDataContentType())
	if err = b.client.Do(ctx, request, &response); err != nil {
		return
	}
	if !ecode.Int(response.Code).Equal(ecode.OK) {
		err = ecode.Int(response.Code)
		return
	}
	location = response.Data.Location
	return
}

// GenWatermark create watermark image by key and text.
func (b *BFS) GenWatermark(ctx context.Context, uploadKey, wmKey, wmText string, vertical bool, distance int) (location string, err error) {
	var (
		params   = url.Values{}
		uri      = b.conf.Host + _pathGenWatermark
		response *struct {
			Code int `json:"code"`
			Data struct {
				Location string `json:"location"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				MD5      string `json:"md5"`
			} `json:"data"`
			Message string `json:"message"`
		}
	)
	params.Set("upload_key", uploadKey)
	params.Set("wm_key", wmKey)
	params.Set("wm_text", wmText)
	params.Set("wm_vertical", fmt.Sprint(vertical))
	params.Set("distance", fmt.Sprint(distance))
	if err = b.client.Post(ctx, uri, "", params, &response); err != nil {
		return
	}
	if !ecode.Int(response.Code).Equal(ecode.OK) {
		err = ecode.Int(response.Code)
		return
	}
	location = response.Data.Location
	return
}

// sign calc appkey and appsecret sign.
func (b *BFS) sign() string {
	key := b.conf.HTTPClient.Key
	secret := b.conf.HTTPClient.Secret
	params := url.Values{}
	params.Set(_appKey, key)
	params.Set(_ts, strconv.FormatInt(time.Now().Unix(), 10))
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var buf bytes.Buffer
	buf.WriteString(tmp)
	buf.WriteString(secret)
	mh := md5.Sum(buf.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(hex.EncodeToString(mh[:]))
	return qb.String()
}
