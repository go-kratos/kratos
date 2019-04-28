package service

import (
	"context"
	"io/ioutil"
	"strings"
	"time"

	"go-common/app/service/main/sms/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

func (s *Service) loadTaskproc() {
	for {
		if !s.c.Sms.PickUpTask {
			log.Warn("service do not pick up new tasks from database")
			return
		}
		task, err := s.pickNewTask()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if task != nil {
			tpl := s.template[task.TemplateCode]
			if tpl == nil {
				log.Error("template not exists, code(%s)", task.TemplateCode)
				continue
			}
			task.TemplateContent = tpl.Template
			if err = s.handleTask(task); err == nil {
				s.dao.UpdateTaskStatus(context.Background(), task.ID, model.TaskStatusSuccess)
			} else {
				s.dao.UpdateTaskStatus(context.Background(), task.ID, model.TaskStatusFailed)
			}
		}
		time.Sleep(time.Duration(s.c.Sms.LoadTaskInteval))
	}
}

func (s *Service) pickNewTask() (task *model.ModelTask, err error) {
	ctx := context.Background()
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTx(ctx); err != nil {
		log.Error("tx.BeginTx() error(%v)", err)
		return
	}
	if task, err = s.dao.TxTask(tx); err != nil {
		tx.Rollback()
		return
	}
	if task == nil {
		tx.Rollback()
		return
	}
	if err = s.dao.TxUpdateTaskStatus(tx, task.ID, model.TaskStatusDoing); err != nil {
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit() task(%+v) error(%v)", task, err)
		return
	}
	return
}

func (s *Service) handleTask(task *model.ModelTask) (err error) {
	bs, err := ioutil.ReadFile(task.FilePath)
	if err != nil {
		log.Error("ioutil.ReadFile(%s) error(v)", task.FilePath, err)
		return
	}
	var (
		counter int
		group   = errgroup.Group{}
		data    []string
	)
	for _, v := range strings.Split(string(bs), "\n") {
		v = strings.Trim(v, " \r\t")
		if v == "" {
			continue
		}
		data = append(data, v)
	}
	for {
		l := len(data)
		if l == 0 {
			break
		}
		n := s.c.Sms.BatchSize
		if l < n {
			n = l
		}
		part := data[:n]
		data = data[n:]
		group.Go(func() error {
			s.sendBatch(task, part)
			return nil
		})
		counter++
		if counter > s.c.Sms.TaskWorker {
			group.Wait()
			counter = 0
		}
	}
	if counter > 0 {
		group.Wait()
	}
	return
}

func (s *Service) sendBatch(task *model.ModelTask, data []string) (err error) {
	send := &model.ModelSend{Type: model.TypeActBatch, Code: task.TemplateCode, Content: task.TemplateContent}
	if task.Type == model.TaskTypeMid {
		send.Mid = strings.Join(data, ",")
	} else if task.Type == model.TaskTypeMobile {
		send.Mobile = strings.Join(data, ",")
	} else {
		log.Error("invalid task type, task(%+v)", task)
		return
	}
	err = s.dao.PubBatch(context.Background(), send)
	return
}
