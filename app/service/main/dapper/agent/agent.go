package agent

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/pkg/deliver"
	"go-common/app/service/main/dapper/pkg/diskqueue"
	"go-common/app/service/main/dapper/server/udpcollect"
	"go-common/library/log"
	spanpb "go-common/library/net/trace/proto"
)

// Agent dapper collect agent
type Agent struct {
	queue diskqueue.DiskQueue
	clt   *udpcollect.UDPCollect
	dlv   *deliver.Deliver
}

// New agent
func New(cfg *conf.AgentConfig, debug bool) (*Agent, error) {
	var options []diskqueue.Option
	if cfg.Queue.BucketBytes != 0 {
		options = append(options, diskqueue.SetBucketByte(cfg.Queue.BucketBytes))
	}
	if cfg.Queue.MemBuckets != 0 {
		options = append(options, diskqueue.SetDynamicMemBucket(cfg.Queue.MemBuckets))
	}
	queue, err := diskqueue.New(cfg.Queue.CacheDir, options...)
	if err != nil {
		return nil, err
	}
	workers := cfg.UDPCollect.Workers
	if workers == 0 {
		if runtime.NumCPU() > 4 {
			workers = 4
		} else {
			workers = runtime.NumCPU()
		}
	}
	clt, err := udpcollect.New(cfg.UDPCollect.Addr, workers, queue.Push)
	if err != nil {
		return nil, fmt.Errorf("new udpcollect error: %s", err)
	}
	if err := clt.Start(); err != nil {
		return nil, err
	}
	ag := &Agent{queue: queue, clt: clt}
	if !debug {
		ag.dlv, err = deliver.New(cfg.Servers, ag.BlockReadFn)
	} else {
		go ag.debugPrinter()
	}
	return ag, err
}

// Reload config, only support reload servers config
func (a *Agent) Reload(cfg *conf.AgentConfig) error {
	dlv, err := deliver.New(cfg.Servers, a.BlockReadFn)
	if err != nil {
		return err
	}
	if err := a.dlv.Close(); err != nil {
		log.Warn("close old deliver error: %s", err)
	}
	a.dlv = dlv
	return nil
}

// Close agent service
func (a *Agent) Close() error {
	var messages []string
	if err := a.clt.Close(); err != nil {
		messages = append(messages, err.Error())
	}
	if a.dlv != nil {
		if err := a.dlv.Close(); err != nil {
			messages = append(messages, err.Error())
		}
	}
	if err := a.queue.Close(); err != nil {
		messages = append(messages, err.Error())
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("close agent error: %s", strings.Join(messages, "\n"))
}

// BlockReadFn wrap queue.Pop block
func (a *Agent) BlockReadFn() ([]byte, error) {
	data, err := a.queue.Pop()
	if err != io.EOF {
		return data, err
	}
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
		}
		data, err = a.queue.Pop()
		if err != io.EOF {
			return data, err
		}
	}
}

func (a *Agent) debugPrinter() {
	for {
		data, err := a.BlockReadFn()
		if err != nil {
			log.Error("read data error: %s", err)
			continue
		}
		span := new(spanpb.Span)
		if err = proto.Unmarshal(data, span); err != nil {
			log.Error("unmarshal error: %s, data: %s", err, data)
			continue
		}
		fmt.Printf("received span: %s\n", span)
		var kind bool
		var component bool
		for _, tag := range span.Tags {
			if tag.Key == "span.kind" {
				kind = true
			}
			if tag.Key == "component" {
				component = true
			}
		}
		if !kind {
			fmt.Printf("tag span.kind missing, Either \"client\" or \"server\" for the appropriate roles in an RPC, and \"producer\" or \"consumer\" for the appropriate roles in a messaging scenario.")
		}
		if !component {
			fmt.Printf("tag component missing, The software package, framework, library, or module that generated the associated Span. E.g., \"grpc\", \"django\", \"JDBI\".")
		}
	}
}
