package eureka

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	statusUp            = "UP"
	statusDown          = "DOWN"
	statusOutOfServeice = "OUT_OF_SERVICE"
	heartbeatRetry      = 3
	maxIdleConns        = 100
	heartbeatTime       = 10
	httpTimeout         = 3
	refreshTime         = 30
)

type Endpoint struct {
	InstanceID     string
	IP             string
	AppID          string
	Port           int
	SecurePort     int
	HomePageURL    string
	StatusPageURL  string
	HealthCheckURL string
	MetaData       map[string]string
}

// ApplicationsRootResponse for /eureka/apps
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

type RequestInstance struct {
	Instance Instance `json:"instance"`
}

type Instance struct {
	InstanceID     string            `json:"instanceId"`
	HostName       string            `json:"hostName"`
	Port           Port              `json:"port"`
	App            string            `json:"app"`
	IPAddr         string            `json:"ipAddr"`
	VipAddress     string            `json:"vipAddress"`
	Status         string            `json:"status"`
	SecurePort     Port              `json:"securePort"`
	HomePageURL    string            `json:"homePageUrl"`
	StatusPageURL  string            `json:"statusPageUrl"`
	HealthCheckURL string            `json:"healthCheckUrl"`
	DataCenterInfo DataCenterInfo    `json:"dataCenterInfo"`
	Metadata       map[string]string `json:"metadata"`
}

type Port struct {
	Port    int    `json:"$"`
	Enabled string `json:"@enabled"`
}

type DataCenterInfo struct {
	Name  string `json:"name"`
	Class string `json:"@class"`
}

var _ APIInterface = (*Client)(nil)

