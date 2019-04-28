# appstatic-job

### v1.0.2
1. 刷完app-resource的GRPC后sleep 2秒

### v1.0.1
1. 新增接broadcast逻辑：
* 拆分dao层，分出cal_diff（增量包计算）和push（请求broadcast推送）两个dao的包出来
* 对接app-resource，计算增量包完成后请求app-resource刷新缓存，成功后再请求broadcast推送

### v1.0.0
1. 项目初始化，从appstatic-admin中迁移出增量包计算逻辑

