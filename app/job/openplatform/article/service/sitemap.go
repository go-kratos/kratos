package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/library/log"
)

type urlNode struct {
	id   int64
	next *urlNode
}

const (
	_sitemapSize     = 20000
	_sitemapInterval = int64(60)
)

func (s *Service) sitemapproc() {
	if s.urlListHead == nil {
		s.initURLList(context.TODO())
	}
	var interval = s.c.Sitemap.Interval
	if interval == 0 {
		interval = _sitemapInterval
	}
	for {
		if s.lastURLNode == nil || s.lastURLNode != s.urlListHead {
			s.updateXML()
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}

// add new url.
func (s *Service) addURLNode(c context.Context, id int64) {
	go s.pushShenma(c, id)
	if s.urlListHead == nil {
		s.initURLList(c)
	}
	if _, ok := s.sitemapMap[id]; ok {
		return
	}
	var (
		mutex sync.Mutex
		empty struct{}
	)
	mutex.Lock()
	node := &urlNode{id: id}
	s.urlListTail.next = node
	s.urlListTail = node
	oldID := s.urlListHead.id
	delete(s.sitemapMap, oldID)
	s.sitemapMap[id] = empty
	s.urlListHead = s.urlListHead.next
	mutex.Unlock()
}

func (s *Service) initURLList(c context.Context) {
	var (
		ids        []int64
		err        error
		size       = s.c.Sitemap.Size
		head, tail *urlNode
		empty      struct{}
		mutex      sync.Mutex
	)
	if size == 0 {
		size = _sitemapSize
	}
	if ids, err = s.dao.LastModIDs(c, size); err != nil {
		log.Error("s.initURLList error(%+v) size(%d)", err, size)
		return
	}
	for i := len(ids); i > 0; i-- {
		id := ids[i-1]
		node := &urlNode{id: id}
		if head == nil {
			head = node
			tail = node
			s.sitemapMap[id] = empty
			continue
		}
		tail.next = node
		tail = node
		s.sitemapMap[id] = empty
	}
	if head != nil && tail != nil {
		mutex.Lock()
		s.urlListHead = head
		s.urlListTail = tail
		mutex.Unlock()
	}
}

func (s *Service) updateXML() {
	var (
		node  = s.urlListHead
		xml   = ""
		mutex sync.Mutex
	)
	for node != nil {
		xml = `<url><loc>http://www.bilibili.com/read/cv` + strconv.FormatInt(node.id, 10) + `/</loc><lastmod>` + time.Now().Format("2006-01-02") + `</lastmod><changefreq>always</changefreq></url>` + xml
		node = node.next
	}
	xml = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.google.com/schemas/sitemap/0.84">` + xml + `</urlset>`
	mutex.Lock()
	s.sitemapXML = xml
	s.lastURLNode = s.urlListHead
	mutex.Unlock()
}

// SitemapXML .
func (s *Service) SitemapXML(c context.Context) (string, string) {
	return s.sitemapXML, strconv.Itoa(len(s.sitemapXML))
}

func (s *Service) pushShenma(c context.Context, id int64) {
	url := `http://data.zhanzhang.sm.cn/push?site=www.bilibili.com&user_name=zjdxgxy@163.com&resource_name=mip_add&token=TI_a1e0216ad1a64b63e9bbadfc7579f131`
	urlNode := "www.bilibili.com/read/mobile/" + strconv.FormatInt(id, 10)
	resp, err := http.Post(url, "text/plain", strings.NewReader(urlNode))
	if err != nil {
		log.Error("s.pushShenma post error(%+v)", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("s.pushShenma ReadAll error(%+v)", err)
		return
	}
	var data struct {
		ReturnCode int    `json:"returnCode"`
		ErrMsg     string `json:"errMsg"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		log.Error("s.pushShenma json unmarshal error(%+v)", err)
		return
	}
	if data.ReturnCode == 200 {
		log.Info("s.pushShenma success, url(%s)", urlNode)
	} else {
		log.Error("s.pushShenma success, errMsg(%s), url(%s)", data.ErrMsg, urlNode)
	}
	return
}
