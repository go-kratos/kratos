package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/infra/canal/conf"
	"go-common/app/infra/canal/dao"
	"go-common/library/conf/env"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil"
	"go-common/library/stat/prom"

	"github.com/siddontang/go-mysql/client"
)

var (
	stats       = prom.New().WithState("go_canal_counter", []string{"type", "addr", "scheme", "table", "action"})
	tblReplacer = regexp.MustCompile("[0-9]+") // NOTE: replace number of sub-table name to space
)

// Canal is canal.
type Canal struct {
	dao            *dao.Dao
	instances      map[string]*Instance
	hbaseInstances map[string]*HBaseInstance
	tidbInstances  map[string]*tidbInstance
	insl           sync.RWMutex
	tidbInsl       sync.RWMutex
	err            error
	backoff        netutil.Backoff
	errCount       int64
	lastErr        time.Time
}

// NewCanal load config and start canal instance.
func NewCanal(config *conf.Config) (c *Canal) {
	c = &Canal{
		dao: dao.New(config),
	}
	cfg, err := conf.LoadCanalConf()
	if err != nil {
		panic(err)
	}
	c.instances = make(map[string]*Instance, len(cfg.Instances))
	c.hbaseInstances = make(map[string]*HBaseInstance, len(cfg.HBaseInstances))
	c.tidbInstances = make(map[string]*tidbInstance, len(cfg.TiDBInstances))
	for _, insc := range cfg.Instances {
		ins, err := NewInstance(insc)
		if err != nil {
			log.Error("new instance error(%+v)", err)
		}
		c.insl.Lock()
		c.instances[insc.Addr] = ins
		c.insl.Unlock()
		if err == nil {
			go ins.Start()
			log.Info("start canal instance(%s) success", ins.String())
		}
	}
	c.backoff = &netutil.BackoffConfig{
		MaxDelay:  120 * time.Second,
		BaseDelay: 1.0 * time.Second,
		Factor:    1.6,
		Jitter:    0.2,
	}
	c.lastErr = time.Now()
	go c.eventproc()
	// hbase
	for _, insc := range cfg.HBaseInstances {
		ins, err := NewHBaseInstance(insc)
		if err != nil {
			log.Error("new hbase instance error(%+v)", err)
		}
		c.insl.Lock()
		c.hbaseInstances[insc.Cluster] = ins
		c.insl.Unlock()
		if err == nil {
			ins.Start()
			log.Info("start hbase instance(%s) success", ins.String())
		}
	}
	go c.hbaseeventproc()

	// tidb
	for _, insc := range cfg.TiDBInstances {
		ins, err := newTiDBInstance(c, insc)
		if err != nil {
			log.Error("new instance error(%+v)", err)
			c.sendWx(fmt.Sprintf("new tidb canal instance(%s) failed error(%v)", ins.String(), err))
			continue
		}
		c.tidbInsl.Lock()
		c.tidbInstances[insc.Name] = ins
		c.tidbInsl.Unlock()
		if err == nil {
			go ins.start()
			log.Info("start tidb canal instance(%s) success", ins.String())
		}
	}
	go c.tidbEventproc()
	go c.monitorproc()
	// errproc
	go c.errproc()
	return
}

// CheckMaster check master status.
func (c *Canal) CheckMaster(addr, user, pwd string) (name string, pos int64, err error) {
	conn, err := client.Connect(addr, user, pwd, "")
	if err != nil {
		return
	}
	rr, err := conn.Execute("SHOW MASTER STATUS")
	if err != nil {
		return
	}
	name, _ = rr.GetString(0, 0)
	pos, _ = rr.GetInt(0, 1)
	if name != "" && pos > 0 {
		return
	}
	return "", 0, fmt.Errorf("check master no name|pos")
}

// PosSync sync newewst bin_pos.
func (c *Canal) PosSync(addr string) (err error) {
	c.insl.Lock()
	old, ok := c.instances[addr]
	c.insl.Unlock()
	if !ok {
		return
	}
	pos, err := old.GetMasterPos()
	if err != nil {
		return
	}
	old.Close()
	old.OnPosSynced(pos, true)

	ins, _ := NewInstance(old.config)
	c.insl.Lock()
	c.instances[addr] = ins
	c.insl.Unlock()
	ins.Start()
	return
}

// Errors returns instance errors.
func (c *Canal) Errors() (ies map[string]string) {
	ies = map[string]string{}
	c.insl.RLock()
	for _, i := range c.instances {
		ies[i.String()] = i.Error()
	}
	for _, i := range c.hbaseInstances {
		ies[i.String()] = i.Error()
	}
	c.insl.RUnlock()
	return
}

// Error returns canal error.
func (c *Canal) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

// errproc check errors.
func (c *Canal) errproc() {
	for {
		time.Sleep(10 * time.Second)
		es := c.Error()
		if es != "" {
			c.sendWx(fmt.Sprintf("canal occur error(%s)", es))
		}
		ies := c.Errors()
		for k, v := range ies {
			if v != "" {
				c.sendWx(fmt.Sprintf("canal instance(%s) occur error(%s)", k, v))
			}
		}
	}
}

// Close close canal instance
func (c *Canal) Close() {
	c.insl.RLock()
	defer c.insl.RUnlock()
	for _, ins := range c.instances {
		ins.Close()
		log.Info("close canal instance(%s) success", ins.String())
	}
	for _, ins := range c.hbaseInstances {
		ins.Close()
		log.Info("close hbase instance(%s) success", ins.String())
	}
	for _, ins := range c.tidbInstances {
		ins.close()
		log.Info("close tidb instance(%s) success", ins.String())
	}
}

