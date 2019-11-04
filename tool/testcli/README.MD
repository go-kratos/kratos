## testcli UT运行环境构建工具
基于 docker-compose 实现跨平台跨语言环境的容器依赖管理方案，以解决运行ut场景下的 (mysql, redis, mc)容器依赖问题。

*这个是testing/lich的二进制工具版本（Go请直接使用库版本：github.com/bilibili/kratos/pkg/testing/lich)*

### 功能和特性
- 自动读取 test 目录下的 yaml 并启动依赖
- 自动导入 test 目录下的 DB 初始化 SQL
- 提供特定容器内的 healthcheck (mysql, mc, redis)
- 提供一站式解决 UT 服务依赖的工具版本 (testcli)

### 编译安装
*使用本工具/库需要前置安装好 docker & docker-compose@v1.24.1^*

#### Method 1. With go get
```shell
go get -u github.com/bilibili/kratos/tool/testcli
$GOPATH/bin/testcli -h
```
#### Method 2. Build with Go
```shell
cd github.com/bilibili/kratos/tool/testcli
go build -o $GOPATH/bin/testcli
$GOPATH/bin/testcli -h
```
#### Method 3. Import with Kratos pkg
```Go
import "github.com/bilibili/kratos/pkg/testing/lich"
```

### 构建数据
#### Step 1. create docker-compose.yml
创建依赖服务的 docker-compose.yml，并把它放在项目路径下的 test 文件夹下面。例如：
```shell
mkdir -p $YOUR_PROJECT/test
```
```yaml
version: "3.7"
 
services:
  db:
    image: mysql:5.6
    ports:
    - 3306:3306
    environment:
    - MYSQL_ROOT_PASSWORD=root
    volumes:
    - .:/docker-entrypoint-initdb.d
    command: [
      '--character-set-server=utf8',
      '--collation-server=utf8_unicode_ci'
    ]
 
  redis:
    image: redis
    ports:
      - 6379:6379
```
一般来讲，我们推荐在项目根目录创建 test 目录，里面存放描述服务的yml，以及需要初始化的数据（database.sql等）。

同时也需要注意，正确的对容器内服务进行健康检测，testcli会在容器的health状态执行UT，其实我们也内置了针对几个较为通用镜像（mysql mariadb mc redis）的健康检测，也就是不写也没事(^^;;

#### Step 2. export database.sql
构造初始化的数据（database.sql等），当然也把它也在 test 文件夹里。
```sql
CREATE DATABASE IF NOT EXISTS `YOUR_DATABASE_NAME`;
 
SET NAMES 'utf8';
USE `YOUR_DATABASE_NAME`;
 
CREATE TABLE IF NOT EXISTS `YOUR_TABLE_NAME` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键'，
  PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='YOUR_TABLE_NAME';
```
这里需要注意，在创建库/表的时候尽量加上 IF NOT EXISTS，以给予一定程度的容错，以及 SET NAMES 'utf8'; 用于解决客户端连接乱码问题。

#### Step 3. change your project mysql config
```toml
[mysql]
    addr = "127.0.0.1:3306"
    dsn = "root:root@tcp(127.0.0.1:3306)/YOUR_DATABASE?timeout=1s&readTimeout=1s&writeTimeout=1s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
    active = 20
    idle = 10
    idleTimeout ="1s"
    queryTimeout = "1s"
    execTimeout = "1s"
    tranTimeout = "1s"
```
在 *Step 1* 我们已经指定了服务对外暴露的端口为3306（这当然也可以是你指定的任何值），那理所应当的我们也要修改项目连接数据库的配置～

Great! 至此你已经完成了运行所需要用到的数据配置，接下来就来运行它。

### 运行
开头也说过本工具支持两种运行方式：testcli 二进制工具版本和 go package 源码包，业务方可以根据需求场景进行选择。
#### Method 1. With testcli tool
*已支持的 flag： -f，--nodown，down，run*
- -f，指定 docker-compose.yaml 文件路径，默认为当前目录下。
- --nodown，指定是否在UT执行完成后保留容器，以供下次复用。
- down，teardown 销毁当前项目下这个 compose 文件产生的容器。
- run，运行你当前语言的单测执行命令（如：golang为 go test -v ./） 

example:
```shell
testcli -f ../../test/docker-compose.yaml run go test -v ./
```
#### Method 2. Import with Kratos pkg
- Step1. 在 Dao|Service 层中的 TestMain 单测主入口中，import  "github.com/bilibili/kratos/pkg/testing/lich" 引入testcli工具的go库版本。
- Step2. 使用  flag.Set("f", "../../test/docker-compose.yaml") 指定 docker-compose.yaml 文件的路径。
- Step3. 在 flag.Parse() 后即可使用 lich.Setup() 安装依赖&初始化数据（注意测试用例执行结束后 lich.Teardown() 回收下～）
- Step4. 运行 `go test -v ./ `看看效果吧～

example:
```Go
package dao
 
 
import (
    "flag"
    "os"
    "strings"
    "testing"
 
    "github.com/bilibili/kratos/pkg/conf/paladin"
    "github.com/bilibili/kratos/pkg/testing/lich"
 )
 
var (
    d *Dao
)
 
func TestMain(m *testing.M) {
    flag.Set("conf", "../../configs")
    flag.Set("f", "../../test/docker-compose.yaml")
    flag.Parse()
    if err := paladin.Init(); err != nil {
        panic(err)
    }
    if err := lich.Setup(); err != nil {
        panic(err)
    }
    defer lich.Teardown()
    d = New()
    if code := m.Run(); code != 0 {
        panic(code)
    }
}
 ```
## 注意
因为启动mysql容器较为缓慢，健康检测的机制会重试3次，每次暂留5秒钟，基本在10s内mysql就能从creating到服务正常启动！

当然你也可以在使用 testcli 时加上 --nodown，使其不用每次跑都新建容器，只在第一次跑的时候会初始化容器，后面都进行复用，这样速度会快很多。

成功启动后就欢乐奔放的玩耍吧～ Good Lucky!