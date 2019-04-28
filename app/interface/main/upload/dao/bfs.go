package dao

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/upload/conf"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Bfs .
type Bfs struct {
	bfsURL         string
	waterMarkURL   string
	imageGenURL    string
	client         *http.Client
	wmClient       *http.Client
	imageGenClient *http.Client
}

// NewBfs .
func NewBfs(c *conf.Config) *Bfs {
	return &Bfs{
		bfsURL:       c.Bfs.BfsURL,
		waterMarkURL: c.Bfs.WaterMarkURL,
		imageGenURL:  c.Bfs.ImageGenURL,
		client: &http.Client{
			Timeout: time.Duration(c.Bfs.TimeOut),
		},
		wmClient: &http.Client{
			Timeout: time.Duration(c.Bfs.WmTimeOut),
		},
		imageGenClient: &http.Client{
			Timeout: time.Duration(c.Bfs.ImageGenTimeOut),
		},
	}
}

// GenImage gen a image.
func (b *Bfs) GenImage(ctx context.Context, wmKey, wmText string, distance int, vertical bool) (res []byte, height, width int, hasher string, err error) {
	var (
		params = url.Values{}
		resp   *http.Response
	)
	params.Set("wm_key", wmKey)
	params.Set("wm_text", wmText)
	params.Set("distance", strconv.Itoa(distance))
	params.Set("vertical", strconv.FormatBool(vertical))
	if resp, err = b.imageGenClient.PostForm(b.imageGenURL, params); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("bfs.genImage.status(%v)", resp.StatusCode)
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		err = ecode.RequestErr
		return
	case http.StatusMethodNotAllowed:
		err = ecode.MethodNotAllowed
		return
	case http.StatusInternalServerError:
		err = ecode.ServerErr
		return
	default:
		err = ecode.ServerErr
		return
	}
	if height, err = strconv.Atoi(resp.Header.Get("O-Height")); err != nil {
		log.Error("bfs.genImage.O-Height,err(%v)", err)
		return
	}
	if width, err = strconv.Atoi(resp.Header.Get("O-Width")); err != nil {
		log.Error("bfs.genImage.O-Width,err(%v)", err)
		return
	}
	hasher = resp.Header.Get("Md5")
	res, err = ioutil.ReadAll(resp.Body)
	return
}

// Watermark generate watermark image by key and text.
func (b *Bfs) Watermark(ctx context.Context, data []byte, contentType, wmKey, wmText string, paddingX, paddingY int, wmScale float64) (res []byte, err error) {
	var (
		resp *http.Response
		bw   io.Writer
		req  *http.Request
		ext  string
	)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	if bw, err = w.CreateFormFile("file", "1.jpg"); err != nil {
		return
	}
	if _, err = bw.Write(data); err != nil {
		return
	}
	w.WriteField("wm_key", wmKey)
	w.WriteField("wm_text", wmText)
	w.WriteField("padding_x", strconv.Itoa(paddingX))
	w.WriteField("padding_y", strconv.Itoa(paddingY))
	w.WriteField("wm_scale", strconv.FormatFloat(wmScale, 'f', 2, 64))
	if ext = path.Base(contentType); ext == "jpeg" {
		ext = "jpg"
	}
	w.WriteField("ext", fmt.Sprintf(".%s", ext))
	if err = w.Close(); err != nil {
		return
	}
	if req, err = http.NewRequest(http.MethodPost, b.waterMarkURL, buf); err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if resp, err = b.wmClient.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("bfs.waterMark.status(%v)", resp.StatusCode)
		err = ecode.ServerErr
		return
	}
	res, err = ioutil.ReadAll(resp.Body)
	return
}

// Upload upload file to bfs.
func (b *Bfs) Upload(ctx context.Context, key, secret, contentType, bucket, dir, filename string, data []byte) (location, etag string, err error) {
	dir = strings.Trim(dir, "/")
	var (
		buf   = new(bytes.Buffer)
		fname = path.Join(dir, filename)
	)
	// 多级目录上传，request url必须以 / 结尾
	if filename == "" && dir != "" {
		fname = fname + "/"
	}
	reqURL := fmt.Sprintf("http://%s/bfs/%s/%s", b.bfsURL, bucket, fname)
	if _, err = buf.Write(data); err != nil {
		log.Error("Upload.buf.Write.error(%v)", err)
		err = ecode.ServerErr
		return
	}
	req, err := http.NewRequest(http.MethodPut, reqURL, buf)
	if err != nil {
		err = ecode.ServerErr
		return
	}
	auth := b.authorize(key, secret, http.MethodPut, bucket, fname, time.Now().Unix())
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", auth)
	resp, err := b.client.Do(req)
	if err != nil {
		log.Error("Upload.client.Do.error(%v)", err)
		err = ecode.ServerErr
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Error("bucket:%s,filename:%s,response code:%+v", bucket, filename, resp.StatusCode)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		err = ecode.RequestErr
		return
	case http.StatusUnauthorized:
		// 验证不通过
		err = ecode.BfsUploadAuthErr
		return
	case http.StatusRequestEntityTooLarge:
		err = ecode.FileTooLarge
		return
	case http.StatusNotFound:
		err = ecode.NothingFound
		return
	case http.StatusMethodNotAllowed:
		err = ecode.MethodNotAllowed
		return
	case http.StatusServiceUnavailable:
		err = ecode.BfsUploadServiceUnavailable
		return
	case http.StatusInternalServerError:
		err = ecode.ServerErr
		return
	default:
		err = ecode.BfsUploadStatusErr
		return
	}
	code, err := strconv.Atoi(resp.Header.Get("Code"))
	if err != nil || code != 200 {
		log.Error("bucket:%s,dir:%s,filename:%s response.Header.Code(%s) error(%v)", bucket, dir, filename, resp.Header.Get("Code"), err)
		err = ecode.BfsUploadCodeErr
		return
	}
	location = resp.Header.Get("location")
	etag = resp.Header.Get("etag")
	return
}

// authorize generate authorization
func (b *Bfs) authorize(key, secret, method, bucket, fileName string, expire int64) string {
	content := fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, fileName, expire)
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s:%s:%d", key, signature, expire)
}
