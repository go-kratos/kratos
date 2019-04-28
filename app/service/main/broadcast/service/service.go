package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/service/main/broadcast/dao"
	"go-common/app/service/main/broadcast/model"
	identify "go-common/app/service/main/identify/api/grpc"
	location "go-common/app/service/main/location/api"
	"go-common/library/cache"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/netutil"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/balancer/wrr"

	"google.golang.org/grpc"
)

const (
	_onlineTick     = time.Second * 10
	_onlineDeadline = time.Minute * 5
)

var (
	_olBackoff = &netutil.BackoffConfig{
		MaxDelay:  5 * time.Second,
		BaseDelay: 1.0 * time.Second,
		Factor:    1.6,
		Jitter:    0.2,
	}
)

// Service struct
type Service struct {
	c     *paladin.Map
	dao   *dao.Dao
	cache *cache.Cache
	dis   *discovery.Discovery
	// online
	totalIPs   int64
	totalConns int64
	roomCount  map[string]int32
	// client
	identifyCli  identify.IdentifyClient
	locationCli  location.LocationClient
	servers      []*model.ServerInfo
	loadBalancer *LoadBalancer
	regions      map[string]string // province -> region
	// infoc
	stats      *infoc.Infoc
	srvBackoff *ServerBackoff
	srvConf    *ServerConfig
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// New new a service and return.
func New(dis *discovery.Discovery) (s *Service) {
	var (
		appConf    = new(paladin.TOML)
		srvConf    = new(ServerConfig)
		srvBackoff = new(ServerBackoff)
		regions    map[string][]string
		infocConf  struct {
			Stats *infoc.Config
		}
		grpcConf struct {
			Client *warden.ClientConfig
		}
	)
	checkErr(paladin.Watch("application.toml", appConf))
	checkErr(appConf.Get("server").UnmarshalTOML(srvConf))
	checkErr(appConf.Get("regions").UnmarshalTOML(&regions))
	checkErr(appConf.Get("backoff").UnmarshalTOML(srvBackoff))
	checkErr(paladin.Get("infoc.toml").UnmarshalTOML(&infocConf))
	if err := paladin.Get("grpc.toml").UnmarshalTOML(&grpcConf); err != nil {
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}
	idtv1, err := identify.NewClient(grpcConf.Client, grpc.WithBalancerName(wrr.Name))
	if err != nil {
		panic(err)
	}
	locConn, err := warden.NewClient(grpcConf.Client).Dial(context.Background(), "discovery://default/location.service")
	if err != nil {
		panic(err)
	}
	s = &Service{
		c:            appConf,
		srvConf:      srvConf,
		srvBackoff:   srvBackoff,
		dis:          dis,
		dao:          dao.New(),
		cache:        cache.New(64, 10240),
		stats:        infoc.New(infocConf.Stats),
		identifyCli:  idtv1,
		locationCli:  location.NewLocationClient(locConn),
		loadBalancer: NewLoadBalancer(),
		regions:      make(map[string]string),
	}
	s.initRegions(regions)
	s.initServer()
	s.loadOnline()
	go s.onlineproc()
	return s
}

func (s *Service) initRegions(regions map[string][]string) {
	for region, ps := range regions {
		for _, province := range ps {
			s.regions[province] = region
		}
	}
}

func (s *Service) initServer() {
	res := s.dis.Build("push.interface.broadcast")
	event := res.Watch()
	select {
	case _, ok := <-event:
		if ok {
			s.newServers(res)
		} else {
			panic("discovery watch failed")
		}
	case <-time.After(10 * time.Second):
		log.Error("discovery start timeout")
	}
	go func() {
		for {
			if _, ok := <-event; !ok {
				return
			}
			s.newServers(res)
		}
	}()
}

func (s *Service) newServers(res naming.Resolver) {
	if zoneIns, ok := res.Fetch(context.Background()); ok {
		var (
			totalConns int64
			totalIPs   int64
			srvs       []*model.ServerInfo
		)
		for _, zins := range zoneIns {
			for _, in := range zins {
				if in.Metadata == nil {
					log.Error("instance metadata is empty(%+v)", in)
					continue
				}
				if ok, _ := strconv.ParseBool(in.Metadata["offline"]); ok {
					continue
				}
				conns, err := strconv.ParseInt(in.Metadata["conns"], 10, 32)
				if err != nil {
					log.Error("strconv.ParseInt(conns:%d) error(%v)", conns, err)
					continue
				}
				ips, err := strconv.ParseInt(in.Metadata["ips"], 10, 32)
				if err != nil {
					log.Error("strconv.ParseInt(ips:%d) error(%v)", ips, err)
					continue
				}
				roomIPs, err := strconv.ParseInt(in.Metadata["room_ips"], 10, 32)
				if err != nil {
					log.Error("strconv.ParseInt(room_ips:%d) error(%v)", ips, err)
					continue
				}
				weight, err := strconv.ParseInt(in.Metadata["weight"], 10, 32)
				if err != nil {
					log.Error("strconv.ParseInt(weight:%d) error(%v)", conns, err)
					continue
				}
				var (
					addrs    []string
					addrsV6  []string
					lat, lng float64
				)
				if in.Metadata["ip_addrs"] != "" {
					addrs = strings.Split(in.Metadata["ip_addrs"], ",")
				}
				if in.Metadata["ip_addrs_v6"] != "" {
					addrsV6 = strings.Split(in.Metadata["ip_addrs_v6"], ",")
				}
				for _, addr := range addrs {
					resp, err := s.locationCli.Info(context.Background(), &location.InfoReq{Addr: addr})
					if err == nil {
						lat, lng = resp.Latitude, resp.Longitude
						break
					}
				}
				overseas, _ := strconv.ParseBool(in.Metadata["overseas"])
				srvs = append(srvs, &model.ServerInfo{
					Region:      in.Region,
					Server:      in.Hostname,
					ConnCount:   int32(conns),
					IPCount:     int32(ips),
					RoomIPCount: int32(roomIPs),
					Weight:      int32(weight),
					IPAddrs:     addrs,
					IPAddrsV6:   addrsV6,
					Updated:     in.LastTs,
					Latitude:    lat,
					Longitude:   lng,
					Overseas:    overseas,
				})
				totalConns += conns
				totalIPs += ips
			}
		}
		s.totalConns = totalConns
		s.totalIPs = totalIPs
		s.servers = srvs
		s.loadBalancer.Update(srvs)
		if ok, _ := s.c.Get("saveServers").Bool(); ok {
			s.dao.SetServers(context.Background(), srvs)
			log.Info("set servers:%d", len(srvs))
		}
	}
}

// Ping Service.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) onlineproc() {
	var retry int
	for {
		if err := s.loadOnline(); err != nil {
			retry++
			time.Sleep(_olBackoff.Backoff(retry))
			continue
		}
		retry = 0
		time.Sleep(_onlineTick)
	}
}

