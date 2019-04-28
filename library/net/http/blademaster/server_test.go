package blademaster

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/net/http/blademaster/tests"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	sonce sync.Once

	once      sync.Once
	curEngine atomic.Value

	ByteContent = []byte(`<html>
	<meta http-equiv="refresh" content="0;uri=http://www.bilibili.com/">
	</html>`)
	CertPEM = `-----BEGIN CERTIFICATE-----
MIIDJzCCAg8CCQDHIbk1Vp7UbzANBgkqhkiG9w0BAQsFADCBkDELMAkGA1UEBhMC
Q04xETAPBgNVBAgMCFNoYW5naGFpMREwDwYDVQQHDAhTaGFuZ2hhaTERMA8GA1UE
CgwIYmlsaWJpbGkxETAPBgNVBAsMCGJpbGliaWxpMRUwEwYDVQQDDAxiaWxpYmls
aS5jb20xHjAcBgkqhkiG9w0BCQEWD2l0QGJpbGliaWxpLmNvbTAeFw0xNzExMDcw
NDI1MzJaFw0xODExMDcwNDI1MzJaMBoxGDAWBgNVBAMMD2xvY2FsaG9zdDoxODA4
MTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALNKS9f7sEE+xx2SiIgI
tfmUVi3HkbWRWiZxZ82gRQnYUvgNsQUenyMH4ukViC7rEqBomscP0YjWhIvunhKe
RzXqWxZyKF86lL0+n1X+USESWdMQe8nYhCAwTE3JIykBrqjEiYMSI5TLwQrqFUJ9
nd7EywdlUgolJFO2pbltU9a8stlto9OOLXo5veb30nAW5tnDF5Q1jlKBRpGV4+Wy
3Tn9V9a6mPaoLQOLQzLWfjIWok0UKdYOWZUwmfboFloI0J0VA8Dn3qr2VGEucUG4
C5pIzV7/ke0Ymca8H2O1Gt5jrhbieY1XLP7NEoic1xdKTa6TLbReWTUEfqErCD3X
b28CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAS+tB9PDV0tzFYtVzaWGrqypfnbEC
l5XoT6kRKn3Pf74MQQDMmuPCNZqf8aznx7+inzf4xeMsM68mkbaWvDkD2W8rSn47
tnFZNPLzlP5sXrt1RgfEK1azVOqX+PSNqDxFetB/smvsVr0EduX0tcmSNMChzx7I
Igy/I93TVf/hzu3MubZlCjTIQEvfek/Qc/eei7SQYS3dauSKaLfOwMdan9U2gmSr
byb4f0vI1wuBSAEceJMrHcPGNgibAUMBMdSOYljYxSgmC0gFW68aD3gdn1Z/KOyd
r1VaEkBHRoXvVUYPrFDFYO4nP65aZBLWgIn/EtilNektlAZhljSzk6bWXA==
-----END CERTIFICATE-----
`
	KeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAs0pL1/uwQT7HHZKIiAi1+ZRWLceRtZFaJnFnzaBFCdhS+A2x
