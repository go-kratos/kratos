package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-common/app/admin/main/cache/conf"
	"go-common/app/admin/main/cache/dao"
	"go-common/app/admin/main/cache/model"
	"go-common/library/ecode"

	"github.com/BurntSushi/toml"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao

	opsMcs []*model.OpsCacheMemcache
	opsRds []*model.OpsCacheRedis
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	go s.loadOpsproc()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) appids(c context.Context, cookie, appid string) (appids []string, err error) {
	msg, err := s.dao.Auth(c, cookie)
	if err != nil {
		err = ecode.AccessDenied
		return
	}
	tmp, ok := msg["token"]
	if !ok {
		err = ecode.NothingFound
		return
	}
	token, ok := tmp.(string)
	if !ok {
		err = ecode.NothingFound
		return
	}
	nodes, err := s.dao.Role(c, token)
	if err != nil {
		return
	}
	if appid == "" {
		for _, node := range nodes.Data {
			appids = append(appids, node.Path)
		}
	} else {
		for _, node := range nodes.Data {
			if appid == node.Path {
				appids = []string{appid}
				break
			}
		}
	}
	return
}

func (s *Service) treeid(c context.Context, cookie, appid string) (treeid int64, err error) {
	if appid == "" {
		err = ecode.AccessDenied
		return
	}
	msg, err := s.dao.Auth(c, cookie)
	if err != nil {
		err = ecode.AccessDenied
		return
	}
	tmp, ok := msg["token"]
	if !ok {
		err = ecode.NothingFound
		return
	}
	token, ok := tmp.(string)
	if !ok {
		err = ecode.NothingFound
		return
	}
	nodes, err := s.dao.Role(c, token)
	if err != nil {
		return
	}
	for _, node := range nodes.Data {
		if appid == node.Path {
			treeid = node.ID
			return
		}
	}
	return
}

// Clusters get clusters.
func (s *Service) Clusters(c context.Context, req *model.ClusterReq) (resp *model.ClusterResp, err error) {
	appids, err := s.appids(c, req.Cookie, req.AppID)
	if err != nil {
		err = ecode.AccessDenied
		return
	}
	resp = new(model.ClusterResp)
	if len(appids) == 0 {
		return
	}
	if err = s.dao.DB.Where("appids in (?) AND zone=? AND type=?", appids, req.Zone, req.Type).Order("id").Offset((req.PN - 1) * req.PS).Limit(req.PS).Find(&resp.Clusters).Error; err != nil {
		return
	}
	var count int64
	s.dao.DB.Model(&model.Cluster{}).Where("appids in (?) AND zone=? AND type=?", appids, req.Zone, req.Type).Count(&count)
	resp.Total = count
	return
}

// AddCluster add new cluster.
func (s *Service) AddCluster(c context.Context, req *model.AddClusterReq) (resp *model.EmpResp, err error) {
	cluster := &model.Cluster{
		Name:             req.Name,
		Type:             req.Type,
		AppID:            req.AppID,
		Zone:             req.Zone,
		HashMethod:       req.HashMethod,
		HashDistribution: req.HashDistribution,
		HashTag:          req.HashTag,
		DailTimeout:      req.DailTimeout,
		ReadTimeout:      req.ReadTimeout,
		WriteTimeout:     req.WriteTimeout,
		NodeConn:         req.NodeConn,
		ListenAddr:       req.ListenAddr,
		ListenProto:      req.ListenProto,
		PingFailLimit:    req.PingFailLimit,
		PingAutoEject:    req.PingAutoEject,
	}

	if req.ID == 0 {
		err = s.dao.DB.Create(cluster).Error
	} else {
		cluster.ID = req.ID
		s.dao.DB.Save(cluster)
	}
	return
}

// DelCluster del cluster of req id.
func (s *Service) DelCluster(c context.Context, req *model.DelClusterReq) (resp *model.EmpResp, err error) {
	err = s.dao.DB.Exec("DELETE FROM cluster where id= ?", req.ID).Error
	if err != nil {
		return
	}
	err = s.dao.DB.Exec("DELETE FROM nodes where cid= ?", req.ID).Error
	return
}

// Cluster search cluster by appid or cluster name.
func (s *Service) Cluster(c context.Context, req *model.ClusterReq) (resp []*model.Cluster, err error) {
	if req.Type != "" {
		err = s.dao.DB.Where("appids=? AND zone=? AND type=?", req.AppID, req.Zone, req.Type).Find(&resp).Error
	} else {
		err = s.dao.DB.Where("appids=? AND zone=?", req.AppID, req.Zone).Find(&resp).Error
	}
	if err != nil {
		return
	}
	for _, clu := range resp {
		err = s.dao.DB.Where("cid = ?", clu.ID).Find(&clu.Nodes).Error
		if err != nil {
			return nil, err
		}
	}
	return
}

// ModifyCluster add or del cluster nodes.
func (s *Service) ModifyCluster(c context.Context, req *model.ModifyClusterReq) (resp *model.EmpResp, err error) {
	var nodes []*model.NodeDtl
	err = json.Unmarshal([]byte(req.Nodes), &nodes)
	if err != nil {
		return
	}
	var id = req.ID
	if req.Name != "" {
		var cluster = &model.Cluster{}
		err = s.dao.DB.Where("name = ?", req.Name).First(cluster).Error
		if err != nil {
			return
		}
		id = cluster.ID
	}
	if req.Action == 2 {
		var alias []string
		for _, ali := range nodes {
			alias = append(alias, ali.Alias)
		}
		//err = s.dao.DB.Delete(&nodes).Error
		err = s.dao.DB.Exec("DELETE FROM nodes WHERE alias in (?) ", strings.Join(alias, ",")).Error
		return
	} else if req.Action == 1 {
		// var nodes []*model.NodeDtl
		for _, node := range nodes {
			node.Cid = id
			err = s.dao.DB.Create(node).Error
		}
		return
	}
	return
}

// ClusterDtl get cluster detail about nodes info.
func (s *Service) ClusterDtl(c context.Context, req *model.ClusterDtlReq) (resp *model.ClusterDtlResp, err error) {
	resp = new(model.ClusterDtlResp)
	err = s.dao.DB.Where("cid = ?", req.ID).Find(&resp.Nodes).Error
	// TODO(lintanghui):get node info
	return
}

// Toml return a toml file of cluster infos.
func (s *Service) Toml(c context.Context, req *model.ClusterReq) (resp []byte, err error) {
	clusters, err := s.Cluster(c, req)
	if err != nil {
		return
	}
	for _, cluster := range clusters {
		for _, node := range cluster.Nodes {
			cluster.Servers = append(cluster.Servers, fmt.Sprintf("%s:%d %s", node.Addr, node.Weight, node.Alias))
		}
	}
	buf := bytes.NewBuffer(resp)
	t := struct {
		Clusters []*model.Cluster `toml:"clusters"`
	}{
		Clusters: clusters,
	}
	err = toml.NewEncoder(buf).Encode(t)
	resp = buf.Bytes()
	return
}
