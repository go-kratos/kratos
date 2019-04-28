package ip

import "testing"

func TestIp(t *testing.T) {
	l, _ := New("./iprepo.txt")
	l, _ = New("./ip_test.txt")
	ips := []string{"000.000.000.001", "127.000.11.57", "183.131.11.57", "255.255.255.255", "3FFF:FFFF:FFFF:FEFF:FFFF:FFFF:FFFF:FFFF", "255.255.255", "0:0:0:0:0:0:0:1"}
	for _, ip := range ips {
		info := l.IP(ip)
		t.Log(info)
		zone := l.Zone(ip)
		t.Log(zone)
	}
	// Zone
	zone := l.Zone("")
	t.Log(zone)
	zone = l.Zone("0:0:0:0:0:0:0:1")
	t.Log(zone)
	// All
	infos := l.All()
	t.Log(infos)
	InternalIP()
	// InetAtoN
	InetAtoN("183.131.11.57")
	InetAtoN("183.131.11")
	InetAtoN("0:0:0:0:0:0:0:1")
	// InetNtoA
	InetNtoA(84549632)
	// ZoneID
	ZoneID("中国", "福建", "莆田")
}

func TestExternalIP(t *testing.T) {
	t.Log(ExternalIP())
}

func BenchmarkIP(b *testing.B) {
	l, _ := New("./iprepo.txt")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.IP("183.131.11.57")
		}
	})
}

func BenchmarkZone(b *testing.B) {
	l, _ := New("./iprepo.txt")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Zone("183.131.11.57")
		}
	})
}