BR6fIwfi6RWILusSoGiaxw/RiNaEi+6eEp5HNepbFnIoXzqUvT6fVf5RIRJZ0xB7
ydiEIDBMTckjKQGuqMSJgxIjlMvBCuoVQn2d3sTLB2VSCiUkU7aluW1T1ryy2W2j
044tejm95vfScBbm2cMXlDWOUoFGkZXj5bLdOf1X1rqY9qgtA4tDMtZ+MhaiTRQp
1g5ZlTCZ9ugWWgjQnRUDwOfeqvZUYS5xQbgLmkjNXv+R7RiZxrwfY7Ua3mOuFuJ5
jVcs/s0SiJzXF0pNrpMttF5ZNQR+oSsIPddvbwIDAQABAoIBACLW2hwTNXHIQCr3
8X31q17fO0vUDvVoVEtKGDC15xR9q8m152MmSygkfSxr2bW8SjdPfFwYL9BWVxVV
/fOCPDY23xJihoPSv1py08WDEMLLbRb9igB0CWCz4e/vmNx8DjOPVWVZ3f4pBc8Y
I59zB31lYkqCnsmH5CI8SMHag8MjPmzc7TtoROUbV0rcyoOx/JwrgwuGsC/UP9GY
Oj021xiLcwP5qD4sBIuXxIPx9zwCDowujjBQ3ViSgcQk92Z1fsPMBvKYUuP2hHcZ
M/Dnuzz/OyzfP0u3Aq+VpQlXHVmEptU6kfFjK8X3J9tWr7PAvAaFqIs2xL5p6gz8
q50gfzkCgYEA5Wa89L30UD17s+2RT2Slg5fmSvn2gO0Dz2Dhav/G2UPfbW+22rAC
iotIfnF4IuA6QQX6tW5L/nVbLYQUtWNxzWsYCQAbRrBQi1BBfh0EOli52ZicuWd3
6rOqeOzqsXRdyrwVnpfKf1Hdh7Dc++zG920ktXbC33jgGDLmSLxnysMCgYEAyBQf
du67z3//yOt3z+zl1Tug3t/ioPSWijrlUtsqbRuOWvxRwUtyOB4Tm171j5sWPiJu
YOisOvjrRI2L9xjdtP4HTAiO4//YiNkCHdFiyHTzMqb1RBbHVQZ6eFmjvdmrgIkG
4vXt3mZ1kQY264APB99W+BKFtLDPvaz+Hgy9xeUCgYEAhPMcA8OrOm3Hqam/k4HD
Ixb/0ug3YtT6Zk/BlN+UAQsDDEu4b9meP2klpJJii+PkHxc2C7xWsqyVITXxQobV
x7WPgnfbVwaMR5FFw69RafdODrwR6Kn8p7tkyxyTkDDewsZqyTUzmMJ7X06zZBX/
4hoRMlIX8qf9SEkHiZQXmz0CgYAGNamsVUh67iwQHk6/o0iWz5z0jdpIyI6Lh7xq
T+cHL93BMSeQajkHSNeu8MmKRXPxRbxLQa1mvyb+H66CYsEuxtuPHozgwqYDyUhp
iIAaXJbXsZrXHCXfm63dYlrUn5bVDGusS5mwV1m6wIif0n+k7OeUF28S5pHr/xx7
7kVNiQKBgDeXBYkZIRlm0PwtmZcoAU9bVK2m7MnLj6G/MebpzSwSFbbFhxZnnwS9
EOUtzcssIiFpJiN+ZKFOV+9S9ylYFSYlxGwrWvE3+nmzm8X04tsYhvjO5Q1w1Egs
U31DFsHpaXEAfBc7ugLimKKKbptJJDPzUfUtXIltCxI9YCzFFRXY
-----END RSA PRIVATE KEY-----
`
)

const (
	SockAddr      = "localhost:18080"
	AdHocSockAddr = "localhost:18081"
)

func cleanup(dir string) {
	matches, err := filepath.Glob(dir)
	if err != nil {
		return
	}

	for _, p := range matches {
		os.RemoveAll(p)
	}
}

func setupLogging() {
	dir, err := ioutil.TempDir("", "blademaster")
	defer func() {
		baseDir := path.Dir(dir)
		cleanup(path.Join(baseDir, "blademaster*"))
	}()

	if err != nil {
		panic(err)
	}
	for l, f := range map[string]string{
		"info.log":    "/dev/stdout",
		"warning.log": "/dev/stdout",
		"error.log":   "/dev/stderr",
	} {
		p := filepath.Join(dir, l)
		if err := os.Symlink(f, p); err != nil {
			panic(err)
		}
	}
	logConf := &log.Config{
		Dir: dir,
	}
	log.Init(logConf)
}

func init() {
	if os.Getenv("TEST_LOGGING") == "1" {
		setupLogging()
	}
}

func startTestServer() {
	_httpDSN = "tcp://0.0.0.0:18090/?maxlisten=20000&timeout=1s&readtimeout=1s&writetimeout=1s"
	engine := NewServer(nil)
	engine.GET("/test", func(ctx *Context) {
		ctx.JSON("", nil)
	})
	engine.GET("/mirrortest", func(ctx *Context) {
		ctx.JSON(strconv.FormatBool(metadata.Bool(ctx, metadata.Mirror)), nil)
	})
	engine.Start()
}

func TestServer2(t *testing.T) {
	sonce.Do(startTestServer)
	resp, err := http.Get("http://localhost:18090/test")
	if err != nil {
		t.Errorf("HTTPServ: get error(%v)", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("http.Get get error code:%d", resp.StatusCode)
	}
	resp.Body.Close()
}

func BenchmarkServer2(b *testing.B) {
	once.Do(startServer)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get("http://localhost:18090/test")
			if err != nil {
				b.Errorf("HTTPServ: get error(%v)", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				b.Errorf("HTTPServ: get error status code:%d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}

func TestDSNParse(t *testing.T) {
	conf := parseDSN("tcp://0.0.0.0:18090/?maxlisten=20000&timeout=1s&readTimeout=1s&writeTimeout=1s")
	assert.Equal(t, ServerConfig{
		Network:      "tcp",
		Addr:         "0.0.0.0:18090",
		Timeout:      xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}, *conf)
	conf = parseDSN("unix:///temp/bm.sock?maxlisten=20000&timeout=2s&readTimeout=1s&writeTimeout=1s")
	assert.Equal(t, ServerConfig{
		Network:      "unix",
		Addr:         "/temp/bm.sock",
		Timeout:      xtime.Duration(time.Second * 2),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}, *conf)
}

type Timestamp struct {
	Now int64 `json:"now"`
}

type Archive struct {
	Mids    []int64 `json:"mids" form:"mids,split" validate:"dive,gt=1,required"`
	Cid     int     `json:"cid" form:"cid" validate:"max=10,min=1"`
	Title   string  `json:"title" form:"title" validate:"required"`
	Content string  `json:"content" form:"content"`
}

type DevResp struct {
	Code int     `json:"code"`
	Data *Device `json:"data"`
}

func uri(base, path string) string {
	return fmt.Sprintf("%s://%s%s", "http", base, path)
}

func now() Timestamp {
	return Timestamp{
		Now: time.Now().Unix(),
	}
}

func setupHandler(engine *Engine) {
	// engine.GET("/", func(ctx *Context) {
	// 	ctx.Status(200)
	// })

	engine.Use(CORS(), CSRF(), Mobile())
	// set the global timeout is 2 second
	engine.conf.Timeout = xtime.Duration(time.Second * 2)

	engine.Ping(func(ctx *Context) {
		ctx.AbortWithStatus(200)
	})
	engine.Register(func(ctx *Context) {
		ctx.JSONMap(map[string]interface{}{
			"region": "aws",
		}, nil)
	})
	engine.HEAD("/head", func(ctx *Context) {
		ctx.Status(200)
	})

	engine.GET("/get", func(ctx *Context) {
		ctx.String(200, "%s", "pong")
	})

	engine.POST("/post", func(ctx *Context) {
		ctx.String(200, "%s", "pong")
	})

	engine.PUT("/put", func(ctx *Context) {
		ctx.String(200, "%s", "pong")
	})

	engine.DELETE("/delete", func(ctx *Context) {
		ctx.String(200, "%s", "pong")
	})

	engine.GET("/json", func(ctx *Context) {
		ctx.JSON(now(), nil)
	})

	engine.GET("/null-json", func(ctx *Context) {
		ctx.JSON(nil, nil)
	})

	engine.GET("/err", func(ctx *Context) {
		ctx.JSON(now(), errors.New("A error raised from handler func"))
	})

	engine.GET("/xml", func(ctx *Context) {
		ctx.XML(now(), nil)
	})

	engine.GET("/bytes", func(ctx *Context) {
		ctx.Bytes(200, "text/html", ByteContent)
	})

	engine.GET("/bench", func(ctx *Context) {
		ctx.JSON(now(), nil)
	})

	engine.GET("/sleep5", func(ctx *Context) {
		time.Sleep(time.Second * 30)
		ctx.JSON(now(), nil)
	})

	engine.GET("/timeout", func(ctx *Context) {
		start := time.Now()
		<-ctx.Done()
		ctx.String(200, "Timeout within %s", time.Since(start))
	})

	engine.GET("/timeout-from-method-config", func(ctx *Context) {
		start := time.Now()
		<-ctx.Done()
		ctx.String(200, "Timeout within %s", time.Since(start))
	})
	engine.SetMethodConfig("/timeout-from-method-config", &MethodConfig{Timeout: xtime.Duration(time.Second * 3)})

	engine.GET("/redirect", func(ctx *Context) {
		ctx.Redirect(301, "/bytes")
	})

	engine.GET("/panic", func(_ *Context) {
		panic("Oupps, Houston, we have a problem")
	})

	engine.GET("/json-map", func(ctx *Context) {
		ctx.JSONMap(map[string]interface{}{
			"tid": 1,
		}, nil)
	})

	engine.GET("/bind", func(ctx *Context) {
		v := new(Archive)

		err := ctx.Bind(v)
		if err != nil {
			return
		}
		ctx.JSON(v, nil)
	})

	engine.POST("/bindwith", func(ctx *Context) {
		v := new(Archive)

		err := ctx.BindWith(v, binding.JSON)
		if err != nil {
			return
		}
		ctx.JSON(v, nil)
	})

	engine.GET("/pb", func(ctx *Context) {
		now := &tests.Time{
			Now: time.Now().Unix(),
		}
		ctx.Protobuf(now, nil)
	})

	engine.GET("/pb-error", func(ctx *Context) {
		ctx.Protobuf(nil, ecode.RequestErr)
	})

	engine.GET("/pb-nildata-nilerr", func(ctx *Context) {
		ctx.Protobuf(nil, nil)
	})

	engine.GET("/pb-data-err", func(ctx *Context) {
		now := &tests.Time{
			Now: time.Now().Unix(),
		}
		ctx.Protobuf(now, ecode.RequestErr)
	})

	engine.GET("/member-blocked", func(ctx *Context) {
		ctx.JSON(nil, ecode.MemberBlocked)
	})

	engine.GET("/device", func(ctx *Context) {
		dev, ok := ctx.Get("device")
		if !ok {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
		ctx.JSON(dev, nil)
	})

	engine.GET("/device-from-meta", func(ctx *Context) {
		dev, ok := metadata.Value(ctx, "device").(*Device)
		if !ok {
			ctx.JSON(nil, ecode.RequestErr)
			return
		}
		ctx.JSON(dev, nil)
	})

	engine.GET("/remote-ip", func(ctx *Context) {
		ctx.JSON(metadata.String(ctx, metadata.RemoteIP), nil)
	})

	m1 := func() HandlerFunc {
		return func(ctx *Context) {
			ctx.Set("m1", "middleware pong")
			ctx.Next()
		}
	}
	m2 := func() func(ctx *Context) {
		return func(ctx *Context) {
			v, isok := ctx.Get("m1")
			if !isok {
				ctx.AbortWithStatus(500)
				return
			}
			ctx.Set("m2", v)
			ctx.Next()
		}
	}
	engine.Use(m1())
	engine.UseFunc(m2())
	engine.GET("/use", func(ctx *Context) {
		v, isok := ctx.Get("m2")
		if !isok {
			ctx.AbortWithStatus(500)
			return
		}
		ctx.String(200, "%s", v.(string))
	})

	r := engine.Group("/group", func(ctx *Context) {
	})
	r.GET("/get", func(ctx *Context) {
		ctx.String(200, "%s", "pong")
	})

	rr := r.Group("/abort", func(ctx *Context) {
		ctx.String(200, "%s", "pong")
		ctx.Abort()
	})
	rr.GET("", func(ctx *Context) {
		ctx.String(500, "never get this echo")
	})

	rrr := r.Group("/test", func(ctx *Context) {
	})
	g1 := func() HandlerFunc {
		return func(ctx *Context) {
			v, isok := ctx.Get("m2")
			if !isok {
				ctx.AbortWithStatus(500)
				return
			}
			ctx.Set("g1", v)
			ctx.Next()
		}
	}
	g2 := func() func(ctx *Context) {
		return func(ctx *Context) {
			v, isok := ctx.Get("g1")
			if !isok {
				ctx.AbortWithStatus(500)
				return
			}
			ctx.Next()
			ctx.String(200, "%s", v.(string))
		}
	}
	rrr.Use(g1()).UseFunc(g2()).GET("/use", func(ctx *Context) {})

	groupInject := engine.Group("/inject")
	engine.Inject("^/inject", func(ctx *Context) {
		ctx.Set("injected", "injected")
	})
	groupInject.GET("/index", func(ctx *Context) {
		injected, _ := ctx.Get("injected")
		injectedString, _ := injected.(string)
		ctx.String(200, strings.Join([]string{"index", injectedString}, "-"))
	})

	engine.GET("/group/test/json-status", func(ctx *Context) {
		ctx.JSON(nil, ecode.MemberBlocked)
	})
	engine.GET("/group/test/xml-status", func(ctx *Context) {
		ctx.XML(nil, ecode.MemberBlocked)
	})
	engine.GET("/group/test/proto-status", func(ctx *Context) {
		ctx.Protobuf(nil, ecode.MemberBlocked)
	})
	engine.GET("/group/test/json-map-status", func(ctx *Context) {
		ctx.JSONMap(map[string]interface{}{}, ecode.MemberBlocked)
	})
}

func startServer() {
	e := Default()
	setupHandler(e)
	go e.Run(SockAddr)
	curEngine.Store(e)
	time.Sleep(time.Second)
}

func TestSetupHandler(t *testing.T) {
	engine := Default()
	setupHandler(engine)
}

func TestServeUnix(t *testing.T) {
	engine := Default()
	setupHandler(engine)
	closed := make(chan struct{})
	defer func() {
		if err := engine.Shutdown(context.TODO()); err != nil {
			t.Errorf("Failed to shutdown engine: %s", err)
		}
		<-closed
	}()
	unixs, err := ioutil.TempFile("", "engine.sock")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}

	go func() {
		if err := engine.RunUnix(unixs.Name()); err != nil {
			if errors.Cause(err) == http.ErrServerClosed {
				t.Logf("Server stopped due to shutting down command")
			} else {
				t.Errorf("Failed to serve with unix socket: %s", err)
			}
		}
		closed <- struct{}{}
	}()
	// connection test required
	time.Sleep(time.Second)
}

func shutdown() {
	if err := curEngine.Load().(*Engine).Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func TestServeTLS(t *testing.T) {
	engine := New()
	setupHandler(engine)
	closed := make(chan struct{})
	defer func() {
		if err := engine.Shutdown(context.TODO()); err != nil {
			t.Errorf("Failed to shutdown engine: %s", err)
		}
		<-closed
	}()

	cert, err := ioutil.TempFile("", "cert.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	key, err := ioutil.TempFile("", "key.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	if _, err = cert.WriteString(CertPEM); err != nil {
		t.Fatalf("Failed to write cert file: %s", err)
	}
	if _, err = key.WriteString(KeyPEM); err != nil {
		t.Fatalf("Failed to write key file: %s", err)
	}

	go func() {
		if rerr := engine.RunTLS(AdHocSockAddr, cert.Name(), key.Name()); rerr != nil {
			if errors.Cause(rerr) == http.ErrServerClosed {
				t.Logf("Server stopped due to shutting down command")
			} else {
				t.Errorf("Failed to serve with tls: %s", rerr)
			}
		}
		closed <- struct{}{}
	}()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fmt.Sprintf("%s://%s%s", "https", AdHocSockAddr, "/get"))
	if err != nil {
		t.Fatalf("Failed to send https request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("http.Get get error code:%d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestServer(t *testing.T) {
	once.Do(startServer)

	resp, err := http.Get(uri(SockAddr, "/get"))
	if err != nil {
		t.Fatalf("BladeMaster: get error(%v)", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("http.Get get error code:%d", resp.StatusCode)
	}
	resp.Body.Close()
}

func pongValidator(expected string, t *testing.T) func(*http.Response) error {
	return func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http.Get get error code: %d", resp.StatusCode)
		}
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Logf("ioutil.ReadAll() failed: %s", err)
			return err
		}
		resp.Body.Close()
		if ps := string(bs); ps != expected {
			return fmt.Errorf("Response body not expected: %s != %s", ps, expected)
		}
		return nil
	}
}

func TestPathHandle(t *testing.T) {
	once.Do(startServer)

	expected := "pong"
	validateFn := pongValidator(expected, t)

	testCase := map[string]map[string]func(*http.Response) error{
		"/monitor/ping": {
			"GET": func(resp *http.Response) error {
				assert.Equal(t, 200, resp.StatusCode)
				return nil
			},
		},

		"/register": {
			"GET": func(resp *http.Response) error {
				assert.Equal(t, 200, resp.StatusCode)
				bs, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)
				md := make(map[string]interface{})
				assert.NoError(t, json.Unmarshal(bs, &md))
				assert.Equal(t, "aws", md["region"])
				assert.Equal(t, 0, int(md["code"].(float64)))
				assert.Equal(t, "0", md["message"])
				return nil
			},
		},

		"/head": {
			"HEAD": func(resp *http.Response) error {
				assert.Equal(t, 200, resp.StatusCode)
				return nil
			},
		},

		"/get": {
			"GET": validateFn,
		},

		"/post": {
			"POST": validateFn,
		},

		"/put": {
			"PUT": validateFn,
		},

		"/delete": {
			"DELETE": validateFn,
		},

		"/not-exist-path": {
			"GET": func(resp *http.Response) error {
				assert.Equal(t, 404, resp.StatusCode)
				return nil
			},
		},
	}

	c := &http.Client{}
	for path := range testCase {
		for method := range testCase[path] {
			validator := testCase[path][method]
			req, err := http.NewRequest(method, uri(SockAddr, path), nil)
			if err != nil {
				t.Fatalf("Failed to build request: %s", err)
			}
			resp, err := c.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %s", err)
			}
			defer resp.Body.Close()

			if err := validator(resp); err != nil {
				t.Errorf("Failed to validate request to `%s %s`: %s", method, path, err)
			}
		}
	}
}

func TestMiddleWare(t *testing.T) {
	once.Do(startServer)

	expected := "middleware pong"
	validateFn := pongValidator(expected, t)

	testCase := map[string]map[string]func(*http.Response) error{
		"/use": {
			"GET": validateFn,
		},

		"/group/test/use": {
			"GET": validateFn,
		},
	}

	c := &http.Client{}
	for path := range testCase {
		for method := range testCase[path] {
			validator := testCase[path][method]
			req, err := http.NewRequest(method, uri(SockAddr, path), nil)
			if err != nil {
				t.Fatalf("Failed to build request: %s", err)
			}
			resp, err := c.Do(req)
			if err != nil {
				t.Errorf("Failed to send request: %s", err)
			}
			defer resp.Body.Close()

			if err := validator(resp); err != nil {
				t.Errorf("Failed to validate request to `%s %s`: %s", method, path, err)
			}
		}
	}
}

func TestRouterGroup(t *testing.T) {
	once.Do(startServer)

	expected := "pong"
	validateFn := pongValidator(expected, t)

	testCase := map[string]map[string]func(*http.Response) error{
		"/group/get": {
			"GET": validateFn,
		},

		"/group/abort": {
			"GET": validateFn,
		},
	}

	c := &http.Client{}
	for path := range testCase {
		for method := range testCase[path] {
			validator := testCase[path][method]
			req, err := http.NewRequest(method, uri(SockAddr, path), nil)
			if err != nil {
				t.Fatalf("Failed to build request: %s", err)
			}
			resp, err := c.Do(req)
			if err != nil {
				t.Errorf("Failed to send request: %s", err)
			}
			defer resp.Body.Close()

			if err := validator(resp); err != nil {
				t.Errorf("Failed to validate request to `%s %s`: %s", method, path, err)
			}
		}
	}
}

func TestMonitor(t *testing.T) {
	once.Do(startServer)

	path := "/metrics"
	resp, err := http.Get(uri(SockAddr, path))
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bs, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, nil, err)
	assert.False(t, len(bs) <= 0)
}

func TestMetadata(t *testing.T) {
	once.Do(startServer)
	path := "/metadata"
	resp, err := http.Get(uri(SockAddr, path))
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bs, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, nil, err)
	assert.False(t, len(bs) <= 0)
}
func TestRender(t *testing.T) {
	once.Do(startServer)

	bodyTestCase := map[string]map[string]func([]byte) error{
		"/json": {
			"GET": func(bs []byte) error {
				var data struct {
					Code int        `json:"code"`
					TTL  int        `json:"ttl"`
					Ts   *Timestamp `json:"data"`
				}
				if err := json.Unmarshal(bs, &data); err != nil {
					t.Logf("json.Unmarshal() failed: %s", err)
				}
				assert.Zero(t, data.Code)
				assert.NotZero(t, data.Ts.Now)
				assert.NotZero(t, data.TTL)
				return nil
			},
		},

		"/null-json": {
			"GET": func(bs []byte) error {
				res := map[string]interface{}{}
				if err := json.Unmarshal(bs, &res); err != nil {
					t.Logf("json.Unmarshal() failed: %s", err)
				}
				if _, ok := res["data"]; ok {
					t.Errorf("Field `data` should be omitted")
				}
				return nil
			},
		},

		"/xml": {
			"GET": func(bs []byte) error {
				var ts Timestamp
				if err := xml.Unmarshal(bs, &ts); err != nil {
					t.Logf("xml.Unmarshal() failed: %s", err)
				}
				if ts.Now == 0 {
					return fmt.Errorf("Timestamp.Now field cannot be zero: %+v", ts)
				}
				return nil
			},
		},

		"/bytes": {
			"GET": func(bs []byte) error {
				respStr := string(bs)
				oriContent := string(ByteContent)
				if respStr != oriContent {
					return fmt.Errorf("Bytes response not expected: %s != %s", respStr, oriContent)
				}
				return nil
			},
		},

		"/json-map": {
			"GET": func(bs []byte) error {
				var data struct {
					Code    int    `json:"code"`
					Tid     int    `json:"tid"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(bs, &data); err != nil {
					t.Logf("json.Unmarshal() failed: %s", err)
				}
				if data.Tid != 1 || data.Code != 0 || data.Message != "0" {
					return fmt.Errorf("Invalid respones: %+v", data)
				}
				return nil
			},
		},
	}
	contentTypeTest := map[string]string{
		"/json":      "application/json; charset=utf-8",
		"/xml":       "application/xml; charset=utf-8",
		"/bytes":     "text/html",
		"/json-map":  "application/json; charset=utf-8",
		"/null-json": "application/json; charset=utf-8",
	}

	c := &http.Client{}
	for path := range bodyTestCase {
		for method := range bodyTestCase[path] {
			validator := bodyTestCase[path][method]
			expHeader := contentTypeTest[path]
			req, err := http.NewRequest(method, uri(SockAddr, path), nil)
			if err != nil {
				t.Fatalf("Failed to build request: %s", err)
			}
			resp, err := c.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("http.%s(%s) get error code: %d", method, path, resp.StatusCode)
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll() failed: %s", err)
			}
			if ct := resp.Header.Get("Content-Type"); ct != expHeader {
				t.Errorf("Unexpected content-type header on path `%s`: expected %s got %s", path, expHeader, ct)
			}
			if err := validator(bs); err != nil {
				t.Errorf("Failed to validate request to `%s %s`: %s", method, path, err)
			}
		}
	}
}

