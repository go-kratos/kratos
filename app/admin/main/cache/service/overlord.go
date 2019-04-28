package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"go-common/app/admin/main/cache/model"
	"go-common/library/ecode"

	"github.com/BurntSushi/toml"
)

// OpsClusterNames .
func (s *Service) OpsClusterNames(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	resp = &model.OverlordResp{}
	if arg.Type == "memcache" {
		for _, opsmc := range s.opsMcs {
			resp.Names = append(resp.Names, opsmc.Labels.Name)
		}
	} else if arg.Type == "redis" {
		for _, opsrd := range s.opsRds {
			if opsrd.Type == "redis_standalone" {
				resp.Names = append(resp.Names, opsrd.Labels.Name)
			}
		}
	} else if arg.Type == "redis_cluster" {
		for _, opsrd := range s.opsRds {
			if opsrd.Type == "redis_cluster" {
				resp.Names = append(resp.Names, opsrd.Labels.Name)
			}
		}
	} else {
		err = fmt.Errorf("unsupport type:%s", arg.Type)
	}
	return
}

// OpsClusterNodes .
func (s *Service) OpsClusterNodes(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	resp = &model.OverlordResp{}
	if arg.Type == "memcache" {
		for _, opsmc := range s.opsMcs {
			if arg.Name == opsmc.Labels.Name {
				resp.Addrs = opsmc.Targets
				return
			}
		}
	} else if arg.Type == "redis" {
		for _, opsrd := range s.opsRds {
			if opsrd.Type == "redis_standalone" && arg.Name == opsrd.Labels.Name {
				resp.Addrs = opsrd.Targets
				return
			}
		}
	} else if arg.Type == "redis_cluster" {
		for _, opsrd := range s.opsRds {
			if opsrd.Type == "redis_cluster" && arg.Name == opsrd.Labels.Name {
				resp.Addrs = opsrd.Targets
				return
			}
		}
	} else {
		err = fmt.Errorf("unsupport type:%s", arg.Type)
	}
	return
}

// ImportOpsCluster .
func (s *Service) ImportOpsCluster(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	exist := 0
	if err = s.dao.DB.Model(&model.OverlordCluster{}).Where("name=?", arg.Name).Count(&exist).Error; err != nil {
		return
	}
	if exist > 0 {
		return
	}
	var targets []string
	if arg.Type == "memcache" {
		for _, opsmc := range s.opsMcs {
			if arg.Name == opsmc.Labels.Name {
				targets = opsmc.Targets
				break
			}
		}
	} else if arg.Type == "redis" {
		for _, opsrd := range s.opsRds {
			if opsrd.Type == "redis_standalone" && arg.Name == opsrd.Labels.Name {
				targets = opsrd.Targets
				break
			}
		}
	} else if arg.Type == "redis_cluster" {
		for _, opsrd := range s.opsRds {
			if opsrd.Type == "redis_cluster" && arg.Name == opsrd.Labels.Name {
				targets = opsrd.Targets
				break
			}
		}
	} else {
		err = fmt.Errorf("unsupport type:%s", arg.Type)
		return
	}
	port := 0
	if err = s.dao.DB.Model(&model.OverlordCluster{}).Where("type=?", arg.Type).Count(&port).Error; err != nil {
		return
	}
	if arg.Type == "memcache" {
		port += 11211
	} else {
		port += 26379
	}
	tranDB := s.dao.DB.Begin()
	oc := &model.OverlordCluster{
		Name:             arg.Name,
		Type:             arg.Type,
		Zone:             arg.Zone,
		HashMethod:       "fnv1a_64",
		HashDistribution: "ketama",
		HashTag:          "",
		ListenProto:      "tcp",
		ListenAddr:       "0.0.0.0:" + strconv.Itoa(port),
		DailTimeout:      1000,
		ReadTimeout:      1000,
		WriteTimeout:     1000,
		NodeConn:         2,
		PingFailLimit:    3,
		PingAutoEject:    true,
	}
	if err = tranDB.Create(oc).Error; err != nil {
		tranDB.Rollback()
		return
	}
	for i, target := range targets {
		ocn := &model.OverlordNode{
			Cid:    oc.ID,
			Alias:  arg.Name + strconv.Itoa(i+1),
			Addr:   target,
			Weight: 1,
		}
		if err = tranDB.Create(ocn).Error; err != nil {
			tranDB.Rollback()
			return
		}
	}
	err = tranDB.Commit().Error
	return
}

