package tcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode"

	"go-common/app/infra/databus/conf"
	"go-common/library/conf/env"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	pb "github.com/gogo/protobuf/proto"
)

func stringify(b []byte) []byte {
	return bytes.Map(
		func(r rune) rune {
			if unicode.IsSymbol(r) || unicode.IsControl(r) {
				return rune('-')
			}
			return r
		},
		b,
	)
}

type proto struct {
	prefix  byte
	integer int
	message string
}

// psCommon is pub sub common
type psCommon struct {
	c      *conn
	err    error
	closed bool
	// kafka
	group   string
	topic   string
	cluster string
	addr    string
	color   []byte
}

func newPsCommon(c *conn, group, topic, color, cluster string) (ps *psCommon) {
	ps = &psCommon{
		c:       c,
		group:   group,
		topic:   topic,
		cluster: cluster,
		color:   []byte(color),
	}
	if c != nil {
		ps.addr = c.conn.RemoteAddr().String()
	}
	return
}

func (ps *psCommon) write(protos ...proto) (err error) {
	for _, p := range protos {
		if err = ps.c.Write(p); err != nil {
			return
		}
	}
	err = ps.c.Flush()
	return
}

func (ps *psCommon) batchWrite(protos []proto) (err error) {
	if err = ps.c.Write(proto{prefix: _protoArray, integer: len(protos)}); err != nil {
		return
	}
	for _, p := range protos {
		if err = ps.c.Write(p); err != nil {
			return
		}
		// FIXME(felix): 因为ops-log性能问题先屏蔽了
		if env.DeployEnv != env.DeployEnvProd {
			log.Info("batchWrite group(%s) topic(%s) cluster(%s) color(%s) addr(%s) consumer(%s) ok", ps.group, ps.topic, ps.cluster, ps.color, ps.addr, stringify([]byte(p.message)))
		}
	}
	err = ps.c.Flush()
	return
}

func (ps *psCommon) pong() (err error) {
	if err = ps.write(proto{prefix: _protoStr, message: _pong}); err != nil {
		return
	}
	log.Info("pong group(%s) topic(%s) cluster(%s) color(%s) addr(%s) ping success", ps.group, ps.topic, ps.cluster, ps.color, ps.addr)
	return
}

func (ps *psCommon) Closed() bool {
	return ps.closed
}

// Close 跟 redis 协议耦合的太紧，加个 sendRedisErr 开关
func (ps *psCommon) Close(sendRedisErr bool) {
	if ps.closed {
		return
	}
	if ps.err == nil {
		ps.err = errConnClosedByServer // when closed by self, send close event to client.
	}
	// write error
	if ps.err != errConnRead && ps.err != errConnClosedByClient && sendRedisErr {
		ps.write(proto{prefix: _protoErr, message: ps.err.Error()})
	}
	if ps.c != nil {
		ps.c.Close()
	}
	ps.closed = true
}

func (ps *psCommon) fatal(err error) {
	if err == nil || ps.closed {
		return
	}
	ps.err = err
	ps.Close(true)
}

// Pub databus producer
type Pub struct {
	*psCommon
	// producer
	producer sarama.SyncProducer
}

// NewPub new databus producer
// http 接口复用此方法，c 传 nil
func NewPub(c *conn, group, topic, color string, pCfg *conf.Kafka) (p *Pub, err error) {
	producer, err := newProducer(group, topic, pCfg)
	if err != nil {
		log.Error("group(%s) topic(%s) cluster(%s) NewPub producer error(%v)", group, topic, pCfg.Cluster, err)
		return
	}
	p = &Pub{
		psCommon: newPsCommon(c, group, topic, color, pCfg.Cluster),
		producer: producer,
	}
	// http 协议的连接不作处理
	if c != nil {
		// set producer read connection timeout
		p.c.readTimeout = _pubReadTimeout
	}
	log.Info("NewPub() success group(%s) topic(%s) color(%s) cluster(%s) addr(%s)", group, topic, color, pCfg.Cluster, p.addr)
	return
}

