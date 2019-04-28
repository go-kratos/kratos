package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"

	"go-common/app/service/main/dapper-query/conf"
	"go-common/app/service/main/dapper-query/model"
	"go-common/library/log"
)

// Order
const (
	TimeDesc     = "time:desc"
	TimeAsc      = "time:asc"
	DurationDesc = "duration:desc"
	DurationAsc  = "duration:asc"
)

// ErrNotFound data not found
var ErrNotFound = errors.New("not found")

// Selector to selector span
type Selector struct {
	Start     int64
	End       int64
	Limit     int
	Offset    int
	OnlyError bool
}

// table name
const (
	DefaultHbaseNameSpace = "ugc"
	DefaultInfluxDatabase = "dapper"
	HbaseRawTraceTable    = "DapperRawtrace"
	HbaseRawTraceFamily   = "pb"
	HbaseListIdxTable     = "DapperListidx"
	HbaseListIdxFamily    = "kind"
	ServiceNameTag        = "service_name"
	OperationNameTag      = "operation_name"
	PeerServiceTag        = "peer.service"
	SpanKindTag           = "span.kind"
	MaxDurationField      = "max_duration"
	MinDurationField      = "min_duration"
	AvgDurationField      = "avg_duration"
	SpanpointMeasurement  = "span_point"
	ErrorsField           = "errors"
)

// Dao dapper dao
type Dao interface {
	// list all ServiceNames
	ServiceNames(ctx context.Context) ([]string, error)
	// list OperationName for specifically service
	OperationNames(ctx context.Context, serviceName string) ([]string, error)
	// QuerySpan by family title and sel
	QuerySpanList(ctx context.Context, serviceName, operationName string, sel *Selector, order string) ([]model.SpanListRef, error)
	// Trace get trace by trace id span sort by start_time
	Trace(ctx context.Context, traceID uint64, spanIDs ...uint64) ([]*model.Span, error)
	// PeerService query all peer service depend by service
	PeerService(ctx context.Context, serviceName string) ([]string, error)
	// Ping pong
	Ping(ctx context.Context) error
	// Close dao
	Close(ctx context.Context) error
	// SupportOrder checkou order support
	SupportOrder(order string) bool
	// MeanOperationNameField
	MeanOperationNameField(ctx context.Context, whereMap map[string]string, field string, start, end int64, groupby []string) ([]model.MeanOperationNameValue, error)
	// SpanSeriesMean
	SpanSeriesMean(ctx context.Context, serviceName, operationName string, fields []string, start, end, interval int64) (*model.Series, error)
	// SpanSeriesCount
	SpanSeriesCount(ctx context.Context, serviceName, operationName string, fields []string, start, end, interval int64) (*model.Series, error)
}

type dao struct {
	hbaseNameSpace string
	hbaseClient    gohbase.Client
	influxDatabase string
	influxdbClient influxdb.Client
}

func (d *dao) Ping(ctx context.Context) error {
	return nil
}

func (d *dao) Close(ctx context.Context) error {
	d.hbaseClient.Close()
	return d.influxdbClient.Close()
}

