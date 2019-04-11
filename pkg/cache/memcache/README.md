# cache/memcache

##### 项目简介
1. 提供protobuf，gob，json序列化方式，gzip的memcache接口

#### 使用方式
```golang
// 初始化 注意这里只是示例 展示用法 不能每次都New 只需要初始化一次
mc := memcache.New(&memcache.Config{})
// 程序关闭的时候调用close方法
defer mc.Close()
// 增加 key
err = mc.Set(c, &memcache.Item{})
// 删除key
err := mc.Delete(c,key)
// 获得某个key的内容
err := mc.Get(c,key).Scan(&v)
// 获取多个key的内容
replies, err := mc.GetMulti(c, keys)
for _, key := range replies.Keys() {
   if err = replies.Scan(key, &v); err != nil {
       return 
    }
}
```
