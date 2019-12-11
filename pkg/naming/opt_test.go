package naming

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Subset(t *testing.T) {
	var inss1 []*Instance
	for i := 0; i < 200; i++ {
		ins := &Instance{
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c1"},
		}
		inss1 = append(inss1, ins)
	}
	var opt BuildOptions
	s := Subset(50)
	s.Apply(&opt)
	sub1 := opt.Subset(inss1, opt.SubsetSize)
	if len(sub1) != 50 {
		t.Fatalf("subset size should be 50")
	}
	sub2 := opt.Subset(inss1, opt.SubsetSize)
	if !reflect.DeepEqual(sub1, sub2) {
		t.Fatalf("two subsets should equal")
	}
}

func Test_FilterClusters(t *testing.T) {
	inss := map[string][]*Instance{
		"sh001": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c1"},
		}, &Instance{
			Addrs:    []string{"http://127.0.0.2:9000"},
			Metadata: map[string]string{MetaCluster: "c1"},
		}, &Instance{
			Addrs:    []string{"grpc://127.0.0.3:9000"},
			Metadata: map[string]string{MetaCluster: "c2"},
		}},
		"sh002": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c3"},
		}, &Instance{
			Addrs:    []string{"zk://127.0.0.2:9000"},
			Metadata: map[string]string{MetaCluster: "c3"},
		}},
	}
	res := map[string][]*Instance{
		"sh001": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c1"},
		}},
		"sh002": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c3"},
		}},
	}
	var opt BuildOptions
	f := Filter("grpc", map[string]struct{}{"c1": struct{}{}, "c3": struct{}{}})
	f.Apply(&opt)
	filtered := opt.Filter(inss)
	equal := reflect.DeepEqual(filtered, res)
	if !equal {
		t.Fatalf("Filter grpc should equal,filtered:%v expected:%v", filtered, res)
	}
}

func Test_FilterInvalidAddr(t *testing.T) {
	inss := map[string][]*Instance{
		"sh001": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c1"},
		}, &Instance{
			Addrs:    []string{"http://127.0.0.2:9000"},
			Metadata: map[string]string{MetaCluster: "c1"},
		}, &Instance{
			Addrs:    []string{"grpc://127.0.0.3:9000"},
			Metadata: map[string]string{MetaCluster: "c2"},
		}},
		"sh002": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c3"},
		}, &Instance{
			Addrs:    []string{"zk://127.0.0.2:9000"},
			Metadata: map[string]string{MetaCluster: "c3"},
		}},
	}
	res := map[string][]*Instance{
		"sh001": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c1"},
		}, &Instance{
			Addrs:    []string{"grpc://127.0.0.3:9000"},
			Metadata: map[string]string{MetaCluster: "c2"},
		}},
		"sh002": []*Instance{&Instance{
			Addrs:    []string{"grpc://127.0.0.1:9000"},
			Metadata: map[string]string{MetaCluster: "c3"},
		}},
	}
	var opt BuildOptions
	f := Filter("grpc", nil)
	f.Apply(&opt)
	filtered := opt.Filter(inss)
	equal := reflect.DeepEqual(filtered, res)
	if !equal {
		t.Fatalf("Filter grpc should equal,filtered:%v expected:%v", filtered, res)
	}
}

