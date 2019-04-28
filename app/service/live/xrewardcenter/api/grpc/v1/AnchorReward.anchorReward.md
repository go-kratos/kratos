## (主播侧)-我的主播奖励(登录态)

`POST http://api.live.bilibili.com/xlive/internal/xrewardcenter/v1/anchorReward/myReward`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer| 页数|
|uid|是|integer| 主播uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": [
            {
                //  id
                "id": 0,
                //  奖励类型 1:ss推荐卡 2:s推荐卡、任意门
                "reward_type": 0,
                //  1:未使用,3:已使用,5:已过期
                "status": 0,
                //  奖励id
                "reward_id": 0,
                //  奖励名称
                "name": "",
                //  奖励图标
                "icon": "",
                //  获得时间datetime
                "achieve_time": "",
                //  过期时间datetime
                "expire_time": "",
                //  来源,1:主播任务,2:小时榜
                "source": 0,
                //  奖励简介
                "reward_intro": ""
            }
        ],
        "page": {
            //  当前页码
            "page": 0,
            //  每页大小
            "page_size": 0,
            //  总页数
            "total_page": 0,
            //  总记录数
            "total_count": 0
        },
        //  过期奖励数量
        "expire_count": 0
    }
}
```

##* (主播侧)-奖励使用记录(登录态)
 

`GET http://api.live.bilibili.com/xlive/xrewardcenter/v1/anchorReward/useRecord`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer| 页数|
|uid|是|integer| 主播uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": [
            {
                //  id
                "id": 0,
                //  奖励id
                "reward_id": 0,
                //  1:未使用,3:已使用,5:已过期
                "status": 0,
                //  奖励名称
                "name": "",
                //  奖励图标
                "icon": "",
                //  获得时间datetime
                "achieve_time": "",
                //  获得时间datetime
                "use_time": "",
                //  过期时间datetime
                "expire_time": "",
                //  来源,1:主播任务,2:小时榜
                "source": 0,
                //  奖励简介
                "reward_intro": ""
            }
        ],
        "page": {
            //  当前页码
            "page": 0,
            //  每页大小
            "page_size": 0,
            //  总页数
            "total_page": 0,
            //  总记录数
            "total_count": 0
        }
    }
}
```

##* (主播侧)-使用奖励(登录态)
 

`GET http://api.live.bilibili.com/xlive/xrewardcenter/v1/anchorReward/useReward`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer| 奖励列表id|
|uid|是|integer| 主播uid|
|use_plat|是|string| 使用平台|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "result": 0
    }
}
```

##* (主播侧)-奖励和任务红点(登录态)
 

`GET http://api.live.bilibili.com/xlive/xrewardcenter/v1/anchorReward/isViewed`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| 主播uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否展示任务红点
        "task_should_notice": 0,
        //  是否展示奖励入口
        "show_reward_entry": 0,
        //  是否展示奖励红点
        "reward_should_notice": 0,
        //  任务状态, 0:没有,1:领取, 5:完成
        "task_status": 0,
        //  是否在首页黑名单中
        "is_blacked": 0,
        //  点击跳转h5链接
        "url": ""
    }
}
```

## (主播侧)-添加主播奖励(内部接口)

`POST http://api.live.bilibili.com/xlive/internal/xrewardcenter/v1/anchorReward/addReward`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|reward_id|是|integer| 奖励id, 1:任意门|
|roomid|是|integer| 房间号|
|source|是|integer| 来源,1:主播任务,2:小时榜,3:周星|
|uid|是|integer| 主播uid|
|order_id|是|string| 流水单号|
|lifespan|否|integer| 礼物有效期小时数,默认24小时|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否成功
        "result": 0
    }
}
```

