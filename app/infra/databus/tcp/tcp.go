package tcp

import (
	"errors"
	lg "log"
	"net"
	"os"
	"sync"
	"time"

	"go-common/app/infra/databus/conf"
	"go-common/app/infra/databus/dsn"
	"go-common/app/infra/databus/model"
	"go-common/app/infra/databus/service"
	"go-common/library/log"

	"github.com/Shopify/sarama"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	// redis proto
	_protoStr   = '+'
	_protoErr   = '-'
	_protoInt   = ':'
	_protoBulk  = '$'
	_protoArray = '*'

	// redis cmd
	_ping = "ping"
	_auth = "auth"
	_quit = "quit"
	_set  = "set"
	_hset = "hset"
	_mget = "mget"
	_ok   = "OK"
	_pong = "PONG"

	// client role
	_rolePub = "pub"
	_roleSub = "sub"

	_listenDelay = 5 * time.Millisecond // how long to sleep on accept failure
	_clearDelay  = 30 * time.Second

	_batchNum      = 100                    // batch write message length
	_batchInterval = 100 * time.Millisecond // batch write interval
	_batchTimeout  = 30 * time.Second       // return empty if timeout

	// connection timeout
	_readTimeout    = 5 * time.Second
	_writeTimeout   = 5 * time.Second
	_pubReadTimeout = 20 * time.Minute
	_subReadTimeout = _batchTimeout + 10*time.Second

	// conn read buffer size 64K
	_readBufSize = 1024 * 64
	// conn write buffer size 8K
	_writeBufSize = 1024 * 8
	// conn max value size(kafka 1M)
	_maxValueSize = 1000000

	// prom operation
	_opAddConsumerRequest = "request_add_comsuner"
	_opCurrentConsumer    = "current_consumer"
	_opAddProducerRequest = "request_add_producer"
	_opAuthError          = "auth_error"
	_opProducerMsgSpeed   = "producer_msg_speed"
	_opConsumerMsgSpeed   = "consumer_msg_speed"
	_opConsumerPartition  = "consumer_partition_speed"
	_opPartitionOffset    = "consumer_partition_offset"
)

var (
	_nullBulk = []byte("-1")
	// kafka header
	_headerColor    = []byte("color")
	_headerMetadata = []byte("metadata")
	// encode type pb/json
	_encodePB = []byte("pb")
)

var (
	errCmdAuthFailed        = errors.New("auth failed")
	errAuthInfo             = errors.New("auth info error")
	errPubParams            = errors.New("pub params error")
	errCmdNotSupport        = errors.New("command not support")
	errClusterNotExist      = errors.New("cluster not exist")
	errClusterNotSupport    = errors.New("cluster not support")
	errConnClosedByServer   = errors.New("connection closed by databus")
	errConnClosedByClient   = errors.New("connection closed by client")
	errClosedMsgChannel     = errors.New("message channel is closed")
	errClosedNotifyChannel  = errors.New("notification channel is closed")
	errNoPubPermission      = errors.New("no publish permission")
	errNoSubPermission      = errors.New("no subscribe permission")
	errConsumerClosed       = errors.New("kafka consumer closed")
	errCommitParams         = errors.New("commit offset params error")
	errMsgFormat            = errors.New("message format must be json")
	errConsumerOver         = errors.New("too many consumers")
	errConsumerTimeout      = errors.New("consumer initial timeout")
	errConnRead             = errors.New("connection read error")
	errUseLessConsumer      = errors.New("useless consumer")
	errKafKaData            = errors.New("err kafka data maybe rebalancing")
	errCousmerCreateLimiter = errors.New("err consumer create limiter")
)

var (
	// tcp listener
	listener net.Listener
	quit     = make(chan struct{})
	// producer snapshot, key:group+topic
	producers = make(map[string]sarama.SyncProducer)
	pLock     sync.RWMutex
	// Pubs
	pubs    = make(map[*Pub]struct{})
	pubLock sync.RWMutex
	// Subs
	subs    = make(map[*Sub]struct{})
	subLock sync.RWMutex
	// service for auth
	svc *service.Service
	// limiter
	consumerLimter = make(chan struct{}, 100)
)

// Init init service
func Init(c *conf.Config, s *service.Service) {
	var err error
	if listener, err = net.Listen("tcp", c.Addr); err != nil {
		panic(err)
	}
	// sarama should be initialized otherwise errors will be ignored when sarama catch error
	sarama.Logger = lg.New(os.Stdout, "[Sarama] ", lg.LstdFlags)
	// sarama metrics disable
	metrics.UseNilMetrics = true
	svc = s
	log.Info("start tcp listen addr: %s", c.Addr)
	go accept()
	go clear()
	go clusterproc()
}