// New dao
func New(cfg *conf.Config) (Dao, error) {
	// disable rpc queue
	hbaseClient := gohbase.NewClient(cfg.HBase.Addrs, gohbase.RpcQueueSize(0))
	hbaseNameSpace := DefaultHbaseNameSpace
	if cfg.HBase.Namespace != "" {
		hbaseNameSpace = cfg.HBase.Namespace
	}

	influxdbCfg := influxdb.HTTPConfig{Addr: cfg.InfluxDB.Addr, Username: cfg.InfluxDB.Username, Password: cfg.InfluxDB.Password}
	influxdbClient, err := influxdb.NewHTTPClient(influxdbCfg)
	if err != nil {
		return nil, err
	}
	influxDatabase := DefaultInfluxDatabase
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

func (d *dao) ServiceNames(ctx context.Context) ([]string, error) {
	where := fmt.Sprintf("%s = '%s'", SpanKindTag, "server")
	return d.showTagValues(ctx, ServiceNameTag, where)
}

func (d *dao) showTagValues(ctx context.Context, tag, where string) ([]string, error) {
	command := fmt.Sprintf(`SHOW TAG VALUES  FROM "%s" WITH KEY = "%s" WHERE %s`,
		SpanpointMeasurement, tag, where)
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
	values := make([]string, 0, len(rows.Values))
	for _, kv := range rows.Values {
		if len(kv) != 2 {
			continue
		}
		if value, ok := kv[1].(string); ok {
			values = append(values, value)
		}
	}
	return values, nil
}

func (d *dao) OperationNames(ctx context.Context, serviceName string) ([]string, error) {
	where := fmt.Sprintf("%s = '%s' AND %s = '%s'", ServiceNameTag, serviceName, SpanKindTag, "server")
	return d.showTagValues(ctx, OperationNameTag, where)
}

func (d *dao) QuerySpanList(ctx context.Context, serviceName string, operationName string, sel *Selector, order string) ([]model.SpanListRef, error) {
	log.V(10).Info("query span list serviceName: %s, operationName: %s, sel: %+v order: %s", serviceName, operationName, sel, order)
	prefix := keyPrefix(serviceName, operationName)
	startKey, stopKey := rangeKey(prefix, sel.Start, sel.End)
	switch order {
	case TimeAsc, TimeDesc:
		return d.querySpanListTimeOrder(ctx, startKey, stopKey, prefix, sel.Limit, sel.Offset, order == TimeDesc, sel.OnlyError)
	case DurationDesc, DurationAsc:
		return d.querySpanListDurationOrder(ctx, startKey, stopKey, prefix, sel.Limit, sel.Offset, order == DurationDesc, sel.OnlyError)
	}
	return nil, fmt.Errorf("unsupport order")
}

func parseSpanListRef(cell *hrpc.Cell) (spanListRef model.SpanListRef, err error) {
	value := cell.Value
	ref := bytes.SplitN(value, []byte(":"), 2)
	if len(ref) != 2 {
		err = fmt.Errorf("invalid ref %s", value)
		return
	}
	if spanListRef.TraceID, err = strconv.ParseUint(string(ref[0]), 16, 64); err != nil {
		return
	}
	spanListRef.SpanID, err = strconv.ParseUint(string(ref[1]), 16, 64)
	if err != nil {
		return
	}
	kd := bytes.SplitN(cell.Qualifier, []byte(":"), 2)
	if len(kd) != 2 {
		err = fmt.Errorf("invalid qualifier %s", cell.Qualifier)
		return
	}
	if bytes.Equal(kd[0], []byte("e")) {
		spanListRef.IsError = true
	}
	spanListRef.Duration, err = strconv.ParseInt(string(kd[1]), 10, 64)
	return
}

type minHeapifyListRef []model.SpanListRef

func (m minHeapifyListRef) push(listRef model.SpanListRef) {
	if m[0].Duration > listRef.Duration {
		return
	}
	m[0] = listRef
	m.minHeapify(0)
}

func (m minHeapifyListRef) minHeapify(i int) {
	var lowest int
	left := (i+1)*2 - 1
	right := (i + 1) * 2
	if left < len(m) && m[left].Duration < m[i].Duration {
		lowest = left
	} else {
		lowest = i
	}
	if right < len(m) && m[right].Duration < m[lowest].Duration {
		lowest = right
	}
	if lowest != i {
		m[i], m[lowest] = m[lowest], m[i]
		m.minHeapify(lowest)
	}
}

func (m minHeapifyListRef) Len() int {
	return len(m)
}

func (m minHeapifyListRef) Less(i, j int) bool {
	return m[i].Duration > m[j].Duration
}

func (m minHeapifyListRef) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (d *dao) querySpanListDurationOrder(ctx context.Context, startKey, stopKey, prefix string, limit, offset int, reverse bool, onlyError bool) ([]model.SpanListRef, error) {
	var options []func(hrpc.Call) error
	options = append(options, hrpc.Filters(filter.NewPrefixFilter([]byte(prefix))))
	if reverse {
		startKey, stopKey = stopKey, startKey
		options = append(options, hrpc.Reversed())
	}
	table := d.hbaseNameSpace + ":" + HbaseListIdxTable
	scan, err := hrpc.NewScanRangeStr(ctx, table, startKey, stopKey, options...)
	if err != nil {
		return nil, err
	}
	scanner := d.hbaseClient.Scan(scan)
	defer scanner.Close()
	spanListRefs := make(minHeapifyListRef, limit+offset)
	for {
		result, err := scanner.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if len(result.Cells) > 0 {
			log.V(10).Info("scan rowkey %s", result.Cells[0].Row)
		}
		for _, cell := range result.Cells {
			if string(cell.Family) != HbaseListIdxFamily {
				continue
			}
			spanListRef, err := parseSpanListRef(cell)
			if err != nil {
				// ignored error?
				return nil, err
			}
			if onlyError && !spanListRef.IsError {
				continue
			}
			if !reverse {
				spanListRef.Duration = math.MaxInt64 - spanListRef.Duration
			}
			spanListRefs.push(spanListRef)
		}
	}
	sort.Sort(spanListRefs)
	for i := range spanListRefs[offset:] {
		if spanListRefs[offset+i].TraceID == 0 {
			spanListRefs = spanListRefs[:offset+i]
			break
		}
	}
	return spanListRefs[offset:], nil
}

func (d *dao) querySpanListTimeOrder(ctx context.Context, startKey, stopKey, prefix string, limit, offset int, reverse bool, onlyError bool) ([]model.SpanListRef, error) {
	var options []func(hrpc.Call) error
	options = append(options, hrpc.Filters(filter.NewPrefixFilter([]byte(prefix))))
	if reverse {
		startKey, stopKey = stopKey, startKey
		options = append(options, hrpc.Reversed())
	}
	table := d.hbaseNameSpace + ":" + HbaseListIdxTable
	scan, err := hrpc.NewScanRangeStr(ctx, table, startKey, stopKey, options...)
	if err != nil {
		return nil, err
	}
	scanner := d.hbaseClient.Scan(scan)
	defer scanner.Close()
	spanListRefs := make([]model.SpanListRef, 0, limit)
	for {
		result, err := scanner.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		if len(result.Cells) > 0 {
			log.V(10).Info("scan rowkey %s", result.Cells[0].Row)
		}
		for _, cell := range result.Cells {
			if string(cell.Family) != HbaseListIdxFamily {
				continue
			}
			spanListRef, err := parseSpanListRef(cell)
			if err != nil {
				// ignored error?
				return nil, err
			}
			if onlyError && !spanListRef.IsError {
				continue
			}
			if offset > 0 {
				offset--
				continue
			}
			if limit <= 0 {
				break
			}
			if err != nil {
				// ignored error?
				return nil, err
			}
			spanListRefs = append(spanListRefs, spanListRef)
			limit--
		}
	}
	return spanListRefs, nil
}

func (d *dao) Trace(ctx context.Context, traceID uint64, spanIDs ...uint64) ([]*model.Span, error) {
	table := d.hbaseNameSpace + ":" + HbaseRawTraceTable
	traceIDStr := strconv.FormatUint(traceID, 16)
	var options []func(hrpc.Call) error
	if len(spanIDs) != 0 {
		filters := make([]filter.Filter, 0, len(spanIDs))
		for _, spanID := range spanIDs {
			spanIDStr := strconv.FormatUint(spanID, 16)
			filters = append(filters, filter.NewColumnPrefixFilter([]byte(spanIDStr)))
		}
		options = append(options, hrpc.Filters(filter.NewList(filter.MustPassOne, filters...)))
	}
	get, err := hrpc.NewGetStr(ctx, table, traceIDStr, options...)
	if err != nil {
		return nil, err
	}
	result, err := d.hbaseClient.Get(get)
	if err != nil {
		return nil, err
	}
	spans := make([]*model.Span, 0, len(result.Cells))
	for _, cell := range result.Cells {
		if string(cell.Family) != HbaseRawTraceFamily {
			continue
		}
		span, err := model.FromProtoSpan(cell.Value)
		if err != nil {
			// TODO: ignore error?
			log.Error("unmarshal protobuf span data rowkey: %x, cf: %s:%s error: %s", traceID, cell.Family, cell.Qualifier, err)
			continue
		}
		spans = append(spans, span)
	}
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].StartTime.UnixNano() > spans[j].StartTime.UnixNano()
	})
	return spans, nil
}