func Test_Schedule(t *testing.T) {
	app := &InstancesInfo{
		Instances: map[string][]*Instance{
			"sh001": []*Instance{&Instance{
				Zone:     "sh001",
				Addrs:    []string{"grpc://127.0.0.1:9000"},
				Metadata: map[string]string{MetaCluster: "c1"},
			}, &Instance{
				Zone:     "sh001",
				Addrs:    []string{"grpc://127.0.0.2:9000"},
				Metadata: map[string]string{MetaCluster: "c1"},
			}, &Instance{
				Zone:     "sh001",
				Addrs:    []string{"grpc://127.0.0.3:9000"},
				Metadata: map[string]string{MetaCluster: "c2"},
			}},
			"sh002": []*Instance{&Instance{
				Zone:     "sh002",
				Addrs:    []string{"grpc://127.0.0.1:9000"},
				Metadata: map[string]string{MetaCluster: "c3"},
			}, &Instance{
				Zone:     "sh002",
				Addrs:    []string{"grpc://127.0.0.2:9000"},
				Metadata: map[string]string{MetaCluster: "c3"},
			}},
		},
		Scheduler: &Scheduler{map[string]*ZoneStrategy{"sh001": &ZoneStrategy{
			Zones: map[string]*Strategy{
				"sh001": &Strategy{10},
				"sh002": &Strategy{20},
			},
		}}},
	}
	var opt BuildOptions
	f := ScheduleNode("sh001")
	f.Apply(&opt)
	err := compareAddr(opt.Scheduler(app), map[string]int{"sh002": 2, "sh001": 1})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_Schedule2(t *testing.T) {
	app := &InstancesInfo{
		Instances: map[string][]*Instance{},
		Scheduler: &Scheduler{map[string]*ZoneStrategy{"sh001": &ZoneStrategy{
			Zones: map[string]*Strategy{
				"sh001": &Strategy{10},
				"sh002": &Strategy{20},
			},
		}}},
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh001",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c1"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh002",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c2"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}
	var opt BuildOptions
	f := ScheduleNode("sh001")
	f.Apply(&opt)
	err := compareAddr(opt.Scheduler(app), map[string]int{"sh002": 30, "sh001": 15})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_Schedule3(t *testing.T) {
	app := &InstancesInfo{
		Instances: map[string][]*Instance{},
		Scheduler: &Scheduler{map[string]*ZoneStrategy{"sh001": &ZoneStrategy{
			Zones: map[string]*Strategy{
				"sh001": &Strategy{1},
				"sh002": &Strategy{30},
			},
		}}},
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh001",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c1"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh002",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c2"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}
	var opt BuildOptions
	f := ScheduleNode("sh001")
	f.Apply(&opt)
	err := compareAddr(opt.Scheduler(app), map[string]int{"sh002": 30, "sh001": 1})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_Schedule4(t *testing.T) {
	app := &InstancesInfo{
		Instances: map[string][]*Instance{},
		Scheduler: &Scheduler{map[string]*ZoneStrategy{"sh001": &ZoneStrategy{
			Zones: map[string]*Strategy{
				"sh001": &Strategy{1},
				"sh002": &Strategy{30},
			},
		}}},
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh001",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c1"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}

	var opt BuildOptions
	f := ScheduleNode("sh001")
	f.Apply(&opt)
	err := compareAddr(opt.Scheduler(app), map[string]int{"sh001": 30, "sh002": 0})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_Schedule5(t *testing.T) {
	app := &InstancesInfo{
		Instances: map[string][]*Instance{},
		Scheduler: &Scheduler{map[string]*ZoneStrategy{"sh001": &ZoneStrategy{
			Zones: map[string]*Strategy{
				"sh002": &Strategy{30},
			},
		}}},
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh001",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c1"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}
	for i := 0; i < 30; i++ {
		ins := &Instance{
			Zone:     "sh002",
			Addrs:    []string{fmt.Sprintf("grpc://127.0.0.%d:9000", i)},
			Metadata: map[string]string{MetaCluster: "c2"},
		}
		app.Instances[ins.Zone] = append(app.Instances[ins.Zone], ins)
	}
	var opt BuildOptions
	f := ScheduleNode("sh001")
	f.Apply(&opt)
	err := compareAddr(opt.Scheduler(app), map[string]int{"sh002": 30, "sh001": 0})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func compareAddr(inss []*Instance, c map[string]int) (err error) {
	for _, ins := range inss {
		c[ins.Zone] = c[ins.Zone] - 1
	}
	for zone, v := range c {
		if v != 0 {
			err = fmt.Errorf("zone(%s) nums is %d", zone, v)
			return
		}
	}
	return
}
