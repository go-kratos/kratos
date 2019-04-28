package report

import (
	"sync"
	"testing"
	"time"

	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

var (
	mnOnce      sync.Once
	mnUatOnce   sync.Once
	userOnce    sync.Once
	userUatOnce sync.Once
)

func newManager() {
	InitManager(nil)
}

func newUatManager() {
	InitManager(&databus.Config{
		Key:          "2511663d546f1413",
		Secret:       "cde3b480836cc76df3d635470f991caa",
		Group:        "LogAudit-MainSearch-P",
		Topic:        "LogAudit-T",
		Action:       "pub",
		Buffer:       10240,
		Name:         "log-audit/log-sub",
		Proto:        "tcp",
		Addr:         "172.18.33.50:6205",
		Active:       100,
		Idle:         100,
		DialTimeout:  xtime.Duration(time.Millisecond * 200),
		ReadTimeout:  xtime.Duration(time.Millisecond * 200),
		WriteTimeout: xtime.Duration(time.Millisecond * 200),
		IdleTimeout:  xtime.Duration(time.Second * 80),
	})
}

func newUser() {
	InitUser(nil)
}

func newUatUser() {
	InitManager(&databus.Config{
		Key:          "2511663d546f1413",
		Secret:       "cde3b480836cc76df3d635470f991caa",
		Group:        "LogUserAction-MainSearch-P",
		Topic:        "LogUserAction-T",
		Action:       "pub",
		Buffer:       10240,
		Name:         "log-user-action/log-sub",
		Proto:        "tcp",
		Addr:         "172.18.33.50:6205",
		Active:       100,
		Idle:         100,
		DialTimeout:  xtime.Duration(time.Millisecond * 200),
		ReadTimeout:  xtime.Duration(time.Millisecond * 200),
		WriteTimeout: xtime.Duration(time.Millisecond * 200),
		IdleTimeout:  xtime.Duration(time.Second * 80),
	})
}

func Test_Manager(b *testing.T) {
	mnOnce.Do(newManager)
	Manager(&ManagerInfo{
		Uname:    "dz",
		UID:      64,
		Business: 0,
		Type:     1,
		Oid:      2,
		Action:   "action",
		Ctime:    time.Now(),
		Index:    []interface{}{5, 6, 7, "a", "b", "c"},
		Content: map[string]interface{}{
			"json": "json",
		},
	})
}

func Test_UatManager(b *testing.T) {
	mnUatOnce.Do(newUatManager)
	Manager(&ManagerInfo{
		Uname:    "dz",
		UID:      64,
		Business: 0,
		Type:     1,
		Oid:      2,
		Action:   "action",
		Ctime:    time.Now(),
		Index:    []interface{}{5, 6, 7, "a", "b", "c"},
		Content: map[string]interface{}{
			"json": "json",
		},
	})
}

func Test_User(b *testing.T) {
	userOnce.Do(newUser)
	User(&UserInfo{
		Mid:      1,
		Platform: "platform",
		Build:    2,
		Buvid:    "buvid",
		Business: 0,
		Type:     3,
		Oid:      4,
		Action:   "action",
		Ctime:    time.Now(),
		IP:       "127.0.0.1",
		// extra
		Index: []interface{}{5, 6, 7, "a", "b", "c"},
		Content: map[string]interface{}{
			"json": "json",
		},
	})
}

func Test_UatUser(b *testing.T) {
	userUatOnce.Do(newUatUser)
	User(&UserInfo{
		Mid:      1,
		Platform: "platform",
		Build:    2,
		Buvid:    "buvid",
		Business: 0,
		Type:     3,
		Oid:      4,
		Action:   "action",
		Ctime:    time.Now(),
		IP:       "127.0.0.1",
		// extra
		Index: []interface{}{5, 6, 7, "a", "b", "c"},
		Content: map[string]interface{}{
			"json": "json",
		},
	})
}
