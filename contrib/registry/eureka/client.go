package eureka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"
)

const STATUS_UP = "UP"
const STATUS_DOWN = "DOWN"
const STATUS_OUT_OF_SERVICE = "OUT_OF_SERVICE"
const EUREKA_PATH = "eureka"

const tpl = `{
  "instance": {
	"instanceId": "{{ .InstanceID }}",
    "hostName": "{{ .IP }}",
    "app": "{{ .AppID }}",
    "ipAddr": "{{ .IP }}",
    "vipAddress": "{{ .AppID }}",
	"overriddenstatus": "UNKNOWN",
    "status":"UP",
    "port": {
      "$":{{ .Port }},
      "@enabled": true
    },
    "securePort": {
      "$":{{ .SecurePort }},
      "@enabled": false
    },
    "homePageUrl" : "http://{{ .IP }}:{{ .Port }}{{ .HomePageUrl }}",
    "statusPageUrl": "http://{{ .IP }}:{{ .Port }}{{ .StatusPageUrl }}",
    "healthCheckUrl": "http://{{ .IP }}:{{ .Port }}{{ .HealthCheckUrl }}",
    "dataCenterInfo" : {
      "@class":"com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
      "name": "MyOwn"
    },
    "metadata": {
			{{- range $key, $value := .MetaData }}
			"{{ $key }}": "{{ $value }}",
	 		{{- end }}
			"agent" : "go-eureka-client"
    }
  }
}`

type Endpoint struct {
	InstanceID     string
	IP             string
	AppID          string
	Port           string
	SecurePort     string
	HomePageUrl    string
	StatusPageUrl  string
	HealthCheckUrl string
	MetaData       map[string]string
}

// Response for /eureka/apps
type ApplicationsRootResponse struct {
	ApplicationsResponse `json:"applications"`
}

type ApplicationsResponse struct {
	Version      string        `json:"versions__delta"`
	AppsHashcode string        `json:"apps__hashcode"`
	Applications []Application `json:"application"`
}

type Application struct {
	Name     string     `json:"name"`
	Instance []Instance `json:"instance"`
}

type Instance struct {
	InstanceId         string            `json:"instanceId"`
	HostName           string            `json:"hostName"`
	Port               Port              `json:"port"`
	App                string            `json:"app"`
	IpAddr             string            `json:"ipAddr"`
	Status             string            `json:"status"`
	SecurePort         SecurePort        `json:"securePort"`
	LastDirtyTimestamp string            `json:"lastDirtyTimestamp"`
	Metadata           map[string]string `json:"metadata"`
}

type Port struct {
	Port int `json:"$"`
}

type SecurePort struct {
	SecurePort int `json:"$"`
}

var _ EurekaApi = new(EurekaClient)

type EurekaApi interface {
	Register(ctx context.Context, ep Endpoint) error
	Deregister(ctx context.Context, appID, instanceID string) error
	Heartbeat(ep Endpoint)
	FetchApps(ctx context.Context) []Application
	FetchAllUpInstances(ctx context.Context) []Instance
	FetchAppInstances(ctx context.Context, appID string) (m Application, err error)
	FetchAppUpInstances(ctx context.Context, appID string) []Instance
	FetchAppInstance(ctx context.Context, appID string, instanceID string) (m Instance, err error)
	FetchInstance(ctx context.Context, instanceID string) (m Instance, err error)
	Out(ctx context.Context, appID, instanceID string) error
	Down(ctx context.Context, appID, instanceID string) error
}

type EurekaOption func(e *EurekaClient)

func WithMaxRetry(max_retry int) EurekaOption {
	return func(e *EurekaClient) { e.max_retry = max_retry }
}

func WithHeartbeatInterval(interval string) EurekaOption {
	return func(e *EurekaClient) {
		if duration := e.toDuration(interval); duration > 0 {
			e.heartbeatInterval = duration
		}
	}
}

func WithCtx(ctx context.Context) EurekaOption {
	return func(e *EurekaClient) { e.ctx = ctx }
}

type EurekaClient struct {
	ctx               context.Context
	urls              []string
	max_retry         int
	heartbeatInterval time.Duration
	client            *http.Client
	keepalive         map[string]chan struct{}
	lock              sync.Mutex
}

func NewEurekaClient(urls []string, opts ...EurekaOption) *EurekaClient {
	tr := &http.Transport{
		MaxIdleConns: 100,
	}

	e := &EurekaClient{
		ctx:               context.Background(),
		urls:              urls,
		max_retry:         len(urls),
		heartbeatInterval: time.Second * 10,
		client:            &http.Client{Transport: tr, Timeout: time.Second * 3},
		keepalive:         make(map[string]chan struct{}),
	}

	for _, o := range opts {
		o(e)
	}

	return e
}

func (e *EurekaClient) FetchApps(ctx context.Context) []Application {
	var m ApplicationsRootResponse
	if err := e.do(ctx, "GET", []string{"apps"}, nil, &m); err != nil {
		return nil
	}

	return m.Applications
}

func (e *EurekaClient) FetchAppInstances(ctx context.Context, appID string) (m Application, err error) {
	err = e.do(ctx, "GET", []string{"apps", appID}, nil, &m)
	return
}

