package http

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"go-common/app/admin/main/sms/conf"
	pb "go-common/app/service/main/sms/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func addTask(ctx *bm.Context) {
	req := new(pb.AddTaskReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(svc.AddTask(ctx, req))
}

func updateTask(ctx *bm.Context) {
	req := new(pb.UpdateTaskReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(svc.UpdateTask(ctx, req))
}

func deleteTask(ctx *bm.Context) {
	req := new(pb.DeleteTaskReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(svc.DeleteTask(ctx, req))
}

func taskInfo(ctx *bm.Context) {
	req := new(pb.TaskInfoReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	res, err := svc.TaskInfo(ctx, req)
	ctx.JSON(res.Info, err)
}

func taskList(ctx *bm.Context) {
	req := new(pb.TaskListReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	res, err := svc.TaskList(ctx, req)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	pager := struct {
		Pn    int32 `json:"page"`
		Ps    int32 `json:"pagesize"`
		Total int32 `json:"total"`
	}{
		Pn:    req.Pn,
		Ps:    req.Ps,
		Total: res.Total,
	}
	data := map[string]interface{}{
		"data":  res.List,
		"pager": pager,
	}
	ctx.JSONMap(data, nil)
}

func upload(ctx *bm.Context) {
	var (
		err error
		req = ctx.Request
	)
	req.ParseMultipartForm(1024 * 1024 * 1024) // 1G
	fileName := req.FormValue("filename")
	if fileName == "" {
		log.Error("filename is empty")
		ctx.JSON(nil, ecode.RequestErr)
		return
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		log.Error("req.FormFile() error(%v)", err)
		ctx.JSON(nil, err)
		return
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		ctx.JSON(nil, err)
		return
	}
	if err = svc.CheckUploadFile(bs); err != nil {
		ctx.JSON(nil, err)
		return
	}
	dir := fmt.Sprintf("%s/%s", strings.TrimSuffix(conf.Conf.Sms.MountDir, "/"), time.Now().Format("20060102"))
	path := fmt.Sprintf("%s/%x", dir, md5.Sum([]byte(fileName)))
	if err = svc.Upload(dir, path, bs); err != nil {
		log.Error("upload file file(%s) error(%v)", path, err)
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}{
		Name: header.Filename,
		Path: path,
	}, nil)
}
