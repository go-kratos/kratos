package service

import (
	"context"
	"time"

	pb "go-common/app/service/main/broadcast/api/grpc/v1"
	"go-common/app/service/main/broadcast/model"
	location "go-common/app/service/main/location/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
)

// ServerBackoff server backoff config.
type ServerBackoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float32
	Jitter    float32
}

// ServerConfig server config.
type ServerConfig struct {
	Domain       string
	HostDomain   string
	TCPPort      int
	WSPort       int
	WSSPort      int
	HeartbeatMax int
	Heartbeat    xtime.Duration
	RegionWeight float64
}

func isIPV6(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return false
		case ':':
			return true
		}
	}
	return false
}

// ServerInfos get servers info.
func (s *Service) ServerInfos(c context.Context) (res []*model.ServerInfo, err error) {
	return s.servers, nil
}

// ServerList get server list.
func (s *Service) ServerList(c context.Context, platform string) *pb.ServerListReply {
	reply := &pb.ServerListReply{
		Domain:       s.srvConf.Domain,
		TcpPort:      int32(s.srvConf.TCPPort),
		WsPort:       int32(s.srvConf.WSPort),
		WssPort:      int32(s.srvConf.WSSPort),
		Heartbeat:    int32(time.Duration(s.srvConf.Heartbeat) / time.Second),
		HeartbeatMax: int32(s.srvConf.HeartbeatMax),
		Backoff: &pb.Backoff{
			MaxDelay:  s.srvBackoff.MaxDelay,
			BaseDelay: s.srvBackoff.BaseDelay,
			Factor:    s.srvBackoff.Factor,
			Jitter:    s.srvBackoff.Jitter,
		},
	}
	domains, addrs := s.nodeAddrs(c)
	if platform == "web" {
		reply.Nodes = domains
		// FIXME cross domain
		reply.Nodes = []string{s.srvConf.Domain}
	} else {
		reply.Nodes = addrs
	}
	if len(reply.Nodes) == 0 {
		reply.Nodes = []string{s.srvConf.Domain}
	}
	return reply
}

// ServerWeight server node details.
func (s *Service) ServerWeight(c context.Context, clientIP string) (interface{}, string, string, error) {
	var (
		region   string
		province string
	)
	if clientIP != "" {
		req := &location.InfoReq{Addr: clientIP}
		resp, err := s.locationCli.Info(c, req)
		if err != nil {
			return nil, "", "", err
		}
		province = resp.Province
		region = s.regions[province]
	}
	return s.loadBalancer.NodeDetails(region, s.srvConf.RegionWeight), region, province, nil
}

func (s *Service) nodeAddrs(c context.Context) (domains, addrs []string) {
	var (
		country  string
		region   string
		province string
		lat, lng float64
		overseas bool
		clientIP = metadata.String(c, metadata.RemoteIP)
	)
	req := &location.InfoReq{Addr: clientIP}
	resp, err := s.locationCli.Info(c, req)
	if err == nil {
		region = s.regions[resp.Province]
		country = resp.Country
		province = resp.Province
		lat, lng = resp.Latitude, resp.Longitude
		// 海外节点判断
		if (country != "" && country != "中国" && country != "局域网") || province == "香港" || province == "澳门" || province == "台湾" {
			overseas = true
		}
	}
	if region != "" || lat == 0 || lng == 0 {
		domains, addrs = s.loadBalancer.NodeAddrsByRegion(region, s.srvConf.HostDomain, s.srvConf.RegionWeight, overseas, isIPV6(clientIP))
	} else {
		domains, addrs = s.loadBalancer.NodeAddrsByGeo(s.srvConf.HostDomain, lat, lng, overseas, isIPV6(clientIP))
	}
	log.Info("nodeAddrs clientIP:%s overseas:%t region:%s country:%s province:%s lat:%6f lng:%6f domains:%v addrs:%v", clientIP, overseas, region, country, province, lat, lng, domains, addrs)
	return
}