func TestRedirect(t *testing.T) {
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header["Referer"] = []string{"http://www.bilibili.com/"}
			return nil
		},
	}
	req, err := http.NewRequest("GET", uri(SockAddr, "/redirect"), nil)
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Redirect test filed: %d", resp.StatusCode)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() failed: %s", err)
	}
	respStr := string(bs)
	oriContent := string(ByteContent)
	if respStr != oriContent {
		t.Errorf("Bytes response not expected: %s != %s", respStr, oriContent)
	}
}

func TestCORSPreflight(t *testing.T) {
	once.Do(startServer)

	origin := "ccc.bilibili.com"
	c := &http.Client{}
	req, err := http.NewRequest("OPTIONS", uri(SockAddr, "/get"), nil)
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}
	req.Header.Set("Referer", "http://www.bilibili.com/")
	req.Header.Set("Origin", origin)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("CORS preflight request status code unexpected: %d", resp.StatusCode)
	}

	preflightHeaders := map[string]string{
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "GET,POST",
		"Access-Control-Allow-Headers":     "Origin,Content-Length,Content-Type",
	}
	for k, v := range preflightHeaders {
		if hv := resp.Header.Get(k); hv != v {
			t.Errorf("Unexpected header %s: %s != %s", k, hv, v)
		}
	}

	varys := map[string]int{
		"Origin":                         0,
		"Access-Control-Request-Method":  0,
		"Access-Control-Request-Headers": 0,
	}
	reqVarys := make(map[string]int)
	rv := resp.Header["Vary"]
	for _, v := range rv {
		reqVarys[v] = 0
	}
	for v := range varys {
		if _, ok := reqVarys[v]; !ok {
			t.Errorf("%s is missed in Vary", v)
		}
	}
}

