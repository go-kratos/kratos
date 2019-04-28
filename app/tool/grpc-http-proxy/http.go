package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"
)

var jsonFormat = "application/json; charset=utf-8"

func splitURL(url string) (addr, path string, err error) {
	arrs := strings.Split(url, "/")
	if len(arrs) < 3 {
		err = ecode.RequestErr
		return
	}
	addr = arrs[1]
	path = fmt.Sprintf("/%s", strings.Join(arrs[2:], "/"))
	return
}

func formatRequest(req url.Values) (res map[string]interface{}) {
	res = map[string]interface{}{}
	for k, v := range req {
		value := v[0]
		res[k] = value
		// 支持数组
		if len(value) > 2 && value[0] == '[' && value[len(value)-1] == ']' {
			value = value[1 : len(value)-1]
			res[k] = strings.Split(value, ",")
		}
	}
	return
}

func handle(c *bm.Context) {
	addr, path, err := splitURL(c.Request.URL.Path)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	var body []byte
	body, _ = ioutil.ReadAll(c.Request.Body)
	if len(body) == 0 || string(body) == "{}" {
		data := formatRequest(c.Request.Form)
		body, err = json.Marshal(data)
		if err != nil {
			c.JSONMap(map[string]interface{}{"message": err.Error()}, err)
			return
		}
	}
	resp, err := callGrpc(addr, path, body)
	if err != nil {
		var intro string
		if strings.Contains(err.Error(), "error unmarshalling request") {
			intro = "由于params都是string没有类型信息 如果含有int/int32等类型 请复制\"请求数据\"到body中 并修改请求的字段类型, 注意Content-Type 不能为 multipart/form-data 不然body无效 还有问题@wangxu01"
		}
		c.JSONMap(map[string]interface{}{
			"错误":   err.Error(),
			"请求目标": addr,
			"请求方法": path,
			"请求数据": string(body),
			"额外说明": intro,
		}, err)
		return
	}
	c.Bytes(http.StatusOK, jsonFormat, resp)
}
