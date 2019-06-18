# 背景

我们需要统一的cache包，用于进行各类缓存操作。

# 概览

* 缓存操作均使用连接池，保证较快的数据读写速度且提高系统的安全可靠性。

# Memcache

提供protobuf，gob，json序列化方式，gzip的memcache接口

[memcache模块说明](cache-mc.md)

# Redis

提供redis操作的各类接口以及各类将redis server返回值转换为golang类型的快捷方法。

[redis模块说明](cache-redis.md)

-------------

[文档目录树](summary.md)