func (s *Service) loadOnline() (err error) {
	const pullRetries = 3
	var roomCount = make(map[string]int32)
	srvs, err := s.dao.Servers(context.Background())
	if err != nil {
		log.Error("s.dao.Servers() error(%v)", err)
		return
	}
	for _, server := range srvs {
		for i := 0; i < roomsShard; i++ {
			var online *model.Online
			for r := 0; r < pullRetries; r++ {
				if online, err = s.dao.ServerOnline(context.Background(), server.Server, i); err != nil {
					log.Error("s.dao.ServerOnline(%s, %d) retries:%d error(%v)", server.Server, i, r, err)
					time.Sleep(_olBackoff.Backoff(r))
					continue
				}
				break
			}
			if err != nil {
				return
			}
			if online == nil {
				continue
			}
			if time.Since(time.Unix(online.Updated, 0)) > _onlineDeadline {
				s.dao.DelServerOnline(context.Background(), server.Server, i)
				continue
			}
			for roomID, count := range online.RoomCount {
				roomCount[roomID] += count
			}
		}
	}
	// FIXME migrate rooms count
	for i := 0; i < roomsShard; i++ {
		var rooms map[string]int32
		for r := 0; r < pullRetries; r++ {
			if rooms, err = s.dao.MigrateRooms(context.Background(), i); err != nil {
				log.Error("s.dao.MigrateRooms(%d) retries:%d error(%v)", i, r, err)
				time.Sleep(_olBackoff.Backoff(r))
				continue
			}
			break
		}
		if err != nil {
			return
		}
		for rid, cnt := range rooms {
			roomCount[rid] += cnt
		}
	}
	s.roomCount = roomCount
	log.Info("loadOnline rooms:%d", len(roomCount))
	return
}