func (e *EurekaClient) FetchAppUpInstances(ctx context.Context, appID string) []Instance {
	app, err := e.FetchAppInstances(ctx, appID)
	if err != nil {
		return nil
	}
	return e.filterUp(app)
}

func (e *EurekaClient) FetchAppInstance(ctx context.Context, appID string, instanceID string) (m Instance, err error) {
	e.do(ctx, "GET", []string{"apps", appID, instanceID}, nil, &m)
	return
}

func (e *EurekaClient) FetchInstance(ctx context.Context, instanceID string) (m Instance, err error) {
	e.do(ctx, "GET", []string{"instances", instanceID}, nil, &m)
	return
}

func (e *EurekaClient) Out(ctx context.Context, appID, instanceID string) error {
	return e.do(ctx, "PUT", []string{"apps", appID, instanceID, fmt.Sprintf("status?value=%s", STATUS_OUT_OF_SERVICE)}, nil, nil)
}

func (e *EurekaClient) Down(ctx context.Context, appID, instanceID string) error {
	return e.do(ctx, "PUT", []string{"apps", appID, instanceID, fmt.Sprintf("status?value=%s", STATUS_DOWN)}, nil, nil)
}

func (e *EurekaClient) FetchAllUpInstances(ctx context.Context) []Instance {
	return e.filterUp(e.FetchApps(ctx)...)
}

func (e *EurekaClient) Register(ctx context.Context, ep Endpoint) error {
	return e.registerEndpoint(ctx, ep)
}

func (e *EurekaClient) Deregister(ctx context.Context, appID, instanceID string) error {
	if err := e.do(ctx, "DELETE", []string{"apps", appID, instanceID}, nil, nil); err != nil {
		return err
	}
	go e.cancelHeartbeat(appID)
	return nil
}

func (e *EurekaClient) registerEndpoint(ctx context.Context, ep Endpoint) error {
	body, err := e.parseTpl(ep)
	if err != nil {
		return err
	}
	return e.do(ctx, "POST", []string{"apps", ep.AppID}, bytes.NewReader(body), nil)
}

func (e *EurekaClient) Heartbeat(ep Endpoint) {
	e.lock.Lock()
	e.keepalive[ep.AppID] = make(chan struct{})
	e.lock.Unlock()

	ticker := time.NewTicker(e.heartbeatInterval)
	defer ticker.Stop()
	retryCount := 0
	for {
		select {
		case <-e.ctx.Done():
		case <-e.keepalive[ep.AppID]:
			return
		case <-ticker.C:
			if err := e.do(e.ctx, "PUT", []string{"apps", ep.AppID, ep.InstanceID}, nil, nil); err != nil {
				if retryCount++; retryCount > 3 {
					e.registerEndpoint(e.ctx, ep)
					retryCount = 0
				}
			}
		}
	}
}

func (e *EurekaClient) cancelHeartbeat(appID string) {
	defer e.lock.Unlock()
	e.lock.Lock()
	if ch, ok := e.keepalive[appID]; ok {
		ch <- struct{}{}
	}
}

func (e *EurekaClient) filterUp(apps ...Application) (res []Instance) {
	for _, app := range apps {
		for _, ins := range app.Instance {
			if ins.Status == STATUS_UP {
				res = append(res, ins)
			}
		}
	}

	return
}

func (e *EurekaClient) pickServer(currentTimes int) string {
	return e.urls[currentTimes%e.max_retry]
}

func (e *EurekaClient) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(e.urls), func(i, j int) {
		e.urls[i], e.urls[j] = e.urls[j], e.urls[i]
	})
}

func (e *EurekaClient) buildApi(currentTimes int, params ...string) string {
	if currentTimes == 0 {
		e.shuffle()
	}
	server := e.pickServer(currentTimes)
	params = append([]string{server, EUREKA_PATH}, params...)
	return strings.Join(params, "/")
}

func (e *EurekaClient) parseTpl(endpoint Endpoint) ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New("eureka").Parse(tpl)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Execute(buf, endpoint); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *EurekaClient) toAppID(serverName string) string {
	return strings.ToUpper(serverName)
}

func (e *EurekaClient) toDuration(interval string) time.Duration {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		return 0
	}
	return duration
}

func (e *EurekaClient) do(ctx context.Context, method string, params []string, input io.Reader, output interface{}) error {
	for i := 0; i < e.max_retry; i++ {
		request, err := http.NewRequest(method, e.buildApi(i, params...), input)
		if err != nil {
			return err
		}
		request = request.WithContext(ctx)
		request.Header.Add("User-Agent", "go-eureka-client")
		request.Header.Add("Accept", "application/json;charset=UTF-8")
		request.Header.Add("Content-Type", "application/json;charset=UTF-8")
		resp, err := e.client.Do(request)
		if err != nil {
			continue
		}
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		if output != nil && resp.StatusCode/100 == 2 {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if err = json.Unmarshal(data, output); err != nil {
				return err
			}
		}

		if resp.StatusCode > 299 {
			return errors.New(fmt.Sprintf("Response Error %d", resp.StatusCode))
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Retry after %d times.", e.max_retry))
}