func newProducer(group, topic string, pCfg *conf.Kafka) (p sarama.SyncProducer, err error) {
	var (
		ok  bool
		key = key(pCfg.Cluster, group, topic)
	)
	pLock.RLock()
	if p, ok = producers[key]; ok {
		pLock.RUnlock()
		return
	}
	pLock.RUnlock()
	// new
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Version = sarama.V1_0_0_0
	if p, err = sarama.NewSyncProducer(pCfg.Brokers, conf); err != nil {
		log.Error("group(%s) topic(%s) cluster(%s) NewSyncProducer error(%v)", group, topic, pCfg.Cluster, err)
		return
	}
	pLock.Lock()
	producers[key] = p
	pLock.Unlock()
	return
}

// Close close all producers and consumers
func Close() {
	close(quit)
	if listener != nil {
		listener.Close()
	}
	// close all consumers
	subLock.RLock()
	for sub := range subs {
		sub.Close()
	}
	subLock.RUnlock()
	pubLock.RLock()
	for pub := range pubs {
		pub.Close(true)
	}
	pubLock.RUnlock()
	pLock.RLock()
	for _, p := range producers {
		p.Close()
	}
	pLock.RUnlock()
}

func accept() {
	var (
		err  error
		ok   bool
		netC net.Conn
		netE net.Error
	)
	for {
		if netC, err = listener.Accept(); err != nil {
			if netE, ok = err.(net.Error); ok && netE.Temporary() {
				log.Error("tcp: Accept error: %v; retrying in %v", err, _listenDelay)
				time.Sleep(_listenDelay)
				continue
			}
			return
		}
		select {
		case <-quit:
			netC.Close()
			return
		default:
		}
		go serveConn(netC)
	}
}

// serveConn serve tcp connect.
func serveConn(nc net.Conn) {
	var (
		err   error
		p     *Pub
		s     *Sub
		d     *dsn.DSN
		cfg   *conf.Kafka
		batch int64
		addr  = nc.RemoteAddr().String()
	)
	c := newConn(nc, _readTimeout, _writeTimeout)
	if d, cfg, batch, err = auth(c); err != nil {
		log.Error("auth failed addr(%s) error(%v)", addr, err)
		c.WriteError(err)
		return
	}
	// auth succeed
	if err = c.Write(proto{prefix: _protoStr, message: _ok}); err != nil {
		log.Error("c.Write() error(%v)", err)
		c.Close()
		return
	}
	if err = c.Flush(); err != nil {
		c.Close()
		return
	}
	log.Info("auth succeed group(%s) topic(%s) color(%s) cluster(%s) addr(%s) role(%s)", d.Group, d.Topic, d.Color, cfg.Cluster, addr, d.Role)
	// command
	switch d.Role {
	case _rolePub: // producer
		svc.CountProm.Incr(_opAddProducerRequest, d.Group, d.Topic)
		if p, err = NewPub(c, d.Group, d.Topic, d.Color, cfg); err != nil {
			c.WriteError(err)
			log.Error("group(%s) topic(%s) color(%s) cluster(%s) addr(%s) NewPub error(%v)", d.Group, d.Topic, d.Color, cfg.Cluster, addr, err)
			return
		}
		pubLock.Lock()
		pubs[p] = struct{}{}
		pubLock.Unlock()
		p.Serve()
	case _roleSub: // consumer
		svc.CountProm.Incr(_opAddConsumerRequest, d.Group, d.Topic)
		select {
		case consumerLimter <- struct{}{}:
		default:
			err = errCousmerCreateLimiter
			c.WriteError(err)
			log.Error("group(%s) topic(%s) color(%s) cluster(%s) addr(%s) error(%v)", d.Group, d.Topic, d.Color, cfg.Cluster, addr, err)
			return
		}
		if s, err = NewSub(c, d.Group, d.Topic, d.Color, cfg, batch); err != nil {
			c.WriteError(err)
			log.Error("group(%s) topic(%s) color(%s) cluster(%s) addr(%s) NewSub error(%v)", d.Group, d.Topic, d.Color, cfg.Cluster, addr, err)
			return
		}
		subLock.Lock()
		subs[s] = struct{}{}
		subLock.Unlock()
		s.Serve()
		svc.CountProm.Incr(_opCurrentConsumer, d.Group, d.Topic)
	default:
		// other command will not be, auth check that.
	}
}

