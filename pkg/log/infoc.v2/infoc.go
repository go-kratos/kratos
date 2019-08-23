package infoc

import (
	"context"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/library/io/recordio"
	"go-common/library/log"
	"go-common/library/log/internal/core"
	"go-common/library/log/internal/filewriter"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

const (
	// LogId key to extract log id from recordio data.
	LogId         = "logId"
	InfoTimestamp = "data_timestamp"

	_infocSpliter  = "\001"
	_infocReplacer = "|"
)

var (
	// ErrFull error chan buffer full.
	ErrFull = errors.New("infoc: chan buffer full")
)

// Config contains basic configuration of infoc.
type Config struct {
	Path    string        // directory where local files are located
	Name    string        // common prefix of local files
	Rotated bool          // control file rotation of local files
	Chan    int           // buffer size
	Buffer  int           // event buffer size
	Flush   time.Duration // interval to flush buffered data to local files
}

// Payload contains log id and data to be collect.
type Payload struct {
	LogID string
	Data  []byte

	timestamp int64
}

// Infoc is interface for collector.
type Infoc interface {
	Info(context.Context, Payload) error
	Close() error
}

// infoCollector is responsible for record buffer to local files.
type infoCollector struct {
	c      *Config
	buffer chan Payload
	waiter sync.WaitGroup
	fw     *filewriter.FileWriter
	wr     *recordio.Writer
}

// New creates infoc and starts to flush buffer to local files.
func New(c *Config) (i Infoc, err error) {
	var (
		opt []filewriter.Option
		ic  *infoCollector
		nc  *Config
	)
	nc = wrapConfig(c)
	ic = &infoCollector{
		c:      nc,
		buffer: make(chan Payload, nc.Chan),
	}

	if !ic.c.Rotated {
		opt = append(opt, filewriter.MaxSize(0))
	}
	opt = append(opt, filewriter.RotateFormat("2006-01-02"))

	if ic.fw, err = filewriter.New(filepath.Join(ic.c.Path, ic.c.Name), opt...); err != nil {
		err = errors.Wrapf(err, "infoc: fail to init file writer on %s/%s", c.Path, c.Name)
		return
	}
	ic.wr = recordio.NewWriter(ic.fw)
	ic.waiter.Add(1)

	go ic.writeproc()
	i = ic
	return
}

// Info records buffer.
func (i *infoCollector) Info(ctx context.Context, payload Payload) (err error) {
	if payload.Data == nil || metadata.Bool(ctx, metadata.Mirror) {
		return
	}
	payload.timestamp = time.Now().UnixNano() / 1000000
	select {
	case i.buffer <- payload:
	default:
		err = ErrFull
	}
	return
}

// Close flush buffer to local files.
func (i *infoCollector) Close() (err error) {
	i.buffer <- Payload{}
	i.waiter.Wait()
	if err = i.wr.Flush(); err != nil {
		log.Error("infoc: fail to collect data to local file: +%v", err)
	}
	return i.fw.Close()
}

// NewLogStreamV create payload with given D
func NewLogStreamV(logID string, args ...log.D) Payload {
	if len(args) == 0 {
		return Payload{}
	}
	var buf []byte
	for i, arg := range args {
		buf = append(buf, []byte(strings.Replace(formatV(arg), _infocSpliter, _infocReplacer, -1))...)
		if i < len(args)-1 {
			buf = append(buf, []byte(_infocSpliter)...)
		}
	}
	return Payload{
		LogID: logID,
		Data:  buf,
	}
}

// NewLogStream create payload with given args
func NewLogStream(logID string, args ...interface{}) Payload {
	return NewLogStreamV(logID, buildV(args)...)
}

// writeproc write buffer to file.
func (i *infoCollector) writeproc() {
	var (
		data Payload
		r    recordio.Record
		meta map[string][]byte
		t    *time.Ticker
		bn   int
	)
	t = time.NewTicker(i.c.Flush)
	meta = make(map[string][]byte, 2)
	bn = 0

	for {
		select {
		case data = <-i.buffer:
			if data.Data == nil {
				i.waiter.Done()
				return
			}
			bn++
			meta[LogId] = []byte(data.LogID)
			meta[InfoTimestamp] = []byte(strconv.FormatInt(data.timestamp, 10))
			r = recordio.Record{
				Meta:    meta,
				Payload: data.Data,
			}
			i.wr.WriteRecord(r)
		case <-t.C:
			if bn > 0 {
				i.wr.Flush()
				bn = 0
			}
		}
	}
}

func wrapConfig(c *Config) *Config {
	nc := &Config{
		Path:    "/data/log/recordio/",
		Name:    "data.log",
		Rotated: true,
		Chan:    1024 * 8,
		Buffer:  1024,
		Flush:   time.Minute * 5,
	}
	if c != nil {
		if c.Path != "" {
			nc.Path = c.Path
		}
		if c.Name != "" {
			nc.Name = c.Name
		}
		nc.Rotated = c.Rotated
		if c.Chan > 0 {
			nc.Chan = c.Chan
		}
		if c.Buffer > 0 {
			nc.Buffer = c.Buffer
		}
		if c.Flush > 0 {
			nc.Flush = c.Flush
		}
	}
	return nc
}

func formatV(arg log.D) (s string) {
	switch arg.Type {
	case core.BoolType, core.StringType:
		s = arg.StringVal
	case core.IntTpye, core.Int64Type:
		s = strconv.FormatInt(arg.Int64Val, 10)
	case core.Float32Type:
		s = strconv.FormatFloat(math.Float64frombits(uint64(arg.Int64Val)), 'f', -1, 32)
	case core.Float64Type:
		s = strconv.FormatFloat(math.Float64frombits(uint64(arg.Int64Val)), 'f', -1, 64)
	default:
		s = fmt.Sprint(arg.Value)
	}
	return
}

func buildV(args []interface{}) (vs []log.D) {
	if len(args) == 0 {
		return []log.D{}
	}
	vs = make([]log.D, len(args))
	for i := 0; i < len(args); i = i + 1 {
		switch v := args[i].(type) {
		case int:
			vs[i] = log.Int(v)
		case int64:
			vs[i] = log.Int64(v)
		case float32:
			vs[i] = log.Float32(v)
		case float64:
			vs[i] = log.Float64(v)
		case string:
			vs[i] = log.String(v)
		case bool:
			vs[i] = log.Bool(v)
		default:
			vs[i] = log.Raw(v)
		}
	}
	return
}
