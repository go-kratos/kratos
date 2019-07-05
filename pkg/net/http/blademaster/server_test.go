package blademaster

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	criticalityPkg "github.com/bilibili/kratos/pkg/net/criticality"
	"github.com/bilibili/kratos/pkg/net/metadata"
	xtime "github.com/bilibili/kratos/pkg/time"

	"github.com/stretchr/testify/assert"
)

var (
	sonce sync.Once

	curEngine atomic.Value
)

func uri(base, path string) string {
	return fmt.Sprintf("%s://%s%s", "http", base, path)
}

func shutdown() {
	if err := curEngine.Load().(*Engine).Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func setupHandler(engine *Engine) {
	// set the global timeout is 2 second
	engine.conf.Timeout = xtime.Duration(time.Second * 2)

	engine.Ping(func(ctx *Context) {
		ctx.AbortWithStatus(200)
	})

	engine.GET("/criticality/api", Criticality(criticalityPkg.Critical), func(ctx *Context) {
		ctx.String(200, "%s", metadata.String(ctx, metadata.Criticality))
	})
	engine.GET("/criticality/none/api", func(ctx *Context) {
		ctx.String(200, "%s", metadata.String(ctx, metadata.Criticality))
	})
}

func startServer(addr string) {
	e := DefaultServer(nil)
	setupHandler(e)
	go e.Run(addr)
	curEngine.Store(e)
	time.Sleep(time.Second)
}

func TestCriticality(t *testing.T) {
	addr := "localhost:18001"
	startServer(addr)
	defer shutdown()

	tests := []*struct {
		path     string
		crtl     criticalityPkg.Criticality
		expected criticalityPkg.Criticality
	}{
		{
			"/criticality/api",
			criticalityPkg.EmptyCriticality,
			criticalityPkg.Critical,
		},
		{
			"/criticality/api",
			criticalityPkg.CriticalPlus,
			criticalityPkg.Critical,
		},
		{
			"/criticality/api",
			criticalityPkg.SheddablePlus,
			criticalityPkg.Critical,
		},
	}
	client := &http.Client{}
	for _, testCase := range tests {
		req, err := http.NewRequest("GET", uri(addr, testCase.path), nil)
		assert.NoError(t, err)
		req.Header.Set("x-bm-metadata-criticality", string(testCase.crtl))
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, criticalityPkg.Criticality(body))
	}
}

func TestNoneCriticality(t *testing.T) {
	addr := "localhost:18002"
	startServer(addr)
	defer shutdown()

	tests := []*struct {
		path     string
		crtl     criticalityPkg.Criticality
		expected criticalityPkg.Criticality
	}{
		{
			"/criticality/none/api",
			criticalityPkg.EmptyCriticality,
			criticalityPkg.Critical,
		},
		{
			"/criticality/none/api",
			criticalityPkg.CriticalPlus,
			criticalityPkg.CriticalPlus,
		},
		{
			"/criticality/none/api",
			criticalityPkg.SheddablePlus,
			criticalityPkg.SheddablePlus,
		},
	}
	client := &http.Client{}
	for _, testCase := range tests {
		req, err := http.NewRequest("GET", uri(addr, testCase.path), nil)
		assert.NoError(t, err)
		req.Header.Set("x-bm-metadata-criticality", string(testCase.crtl))
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expected, criticalityPkg.Criticality(body))
	}
}
