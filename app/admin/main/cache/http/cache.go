package http

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"

	"go-common/app/admin/main/cache/model"
	"go-common/library/ecode"
	bm "go-common/library/net/http/blademaster"

	yaml "gopkg.in/yaml.v2"
)

// @params ClustersReq
// @router get /x/admin/cache/clusters
// @response ClustersResp
func clusters(ctx *bm.Context) {
	req := new(model.ClusterReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	req.Cookie = ctx.Request.Header.Get("Cookie")
	ctx.JSON(srv.Clusters(ctx, req))
}

// @params AddClusterReq
// @router post /x/admin/cache/cluster/add
// @response EmpResp
func addCluster(ctx *bm.Context) {
	req := new(model.AddClusterReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.AddCluster(ctx, req))
}

// @params delClusterReq
// @router post /x/admin/cache/cluster/del
// @response EmpResp
func delCluster(ctx *bm.Context) {
	req := new(model.DelClusterReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.DelCluster(ctx, req))
}

// @params ClusterReq
// @router get /x/admin/cache/cluster
// @response []Cluster
func cluster(ctx *bm.Context) {
	req := new(model.ClusterReq)
	if err := ctx.Bind(req); err != nil {
		return
	}

	if req.AppID != "" {
		resp, err := srv.Cluster(ctx, req)
		if err != nil {
			ctx.JSON(nil, err)
			return
		}
		ctx.JSON(&model.ClusterResp{Clusters: resp}, nil)
	} else {
		req.Cookie = ctx.Request.Header.Get("Cookie")
		ctx.JSON(srv.Clusters(ctx, req))
	}
}

// @params ClusterDtlReq
// @router get /x/admin/cache/cluster/detail
// @response ClusterDtlResp
func clusterDtl(ctx *bm.Context) {
	req := new(model.ClusterDtlReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.ClusterDtl(ctx, req))
}

// @params ModifyClusterReq
// @router post /x/admin/cache/cluster/modify
// @response EmpResp
func modifyCluster(ctx *bm.Context) {
	req := new(model.ModifyClusterReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	ctx.JSON(srv.ModifyCluster(ctx, req))
}

func toml(ctx *bm.Context) {
	req := new(model.ClusterReq)
	if err := ctx.Bind(req); err != nil {
		return
	}
	resp, err := srv.Toml(ctx, req)
	if err != nil {
		ctx.Status(500)
		return
	}
	ctx.Writer.Write(resp)
}

// @params addFromYml
// @router post /x/admin/cache/cluster/from/yml
// @response EmpResp
func addFromYml(ctx *bm.Context) {
	req := new(model.ClusterFromYml)
	if err := ctx.Bind(req); err != nil {
		ctx.JSONMap(map[string]interface{}{
			"message": "参数有问题：app_id,zone,tw_yml",
		}, ecode.RequestErr)
		return
	}
	type server struct {
		AutoEjectHosts     bool     `yaml:"auto_eject_hosts"`
		Backlog            int      `yaml:"backlog"`
		Distribution       string   `yaml:"distribution"`
		Hash               string   `yaml:"hash"`
		Listen             string   `yaml:"listen"`
		Preconnect         bool     `yaml:"preconnect"`
		Timeout            int      `yaml:"timeout"`
		Redis              bool     `yaml:"redis"`
		ServerConnections  int      `yaml:"server_connections"`
		ServerFailureLimit int      `yaml:"server_failure_limit"`
		ServerRetryTimeout int      `yaml:"server_retry_timeout"`
		Servers            []string `yaml:"servers"`
	}
	type node struct {
		Addr   string `json:"addr"`
		Weigth int64  `json:"weight"`
		Alias  string `json:"alias"`
	}
	confs := make(map[string]server)
	err := yaml.Unmarshal([]byte(req.TwYml), &confs)
	if err != nil {
		ctx.JSONMap(map[string]interface{}{
			"message": "解析twemproxy.yml文件失败",
		}, ecode.RequestErr)
		return
	}
	mcPort := 11211
	rdPort := 26379
	for name, conf := range confs {
		ctp := "memcache"
		if conf.Redis {
			ctp = "redis"
		}
		if conf.Hash != "fnv1a_64" {
			ctx.JSONMap(map[string]interface{}{
				"message": "不支持除了fnv1a_64之外的hash方法",
			}, ecode.RequestErr)
			return
		}
		addr := "0.0.0.0:"
		_, port, err := net.SplitHostPort(conf.Listen)
		if err == nil {
			addr = addr + port
		} else {
			if conf.Redis {
				addr = addr + strconv.Itoa(rdPort)
				rdPort++
			} else {
				addr = addr + strconv.Itoa(mcPort)
				mcPort++
			}
		}
		clst := &model.AddClusterReq{
			Type:             ctp,
			AppID:            req.AppID,
			Zone:             req.Zone,
			HashMethod:       "fnv1a_64",
			HashDistribution: "ketama",
			HashTag:          "",
			Name:             name,
			DailTimeout:      1000,
			ReadTimeout:      1000,
			WriteTimeout:     1000,
			NodeConn:         2,
			PingFailLimit:    3,
			PingAutoEject:    true,
			ListenProto:      "tcp",
			ListenAddr:       addr,
		}
		if _, err := srv.AddCluster(ctx, clst); err != nil {
			ctx.JSONMap(map[string]interface{}{
				"message": "添加cluster失败:" + name + " " + err.Error(),
			}, ecode.RequestErr)
			return
		}
		var nodes []*node
		for _, n := range conf.Servers {
			ss := strings.Split(n, " ")
			idx := strings.LastIndex(ss[0], ":")
			weight, _ := strconv.ParseInt(ss[0][idx+1:], 10, 64)
			nodes = append(nodes, &node{Addr: ss[0][:idx], Weigth: weight, Alias: ss[1]})
		}
		ns, _ := json.Marshal(nodes)
		cn := &model.ModifyClusterReq{
			Name:   name,
			Action: 1,
			Nodes:  string(ns),
		}
		if _, err := srv.ModifyCluster(ctx, cn); err != nil {
			ctx.JSONMap(map[string]interface{}{
				"message": "添加cluster node失败:" + name + " " + err.Error(),
			}, ecode.RequestErr)
			return
		}
	}
}