func TestCORSNormal(t *testing.T) {
	once.Do(startServer)

	origin := "ccc.bilibili.com"
	c := &http.Client{}
	req, err := http.NewRequest("GET", uri(SockAddr, "/get"), nil)
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}
	req.Header.Set("Referer", "http://www.bilibili.com/")
	req.Header.Set("Origin", origin)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	normalHeaders := map[string]string{
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "GET,POST",
	}
	for k, v := range normalHeaders {
		if hv := resp.Header.Get(k); hv != v {
			t.Errorf("Unexpected header %s: %s != %s", k, hv, v)
		}
	}
}

func TestJSONP(t *testing.T) {
	once.Do(startServer)

	origin := "ccc.bilibili.com"
	r := regexp.MustCompile(`onsuccess\((.*)\)`)
	c := &http.Client{}
	req, err := http.NewRequest("GET", uri(SockAddr, "/json?cross_domain=true&jsonp=jsonp&callback=onsuccess"), nil)
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}
	req.Header.Set("Referer", "http://www.bilibili.com/")
	req.Header.Set("Origin", origin)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() failed: %s", err)
	}
	matched := r.FindSubmatch(bs)
	if matched == nil {
		t.Errorf("Response not matched pattern: %s", r)
	}

	var data struct {
		Code int        `json:"code"`
		Ts   *Timestamp `json:"data"`
	}
	if err := json.Unmarshal(matched[1], &data); err != nil {
		t.Logf("json.Unmarshal() failed: %s", err)
	}
	if data.Code != 0 || data.Ts.Now == 0 {
		t.Errorf("Request should succeed but got blocked with code(%d) or Timestamp.Now field cannot be zero: %+v",
			data.Code, data)
	}
}

