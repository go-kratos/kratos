package service

import (
	"go-common/app/service/main/upcredit/model/upcrmmodel"
	"time"

	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/siddontang/go-mysql/mysql"
	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/dao/upcrmdao"
	"go-common/app/service/main/upcredit/model/calculator"
	"go-common/library/log"
	"sync"
)

//CalcService calc service
type CalcService struct {
	c             *conf.Config
	JobChannel    chan BaseJobInterface
	OutputChan    chan *upcrmmodel.UpScoreHistory
	calc          *calculator.CreditScoreCalculator
	wg            sync.WaitGroup
	isRunning     bool
	mapLock       sync.Mutex
	jobsMap       map[int]*CalcJob
	jobID         int
	finishJobChan chan *CalcJob
	crmdao        *upcrmdao.Dao
}

//NewCalc create calc service
func NewCalc(c *conf.Config, outChan chan *upcrmmodel.UpScoreHistory, crmdao *upcrmdao.Dao) *CalcService {
	var s = &CalcService{
		c:             c,
		JobChannel:    make(chan BaseJobInterface, 100),
		calc:          calculator.New(outChan),
		isRunning:     true,
		jobsMap:       make(map[int]*CalcJob),
		jobID:         1,
		finishJobChan: make(chan *CalcJob),
		crmdao:        crmdao,
	}
	return s
}

//BaseJobInterface job interface
type BaseJobInterface interface {
	Run() (err error)
	Description() string
}

//CalcJob calculate job
type CalcJob struct {
	ID int
	// 数据日
	Date time.Time
	// 对应的数据
	TableNum int
	Overall  *calculator.OverAllStatistic
	IsDone   bool
	Svc      *CalcService
}

//Run run job
func (job *CalcJob) Run() (err error) {
	err = job.Svc.calc.CalcLogTable(job.TableNum, job.Date, job.Overall)
	if err != nil {
		log.Error("calculate job error, job=%+v, err=%+v", job, err)
	} else {
		log.Info("calculate job finish, job=%+v", job)
	}
	job.IsDone = true
	job.Svc.finishJobChan <- job
	return
}

//Description descrpit job
func (job *CalcJob) Description() (dest string) {
	return fmt.Sprintf("calc table=%d, date=%s", job.TableNum, job.Date.Format(mysql.TimeFormat))
}

//CalcStatisticJob statis job
type CalcStatisticJob struct {
	ID int
	// 数据日
	Date time.Time
	Svc  *CalcService
}

//Run run
func (j *CalcStatisticJob) Run() (err error) {
	return j.Svc.CalcScoreSectionData()
}

//Description desc
func (j *CalcStatisticJob) Description() (dest string) {
	return fmt.Sprintf("calc statis, date=%s", j.Date.Format(mysql.TimeFormat))
}

//AddCalcJob add job
func (c *CalcService) AddCalcJob(date time.Time) {
	var count = upcrmmodel.CreditLogTableCount
	for i := 0; i < count; i++ {
		var job = &CalcJob{
			ID:       c.getJobID(),
			Date:     date,
			TableNum: i,
			Overall:  calculator.NewOverAllStatistic(1000, 10),
		}
		job.Svc = c
		log.Info("send calculate job, job=%+v", job)
		c.JobChannel <- job
		c.mapLock.Lock()
		c.jobsMap[job.ID] = job
		c.mapLock.Unlock()
	}
}

func (c *CalcService) getJobID() int {
	c.jobID++
	return c.jobID
}

//Run fun service
func (c *CalcService) Run() {
	go c.ScheduleJob()
	go c.jobFinishCheck()

	// run workers
	for i := 0; i < c.c.RunStatJobConf.WorkerNumber; i++ {
		c.wg.Add(1)
		go c.worker()
	}
	log.Info("run worker count=%d", c.c.RunStatJobConf.WorkerNumber)
}

//Close close service
func (c *CalcService) Close() {
	c.isRunning = false
	//c.wg.Wait()
}

