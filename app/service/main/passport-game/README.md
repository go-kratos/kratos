## passport-game-service

#### 项目简介
> 1.游戏云账号服务

#### 编译环境
> 请只用golang v1.8.x以上版本编译执行。

#### 依赖包
> 1.公共包go-common

#### 特别说明
> 1.model目录可能会被其他项目引用，请谨慎更改并通知各方。

#### 部署注意
> 1.确定配置文件中 httpClient 节点中的 key 在 aso_apps 表中存在，且 issystem字段为 1。
> 2.确定使用上述 key 访问源站 /api/login 接口返回正常（非 -1、-403、-105）。
> 3.确定使用上述 key 访问源站 /api/login/oauth、/api/login/renewToken 接口返回正常（非 -1、-403）
> 2.确定配置文件中 db.cloud 节点中所指向的数据库的 app 表包含了 identify-game-service 的 appkey。

#### 部署检查
```
ssh root@online-host
# expected -629
curl 'http://passport.bilibili.com/api/login?appkey=868fb9ea57619022&ts=1514374489&sign=12abbf6d77076cbb6627d4b8ecea1e43'
# expected -101
curl 'http://passport.bilibili.com/api/oauth?appkey=868fb9ea57619022&ts=1514374489&sign=12abbf6d77076cbb6627d4b8ecea1e43'
# expected -101
curl 'http://passport.bilibili.com/api/login/renewToken?appkey=868fb9ea57619022&ts=1514374489&sign=12abbf6d77076cbb6627d4b8ecea1e43'

# expected -101
curl 'http://api.bilibili.co/x/internal/identify-game/oauth?appkey=868fb9ea57619022&ts=1514374489&sign=12abbf6d77076cbb6627d4b8ecea1e43'
# expected -101
curl 'http://api.bilibili.co/x/internal/identify-game/renewtoken?appkey=868fb9ea57619022&ts=1514374489&sign=12abbf6d77076cbb6627d4b8ecea1e43'
```

#### 接入注意
> 1.使用本服务，key 接口和 login 接口需要配套使用，如果用来加密密码的公钥与 key 接口返回的公钥不相同，会返回错误 -500。