func TestCORSFailed(t *testing.T) {
	once.Do(startServer)

	c := &http.Client{}
	req, err := http.NewRequest("GET", uri(SockAddr, "/json?cross_domain=true&jsonp=jsonp&callback=onsuccess"), nil)
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}

	req.Header.Set("Referer", "http://www.bilibili2.com/")
	req.Header.Set("Origin", "ccc.bilibili2.com")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		t.Errorf("This request should be blocked, but got status %d", resp.StatusCode)
	}
}

func TestCSRF(t *testing.T) {
	once.Do(startServer)

	allowed := []string{
		"http://www.bilibili.com/",
		"http://www.biligame.com/",
		"http://www.im9.com/",
		"http://www.acg.tv/",
		"http://www.hdslb.com/",
		"http://www.bilibili.co/",
		// should match by appid
		"https://servicewechat.com/wx7564fd5313d24844/devtools/page-frame.html",
		"http://servicewechat.com/wx7564fd5313d24844/",
		"http://servicewechat.com/wx7564fd5313d24844",
		"http://servicewechat.com/wx618ca8c24bf06c33",

		// "http://bilibili.co/",
		// "http://hdslb.com/",
		// "http://acg.tv/",
		// "http://im9.com/",
		// "http://biligame.com/",
		// "http://bilibili.com/",
	}
	notAllowed := []string{
		"http://www.bilibili2.com/",
		"http://www.biligame2.com/",
		"http://www.im92.com/",
		"http://www.acg2.tv/",
		"http://www.hdslb2.com/",
		"http://servicewechat.com/",
		"http://servicewechat.com/wx7564fd5313d24842",
		"http://servicewechat.com/wx618ca8c24bf06c34",
	}

	c := &http.Client{}
	for _, r := range allowed {
		req, err := http.NewRequest("GET", uri(SockAddr, "/json"), nil)
		assert.Nil(t, err)

		req.Header.Set("Referer", r)
		resp, err := c.Do(req)
		assert.Nil(t, err)
		resp.Body.Close()
		assert.Equal(t, 200, resp.StatusCode)
	}

	for _, r := range notAllowed {
		req, err := http.NewRequest("GET", uri(SockAddr, "/json"), nil)
		assert.Nil(t, err)

		req.Header.Set("Referer", r)
		resp, err := c.Do(req)
		assert.Nil(t, err)
		resp.Body.Close()
		assert.Equal(t, 403, resp.StatusCode)
	}

	req, err := http.NewRequest("GET", uri(SockAddr, "/json"), nil)
	assert.Nil(t, err)
	resp, err := c.Do(req)
	assert.Nil(t, err)
	resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	req, err = http.NewRequest("GET", uri(SockAddr, "/json?callback=123&jsonp=jsonp"), nil)
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Nil(t, err)
	resp.Body.Close()
	assert.Equal(t, 403, resp.StatusCode)

	req, err = http.NewRequest("GET", uri(SockAddr, "/json?cross_domain=123"), nil)
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Nil(t, err)
	resp.Body.Close()
	assert.Equal(t, 403, resp.StatusCode)
}

