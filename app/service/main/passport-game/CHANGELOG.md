## passport-game-service

#### Version 1.13.1
> 1. 修复myinfo bug

#### Version 1.13.0
> 1. 增加用户手机邮箱信息

#### Version 1.12.0
> 1. 增加支持重置密码短信发送

#### Version 1.11.0
> 1. 缓存来自 origin 的结果

#### Version 1.10.1
> 1. 增加ut

#### Version 1.9.1
> 1. add x-forward header

#### Version 1.9.0
> 1. 增加regv3接口

#### Version 1.8.0
> 1. add other region db
> 2. update RemoteIP

#### Version 1.7.0
> 1. bm && mv dir

#### Version 1.6.2
> 1.renew token on origin in priority

#### Version 1.6.1
> 1.fix password not matches when salt is empty

#### Version 1.6.0
> 1.remove cold 30s when back to origin password error
> 2.add info log when back to origin ok

#### Version 1.5.0
> 1.use origin public and private key

#### Version 1.4.0
> 1.generate new public key and private key instead of using origin's
> 2.remove account table usage
> 3.change perm, info model to pb
> 4.add pb cache reading APIs
> 5.reduce cache reading times of api myinfo and oauth
> 6.add info API

#### Version 1.3.0
> 1.move to kratos
> 2.add oauth and renewToken dispatch

#### Version 1.2.1
> 1.token proc bugfix: unmarshall `ctime` `mtime` error

#### Version 1.2.0
> 1.token 消费优化

#### Version 1.1.0
> 1.api login
> 2.api get key

#### Version 1.0.0
> 1.基础api