// OverlordClusters .
func (s *Service) OverlordClusters(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	resp = &model.OverlordResp{}
	if arg.Name != "" {
		err = s.dao.DB.Where("zone=? AND type=? AND name like ?", arg.Zone, arg.Type, "%"+arg.Name+"%").Order("id desc").Offset((arg.PN - 1) * arg.PS).Limit(arg.PS).Find(&resp.Clusters).Error
		s.dao.DB.Model(&model.OverlordCluster{}).Where("zone=? AND type=? AND name like ?", arg.Zone, arg.Type, arg.Name).Count(&resp.Total)
	} else {
		err = s.dao.DB.Where("zone=? AND type=?", arg.Zone, arg.Type).Order("id desc").Offset((arg.PN - 1) * arg.PS).Limit(arg.PS).Find(&resp.Clusters).Error
		s.dao.DB.Model(&model.OverlordCluster{}).Where("zone=? AND type=?", arg.Zone, arg.Type).Count(&resp.Total)
	}
	if err != nil {
		return
	}
	for _, cluster := range resp.Clusters {
		var ens *model.OverlordResp
		if ens, err = s.ExistOverlordNodes(c, &model.OverlordReq{Name: cluster.Name}); err != nil {
			return
		}
		cluster.Nodes = ens.Nodes
	}
	return
}

// ExistOverlordNodes .
func (s *Service) ExistOverlordNodes(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	cluster := &model.OverlordCluster{}
	if err = s.dao.DB.Model(cluster).Where("name=?", arg.Name).First(cluster).Error; err != nil {
		fmt.Printf("get cluster err %v\n", err)
		return
	}
	if cluster.ID == 0 {
		err = fmt.Errorf("cluster not exist:%s", arg.Name)
		return
	}
	var exists []*model.OverlordNode
	if err = s.dao.DB.Where("cid=?", cluster.ID).Order("id").Find(&exists).Error; err != nil {
		return
	}
	resp = &model.OverlordResp{}
	resp.Cluster = cluster
	resp.Nodes = exists
	return
}

// NotExistOverlordAddrs .
func (s *Service) NotExistOverlordAddrs(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	cluster := &model.OverlordCluster{}
	if err = s.dao.DB.Model(cluster).Where("name=?", arg.Name).First(cluster).Error; err != nil {
		return
	}
	if cluster.ID == 0 {
		err = fmt.Errorf("cluster not exist:%s", arg.Name)
		return
	}
	var targets []string
	tp := arg.Type
	if tp == "memcache" {
		for _, opsmc := range s.opsMcs {
			if arg.Name == opsmc.Labels.Name {
				targets = opsmc.Targets
				break
			}
		}
	} else if tp == "redis" {
		for _, opsrd := range s.opsRds {
			if arg.Name == opsrd.Labels.Name {
				targets = opsrd.Targets
				if opsrd.Type == "redis_cluster" {
					tp = "redis_cluster"
				}
				break
			}
		}
	} else {
		err = fmt.Errorf("unsupport type:%s", arg.Type)
		return
	}
	var exists []*model.OverlordNode
	if err = s.dao.DB.Where("cid=?", cluster.ID).Order("id").Find(&exists).Error; err != nil {
		return
	}
	resp = &model.OverlordResp{}
NEXT:
	for _, target := range targets {
		for _, exist := range exists {
			if target == exist.Addr {
				continue NEXT
			}
		}
		resp.Addrs = append(resp.Addrs, target)
	}
	return
}

// ImportOpsNode .
func (s *Service) ImportOpsNode(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	nen, err := s.NotExistOverlordAddrs(c, arg)
	if err != nil {
		return
	}
	en, err := s.ExistOverlordNodes(c, arg)
	if err != nil {
		return
	}
	i := len(en.Nodes)
	tranDB := s.dao.DB.Begin()
	for _, target := range nen.Addrs {
		ocn := &model.OverlordNode{
			Cid:    en.Cluster.ID,
			Alias:  arg.Name + strconv.Itoa(i+1),
			Addr:   target,
			Weight: 1,
		}
		if err = tranDB.Create(ocn).Error; err != nil {
			tranDB.Rollback()
			return
		}
		i++
	}
	err = tranDB.Commit().Error
	return
}

