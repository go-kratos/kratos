package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go-common/app/admin/main/dm/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_csvTaskTitle = []string{"dmid", "oid", "mid", "state", "msg", "ip", "ctime"}
)

// TaskList .
func (s *Service) TaskList(c context.Context, v *model.TaskListArg) (res *model.TaskList, err error) {
	taskSQL := make([]string, 0)
	if v.Creator != "" {
		taskSQL = append(taskSQL, fmt.Sprintf("creator LIKE %q", "%"+v.Creator+"%"))
	}
	if v.Reviewer != "" {
		taskSQL = append(taskSQL, fmt.Sprintf("reviewer LIKE %q", "%"+v.Reviewer+"%"))
	}
	if v.State >= 0 {
		taskSQL = append(taskSQL, fmt.Sprintf("state=%d", v.State))
	}
	if v.Title != "" {
		taskSQL = append(taskSQL, fmt.Sprintf("title LIKE %q", "%"+v.Title+"%"))
	}
	if v.Ctime != "" {
		taskSQL = append(taskSQL, fmt.Sprintf("ctime>=%q", v.Ctime))
	}
	tasks, total, err := s.dao.TaskList(c, taskSQL, v.Pn, v.Ps)
	if err != nil {
		return
	}
	res = &model.TaskList{
		Result: tasks,
		Page: &model.PageInfo{
			Num:   v.Pn,
			Size:  v.Ps,
			Total: total,
		},
	}
	return
}

// AddTask .
func (s *Service) AddTask(c context.Context, v *model.AddTaskArg) (err error) {
	var taskID int64
	var sub int32
	tx, err := s.dao.BeginBiliDMTrans(c)
	if err != nil {
		log.Error("tx.BeginBiliDMTrans error(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if v.Regex != "" {
		if len([]rune(v.Regex)) > model.TaskRegexLen {
			err = ecode.DMTaskRegexTooLong
			return
		}
		if _, err = regexp.Compile(v.Regex); err != nil {
			log.Error("regexp.Compile(v.Regex:%s) error(%v)", v.Regex, err)
			err = ecode.DMTaskRegexIllegal
			return
		}
	}
	if v.Operation >= 0 {
		sub = 1
	}
	if taskID, err = s.dao.AddTask(tx, v, sub); err != nil {
		return
	}
	if sub > 0 {
		_, err = s.dao.AddSubTask(tx, taskID, v.Operation, v.OpTime, v.OpRate)
	}
	return
}

// ReviewTask .
func (s *Service) ReviewTask(c context.Context, v *model.ReviewTaskArg) (err error) {
	if v.State == model.TaskReviewPass {
		var (
			task         *model.TaskView
			sTime, eTime time.Time
		)
		if task, err = s.dao.TaskView(c, v.ID); err != nil {
			return
		}
		taskSQL := make([]string, 0)
		if task.Regex != "" {
			taskSQL = append(taskSQL, fmt.Sprintf("content.msg regexp %q", task.Regex))
		}
		if task.KeyWords != "" {
			taskSQL = append(taskSQL, fmt.Sprintf("content.msg like %q", "%"+task.KeyWords+"%"))
		}
		if task.IPs != "" {
			ips := xstr.JoinInts(ipsToInts(task.IPs))
			taskSQL = append(taskSQL, fmt.Sprintf("content.ip in (%s)", ips))
		}
		if task.Mids != "" {
			taskSQL = append(taskSQL, fmt.Sprintf("index.mid in (%s)", task.Mids))
		}
		if task.Cids != "" {
			taskSQL = append(taskSQL, fmt.Sprintf("index.oid in (%s)", task.Cids))
		}
		if task.Start != "" {
			if sTime, err = time.ParseInLocation("2006-01-02 15:04:05", task.Start, time.Local); err != nil {
				return
			}
			taskSQL = append(taskSQL, fmt.Sprintf("content.log_date>=%s", sTime.Format("20060102")))
			taskSQL = append(taskSQL, fmt.Sprintf("content.ctime>=%q", task.Start))
		}
		if task.End != "" {
			if eTime, err = time.ParseInLocation("2006-01-02 15:04:05", task.End, time.Local); err != nil {
				return
			}
			taskSQL = append(taskSQL, fmt.Sprintf("content.log_date<=%s", eTime.Format("20060102")))
			taskSQL = append(taskSQL, fmt.Sprintf("content.ctime<=%q", task.End))
		}
		if v.Topic, err = s.dao.SendTask(c, taskSQL); err != nil {
			//	err = nil
			//	v.State = model.TaskStateFailed
			return
		}
	}
	_, err = s.dao.ReviewTask(c, v)
	return
}

// EditTaskState .
func (s *Service) EditTaskState(c context.Context, v *model.EditTasksStateArg) (err error) {
	if _, err = s.dao.EditTaskState(c, v); err != nil {
		return
	}
	if v.State == model.TaskStateRun {
		_, err = s.dao.EditTaskPriority(c, v.IDs, time.Now().Unix())
	}
	return
}

// TaskView .
func (s *Service) TaskView(c context.Context, v *model.TaskViewArg) (task *model.TaskView, err error) {
	if task, err = s.dao.TaskView(c, v.ID); err != nil || task == nil {
		return
	}
	if task.SubTask, err = s.dao.SubTask(c, task.ID); err != nil || task.SubTask == nil {
		return
	}
	task.Tcount = task.SubTask.Tcount
	return
}

// TaskCsv .
func (s *Service) TaskCsv(c context.Context, id int64) (bs []byte, err error) {
	var (
		task      *model.TaskView
		buf       *bytes.Buffer
		csvWriter *csv.Writer
	)
	if task, err = s.dao.TaskView(c, id); err != nil {
		return
	}
	if task == nil || len(task.Result) == 0 {
		err = ecode.NothingFound
		return
	}
	res, err := http.Get(task.Result)
	if err != nil {
		log.Error("s.HttpGet(%s) error(%v)", task.Result, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = ecode.NothingFound
		return
	}
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("s.ioutilRead error(%v)", err)
		return
	}
	buf = &bytes.Buffer{}
	csvWriter = csv.NewWriter(buf)
	if err = csvWriter.Write(_csvTaskTitle); err != nil {
		log.Error("csvWriter.Write(%v) erorr(%v)", _csvTaskTitle, err)
		return
	}
	lines := bytes.Split(resp, []byte("\n"))
	if len(lines) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, line := range lines {
		var (
			items []string
		)
		fields := bytes.Split(line, []byte("\001"))
		if len(fields) < 7 {
			log.Error("fields lenth too small:%d", len(fields))
			continue
		}
		for _, field := range fields {
			items = append(items, string(field))
		}
		if err = csvWriter.Write(items); err != nil {
			log.Error("csvWriter.Write(%v) erorr(%v)", items, err)
			return
		}
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		log.Error("csvWriter.Error(%v)", err)
		return
	}
	bs = buf.Bytes()
	if len(bs) == 0 {
		err = ecode.NothingFound
	}
	return
}

func ipsToInts(ips string) (ipInts []int64) {
	ipStrs := strings.Split(ips, ",")
	ipInts = make([]int64, 0, len(ipStrs))
	for _, ipStr := range ipStrs {
		ipInts = append(ipInts, ipToInt(ipStr))
	}
	return
}

func ipToInt(ip string) (ipInt int64) {
	ret := big.NewInt(0)
	if net.ParseIP(ip) == nil {
		return
	}
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}