func (d *dao) SupportOrder(order string) bool {
	switch order {
	case TimeAsc, TimeDesc, DurationAsc, DurationDesc:
		return true
	}
	return false
}

func (d *dao) MeanOperationNameField(ctx context.Context, whereMap map[string]string, field string, start, end int64, orderby []string) ([]model.MeanOperationNameValue, error) {
	var wheres []string
	for k, v := range whereMap {
		wheres = append(wheres, fmt.Sprintf(`%s='%s'`, k, v))
	}
	command := fmt.Sprintf(`SELECT mean("%s") AS mean_%s FROM "%s" WHERE %s AND time > %ds AND time < %ds GROUP BY %s`,
		field,
		field,
		SpanpointMeasurement,
		strings.Join(wheres, " AND "),
		start, end,
		strings.Join(orderby, ", "),
	)
	log.V(10).Info("query command %s", command)
	query := influxdb.NewQuery(command, d.influxDatabase, "1s")
	resp, err := d.influxdbClient.Query(query)
	if err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 || len(resp.Results[0].Series) == 0 {
		return make([]model.MeanOperationNameValue, 0), nil
	}
	values := make([]model.MeanOperationNameValue, 0, len(resp.Results[0].Series))
	for _, row := range resp.Results[0].Series {
		value := model.MeanOperationNameValue{Tag: row.Tags}
		if len(row.Values) == 0 || len(row.Values[0]) < 2 {
			continue
		}
		valStr, ok := row.Values[0][1].(json.Number)
		if !ok {
			continue
		}
		// 相信 ifluxdb 不会瞎返回的
		value.Value, _ = valStr.Float64()
		values = append(values, value)
	}
	return values, nil
}

