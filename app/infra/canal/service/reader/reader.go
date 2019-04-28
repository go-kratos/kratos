// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package reader

import (
	"fmt"

	"go-common/library/log"

	"github.com/Shopify/sarama"
	pb "github.com/pingcap/tidb-tools/tidb_binlog/slave_binlog_proto/go-binlog"
	pkgerr "github.com/pkg/errors"
)

func init() {
	// log.SetLevel(log.LOG_LEVEL_NONE)
	sarama.MaxResponseSize = 1 << 30
}

// Config for Reader
type Config struct {
	KafkaAddr []string
	// the CommitTs of binlog return by reader will bigger than the config CommitTs
	CommitTS  int64
	Offset    int64 // start at kafka offset
	ClusterID string
	Name      string
}

// Message read from reader
type Message struct {
	Binlog *pb.Binlog
	Offset int64 // kafka offset
}

// Reader to read binlog from kafka
type Reader struct {
	cfg    *Config
	client sarama.Client

	msgs      chan *Message
	stop      chan struct{}
	clusterID string
}

func (r *Reader) getTopic() (string, int32) {
	return r.cfg.ClusterID + "_obinlog", 0
}

func (r *Reader) name() string {
	return fmt.Sprintf("%s-%s", r.cfg.Name, r.cfg.ClusterID)
}

// NewReader creates an instance of Reader
func NewReader(cfg *Config) (r *Reader, err error) {
	r = &Reader{
		cfg:       cfg,
		stop:      make(chan struct{}),
		msgs:      make(chan *Message, 1024),
		clusterID: cfg.ClusterID,
	}

	r.client, err = sarama.NewClient(r.cfg.KafkaAddr, nil)
	if err != nil {
		err = pkgerr.WithStack(err)
		r = nil
		return
	}
	if (r.cfg.Offset == 0) && (r.cfg.CommitTS > 0) {
		r.cfg.Offset, err = r.getOffsetByTS(r.cfg.CommitTS)
		if err != nil {
			err = pkgerr.WithStack(err)
			r = nil
			return
		}
		log.Info("tidb %s: set offset to: %v", r.name(), r.cfg.Offset)
	}
	return
}

// Close shuts down the reader
func (r *Reader) Close() {
	close(r.stop)

	r.client.Close()
}

// Messages returns a chan that contains unread buffered message
func (r *Reader) Messages() (msgs <-chan *Message) {
	return r.msgs
}

func (r *Reader) getOffsetByTS(ts int64) (offset int64, err error) {
	seeker, err := NewKafkaSeeker(r.cfg.KafkaAddr, nil)
	if err != nil {
		err = pkgerr.WithStack(err)
		return
	}

	topic, partition := r.getTopic()
	offsets, err := seeker.Seek(topic, ts, []int32{partition})
	if err != nil {
		err = pkgerr.WithStack(err)
		return
	}

	offset = offsets[0]

	return
}

// Run start consume msg
func (r *Reader) Run() {
	offset := r.cfg.Offset
	log.Info("tidb %s start at offset: %v", r.name(), offset)

	consumer, err := sarama.NewConsumerFromClient(r.client)
	if err != nil {
		log.Error("tidb %s NewConsumerFromClient err: %v", r.name(), err)
		return
	}
	defer consumer.Close()
	topic, partition := r.getTopic()
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Error("tidb %s ConsumePartition err: %v", r.name(), err)
		return
	}
	defer partitionConsumer.Close()
	for {
		select {
		case <-r.stop:
			partitionConsumer.Close()
			close(r.msgs)
			log.Info("tidb %s reader stop to run", r.name())
			return
		case kmsg, ok := <-partitionConsumer.Messages():
			if !ok {
				close(r.msgs)
				log.Info("tidb %s reader stop to run because partitionConsumer close", r.name())
				return
			}
			if kmsg == nil {
				continue
			}
			log.Info("tidb %s get kmsg offset: %v", r.name(), kmsg.Offset)
			binlog := new(pb.Binlog)
			err := binlog.Unmarshal(kmsg.Value)
			if err != nil {
				log.Warn("%s unmarshal err %+v", r.name(), err)
				continue
			}
			if r.cfg.CommitTS > 0 && binlog.CommitTs <= r.cfg.CommitTS {
				log.Warn("%s skip binlog CommitTs: ", r.name(), binlog.CommitTs)
				continue
			}
			r.msgs <- &Message{
				Binlog: binlog,
				Offset: kmsg.Offset,
			}
		}
	}
}