// Serve databus producer goroutine
func (p *Pub) Serve() {
	var (
		err  error
		cmd  string
		args [][]byte
	)
	for {
		if cmd, args, err = p.c.Read(); err != nil {
			if err != io.EOF {
				log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) read error(%v)", p.group, p.topic, p.cluster, p.color, p.addr, err)
			}
			p.fatal(errConnRead)
			return
		}
		if p.Closed() {
			log.Warn("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) p.Closed()", p.group, p.topic, p.cluster, p.color, p.addr)
			return
		}
		select {
		case <-quit:
			p.fatal(errConnClosedByServer)
			return
		default:
		}
		switch cmd {
		case _auth:
			err = p.write(proto{prefix: _protoStr, message: _ok})
		case _ping:
			err = p.pong()
		case _set:
			if len(args) != 2 {
				p.write(proto{prefix: _protoErr, message: errPubParams.Error()})
				continue
			}
			err = p.publish(args[0], nil, args[1])
		case _hset:
			if len(args) != 3 {
				p.write(proto{prefix: _protoErr, message: errPubParams.Error()})
				continue
			}
			err = p.publish(args[0], args[1], args[2])
		case _quit:
			err = errConnClosedByClient
		default:
			err = errCmdNotSupport
		}
		if err != nil {
			p.fatal(err)
			log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) serve error(%v)", p.group, p.topic, p.cluster, p.color, p.addr, p.err)
			return
		}
	}
}

func (p *Pub) publish(key, header, value []byte) (err error) {
	if _, _, err = p.Publish(key, header, value); err != nil {
		return
	}

	return p.write(proto{prefix: _protoStr, message: _ok})
}

// Publish 发送消息 redis 和 http 协议共用
func (p *Pub) Publish(key, header, value []byte) (partition int32, offset int64, err error) {
	var message = &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
		Headers: []sarama.RecordHeader{
			{Key: _headerColor, Value: p.color},
			{Key: _headerMetadata, Value: header},
		},
	}
	now := time.Now()
	// TODO(felix): support RecordHeader
	if partition, offset, err = p.producer.SendMessage(message); err != nil {
		log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) publish(%v) error(%v)", p.group, p.topic, p.cluster, p.color, p.addr, message, err)
		return
	}
	if svc != nil {
		svc.TimeProm.Timing(p.group, int64(time.Since(now)/time.Millisecond))
		svc.CountProm.Incr(_opProducerMsgSpeed, p.group, p.topic)
	}
	// FIXME(felix): 因为ops-log性能问题先屏蔽了
	if env.DeployEnv != env.DeployEnvProd {
		log.Info("publish group(%s) topic(%s) cluster(%s) color(%s) addr(%s) key(%s) header(%s) value(%s) ok", p.group, p.topic, p.cluster, p.color, p.addr, key, stringify(header), stringify(value))
	}
	return
}

// Sub databus consumer
type Sub struct {
	*psCommon
	// kafka consumer
	consumer    *cluster.Consumer
	waitClosing bool
	batch       int
	// ticker
	ticker *time.Ticker
}

