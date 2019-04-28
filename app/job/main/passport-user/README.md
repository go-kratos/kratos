## passport-user

#### 项目简介
> 1.数据同步job，主要用于同步全量历史数据和订阅源站变更的数据

#### 编译环境
> 请只用golang v1.8.x以上版本编译执行。

#### 依赖包
> 1.公共包go-common

#### 数据同步流程
> 1.user_secret添加aes key
> 2.country_code导入数据
> 3.全量同步aso_account
    全量同步aso_account_sns
> 4.全量同步aso_account_info
    全量同步aso_tel_bind_log
> 5.全量同步aso_account_reg
> 6.增量同步
