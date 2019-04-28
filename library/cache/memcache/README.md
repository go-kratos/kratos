# go-common/cache/memcache

##### 项目简介
> 1. 提供protobuf，gob，json序列化方式，gzip的memcache接口

##### 编译环境
> 1. 请只用golang v1.7.x以上版本编译执行。

##### 测试
> 1. 执行当前目录下所有测试文件，测试所有功能

##### 特别说明
> 1. 使用protobuf需要在pb文件目录下运行business/make.sh脚本生成go文件才能使用

#### 使用方式
```golang
// 初始化
mc := memcache.New(&memcache.Config{})
// 增加 key
err = mc.Set(c, &memcache.Item{})
// 删除key
err := mc.Delete(c,key)
// 获得某个key的内容
err := mc.Get(c,key).Scan(&v)
// 获取多个key的内容
replies, err := mc.GetMulti(c, keys)
for _, key := range replies.Keys() {
   if err = rows.Scan(key, &v); err != nil {
       return 
    }
}
```