#### App wall 移动端商务对接接口

### Version 2.6.6

> 1.修复联通问题数据

### Version 2.6.5

> 1.修复日志拉取日期判断

### Version 2.6.4

> 1.修复日志拉取日期判断

### Version 2.6.3

> 1.不等ipv4 return

### Version 2.6.2

> 1.多机房

### Version 2.6.1

> 1.多机房

### Version 2.5.19

> 1.福利社日志修改

### Version 2.5.18

1.去掉net/ip调用

### Version 2.5.17

> 1.福利社日志修改

### Version 2.5.16

> 1.福利社用户日志区分

### Version 2.5.15

> 1.福利社绑定状态

### Version 2.5.14

> 1.大会员requestNo int64

### Version 2.5.13

> 1.ip方法更换

### Version 2.5.12

> 1.修复福利社日志问题

### Version 2.5.11

> 1.build

### Version 2.5.10

> 1.福利社用户日志

### Version 2.5.9

> 1.csrf false

### Version 2.5.8

> 1.使用grpc auth

### Version 2.5.7

> 1.去除MC KEY找不到的错误

### Version 2.5.6

> 1.福利社绑定用户
> 2.缓存修改

### Version 2.5.5

> 1.福利社日志查询

### Version 2.5.4

> 1.M站IP查询用户伪码解密

### Version 2.5.3

> 1.M站接口增加IP判断

### Version 2.5.2

> 1.M站接口改用联通加密

### Version 2.5.1

> 1.M站接口开新路由

### Version 2.5.0

> 1. update infoc sdk

##### Version 2.4.3

> 1.广点通

##### Version 2.4.2

> 1.增加流量领取间隔时间
> 2.增加只能流量卡领取限制
> 3.增加积分限制

##### Version 2.4.1

> 1.http active

##### Version 2.4.0

> 1.广点通

##### Version 2.3.9

> 1.seq server

##### Version 2.3.8

>1.修复 url


##### Version 2.3.7

>1.fix bug

##### Version 2.3.6

>2.慢查询

##### Version 2.3.5

>2.gdt重构

##### Version 2.3.4

>2.gdt response返回ret 0

##### Version 2.3.3

>2.gdt 新增 advertiser_id

##### Version 2.3.2

> 1.部分接口bm.CORS

##### Version 2.3.1

> 1.bm cors

##### Version 2.3.0

> 1.http层换成BM

##### Version 2.2.9

> 1.流量卡不能订购流量包

##### Version 2.2.8

> 1.联通流量包强行删除缓存

##### Version 2.2.7

> 1.联通IP异步同步下沉

##### Version 2.2.6

> 1.异步消费逻辑放到job

##### Version 2.2.5

> 1.修复积分异常添加问题

##### Version 2.2.4

> 1.增加礼包领取日志

##### Version 2.2.3

> 1.修复联通礼包BUG
> 2.修复联通订购关系问题

##### Version 2.2.2

> 1.修复联通订购关系BUG
> 2.fix头条重复激活

##### Version 2.2.1

> 1.fix头条重复激活

##### Version 2.2.0

> 1.联通礼包

##### Version 2.1.9

> 1.头条fix

##### Version 2.1.6

> 1.广点通广告投放点击上报

##### Version 2.1.5

> 1.今日头条广告投放点击上报

##### Version 2.1.4

> 1.drop statsd

##### Version 2.1.3

> 1.移动新订单默认流量100%

##### Version 2.1.2

> 1.增加运营商数据回掉结果日志

##### Version 2.1.1

> 1.增加红点

##### Version 2.1.0

> 1.联通订购数据查询修改

##### Version 2.0.9

> 1.移动流量包增加产品类型字段

##### Version 2.0.8

> 1.联通IP同步改为异步更新

##### Version 2.0.7

> 1.联通流量包订购关系接口下沉到服务端
> 2.去除手机号

##### Version 2.0.6

> 1.更新缓存逻辑修改
> 2.删除无用的日志

##### Version 2.0.5

> 1.运营商用户数据缓存修改

##### Version 2.0.4

> 1.缓存增加回去

##### Version 2.0.3

> 1.直播礼包切回原地址

##### Version 2.0.2

> 1.直播接口切换

##### Version 2.0.1

> 1.IOS客户端需要线上测试，暂时屏蔽缓存逻辑

##### Version 2.0.0

> 1.修复电信文档与实际请求不符合问题 

##### Version 1.9.9

> 1.判断电信流量包是否是有效的

##### Version 1.9.8

> 1.errgroup error return
> 2.验证电信接口状态是否正确

##### Version 1.9.7

> 1.电信接口增加errgroup减少接口超时时间  

##### Version 1.9.6

> 1.删除多余的 ecode.NoLogin  
> 2.电信端口改成string转int
> 3.增加错误日志

##### Version 1.9.5

> 1.电信增加提示  
> 2.电信增加短信模板  
> 3.增加流水号和手机号缓存   

##### Version 1.9.4

> 1.修复error没有return问题   

##### Version 1.9.3

> 1.电信接口返回error修改 

##### Version 1.9.2

> 1.电信用户状态接口改成Get请求  
> 2.电信用户许可新增状态  
> 3.电信支付接口返回订单流水号  

##### Version 1.9.1

> 1.电信接口修改   

##### Version 1.9.0

> 1.监控图名字修改  

##### Version 1.8.9

> 1.监控图名字修改  

##### Version 1.8.8

> 1.电信数据同步地址修改   