// ReplaceOpsNode .
func (s *Service) ReplaceOpsNode(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	nen, err := s.NotExistOverlordAddrs(c, arg)
	if err != nil {
		return
	}
	if len(nen.Addrs) == 0 {
		err = fmt.Errorf("cluster have not new node:%s", arg.Name)
		return
	}
	en, err := s.ExistOverlordNodes(c, arg)
	if err != nil {
		return
	}
	for _, n := range en.Nodes {
		if n.Alias != arg.Alias && n.Addr == arg.Addr {
			err = fmt.Errorf("cluster:%s node:%s used by other node:%s ", arg.Name, arg.Addr, n.Alias)
			return
		}
	}
	node := &model.OverlordNode{}
	if err = s.dao.DB.Model(node).Where("cid=? AND alias=?", en.Cluster.ID, arg.Alias).First(node).Error; err != nil {
		return
	}
	if node.Addr == arg.Addr {
		return
	}
	err = s.dao.DB.Model(node).Where("alias=? AND addr=?", node.Alias, node.Addr).Update("addr", arg.Addr).Error
	return
}

// DelOverlordCluster .
func (s *Service) DelOverlordCluster(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	en, err := s.ExistOverlordNodes(c, arg)
	if err != nil {
		return
	}
	if err = s.dao.DB.Delete(en.Cluster).Error; err != nil {
		return
	}
	err = s.dao.DB.Delete(&model.OverlordNode{}, "cid=?", en.Cluster.ID).Error
	return
}

// DelOverlordNode .
func (s *Service) DelOverlordNode(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	en, err := s.ExistOverlordNodes(c, arg)
	if err != nil {
		return
	}
	err = s.dao.DB.Delete(&model.OverlordNode{}, "cid=? AND alias=? AND addr=?", en.Cluster.ID, arg.Alias, arg.Addr).Error
	return
}

// OverlordAppClusters .
func (s *Service) OverlordAppClusters(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	appids, err := s.appids(c, arg.Cookie, arg.AppID)
	if err != nil {
		err = ecode.AccessDenied
		return
	}
	resp = &model.OverlordResp{}
	if len(appids) <= 1 {
		err = s.dao.DB.Where("app_id like ?", "%"+arg.AppID+"%").Order("id desc").Offset((arg.PN - 1) * arg.PS).Limit(arg.PS).Find(&resp.Apps).Error
		s.dao.DB.Model(&model.OverlordApp{}).Where("app_id like ?", arg.AppID).Count(&resp.Total)
	} else if len(appids) > 1 {
		err = s.dao.DB.Where("app_id in (?)", appids).Order("id desc").Offset((arg.PN - 1) * arg.PS).Limit(arg.PS).Find(&resp.Apps).Error
		s.dao.DB.Model(&model.OverlordApp{}).Where("app_id in (?)", appids).Count(&resp.Total)
	}
	if err != nil || len(resp.Apps) == 0 {
		return
	}
	var cids []int64
	for _, app := range resp.Apps {
		cids = append(cids, app.Cid)
	}
	var clusters []*model.OverlordCluster
	if err = s.dao.DB.Find(&clusters, "id in (?)", cids).Error; err != nil {
		return
	}
	for _, cluster := range clusters {
		var ens *model.OverlordResp
		if ens, err = s.ExistOverlordNodes(c, &model.OverlordReq{Name: cluster.Name}); err != nil {
			return
		}
		cluster.Nodes = ens.Nodes
		for _, app := range resp.Apps {
			if cluster.ID == app.Cid {
				app.Cluster = cluster
			}
		}
	}
	if len(appids) <= 1 {
		// 当使用appid查询时，填充overlord-mesos的数据
		if ocs, err := s.dao.OverlordClusters(c, "", arg.AppID); err == nil {
			clusters = append(clusters, ocs...)
		}
	}
	return
}

// OverlordAppCanBindClusters .
func (s *Service) OverlordAppCanBindClusters(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	resp = &model.OverlordResp{}
	err = s.dao.DB.Where("zone=? AND type=?", arg.Zone, arg.Type).Find(&resp.Clusters).Error
	if err != nil {
		return
	}
	for _, cluster := range resp.Clusters {
		resp.Names = append(resp.Names, cluster.Name)
	}
	resp.Clusters = nil
	return
}

// OverlordAppClusterBind .
func (s *Service) OverlordAppClusterBind(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	treeid, err := s.treeid(c, arg.Cookie, arg.AppID)
	if err != nil || treeid == 0 {
		err = ecode.AccessDenied
		return
	}
	cluster := &model.OverlordCluster{}
	if err = s.dao.DB.Model(cluster).Where("zone=? AND type=? AND name=?", arg.Zone, arg.Type, arg.Name).First(cluster).Error; err != nil {
		return
	}
	app := &model.OverlordApp{
		TreeID: treeid,
		AppID:  arg.AppID,
		Cid:    cluster.ID,
	}
	err = s.dao.DB.Create(app).Error
	return
}

