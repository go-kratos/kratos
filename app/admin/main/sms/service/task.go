package service

import (
	"context"
	"os"
	"strconv"
	"strings"

	pb "go-common/app/service/main/sms/api"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const _tableTask = "sms_tasks"

func (s *Service) checkActTemplate(ctx context.Context, code string) (err error) {
	tpl, err := s.templateByCode(ctx, code)
	if err != nil {
		return
	}
	if tpl == nil {
		return ecode.SmsTemplateNotExist
	}
	if tpl.Stype != smsmdl.TypeActSms {
		return ecode.SmsTemplateNotAct
	}
	return
}

// AddTask add task
func (s *Service) AddTask(ctx context.Context, req *pb.AddTaskReq) (res *pb.AddTaskReply, err error) {
	if err = s.checkActTemplate(ctx, req.TemplateCode); err != nil {
		return
	}
	t := &smsmdl.ModelTask{
		Type:         req.Type,
		BusinessID:   req.BusinessID,
		TemplateCode: req.TemplateCode,
		Desc:         req.Desc,
		FileName:     req.FileName,
		FilePath:     req.FilePath,
		SendTime:     xtime.Time(req.SendTime),
		Status:       smsmdl.TaskStatusPrepared,
	}
	if err = s.db.Table(_tableTask).Create(t).Error; err != nil {
		log.Error("s.AddTask(%+v) error(%v)", req, err)
	}
	return
}

// UpdateTask update task
func (s *Service) UpdateTask(ctx context.Context, req *pb.UpdateTaskReq) (res *pb.UpdateTaskReply, err error) {
	if err = s.checkActTemplate(ctx, req.TemplateCode); err != nil {
		return
	}
	data := map[string]interface{}{
		"type":          req.Type,
		"business_id":   req.BusinessID,
		"template_code": req.TemplateCode,
		"desc":          req.Desc,
		"file_name":     req.FileName,
		"file_path":     req.FilePath,
		"send_time":     xtime.Time(req.SendTime),
	}
	if err = s.db.Table(_tableTask).Where("id=?", req.ID).Update(data).Error; err != nil {
		log.Error("s.UpdateTask(%+v) error(%v)", req, err)
	}
	return
}

// DeleteTask delete task
func (s *Service) DeleteTask(ctx context.Context, req *pb.DeleteTaskReq) (res *pb.DeleteTaskReply, err error) {
	if err = s.db.Table(_tableTask).Where("id=?", req.ID).Update("status", smsmdl.TaskStatusStop).Error; err != nil {
		log.Error("s.DeleteTask(%+v) error(%v)", req, err)
	}
	return
}

// TaskInfo get task info
func (s *Service) TaskInfo(ctx context.Context, req *pb.TaskInfoReq) (res *pb.TaskInfoReply, err error) {
	res = &pb.TaskInfoReply{Info: new(smsmdl.ModelTask)}
	if err = s.db.Table(_tableTask).Where("id=?", req.ID).Find(&res.Info).Error; err != nil {
		log.Error("s.TaskInfo(%+v) error(%v)", req, err)
		return
	}
	tpl, err := s.templateByCode(ctx, res.Info.TemplateCode)
	if err != nil || tpl == nil {
		return
	}
	res.Info.TemplateCode = tpl.Code
	res.Info.TemplateContent = tpl.Template
	return
}

// TaskList get task list
func (s *Service) TaskList(ctx context.Context, req *pb.TaskListReq) (res *pb.TaskListReply, err error) {
	res = &pb.TaskListReply{List: make([]*smsmdl.ModelTask, 0)}
	start := (req.Pn - 1) * req.Ps
	if err = s.db.Table(_tableTask).Order("id desc").Offset(start).Limit(req.Ps).Find(&res.List).Error; err != nil {
		log.Error("s.TaskList(%d,%d) error(%v)", req.Pn, req.Ps, err)
		return
	}
	if err = s.db.Table(_tableTask).Count(&res.Total).Error; err != nil {
		log.Error("s.TaskList(%d,%d) count error(%v)", req.Pn, req.Ps, err)
		return
	}
	return
}

// 上传说明：
// 前端是批量上传，会随机按内容长度切割文件进行分批上传，有可能会切断原始行内容
// 如果想在上传的时候同步判断文件格式，产生错误时需要忽略首行和末行

// CheckUploadFile checks uploaded content validation.
func (s *Service) CheckUploadFile(data []byte) (err error) {
	var (
		val     int64
		lineNum int
		lines   = strings.Split(string(data), "\n")
		total   = len(lines)
	)
	for _, line := range lines {
		lineNum++
		line = strings.Trim(line, " \r\t")
		if line == "" {
			continue
		}
		if val, err = strconv.ParseInt(line, 10, 64); err != nil {
			log.Error("CheckUploadMid data(%s) error(%v)", line, err)
			return ecode.PushUploadInvalidErr
		}
		if val <= 0 {
			if lineNum == 1 || lineNum == total {
				continue
			}
			log.Error("CheckUploadMid data(%s) error(%v)", line, err)
			return ecode.PushUploadInvalidErr
		}
	}
	return
}

// Upload add mids file.
func (s *Service) Upload(dir, path string, data []byte) (err error) {
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
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("s.Upload(%s) error(%v)", path, err)
		return
	}
	if _, err = f.Write(data); err != nil {
		log.Error("f.Write() error(%v)", err)
		return
	}
	if err = f.Close(); err != nil {
		log.Error("f.Close() error(%v)", err)
	}
	return
}
