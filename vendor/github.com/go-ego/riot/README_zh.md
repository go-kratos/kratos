# [Riot 搜索引擎](https://github.com/go-ego/riot)

<!--<img align="right" src="https://raw.githubusercontent.com/go-ego/ego/master/logo.jpg">-->
<!--<a href="https://circleci.com/gh/go-ego/ego/tree/dev"><img src="https://img.shields.io/circleci/project/go-ego/ego/dev.svg" alt="Build Status"></a>-->
[![CircleCI Status](https://circleci.com/gh/go-ego/riot.svg?style=shield)](https://circleci.com/gh/go-ego/riot)
![Appveyor](https://ci.appveyor.com/api/projects/status/github/go-ego/riot?branch=master&svg=true)
[![codecov](https://codecov.io/gh/go-ego/riot/branch/master/graph/badge.svg)](https://codecov.io/gh/go-ego/riot)
[![Build Status](https://travis-ci.org/go-ego/riot.svg)](https://travis-ci.org/go-ego/riot)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-ego/riot)](https://goreportcard.com/report/github.com/go-ego/riot)
[![GoDoc](https://godoc.org/github.com/go-ego/riot?status.svg)](https://godoc.org/github.com/go-ego/riot)
[![Release](https://github-release-version.herokuapp.com/github/go-ego/riot/release.svg?style=flat)](https://github.com/go-ego/riot/releases/latest)
[![Join the chat at https://gitter.im/go-ego/ego](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-ego/ego?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<!--<a href="https://github.com/go-ego/ego/releases"><img src="https://img.shields.io/badge/%20version%20-%206.0.0%20-blue.svg?style=flat-square" alt="Releases"></a>-->


Go Open Source, Distributed, Simple and efficient full text search engine.

# Features

* [高效索引和搜索](/docs/zh/benchmarking.md)（1M 条微博 500M 数据28秒索引完，1.65毫秒搜索响应时间，19K 搜索 QPS）
* 支持中文分词（使用 [gse 分词包](https://github.com/go-ego/gse)并发分词，速度 27MB/秒）
* 支持[逻辑搜索](https://github.com/go-ego/riot/blob/master/docs/zh/logic.md)
* 支持中文转拼音搜索(使用 [gpy](https://github.com/go-ego/gpy) 中文转拼音)
* 支持计算关键词在文本中的[紧邻距离](/docs/zh/token_proximity.md)（token proximity）
* 支持计算[BM25相关度](/docs/zh/bm25.md)
* 支持[自定义评分字段和评分规则](/docs/zh/custom_scoring_criteria.md)
* 支持[在线添加、删除索引](/docs/zh/realtime_indexing.md)
* 支持多种[持久存储](/docs/zh/persistent_storage.md)
* 支持 heartbeat
* 支持[分布式索引和搜索](https://github.com/go-ego/riot/tree/master/data)
* 可实现[分布式索引和搜索](/docs/zh/distributed_indexing_and_search.md)
* 采用对商业应用友好的[Apache License v2](/LICENSE)发布

* [查看分词规则](https://github.com/go-ego/riot/blob/master/docs/zh/segmenter.md)

Riot v0.10.0 was released in Nov 2017, check the [Changelog](https://github.com/go-ego/riot/blob/master/docs/CHANGELOG.md) for the full details.

QQ 群: 120563750

## 安装/更新

```
go get -u github.com/go-ego/riot
```

## Requirements

需要 Go 版本至少 1.8

### Vendored Dependencies

Riot 使用 [dep](https://github.com/golang/dep) 管理 vendor 依赖, but we don't commit the vendored packages themselves to the Riot git repository. Therefore, a simple go get is not supported because the command is not vendor aware. 

请用 dep 管理它, 运行 `dep ensure` 克隆依赖.

## [Build-tools](https://github.com/go-ego/re)
```
go get -u github.com/go-ego/re 
```
### re riot
创建 riot 项目

```
$ re riot my-riotapp
```

### re run

运行我们创建的 riot 项目, 你可以导航到应用程序文件夹并执行:
```
$ cd my-riotapp && re run
```

## 使用

先看一个例子（来自 [simplest_example.go](/examples/simple/zh/main.go)）
```go
package main

import (
	"log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	// searcher 是协程安全的
	searcher = riot.Engine{}
)

func main() {
	// 初始化
	searcher.Init(types.EngineOpts{
		Using:             3,
		GseDict: "zh",
		// GseDict: "your gopath"+"/src/github.com/go-ego/riot/data/dict/dictionary.txt",
	})
	defer searcher.Close()

	text := "此次百度收购将成中国互联网最大并购"
	text1 := "百度宣布拟全资收购91无线业务"
	text2 := "百度是中国最大的搜索引擎"
	
	// 将文档加入索引，docId 从1开始
	searcher.Index(1, types.DocData{Content: text})
	searcher.Index(2, types.DocData{Content: text1}, false)
	searcher.Index(3, types.DocData{Content: text2}, true)

	// 等待索引刷新完毕
	searcher.Flush()
	// engine.FlushIndex()

	// 搜索输出格式见 types.SearchResp 结构体
	log.Print(searcher.Search(types.SearchReq{Text:"百度中国"}))
}
```

是不是很简单！

然后看看一个[入门教程](/docs/zh/codelab.md)，教你用不到200行 Go 代码实现一个微博搜索网站。

### 使用默认引擎:

```Go
package main

import (
	"log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

var (
	searcher = riot.New("zh")
)

func main() {
	data := types.DocData{Content: `I wonder how, I wonder why
		, I wonder where they are`}
	data1 := types.DocData{Content: "所以, 你好, 再见"}
	data2 := types.DocData{Content: "没有理由"}
	searcher.Index(1, data)
	searcher.Index(2, data1)
	searcher.IndexDoc(3, data2)
	searcher.Flush()

	req := types.SearchReq{Text: "你好"}
	search := searcher.Search(req)
	log.Println("search...", search)
}
```

#### [查看更多例子](https://github.com/go-ego/riot/tree/master/examples)

#### [持久化的例子](https://github.com/go-ego/riot/blob/master/examples/store/main.go)
#### [逻辑搜索的例子](https://github.com/go-ego/riot/blob/master/examples/logic/main.go)

#### [拼音搜索的例子](https://github.com/go-ego/riot/blob/master/examples/pinyin/main.go)

#### [不同字典和语言例子](https://github.com/go-ego/riot/blob/master/examples/dict/main.go)

#### [benchmark](https://github.com/go-ego/riot/blob/master/examples/benchmark/benchmark.go)

#### [Riot 搜索模板, 客户端和字典](https://github.com/go-ego/riot/tree/master/data)

## 主要改进:

- 增加逻辑搜索 
- 增加拼音搜索 
- 增加分布式 
- 分词等改进 
- 增加更多 api
- 支持 heartbeat
- 修复 bug
- 删除依赖 cgo 的存储引擎, 增加 badger和 leveldb 持久化引擎

## Donate

支持 riot, [buy me a coffee](https://github.com/go-vgo/buy-me-a-coffee).

#### Paypal

Donate money by [paypal](https://www.paypal.me/veni0/25) to my account [vzvway@gmail.com](vzvway@gmail.com)

## 其它

* [为什么要有 riot 引擎](/docs/zh/why_riot.md)
* [联系方式](/docs/zh/feedback.md)

## License

Riot is primarily distributed under the terms of the Apache License (Version 2.0), base on [wukong](https://github.com/huichen/wukong).
