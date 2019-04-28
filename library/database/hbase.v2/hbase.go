package hbase

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"

	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/stat"
)

var stats = stat.DB

const (
	_family = "hbase_client_v2"
)

// Client hbase client.
type Client struct {
	hc     gohbase.Client
	addr   string
	config *Config
}

func (c *Client) setTrace(ctx context.Context, call hrpc.Call, perr *error) func() {
	now := time.Now()
	if t, ok := trace.FromContext(ctx); ok {
		t = t.Fork(_family, call.Name())
		t.SetTag(trace.String(trace.TagAddress, c.addr), trace.String(trace.TagComment, string(call.Table())+"."+string(call.Key())))
		return func() {
			t.Finish(perr)
			stats.Timing("hbase:"+call.Name(), int64(time.Since(now)/time.Millisecond))
		}
	}
	return func() {
		stats.Timing("hbase:"+call.Name(), int64(time.Since(now)/time.Millisecond))
	}
}

// NewClient new a hbase client.
func NewClient(config *Config, options ...gohbase.Option) *Client {
	zkquorum := strings.Join(config.Zookeeper.Addrs, ",")
	if config.Zookeeper.Root != "" {
		options = append(options, gohbase.ZookeeperRoot(config.Zookeeper.Root))
	}
	if config.Zookeeper.Timeout != 0 {
		options = append(options, gohbase.ZookeeperTimeout(time.Duration(config.Zookeeper.Timeout)))
	}

	if config.RPCQueueSize != 0 {
		log.Warn("RPCQueueSize configuration be ignored")
	}
	// force RpcQueueSize = 1, don't change it !!! it has reason  (゜-゜)つロ
	options = append(options, gohbase.RpcQueueSize(1))

	if config.FlushInterval != 0 {
		options = append(options, gohbase.FlushInterval(time.Duration(config.FlushInterval)))
	}
	if config.EffectiveUser != "" {
		options = append(options, gohbase.EffectiveUser(config.EffectiveUser))
	}
	if config.RegionLookupTimeout != 0 {
		options = append(options, gohbase.RegionLookupTimeout(time.Duration(config.RegionLookupTimeout)))
	}
	if config.RegionReadTimeout != 0 {
		options = append(options, gohbase.RegionReadTimeout(time.Duration(config.RegionReadTimeout)))
	}
	hc := gohbase.NewClient(zkquorum, options...)
	return &Client{
		hc:     hc,
		addr:   zkquorum,
		config: config,
	}
}