func (d *dao) SpanSeriesMean(ctx context.Context, serviceName, operationName string, fields []string, start, end, interval int64) (*model.Series, error) {
	return d.spanSeries(ctx, serviceName, operationName, "mean", fields, start, end, interval)
}

func (d *dao) SpanSeriesCount(ctx context.Context, serviceName, operationName string, fields []string, start, end, interval int64) (*model.Series, error) {
	return d.spanSeries(ctx, serviceName, operationName, "count", fields, start, end, interval)
}

func (d *dao) spanSeries(ctx context.Context, serviceName, operationName, fn string, fields []string, start, end, interval int64) (*model.Series, error) {
	var selects []string
	for _, field := range fields {
		selects = append(selects, fmt.Sprintf(`%s("%s") AS "%s_%s"`, fn, field, fn, field))
	}
	command := fmt.Sprintf(`SELECT %s FROM %s WHERE "%s"='%s' AND "%s"='%s' AND time > %ds AND time < %ds GROUP BY time(%ds) FILL(null)`,
		strings.Join(selects, ", "),
		SpanpointMeasurement,
		ServiceNameTag, serviceName,
		OperationNameTag, operationName,
		start, end, interval)
	log.V(10).Info("query command %s", command)
	query := influxdb.NewQuery(command, d.influxDatabase, "1s")
	resp, err := d.influxdbClient.Query(query)
	if err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 || len(resp.Results[0].Series) == 0 {
		return new(model.Series), nil
	}
	series := new(model.Series)
	fieldMap := make(map[int]*model.SeriesItem)
	for _, row := range resp.Results[0].Series {
		// first colums is time
		for i, name := range row.Columns {
			if name == "time" {
				continue
			}
			fieldMap[i] = &model.SeriesItem{Field: name}
		}
		for _, value := range row.Values {
			for i, val := range value {
				if i == 0 {
					timestamp, _ := val.(json.Number).Int64()
					series.Timestamps = append(series.Timestamps, timestamp)
					continue
				}
				n, ok := val.(json.Number)
				if !ok {
					fieldMap[i].Rows = append(fieldMap[i].Rows, nil)
				}
				v, _ := n.Float64()
				fieldMap[i].Rows = append(fieldMap[i].Rows, &v)
			}
		}
	}
	for _, v := range fieldMap {
		series.Items = append(series.Items, v)
	}
	return series, nil
}

func (d *dao) PeerService(ctx context.Context, serviceName string) ([]string, error) {
	where := fmt.Sprintf("%s = '%s' AND %s = '%s'", ServiceNameTag, serviceName, SpanKindTag, "client")
	return d.showTagValues(ctx, PeerServiceTag, where)
}
