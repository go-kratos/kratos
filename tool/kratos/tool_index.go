package main

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "kratos",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos",
		Summary:   "Kratos工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
		URL:       "wiki",
	},
	{
		Name:      "kprotoc",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "bash -c ${GOPATH}/src/github.com/bilibili/kratos/tool/kprotoc/install_kprotoc.sh",
		Summary:   "快速方便生成pb.go的protoc封装",
		Platform:  []string{"darwin", "linux"},
		Author:    "kratos",
		URL:       "wiki",
	},
	{
		Name:      "cachegen",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/cachegen",
		Summary:   "缓存回源代码生成",
		Platform:  []string{"darwin", "linux"},
		Author:    "kratos",
		URL:       "wiki",
	},
	{
		Name:      "mcgen",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/cachegen/mcgen",
		Summary:   "memcache代码生成",
		Platform:  []string{"darwin", "linux"},
		Author:    "kratos",
		URL:       "wiki",
	},
}
