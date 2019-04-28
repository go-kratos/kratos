package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/infra/notify/conf"
	"go-common/app/infra/notify/dao"
	"go-common/app/infra/notify/model"
	"go-common/library/log"
	"go-common/library/net/netutil"
	"go-common/library/stat/prom"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/rcrowley/go-metrics"
)

func init() {
	// sarama metrics disable
	metrics.UseNilMetrics = true
}

var (
	errClusterNotSupport   = errors.New("cluster not support")
	errClosedNotifyChannel = errors.New("notification channel is closed")
	errConsumerOver        = errors.New("too many consumers")
	errCallbackParse       = errors.New("parse callback error")

	statProm  = prom.New().WithState("go_notify_state", []string{"role", "group", "topic", "partition"})
	countProm = prom.New().WithState("go_notify_counter", []string{"operation", "group", "topic"})
	// prom operation
	_opCurrentConsumer   = "current_consumer"
	_opProducerMsgSpeed  = "producer_msg_speed"
	_opConsumerMsgSpeed  = "consumer_msg_speed"
	_opConsumerPartition = "consumer_partition_speed"
	_opPartitionOffset   = "consumer_partition_offset"
	_opConsumerFail      = "consumer_fail"
)

const (
	_defRoutine = 1
	_retry      = 5
	_syncCall   = 1
	_asyncCall  = 2
)

// Sub notify instance
type Sub struct {
	c        *conf.Config
	ctx      context.Context
	cancel   context.CancelFunc
	w        *model.Watcher
	cluster  *conf.Kafka
	clients  *Clients
	dao      *dao.Dao
	consumer *cluster.Consumer
	filter   func(msg []byte) bool
	routine  int
	backoff  netutil.BackoffConfig
	asyncRty chan *rtyMsg
	ticker   *time.Ticker
	stop     bool
	closed   bool
	once     sync.Once
}

type rtyMsg struct {
	id    int64
	msg   string
	index int
}

// NewSub create notify instance and return it.
func NewSub(w *model.Watcher, d *dao.Dao, c *conf.Config) (n *Sub, err error) {
	n = &Sub{
		c:        c,
		w:        w,
		routine:  _defRoutine,
		backoff:  netutil.DefaultBackoffConfig,
		asyncRty: make(chan *rtyMsg, 100),
		dao:      d,
		ticker:   time.NewTicker(time.Minute),
	}
	n.ctx, n.cancel = context.WithCancel(context.Background())
	if clu, ok := c.Clusters[w.Cluster]; ok {
		n.cluster = clu
	} else {
		err = errClusterNotSupport
		return
	}
	if len(w.Filters) != 0 {
		n.parseFilter()
	}
	err = n.parseCallback()
	if err != nil {
		err = errCallbackParse
		return
	}
	// init clients
	n.clients = NewClients(c, w)
	err = n.dial()
	if err != nil {
		return
	}
	if w.Concurrent != 0 {
		n.routine = w.Concurrent
	}
	go n.asyncRtyproc()
	for i := 0; i < n.routine; i++ {
		go n.serve()
	}
	countProm.Incr(_opCurrentConsumer, w.Group, w.Topic)
	return
}

func (n *Sub) parseFilter() {
	n.filter = func(b []byte) bool {
		nmsg := new(model.Message)
		err := json.Unmarshal(b, nmsg)
		if err != nil {
			log.Error("json err %v", err)
			return true
		}
		for _, f := range n.w.Filters {
			switch {
			case f.Field == "table":
				if f.Condition == model.ConditionEq && f.Value == nmsg.Table {
					return false
				}
				if f.Condition == model.ConditionPre && strings.HasPrefix(nmsg.Table, f.Value) {
					return false
				}
			case f.Field == "action":
				if f.Condition == model.ConditionEq && f.Value == nmsg.Action {
					return false
				}
			}
		}
		return true
	}
}

// parseCallback parse each watcher's callback urls
func (n *Sub) parseCallback() (err error) {
	var notifyURL *model.NotifyURL
	cbm := make(map[string]int8)
	log.Info("callback(%v), topic(%s), group(%s)", n.w.Callback, n.w.Topic, n.w.Group)
	err = json.Unmarshal([]byte(n.w.Callback), &cbm)
	if err != nil {
		log.Error(" Notify.parseCallback sub parse callback err %v, topic(%s), group(%s), callback(%s)",
			err, n.w.Topic, n.w.Group, n.w.Callback)
		return
	}
	cbs := make([]*model.Callback, 0, len(cbm))
	for u, p := range cbm {
		notifyURL, err = parseNotifyURL(u)
		if err != nil {
			log.Error("Notify.parseCallback url parse error(%v), url(%s), topic(%s), group(%s)",
				err, u, n.w.Topic, n.w.Group)
			return
		}
		cbs = append(cbs, &model.Callback{URL: notifyURL, Priority: p})
	}
	sort.Slice(cbs, func(i, j int) bool { return cbs[i].Priority > cbs[j].Priority })
	n.w.Callbacks = cbs
	return
}

