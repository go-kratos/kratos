package service

import (
	"bufio"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/job/main/push/dao"
	pb "go-common/app/service/main/push/api/grpc/v1"
	pushmdl "go-common/app/service/main/push/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_delTaskLimit = 5000
)

var (
	errEmptyLine    = errors.New("empty line")
	errInvalidMid   = errors.New("invalid mid format")
	errInvalidToken = errors.New("invalid token format")
)

func (s *Service) addTaskproc() {
	defer s.addTaskWg.Done()
	var err error
	for {
		task, ok := <-s.addTaskCh
		if !ok {
			log.Info("add task channel exit")
			return
		}
		if task == nil {
			continue
		}
		task.Status = pushmdl.TaskStatusPrepared
		for i := 0; i < _retry; i++ {
			if err = s.dao.AddTask(context.Background(), task); err == nil {
				break
			}
		}
		if err != nil {
			log.Error("add task(%+v) error(%v)", task, err)
			s.cache.Save(func() {
				s.dao.SendWechat(fmt.Sprintf("add task(%d)", task.Job))
			})
			continue
		}
		dao.PromInfo("add task")
		time.Sleep(time.Millisecond)
	}
}

func (s *Service) delTasksproc() {
	for {
		now := time.Now()
		// 每天2点时删除一个月前的task数据
		if now.Hour() != 2 {
			time.Sleep(time.Minute)
			continue
		}
		var (
			err     error
			deleted int64
			b       = now.Add(time.Duration(-s.c.Job.DelTaskInterval*24) * time.Hour)
			loc, _  = time.LoadLocation("Local")
			t       = time.Date(b.Year(), b.Month(), b.Day(), 23, 59, 59, 0, loc)
		)
		for {
			if deleted, err = s.dao.DelTasks(context.TODO(), t, _delTaskLimit); err != nil {
				log.Error("s.delTasks(%v) error(%v)", t, err)
				s.dao.SendWechat("DB操作失败:push-job删除task数据错误")
				time.Sleep(time.Second)
				continue
			}
			if deleted < _delTaskLimit {
				break
			}
			time.Sleep(time.Second)
		}
		time.Sleep(time.Hour)
	}
}

func (s *Service) pretreatTaskproc() {
	defer s.waiter.Done()
	for {
		if s.closed {
			return
		}
		task, err := s.pickPretreatmentTask()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if task != nil {
			log.Info("pretreat task job(%d) id(%s)", task.Job, task.ID)
			if err = s.pretreatTask(task); err != nil {
				log.Error("pretreat task(%+v) error(%v)", task, err)
				s.cache.Save(func() { s.dao.SendWechat(fmt.Sprintf("pretreat task(%s) error", task.ID)) })
			}
		}
		time.Sleep(time.Duration(s.c.Job.LoadTaskInteval))
	}
}

