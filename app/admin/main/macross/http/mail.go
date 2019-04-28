package http

import (
	"encoding/json"
	"fmt"
	"go-common/app/admin/main/macross/model/mail"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"io/ioutil"
)

// sendMail  send mail
func sendmail(c *bm.Context) {
	req := c.Request
	res := map[string]interface{}{}
	res["message"] = "success"

	var attach *mail.Attach
	// 附件
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		defer file.Close()
		attach = &mail.Attach{}
		attach.Name = header.Filename
		attach.File = file
		unzip := c.Request.Form.Get("unzip")
		if unzip != "" && unzip != "0" {
			attach.ShouldUnzip = true
		} else {
			attach.ShouldUnzip = false
		}
	}

	var bs []byte
	if attach == nil {
		bs, err = ioutil.ReadAll(req.Body)
	} else {
		// 使用 multipart 上传附件时，body 并不是 json，因此原来的 json 放在 form 的 json_body 中
		jsonBody := c.Request.Form.Get("json_body")
		bs = []byte(jsonBody)
	}
	if err != nil {
		log.Error("ioutil.ReadAll() error(%v)", err)
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	req.Body.Close()

	// params
	var m = &mail.Mail{}
	if err = json.Unmarshal(bs, m); err != nil {
		log.Error("http sendmail() json.Unmarshal(%s) error(%v)", string(bs), err)
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if m.Subject == "" || m.Body == "" || len(m.ToAddresses) == 0 {
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = svr.SendMail(c, m, attach); err != nil {
		res["message"] = fmt.Sprintf("%v", err)
		c.JSONMap(res, err)
		return
	}
	c.JSONMap(res, nil)
}