func (n *Sub) dial() (err error) {
	cfg := cluster.NewConfig()
	cfg.ClientID = fmt.Sprintf("%s-%s", n.w.Topic, n.w.Group)
	cfg.Net.KeepAlive = time.Second
	cfg.Consumer.Offsets.CommitInterval = time.Second
	cfg.Consumer.MaxWaitTime = time.Millisecond * 250
	cfg.Consumer.MaxProcessingTime = time.Millisecond * 50
	cfg.Consumer.Return.Errors = true
	cfg.Version = sarama.V1_0_0_0
	cfg.Group.Return.Notifications = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	if n.consumer, err = cluster.NewConsumer(n.cluster.Brokers, n.w.Group, []string{n.w.Topic}, cfg); err != nil {
		log.Error("group(%s) topic(%s) cluster(%s)  cluster.NewConsumer() error(%v)", n.w.Group, n.w.Topic, n.cluster.Cluster, err)
	} else {
		log.Info("group(%s) topic(%s) cluster(%s) cluster.NewConsumer() ok", n.w.Group, n.w.Topic, n.cluster.Cluster)
	}
	return
}

func (n *Sub) serve() {
	var (
		msg    *sarama.ConsumerMessage
		err    error
		ok     bool
		notify *cluster.Notification
	)
	defer n.once.Do(func() {
		n.cancel()
		n.Close()
	})
	for {
		select {
		case <-n.ctx.Done():
			log.Error("sub cancel")
			return
		case err = <-n.consumer.Errors():
			log.Error("group(%s) topic(%s) cluster(%s)  catch error(%v)", n.w.Group, n.w.Topic, n.cluster.Cluster, err)
			return
		case notify, ok = <-n.consumer.Notifications():
			if !ok {
				err = errClosedNotifyChannel
				log.Info("notification notOk group(%s) topic(%s) cluster(%s) catch error(%v)", n.w.Group, n.w.Topic, n.cluster.Cluster, err)
				return
			}
			switch notify.Type {
			case cluster.UnknownNotification, cluster.RebalanceError:
				err = errClosedNotifyChannel
				log.Error("notification(%s) group(%s) topic(%s) cluster(%s)  catch error(%v)", notify.Type, n.w.Group, n.w.Topic, n.cluster.Cluster, err)
				return
			case cluster.RebalanceStart:
				log.Info("notification(%s) group(%s) topic(%s) cluster(%s)  catch error(%v)", notify.Type, n.w.Group, n.w.Topic, n.cluster.Cluster, err)
				continue
			case cluster.RebalanceOK:
				log.Info("notification(%s) group(%s) topic(%s) cluster(%s) catch error(%v)", notify.Type, n.w.Group, n.w.Topic, n.cluster.Cluster, err)
			}
			if len(notify.Current[n.w.Topic]) == 0 {
				err = errConsumerOver
				log.Warn("notification(%s) no topic group(%s) topic(%s) cluster(%s)  catch error(%v)", notify.Type, n.w.Group, n.w.Topic, n.cluster.Cluster, err)
				return
			}
		case msg, ok = <-n.consumer.Messages():
			if !ok {
				log.Error("group(%s) topic(%s) cluster(%s)  message channel closed", n.w.Group, n.w.Topic, n.cluster.Cluster)
				return
			}
			n.push(msg.Value)
			n.consumer.MarkPartitionOffset(msg.Topic, msg.Partition, msg.Offset, "")
			statProm.State(_opPartitionOffset, msg.Offset, n.w.Group, n.w.Topic, strconv.Itoa(int(msg.Partition)))
			countProm.Incr(_opConsumerMsgSpeed, n.w.Group, n.w.Topic)
			statProm.Incr(_opConsumerPartition, n.w.Group, n.w.Topic, strconv.Itoa(int(msg.Partition)))
		}
	}
}

// push call retry for each consumer group callback
// will push to retry channel if failed.
func (n *Sub) push(nmsg []byte) {
	if n.filter != nil && n.filter(nmsg) {
		return
	}
	for n.stop {
		time.Sleep(time.Minute)
	}
	msg := string(nmsg)
	for i := 0; i < len(n.w.Callbacks); i++ {
		cb := n.w.Callbacks[i]
		if err := n.retry(cb.URL, string(nmsg), _syncCall); err != nil {
			id, err := n.backupMsg(msg, i)
			if err != nil {
				log.Error("group(%s) topic(%s) add msg(%s) backup fail err %v", n.w.Group, n.w.Topic, string(nmsg), err)
			}
			n.addAsyncRty(id, msg, i)
			return
		}
	}
}

