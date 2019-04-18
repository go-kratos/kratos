package main

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "kratos",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/Kratos/tool/kratos",
		Summary:   "Kratos工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
		URL:       "wiki",
	},
	{
		Name:      "cachegen",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/Kratos/tool/cachegen",
		Summary:   "缓存回源逻辑代码生成器",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
		URL:       "wiki",
	},
}