func (s *Service) pickPretreatmentTask() (t *pushmdl.Task, err error) {
	c := context.Background()
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTx(c); err != nil {
		log.Error("tx.BeginTx() error(%v)", err)
		return
	}
	if t, err = s.dao.TxTaskByStatus(tx, pushmdl.TaskStatusPretreatmentPrepared); err != nil {
		if e := tx.Rollback(); e != nil {
			dao.PromError("task:获取新任务")
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if t == nil {
		if e := tx.Rollback(); e != nil {
			dao.PromError("task:获取新任务")
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if err = s.dao.TxUpdateTaskStatus(tx, t.ID, pushmdl.TaskStatusPretreatmentDoing); err != nil {
		if e := tx.Rollback(); e != nil {
			dao.PromError("task:更新任务状态")
			log.Error("tx.Rollback() error(%v)", e)
		}
		return
	}
	if err = tx.Commit(); err != nil {
		dao.PromError("task:获取新任务commit")
		log.Error("tx.Commit() error(%v)", err)
	}
	return
}

func (s *Service) pretreatTask(t *pushmdl.Task) (err error) {
	id, _ := strconv.ParseInt(t.ID, 10, 64)
	switch t.Type {
	case pushmdl.TaskTypeAll:
		err = s.pretreatTaskAll(t)
	case pushmdl.TaskTypeMngToken, pushmdl.TaskTypeDataPlatformToken, pushmdl.TaskTypeDataPlatformMid:
		err = s.pretreatTaskToken(t)
	case pushmdl.TaskTypeStrategyMid, pushmdl.TaskTypeMngMid:
		err = s.pretreatTaskMid(t)
	default:
		log.Error("invalid task type, (%+v)", t)
	}
	if err != nil {
		err = s.dao.UpdateTaskStatus(context.Background(), id, pushmdl.TaskStatusPretreatmentFailed)
		return
	}
	err = s.dao.UpdateTaskStatus(context.Background(), id, pushmdl.TaskStatusPretreatmentDone)
	return
}

func (s *Service) pretreatTaskAll(t *pushmdl.Task) (err error) {
	log.Info("AddTaskAll start, task(%+v)", t)
	var (
		maxID int64
		group = errgroup.Group{}
	)
	maxID, err = s.dao.ReportLastID(context.Background())
	if err != nil || maxID <= 0 {
		log.Error("s.pretreatTaskAll() error(%v)", err)
		s.cache.Save(func() {
			s.dao.SendWechat(fmt.Sprintf("pretreatTaskAll(%v) ReportLastID(%d) error", t.ID, maxID))
		})
		return
	}
	log.Info("AddTaskAll get last report ID(%d)", maxID)
	buildCount := len(t.Build)
	batch := maxID / int64(s.c.Job.TaskGoroutines)
	for j := 0; j < s.c.Job.TaskGoroutines; j++ {
		begin := int64(j) * batch
		end := begin + batch
		group.Go(func() (e error) {
			var (
				path   string
				rows   *xsql.Rows
				tokens = make(map[int][]string)
			)
			for {
				if begin >= end {
					break
				}
				l := begin + int64(_dbBatch)
				if l >= end {
					l = end
				}
				log.Info("AddTaskAll load reports start(%d) end(%d)", begin, l)
				if rows, e = s.dao.ReportsTaskAll(context.Background(), begin, l, t.APPID); e != nil {
					log.Error("s.dao.ReportsTaskAll(%d,%d,%d) error(%v)", begin, l, t.APPID)
					s.cache.Save(func() {
						s.dao.SendWechat(fmt.Sprintf("pretreatTaskAll(%v) ReportsTaskAll(%d,%d,%d) error", t.ID, begin, l, t.APPID))
					})
					return
				}
				for rows.Next() {
					var (
						platformID int
						build      int
						token      string
					)
					if e = rows.Scan(&platformID, &token, &build); e != nil {
						log.Error("AddTaskAll rows.Scan() error(%v)", e)
						s.cache.Save(func() {
							s.dao.SendWechat(fmt.Sprintf("pretreatTaskAll(%v) ReportsTaskAll(%d,%d,%d) error", t.ID, begin, l, t.APPID))
						})
						return
					}
					if buildCount > 0 && !pushmdl.ValidateBuild(platformID, build, t.Build) {
						continue
					}
					tokens[platformID] = append(tokens[platformID], token)
					if len(tokens[platformID]) >= s.c.Job.LimitPerTask {
						if path, e = s.saveFile(tokens[platformID]); e != nil {
							log.Error("AddTaskAll s.saveTokens error(%v)", e)
							s.cache.Save(func() {
								s.dao.SendWechat(fmt.Sprintf("pretreatTaskAll(%v) saveTokens error(%v)", t.ID, e))
							})
							return
						}
						tokens[platformID] = []string{}
						task := *t
						task.MidFile = path
						task.PlatformID = platformID
						s.addTaskCh <- &task
					}
				}
				begin = l
			}
			for p, v := range tokens {
				if len(v) == 0 {
					continue
				}
				if path, e = s.saveFile(v); e == nil {
					task := *t
					task.MidFile = path
					task.PlatformID = p
					s.addTaskCh <- &task
				}
			}
			return
		})
	}
	if err = group.Wait(); err != nil {
		log.Error("add task all, task(%+v) error(%v)", t, err)
		s.cache.Save(func() {
			s.dao.SendWechat(fmt.Sprintf("pretreatTaskAll(%v) error(%v)", t.ID, err))
		})
		return
	}
	log.Info("AddTaskAll end, task(%+v)", t)
	s.cache.Save(func() {
		s.dao.SendWechat(fmt.Sprintf("add task all success, job(%d)", t.Job))
	})
	return
}

func (s *Service) pretreatTaskMid(t *pushmdl.Task) (err error) {
	f, err := os.Open(t.MidFile)
	if err != nil {
		log.Error("pretreatTaskMid(%+v) open file error(%v)", t, err)
		return
	}
	defer f.Close()
	var (
		exit     bool
		line     string
		path     string
		mid      int64
		counter  int
		midTotal int64
		midValid int64
		mu       sync.Mutex
		mids     []int64
		tokens   = make(map[int][]string)
		group    = errgroup.Group{}
		reader   = bufio.NewReader(f)
	)
	for {
		if exit {
			break
		}
		if line, err = reader.ReadString('\n'); err != nil {
			if err == io.EOF {
				exit = true
			} else {
				log.Error("read file error(%v)", err)
				continue
			}
		}
		if mid, err = parseMidLine(line); err != nil {
			log.Error("parse mid line(%s) error(%v)", line, err)
			continue
		}
		midTotal++
		mids = append(mids, mid)
		if len(mids) >= s.c.Job.PushPartSize {
			midsCp := make([]int64, len(mids))
			copy(midsCp, mids)
			mids = []int64{}
			group.Go(func() (e error) {
				ts, valid, e := s.tokensByMids(t, midsCp)
				if e != nil {
					log.Error("s.tokensByMids(%v) error(%v)", t.ID, e)
					return
				}
				tcopy := make(map[int][]string)
				mu.Lock()
				midValid += valid
				for p, v := range ts {
					tokens[p] = append(tokens[p], v...)
					if len(tokens[p]) >= s.c.Job.LimitPerTask {
						tcopy[p] = append(tcopy[p], tokens[p]...)
						tokens[p] = []string{}
					}
				}
				mu.Unlock()
				for p, v := range tcopy {
					if path, err = s.saveFile(v); err != nil {
						log.Error("pretreatTaskMid s.saveFild error(%v)", err)
						s.cache.Save(func() {
							s.dao.SendWechat(fmt.Sprintf("pretreatTaskMid(%v) saveTokens error(%v)", t.ID, err))
						})
						return
					}
					task := *t
					task.MidFile = path
					task.PlatformID = p
					s.addTaskCh <- &task
				}
				return
			})
			counter++
			if counter == s.c.Job.PushPartChanSize {
				group.Wait()
				counter = 0
			}
		}
	}
	if counter > 0 {
		group.Wait()
	}
	if len(mids) > 0 {
		var (
			valid int64
			ts    map[int][]string
		)
		if ts, valid, err = s.tokensByMids(t, mids); err == nil {
			midValid += valid
			for p, v := range ts {
				tokens[p] = append(tokens[p], v...)
			}
		} else {
			log.Error("s.tokensByMids(%+v) error(%v)", t, err)
		}
	}
	s.cache.Save(func() {
		arg := &pb.AddMidProgressRequest{Task: t.ID, MidTotal: midTotal, MidValid: midValid}
		if _, e := s.pushRPC.AddMidProgress(context.Background(), arg); e != nil {
			log.Error("s.pushRPC.AddMidProgress(%+v) error(%v)", arg, e)
		}
	})
	for p, v := range tokens {
		if len(v) == 0 {
			continue
		}
		if path, err = s.saveFile(v); err != nil {
			log.Error("pretreatTaskMid s.saveFild error(%v)", err)
			return
		}
		task := *t
		task.MidFile = path
		task.PlatformID = p
		s.addTaskCh <- &task
	}
	log.Info("pretreatTaskMid task(%+v)", t)
	return
}

func (s *Service) pretreatTaskToken(t *pushmdl.Task) (err error) {
	f, err := os.Open(t.MidFile)
	if err != nil {
		log.Error("pretreatTaskToken(%+v) open file error(%v)", t, err)
		return
	}
	defer f.Close()
	var (
		exit   bool
		plat   int
		line   string
		token  string
		path   string
		tokens = make(map[int][]string)
		reader = bufio.NewReader(f)
	)
	for {
		if exit {
			break
		}
		if line, err = reader.ReadString('\n'); err != nil {
			if err == io.EOF {
				exit = true // no 'continue', solve the last line whitout '\n'
			} else {
				log.Error("read file error(%v)", err)
				continue
			}
		}
		if plat, token, err = parseTokenLine(line); err != nil {
			log.Error("parse token line(%s) error(%v)", line, err)
			continue
		}
		tokens[plat] = append(tokens[plat], token)
		if len(tokens[plat]) >= s.c.Job.LimitPerTask {
			if path, err = s.saveFile(tokens[plat]); err != nil {
				log.Error("pretreatTaskToken s.saveFile error(%v)", err)
				s.cache.Save(func() {
					s.dao.SendWechat(fmt.Sprintf("pretreatTaskToken(%v) saveTokens error(%v)", t.ID, err))
				})
				return
			}
			tokens[plat] = []string{}
			task := *t
			task.MidFile = path
			task.PlatformID = plat
			s.addTaskCh <- &task
		}
	}
	for p, v := range tokens {
		if len(v) == 0 {
			continue
		}
		if path, err = s.saveFile(v); err == nil {
			task := *t
			task.MidFile = path
			task.PlatformID = p
			s.addTaskCh <- &task
		}
	}
	log.Info("pretreatTaskToken task(%+v)", t)
	return
}

func parseTokenLine(line string) (plat int, token string, err error) {
	line = strings.Trim(line, " \r\n")
	if line == "" {
		err = errEmptyLine
		return
	}
	res := strings.Split(line, "\t")
	if len(res) != 2 {
		err = errInvalidToken
		return
	}
	if res[0] == "" || res[1] == "" {
		err = errInvalidToken
		return
	}
	if plat, err = strconv.Atoi(res[0]); err != nil || plat <= 0 {
		err = errInvalidToken
		return
	}
	token = res[1]
	return
}

func parseMidLine(line string) (mid int64, err error) {
	line = strings.Trim(line, " \r\t\n")
	if line == "" {
		err = errEmptyLine
		return
	}
	if mid, err = strconv.ParseInt(line, 10, 64); err != nil || mid <= 0 {
		err = errInvalidMid
	}
	return
}

func (s *Service) saveFile(tokens []string) (path string, err error) {
	name := strconv.FormatInt(time.Now().UnixNano(), 10) + tokens[0]
	data := []byte(strings.Join(tokens, "\n"))
	for i := 0; i < _retry; i++ {
		if path, err = s.saveNASFile(name, data); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

// saveNASFile writes data into NAS.
func (s *Service) saveNASFile(name string, data []byte) (path string, err error) {
	name = fmt.Sprintf("%x", md5.Sum([]byte(name)))
	dir := fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(s.c.Job.MountDir, "/"), time.Now().Format("20060102"), name[:2])
	if _, err = os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			log.Error("os.IsNotExist(%s) error(%v)", dir, err)
			return
		}
		if err = os.MkdirAll(dir, 0777); err != nil {
			log.Error("os.MkdirAll(%s) error(%v)", dir, err)
			return
		}
	}
	path = fmt.Sprintf("%s/%s", dir, name)
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("s.saveNASFile(%s) OpenFile() error(%v)", path, err)
		return
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		log.Error("s.saveNASFile(%s) f.Write() error(%v)", err)
	}
	return
}
