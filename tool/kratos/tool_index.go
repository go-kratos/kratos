package main

import "time"

var toolIndexs = []*Tool{
	&Tool{
		Name:      "kratos",
		Alias:     "kratos",
		BuildTime: time.Date(2019, 4, 2, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos",
		Summary:   "Kratos工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
	},
	&Tool{
		Name:      "protoc",
		Alias:     "kratos-protoc",
		BuildTime: time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos-protoc",
		Summary:   "快速方便生成pb.go的protoc封装，windows、Linux请先安装protoc工具",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
	},
	&Tool{
		Name:      "swagger",
		Alias:     "swagger",
		BuildTime: time.Date(2019, 5, 5, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/go-swagger/go-swagger/cmd/swagger",
		Summary:   "swagger api文档",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "goswagger.io",
	},
}
