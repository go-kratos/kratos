// Package dsn implements dsn parse with struct bind
/*
DSN 格式类似 URI, DSN 结构如下图

	network:[//[username[:password]@]address[:port][,address[:port]]][/path][?query][#fragment]

与 URI 的主要区别在于 scheme 被替换为 network, host 被替换为 address 并且支持多个 address.
network 与 net 包中 network 意义相同, tcp、udp、unix 等, address 支持多个使用 ',' 分割, 如果
network 为 unix 等本地 sock 协议则使用 Path, 有且只有一个

dsn 包主要提供了 Parse, Bind 和 validate 功能

Parse 解析 dsn 字符串成 DSN struct, DSN struct 与 url.URL 几乎完全一样

Bind 提供将 DSN 数据绑定到一个 struct 的功能, 通过 tag dsn:"key,[default]" 指定绑定的字段, 目前支持两种类型的数据绑定

内置变量 key:
	network string tcp, udp, unix 等, 参考 net 包中的 network
	username string
	password string
	address string or []string address 可以绑定到 string 或者 []string, 如果为 string 则取 address 第一个

Query: 通过 query.name 可以取到 query 上的数据

	数组可以通过传递多个获得

	array=1&array=2&array3 -> []int `tag:"query.array"`

	struct 支持嵌套

	foo.sub.name=hello&foo.tm=hello

	struct Foo {
		Tm string `dsn:"query.tm"`
		Sub struct {
			Name string `dsn:"query.name"`
		} `dsn:"query.sub"`
	}

默认值: 通过 dsn:"key,[default]" 默认值暂时不支持数组

忽略 Bind: 通过 dsn:"-" 忽略 Bind

自定义 Bind: 可以同时实现 encoding.TextUnmarshaler 自定义 Bind 实现

Validate: 参考 https://github.com/go-playground/validator

使用参考: example_test.go

DSN 命名规范:

没有历史遗留的情况下，尽量使用 Address, Network, Username, Password 等命名，代替之前的 Proto 和 Addr 等命名

Query 命名参考, 使用驼峰小写开头:

	timeout 通用超时
	dialTimeout 连接建立超时
	readTimeout 读操作超时
	writeTimeout 写操作超时
	readsTimeout 批量读超时
	writesTimeout 批量写超时
*/
package dsn
