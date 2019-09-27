## testgen UT代码自动生成器
解放你的双手，让你的UT一步到位！

### 功能和特性
- 支持生成 Dao|Service 层UT代码功能（每个方法包含一个正向用例）
- 支持生成 Dao|Service 层测试入口文件dao_test.go, service_test.go（用于控制初始化，控制测试流程等）
- 支持生成Mock代码（使用GoMock框架）
- 支持选择不同模式生成不同代码（使用"–m mode"指定）
- 生成单元测试代码时，同时支持传入目录或文件
- 支持指定方法追加生成测试用例（使用"–func funcName"指定）

### 编译安装
#### Method 1. With go get
```shell
go get -u github.com/bilibili/kratos/tool/testgen
$GOPATH/bin/testgen -h
```
#### Method 2. Build with Go
```shell
cd github.com/bilibili/kratos/tool/testgen
go build -o $GOPATH/bin/testgen
$GOPATH/bin/testgen -h
```
### 运行
#### 生成Dao/Service层单元UT
```shell
$GOPATH/bin/testgen YOUR_PROJECT/dao # default mode 
$GOPATH/bin/testgen --m test path/to/your/pkg
$GOPATH/bin/testgen --func functionName path/to/your/pkg
```

#### 生成接口类型
```shell
$GOPATH/bin/testgen --m interface YOUR_PROJECT/dao #当前仅支持传目录，如目录包含子目录也会做处理
```

#### 生成Mock代码
 ```shell
$GOPATH/bin/testgen --m mock YOUR_PROJECT/dao #仅传入包路径即可
```

#### 生成Monkey代码
```shell
$GOPATH/bin/testgen --m monkey yourCodeDirPath #仅传入包路径即可
```
### 赋诗一首
```
莫生气 莫生气
代码辣鸡非我意 
自己动手分田地
谈笑风生活长命
```