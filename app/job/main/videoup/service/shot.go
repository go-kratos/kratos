package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/redis"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	videoshotPath      = "http://bfs.bilibili.co/bfs/"
	videoshotBucket    = "videoshot"
	videoshotBfsKey    = "624480b94ec1e0eb"
	videoshotBfsSecret = "303c6cfa85cf2d550fa5050d123cd2"

	archiveBucket    = "archive"
	archiveBfsKey    = "8d4e593ba7555502"
	archiveBfsSecret = "0bdbd4c7caeeddf587c3c4daec0475"
)

var (
	videoshotHTTPClient = &http.Client{Timeout: time.Second * 5}
)

// videoshotSHConsumer is videoshot consumer in shanghai 云立方.
func (s *Service) videoshotSHConsumer() {
	defer s.wg.Done()
	var (
		msgs = s.videoshotSub2.Messages()
	)
	for {
		msg, ok := <-msgs
		if !ok {
			log.Error("s.videoshotSub2.Message closed")
			return
		}
		msg.Commit()
		var res struct {
			Cid   int64    `json:"cid"`
			Image []string `json:"image"`
			Bin   string   `json:"bin"`
		}
		if err := json.Unmarshal([]byte(msg.Value), &res); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}
		log.Info("videoshotSHConsumer videoshotMessage in shanghai yunlifang key(%s) cid(%d) count(%d) partition(%d) offset(%d) commit", msg.Key, res.Cid, len(res.Image), msg.Partition, msg.Offset)
		s.videoshotAdd(res.Cid, res.Bin, res.Image)
	}
}

func (s *Service) videocoverCopy(aid, count int64, a *archive.Archive) {
	if a == nil {
		log.Warn("archive(%d) is not exist", aid)
		return
	}
	cover := a.Cover
	_, err := url.ParseRequestURI(cover)
	if err != nil {
		log.Error("videocoverCopy aid(%s)  cover(%s) parse error(%v)", aid, cover, err)
		return
	}
	count = count + 1
	_, bs, err := s.videoshotDown(cover)
	if err != nil {
		log.Error("videocoverCopy aid(%d) cover(%s) videoshotDown err(%v)", aid, cover, err)
		if err == ecode.NothingFound {
			return
		}
		s.syncRetry(context.TODO(), aid, count, redis.ActionForVideocovers, cover, cover)
		return
	}
	if cover, err = s.videocoverUp("", bs); err != nil {
		log.Error("videocoverCopy aid(%d) cover(%s) videocoverUp err(%v)", aid, cover, err)
		s.syncRetry(context.TODO(), aid, count, redis.ActionForVideocovers, a.Cover, a.Cover)
		return
	}
	u, err := url.Parse(cover)
	if err != nil {
		log.Error("videocoverCopy aid(%d) url(%s) Parse(%s) error", aid, cover, err)
		return
	}
	if _, err := s.arc.UpCover(context.TODO(), aid, u.Path); err != nil {
		s.syncRetry(context.TODO(), aid, count, redis.ActionForVideocovers, cover, cover)
		return
	}
	log.Info("videocoverCopy aid(%d) new cover(%s) sub(%s) success", aid, cover, u.Path)
}

func (s *Service) videoshotAdd(cid int64, bin string, imgs []string) {
	_, err := url.ParseRequestURI(bin)
	if err != nil {
		log.Error("videoshotAdd cid(%s) add bin(%s) parse error(%v)", cid, bin, err)
		return
	}
	fn, bs, err := s.videoshotDown(bin)
	if err != nil {
		if err == ecode.NothingFound {
			return
		}
		s.syncRetry(context.TODO(), cid, 0, redis.ActionForVideoshot, bin, strings.Join(imgs, ","))
		return
	}
	if _, err := s.videoshotUp(fn, bs); err != nil {
		s.syncRetry(context.TODO(), cid, 0, redis.ActionForVideoshot, bin, strings.Join(imgs, ","))
		return
	}
	for _, img := range imgs {
		if len(img) == 0 {
			continue
		}
		fn, bs, err := s.videoshotDown(img)
		if err != nil {
			if err == ecode.NothingFound {
				return
			}
			s.syncRetry(context.TODO(), cid, 0, redis.ActionForVideoshot, bin, strings.Join(imgs, ","))
			return
		}
		if _, err := s.videoshotUp(fn, bs); err != nil {
			s.syncRetry(context.TODO(), cid, 0, redis.ActionForVideoshot, bin, strings.Join(imgs, ","))
			return
		}
	}
	if _, err := s.arc.AddVideoShot(context.TODO(), cid, len(imgs)); err != nil {
		s.syncRetry(context.TODO(), cid, 0, "", bin, strings.Join(imgs, ","))
		return
	}
	log.Info("videoshotAdd cid(%d) success", cid)
}

