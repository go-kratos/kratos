## 单元测试辅助工具
在单元测试中，我们希望每个测试用例都是独立的。这时候就需要Stub, Mock, Fakes等工具来帮助我们进行用例和依赖之间的隔离。

同时通过对错误情况的 Mock 也可以帮我们检查代码多个分支结果，从而提高覆盖率。

以下工具已加入到 Kratos 框架 go modules，可以借助 testgen 代码生成器自动生成部分工具代码，请放心食用。更多使用方法还欢迎大家多多探索。

### GoConvey
GoConvey是一套针对golang语言的BDD类型的测试框架。提供了良好的管理和执行测试用例的方式，包含丰富的断言函数，而且同时有测试执行和报告Web界面的支持。

#### 使用特性
为了更好的使用 GoConvey 来编写和组织测试用例，需要注意以下几点特性：

1. Convey方法和So方法的使用
> - Convey方法声明了一种规格的组织，每个组织内包含一句描述和一个方法。在方法内也可以嵌套其他Convey语句和So语句。
```Go
// 顶层Convey方法，需引入*testing.T对象
Convey(description string, t *testing.T, action func())

// 其他嵌套Convey方法，无需引入*testing.T对象
Convey(description string, action func())
```
注：同一Scope下的Convey语句描述不可以相同！
> - So方法是断言方法，用于对执行结果进行比对。GoConvey官方提供了大量断言，同时也可以自定义自己的断言（[戳这里了解官方文档](https://github.com/smartystreets/goconvey/wiki/Assertions)）
```Go
// A=B断言
So(A, ShouldEqual, B)
 
// A不为空断言
So(A, ShouldNotBeNil)
```

2. 执行次序
> 假设有以下Convey伪代码，执行次序将为A1B2A1C3。将Convey方法类比树的结点的话，整体执行类似树的遍历操作。  
> 所以Convey A部分可在组织测试用例时，充当“Setup”的方法。用于初始化等一些操作。
```Go
Convey伪代码
Convey A
    So 1
    Convey B
        So 2
    Convey C
        So 3
```

3. Reset方法
> GoConvey提供了Reset方法来进行“Teardown”的操作。用于执行完测试用例后一些状态的回收，连接关闭等操作。Reset方法不可与顶层Convey语句在同层。
```Go 
// Reset
Reset func(action func())
```
假设有以下带有Reset方法的伪代码，同层Convey语句执行完后均会执行同层的Reset方法。执行次序为A1B2C3EA1D4E。
```Go
Convey A
    So 1
    Convey B
        So 2
        Convey C
            So 3
    Convey D
        So 4
    Reset E
```

4. 自然语言逻辑到测试用例的转换
> 在了解了Convey方法的特性和执行次序后，我们可以通过这些性质把对一个方法的测试用例按照日常逻辑组织起来。尤其建议使用Given-When-Then的形式来组织
> - 比较直观的组织示例
```Go
Convey("Top-level", t, func() {
 
    // Setup 工作，在本层内每个Convey方法执行前都会执行的部分：
    db.Open()
    db.Initialize()
 
    Convey("Test a query", func() {
        db.Query()
        // TODO: assertions here
    })
 
    Convey("Test inserts", func() {
        db.Insert()
        // TODO: assertions here
    })
 
    Reset(func() {
        // Teardown工作，在本层内每个Convey方法执行完后都会执行的部分：
        db.Close()
    })
 
})
```
> - 定义单独的包含Setup和Teardown的帮助方法
```Go
package main
 
import (
    "database/sql"
    "testing"
 
    _ "github.com/lib/pq"
    . "github.com/smartystreets/goconvey/convey"
)
 
// 帮助方法，将原先所需的处理方法以参数（f）形式传入
func WithTransaction(db *sql.DB, f func(tx *sql.Tx)) func() {
    return func() {
        // Setup工作
        tx, err := db.Begin()
        So(err, ShouldBeNil)
 
        Reset(func() {
            // Teardown工作
            /* Verify that the transaction is alive by executing a command */
            _, err := tx.Exec("SELECT 1")
            So(err, ShouldBeNil)
 
            tx.Rollback()
        })
 
        // 调用传入的闭包做实际的事务处理
        f(tx)
    }
}
 
func TestUsers(t *testing.T) {
    db, err := sql.Open("postgres", "postgres://localhost?sslmode=disable")
    if err != nil {
        panic(err)
    }
 
    Convey("Given a user in the database", t, WithTransaction(db, func(tx *sql.Tx) {
        _, err := tx.Exec(`INSERT INTO "Users" ("id", "name") VALUES (1, 'Test User')`)
        So(err, ShouldBeNil)
 
        Convey("Attempting to retrieve the user should return the user", func() {
             var name string
 
             data := tx.QueryRow(`SELECT "name" FROM "Users" WHERE "id" = 1`)
             err = data.Scan(&name)
 
             So(err, ShouldBeNil)
             So(name, ShouldEqual, "Test User")
        })
    }))
}
```

#### 使用建议
强烈建议使用 [testgen](https://github.com/bilibili/kratos/blob/master/doc/wiki-cn/ut-testgen.md) 进行测试用例的生成，生成后每个方法将包含一个符合以下规范的正向用例。

用例规范：
1. 每个方法至少包含一个测试方法（命名为Test[PackageName][FunctionName]）
2. 每个测试方法包含一个顶层Convey语句，仅在此引入admin *testing.T类型的对象，在该层进行变量声明。
3. 每个测试方法不同的用例用Convey方法组织
4. 每个测试用例的一组断言用一个Convey方法组织
5. 使用convey.C保持上下文一致

### MonkeyPatching

#### 特性和使用条件
1. Patch()对任何无接收者的方法均有效
2. PatchInstanceMethod()对有接收者的包内/私有方法无法工作（因使用到了反射机制）。可以采用给私有方法的下一级打补丁，或改为无接收者的方法，或将方法转为公有

#### 适用场景（建议）
项目代码中上层对下层包依赖时，下层包方法Mock（例如service层对dao层方法依赖时）
基础库（MySql, Memcache, Redis）错误Mock
其他标准库，基础库以及第三方包方法Mock

#### 使用示例
1. 上层包对下层包依赖示例
Service层对Dao层依赖：
```GO
// 原方法
func (s *Service) realnameAlipayApply(c context.Context, mid int64) (info *model.RealnameAlipayApply, err error) {
    if info, err = s.mbDao.RealnameAlipayApply(c, mid); err != nil {
        return
    }
    ...
    return
}
  
// 测试方法
func TestServicerealnameAlipayApply(t *testing.T) {
    convey.Convey("realnameAlipayApply", t, func(ctx convey.C) {
        ...
        ctx.Convey("When everything goes positive", func(ctx convey.C) {
            guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.mbDao), "RealnameAlipayApply", func(_ *dao.Dao, _ context.Context, _ int64) (*model.RealnameAlipayApply, error) {
                return nil, nil
            })
            defer guard.Unpatch()
            info, err := s.realnameAlipayApply(c, mid)
            ctx.Convey("Then err should be nil,info should not be nil", func(ctx convey.C) {
                ctx.So(info, convey.ShouldNotBeNil)
                ctx.So(err, convey.ShouldBeNil)
            })
        })
    })
}
```
2. 基础库错误Mock示例
```Go
  
// 原方法（部分）
func (d *Dao) BaseInfoCache(c context.Context, mid int64) (info *model.BaseInfo, err error) {
    ...
    conn := d.mc.Get(c)
    defer conn.Close()
    item, err := conn.Get(key)
    if err != nil {
        log.Error("conn.Get(%s) error(%v)", key, err)
        return
    }
    ...
    return
}
 
 
// 测试方法（错误Mock部分）
func TestDaoBaseInfoCache(t *testing.T) {
    convey.Convey("BaseInfoCache", t, func(ctx convey.C) {
        ...
        Convey("When conn.Get gets error", func(ctx convey.C) {
            guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.mc), "Get", func(_ *memcache.Pool, _ context.Context) memcache.Conn {
                return memcache.MockWith(memcache.ErrItemObject)
            })
            defer guard.Unpatch()
            _, err := d.BaseInfoCache(c, mid)
            ctx.Convey("Error should be equal to memcache.ErrItemObject", func(ctx convey.C) {
                ctx.So(err, convey.ShouldEqual, memcache.ErrItemObject)
            })
        })
    })
}
```
#### 注意事项
- Monkey非线程安全
- Monkey无法针对Inline方法打补丁，在测试时可以使用go test -gcflags=-l来关闭inline编译的模式（一些简单的go inline介绍戳这里）
- Monkey在一些面向安全不允许内存页写和执行同时进行的操作系统上无法工作
- 更多详情请戳：https://github.com/bouk/monkey



### Gock——HTTP请求Mock工具

#### 特性和使用条件

#### 工作原理
1. 截获任意通过 http.DefaultTransport或者自定义http.Transport对外的http.Client请求
2. 以“先进先出”原则将对外需求和预定义好的HTTP Mock池中进行匹配
3. 如果至少一个Mock被匹配，将按照2中顺序原则组成Mock的HTTP返回
4. 如果没有Mock被匹配，若实际的网络可用，将进行实际的HTTP请求。否则将返回错误

#### 特性
- 内建帮助工具实现JSON/XML简单Mock
- 支持持久的、易失的和TTL限制的Mock
- 支持HTTP Mock请求完整的正则表达式匹配
- 可通过HTTP方法，URL参数，请求头和请求体匹配
- 可扩展和可插件化的HTTP匹配规则
- 具备在Mock和实际网络模式之间切换的能力
- 具备过滤和映射HTTP请求到正确的Mock匹配的能力
- 支持映射和过滤可以更简单的掌控Mock
- 通过使用http.RoundTripper接口广泛兼容HTTP拦截器
- 可以在任意net/http兼容的Client上工作
- 网络延迟模拟（beta版本）
- 无其他依赖

#### 适用场景（建议）
任何需要进行HTTP请求的操作，建议全部用Gock进行Mock，以减少对环境的依赖。

使用示例：
1. net/http 标准库 HTTP 请求Mock
```Go
import  gock "gopkg.in/h2non/gock.v1"
  
// 原方法
 func (d *Dao) Upload(c context.Context, fileName, fileType string, expire int64, body io.Reader) (location string, err error) {
    ...
    resp, err = d.bfsClient.Do(req) //d.bfsClient类型为*http.client
    ...
    if resp.StatusCode != http.StatusOK {
        ...
    }
    header = resp.Header
    code = header.Get("Code")
    if code != strconv.Itoa(http.StatusOK) {
        ...
    }
    ...
    return
}
 
 
// 测试方法
func TestDaoUpload(t *testing.T) {
    convey.Convey("Upload", t, func(ctx convey.C) {
        ...
        // d.client 类型为 *http.client  根据Gock包描述需要设置http.Client的Transport情况。也可在TestMain中全局设置，则所有的HTTP请求均通过Gock来解决
        d.client.Transport = gock.DefaultTransport // ！注意：进行httpMock前需要对http 请求进行拦截，否则Mock失败
        // HTTP请求状态和Header都正确的Mock
        ctx.Convey("When everything is correct", func(ctx convey.C) {
            httpMock("PUT", url).Reply(200).SetHeaders(map[string]string{
                "Code":     "200",
                "Location": "SomePlace",
            })
            location, err := d.Upload(c, fileName, fileType, expire, body)
            ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
                ctx.So(err, convey.ShouldBeNil)
                ctx.So(location, convey.ShouldNotBeNil)
            })
        })
        ...
        // HTTP请求状态错误Mock
        ctx.Convey("When http request status != 200", func(ctx convey.C) {
            d.client.Transport = gock.DefaultTransport
            httpMock("PUT", url).Reply(404)
            _, err := d.Upload(c, fileName, fileType, expire, body)
            ctx.Convey("Then err should not be nil", func(ctx convey.C) {
                ctx.So(err, convey.ShouldNotBeNil)
            })
        })
        // HTTP请求Header中Code值错误Mock
        ctx.Convey("When http request Code in header != 200", func(ctx convey.C) {
            d.client.Transport = gock.DefaultTransport
            httpMock("PUT", url).Reply(404).SetHeaders(map[string]string{
                "Code":     "404",
                "Location": "SomePlace",
            })
            _, err := d.Upload(c, fileName, fileType, expire, body)
            ctx.Convey("Then err should not be nil", func(ctx convey.C) {
                ctx.So(err, convey.ShouldNotBeNil)
            })
        })
  
        // 由于同包内有其他进行实际HTTP请求的测试。所以再每次用例结束后，进行现场恢复（关闭Gock设置默认的Transport）
        ctx.Reset(func() {
           gock.OffAll()
           d.client.Transport = http.DefaultClient.Transport
        })
 
 
    })
}
  
func httpMock(method, url string) *gock.Request {
    r := gock.New(url)
    r.Method = strings.ToUpper(method)
    return r
}
```
2. blademaster库HTTP请求Mock
```Go
// 原方法
func (d *Dao) SendWechatToGroup(c context.Context, chatid, msg string) (err error) {
    ...
    if err = d.client.Do(c, req, &res); err != nil {
        ...
    }
    if res.Code != 0 {
        ...
    }
    return
}
  
// 测试方法
func TestDaoSendWechatToGroup(t *testing.T) {
    convey.Convey("SendWechatToGroup", t, func(ctx convey.C) {
        ...
        // 根据Gock包描述需要设置bm.Client的Transport情况。也可在TestMain中全局设置，则所有的HTTP请求均通过Gock来解决。
        // d.client 类型为 *bm.client
        d.client.SetTransport(gock.DefaultTransport) // ！注意：进行httpMock前需要对http 请求进行拦截，否则Mock失败
        // HTTP请求状态和返回内容正常Mock
        ctx.Convey("When everything gose postive", func(ctx convey.C) {
            httpMock("POST", _sagaWechatURL+"/appchat/send").Reply(200).JSON(`{"code":0,"message":"0"}`)
            err := d.SendWechatToGroup(c, d.c.WeChat.ChatID, msg)
            ...
        })
        // HTTP请求状态错误Mock
        ctx.Convey("When http status != 200", func(ctx convey.C) {
            httpMock("POST", _sagaWechatURL+"/appchat/send").Reply(404)
            err := d.SendWechatToGroup(c, d.c.WeChat.ChatID, msg)
            ...
        })
        // HTTP请求返回值错误Mock
        ctx.Convey("When http response code != 0", func(ctx convey.C) {
            httpMock("POST", _sagaWechatURL+"/appchat/send").Reply(200).JSON(`{"code":-401,"message":"0"}`)
            err := d.SendWechatToGroup(c, d.c.WeChat.ChatID, msg)
            ...
        })
        // 由于同包内有其他进行实际HTTP请求的测试。所以再每次用例结束后，进行现场恢复（关闭Gock设置默认的Transport）。
        ctx.Reset(func() {
            gock.OffAll()
            d.client.SetTransport(http.DefaultClient.Transport)
        })
    })
}
 
func httpMock(method, url string) *gock.Request {
    r := gock.New(url)
    r.Method = strings.ToUpper(method)
    return r
}
```

#### 注意事项
- Gock不是完全线程安全的
- 如果执行并发代码，在配置Gock和解释定制的HTTP clients时，要确保Mock已经事先声明好了来避免不需要的竞争机制
- 更多详情请戳：https://github.com/h2non/gock


### GoMock

#### 使用条件
只能对公有接口（interface）定义的代码进行Mock，并仅能在测试过程中进行

#### 使用方法
- 官方安装使用步骤
```shell
## 获取GoMock包和自动生成Mock代码工具mockgen
go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen
  
## 生成mock文件
## 方法1：生成对应文件下所有interface
mockgen -source=path/to/your/interface/file.go
  
## 方法2：生成对应包内指定多个interface，并用逗号隔开
mockgen database/sql/driver Conn,Driver
  
## 示例：
mockgen -destination=$GOPATH/kratos/app/xxx/dao/dao_mock.go -package=dao kratos/app/xxx/dao DaoInterface
```
- testgen 使用步骤（GoMock生成功能已集成在Creater工具中，无需额外安装步骤即可直接使用）
```shell
## 直接给出含有接口类型定义的包路径，生成Mock文件将放在包目录下一级mock/pkgName_mock.go中
./creater --m mock absolute/path/to/your/pkg
```
- 测试代码内使用方法
```Go
// 测试用例内直接使用
// 需引入的包
import (
    ...
    "github.com/otokaze/mock/gomock"
    ...
)
  
func TestPkgFoo(t *testing.T) {
    convey.Convey("Foo", t, func(ctx convey.C) {
        ...       
        ctx.Convey("Mock Interface to test", func(ctx convey.C) {
            // 1. 使用gomock.NewController新增一个控制器
            mockCtrl := gomock.NewController(t)
            // 2. 测试完成后关闭控制器
            defer mockCtrl.Finish()
            // 3. 以控制器为参数生成Mock对象
            yourMock := mock.NewMockYourClient(mockCtrl)
            // 4. 使用Mock对象替代原代码中的对象
            yourClient = yourMock
            // 5. 使用EXPECT().方法名(方法参数).Return(返回值)来构造所需输入/输出
            yourMock.EXPECT().YourMethod(gomock.Any()).Return(nil)
            res:= Foo(params)
            ...
        })
        ...
    })
}
 
// 可以利用Convey执行顺序方式适当调整以简化代码
func TestPkgFoo(t *testing.T) {
    convey.Convey("Foo", t, func(ctx convey.C) {
        ...
        mockCtrl := gomock.NewController(t)
        yourMock := mock.NewMockYourClient(mockCtrl) 
        ctx.Convey("Mock Interface to test1", func(ctx convey.C) {
            yourMock.EXPECT().YourMethod(gomock.Any()).Return(nil)
            ...
        })
        ctx.Convey("Mock Interface to test2", func(ctx convey.C) {
            yourMock.EXPECT().YourMethod(args).Return(res)
            ...
        })
        ...
        ctx.Reset(func(){
            mockCtrl.Finish()
        })
    })
}
```

#### 适用场景（建议）
1. gRPC中的Client接口
2. 也可改造现有代码构造Interface后使用（具体可配合Creater的功能进行Interface和Mock的生成）
3. 任何对接口中定义方法依赖的场景 

#### 注意事项 
- 如有Mock文件在包内，在执行单元测试时Mock代码会被识别进行测试。请注意Mock文件的放置。
- 更多详情请戳：https://github.com/golang/mock