// OverlordAppClusterDel .
func (s *Service) OverlordAppClusterDel(c context.Context, arg *model.OverlordReq) (resp *model.EmpResp, err error) {
	treeid, err := s.treeid(c, arg.Cookie, arg.AppID)
	if err != nil || treeid == 0 {
		err = ecode.AccessDenied
		return
	}
	cluster := &model.OverlordCluster{}
	if err = s.dao.DB.Model(cluster).Where("zone=? AND type=? AND name=?", arg.Zone, arg.Type, arg.Name).First(cluster).Error; err != nil {
		return
	}
	app := &model.OverlordApp{}
	if err = s.dao.DB.Model(app).Where("app_id=? AND cid=?", arg.AppID, cluster.ID).First(app).Error; err != nil {
		return
	}
	err = s.dao.DB.Table(app.TableName()).Delete(app).Error
	return
}

// OverlordAppAppIDs .
func (s *Service) OverlordAppAppIDs(c context.Context, arg *model.OverlordReq) (resp *model.OverlordResp, err error) {
	appids, err := s.appids(c, arg.Cookie, "")
	if err != nil {
		err = ecode.AccessDenied
		return
	}
	resp = &model.OverlordResp{}
	resp.AppIDs = appids
	return
}

// OverlordToml return a toml file of cluster infos.
func (s *Service) OverlordToml(c context.Context, arg *model.OverlordReq) (resp []byte, err error) {
	var apps []*model.OverlordApp
	if err = s.dao.DB.Where("app_id=?", arg.AppID).Find(&apps).Error; err != nil {
		return
	}
	var cids []int64
	for _, app := range apps {
		cids = append(cids, app.Cid)
	}
	var clusters []*model.OverlordCluster
	// TODO(felix): 待都走overlord-mesos后干掉
	if err = s.dao.DB.Where("zone=? AND id in (?)", arg.Zone, cids).Find(&clusters).Error; err != nil {
		return
	}
	if len(clusters) == 0 {
		// TODO(felix): 待都走overlord-mesos后干掉
		if err = s.dao.DB.Where("zone='sh001' AND id in (?)", cids).Find(&clusters).Error; err != nil {
			return
		}
	}
	var ocs []*model.OverlordCluster
	if ocs, err = s.dao.OverlordClusters(c, arg.Zone, arg.AppID); err == nil {
		if len(ocs) == 0 {
			ocs, err = s.dao.OverlordClusters(c, "sh001", arg.AppID)
		}
		if len(ocs) > 0 {
			clusters = append(clusters, ocs...)
		}
	}
	t := struct {
		Clusters []*model.OverlordToml `toml:"clusters"`
	}{}
	for _, cluster := range clusters {
		ot := &model.OverlordToml{
			Name:             cluster.Name,
			Type:             cluster.Type,
			HashMethod:       cluster.HashMethod,
			HashDistribution: cluster.HashDistribution,
			HashTag:          cluster.HashTag,
			ListenProto:      cluster.ListenProto,
			ListenAddr:       cluster.ListenAddr,
			DailTimeout:      cluster.DailTimeout,
			ReadTimeout:      cluster.ReadTimeout,
			WriteTimeout:     cluster.WriteTimeout,
			NodeConn:         cluster.NodeConn,
			PingFailLimit:    cluster.PingFailLimit,
			PingAutoEject:    cluster.PingAutoEject,
		}
		var nodes []*model.OverlordNode
		if len(cluster.Nodes) == 0 {
			if err = s.dao.DB.Where("cid=?", cluster.ID).Order("id").Find(&nodes).Error; err != nil {
				return
			}
		} else {
			nodes = cluster.Nodes
		}
		var servers []string
		for _, node := range nodes {
			var server string
			if cluster.Type == "redis_cluster" {
				server = node.Addr
			} else {
				server = fmt.Sprintf("%s:%d %s", node.Addr, node.Weight, node.Alias)
			}
			servers = append(servers, server)
		}
		ot.Servers = servers
		t.Clusters = append(t.Clusters, ot)
	}
	buf := bytes.NewBuffer(resp)
	err = toml.NewEncoder(buf).Encode(t)
	resp = buf.Bytes()
	return
}
