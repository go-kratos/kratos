package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-common/library/cache/memcache"
	"go-common/library/container/pool"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/cache/store"
	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
)

const (
	SockAddr   = "127.0.0.1:18080"
	McSockAddr = "172.16.33.54:11211"
)

func uri(base, path string) string {
	return fmt.Sprintf("%s://%s%s", "http", base, path)
}

func init() {
	log.Init(nil)
}

func newMemcache() (*Cache, func()) {
	s := store.NewMemcache(&memcache.Config{
		Config: &pool.Config{
			Active:      10,
			Idle:        2,
			IdleTimeout: xtime.Duration(time.Second),
		},
		Name:         "test",
		Proto:        "tcp",
		Addr:         McSockAddr,
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	})
	return New(s), func() {}
}

func newFile() (*Cache, func()) {
	path, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		panic("Failed to create cache directory")
	}
	s := store.NewFile(&store.FileConfig{
		RootDir: path,
	})
	remove := func() {
		os.RemoveAll(path)
	}
	return New(s), remove
}

func TestPage(t *testing.T) {
	memcache, remove1 := newMemcache()
	filestore, remove2 := newFile()
	defer func() {
		remove1()
		remove2()
	}()

	t.Run("Memcache Store", pageCase(memcache, true))
	t.Run("File Store", pageCase(filestore, false))
}

func TestControl(t *testing.T) {
	memcache, remove1 := newMemcache()
	filestore, remove2 := newFile()
	defer func() {
		remove1()
		remove2()
	}()

	t.Run("Memcache Store", controlCase(memcache, true))
	t.Run("File Store", controlCase(filestore, false))
}

func TestPageCacheMultiWrite(t *testing.T) {
	memcache, remove1 := newMemcache()
	filestore, remove2 := newFile()
	defer func() {
		remove1()
		remove2()
	}()

	t.Run("Memcache Store", pageMultiWriteCase(memcache))
	t.Run("File Store", pageMultiWriteCase(filestore))
}

func TestDegrade(t *testing.T) {
	memcache, remove1 := newMemcache()
	filestore, remove2 := newFile()
	defer func() {
		remove1()
		remove2()
	}()

	t.Run("Memcache Store", degradeCase(memcache))
	t.Run("File Store", degradeCase(filestore))
}

func pageCase(cache *Cache, testExpire bool) func(t *testing.T) {
	return func(t *testing.T) {
		expire := int32(3)
		pc := NewPage(expire)
		engine := bm.Default()
		engine.GET("/page-cache", cache.Cache(pc, nil), func(ctx *bm.Context) {
			ctx.Writer.Header().Set("X-Hello", "World")
			ctx.String(203, "%s\n", time.Now().String())
		})

		go engine.Run(SockAddr)
		defer func() {
			engine.Server().Shutdown(context.Background())
		}()
		time.Sleep(time.Second)

		code1, content1, headers1, err1 := httpGet(uri(SockAddr, "/page-cache"))
		code2, content2, headers2, err2 := httpGet(uri(SockAddr, "/page-cache"))
		assert.Nil(t, err1)
		assert.Nil(t, err2)

		assert.Equal(t, code1, 203)
		assert.Equal(t, code2, 203)

		assert.NotNil(t, content1)
		assert.NotNil(t, content2)

		assert.Equal(t, headers1["X-Hello"], []string{"World"})
		assert.Equal(t, headers2["X-Hello"], []string{"World"})

		assert.Equal(t, string(content1), string(content2))

		if !testExpire {
			return
		}
		// test if the last caching is expired
		t.Logf("Waiting %d seconds for caching expire test", expire+1)
		time.Sleep(time.Second * time.Duration(expire+1))

		_, content3, _, err3 := httpGet(uri(SockAddr, "/page-cache"))
		_, content4, _, err4 := httpGet(uri(SockAddr, "/page-cache"))

		assert.Nil(t, err3)
		assert.Nil(t, err4)
		assert.NotNil(t, content1)
		assert.NotNil(t, content2)

		assert.NotEqual(t, string(content1), string(content3))
		assert.Equal(t, string(content3), string(content4))
	}
}

func pageMultiWriteCase(cache *Cache) func(t *testing.T) {
	return func(t *testing.T) {
		chunks := []string{
			"Hello",
			"World",
			"Hello",
			"World",
			"Hello",
			"World",
			"Hello",
			"World",
		}
		pc := NewPage(3)
		engine := bm.Default()
		engine.GET("/page-cache-write", cache.Cache(pc, nil), func(ctx *bm.Context) {
			ctx.Writer.Header().Set("X-Hello", "World")
			ctx.Writer.WriteHeader(203)
			for _, chunk := range chunks {
				ctx.Writer.Write([]byte(chunk))
			}
		})

		go engine.Run(SockAddr)
		defer func() {
			engine.Server().Shutdown(context.Background())
		}()
		time.Sleep(time.Second)

		code1, content1, headers1, err1 := httpGet(uri(SockAddr, "/page-cache-write"))
		code2, content2, headers2, err2 := httpGet(uri(SockAddr, "/page-cache-write"))
		assert.Nil(t, err1)
		assert.Nil(t, err2)

		assert.Equal(t, code1, 203)
		assert.Equal(t, code2, 203)

		assert.NotNil(t, content1)
		assert.NotNil(t, content2)

		assert.Equal(t, headers1["X-Hello"], []string{"World"})
		assert.Equal(t, headers2["X-Hello"], []string{"World"})

		assert.Equal(t, strings.Join(chunks, ""), string(content1))
		assert.Equal(t, strings.Join(chunks, ""), string(content2))
		assert.Equal(t, string(content1), string(content2))
	}
}

