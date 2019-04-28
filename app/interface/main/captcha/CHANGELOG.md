### captcha 图片验证码服务

#### V1.0.6
> 1.修复http层接口返回多json.
> 2.修复map并发锁.

#### V1.0.5
> 1.使用新的verify
> 2.使用Toml2()

#### V1.0.4
> 1.迁移项目至main目录
> 2.修改bind使用方法

#### V1.0.3
> 1.将captcha client 集成于captcha-interface，接口返回url

#### V1.0.2
> 1.将/x/v1/captcha/get接口图片encode逻辑转移到service层提前处理
> 2.HTTP初始化时增加请求限速功能
> 3./x/internal/v1/captcha/token、/x/internal/v1/captcha/verify接口增加identify认证

#### V1.0.1
> 1.修改URI /x/internal/v1

#### V1.0.0
> 1.迁移大仓库 
> 2.使用新的HTTP框架
