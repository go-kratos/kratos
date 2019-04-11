package dsn_test

import (
	"log"

	"github.com/bilibili/kratos/pkg/conf/dsn"
	xtime "github.com/bilibili/kratos/pkg/time"
)

// Config struct
type Config struct {
	Network  string         `dsn:"network" validate:"required"`
	Host     string         `dsn:"host" validate:"required"`
	Username string         `dsn:"username" validate:"required"`
	Password string         `dsn:"password" validate:"required"`
	Timeout  xtime.Duration `dsn:"query.timeout,1s"`
	Offset   int            `dsn:"query.offset" validate:"gte=0"`
}

func ExampleParse() {
	cfg := &Config{}
	d, err := dsn.Parse("tcp://root:toor@172.12.12.23:2233?timeout=10s")
	if err != nil {
		log.Fatal(err)
	}
	_, err = d.Bind(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", cfg)
}