// ScanAll do scan command and return all result
// NOTE: if err != nil the results is safe for range operate even not result found
func (c *Client) ScanAll(ctx context.Context, table []byte, options ...func(hrpc.Call) error) (results []*hrpc.Result, err error) {
	cursor, err := c.Scan(ctx, table, options...)
	if err != nil {
		return nil, err
	}
	for {
		result, err := cursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

type scanTrace struct {
	hrpc.Scanner
	err         error
	cancelTrace func()
}

func (s *scanTrace) Next() (*hrpc.Result, error) {
	var result *hrpc.Result
	result, s.err = s.Scanner.Next()
	if s.err != nil {
		if s.err == io.EOF {
			// reset error for trace
			s.err = nil
			return result, io.EOF
		}
		s.cancelTrace()
		return result, s.err
	}
	return result, s.err
}

func (s *scanTrace) Close() error {
	defer s.cancelTrace()
	s.err = s.Scanner.Close()
	return s.err
}

// Scan do a scan command.
func (c *Client) Scan(ctx context.Context, table []byte, options ...func(hrpc.Call) error) (scanner hrpc.Scanner, err error) {
	var scan *hrpc.Scan
	scan, err = hrpc.NewScan(ctx, table, options...)
	if err != nil {
		return nil, err
	}

	scanner = c.hc.Scan(scan)
	st := &scanTrace{
		Scanner: scanner,
	}
	st.cancelTrace = c.setTrace(ctx, scan, &st.err)
	return st, nil
}

// ScanStr scan string
func (c *Client) ScanStr(ctx context.Context, table string, options ...func(hrpc.Call) error) (hrpc.Scanner, error) {
	return c.Scan(ctx, []byte(table), options...)
}

// ScanStrAll scan string
// NOTE: if err != nil the results is safe for range operate even not result found
func (c *Client) ScanStrAll(ctx context.Context, table string, options ...func(hrpc.Call) error) ([]*hrpc.Result, error) {
	return c.ScanAll(ctx, []byte(table), options...)
}

// ScanRange get a scanner for the given table and key range.
// The range is half-open, i.e. [startRow; stopRow[ -- stopRow is not
// included in the range.
func (c *Client) ScanRange(ctx context.Context, table, startRow, stopRow []byte, options ...func(hrpc.Call) error) (scanner hrpc.Scanner, err error) {
	var scan *hrpc.Scan
	scan, err = hrpc.NewScanRange(ctx, table, startRow, stopRow, options...)
	if err != nil {
		return nil, err
	}
	scanner = c.hc.Scan(scan)
	st := &scanTrace{
		Scanner: scanner,
	}
	st.cancelTrace = c.setTrace(ctx, scan, &st.err)
	return st, nil
}

// ScanRangeStr get a scanner for the given table and key range.
// The range is half-open, i.e. [startRow; stopRow[ -- stopRow is not
// included in the range.
func (c *Client) ScanRangeStr(ctx context.Context, table, startRow, stopRow string, options ...func(hrpc.Call) error) (hrpc.Scanner, error) {
	return c.ScanRange(ctx, []byte(table), []byte(startRow), []byte(stopRow), options...)
}

// Get get result for the given table and row key.
// NOTE: if err != nil then result != nil, if result not exists result.Cells length is 0
func (c *Client) Get(ctx context.Context, table, key []byte, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	var get *hrpc.Get
	get, err = hrpc.NewGet(ctx, table, key, options...)
	if err != nil {
		return nil, err
	}
	defer c.setTrace(ctx, get, &err)()
	return c.hc.Get(get)
}

// GetStr do a get command.
// NOTE: if err != nil then result != nil, if result not exists result.Cells length is 0
func (c *Client) GetStr(ctx context.Context, table, key string, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	return c.Get(ctx, []byte(table), []byte(key), options...)
}

// PutStr insert the given family-column-values in the given row key of the given table.
func (c *Client) PutStr(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	var put *hrpc.Mutate
	put, err = hrpc.NewPutStr(ctx, table, key, values, options...)
	if err != nil {
		return nil, err
	}
	defer c.setTrace(ctx, put, &err)()
	return c.hc.Put(put)
}

// Delete is used to perform Delete operations on a single row.
// To delete entire row, values should be nil.
//
// To delete specific families, qualifiers map should be nil:
//  map[string]map[string][]byte{
//		"cf1": nil,
//		"cf2": nil,
//  }
//
// To delete specific qualifiers:
//  map[string]map[string][]byte{
//      "cf": map[string][]byte{
//			"q1": nil,
//			"q2": nil,
//		},
//  }
//
// To delete all versions before and at a timestamp, pass hrpc.Timestamp() option.
// By default all versions will be removed.
//
// To delete only a specific version at a timestamp, pass hrpc.DeleteOneVersion() option
// along with a timestamp. For delete specific qualifiers request, if timestamp is not
// passed, only the latest version will be removed. For delete specific families request,
// the timestamp should be passed or it will have no effect as it's an expensive
// operation to perform.
func (c *Client) Delete(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	var delete *hrpc.Mutate
	delete, err = hrpc.NewDelStr(ctx, table, key, values, options...)
	defer c.setTrace(ctx, delete, &err)()
	return c.hc.Delete(delete)
}

// Append do a append command.
func (c *Client) Append(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (result *hrpc.Result, err error) {
	var append *hrpc.Mutate
	append, err = hrpc.NewAppStr(ctx, table, key, values, options...)
	defer c.setTrace(ctx, append, &err)()
	return c.hc.Append(append)
}

// Increment the given values in HBase under the given table and key.
func (c *Client) Increment(ctx context.Context, table string, key string, values map[string]map[string][]byte, options ...func(hrpc.Call) error) (result int64, err error) {
	var increment *hrpc.Mutate
	increment, err = hrpc.NewIncStr(ctx, table, key, values, options...)
	if err != nil {
		return 0, err
	}
	defer c.setTrace(ctx, increment, &err)()
	return c.hc.Increment(increment)
}

// IncrementSingle increment the given value by amount in HBase under the given table, key, family and qualifier.
func (c *Client) IncrementSingle(ctx context.Context, table string, key string, family string, qualifier string, amount int64, options ...func(hrpc.Call) error) (result int64, err error) {
	var increment *hrpc.Mutate
	increment, err = hrpc.NewIncStrSingle(ctx, table, key, family, qualifier, amount, options...)
	if err != nil {
		return 0, err
	}
	defer c.setTrace(ctx, increment, &err)()
	return c.hc.Increment(increment)
}

// Ping ping.
func (c *Client) Ping(ctx context.Context) (err error) {
	testRowKey := "test"
	if c.config.TestRowKey != "" {
		testRowKey = c.config.TestRowKey
	}
	values := map[string]map[string][]byte{"test": map[string][]byte{"test": []byte("test")}}
	_, err = c.PutStr(ctx, "test", testRowKey, values)
	return
}

// Close close client.
func (c *Client) Close() error {
	c.hc.Close()
	return nil
}
