package notify

import (
	"sync"

	"go-common/app/infra/notify/conf"
	"go-common/app/infra/notify/model"
	"go-common/library/log"

	"github.com/Shopify/sarama"
)

// Pub define producer.
type Pub struct {
	group     string
	topic     string
	cluster   string
	appsecret string
	producer  sarama.SyncProducer
}

var (
	// producer snapshot, key:group+topic
	producers = make(map[string]sarama.SyncProducer)
	pLock     sync.RWMutex
)

// NewPub new kafka producer.
func NewPub(w *model.Pub, c *conf.Config) (p *Pub, err error) {
	cluster, ok := c.Clusters[w.Cluster]
	if !ok {
		err = errClusterNotSupport
		return
	}
	producer, err := producer(w.Group, w.Topic, cluster)
	if err != nil {
		return
	}
	p = &Pub{
		group:     w.Group,
		topic:     w.Topic,
		cluster:   w.Cluster,
		appsecret: w.AppSecret,
		producer:  producer,
	}
	return
}

func producer(group, topic string, pCfg *conf.Kafka) (p sarama.SyncProducer, err error) {
	var (
		ok  bool
		key = key(group, topic)
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

// Send publish kafka message.
func (p *Pub) Send(key, value []byte) (err error) {
	countProm.Incr(_opProducerMsgSpeed, p.group, p.topic)
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = p.producer.SendMessage(message)
	return
}
func key(group, topic string) string {
	return group + ":" + topic
}

// Auth check user pub permission.
func (p *Pub) Auth(secret string) bool {
	return p.appsecret == secret
}
