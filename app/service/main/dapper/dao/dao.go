package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgryski/go-farm"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"

	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/model"
	"go-common/library/log"
)

const (
	defaultHbaseNameSpace = "ugc"
	defaultInfluxDatabase = "dapper"
	hbaseRawTraceTable    = "DapperRawtrace"
	hbaseRawTraceFamily   = "pb"
	hbaseListIdxTable     = "DapperListidx"
	hbaseListIdxFamily    = "kind"
	serviceNameTag        = "service_name"
	operationNameTag      = "operation_name"
	peerServiceTag        = "peer.service"
	spanKindTag           = "span.kind"
	maxDurationField      = "max_duration"
	minDurationField      = "min_duration"
	avgDurationField      = "avg_duration"
	spanpointMeasurement  = "span_point"
	errorsField           = "errors"
)

// Dao interface
type Dao interface {
	// WriteRawSpan to hbase
	WriteRawTrace(ctx context.Context, rowKey string, values map[string][]byte) error
	// BatchWriteSpanPoint
	BatchWriteSpanPoint(ctx context.Context, spanPoints []*model.SpanPoint) error
	// Fetch ServiceName
	FetchServiceName(ctx context.Context) ([]string, error)
	// Fetch OperationName
	FetchOperationName(ctx context.Context, serviceName string) ([]string, error)
}

// New dao
func New(cfg *conf.Config) (Dao, error) {
	// disable rpc queue
	hbaseClient := gohbase.NewClient(cfg.HBase.Addrs, gohbase.RpcQueueSize(0))
	hbaseNameSpace := defaultHbaseNameSpace
	if cfg.HBase.Namespace != "" {
		hbaseNameSpace = cfg.HBase.Namespace
	}

	influxdbCfg := influxdb.HTTPConfig{Addr: cfg.InfluxDB.Addr, Username: cfg.InfluxDB.Username, Password: cfg.InfluxDB.Password}
	influxdbClient, err := influxdb.NewHTTPClient(influxdbCfg)
	if err != nil {
		return nil, err
	}
	influxDatabase := defaultInfluxDatabase
	if cfg.InfluxDB.Database != "" {
		influxDatabase = cfg.InfluxDB.Database
	}

	return &dao{
		hbaseNameSpace: hbaseNameSpace,
		hbaseClient:    hbaseClient,
		influxDatabase: influxDatabase,
		influxdbClient: influxdbClient,
	}, nil
}

var _ Dao = &dao{}

type dao struct {
	hbaseNameSpace string
	hbaseClient    gohbase.Client
	influxDatabase string
	influxdbClient influxdb.Client
}

func (d *dao) WriteRawTrace(ctx context.Context, rowKey string, values map[string][]byte) error {
	table := d.hbaseNameSpace + ":" + hbaseRawTraceTable
	put, err := hrpc.NewPutStr(ctx, table, rowKey, map[string]map[string][]byte{hbaseRawTraceFamily: values})
	if err != nil {
		return err
	}
	_, err = d.hbaseClient.Put(put)
	return err
}

func (d *dao) BatchWriteSpanPoint(ctx context.Context, spanPoints []*model.SpanPoint) error {
	var messages []string
	batchPoint, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{Database: d.influxDatabase, Precision: "1s"})
	if err != nil {
		return err
	}
	for _, spanPoint := range spanPoints {
		if err := d.writeSamplePoint(ctx, spanPoint); err != nil {
			messages = append(messages, err.Error())
		}
		if point, err := toInfluxDBPoint(spanPoint); err != nil {
			messages = append(messages, err.Error())
		} else {
			batchPoint.AddPoint(point)
		}
	}
	if err := d.influxdbClient.Write(batchPoint); err != nil {
		messages = append(messages, err.Error())
	}
	if len(messages) != 0 {
		return fmt.Errorf("%s", strings.Join(messages, "\n"))
	}
	return nil
}

