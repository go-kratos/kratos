package databus

import (
	"time"

	"go-common/library/conf/env"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

type conf struct {
	Key    string
	Secret string
	Addr   string
}

var (
	_defaultAddrConfig = map[string]*conf{
		env.DeployEnvUat: {
			Key:    "4c76cbb7a985ac90",
			Secret: "43bb22ce34a6b13e7814f09cb8116522",
			Addr:   "172.18.33.50:6205",
		},
		env.DeployEnvPre: {
			Key:    "4c76cbb7a985ac90",
			Secret: "8c22a196bf00d623b69cba4e61b1f7dc",
			Addr:   "172.18.21.90:6205",
		},
		env.DeployEnvProd: {
			Key:    "4c76cbb7a985ac90",
			Secret: "8c22a196bf00d623b69cba4e61b1f7dc",
			Addr:   "172.18.21.90:6205",
		},
	}

	_defaultConfig = &databus.Config{
		Key:          "0PtMsLLxWyyvoTgAyLCD",
		Secret:       "0PtMsLLxWyyvoTgAyLCE",
		Group:        "AegisResource-MainArchive-P",
		Topic:        "AegisResource-T",
		Action:       "pub",
		Buffer:       10240,
		Name:         "aegis-admin/databus",
		Proto:        "tcp",
		Addr:         "172.16.33.158:6205",
		Active:       100,
		Idle:         10,
		DialTimeout:  xtime.Duration(time.Millisecond * 200),
		ReadTimeout:  xtime.Duration(time.Millisecond * 200),
		WriteTimeout: xtime.Duration(time.Millisecond * 200),
		IdleTimeout:  xtime.Duration(time.Second * 80),
	}
)
