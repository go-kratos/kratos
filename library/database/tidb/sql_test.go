package tidb

import "testing"

func TestParseDSNAddrr(t *testing.T) {
	dsn := "u:p@tcp(127.0.0.1:3306)/db?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8,utf8mb4"
	addr := parseDSNAddr(dsn)
	if addr != "127.0.0.1:3306" {
		t.Errorf("expect 127.0.0.1:3306 got: %s", addr)
	}
}

func TestGenDSN(t *testing.T) {
	dsn := "u:p@tcp(127.0.0.1:3306)/db?loc=Local&parseTime=true&readTimeout=5s&timeout=5s&writeTimeout=5s&charset=utf8mb4"
	addr := "127.0.0.2:3308"
	expect := "u:p@tcp(127.0.0.2:3308)/db?loc=Local&parseTime=true&readTimeout=5s&timeout=5s&writeTimeout=5s&charset=utf8mb4"
	newdsn := genDSN(dsn, addr)
	if newdsn != expect {
		t.Errorf("expect %s got: %s", expect, newdsn)
	}
}