func (d *dao) FetchServiceName(ctx context.Context) ([]string, error) {
	command := fmt.Sprintf("SHOW TAG VALUES  FROM span_point WITH KEY = %s", serviceNameTag)
	log.V(10).Info("query command %s", command)
	query := influxdb.NewQuery(command, d.influxDatabase, "1s")
	resp, err := d.influxdbClient.Query(query)
	if err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 || len(resp.Results[0].Series) == 0 {
		return make([]string, 0), nil
	}
	rows := resp.Results[0].Series[0]
	serviceNames := make([]string, 0, len(rows.Values))
	for _, kv := range rows.Values {
		if len(kv) != 2 {
			continue
		}
		if serviceName, ok := kv[1].(string); ok {
			serviceNames = append(serviceNames, serviceName)
		}
	}
	return serviceNames, nil
}

func (d *dao) FetchOperationName(ctx context.Context, serviceName string) ([]string, error) {
	command := fmt.Sprintf("SHOW TAG VALUES  FROM %s WITH KEY = %s WHERE %s = '%s' AND %s = '%s'",
		spanpointMeasurement, operationNameTag, serviceNameTag, serviceName, spanKindTag, "server")
	log.V(10).Info("query command %s", command)
	query := influxdb.NewQuery(command, d.influxDatabase, "1s")
	resp, err := d.influxdbClient.Query(query)
	if err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 || len(resp.Results[0].Series) == 0 {
		return make([]string, 0), nil
	}
	rows := resp.Results[0].Series[0]
	operationNames := make([]string, 0, len(rows.Values))
	for _, kv := range rows.Values {
		if len(kv) != 2 {
			continue
		}
		if operationName, ok := kv[1].(string); ok {
			operationNames = append(operationNames, operationName)
		}
	}
	return operationNames, nil
}

func (d *dao) writeSamplePoint(ctx context.Context, spanPoint *model.SpanPoint) error {
	table := d.hbaseNameSpace + ":" + hbaseListIdxTable
	rowKey := listIdxKey(spanPoint)
	values := make(map[string][]byte)
	values = fuelDurationSamplePoint(values, spanPoint.MaxDuration, spanPoint.AvgDuration, spanPoint.MinDuration)
	values = fuelErrrorSamplePoint(values, spanPoint.Errors...)
	put, err := hrpc.NewPutStr(ctx, table, rowKey, map[string]map[string][]byte{hbaseListIdxFamily: values})
	if err != nil {
		return err
	}
	_, err = d.hbaseClient.Put(put)
	return err
}

func listIdxKey(spanPoint *model.SpanPoint) string {
	serviceNameHash := farm.Hash32([]byte(spanPoint.ServiceName))
	operationNameHash := farm.Hash32([]byte(spanPoint.OperationName))
	return fmt.Sprintf("%x%x%d", serviceNameHash, operationNameHash, spanPoint.Timestamp)
}

func fuelDurationSamplePoint(values map[string][]byte, points ...model.SamplePoint) map[string][]byte {
	for i := range points {
		key := "d:" + strconv.FormatInt(points[i].Value, 10)
		values[key] = []byte(fmt.Sprintf("%x:%x", points[i].TraceID, points[i].SpanID))
	}
	return values
}

func fuelErrrorSamplePoint(values map[string][]byte, points ...model.SamplePoint) map[string][]byte {
	for i := range points {
		key := "e:" + strconv.FormatInt(points[i].Value, 10)
		values[key] = []byte(fmt.Sprintf("%x:%x", points[i].TraceID, points[i].SpanID))
	}
	return values
}

func toInfluxDBPoint(spanPoint *model.SpanPoint) (*influxdb.Point, error) {
	tags := map[string]string{
		serviceNameTag:   spanPoint.ServiceName,
		operationNameTag: spanPoint.OperationName,
		spanKindTag:      spanPoint.SpanKind,
		peerServiceTag:   spanPoint.PeerService,
	}
	fields := map[string]interface{}{
		maxDurationField: spanPoint.MaxDuration.Value,
		minDurationField: spanPoint.MinDuration.Value,
		avgDurationField: spanPoint.AvgDuration.Value,
		errorsField:      len(spanPoint.Errors),
	}
	return influxdb.NewPoint(spanpointMeasurement, tags, fields, time.Unix(spanPoint.Timestamp, 0))
}
