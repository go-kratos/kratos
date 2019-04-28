package dao

import (
	"github.com/ipipdotnet/ipdb-go"
)

const (
	ip4dbaddr = "/data/conf/v4.ipdb"
	ip6dbaddr = "/data/conf/v6.ipdb"
)

var (
	//IP4db ip4地址库
	IP4db *ipdb.City
	//IP6db ip6地址库
	IP6db *ipdb.City
)

//InitIPdb 初始化IPdb
func InitIPdb() {
	var err error
	IP4db, err = ipdb.NewCity(ip4dbaddr)
	if err != nil {
		panic(err)
	}
	IP6db, err = ipdb.NewCity(ip6dbaddr)
	if err != nil {
		panic(err)
	}
}
