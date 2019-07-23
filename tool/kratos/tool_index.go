package main

import "time"

var toolIndexs = []*Tool{
	&Tool{
		Name:      "kratos",
		Alias:     "kratos",
		BuildTime: time.Date(2019, 6, 21, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos",
		Summary:   "Kratos工具集本体",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
	},
	&Tool{
		Name:      "protoc",
		Alias:     "kratos-protoc",
		BuildTime: time.Date(2019, 6, 21, 0, 0, 0, 0, time.Local),
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
	&Tool{
		Name:      "genbts",
		Alias:     "kratos-gen-bts",
		BuildTime: time.Date(2019, 7, 23, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos-gen-bts",
		Summary:   "缓存回源逻辑代码生成器",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
	},
	&Tool{
		Name:      "genmc",
		Alias:     "kratos-gen-mc",
		BuildTime: time.Date(2019, 7, 23, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/bilibili/kratos/tool/kratos-gen-mc",
		Summary:   "mc缓存代码生成",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "kratos",
	},
}