//ScheduleJob schedule job
func (c *CalcService) ScheduleJob() {
	var now = time.Now()
	var startTime, err = time.Parse("15:04:05", c.c.RunStatJobConf.StartTime)
	if err != nil {
		log.Error("schedule job start time error, config=%s, should like '12:00:00'")
		startTime = time.Date(2000, 1, 1, 3, 0, 0, 0, now.Location())
	}
	n := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), 0, now.Location())
	d := n.Sub(now)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(now)
	}
	for c.isRunning {
		time.Sleep(d)
		d = 24 * time.Hour
		// 只有master来计算任务，做的比较简单，正常应该是master发任务，follower来做任务，然而简单起见，先由master自己来完成任务
		if !conf.IsMaster {
			continue
		}
		c.AddCalcJob(time.Now())
	}
}

func (c *CalcService) worker() {
	defer func() {
		c.wg.Done()
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("worker Runtime error caught, try recover: %+v", r)
			c.wg.Add(1)
			go c.worker()
		}
	}()
	for job := range c.JobChannel {
		var err = job.Run()
		if err != nil {
			log.Error("job run, err=%+v", err)
		}
	}
}

func (c *CalcService) jobFinishCheck() {
	for range c.finishJobChan {
		if c.isAllJobFinish() {
			c.onAllJobFinish()
		}
	}
}

func (c *CalcService) onAllJobFinish() {
	log.Info("all job is done, add calc statistic job")

	var job = &CalcStatisticJob{
		ID:   c.getJobID(),
		Date: time.Now().AddDate(0, 0, -1),
		Svc:  c,
	}
	c.JobChannel <- job
}
func (c *CalcService) isAllJobFinish() bool {
	c.mapLock.Lock()
	defer c.mapLock.Unlock()

	if len(c.jobsMap) == 0 {
		return true
	}
	for _, j := range c.jobsMap {
		if !j.IsDone {
			return false
		}
	}
	return true
}

//CalcScoreSectionData calc score section
func (c *CalcService) CalcScoreSectionData() (err error) {
	var crmdb, e = gorm.Open("mysql", conf.Conf.DB.Upcrm.DSN)
	err = e
	if e != nil {
		log.Error("fail to open crm db, for score section")
		return
	}
	log.Info("start calculate crm score section")
	defer crmdb.Close()
	var limit = 1000
	var total = 0
	var prSection = calculator.NewOverAllStatistic(1000, 10)
	var qualitySection = calculator.NewOverAllStatistic(1000, 10)
	var creditSection = calculator.NewOverAllStatistic(1000, 10)
	var lastID = int32(0)

	for {
		var upInfos []upcrmmodel.UpBaseInfo
		// 1. 从数据库中取limit条数据
		err = crmdb.Table("up_base_info").Select("id, credit_score, pr_score, quality_score").Where("id>?", lastID).Limit(limit).Find(&upInfos).Error
		if err != nil {
			log.Error("get data err, e=%+v", e)
			break
		}

		// 2. 计算各分数的分数段
		for _, u := range upInfos {
			prSection.AddScore(u.PrScore, -1)
			qualitySection.AddScore(u.QualityScore, -1)
			creditSection.AddScore(u.CreditScore, 500)
			lastID = u.ID
		}

		var thisCount = len(upInfos)
		total += thisCount
		if thisCount < limit {
			log.Info("table[up_base_info] total read record, num=%d", total)
			break
		}
	}
	// 输出到score表中
	var now = time.Now()
	err = c.crmdao.InsertScoreSection(*prSection, upcrmmodel.ScoreTypePr, now)
	if err != nil {
		log.Error("fail update pr score section, err=%+v", err)
	}
	err = c.crmdao.InsertScoreSection(*qualitySection, upcrmmodel.ScoreTypeQuality, now)
	if err != nil {
		log.Error("fail update quality score section, err=%+v", err)
	}
	err = c.crmdao.InsertScoreSection(*creditSection, upcrmmodel.ScoreTypeCredit, now)
	if err != nil {
		log.Error("fail update credit score section, err=%+v", err)
	}
	log.Info("finish calculate crm score section, totalcount=%d", total)

	return
}