// NewSub new databus consumer
func NewSub(c *conn, group, topic, color string, sCfg *conf.Kafka, batch int64) (s *Sub, err error) {
	select {
	case <-consumerLimter:
	default:
	}
	// NOTE color 用于染色消费消息过虑
	if color != "" {
		group = fmt.Sprintf("%s-%s", group, color)
	}
	if err = validate(group, topic, sCfg.Brokers); err != nil {
		return
	}
	s = &Sub{
		psCommon: newPsCommon(c, group, topic, color, sCfg.Cluster),
		ticker:   time.NewTicker(_batchInterval),
	}
	if batch == 0 {
		s.batch = _batchNum
	} else {
		s.batch = int(batch)
	}
	// set consumer read connection timeout
	s.c.readTimeout = _subReadTimeout
	// cluster config
	cfg := cluster.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	cfg.ClientID = fmt.Sprintf("%s-%s", group, topic)
	cfg.Net.KeepAlive = 30 * time.Second
	// NOTE cluster auto commit offset interval
	cfg.Consumer.Offsets.CommitInterval = time.Second * 1
	// NOTE set fetch.wait.max.ms
	cfg.Consumer.MaxWaitTime = time.Millisecond * 250
	cfg.Consumer.MaxProcessingTime = 50 * time.Millisecond
	// NOTE errors that occur during offset management,if enabled, c.Errors channel must be read
	cfg.Consumer.Return.Errors = true
	// NOTE notifications that occur during consumer, if enabled, c.Notifications channel must be read
	cfg.Group.Return.Notifications = true
	// The initial offset to use if no offset was previously committed.
	// default: OffsetOldest
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	if s.consumer, err = cluster.NewConsumer(sCfg.Brokers, group, []string{topic}, cfg); err != nil {
		log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) cluster.NewConsumer() error(%v)", s.group, s.topic, s.cluster, s.color, s.addr, err)
	} else {
		log.Info("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) cluster.NewConsumer() ok", s.group, s.topic, s.cluster, s.color, s.addr)
	}
	return
}

func validate(group, topic string, brokers []string) (err error) {
	var (
		cli              *cluster.Client
		c                *cluster.Config
		broker           *sarama.Broker
		gresp            *sarama.DescribeGroupsResponse
		memberAssignment *sarama.ConsumerGroupMemberAssignment
		consumerNum      int
		partitions       []int32
	)
	c = cluster.NewConfig()
	c.Version = sarama.V0_10_0_1
	if cli, err = cluster.NewClient(brokers, c); err != nil {
		log.Error("group(%s) topic(%s) cluster.NewClient() error(%v)", group, topic, err)
		err = errKafKaData
		return
	}
	defer cli.Close()
	if partitions, err = cli.Partitions(topic); err != nil {
		log.Error("group(%s) topic(%s) cli.Partitions error(%v)", group, topic, err)
		err = errKafKaData
		return
	}
	if len(partitions) <= 0 {
		err = errKafKaData
		return
	}
	if err = cli.RefreshCoordinator(group); err != nil {
		log.Error("group(%s) topic(%s) cli.RefreshCoordinator error(%v)", group, topic, err)
		err = errKafKaData
		return
	}
	if broker, err = cli.Coordinator(group); err != nil {
		log.Error("group(%s) topic(%s) cli.Coordinator error(%v)", group, topic, err)
		err = errKafKaData
		return
	}
	defer broker.Close()
	if gresp, err = broker.DescribeGroups(&sarama.DescribeGroupsRequest{
		Groups: []string{group},
	}); err != nil {
		log.Error("group(%s) topic(%s) cli.DescribeGroups error(%v)", group, topic, err)
		err = errKafKaData
		return
	}
	if len(gresp.Groups) != 1 {
		err = errKafKaData
		return
	}
	for _, member := range gresp.Groups[0].Members {
		if memberAssignment, err = member.GetMemberAssignment(); err != nil {
			log.Error("group(%s) topic(%s) member.GetMemberAssignment error(%v)", group, topic, err)
			err = errKafKaData
			return
		}
		for mtopic := range memberAssignment.Topics {
			if mtopic == topic {
				consumerNum++
				break
			}
		}
	}
	if consumerNum >= len(partitions) {
		err = errUseLessConsumer
		return
	}
	return nil
}

