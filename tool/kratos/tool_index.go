package main

import "time"

var toolIndexs = []*Tool{
	&Tool{
		Name:      "kratos",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos",
		Summary:   "Kratos工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
		URL:       "wiki",
	},
}
