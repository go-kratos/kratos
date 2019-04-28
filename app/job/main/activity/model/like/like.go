package like

import (
	"database/sql/driver"
	"go-common/app/service/main/archive/api"
	"time"
)

// Like like
type Like struct {
	ID      int64    `json:"id"`
	Wid     int64    `json:"wid"`
	Archive *api.Arc `json:"archive,omitempty"`
}

// Item like item struct.
type Item struct {
	ID       int64     `json:"id"`
	Wid      int64     `json:"wid"`
	Ctime    wocaoTime `json:"ctime"`
	Sid      int64     `json:"sid"`
	Type     int       `json:"type"`
	Mid      int64     `json:"mid"`
	State    int       `json:"state"`
	StickTop int       `json:"stick_top"`
	Mtime    wocaoTime `json:"mtime"`
}

// Content  like_content.
type Content struct {
	ID      int64     `json:"id"`
	Message string    `json:"message"`
	IP      int64     `json:"ip"`
	Plat    int       `json:"plat"`
	Device  int       `json:"device"`
	Ctime   wocaoTime `json:"ctime"`
	Mtime   wocaoTime `json:"mtime"`
	Image   string    `json:"image"`
	Reply   string    `json:"reply"`
	Link    string    `json:"link"`
	ExName  string    `json:"ex_name"`
}

// WebData act web data.
type WebData struct {
	ID   int64  `json:"id"`
	Vid  int64  `json:"vid"`
	Data string `json:"data"`
}

// Action like_action .
type Action struct {
	ID     int64     `json:"id"`
	Lid    int64     `json:"lid"`
	Mid    int64     `json:"mid"`
	Action int64     `json:"action"`
	Ctime  wocaoTime `json:"ctime"`
	Mtime  wocaoTime `json:"mtime"`
	Sid    int64     `json:"sid"`
	IP     int64     `json:"ip"`
}

// Extend .
type Extend struct {
	ID    int64     `json:"id"`
	Lid   int64     `json:"lid"`
	Like  int64     `json:"like"`
	Ctime wocaoTime `json:"ctime"`
	Mtime wocaoTime `json:"mtime"`
}

// LastTmStat .
type LastTmStat struct {
	Last int64
}

type wocaoTime string

// Scan scan time.
func (jt *wocaoTime) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*jt = wocaoTime(sc.Format("2006-01-02 15:04:05"))
	case string:
		*jt = wocaoTime(sc)
	}
	return
}

// Value get time value.
func (jt wocaoTime) Value() (driver.Value, error) {
	return time.Parse("2006-01-02 15:04:05", string(jt))
}