// Serve databus consumer goroutine
func (s *Sub) Serve() {
	var (
		err  error
		cmd  string
		args [][]byte
	)
	defer func() {
		svc.CountProm.Decr(_opCurrentConsumer, s.group, s.topic)
	}()
	log.Info("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) begin serve", s.group, s.topic, s.cluster, s.color, s.addr)
	for {
		if cmd, args, err = s.c.Read(); err != nil {
			if err != io.EOF {
				log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) read error(%v)", s.group, s.topic, s.cluster, s.color, s.addr, err)
			}
			s.fatal(errConnRead)
			return
		}
		if s.consumer == nil {
			log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) s.consumer is nil", s.group, s.topic, s.cluster, s.color, s.addr)
			s.fatal(errConsumerClosed)
			return
		}
		if s.Closed() {
			log.Warn("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) s.Closed()", s.group, s.topic, s.cluster, s.color, s.addr)
			return
		}
		switch cmd {
		case _auth:
			err = s.write(proto{prefix: _protoStr, message: _ok})
		case _ping:
			err = s.pong()
		case _mget:
			var enc []byte
			if len(args) > 0 {
				enc = args[0]
			}
			err = s.message(enc)
		case _set:
			err = s.commit(args)
		case _quit:
			err = errConnClosedByClient
		default:
			err = errCmdNotSupport
		}
		if err != nil {
			s.fatal(err)
			log.Error("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) serve error(%v)", s.group, s.topic, s.cluster, s.color, s.addr, err)
			return
		}
	}
}

func (s *Sub) message(enc []byte) (err error) {
	var (
		msg    *sarama.ConsumerMessage
		notify *cluster.Notification
		protos []proto
		ok     bool
		bs     []byte
		last   = time.Now()
		ret    = &databus.MessagePB{}
		p      = proto{prefix: _protoBulk}
	)
	for {
		select {
		case err = <-s.consumer.Errors():
			log.Error("group(%s) topic(%s) cluster(%s) addr(%s) catch error(%v)", s.group, s.topic, s.cluster, s.addr, err)
			return
		case notify, ok = <-s.consumer.Notifications():
			if !ok {
				log.Info("notification notOk group(%s) topic(%s) cluster(%s) addr(%s) catch error(%v)", s.group, s.topic, s.cluster, s.addr, err)
				err = errClosedNotifyChannel
				return
			}
			switch notify.Type {
			case cluster.UnknownNotification, cluster.RebalanceError:
				log.Error("notification(%s) group(%s) topic(%s) cluster(%s) addr(%s) catch error(%v)", notify.Type, s.group, s.topic, s.cluster, s.addr, err)
				err = errClosedNotifyChannel
				return
			case cluster.RebalanceStart:
				log.Info("notification(%s) group(%s) topic(%s) cluster(%s) addr(%s) catch error(%v)", notify.Type, s.group, s.topic, s.cluster, s.addr, err)
				continue
			case cluster.RebalanceOK:
				log.Info("notification(%s) group(%s) topic(%s) cluster(%s) addr(%s) catch error(%v)", notify.Type, s.group, s.topic, s.cluster, s.addr, err)
			}
			if len(notify.Current[s.topic]) == 0 {
				log.Warn("notification(%s) no topic group(%s) topic(%s) cluster(%s) addr(%s) catch error(%v)", notify.Type, s.group, s.topic, s.cluster, s.addr, err)
				err = errConsumerOver
				return
			}
		case msg, ok = <-s.consumer.Messages():
			if !ok {
				log.Error("group(%s) topic(%s) cluster(%s) addr(%s) message channel closed", s.group, s.topic, s.cluster, s.addr)
				err = errClosedMsgChannel
				return
			}
			// reset timestamp
			last = time.Now()
			ret.Key = string(msg.Key)
			ret.Value = msg.Value
			ret.Topic = s.topic
			ret.Partition = msg.Partition
			ret.Offset = msg.Offset
			ret.Timestamp = msg.Timestamp.Unix()
			if len(msg.Headers) > 0 {
				var notMatchColor bool
				for _, h := range msg.Headers {
					if bytes.Equal(h.Key, _headerColor) && !bytes.Equal(h.Value, s.color) {
						// match color
						notMatchColor = true
					} else if bytes.Equal(h.Key, _headerMetadata) && h.Value != nil {
						// parse metadata
						dh := new(databus.Header)
						if err = pb.Unmarshal(h.Value, dh); err != nil {
							log.Error("pb.Unmarshal(%s) error(%v)", h.Value, err)
							err = nil
						} else {
							ret.Metadata = dh.Metadata
						}
					}
				}
				if notMatchColor {
					continue
				}
			}
			if bytes.Equal(enc, _encodePB) {
				// encode to pb bytes
				if bs, err = pb.Marshal(ret); err != nil {
					log.Error("proto.Marshal(%v) error(%v)", ret, err)
					s.consumer.MarkPartitionOffset(s.topic, msg.Partition, msg.Offset, "")
					return s.write(proto{prefix: _protoErr, message: errMsgFormat.Error()})
				}
			} else {
				// encode to json bytes
				if bs, err = json.Marshal(ret); err != nil {
					log.Error("json.Marshal(%v) error(%v)", ret, err)
					s.consumer.MarkPartitionOffset(s.topic, msg.Partition, msg.Offset, "")
					return s.write(proto{prefix: _protoErr, message: errMsgFormat.Error()})
				}
			}
			svc.StatProm.State(_opPartitionOffset, msg.Offset, s.group, s.topic, strconv.Itoa(int(msg.Partition)))
			svc.CountProm.Incr(_opConsumerMsgSpeed, s.group, s.topic)
			svc.StatProm.Incr(_opConsumerPartition, s.group, s.topic, strconv.Itoa(int(msg.Partition)))
			p.message = string(bs)
			protos = append(protos, p)
			if len(protos) >= s.batch {
				return s.batchWrite(protos)
			}
		case <-s.ticker.C:
			if len(protos) != 0 {
				return s.batchWrite(protos)
			}
			if time.Since(last) < _batchTimeout {
				continue
			}
			if s.waitClosing {
				log.Info("consumer group(%s) topic(%s) cluster(%s) addr(%s) wait closing then exit,maybe cluster changed", s.group, s.topic, s.cluster, s.addr)
				err = errConsumerTimeout
				return
			}
			return s.batchWrite(protos)
		}
	}
}

