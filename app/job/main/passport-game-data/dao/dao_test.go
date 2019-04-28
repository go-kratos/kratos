package dao

import (
	"fmt"
	"sync"

	"go-common/app/job/main/passport-game-data/conf"
)

var (
	once sync.Once
	d    *Dao
)

func startDao() {
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	d = New(conf.Conf)
}
