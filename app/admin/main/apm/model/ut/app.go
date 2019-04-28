package ut

import (
	"go-common/library/time"
	"sync"
)

// TableName .
func (*App) TableName() string {
	return "ut_app"
}

// App ..
type App struct {
	ID       int64     `gorm:"column:id" json:"id"`
	Path     string    `gorm:"column:path" json:"path"`
	Owner    string    `gorm:"column:owner" json:"owner"`
	HasUt    int       `gorm:"column:has_ut" json:"has_ut"`
	Link     string    `gorm:"-" json:"link"`
	Coverage float64   `gorm:"coverage"`
	CTime    time.Time `gorm:"column:ctime" json:"ctime"`
	MTime    time.Time `gorm:"column:mtime" json:"mtime"`
}

// AppReq .
type AppReq struct {
	HasUt int    `form:"has_ut"`
	Path  string `form:"path"`
	Pn    int    `form:"pn" default:"1"`
	Ps    int    `form:"ps" defalut:"20"`
}

// Department .
type Department struct {
	Name     string
	Total    int64
	Access   int64
	Coverage float64
}

// AppsCache apps cache.
type AppsCache struct {
	Slice []*App
	Map   map[string]*App
	Owner map[string][]*App
	Dept  map[string]*Department
	sync.Mutex
}

//PathsByOwner get app paths by owner.
func (apps *AppsCache) PathsByOwner(owner string) (paths []string) {
	for _, app := range apps.Owner[owner] {
		paths = append(paths, app.Path)
	}
	return
}