func TestOnError(t *testing.T) {
	once.Do(startServer)
	c := &http.Client{}
	resp, err := c.Get(uri(SockAddr, "/err"))
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() failed: %s", err)
	}
	var data struct {
		Code int        `json:"code"`
		Ts   *Timestamp `json:"data"`
	}
	if err := json.Unmarshal(bs, &data); err != nil {
		t.Logf("json.Unmarshal() failed: %s", err)
	}
	if data.Code != ecode.ServerErr.Code() {
		t.Errorf("Error code is not expected: %d != %d", data.Code, ecode.ServerErr.Code())
	}
}

func TestRecovery(t *testing.T) {
	once.Do(startServer)
	c := &http.Client{}
	resp, err := c.Get(uri(SockAddr, "/panic"))
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expect status code 500 get %d, maybe recovery not working as expected", resp.StatusCode)
	}
}

func TestGlobalTimeout(t *testing.T) {
	once.Do(startServer)

	t.Run("Should timeout by default", func(t *testing.T) {
		c := &http.Client{}
		req, err := http.NewRequest("GET", uri(SockAddr, "/timeout"), nil)
		assert.Nil(t, err)
		start := time.Now()
		resp, err := c.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.InDelta(t, float64(time.Second*2), float64(time.Since(start)), float64(time.Second))
	})
	t.Run("Should timeout by delivered", func(t *testing.T) {
		c := &http.Client{}
		req, err := http.NewRequest("GET", uri(SockAddr, "/timeout"), nil)
		assert.Nil(t, err)
		td := int64(time.Second / time.Millisecond)
		req.Header.Set(_httpHeaderTimeout, strconv.FormatInt(td, 10))
		start := time.Now()
		resp, err := c.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.InDelta(t, float64(time.Second), float64(time.Since(start)), float64(time.Second))
	})
	t.Run("Should not timeout by delivered", func(t *testing.T) {
		c := &http.Client{}
		req, err := http.NewRequest("GET", uri(SockAddr, "/timeout"), nil)
		assert.Nil(t, err)
		td := int64(time.Second * 10 / time.Millisecond)
		req.Header.Set(_httpHeaderTimeout, strconv.FormatInt(td, 10))
		start := time.Now()
		resp, err := c.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.InDelta(t, float64(time.Second*2), float64(time.Since(start)), float64(time.Second))
	})
	t.Run("Should timeout by method config", func(t *testing.T) {
		c := &http.Client{}
		req, err := http.NewRequest("GET", uri(SockAddr, "/timeout-from-method-config"), nil)
		assert.Nil(t, err)
		start := time.Now()
		resp, err := c.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.InDelta(t, float64(time.Second*3), float64(time.Since(start)), float64(time.Second))
	})
}

