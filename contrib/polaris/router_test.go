package polaris

import (
	"testing"
)

func TestRouter(t *testing.T) {
	//sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//p := New(sdk)
	//nodes := []selector.Node{
	//	selector.NewNode(
	//		"grpc",
	//		"192.168.1.1:1234",
	//		&registry.ServiceInstance{
	//			ID:        "1",
	//			Name:      "test",
	//			Version:   "v2.0.0",
	//			Metadata:  map[string]string{"az": "1"},
	//			Endpoints: []string{"grpc://192.168.1.2:1234"},
	//		},
	//	),
	//	selector.NewNode(
	//		"grpc",
	//		"192.168.1.2:1234",
	//		&registry.ServiceInstance{
	//			ID:        "2",
	//			Name:      "test",
	//			Version:   "v2.0.0",
	//			Metadata:  map[string]string{"az": "1"},
	//			Endpoints: []string{"grpc://192.168.1.2:1234"},
	//		},
	//	),
	//	selector.NewNode(
	//		"grpc",
	//		"192.168.1.3:1234",
	//		&registry.ServiceInstance{
	//			ID:        "3",
	//			Name:      "test",
	//			Version:   "v2.0.0",
	//			Metadata:  map[string]string{"az": "2"},
	//			Endpoints: []string{"grpc://192.168.1.2:1234"},
	//		},
	//	),
	//}
	//f, err := p.NodeFilter("test", "default")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//ctx := kratos.NewContext(context.Background(), &mockApp{})
	//n := f(ctx, nodes)
	//
	//for _, node := range n {
	//	t.Log(node)
	//}
}

func BenchmarkPolaris_NodeFilter(b *testing.B) {
	//sdk, err := polaris.NewSDKContextByAddress("127.0.0.1:8091")
	//if err != nil {
	//	b.Fatal(err)
	//}
	//p := New(sdk)
	//
	//var instances []*registry.ServiceInstance
	//for i := 0; i < 10; i++ {
	//	s := &registry.ServiceInstance{
	//		ID:        strconv.Itoa(i),
	//		Name:      "test",
	//		Version:   "v2.0.0",
	//		Endpoints: []string{fmt.Sprintf("grpc://192.168.0.1:%d", i+1)},
	//		Metadata:  map[string]string{},
	//	}
	//	instances = append(instances, s)
	//}
	//
	//time.Sleep(3 * time.Second)
	//
	//var nodes []selector.Node
	//
	//for _, instance := range instances {
	//	nodes = append(nodes, selector.NewNode(
	//		"grpc",
	//		instance.Endpoints[0],
	//		instance,
	//	))
	//}
	//
	//req := &polaris.ProcessRoutersRequest{
	//	ProcessRoutersRequest: model.ProcessRoutersRequest{
	//		SourceService: model.ServiceInfo{
	//			Service:   "kratos",
	//			Namespace: "default",
	//		},
	//		DstInstances: m,
	//	},
	//}
	//b.ResetTimer()
	//for i := 0; i < b.N; i++ {
	//	_, err := p.router.ProcessRouters(req)
	//	if err != nil {
	//		b.Fatal(err)
	//	}
	//}
}

type mockApp struct {
}

func (m mockApp) ID() string {
	return "1"
}

func (m mockApp) Name() string {
	//TODO implement me
	return "kratos"
}

func (m mockApp) Version() string {
	//TODO implement me
	return "v2.0.0"
}

func (m mockApp) Metadata() map[string]string {
	//TODO implement me
	return map[string]string{}
}

func (m mockApp) Endpoint() []string {
	//TODO implement me
	return []string{"grpc://123.123.123.123:9090"}
}