func auth(c *conn) (d *dsn.DSN, cfg *conf.Kafka, batch int64, err error) {
	var (
		args [][]byte
		cmd  string
		addr = c.conn.RemoteAddr().String()
	)
	if cmd, args, err = c.Read(); err != nil {
		log.Error("c.Read addr(%s) error(%v)", addr, err)
		return
	}
	if cmd != _auth || len(args) != 1 {
		log.Error("c.Read addr(%s) first cmd(%s) not auth or have not enough args(%v)", addr, cmd, args)
		err = errCmdAuthFailed
		return
	}
	// key:secret@group/topic=?&role=?&offset=?
	if d, err = dsn.ParseDSN(string(args[0])); err != nil {
		log.Error("auth failed arg(%s) is illegal,addr(%s) error(%v)", args[0], addr, err)
		return
	}

	cfg, batch, err = Auth(d, addr)
	return
}

// Auth 校验认证信息并反回相应配置
// 与 http 接口共用，不要在此方法执行 io 操作
func Auth(d *dsn.DSN, addr string) (cfg *conf.Kafka, batch int64, err error) {
	var (
		a  *model.Auth
		ok bool
	)

	if a, ok = svc.AuthApp(d.Group); !ok {
		log.Error("addr(%s) group(%s) cant not be found", addr, d.Group)
		svc.CountProm.Incr(_opAuthError, d.Group, d.Topic)
		err = errAuthInfo
		return
	}
	batch = a.Batch
	if err = a.Auth(d.Group, d.Topic, d.Key, d.Secret); err != nil {
		log.Error("a.Auth addr(%s) group(%s) topic(%s) color(%s) key(%s) secret(%s) error(%v)", addr, d.Group, d.Topic, d.Color, d.Key, d.Secret, err)
		svc.CountProm.Incr(_opAuthError, d.Group, d.Topic)
		return
	}
	switch d.Role {
	case _rolePub:
		if !a.CanPub() {
			err = errNoPubPermission
			return
		}
	case _roleSub:
		if !a.CanSub() {
			err = errNoSubPermission
			return
		}
	default:
		err = errCmdNotSupport
		return
	}
	if len(conf.Conf.Clusters) == 0 {
		err = errClusterNotExist
		return
	}
	if cfg, ok = conf.Conf.Clusters[a.Cluster]; !ok || cfg == nil {
		log.Error("a.Auth addr(%s) group(%s) topic(%s) color(%s) key(%s) secret(%s) cluster(%s) not support", addr, d.Group, d.Topic, d.Color, d.Key, d.Secret, a.Cluster)
		err = errClusterNotSupport
	}
	// TODO check ip addr
	// rAddr = conn.RemoteAddr().String()
	return
}

// ConsumerAddrs returns consumer addrs.
func ConsumerAddrs(group string) (addrs []string, err error) {
	subLock.RLock()
	for sub := range subs {
		if sub.group == group {
			addrs = append(addrs, sub.addr)
		}
	}
	subLock.RUnlock()
	return
}

func key(cluster, group, topic string) string {
	return cluster + ":" + group + ":" + topic
}

func clear() {
	for {
		time.Sleep(_clearDelay)
		t := time.Now()
		log.Info("clear proc start,id(%d)", t.Nanosecond())
		subLock.Lock()
		for sub := range subs {
			if sub.Closed() {
				delete(subs, sub)
			}
		}
		subLock.Unlock()
		pubLock.Lock()
		for pub := range pubs {
			if pub.Closed() {
				delete(pubs, pub)
			}
		}
		pubLock.Unlock()
		log.Info("clear proc end,id(%d) used(%d)", t.Nanosecond(), time.Since(t))
	}
}

func clusterproc() {
	for {
		oldAuth, ok := <-svc.ClusterEvent()
		if !ok {
			return
		}
		log.Info("cluster changed event group(%s) topic(%s) cluster(%s)", oldAuth.Group, oldAuth.Topic, oldAuth.Cluster)

		k := key(oldAuth.Cluster, oldAuth.Group, oldAuth.Topic)
		pLock.Lock()
		if p, ok := producers[k]; ok {
			// renew producer
			if newAuth, ok := svc.AuthApp(oldAuth.Group); ok {
				pLock.Unlock()
				np, err := newProducer(newAuth.Group, newAuth.Topic, conf.Conf.Clusters[newAuth.Cluster])
				pLock.Lock()
				// check pubs
				pubLock.Lock()
				for pub := range pubs {
					if pub.group == oldAuth.Group && pub.topic == oldAuth.Topic {
						if err != nil {
							pub.Close(true)
						} else {
							pub.producer = np
						}
					}
				}
				pubLock.Unlock()
			}
			// close unused producer
			p.Close()
			delete(producers, k)
		}
		pLock.Unlock()
		// wait closing subs
		subLock.Lock()
		for sub := range subs {
			if sub.group == oldAuth.Group && sub.topic == oldAuth.Topic {
				sub.WaitClosing()
			}
		}
		subLock.Unlock()
	}
}
