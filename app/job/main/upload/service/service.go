package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"go-common/app/job/main/upload/conf"
	"go-common/app/job/main/upload/dao"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
)

const (
	_downloadFmt       = "%s/bfs/%s/%s"
	_uploadAdminAddFmt = "%s/add"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// Meta describe databus message value format from bfs proxy
type Meta struct {
	Bucket   string `json:"bucket"`
	Filename string `json:"filename"`
	Mine     string `json:"mine"`
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	Run(c)
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

// Job .
type Job struct {
	downloadHost             string
	uploadHost               string
	uploadAdminHost          string
	xclient                  *xhttp.Client
	sub                      *databus.Databus
	ch                       chan *Meta
	routineCount             int
	AIYellowingConsumer      *databus.Databus //consume response from ai
	AIYellowingProducer      *databus.Databus //produce request to ai
	AIYellowingExceptBuckets map[string]bool
	Threshold                *conf.Threshold //score threshold
}

// AIReqMessage defined send and receive databus message format with AI
type AIReqMessage struct {
	Bucket   string `json:"bucket"`
	FileName string `json:"file_name"`
	URL      string `json:"url"`
	IsYellow bool   `json:"is_yellow"`
}

//AIRespMessage defined response databus message format from AI
type AIRespMessage struct {
	URL       string  `json:"url"`
	FileName  string  `json:"file_name"`
	Bucket    string  `json:"bucket"`
	Sex       float64 `json:"sex"`
	Violent   float64 `json:"violent"`
	Blood     float64 `json:"blood"`
	Politics  float64 `json:"politics"`
	IsYellow  bool    `json:"is_yellow"`
	ErrorCode int64   `json:"error_code"`
	ErrorMsg  string  `json:"error_msg"`
}

// Add describe params of /x/admin/bfs-upload/add
type Add struct {
	Bucket   string `json:"bucket"`
	FileName string `json:"filename"`
}

// AddResp describe response of /x/admin/bfs-upload/add
type AddResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

// Run .
func Run(cfg *conf.Config) {
	exceptBuckets := make(map[string]bool)
	for _, bucket := range cfg.AIYellowing.ExceptBuckets {
		exceptBuckets[bucket] = false
	}
	if cfg.Threshold == nil {
		cfg.Threshold = &conf.Threshold{
			Sex:      5000,
			Politics: 5000,
			Blood:    5000,
			Violent:  5000,
		}
	}
	if cfg.Threshold.Sex == 0 {
		cfg.Threshold.Sex = 5000
	}
	if cfg.Threshold.Politics == 0 {
		cfg.Threshold.Politics = 5000
	}
	if cfg.Threshold.Blood == 0 {
		cfg.Threshold.Blood = 5000
	}
	if cfg.Threshold.Violent == 0 {
		cfg.Threshold.Violent = 5000
	}

	j := &Job{
		downloadHost:             cfg.DonwloadHost,
		uploadHost:               cfg.UploadHost,
		uploadAdminHost:          cfg.UploadAdminHost,
		routineCount:             cfg.RoutineCount,
		xclient:                  bm.NewClient(cfg.HTTPClient),
		sub:                      databus.New(cfg.Databus),
		ch:                       make(chan *Meta, 1024),
		AIYellowingConsumer:      databus.New(cfg.AIYellowing.Consumer),
		AIYellowingProducer:      databus.New(cfg.AIYellowing.Producer),
		AIYellowingExceptBuckets: exceptBuckets,
		Threshold:                cfg.Threshold,
	}
	go j.aireqproc()
	go j.proxysyncproc(context.Background())
	go j.aiResp(context.Background()) // deal ai response
}

func (j *Job) aireqproc() {
	for i := 0; i < 10; i++ {
		go func() {
			for meta := range j.ch {
				j.aiReq(context.Background(), meta)
			}
		}()
	}
}

// Run .
func (j *Job) proxysyncproc(ctx context.Context) {
	var (
		err  error
		meta *Meta
	)
	for {
		msg, ok := <-j.sub.Messages()
		if !ok {
			log.Error("consume msg error, uploadproc exit!")
			return
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit() error(%v)", err)
			continue
		}
		meta = new(Meta)
		if err = json.Unmarshal(msg.Value, meta); err != nil {
			log.Error("Job.run.Unmarshal(key:%v),err(%v)", msg.Key, err)
			continue
		}
		// if bucket need yellow inspects
		_, ok = j.AIYellowingExceptBuckets[meta.Bucket]
		if !ok {
			log.Info("bfs-proxy sync a msg:(%+v) meta(%+v) need to req ai", msg, meta)
			j.add(meta)
		}
	}
}

func (j *Job) add(item *Meta) {
	j.ch <- item
}

// addRecord add a record into upload-admin if a pic is yellow
func (j *Job) addRecord(aiMsg *AIRespMessage) (err error) {
	var (
		uploadAdminAddURL string
	)

	addResp := new(AddResp)
	params := url.Values{}
	params.Add("bucket", aiMsg.Bucket)
	params.Add("filename", aiMsg.FileName)
	params.Add("url", aiMsg.URL)
	params.Add("sex", strconv.Itoa(int(math.Round(aiMsg.Sex*10000))))
	params.Add("politics", strconv.Itoa(int(math.Round(aiMsg.Politics*10000))))

	uploadAdminAddURL = fmt.Sprintf(_uploadAdminAddFmt, j.uploadAdminHost)

	if err = j.xclient.Post(context.TODO(), uploadAdminAddURL, "", params, addResp); err != nil {
		return
	}
	if addResp.Code != 0 {
		log.Error("call /x/admin/bfs-upload/add code error, code(%d),message(%s)", addResp.Code, addResp.Message)
		return
	}
	log.Info("upload-admin add success(%+v)", aiMsg)
	return
}

// RetryAddRecord try to add a record to upload-admin
func (j *Job) RetryAddRecord(aiMsg *AIRespMessage) (err error) {
	attempts := 3
	for i := 0; i < attempts; i++ {
		err = j.addRecord(aiMsg)
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
	return
}
