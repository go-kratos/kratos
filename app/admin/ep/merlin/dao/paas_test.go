package dao

import (
	"testing"

	"go-common/app/admin/ep/merlin/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	gmRequest = model.PaasGenMachineRequest{
		Env: model.Env{
			ClusterID: 3,
			NetworkID: 7,
		},
		Machines: []model.PaasMachine{
			{
				Name:          "vm1",
				Image:         "docker-reg.bilibili.co/zccdebian:1.0",
				CPURequest:    1000,
				MemoryRequest: 500,
				DiskRequest:   20,
				VolumnMount:   "",
			},
			{
				Name:          "vm2",
				Image:         "docker-reg.bilibili.co/debian:latest",
				CPURequest:    1000,
				MemoryRequest: 500,
				DiskRequest:   20,
				VolumnMount:   "",
			},
		},
	}

	qdmRequest = model.PaasQueryAndDelMachineRequest{
		BusinessUnit: "ops",
		Project:      "zcc",
		App:          "zcc",
		ClusterID:    3,
		Name:         "vm1",
	}
)

func Test_GenPaasMachines(t *testing.T) {
	Convey("create two machines in paas", t, WithPaasToken(func() {
		data := `{"status": 200,"message":"success","data":[{"instance_name": "vm1","instance_create_status": 2},{"instance_name": "vm2","instance_create_status": 1}]}`
		gmURL := d.c.Paas.Host + _genPaasMachines
		httpMock("POST", gmURL).Reply(200).JSON(data)
		res, err := d.GenPaasMachines(ctx(), &gmRequest)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, []*model.CreateInstance{
			{
				Instance:             model.Instance{InstanceName: "vm1"},
				InstanceCreateStatus: 2,
			},
			{
				Instance:             model.Instance{InstanceName: "vm2"},
				InstanceCreateStatus: 1,
			},
		})
	}))
}

func Test_DelPaasMachine(t *testing.T) {
	Convey("Delete machine from paas", t, WithPaasToken(func() {
		data := `{"status": 200,"message": "success","data": {"instance_name": "vm1","instance_release_status": 1}}`
		dmURL := d.c.Paas.Host + _delPaasMachine
		httpMock("POST", dmURL).Reply(200).JSON(data)
		res, err := d.DelPaasMachine(ctx(), &qdmRequest)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, &model.ReleaseInstance{
			Instance:              model.Instance{InstanceName: "vm1"},
			InstanceReleaseStatus: 1,
		})
	}))
}

func Test_QueryPaasMachineStatus(t *testing.T) {
	Convey("query machine status", t, WithPaasToken(func() {
		data := `{"status": 200,"message": "success","data": {"condition": "Running","message": "","detail_conditions": {"Initialized": "True","PodScheduled": "True","Ready": "True"},"instance_ip": "172.16.62.84","restart_count": 0}}`
		qmsURL := d.c.Paas.Host + _queryPaasMachineStatus
		httpMock("POST", qmsURL).Reply(200).JSON(data)
		res, err := d.QueryPaasMachineStatus(ctx(), &qdmRequest)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, &model.MachineStatus{
			Condition: "Running",
			Message:   "",
			DetailCondition: model.DetailCondition{
				Initialized:  "True",
				PodScheduled: "True",
				Ready:        "True",
			},
			InstanceIP:   "172.16.62.84",
			RestartCount: 0,
		})
	}))
}

func Test_QueryPaasMachine(t *testing.T) {
	Convey("query machine details", t, WithPaasToken(func() {
		data := `{"status": 200,"message": "success","data": {"condition": "Running","name": "vm1","image": "docker-reg.bilibili.co/zccdebian:1.0","cpu_request": 1000,"memory_request": 512,"disk_request": 0,"volumn_mount": "","cluster_name": "shyp-fat-k8s","env": "fat1","ip": "172.16.62.84"}}`
		qmdURL := d.c.Paas.Host + _queryPaasMachine
		httpMock("POST", qmdURL).Reply(200).JSON(data)
		res, err := d.QueryPaasMachine(ctx(), &qdmRequest)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, &model.PaasMachineDetail{
			Condition:     "Running",
			Name:          "vm1",
			Image:         "docker-reg.bilibili.co/zccdebian:1.0",
			CPURequest:    1000,
			MemoryRequest: 512,
			DiskRequest:   0,
			VolumnMount:   "",
			ClusterName:   "shyp-fat-k8s",
			Env:           "fat1",
			IP:            "172.16.62.84",
		})

	}))
}

func Test_QueryClusters(t *testing.T) {
	Convey("query clusters from paas", t, WithPaasToken(func() {
		data := `{"status": 200,"message": "success","count": 4,"data": {"items": [{"id": 1,"name": "shylf-uat-k8s","desc": "上海云立集成环境k8s集群","idc": "shylf","networks": [{"id": 16,"name": "shylf_uat_vlan1035","subnet": "172.22.35.0/24","capacity": 0.3346774193548387}],"resources": {"cpu_usage": 0.5977777777777777,"mem_usage": 0.19947637451689315,"pod_total": 439,"pod_capacity": 990,"nodes_num": 9}}]}}`
		qcURL := d.c.Paas.Host + _queryPaasClusters
		httpMock("GET", qcURL).Reply(200).JSON(data)
		res, err := d.QueryClusters(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldResemble, []*model.Cluster{
			{
				ID:   1,
				Name: "shylf-uat-k8s",
				Desc: "上海云立集成环境k8s集群",
				IDc:  "shylf",
				Networks: []model.Network{
					{
						ID:       16,
						Name:     "shylf_uat_vlan1035",
						Subnet:   "172.22.35.0/24",
						Capacity: 0.3346774193548387,
					},
				},
				Resources: model.Resource{
					CPUUsage:    0.5977777777777777,
					MemUsage:    0.19947637451689315,
					PodTotal:    439,
					PodCapacity: 990,
					NodesNum:    9,
				},
			},
		})

	}))
}
