<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/capsule/get_coin_list](#xliveinternallive-adminv1capsuleget_coin_list) 
- [/xlive/internal/live-admin/v1/capsule/update_coin_config](#xliveinternallive-adminv1capsuleupdate_coin_config) 
- [/xlive/internal/live-admin/v1/capsule/update_coin_status](#xliveinternallive-adminv1capsuleupdate_coin_status) 
- [/xlive/internal/live-admin/v1/capsule/delete_coin](#xliveinternallive-adminv1capsuledelete_coin) 
- [/xlive/internal/live-admin/v1/capsule/get_pool_list](#xliveinternallive-adminv1capsuleget_pool_list) 
- [/xlive/internal/live-admin/v1/capsule/update_pool](#xliveinternallive-adminv1capsuleupdate_pool) 
- [/xlive/internal/live-admin/v1/capsule/delete_pool](#xliveinternallive-adminv1capsuledelete_pool) 
- [/xlive/internal/live-admin/v1/capsule/update_pool_status](#xliveinternallive-adminv1capsuleupdate_pool_status) 
- [/xlive/internal/live-admin/v1/capsule/get_pool_prize](#xliveinternallive-adminv1capsuleget_pool_prize) 
- [/xlive/internal/live-admin/v1/capsule/get_prize_type](#xliveinternallive-adminv1capsuleget_prize_type) 
- [/xlive/internal/live-admin/v1/capsule/get_prize_expire](#xliveinternallive-adminv1capsuleget_prize_expire) 
- [/xlive/internal/live-admin/v1/capsule/update_pool_prize](#xliveinternallive-adminv1capsuleupdate_pool_prize) 
- [/xlive/internal/live-admin/v1/capsule/delete_pool_prize](#xliveinternallive-adminv1capsuledelete_pool_prize) 
- [/xlive/internal/live-admin/v1/capsule/get_coupon_list](#xliveinternallive-adminv1capsuleget_coupon_list) 

## /xlive/internal/live-admin/v1/capsule/get_coin_list
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|是|integer|页码，从1开始|
|page_size|是|integer|页面的大小|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "total": 0,
        "total_page": 0,
        "list": [
            {
                // 扭蛋ID
                "id": 0,
                // 名称 普通扭蛋, 梦幻扭蛋
                "title": "",
                // 转化数量
                "change_num": 0,
                // 开始时间
                "start_time": "",
                // 结束时间
                "end_time": "",
                // 状态 0为下线，1为上线
                "status": 0,
                // 获得方式 1为所有瓜子道具，2为所有金瓜子道具，3为指定道具ID
                "gift_type": 0,
                //  道具的ID
                "gift_config": "",
                // 活动分区
                "area_ids": [
                    {
                        //  父分区ID
                        "parent_id": 0,
                        //  是否全选
                        "is_all": 0,
                        //  子分区ID
                        "list": [
                            0
                        ]
                    }
                ]
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/capsule/update_coin_config
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer|扭蛋ID|
|title|是|string|名称 普通扭蛋, 梦幻扭蛋|
|change_num|是|integer|转化数量|
|start_time|是|string|开始时间|
|end_time|是|string|结束时间|
|status|是|integer|状态 0为下线，1为上线|
|gift_type|是|integer|获得方式 1为所有瓜子道具，2为所有金瓜子道具，3为指定道具ID|
|gift_config|否|string| 道具的ID|
|area_ids|是|string|里面是父分区ID，是否全选，分区ID|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/update_coin_status
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer|扭蛋币id|
|status|否|integer|状态 0为下线，1为上线|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/delete_coin
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer|扭蛋币id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/get_pool_list
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|是|integer|页码|
|page_size|是|integer|页面的大小|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 总数
        "total": 0,
        // 总页数
        "total_page": 0,
        "list": [
            {
                //  奖池id
                "id": 0,
                //  扭蛋名称
                "coin_id": 0,
                //  奖池名称
                "title": "",
                //  奖池名称
                "coin_title": "",
                // 开始时间
                "start_time": "",
                // 结束时间
                "end_time": "",
                // 状态 0为下线，1为上线
                "status": 0,
                // 描述
                "rule": ""
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/capsule/update_pool
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer| 奖池id|
|coin_id|是|integer| 扭蛋名称|
|title|是|string|奖池名称|
|start_time|是|string|开始时间|
|end_time|是|string|结束时间|
|rule|是|string|描述|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/delete_pool
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer|奖池id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/update_pool_status
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer|奖池id|
|status|否|integer|状态 0为未上线，1为上线|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/get_pool_prize
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|pool_id|是|integer|奖池id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                // 奖励id
                "id": 0,
                // 奖池id
                "pool_id": 0,
                // 奖品类型 1为道具，2为头衔，3为经验原石，4为经验曜石，5为贤者之石，6为小号小电视，7为舰长守护，8为提督守护，9为总督守护
                "type": 0,
                // 数量
                "num": 0,
                // 奖品真实id
                "object_id": 0,
                // web端图片
                "web_url": "",
                // 移动端图片
                "mobile_url": "",
                // 奖励描述
                "description": "",
                // 跳转地址
                "jump_url": "",
                // 概率类型 1为普通，2为固定每天，3为固定每周
                "pro_type": 0,
                // 概率，3位小数，''为另一种概率模式
                "chance": "",
                // 循环的数量 0为另一种概率模式
                "loop": 0,
                // 限制数量 0为另一种概率模式
                "limit": 0,
                //  奖励名称
                "name": "",
                //  权重
                "weight": 0,
                //  白名单用户
                "white_uids": "",
                //  过期类型
                "expire": 0
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/capsule/get_prize_type
### 无标题

#### 方法：GET

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "type": 0,
                "name": ""
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/capsule/get_prize_expire
### 无标题

#### 方法：GET

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "expire": 0,
                "name": ""
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/capsule/update_pool_prize
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer|奖励id|
|pool_id|否|integer|奖池id|
|type|是|integer|奖品类型 1为道具，2为头衔，3为经验原石，4为经验曜石，5为贤者之石，6为小号小电视，7为舰长守护，8为提督守护，9为总督守护|
|num|是|integer|数量|
|object_id|否|integer|奖品真实id|
|expire|否|integer|过期时间|
|web_url|是|string|web端图片|
|mobile_url|是|string|移动端图片|
|description|是|string|奖励描述|
|jump_url|否|string|跳转地址|
|pro_type|是|integer|概率类型 1为普通，2为固定每天，3为固定每周，4位白名单|
|chance|否|string|概率，3位小数，''为另一种概率模式|
|loop|否|integer|循环的数量 0为另一种概率模式|
|limit|否|integer|限制数量 0为另一种概率模式|
|weight|否|integer| 权重|
|white_uids|否|string| 白名单用户|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true,
        // 新增id
        "prize_id": 0
    }
}
```


## /xlive/internal/live-admin/v1/capsule/delete_pool_prize
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer|奖励id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // 状态
        "status": true
    }
}
```


## /xlive/internal/live-admin/v1/capsule/get_coupon_list
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "uid": 0,
                //  中奖时间
                "award_time": "",
                //  奖品名称
                "award_name": "",
                //  券码
                "award_code": "",
                //  0 重试 1 成功
                "status": 0,
                //  上次重试时间
                "retry_time": ""
            }
        ]
    }
}
```

