package zlimit

import (
	"context"
	"testing"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/container/pool"
	"go-common/library/database/sql"
	xtime "go-common/library/time"
)

func TestZlimit(t *testing.T) {
	c := &Config{
		DB: &sql.Config{
			Addr:   "172.16.33.54:3306",
			DSN:    "test:test@tcp(172.16.33.54:3306)/bilibili_oversea?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8,utf8mb4",
			Active: 5,
			Idle:   2,
		},
		Redis: &Redis{
			Config: &redis.Config{
				Config: &pool.Config{
					Active:      20,
					Idle:        10,
					IdleTimeout: xtime.Duration(80 * time.Millisecond),
				},
				Name:         "business/zlimit",
				Proto:        "tcp",
				Addr:         "172.16.33.54:6380",
				DialTimeout:  xtime.Duration(50 * time.Millisecond),
				ReadTimeout:  xtime.Duration(100 * time.Millisecond),
				WriteTimeout: xtime.Duration(100 * time.Millisecond),
			},
			Expire: xtime.Duration(1 * time.Second),
		},
		IPFile:    "/data/conf/iprepo/iprepo.txt",
		FlushTime: xtime.Duration(1 * time.Second),
	}
	s := New(c)
	testFind(t, s, 0, "", "")
	testForbid(t, s, "", "")
	testFindByPid(t, s, []int64{0}, "")
	testFindByGid(t, s, 0, "", "")
}

func testFind(t *testing.T, s *Service, aid int64, ipaddr, cdnip string) {
	ret, retdown, err := s.Find(context.TODO(), aid, ipaddr, cdnip)
	if err != nil {
		t.Errorf("Service: find err: %v", err)
	} else {
		t.Logf("Service: find %d,%d", ret, retdown)
	}
}

func testForbid(t *testing.T, s *Service, pstr, ipaddr string) {
	if err := s.Forbid(context.TODO(), pstr, ipaddr); err != nil {
		t.Errorf("Service: forbid err: %v", err)
	}
}

func testFindByPid(t *testing.T, s *Service, pids []int64, ipaddr string) {
	ret, retdown := s.FindByPid(context.TODO(), pids, ipaddr)
	t.Logf("Service: findByPid %d,%d", ret, retdown)
}

func testFindByGid(t *testing.T, s *Service, gid int64, ipaddr, cdnip string) {
	ret, retdown := s.FindByGid(context.TODO(), gid, ipaddr, cdnip)
	t.Logf("Service: findByGid %d,%d", ret, retdown)
}
