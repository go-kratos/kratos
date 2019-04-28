/*Package log 是kratos日志库.

一、主要功能：

	1. 日志打印到elk
	2. 日志打印到本地，内部使用log4go
	3. 日志打印到标准输出
	4. verbose日志实现，参考glog实现，可通过设置不同verbose级别，默认不开启

二、日志配置

1. 默认agent配置

	目前日志已经实现默认配置，可以根据env自动切换远程日志。可以直接使用以下方式：
	log.Init(nil)

2. 启动参数 or 环境变量

	启动参数		环境变量		说明
	log.stdout	LOG_STDOUT	是否开启标准输出
	log.agent	LOG_AGENT	远端日志地址：unixpacket:///var/run/lancer/collector_tcp.sock?timeout=100ms&chan=1024
	log.dir		LOG_DIR		文件日志路径
	log.v		LOG_V		verbose日志级别
	log.module	LOG_MODULE	可单独配置每个文件的verbose级别：file=1,file2=2
	log.filter	LOG_FILTER	配置需要过滤的字段：field1,field2

3. 配置文件
但是如果有特殊需要可以走一下格式配置：
	[log]
		family = "xxx-service"
		dir = "/data/log/xxx-service/"
		stdout = true
		vLevel = 3
		filter = ["fileld1", "field2"]
	[log.module]
		"dao_user" = 2
		"servic*" = 1
	[log.agent]
		taskID = "00000x"
		proto = "unixpacket"
		addr = "/var/run/lancer/collector_tcp.sock"
		chanSize = 10240

三、配置说明

1.log

	family		项目名，默认读环境变量$APPID
	studout		标准输出，prod环境不建议开启
	filter		配置需要过滤掉的字段，以“***”替换
	dir		文件日志地址，prod环境不建议开启
	v		开启verbose级别日志，可指定全局级别

2. log.module

	可单独配置每个文件的verbose级别

3. log.agent
远端日志配置项
	taskID		lancer分配的taskID
	proto		网络协议，常见：tcp, udp, unixgram
	addr		网络地址，常见：ip:prot, sock
	chanSize	日志队列长度
*/
package log