func (c *Canal) eventproc() {
	ech := conf.Event()
	for {
		insc := <-ech
		if insc == nil {
			continue
		}
		ins, err := NewInstance(insc)
		if err != nil {
			log.Error("new instance error(%v)", err)
			c.sendWx(fmt.Sprintf("reload canal instance(%s) failed error(%v)", ins.String(), err))
			c.insl.Lock()
			if old, ok := c.instances[insc.Addr]; ok {
				old.Close()
			}
			c.instances[insc.Addr] = ins
			c.insl.Unlock()
			continue
		}
		c.insl.Lock()
		if old, ok := c.instances[insc.Addr]; ok {
			old.Close()
		}
		c.instances[insc.Addr] = ins
		c.insl.Unlock()
		go ins.Start()
		log.Info("reload canal instance(%s) success", ins.String())
		//	c.sendWx(fmt.Sprintf("reload canal instance(%s) success", ins.String()))
	}
}

func (c *Canal) hbaseeventproc() {
	ech := conf.HBaseEvent()
	for {
		insc := <-ech
		if insc == nil {
			continue
		}
		ins, err := NewHBaseInstance(insc)
		if err != nil {
			log.Error("new instance error(%v)", err)
			c.insl.Lock()
			if old, ok := c.hbaseInstances[insc.Cluster]; ok {
				old.Close()
			}
			c.hbaseInstances[insc.Cluster] = ins
			c.insl.Unlock()
			continue
		}
		c.insl.Lock()
		if old, ok := c.hbaseInstances[insc.Cluster]; ok {
			old.Close()
		}
		c.hbaseInstances[insc.Cluster] = ins
		c.insl.Unlock()
		ins.Start()
		log.Info("reload hbase instance(%s) success", ins.String())
	}
}

// monitorproc monitor instance delay.
func (c *Canal) monitorproc() {
	if env.DeployEnv != env.DeployEnvProd {
		return
	}
	const delay = 2 * time.Minute
	for {
		time.Sleep(delay)
		c.insl.RLock()
		for _, ins := range c.instances {
			if ins.closed {
				continue
			}
			threshold := int64(time.Duration(ins.config.MonitorPeriod) / time.Second)
			if threshold <= 0 {
				threshold = int64(delay / time.Second)
			}
			dt := ins.delay()
			if ins.config != nil && !ins.config.MonitorOff && dt > threshold {
				for _, db := range ins.config.Databases {
					c.sendWx(fmt.Sprintf("canal env(%s) 数据库(%s)地址(%s) 同步延迟时间超过阈值:%d秒 当前超过:%d秒", env.DeployEnv, db.Schema, ins.config.Addr, threshold, dt))
				}
			}
		}
		for _, ins := range c.hbaseInstances {
			if ins.closed {
				continue
			}
			threshold := int64(time.Duration(ins.config.MonitorPeriod) / time.Second)
			if threshold <= 0 {
				threshold = int64(delay / time.Second)
			}
			dt := ins.delay()
			if ins.config != nil && !ins.config.MonitorOff && dt > threshold {
				c.sendWx(fmt.Sprintf("canal env(%s) hbase集群(%s) 同步延迟时间超过阈值:%d秒 当前超过:%d秒", env.DeployEnv, ins.config.Cluster, threshold, dt))
			}
		}
		c.insl.RUnlock()
		c.tidbInsl.RLock()
		for _, ins := range c.tidbInstances {
			if ins.closed {
				continue
			}
			threshold := int64(time.Duration(ins.config.MonitorPeriod) / time.Second)
			if threshold <= 0 {
				threshold = int64(delay / time.Second)
			}
			dt := ins.delay()
			if ins.config != nil && dt > threshold {
				for _, db := range ins.config.Databases {
					c.sendWx(fmt.Sprintf("tidb canal env(%s) 数据库(%s)地址(%s) 同步延迟时间超过阈值:%d秒 当前超过:%d秒", env.DeployEnv, db.Schema, ins.config.Name, threshold, dt))
				}
			}
		}
		c.tidbInsl.RUnlock()
	}
}

func (c *Canal) sendWx(msg string) {
	count := atomic.LoadInt64(&c.errCount)
	atomic.AddInt64(&c.errCount, 1)
	duration := c.backoff.Backoff(int(count))
	if time.Since(c.lastErr) < duration {
		return
	}
	c.lastErr = time.Now()
	sendWechat(msg, conf.Conf.Monitor.User)
}

func sendWechat(msg string, user string) {
	params := url.Values{}
	params.Set("Action", "CreateWechatMessage")
	params.Set("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0etracehubtuhui@bilibili.com")
	params.Set("UserName", user)
	params.Set("Content", msg)
	params.Set("TreeId", "bilibili.main.common-arch.canal")
	params.Set("Title", "canal 监控报警")
	params.Set("Signature", "1")
	req, _ := http.NewRequest("POST", "http://merak.bilibili.co", strings.NewReader(params.Encode()))
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	var v struct {
		ReqID    string `json:"reqId"`
		Status   int64  `json:"status"`
		Response struct {
			status int
		} `json:"Response"`
	}
	if err := xhttp.NewClient(conf.Conf.HTTPClient).Do(context.TODO(), req, &v); err != nil {
		log.Error("send wechat monitor status(%d) msg(%v,%v) error(%v)", v.Status, v.Response, v.Response.status, err)
	}
}