func TestServerConfigReload(t *testing.T) {
	engine := Default()
	setupHandler(engine)

	engine.lock.RLock()
	startTm := engine.conf.Timeout
	engine.lock.RUnlock()
	toTm := startTm * 5

	wg := &sync.WaitGroup{}
	closed := make(chan struct{})
	go func() {
		if err := engine.Run(AdHocSockAddr); err != nil {
			if errors.Cause(err) == http.ErrServerClosed {
				t.Logf("Server stopped due to shutting down command")
			} else {
				t.Errorf("Failed to serve with tls: %s", err)
			}
		}
		closed <- struct{}{}
	}()
	time.Sleep(time.Second)

	clientNum := 20
	for i := 0; i < clientNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := &http.Client{}
			resp, err := c.Get(uri(AdHocSockAddr, "/get"))
			if err != nil {
				assert.Nil(t, err)
				return
			}
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}
	for i := 2; i < clientNum; i++ {
		wg.Add(1)
		conf := &ServerConfig{
			Timeout: toTm,
		}
		if i%2 == 0 {
			conf.Timeout = 0
		}
		go func() {
			defer wg.Done()
			err := engine.SetConfig(conf)
			if conf.Timeout <= 0 {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		}()
	}

	wg.Wait()

	engine.lock.RLock()
	endTm := engine.conf.Timeout
	engine.lock.RUnlock()

	assert.NotEqual(t, toTm, startTm)
	assert.Equal(t, toTm, endTm)

	if err := engine.Shutdown(context.TODO()); err != nil {
		assert.Nil(t, err)
	}
	<-closed
}

func TestGracefulShutdown(t *testing.T) {
	engine := Default()
	setupHandler(engine)

	closed := make(chan struct{})
	go func() {
		if err := engine.Run(AdHocSockAddr); err != nil {
			if errors.Cause(err) == http.ErrServerClosed {
				t.Logf("Server stopped due to shutting down command")
			} else {
				t.Errorf("Failed to serve with tls: %s", err)
			}
		}
		closed <- struct{}{}
	}()
	time.Sleep(time.Second)

	clientNum := 10
	wg := &sync.WaitGroup{}
	for i := 0; i < clientNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := &http.Client{}
			resp, err := c.Get(uri(AdHocSockAddr, "/sleep5"))
			if err != nil {
				t.Error(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: %d", resp.StatusCode)
				return
			}
			t.Logf("Request finished at: %v", time.Now())
		}()
	}
	time.Sleep(time.Second)

	t.Logf("Invoke Shutdown method at: %v", time.Now())
	if err := engine.Shutdown(context.TODO()); err != nil {
		t.Fatalf("Failed to shutdown engine: %s", err)
	}
	wg.Wait()
	<-closed
}

func TestProtobuf(t *testing.T) {
	once.Do(startServer)
	c := &http.Client{}

	t.Run("On-Success", func(t *testing.T) {
		req, err := http.NewRequest("GET", uri(SockAddr, "/pb"), nil)
		assert.Nil(t, err)

		resp, err := c.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/x-protobuf")

		bs, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		d := &render.PB{}
		err = proto.Unmarshal(bs, d)
		assert.Nil(t, err)
		assert.Equal(t, int(d.Code), 0)
		assert.NotNil(t, d.Data)

		tt := &tests.Time{}
		err = proto.Unmarshal(d.Data.Value, tt)
		assert.Nil(t, err)
		assert.NotZero(t, tt.Now)
	})
	t.Run("On-Error", func(t *testing.T) {
		req, err := http.NewRequest("GET", uri(SockAddr, "/pb-error"), nil)
		assert.Nil(t, err)

		resp, err := c.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/x-protobuf")

		bs, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		d := &render.PB{}
		err = proto.Unmarshal(bs, d)
		assert.Nil(t, err)
		assert.Equal(t, int(d.Code), ecode.RequestErr.Code())
	})

	t.Run("On-NilData-NilErr", func(t *testing.T) {
		req, err := http.NewRequest("GET", uri(SockAddr, "/pb-nildata-nilerr"), nil)
		assert.Nil(t, err)

		resp, err := c.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/x-protobuf")

		bs, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		d := &render.PB{}
		err = proto.Unmarshal(bs, d)
		assert.Nil(t, err)
		assert.Equal(t, int(d.Code), 0)
	})

	t.Run("On-Data-Err", func(t *testing.T) {
		req, err := http.NewRequest("GET", uri(SockAddr, "/pb-data-err"), nil)
		assert.Nil(t, err)

		resp, err := c.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/x-protobuf")

		bs, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		d := &render.PB{}
		err = proto.Unmarshal(bs, d)
		assert.Nil(t, err)
		assert.Equal(t, int(d.Code), ecode.RequestErr.Code())
	})
}

func TestBind(t *testing.T) {
	once.Do(startServer)

	c := &http.Client{}
	req, err := http.NewRequest("GET", uri(SockAddr, "/bind?mids=2,3,4&title=hello&content=world&cid=8"), nil)
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	assert.Equal(t, resp.StatusCode, 200)

	bs, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var p struct {
		Code int      `json:"code"`
		Data *Archive `json:"data"`
	}
	if err := json.Unmarshal(bs, &p); err != nil {
		t.Fatalf("Failed to json.Unmarshal: %v", resp.StatusCode)
	}
	assert.Equal(t, p.Code, 0)
	assert.Equal(t, []int64{2, 3, 4}, p.Data.Mids)
	assert.Equal(t, "hello", p.Data.Title)
	assert.Equal(t, "world", p.Data.Content)
	assert.Equal(t, 8, p.Data.Cid)
}

func TestBindWith(t *testing.T) {
	once.Do(startServer)

	a := &Archive{
		Mids:    []int64{2, 3, 4},
		Title:   "hello",
		Content: "world",
		Cid:     8,
	}
	d, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Failed to json.Marshal: %v", err)
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", uri(SockAddr, "/bindwith"), bytes.NewBuffer(d))
	if err != nil {
		t.Fatalf("Failed to build request: %s", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %s", err)
	}
	assert.Equal(t, resp.StatusCode, 200)

	bs, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	var p struct {
		Code int      `json:"code"`
		Data *Archive `json:"data"`
	}
	if err := json.Unmarshal(bs, &p); err != nil {
		t.Fatalf("Failed to json.Unmarshal: %v", resp.StatusCode)
	}
	assert.Equal(t, p.Code, 0)
	assert.Equal(t, a.Mids, p.Data.Mids)
	assert.Equal(t, a.Title, p.Data.Title)
	assert.Equal(t, a.Content, p.Data.Content)
	assert.Equal(t, a.Cid, p.Data.Cid)
}

func TestMethodNotAllowed(t *testing.T) {
	once.Do(startServer)

	resp, err := http.Get(uri(SockAddr, "/post"))
	assert.NoError(t, err)
	bs, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, 405, resp.StatusCode)
	assert.Equal(t, http.StatusText(405), strings.TrimSpace(string(bs)))
}

func TestDevice(t *testing.T) {
	testDevice(t, "/device")
}

