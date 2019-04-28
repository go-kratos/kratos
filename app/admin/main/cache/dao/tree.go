package dao

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-common/app/admin/main/cache/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	tokenURI = "http://easyst.bilibili.co/v1/token"
	dataURI  = "http://easyst.bilibili.co/v1/node/apptree"
	authURI  = "http://easyst.bilibili.co/v1/auth"
	nodeURI  = "http://easyst.bilibili.co/v1/node/bilibili%s"
	appsURI  = "http://easyst.bilibili.co/v1/node/role/app"
	prefix   = []byte("bilibili.")
)

// Token get Token.
func (d *Dao) Token(c context.Context, body string) (msg map[string]interface{}, err error) {
	var (
		req *http.Request
	)
	if req, err = http.NewRequest("POST", tokenURI, strings.NewReader(body)); err != nil {
		log.Error("Token url(%s) error(%v)", tokenURI, err)
		return
	}
	var res struct {
		Code    int                    `json:"code"`
		Data    map[string]interface{} `json:"data"`
		Message string                 `json:"message"`
		Status  int                    `json:"status"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) res(%+v) err(%v)", tokenURI, res, err)
		return
	}
	if res.Code != 90000 {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", tokenURI, res)
		return
	}
	msg = res.Data
	return
}

// Auth get Token.
func (d *Dao) Auth(c context.Context, cookie string) (msg map[string]interface{}, err error) {
	var (
		req *http.Request
	)
	if req, err = http.NewRequest("GET", authURI, nil); err != nil {
		log.Error("Token url(%s) error(%v)", tokenURI, err)
		return
	}
	req.Header.Set("Cookie", cookie)
	var res struct {
		Code    int                    `json:"code"`
		Data    map[string]interface{} `json:"data"`
		Message string                 `json:"message"`
		Status  int                    `json:"status"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.Token url(%s) res(%s) err(%v)", tokenURI, res, err)
		return
	}
	if res.Code != 90000 {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", tokenURI, res)
		return
	}
	msg = res.Data
	return
}

// Tree get service tree.
func (d *Dao) Tree(c context.Context, token string) (data interface{}, err error) {
	var (
		req *http.Request
		tmp map[string]interface{}
		ok  bool
	)
	if req, err = http.NewRequest("GET", dataURI, nil); err != nil {
		log.Error("Status url(%s) error(%v)", dataURI, err)
		return
	}
	req.Header.Set("X-Authorization-Token", token)
	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code    int                               `json:"code"`
		Data    map[string]map[string]interface{} `json:"data"`
		Message string                            `json:"message"`
		Status  int                               `json:"status"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.Status url(%s) res($s) err(%v)", dataURI, res, err)
		return
	}
	if res.Code != 90000 {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", dataURI, res)
		return
	}
	if tmp, ok = res.Data["bilibili"]; ok {
		data, ok = tmp["children"]
	}
	if !ok {
		err = ecode.NothingFound
	}
	return
}

// Role get service tree.
func (d *Dao) Role(c context.Context, token string) (nodes *model.CacheData, err error) {
	var (
		req *http.Request
	)
	if req, err = http.NewRequest("GET", appsURI, nil); err != nil {
		log.Error("Status url(%s) error(%v)", dataURI, err)
		return
	}
	req.Header.Set("X-Authorization-Token", token)
	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code    int               `json:"code"`
		Data    []*model.RoleNode `json:"data"`
		Message string            `json:"message"`
		Status  int               `json:"status"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.Status url(%s) res($s) err(%v)", dataURI, res, err)
		return
	}
	if res.Code != 90000 {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) res(%v)", dataURI, res)
		return
	}
	nodes = &model.CacheData{Data: make(map[int64]*model.RoleNode)}
	nodes.CTime = time.Now()
	for _, node := range res.Data {
		if bytes.Equal(prefix, []byte(node.Path)[0:9]) {
			node.Path = string([]byte(node.Path)[9:])
		}
		nodes.Data[node.ID] = node
	}
	return
}

// NodeTree get service tree.
func (d *Dao) NodeTree(c context.Context, token, bu, team string) (nodes []*model.Node, err error) {
	var (
		req  *http.Request
		node string
	)
	if len(bu) != 0 {
		node = "." + bu
	}
	if len(team) != 0 {
		node = "." + team
	}
	if len(node) == 0 {
		nodes = append(nodes, &model.Node{Name: "main", Path: "main"})
		nodes = append(nodes, &model.Node{Name: "ai", Path: "ai"})
		return
	}
	if req, err = http.NewRequest("GET", fmt.Sprintf(nodeURI, node), nil); err != nil {
		log.Error("Status url(%s) error(%v)", dataURI, err)
		return
	}
	req.Header.Set("X-Authorization-Token", token)
	req.Header.Set("Content-Type", "application/json")
	var res struct {
		Code    int        `json:"code"`
		Data    *model.Res `json:"data"`
		Message string     `json:"message"`
		Status  int        `json:"status"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.Status url(%s) res($s) err(%v)", dataURI, res, err)
		return
	}
	if res.Code != 90000 {
		err = fmt.Errorf("error code :%d", res.Code)
		log.Error("Status url(%s) error(%v)", dataURI, err)
		return
	}
	for _, tnode := range res.Data.Data {
		if bytes.Equal(prefix, []byte(tnode.Path)[0:9]) {
			tnode.Path = string([]byte(tnode.Path)[9:])
		}
		nodes = append(nodes, &model.Node{Name: tnode.Name, Path: tnode.Path})
	}
	return
}