func (s *Service) videoshotDown(uri string) (filename string, bs []byte, err error) {
	fileURL, err := url.ParseRequestURI(uri)
	if err != nil {
		log.Error("videoshotDown(%s) url parse error(%v)", uri, err)
		return
	}
	filename = filepath.Base(fileURL.Path)

	resp, err := videoshotHTTPClient.Get(uri)
	if err != nil {
		log.Error("videoshotDown download from url(%s) error(%v)", uri, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("videoshowDown(%s) status error(%d)", uri, resp.StatusCode)
		log.Error("%v", err)
		if resp.StatusCode == http.StatusNotFound {
			err = ecode.NothingFound
		}
		return
	}
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("videoshowDown(%s) read body error(%v)", uri, err)
	}
	return
}

func (s *Service) videoshotUp(filename string, bs []byte) (url string, err error) {
	var (
		uri  = fmt.Sprintf("%s%s/%s", videoshotPath, videoshotBucket, filename)
		req  *http.Request
		resp *http.Response
		code int
	)
	if req, err = http.NewRequest(http.MethodPut, uri, bytes.NewReader(bs)); err != nil {
		return
	}
	var sign = func() string {
		expire := time.Now().Unix()
		content := fmt.Sprintf("%s\n%s\n%s\n%d\n", http.MethodPut, videoshotBucket, filename, expire)
		mac := hmac.New(sha1.New, []byte(videoshotBfsSecret))
		mac.Write([]byte(content))
		return fmt.Sprintf("%s:%s:%d", videoshotBfsKey, base64.StdEncoding.EncodeToString(mac.Sum(nil)), expire)
	}
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Authorization", sign())
	if resp, err = videoshotHTTPClient.Do(req); err != nil {
		log.Error("videoshotUp client.Do(%s) error(%v)", filename, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("videoshotUp client.Do(%s) status: %d", filename, resp.StatusCode)
		log.Error("%v", err)
		return
	}
	if code, err = strconv.Atoi(resp.Header.Get("code")); err != nil {
		err = fmt.Errorf("videoshotUp conv(%s) code to int error(%v)", filename, err)
		log.Error("%v", err)
		return
	}
	if code == http.StatusOK {
		url = resp.Header.Get("location")
	} else {
		err = fmt.Errorf("videoshotUp client.Do(%s) code: %d", filename, code)
		log.Error("%v", err)
	}
	return
}

func (s *Service) videocoverUp(filename string, bs []byte) (url string, err error) {
	var (
		uri  = fmt.Sprintf("%s%s/%s", videoshotPath, archiveBucket, filename)
		req  *http.Request
		resp *http.Response
		code int
	)
	if req, err = http.NewRequest(http.MethodPut, uri, bytes.NewReader(bs)); err != nil {
		return
	}
	var sign = func() string {
		expire := time.Now().Unix()
		content := fmt.Sprintf("%s\n%s\n%s\n%d\n", http.MethodPut, archiveBucket, filename, expire)
		mac := hmac.New(sha1.New, []byte(archiveBfsSecret))
		mac.Write([]byte(content))
		return fmt.Sprintf("%s:%s:%d", archiveBfsKey, base64.StdEncoding.EncodeToString(mac.Sum(nil)), expire)
	}
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Authorization", sign())
	if resp, err = videoshotHTTPClient.Do(req); err != nil {
		log.Error("videocoverUp client.Do(%s) error(%v)", filename, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("videocoverUp client.Do(%s) status: %d", filename, resp.StatusCode)
		log.Error("%v", err)
		return
	}
	if code, err = strconv.Atoi(resp.Header.Get("code")); err != nil {
		err = fmt.Errorf("videocoverUp conv(%s) code to int error(%v)", filename, err)
		log.Error("%v", err)
		return
	}
	if code == http.StatusOK {
		url = resp.Header.Get("location")
	} else {
		err = fmt.Errorf("videocoverUp client.Do(%s) code: %d", filename, code)
		log.Error("%v", err)
	}
	return
}
