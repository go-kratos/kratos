package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/live/resource/api/http/v1"
	"go-common/app/service/live/resource/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strings"
	"time"
)

const (
	_treeNodes = "treeNodes_name:%s"
)

// Nodes node.
func (s *TitansService) Nodes(ctx context.Context, user, node, team, cookie string) (res []*v1.Node, err error) {
	var nodes *model.CacheData

	if nodes, err = s.AuthApps(ctx, user, cookie); err != nil {
		return
	}
	//node list.
	if node == "" && team == "" {

		tmp := make(map[string]struct{})
		for _, app := range nodes.Data {
			idx := strings.Index(app.Path, ".")
			bu := string([]byte(app.Path)[:idx])
			if _, ok := tmp[bu]; ok {
				continue
			}
			n := new(v1.Node)
			n.Name = bu
			n.Path = bu
			res = append(res, n)
			tmp[bu] = struct{}{}
		}
		return
	}
	//team list.
	if node != "" && team == "" {
		tmp := make(map[string]struct{})
		for _, app := range nodes.Data {
			s := []byte(app.Path)
			sep := []byte(".")
			fidx := bytes.Index(s, sep)
			lidx := bytes.LastIndex(s, sep)
			team = string(s[:lidx])
			if node == string(s[:fidx]) {
				if _, ok := tmp[team]; ok {
					continue
				}
				n := new(v1.Node)
				n.Name = string([]byte(app.Path)[fidx+1 : lidx])
				n.Path = team
				tmp[team] = struct{}{}
				res = append(res, n)
			}
		}
		return
	}
	//app list.
	if team == "" {
		return
	}
	for _, app := range nodes.Data {
		s := []byte(app.Path)
		sep := []byte(".")
		lidx := bytes.LastIndex(s, sep)
		t := string(s[:lidx])
		if team != t {
			continue
		}
		n := new(v1.Node)
		n.Name = string(s[lidx+1:])
		n.Path = app.Path
		n.TreeId = app.ID
		res = append(res, n)
	}
	return
}

// AuthApps 获取用户节点权限.
func (s *TitansService) AuthApps(ctx context.Context, user string, cookie string) (nodes *model.CacheData, err error) {
	if len(user) == 0 {
		err = ecode.NothingFound
		return
	}
	var ok bool
	cacheKey := fmt.Sprintf(_treeNodes, user)
	nodesCache, ok := s.treeCache.Get(cacheKey)
	cacheStr, _ := json.Marshal(nodesCache)
	nodes = &model.CacheData{}
	json.Unmarshal(cacheStr, nodes)
	if !ok || (time.Since(nodes.CTime) > 60*time.Second) {
		log.Info("[Titans][Tree] miss lruCache call for tree service")
		nodes, err = s.SyncTree(ctx, user, cookie)
		s.treeCache.Put(cacheKey, nodes)
	}
	return
}

//SyncTree syncTree.
func (s *TitansService) SyncTree(ctx context.Context, user string, cookie string) (nodes *model.CacheData, err error) {
	var (
		msg   map[string]interface{}
		tmp   interface{}
		token string
		ok    bool
	)
	if msg, err = s.dao.Auth(ctx, cookie); err != nil {
		return
	}
	if tmp, ok = msg["token"]; !ok {
		err = ecode.NothingFound
		return
	}
	if token, ok = tmp.(string); !ok {
		err = ecode.NothingFound
		return
	}
	//data, _ := s.dao.Tree(ctx, token)
	if nodes, err = s.dao.Role(ctx, user, token); err != nil {
		return
	}

	return
}
