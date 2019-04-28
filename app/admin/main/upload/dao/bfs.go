package dao

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"time"

	"go-common/app/admin/main/upload/conf"
	"go-common/app/admin/main/upload/model"
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

func (b *Bfs) waterMark(ctx context.Context, data []byte, contentType, wmKey, wmText string, paddingX, paddingY int, wmScale float64) (res []byte, err error) {
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
		if resp.StatusCode == http.StatusRequestEntityTooLarge {
			err = ecode.BfsUploadFileTooLarge
		} else {
			err = ecode.ServerErr
		}
		return
	}
	res, err = ioutil.ReadAll(resp.Body)
	return
}

// Upload with watermark
func (b *Bfs) Upload(ctx context.Context, up *model.UploadParam, data []byte) (location, etag string, err error) {
	var (
		res []byte
	)
	if up.WmKey != "" || up.WmText != "" {
		if res, err = b.waterMark(ctx, data, up.ContentType, up.WmKey, up.WmText, up.WmPaddingX, up.WmPaddingY, up.WmScale); err != nil {
			log.Error("b.waterMark(%+v) error(%v)", up, err)
			return
		}
		data = res
	}
	location, etag, err = b.upload(ctx, up.ContentType, up.Auth, up.Bucket, up.FileName, data)
	return
}

// SpecificUpload .
func (b *Bfs) SpecificUpload(ctx context.Context, contentType, auth, bucket, fileName string, data []byte) (location, etag string, err error) {
	location, etag, err = b.upload(ctx, contentType, auth, bucket, fileName, data)
	return
}

func (b *Bfs) upload(ctx context.Context, contentType, auth, bucket, fileName string, data []byte) (location, etag string, err error) {
	reqURL := fmt.Sprintf("http://%s/bfs/%s/%s", b.bfsURL, bucket, fileName)
	buf := new(bytes.Buffer)
	if _, err = buf.Write(data); err != nil {
		log.Error("buf.Write error(%v)", err)
		return
	}
	req, err := http.NewRequest("PUT", reqURL, buf)
	if err != nil {
		log.Error("http.NewRequest() error(%v)", err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", auth)
	resp, err := b.client.Do(req)
	if err != nil {
		log.Error("client.Do error(%v)", err)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		err = ecode.BfsRequestErr
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
	code, err := strconv.Atoi(resp.Header.Get("code"))
	if err != nil || code != 200 {
		err = ecode.BfsUploadCodeErr
		return
	}
	location = resp.Header.Get("location")
	etag = resp.Header.Get("etag")
	return
}

// Authorize .
func (b *Bfs) Authorize(key, secret, method, bucket, fileName string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, fileName, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}
