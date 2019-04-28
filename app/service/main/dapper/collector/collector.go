package collector

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/dao"
	"go-common/app/service/main/dapper/model"
	"go-common/app/service/main/dapper/pkg/batchwrite"
	"go-common/app/service/main/dapper/pkg/collect"
	"go-common/app/service/main/dapper/pkg/collect/kafkacollect"
	"go-common/app/service/main/dapper/pkg/collect/tcpcollect"
	"go-common/app/service/main/dapper/pkg/pointwrite"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

// Collector dapper collector receving trace data from tcp and write
// to hbase and influxdb
type Collector struct {
	daoImpl dao.Dao
	cfg     *conf.Config
	bw      batchwrite.BatchWriter
	pw      pointwrite.PointWriter
	process []Processer
	clts    []collect.Collecter
	tcpClt  *tcpcollect.TCPCollect
}

func (c *Collector) checkRetention(span *model.Span) error {
	if c.cfg.Dapper.RetentionDay != 0 {
		retentionSecond := c.cfg.Dapper.RetentionDay * 3600 * 24
		if (time.Now().Unix() - span.StartTime.Unix()) > int64(retentionSecond) {
			return fmt.Errorf("span beyond retention policy ignored")
		}
	}
	return nil
}

// Process dispatch protoSpan
func (c *Collector) Process(ctx context.Context, protoSpan *model.ProtoSpan) error {
	log.V(5).Info("dispatch span form serviceName: %s, operationName: %s", protoSpan.ServiceName, protoSpan.OperationName)
	// ignored serviceName empry span
	if protoSpan.ServiceName == "" {
		log.Warn("span miss servicename ignored")
		return nil
	}

	// fix operationName
	if protoSpan.OperationName == "" {
		protoSpan.OperationName = "missing operationname !!"
	}
	if protoSpan.ServiceName == "" {
		log.Warn("span miss servicename ignored")
		return nil
	}

	// convert to model.Span
	span, err := model.FromProtoSpan(protoSpan, false)
	if err != nil {
		return err
	}
	// run process
	for _, p := range c.process {
		if err := p.Process(span); err != nil {
			return err
		}
	}

	if c.writeSamplePoint(span); err != nil {
		log.Error("write sample point error: %s", err)
	}
	if c.writeRawTrace(span); err != nil {
		log.Error("write sample point error: %s", err)
	}
	return nil
}

func (c *Collector) writeSamplePoint(span *model.Span) error {
	return c.pw.WriteSpan(span)
}

func (c *Collector) writeRawTrace(span *model.Span) error {
	return c.bw.WriteSpan(span)
}

// ListenAndStart listen and start collector server
func (c *Collector) ListenAndStart() error {
	tcpClt := tcpcollect.New(c.cfg.Collect)
	// NOTE: remove this future
	c.tcpClt = tcpClt
	c.clts = append(c.clts, tcpClt)
	// if KafkaCollect is configured enable kafka collect
	if c.cfg.KafkaCollect != nil {
		kafkaClt, err := kafkacollect.New(c.cfg.KafkaCollect.Topic, c.cfg.KafkaCollect.Addrs)
		if err != nil {
			return err
		}
		c.clts = append(c.clts, kafkaClt)
	}
	for _, clt := range c.clts {
		clt.RegisterProcess(c)
		if err := clt.Start(); err != nil {
			return err
		}
	}
	if c.cfg.Dapper.APIListen != "" {
		c.listenAndStartAPI()
	}
	return nil
}

func (c *Collector) listenAndStartAPI() {
	engine := bm.NewServer(nil)
	engine.GET("/x/internal/dapper-collector/client-status", c.clientStatusHandler)
	engine.Start()
}

// Close collector
func (c *Collector) Close() error {
	for _, clt := range c.clts {
		if err := clt.Close(); err != nil {
			log.Error("close collect error: %s", err)
		}
	}
	c.bw.Close()
	c.pw.Close()
	return nil
}

// ClientStatus .
func (c *Collector) clientStatusHandler(bc *bm.Context) {
	resp := &model.ClientStatusResp{QueueLen: c.bw.QueueLen()}
	for _, cs := range c.tcpClt.ClientStatus() {
		resp.Clients = append(resp.Clients, &model.ClientStatus{
			Addr:     cs.Addr,
			UpTime:   cs.UpTime,
			ErrCount: cs.ErrorCounter.Value(),
			Rate:     cs.Counter.Value(),
		})
	}
	bc.JSON(resp, nil)
}

// New new dapper collector
func New(cfg *conf.Config) (*Collector, error) {
	daoImpl, err := dao.New(cfg)
	if err != nil {
		return nil, err
	}
	bw := batchwrite.NewRawDataBatchWriter(
		daoImpl.WriteRawTrace,
		cfg.BatchWriter.RawBufSize,
		cfg.BatchWriter.RawChanSize,
		cfg.BatchWriter.RawWorkers,
		0,
	)

	detectData := make(map[string]map[string]struct{})
	serviceNames, err := daoImpl.FetchServiceName(context.Background())
	if err != nil {
		return nil, nil
	}
	log.V(10).Info("fetch serviceNames get %d", len(serviceNames))
	for _, serviceName := range serviceNames {
		operationNames, err := daoImpl.FetchOperationName(context.Background(), serviceName)
		if err != nil {
			log.Error("fetch operationName for %s error: %s", serviceName, err)
			continue
		}
		log.V(10).Info("fetch operationName for %s get %d", serviceName, len(operationNames))
		detectData[serviceName] = make(map[string]struct{})
		for _, operationName := range operationNames {
			detectData[serviceName][operationName] = struct{}{}
		}
	}
	// TODO configable
	pw := pointwrite.New(daoImpl.BatchWriteSpanPoint, 5, 10*time.Second)

	// register processer
	var process []Processer
	process = append(process, NewOperationNameProcess())
	// FIXME: breaker not work
	process = append(process, NewServiceBreakerProcess(4096))
	process = append(process, NewPeerServiceDetectProcesser(detectData))
	return &Collector{
		cfg:     cfg,
		daoImpl: daoImpl,
		bw:      bw,
		pw:      pw,
		process: process,
	}, nil
}
