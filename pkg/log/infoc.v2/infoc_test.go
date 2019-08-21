package infoc

import (
	"context"
	"os"
	"sync"
	"testing"

	"go-common/library/io/recordio"
	"go-common/library/log"

	"github.com/stretchr/testify/assert"
)

var (
	once sync.Once
	i    Infoc
)

func TestInfo(t *testing.T) {
	var (
		ctx   context.Context
		err   error
		event Payload
		index int
		meta  map[string]string
	)
	ctx = context.TODO()
	i, _ = New(&Config{
		Path:    "./",
		Name:    "infoc",
		Rotated: false,
		Chan:    4,
	})

	meta = make(map[string]string)
	meta["a"] = "b"
	event = NewLogStreamV("pb", log.Bool(true), log.Int(1), log.Int64(2), log.String("data"), log.Float32(1.2), log.Float64(3.1415926535))
	for index = 0; index < 4; index++ {
		err = i.Info(ctx, event)
		if err != nil {
			t.Fatal(err)
		}
	}

	i.Close()
}

func TestCancel(t *testing.T) {
	var (
		ctx   context.Context
		event Payload
		meta  map[string]string
		file  *os.File
		rr    *recordio.Reader
		r     recordio.Record
		err   error
	)
	ctx = context.Background()
	i, _ = New(&Config{
		Path:    "./",
		Name:    "infoc-cancel",
		Rotated: false,
		Chan:    10000,
	})

	meta = make(map[string]string)
	meta["a"] = "b"
	event = NewLogStream("pb", true, 1, int64(2), "newdata")
	i.Info(ctx, event)
	meta = make(map[string]string)
	meta["a"] = "c"
	event = Payload{
		LogID: "pb",
		Data:  []byte("old"),
	}
	i.Info(ctx, event)

	i.Close()

	file, err = os.OpenFile("infoc-cancel", os.O_RDWR|os.O_CREATE, 0644)
	rr = recordio.NewReader(file)

	r, err = rr.ReadRecord()
	assert.Nil(t, err)
	assert.Equal(t, string(r.Payload[0:4]), "true")
	assert.Equal(t, string(r.Payload[5]), "1")
	assert.Equal(t, string(r.Payload[7]), "2")
	assert.Equal(t, string(r.Payload[9:]), "newdata")

	r, err = rr.ReadRecord()
	assert.Nil(t, err)
	assert.Equal(t, string(r.Payload), "old")
}

func BenchmarkInfo(b *testing.B) {
	var (
		ctx   context.Context
		event Payload
		meta  map[string]string
		j     int
		logId string
	)
	ctx = context.Background()
	i, _ = New(&Config{
		Path:    "./",
		Name:    "infoc-bench",
		Rotated: false,
		Chan:    10000,
	})

	for j = 0; j < b.N; j++ {
		logId = string(j % 10)
		meta = make(map[string]string)
		meta["logId"] = logId
		event = Payload{
			LogID: logId,
			Data:  []byte(logId),
		}
		i.Info(ctx, event)
	}
	i.Close()
}
