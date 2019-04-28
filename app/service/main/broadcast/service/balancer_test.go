package service

import (
	"sort"
	"testing"

	"go-common/app/service/main/broadcast/model"
)

func TestCurrentWeightByConns(t *testing.T) {
	nodes := []*weightedNode{
		{fixedWeight: 10, currentWeight: 0, currentConns: 100000},
		{fixedWeight: 10, currentWeight: 0, currentConns: 100000},
		{fixedWeight: 10, currentWeight: 0, currentConns: 100000},
	}
	totalWeight := nodes[0].fixedWeight + nodes[1].fixedWeight + nodes[2].fixedWeight
	totalConns := nodes[0].currentConns + nodes[1].currentConns + nodes[2].currentConns
	for i := 0; i < 1000000; i++ {
		for _, n := range nodes {
			n.calcuateWeightByConns(totalWeight, totalConns, n.currentConns, 1.0)
		}
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].currentWeight > nodes[j].currentWeight
		})
		nodes[0].chosen()
		totalConns++
	}
	ft := float64(nodes[0].fixedWeight + nodes[1].fixedWeight + nodes[2].fixedWeight)
	ct := float64(nodes[0].currentConns + nodes[1].currentConns + nodes[2].currentConns)
	for _, n := range nodes {
		k, j := float64(n.fixedWeight)/ft*10, float64(n.currentConns)/ct*10
		if j/k < 1/3 {
			t.Errorf("unmatch %d:%d", int(k), int(j))
		}
		t.Logf("node:%+v ratio %d:%d", n, int(k), int(j))
	}
}

func TestFixedWeightByConns(t *testing.T) {
	nodes := []*weightedNode{
		{fixedWeight: 1, currentWeight: 0, currentConns: 0},
		{fixedWeight: 2, currentWeight: 0, currentConns: 0},
		{fixedWeight: 3, currentWeight: 0, currentConns: 0},
	}
	totalWeight := nodes[0].fixedWeight + nodes[1].fixedWeight + nodes[2].fixedWeight
	totalConns := nodes[0].currentConns + nodes[1].currentConns + nodes[2].currentConns
	for i := 0; i < 1000000; i++ {
		for _, n := range nodes {
			n.calcuateWeightByConns(totalWeight, totalConns, n.currentConns, 1.0)
		}
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].currentWeight > nodes[j].currentWeight
		})
		nodes[0].chosen()
		totalConns++
	}
	ft := float64(nodes[0].fixedWeight + nodes[1].fixedWeight + nodes[2].fixedWeight)
	ct := float64(nodes[0].currentConns + nodes[1].currentConns + nodes[2].currentConns)
	for _, n := range nodes {
		k, j := float64(n.fixedWeight)/ft*10, float64(n.currentConns)/ct*10
		if j/k < 1/3 {
			t.Errorf("unmatch %d:%d", int(k), int(j))
		}
		t.Logf("node:%+v ratio %d:%d", n, int(k), int(j))
	}
}

func TestCurrentWeightByGeo(t *testing.T) {
	nodes := []*weightedNode{
		{region: "bj", fixedWeight: 10, currentWeight: 0, currentConns: 0, lat: 39.904989, lng: 116.405285},
		{region: "sh", fixedWeight: 10, currentWeight: 0, currentConns: 0, lat: 31.231706, lng: 121.472644},
		{region: "gz", fixedWeight: 10, currentWeight: 0, currentConns: 0, lat: 23.125178, lng: 113.280637},
	}
	totalWeight := nodes[0].fixedWeight + nodes[1].fixedWeight + nodes[2].fixedWeight
	totalConns := nodes[0].currentConns + nodes[1].currentConns + nodes[2].currentConns
	lat, lng := float64(31), float64(121)
	for i := 0; i < 1000000; i++ {
		var totalDist float64
		for _, n := range nodes {
			if n.lat != 0 && n.lng != 0 {
				totalDist += distance(n.lat, n.lng, lat, lng)
			}
		}
		for _, n := range nodes {
			n.calcuateWeightByGeo(totalWeight, totalConns, n.currentConns, totalDist, distance(n.lat, n.lng, lat, lng))
		}
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].currentWeight > nodes[j].currentWeight
		})
		nodes[0].chosen()
		totalConns++
	}
	ft := float64(nodes[0].fixedWeight + nodes[1].fixedWeight + nodes[2].fixedWeight)
	ct := float64(nodes[0].currentConns + nodes[1].currentConns + nodes[2].currentConns)
	for _, n := range nodes {
		k, j := float64(n.fixedWeight)/ft*10, float64(n.currentConns)/ct*10
		if j/k < 1/3 {
			t.Errorf("unmatch %d:%d", int(k), int(j))
		}
		t.Logf("node:%+v ratio %d:%d", n, int(k), int(j))
	}
}

func TestLoadBalancer(t *testing.T) {
	var (
		ss = []*model.ServerInfo{
			{
				Region:    "bj",
				Server:    "01",
				Weight:    10,
				ConnCount: 300000,
				IPCount:   100000,
				IPAddrs:   []string{"ip_bj"},
				IPAddrsV6: []string{"ip_bj_v6"},
				Latitude:  39.904989,
				Longitude: 116.405285,
			},
			{
				Region:    "sh",
				Server:    "02",
				Weight:    10,
				ConnCount: 300000,
				IPCount:   10,
				IPAddrs:   []string{"ip_sh"},
				IPAddrsV6: []string{"ip_sh_v6"},
				Latitude:  31.231706,
				Longitude: 121.472644,
			},
			{
				Region:    "gz",
				Server:    "03",
				Weight:    10,
				ConnCount: 300000,
				IPCount:   100000,
				IPAddrs:   []string{"ip_gz"},
				IPAddrsV6: []string{"ip_gz_v6"},
				Latitude:  23.125178,
				Longitude: 113.280637,
				Overseas:  true,
			},
		}
	)
	lb := NewLoadBalancer()
	lb.Update(ss)
	for i := 0; i < 10; i++ {
		t.Log(lb.NodeAddrsByGeo(".text", 39, 116, false, false))
		t.Log(lb.NodeAddrsByGeo(".text", 31, 121, false, false))
		t.Log(lb.NodeAddrsByGeo(".text", 23, 121, true, false))
	}
	t.Log(lb.nodes)
	for i := 0; i < 10; i++ {
		t.Log(lb.NodeAddrsByRegion("sh", ".test", 1.5, false, false))
	}
	t.Log(lb.nodes)
	for i := 0; i < 10; i++ {
		t.Log(lb.NodeAddrsByRegion("sh", ".test", 1.5, true, true))
	}
	t.Log(lb.nodes)
}