func TestDeviceMeta(t *testing.T) {
	testDevice(t, "/device-from-meta")
}

func testDevice(t *testing.T, path string) {
	once.Do(startServer)

	req, err := http.NewRequest("GET", uri(SockAddr, path), nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	req.Header.Set("Buvid", "9346b9ca66dangerous4764eede8bb2")
	req.AddCookie(&http.Cookie{Name: "buvid3", Value: "25213BD4-841C-449F-8BBF-B96B58A8Fdangerousinfoc"})
	req.AddCookie(&http.Cookie{Name: "sid", Value: "70***dpi"})
	query := req.URL.Query()
	query.Set("build", "6280")
	query.Set("device", "phone")
	query.Set("mobi_app", "iphone")
	query.Set("platform", "ios")
	query.Set("channel", "appstore")
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	rdev := new(DevResp)
	bs, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(bs, rdev))
	dev := rdev.Data
	assert.Equal(t, "70***dpi", dev.Sid)
	assert.Equal(t, int64(6280), dev.Build)
	assert.Equal(t, "phone", dev.Device)
	assert.Equal(t, "iphone", dev.RawMobiApp)
	assert.Equal(t, "ios", dev.RawPlatform)
	assert.Equal(t, "appstore", dev.Channel)
	assert.Equal(t, "9346b9ca66dangerous4764eede8bb2", dev.Buvid)
	assert.Equal(t, "25213BD4-841C-449F-8BBF-B96B58A8Fdangerousinfoc", dev.Buvid3)
	assert.Equal(t, PlatIPhone, dev.Plat())
	assert.Equal(t, false, dev.IsAndroid())
	assert.Equal(t, true, dev.IsIOS())
	assert.Equal(t, false, dev.IsOverseas())
	assert.Equal(t, false, dev.InvalidChannel("*"))
	assert.Equal(t, false, dev.InvalidChannel("appstore"))
	assert.Equal(t, "iphone", dev.MobiApp())
	assert.Equal(t, "iphone", dev.MobiAPPBuleChange())

	dev.RawPlatform = "android"
	dev.RawMobiApp = "android"
	dev.Channel = "test.channel"
	assert.Equal(t, true, dev.InvalidChannel("test.channel2"))
}

func TestRemoteIPFromContext(t *testing.T) {
	once.Do(startServer)

	ip := "192.168.22.33"
	req, err := http.NewRequest("GET", uri(SockAddr, "/remote-ip"), nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	req.Header.Set(_httpHeaderRemoteIP, ip)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var response struct {
		Data string `json:"data"`
	}
	err = json.Unmarshal(bs, &response)
	assert.NoError(t, err)
	assert.Equal(t, ip, response.Data)
}

func TestNewServerMiddleware(t *testing.T) {
	e := DefaultServer(nil)
	// should contain 4 default handlers
	assert.Len(t, e.RouterGroup.Handlers, 5)
}

func TestMethodConfig(t *testing.T) {
	e := New()
	tg := e.Group("/timeout").SetMethodConfig(&MethodConfig{Timeout: xtime.Duration(time.Second * 10)})
	tg.GET("/default", func(ctx *Context) {})

	e.GET("/timeout/5s", func(ctx *Context) {})
	e.SetMethodConfig("/timeout/5s", &MethodConfig{Timeout: xtime.Duration(time.Second * 5)})

	pc := e.methodConfig("/timeout/default")
	assert.NotNil(t, pc)
	assert.Equal(t, xtime.Duration(time.Second*10), pc.Timeout)

	pc = e.methodConfig("/timeout/5s")
	assert.NotNil(t, pc)
	assert.Equal(t, xtime.Duration(time.Second*5), pc.Timeout)
}

func BenchmarkServer(b *testing.B) {
	startServer()
	defer shutdown()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get(uri(SockAddr, "/bench"))
			if err != nil {
				b.Errorf("HTTPServ: get error(%v)", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				b.Errorf("HTTPServ: get error status code:%d", resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}

func TestServerWithMirror(t *testing.T) {

	var response struct {
		Data string `json:"data"`
	}

	//sonce.Do(startMirrorServer)
	client := &http.Client{}

	reqest, _ := http.NewRequest("GET", "http://localhost:18090/mirrortest", nil)
	reqest.Header.Add("x1-bilispy-mirror", "1")

	resp, err := client.Do(reqest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	bs, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bs, &response)
	assert.NoError(t, err)
	assert.Equal(t, "true", response.Data)
	resp.Body.Close()

	reqest, _ = http.NewRequest("GET", "http://localhost:18090/mirrortest", nil)
	reqest.Header.Add("x1-bilispy-mirror", "0")
	resp, err = client.Do(reqest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	bs, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bs, &response)
	assert.NoError(t, err)
	assert.Equal(t, "false", response.Data)
	resp.Body.Close()

	reqest, _ = http.NewRequest("GET", "http://localhost:18090/mirrortest", nil)
	reqest.Header.Add("x1-bilispy-mirror", "xxxxxxxxxx")
	resp, err = client.Do(reqest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	bs, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bs, &response)
	assert.NoError(t, err)
	assert.Equal(t, "false", response.Data)
	resp.Body.Close()

	reqest, _ = http.NewRequest("GET", "http://localhost:18090/mirrortest", nil)
	resp, err = client.Do(reqest)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	bs, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bs, &response)
	assert.NoError(t, err)
	assert.Equal(t, "false", response.Data)
	resp.Body.Close()

}

func TestEngineInject(t *testing.T) {
	once.Do(startServer)

	client := &http.Client{}
	resp, err := client.Get(uri(SockAddr, "/inject/index"))
	assert.NoError(t, err)
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "index-injected", string(bs))
}

func TestBiliStatusCode(t *testing.T) {
	once.Do(startServer)

	tests := map[string]int{
		"/group/test/json-status":     ecode.MemberBlocked.Code(),
		"/group/test/xml-status":      ecode.MemberBlocked.Code(),
		"/group/test/proto-status":    ecode.MemberBlocked.Code(),
		"/group/test/json-map-status": ecode.MemberBlocked.Code(),
	}
	client := &http.Client{}
	for path, code := range tests {
		resp, err := client.Get(uri(SockAddr, path))
		assert.NoError(t, err)
		assert.Equal(t, resp.Header.Get("bili-status-code"), strconv.FormatInt(int64(code), 10))
	}
}