##### Version 1.8.7

> 1.增加缓存监控  

##### Version 1.8.6

> 1.联通、电信、移动用户增加缓存  
> 2.新增电信流量包服务  

##### Version 1.8.5

> 1.联通直播礼包接口日志修改  

##### Version 1.8.4

> 1.联通直播礼包接口  

##### Version 1.8.3

> 1.联通相关接口infoc数据上报   

##### Version 1.8.2

> 1.联通直播礼包实时查询数据库     

##### Version 1.8.1

> 1.移动用户IP判断  

##### Version 1.8.0

> 1.联通用户IP判断  

##### Version 1.7.9

> 1.移动逻辑修改  

##### Version 1.7.8

> 1.移动接口数据同步合并成一个接口   

##### Version 1.7.7

> 1.中国移动流量包   

##### Version 1.7.6

> 1.实时查库去除异步读数据库   
> 2.打印联通流量包状态日志   

##### Version 1.7.5

> 1.ip限制放入配置文件  

##### Version 1.7.4

> 1.修复BUG   

##### Version 1.7.3

> 1.修复message="0"的bug   

##### Version 1.7.2

> 1.合并大仓库   

##### Version 1.7.1

> 1.dotinapp渠道  

##### Version 1.7.0

> 1.添加IP白名单  

##### Version 1.6.9

> 1.联通用户状态处理  

##### Version 1.6.8

> 1.联通用户状态处理    

##### Version 1.6.7

> 1.添加IP白名单        

##### Version 1.6.6

> 1.unicomtype在其他接口不返回      

##### Version 1.6.5

> 1.修改显示订购状态   

##### Version 1.6.4

> 1.修改显示订购状态   

##### Version 1.6.3

> 1.增加状态接口   

##### Version 1.6.2

> 1.增加状态接口   

##### Version 1.6.1

> 1.增加预开户数据同步接口   
> 2.更新vendor  

##### Version 1.6.0

> 1.增加数据同步日志  

##### Version 1.5.9

> 1.修复bug  

##### Version 1.5.8

> 1.增加错误吗提示  

##### Version 1.5.7

> 1.增加H5查看用户状态，修改   

##### Version 1.5.6

> 1.增加H5查看用户状态，修改   

##### Version 1.5.5

> 1.更新vendor  

##### Version 1.5.4

> 1.增加H5查看用户状态，改成GET   

##### Version 1.5.3

> 1.增加同步接口IP白名单  

##### Version 1.5.2

> 1.增加H5查看用户状态   

##### Version 1.5.1

> 1.升级vendor  

##### Version 1.5.0

> 1.接入新的配置中心

##### Version 1.4.9

> 1.接入新的配置中心

##### Version 1.4.8

> 1.实时查询订购关系逻辑修改     

##### Version 1.4.7

> 1.MonitorPing   

##### Version 1.4.6

> 1.本地TW   

##### Version 1.4.5

> 1.vendor升级   

##### Version 1.4.4

> 1.IP同步SQL逻辑修改   

##### Version 1.4.3

> 1.增加数据同步白名单IP  

##### Version 1.4.2

> 1.增加message提示  

##### Version 1.4.1

> 1.增加message提示  

##### Version 1.4.0

> 1.增加错误message  

##### Version 1.3.9

> 1.增加ecode  

##### Version 1.3.8

> 1.判断是否是联通IP  
> 2.增加特权礼包接口  

##### Version 1.3.7

> 1.接入平滑发布  

##### Version 1.3.6

> 1.增加spid对应卡的类型   

##### Version 1.3.5

> 1.vendor升级   

##### Version 1.3.4

> 1.增加错误码返回   

##### Version 1.3.3

> 1.增加日志信息   

##### Version 1.3.2

> 1.ordertype int改成string   

##### Version 1.3.1

> 1.增加返回字段spid   

##### Version 1.3.0

> 1.升级go-business   

##### Version 1.2.9

> 1.修改项目上报  

##### Version 1.2.8

> 1.修改联通数据同步缓存    

##### Version 1.2.7

> 1.vendor升级   
> 2.增加过期判断   
> 3.增加返回字段   
> 4.增加用户状态接口   

##### Version 1.2.6

> 1.ci更新服务镜像   

##### Version 1.2.5

> 1.配置文件支持本地读取  

##### Version 1.2.4

> 1.删除Inner  
> 2.升级vendor     

##### Version 1.2.3

> 1.删除无用的接口      

##### Version 1.2.2

> 1.vendor升级   

##### Version 1.2.1

> 1.vendor升级   

##### Version 1.2.0

> 1.vendor升级   

##### Version 1.1.9

> 1.Success int改成string  

##### Version 1.1.8

> 1.shike接口改成post请求  
> 2.vendor升级  

##### Version 1.1.7

> 1.增加联通IP      

##### Version 1.1.6

> 1.联通IP同步接口Ipbegion改成ipbegin   
> 2.增加联通IP      

##### Version 1.1.5

> 1.联通IP同步接口开始ip和结束ip字段改成string   

##### Version 1.1.4

> 1.联通IP同步接口   

##### Version 1.1.3

> 1.联通流量接口请求改成post json     

##### Version 1.1.2

> 1.usermob解密  

##### Version 1.1.1

> 1.推广渠道编号可以为空  

##### Version 1.1.0

> 1.联通流量包接口对接  

##### Version 1.0.1

> 1.修复BUG  
> 2.联通信息同步接口  

##### Version 1.0.0

> 1.初始化项目  