func degradeCase(cache *Cache) func(t *testing.T) {
	return func(t *testing.T) {
		wg := sync.WaitGroup{}
		i := int32(0)
		degrade := NewDegrader(10)
		engine := bm.Default()
		engine.GET("/scheduled/error", cache.Cache(degrade.Args("name", "age"), nil), func(c *bm.Context) {
			code := atomic.AddInt32(&i, 1)
			if code == 5 {
				c.JSON("succeed", nil)
				return
			}
			if code%2 == 0 {
				c.JSON("", ecode.Degrade)
				return
			}
			c.JSON(fmt.Sprintf("Code: %d", code), ecode.Int(int(code)))
		})

		wg.Add(1)
		go func() {
			engine.Run(":18080")
			wg.Done()
		}()
		defer func() {
			engine.Server().Shutdown(context.TODO())
			wg.Wait()
		}()
		time.Sleep(time.Second)

		for index := 1; index < 10; index++ {
			_, content, _, _ := httpGet(uri(SockAddr, "/scheduled/error?name=degrader&age=26"))
			t.Log(index, string(content))
			var res struct {
				Data string `json:"data"`
			}
			err := json.Unmarshal(content, &res)
			assert.Nil(t, err)

			if index == 5 {
				// ensure response is write to cache
				time.Sleep(time.Second)
			}
			if index > 5 && index%2 == 0 {
				if res.Data != "succeed" {
					t.Fatalf("Failed to degrade at index: %d", index)
				} else {
					t.Logf("This request is degraded at index: %d", index)
				}
			}
		}
	}
}

func controlCase(cache *Cache, testExpire bool) func(t *testing.T) {
	return func(t *testing.T) {
		wg := sync.WaitGroup{}
		i := int32(0)
		expire := int32(30)
		control := NewControl(expire)
		filter := func(ctx *bm.Context) bool {
			if ctx.Request.Form.Get("cache") == "false" {
				return false
			}
			return true
		}
		engine := bm.Default()
		engine.GET("/large/response", cache.Cache(control, filter), func(c *bm.Context) {
			c.JSON(map[string]interface{}{
				"index":  atomic.AddInt32(&i, 1),
				"Hello0": "World",
				"Hello1": "World",
				"Hello2": "World",
				"Hello3": "World",
				"Hello4": "World",
				"Hello5": "World",
				"Hello6": "World",
				"Hello7": "World",
				"Hello8": "World",
			}, nil)
		})
		engine.GET("/large/response/error", cache.Cache(control, filter), func(c *bm.Context) {
			c.JSON(nil, ecode.RequestErr)
		})

		wg.Add(1)
		go func() {
			engine.Run(":18080")
			wg.Done()
		}()
		defer func() {
			engine.Server().Shutdown(context.TODO())
			wg.Wait()
		}()
		time.Sleep(time.Second)

		code, content, headers, err := httpGet(uri(SockAddr, "/large/response?name=hello&age=1"))
		assert.NoError(t, err)
		assert.Equal(t, 200, code)
		assert.NotEmpty(t, content)
		assert.Equal(t, "max-age=30", headers.Get("Cache-Control"))
		exp, err := http.ParseTime(headers.Get("Expires"))
		assert.NoError(t, err)
		assert.InDelta(t, 30, exp.Unix()-time.Now().Unix(), 5)

		code, content, headers, err = httpGet(uri(SockAddr, "/large/response/error?name=hello&age=1&cache=false"))
		assert.NoError(t, err)
		assert.Equal(t, 200, code)
		assert.NotEmpty(t, content)
		assert.Empty(t, headers.Get("Expires"))
		assert.Empty(t, headers.Get("Cache-Control"))

		code, content, headers, err = httpGet(uri(SockAddr, "/large/response/error?name=hello&age=1"))
		assert.NoError(t, err)
		assert.Equal(t, 200, code)
		assert.NotEmpty(t, content)
		assert.Empty(t, headers.Get("Expires"))
		assert.Empty(t, headers.Get("Cache-Control"))
	}
}

func httpGet(url string) (code int, content []byte, headers http.Header, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if content, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	code = resp.StatusCode
	headers = resp.Header
	return
}