// asyncRtyproc async retry proc
func (n *Sub) asyncRtyproc() {
	var err error
	for {
		if n.Closed() {
			return
		}
		rty, ok := <-n.asyncRty
		countProm.Decr(_opConsumerFail, n.w.Group, n.w.Topic)
		if !ok {
			log.Error("async chan close ")
			return
		}
		for i := rty.index; i < len(n.w.Callbacks); i++ {
			err = n.retry(n.w.Callbacks[i].URL, rty.msg, _asyncCall)
			if err != nil {
				n.addAsyncRty(rty.id, rty.msg, i)
				break
			}
		}
		if err == nil {
			// if ok,restart consumer.
			n.stop = false
			n.delBackup(rty.id)
		}
	}
}

// retry Sub do callback with retry
func (n *Sub) retry(uri *model.NotifyURL, msg string, source int) (err error) {
	log.Info("Notify.retry do callback url(%v), msg(%s), source(%d)", uri, msg, source)
	for i := 0; i < _retry; i++ {
		err = n.clients.Post(context.TODO(), uri, msg)
		if err != nil {
			time.Sleep(n.backoff.Backoff(i))
			continue
		} else {
			log.Info("Notify.retry callback success group(%s), topic(%s), retry(%d), msg(%s), source(%d)",
				n.w.Group, n.w.Topic, i, msg, source)
			return
		}
	}
	if err != nil {
		log.Error("Notify.retry callback error(%v), uri(%s), msg(%s), source(%d)",
			err, uri, msg, source)
	}
	return
}

// addAsyncRty asycn retry from last fail callback index.
func (n *Sub) addAsyncRty(id int64, nmsg string, cbi int) {
	if n.Closed() {
		return
	}
	select {
	case n.asyncRty <- &rtyMsg{id: id, msg: nmsg, index: cbi}:
		countProm.Incr(_opConsumerFail, n.w.Group, n.w.Topic)
	case <-n.ticker.C:
		// async chan full,stop consumer until retry sucess.
		n.stop = true
	}
}

// AddRty add retry msg to asyncretry chan by global service.
func (n *Sub) AddRty(nmsg string, id, cbi int64) {
	select {
	case n.asyncRty <- &rtyMsg{id: id, msg: nmsg, index: int(cbi)}:
		countProm.Incr(_opConsumerFail, n.w.Group, n.w.Topic)
	default:
		log.Error("sub topic %s group %s,async chan full ", n.w.Topic, n.w.Group)
	}
}

// backupMsg add failed message into db
func (n *Sub) backupMsg(msg string, cbi int) (id int64, err error) {
	id, err = n.dao.AddFailBk(context.Background(), n.w.Topic, n.w.Group, n.w.Cluster, msg, int64(cbi))
	if err != nil {
		log.Error("group(%s) topic(%s) add backup msg fail (%s)", n.w.Group, n.w.Topic, msg)
	}
	return
}

// delBackup delete failed message from db after async retry success
func (n *Sub) delBackup(id int64) {
	_, err := n.dao.DelFailBk(context.TODO(), id)
	if err != nil {
		log.Error("group(%s) topic(%s) del backup msg err(%v)", n.w.Group, n.w.Topic, err)
	}
}

// Closed return if sub had been closed.
func (n *Sub) Closed() bool {
	return n.closed
}

// Close close sub.
func (n *Sub) Close() {
	if !n.closed {
		n.closed = true
		n.consumer.Close()
		//	close(n.asyncRty)
		countProm.Decr(_opCurrentConsumer, n.w.Group, n.w.Topic)
	}
}

// IsUpdate check watch metadata if update.
func (n *Sub) IsUpdate(w *model.Watcher) bool {
	return n.w.Mtime.Unix() != w.Mtime.Unix()
}

// ParseNotifyURL Parse callback url struct by url string
func parseNotifyURL(u string) (notifyURL *model.NotifyURL, err error) {
	var parsedURL *url.URL
	parsedURL, err = url.Parse(u)
	if err != nil {
		return
	}
	notifyURL = &model.NotifyURL{
		RawURL: u,
		Schema: parsedURL.Scheme,
		Host:   parsedURL.Host,
		Path:   parsedURL.Path,
		Query:  parsedURL.Query(),
	}
	return
}
