package report

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	xhttp "net/http"
	"strings"
	"sync"
	"time"

	"go-common/app/job/main/tv/conf"
	"go-common/app/job/main/tv/dao/report"
	mdlrep "go-common/app/job/main/tv/model/report"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"

	"github.com/robfig/cron"
)

const (
	_retry      = 3
	_jobRunning = 3
	_startJob   = 4
	_readSize   = 1024
)

// Service struct of service .
type Service struct {
	c        *conf.Config
	ch       chan bool
	dao      *report.Dao
	respURL  map[string]interface{}
	cache    *fanout.Fanout
	lock     sync.Mutex
	labelRes map[int]map[string]int
	readSize int
	// cron
	cron *cron.Cron
}

// New creates a Service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      report.New(c),
		respURL:  make(map[string]interface{}),
		cache:    fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
		labelRes: make(map[int]map[string]int),
		readSize: c.Report.ReadSize * _readSize,
		cron:     cron.New(),
		ch:       make(chan bool, c.Report.RoutineCount),
	}
	if err := s.cron.AddFunc(s.c.Report.CronAc, s.oneWork(mdlrep.ArchiveClick)); err != nil { // corn report run
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Report.CronAd, s.oneWork(mdlrep.ActiveDuration)); err != nil { // corn report run
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Report.CronPd, s.oneWork(mdlrep.PlayDuration)); err != nil { // corn report run
		panic(err)
	}
	if err := s.cron.AddFunc(s.c.Report.CronVe, s.oneWork(mdlrep.VisitEvent)); err != nil { // corn report run
		panic(err)
	}
	s.cron.Start()
	s.readCache()      // data report
	s.readLabelCache() // label style
	go s.reportCon()   // data report
	go s.showStyle()   // label style
	go s.showLabel()   // label style
	return
}

func (s *Service) oneWork(table string) func() {
	return func() {
		if s.c.Report.Env != "prod" {
			return
		}
		var (
			res string
			err error
		)
		if res, err = s.requestURL(table); err != nil {
			log.Error("reportPro s.requestURL() error(%v)", err)
			return
		}
		if res == "" {
			return
		}
		s.lock.Lock()
		s.respURL[res] = struct{}{}
		s.lock.Unlock()
		s.setCache()
	}
}

func (s *Service) reportCon() {
	if s.c.Report.Env != "prod" {
		return
	}
	var (
		err  error
		info *mdlrep.DpCheckJobResult
	)
	for {
		var (
			flags, failStr []string
		)
		s.lock.Lock()
		for k := range s.respURL {
			flags = append(flags, k)
		}
		s.respURL = make(map[string]interface{})
		s.lock.Unlock()
		for _, v := range flags {
			if v == "" {
				continue
			}
			// loop send http request and return result
			if info, err = s.check(v); err == nil && len(info.Files) > 0 {
				now := time.Now()
				s.upReport(info)
				log.Warn("report success fileNum(%d) url(%s) 本次上报数据耗时: %s", len(info.Files), v, time.Since(now))
				continue
			}
			if info.StatusID == _jobRunning || info.StatusID == _startJob {
				failStr = append(failStr, v)
			}
		}
		s.lock.Lock()
		for _, v := range failStr {
			s.respURL[v] = struct{}{}
		}
		s.lock.Unlock()
		s.setCache()
		time.Sleep(3 * time.Second)
	}
}

func (s *Service) readFile(path string) {
	var (
		n       int
		err     error
		resdata []map[string]interface{}
		resp    *xhttp.Response
		buf     = make([]byte, 1024)
		chunks  []byte
		req     *xhttp.Request
		fileCnt = 0
	)
	client := &xhttp.Client{
		Transport: &xhttp.Transport{
			DisableKeepAlives: true,
		},
	}
	req, err = xhttp.NewRequest("GET", path, strings.NewReader(""))
	if err != nil {
		log.Error("[url(%s)] xhttp.NewRequest error(%v)", path, err)
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Error("[url(%s)] client.Do error(%v)", path, err)
		return
	}
	defer resp.Body.Close()
	for {
		n, err = resp.Body.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("resp.Body.Read error(%v)", err)
			return
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
		if len(chunks) > s.readSize { // 500K
			lastPos := bytes.LastIndex(chunks, []byte("\n"))
			if lastPos < 0 {
				continue
			}
			fileCnt = fileCnt + 1
			results := append([]byte{}, chunks[:lastPos]...)
			chunks = append([]byte{}, chunks[lastPos:]...)
			bsdata := bytes.Split(results, []byte("\n"))
			for _, bs := range bsdata {
				n := bytes.Split(bs, []byte("\u0001"))
				m := mdlrep.ArcClickParam(n)
				resdata = append(resdata, m)
			}
			if err = s.postData(resdata); err != nil {
				log.Error("[url(%s)] s.postData error(%v)", path, err)
			}
			resdata = make([]map[string]interface{}, 0)
		}
	}
	if len(chunks) > 0 {
		bsdata := bytes.Split(chunks, []byte("\n"))
		for _, bs := range bsdata {
			n := bytes.Split(bs, []byte("\u0001"))
			m := mdlrep.ArcClickParam(n)
			resdata = append(resdata, m)
		}
		if err = s.postData(resdata); err != nil {
			log.Error("[url(%s)] s.postData error(%v)", path, err)
		}
	}
}

func (s *Service) requestURL(table string) (res string, err error) {
	for i := 0; i < _retry; i++ {
		if res, err = s.dao.Report(context.Background(), table); err == nil {
			break
		}
	}
	return
}

func (s *Service) check(res string) (info *mdlrep.DpCheckJobResult, err error) {
	for i := 0; i < _retry; i++ {
		if info, err = s.dao.CheckJob(context.Background(), res); err == nil {
			break
		}
	}
	return
}

// upReport .
func (s *Service) upReport(info *mdlrep.DpCheckJobResult) {
	for _, v := range info.Files {
		s.readFile(v)
	}
}

func (s *Service) postData(param []map[string]interface{}) (err error) {
	for _, v := range param {
		s.ch <- true
		go s.sendOnce(v)
	}
	return
}

func (s *Service) sendOnce(v map[string]interface{}) (err error) {
	var (
		body string
		data []byte
	)
	defer func() {
		<-s.ch
	}()
	if data, err = json.Marshal(v); err != nil {
		log.Error("Service postData json.Marshal error(%v)", err)
		return
	}
	body = body + string(data) + ","
	s.dealBody(body)
	return
}

func (s *Service) readCache() {
	if s.c.Report.Env != "prod" {
		return
	}
	var (
		err   error
		btRes = make(map[string]interface{})
	)
	if btRes, err = s.dao.GetReportCache(context.Background()); err != nil {
		log.Error("s.dao.GetReportCache error(%v)", err)
		panic(err)
	}
	s.respURL = btRes
}

func (s *Service) dealBody(body string) {
	body = strings.Replace(body, `\\N`, "", -1)
	body = strings.TrimSuffix(body, ",")
	body = `{"code": 0,"message": "0","ttl": 1,"data":[` + body + `]}`
	if err := s.dao.PostRequest(context.Background(), body); err != nil {
		log.Error("s.dao.PostRequest error(%v)", err)
	}
}

func (s *Service) setCache() {
	s.cache.Do(context.Background(), func(c context.Context) {
		s.dao.SetReportCache(c, s.respURL)
	})
}