type APIInterface interface {
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

type ClientOption func(e *Client)

func WithMaxRetry(maxRetry int) ClientOption {
	return func(e *Client) { e.maxRetry = maxRetry }
}

func WithHeartbeatInterval(interval time.Duration) ClientOption {
	return func(e *Client) {
		e.heartbeatInterval = interval
	}
}

func WithClientContext(ctx context.Context) ClientOption {
	return func(e *Client) { e.ctx = ctx }
}

func WithNamespace(path string) ClientOption {
	return func(e *Client) { e.eurekaPath = path }
}

type Client struct {
	ctx               context.Context
	urls              []string
	eurekaPath        string
	maxRetry          int
	heartbeatInterval time.Duration
	client            *http.Client
	keepalive         map[string]chan struct{}
	lock              sync.Mutex
}

func NewClient(urls []string, opts ...ClientOption) *Client {
	tr := &http.Transport{
		MaxIdleConns: maxIdleConns,
	}

	e := &Client{
		ctx:               context.Background(),
		urls:              urls,
		eurekaPath:        "eureka/v2",
		maxRetry:          len(urls),
		heartbeatInterval: time.Second * heartbeatTime,
		client:            &http.Client{Transport: tr, Timeout: time.Second * httpTimeout},
		keepalive:         make(map[string]chan struct{}),
	}

	for _, o := range opts {
		o(e)
	}

	return e
}

func (e *Client) FetchApps(ctx context.Context) []Application {
	var m ApplicationsRootResponse
	if err := e.do(ctx, "GET", []string{"apps"}, nil, &m); err != nil {
		return nil
	}

	return m.Applications
}

func (e *Client) FetchAppInstances(ctx context.Context, appID string) (m Application, err error) {
	err = e.do(ctx, "GET", []string{"apps", appID}, nil, &m)
	return
}

func (e *Client) FetchAppUpInstances(ctx context.Context, appID string) []Instance {
	app, err := e.FetchAppInstances(ctx, appID)
	if err != nil {
		return nil
	}
	return e.filterUp(app)
}

func (e *Client) FetchAppInstance(ctx context.Context, appID string, instanceID string) (m Instance, err error) {
	err = e.do(ctx, "GET", []string{"apps", appID, instanceID}, nil, &m)
	return
}

func (e *Client) FetchInstance(ctx context.Context, instanceID string) (m Instance, err error) {
	err = e.do(ctx, "GET", []string{"instances", instanceID}, nil, &m)
	return
}

func (e *Client) Out(ctx context.Context, appID, instanceID string) error {
	return e.do(ctx, "PUT", []string{"apps", appID, instanceID, fmt.Sprintf("status?value=%s", statusOutOfServeice)}, nil, nil)
}

func (e *Client) Down(ctx context.Context, appID, instanceID string) error {
	return e.do(ctx, "PUT", []string{"apps", appID, instanceID, fmt.Sprintf("status?value=%s", statusDown)}, nil, nil)
}

func (e *Client) FetchAllUpInstances(ctx context.Context) []Instance {
	return e.filterUp(e.FetchApps(ctx)...)
}

func (e *Client) Register(ctx context.Context, ep Endpoint) error {
	return e.registerEndpoint(ctx, ep)
}

func (e *Client) Deregister(ctx context.Context, appID, instanceID string) error {
	if err := e.do(ctx, "DELETE", []string{"apps", appID, instanceID}, nil, nil); err != nil {
		return err
	}
	go e.cancelHeartbeat(appID)
	return nil
}

func (e *Client) registerEndpoint(ctx context.Context, ep Endpoint) error {
	instance := RequestInstance{
		Instance: Instance{
			InstanceID: ep.InstanceID,
			HostName:   ep.AppID,
			Port: Port{
				Port:    ep.Port,
				Enabled: "true",
			},
			App:        ep.AppID,
			IPAddr:     ep.IP,
			VipAddress: ep.AppID,
			Status:     statusUp,
			SecurePort: Port{
				Port:    ep.SecurePort,
				Enabled: "false",
			},
			HomePageURL:    ep.HomePageURL,
			StatusPageURL:  ep.StatusPageURL,
			HealthCheckURL: ep.HealthCheckURL,
			DataCenterInfo: DataCenterInfo{
				Name:  "MyOwn",
				Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			},
			Metadata: ep.MetaData,
		},
	}

	body, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	return e.do(ctx, "POST", []string{"apps", ep.AppID}, bytes.NewReader(body), nil)
}

func (e *Client) Heartbeat(ep Endpoint) {
	e.lock.Lock()
	e.keepalive[ep.AppID] = make(chan struct{})
	e.lock.Unlock()

	ticker := time.NewTicker(e.heartbeatInterval)
	defer ticker.Stop()
	retryCount := 0
	for {
		select {
		case <-e.ctx.Done():
			return
		case <-e.keepalive[ep.AppID]:
			return
		case <-ticker.C:
			if err := e.do(e.ctx, "PUT", []string{"apps", ep.AppID, ep.InstanceID}, nil, nil); err != nil {
				if retryCount++; retryCount > heartbeatRetry {
					_ = e.registerEndpoint(e.ctx, ep)
					retryCount = 0
				}
			}
		}
	}
}

func (e *Client) cancelHeartbeat(appID string) {
	defer e.lock.Unlock()
	e.lock.Lock()
	if ch, ok := e.keepalive[appID]; ok {
		ch <- struct{}{}
	}
}

func (e *Client) filterUp(apps ...Application) (res []Instance) {
	for _, app := range apps {
		for _, ins := range app.Instance {
			if ins.Status == statusUp {
				res = append(res, ins)
			}
		}
	}

	return
}

func (e *Client) pickServer(currentTimes int) string {
	return e.urls[currentTimes%e.maxRetry]
}

func (e *Client) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(e.urls), func(i, j int) {
		e.urls[i], e.urls[j] = e.urls[j], e.urls[i]
	})
}

func (e *Client) buildAPI(currentTimes int, params ...string) string {
	if currentTimes == 0 {
		e.shuffle()
	}
	server := e.pickServer(currentTimes)
	params = append([]string{server, e.eurekaPath}, params...)
	return strings.Join(params, "/")
}

func (e *Client) do(ctx context.Context, method string, params []string, input io.Reader, output interface{}) error {
	for i := 0; i < e.maxRetry; i++ {
		request, err := http.NewRequest(method, e.buildAPI(i, params...), input)
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
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		if output != nil && resp.StatusCode/100 == 2 {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			if err = json.Unmarshal(data, output); err != nil {
				return err
			}
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("response Error %d", resp.StatusCode)
		}

		return nil
	}
	return fmt.Errorf("retry after %d times", e.maxRetry)
}