func (s *Sub) commit(args [][]byte) (err error) {
	var (
		partition, offset int64
	)
	if len(args) != 2 {
		log.Error("group(%v) topic(%v) cluster(%s) addr(%s) commit offset error, args(%v) is illegal", s.group, s.topic, s.cluster, s.addr, args)
		// write error
		return s.write(proto{prefix: _protoErr, message: errCommitParams.Error()})
	}
	if partition, err = strconv.ParseInt(string(args[0]), 10, 32); err != nil {
		return s.write(proto{prefix: _protoErr, message: errCommitParams.Error()})
	}
	if offset, err = strconv.ParseInt(string(args[1]), 10, 64); err != nil {
		return s.write(proto{prefix: _protoErr, message: errCommitParams.Error()})
	}
	// mark partition offset
	s.consumer.MarkPartitionOffset(s.topic, int32(partition), offset, "")
	// FIXME(felix): 因为ops-log性能问题先屏蔽了
	if env.DeployEnv != env.DeployEnvProd {
		log.Info("commit group(%s) topic(%s) cluster(%s) color(%s) addr(%s) partition(%d) offset(%d) mark offset succeed", s.group, s.topic, s.cluster, s.color, s.addr, partition, offset)
	}
	return s.write(proto{prefix: _protoStr, message: _ok})
}

// Closed judge if consumer is closed
func (s *Sub) Closed() bool {
	return s.psCommon != nil && s.psCommon.Closed()
}

// Close close consumer
func (s *Sub) Close() {
	if !s.psCommon.Closed() {
		s.psCommon.Close(true)
	}
	if s.consumer != nil {
		s.consumer.Close()
		s.consumer = nil
	}
	if s.ticker != nil {
		s.ticker.Stop()
	}
	log.Info("group(%s) topic(%s) cluster(%s) color(%s) addr(%s) consumer exit", s.group, s.topic, s.cluster, s.color, s.addr)
}

// WaitClosing marks closing state and close when consumer stoped until 30s.
func (s *Sub) WaitClosing() {
	s.waitClosing = true
}

func (s *Sub) fatal(err error) {
	if err == nil || s.closed {
		return
	}
	if s.psCommon != nil {
		s.psCommon.fatal(err)
	}
	s.Close()
}
