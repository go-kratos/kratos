package polaris

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
)

func TestRouter(t *testing.T) {
	token, err := getToken("http://127.0.0.1:8090")
	if err != nil {
		t.Fatal(err)
	}
	data := `
[
   {
       "name":"kratos",
       "enable":false,
       "description":"123",
       "priority":2,
       "routing_config":{
           "@type":"type.googleapis.com/v2.RuleRoutingConfig",
           "sources":[
               {
                   "service":"*",
                   "namespace":"*",
                   "arguments":[

                   ]
               }
           ],
           "destinations":[
               {
                   "labels":{
                       "az":{
                           "value":"1",
                           "value_type":"TEXT",
                           "type":"EXACT"
                       }
                   },
                   "weight":100,
                   "priority":1,
                   "isolate":false,
                   "name":"实例分组1",
                   "namespace":"default",
                   "service":"test-ut"
               }
           ]
       }
   }
]
`
	res, err := makeJSONRequest("http://127.0.0.1:8090/naming/v2/routings", data, http.MethodPost, map[string]string{
		"X-Polaris-Token": token,
	})
	if err != nil {
		t.Fatal(err)
	}

	resJSON := struct {
		Code      int `json:"code"`
		Responses []struct {
			Data struct {
				ID string `json:"id"`
			}
		} `json:"responses"`
	}{}

	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		t.Fatal(err, string(res))
	}
	if resJSON.Code != 200000 {
		t.Fatal("create failed", string(res))
	}

	// enable router
	enableData := fmt.Sprintf(`[{"id":"%s","enable":true}]`, resJSON.Responses[0].Data.ID)
	res, err = makeJSONRequest("http://127.0.0.1:8090/naming/v2/routings/enable", enableData, http.MethodPut, map[string]string{
		"X-Polaris-Token": token,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(res, &resJSON)
	if err != nil {
		t.Fatal(err)
	}
	if resJSON.Code != 200000 {
		t.Fatal("enable failed", string(res))
	}

	t.Cleanup(func() {
		enableData := fmt.Sprintf(`[{"id":"%s"}]`, resJSON.Responses[0].Data.ID)
		res, err = makeJSONRequest("http://127.0.0.1:8090/naming/v2/routings/delete", enableData, http.MethodPost, map[string]string{
			"X-Polaris-Token": token,
		})
		resJSON := &commonRes{}
		err = json.Unmarshal(res, resJSON)
		if err != nil {
			t.Fatal(err, string(res))
		}
		if resJSON.Code != 200000 {
			t.Fatal("delete failed", string(res))
		}
	})

	sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	if err != nil {
		t.Fatal(err)
	}

	p := New(sdk)

	r := p.Registry(
		WithRegistryTimeout(time.Second),
		WithRegistryHealthy(true),
		WithRegistryIsolate(false),
		WithRegistryRetryCount(0),
		WithRegistryWeight(100),
		WithRegistryTTL(10),
	)

	ins := &registry.ServiceInstance{
		ID:      "kratos",
		Name:    "kratos",
		Version: "v1.0.0",
		Endpoints: []string{
			"grpc://127.0.0.1:8080",
			"http://127.0.0.1:9090",
		},
	}

	err = r.Register(context.Background(), ins)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	nodes := []selector.Node{
		selector.NewNode("grpc", "127.0.0.1:9000", &registry.ServiceInstance{
			ID:        "123",
			Name:      "test-ut",
			Version:   "v1.0.0",
			Metadata:  map[string]string{"weight": "100", "az": "1"},
			Endpoints: []string{"grpc://127.0.0.1:9000"},
		}),
		selector.NewNode("grpc", "127.0.0.2:9000", &registry.ServiceInstance{
			ID:        "123",
			Name:      "test-ut",
			Version:   "v1.0.0",
			Metadata:  map[string]string{"weight": "100", "az": "2"},
			Endpoints: []string{"grpc://127.0.0.2:9000"},
		}),
		selector.NewNode("grpc", "127.0.0.3:9000", &registry.ServiceInstance{
			ID:        "123",
			Name:      "test-ut",
			Version:   "v1.0.0",
			Metadata:  map[string]string{"weight": "100", "az": "1"},
			Endpoints: []string{"grpc://127.0.0.3:9000"},
		}),
	}

	f := p.NodeFilter()
	ctx := kratos.NewContext(context.Background(), &mockApp{})
	n := f(ctx, nodes)
	for _, node := range n {
		if node.Metadata()["az"] != "1" {
			t.Fatal("node filter result wrong")
		}
		t.Log(node)
	}
	if len(n) != 2 {
		t.Fatal("node filter result wrong")
	}
}

type mockApp struct{}

func (m mockApp) ID() string {
	return "1"
}

func (m mockApp) Name() string {
	return "kratos"
}

func (m mockApp) Version() string {
	return "v2.0.0"
}

func (m mockApp) Metadata() map[string]string {
	return map[string]string{}
}

func (m mockApp) Endpoint() []string {
	return []string{"grpc://123.123.123.123:9090"}
}
