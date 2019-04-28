<!-- package=live.xlottery.v1 -->
- [/live.xlottery.v1.Capsule/get_detail](#live.xlottery.v1.Capsuleget_detail) 
- [/live.xlottery.v1.Capsule/open_capsule](#live.xlottery.v1.Capsuleopen_capsule) 
- [/live.xlottery.v1.Capsule/get_coin_list](#live.xlottery.v1.Capsuleget_coin_list) 
- [/live.xlottery.v1.Capsule/update_coin_config](#live.xlottery.v1.Capsuleupdate_coin_config) 
- [/live.xlottery.v1.Capsule/update_coin_status](#live.xlottery.v1.Capsuleupdate_coin_status) 
- [/live.xlottery.v1.Capsule/delete_coin](#live.xlottery.v1.Capsuledelete_coin) 
- [/live.xlottery.v1.Capsule/get_pool_list](#live.xlottery.v1.Capsuleget_pool_list) 
- [/live.xlottery.v1.Capsule/update_pool](#live.xlottery.v1.Capsuleupdate_pool) 
- [/live.xlottery.v1.Capsule/delete_pool](#live.xlottery.v1.Capsuledelete_pool) 
- [/live.xlottery.v1.Capsule/update_pool_status](#live.xlottery.v1.Capsuleupdate_pool_status) 
- [/live.xlottery.v1.Capsule/get_pool_prize](#live.xlottery.v1.Capsuleget_pool_prize) 
- [/live.xlottery.v1.Capsule/get_prize_type](#live.xlottery.v1.Capsuleget_prize_type) 
- [/live.xlottery.v1.Capsule/get_prize_expire](#live.xlottery.v1.Capsuleget_prize_expire) 
- [/live.xlottery.v1.Capsule/update_pool_prize](#live.xlottery.v1.Capsuleupdate_pool_prize) 
- [/live.xlottery.v1.Capsule/delete_pool_prize](#live.xlottery.v1.Capsuledelete_pool_prize) 
- [/live.xlottery.v1.Capsule/get_capsule_info](#live.xlottery.v1.Capsuleget_capsule_info) 
- [/live.xlottery.v1.Capsule/open_capsule_by_type](#live.xlottery.v1.Capsuleopen_capsule_by_type) 
- [/live.xlottery.v1.Capsule/get_coupon_list](#live.xlottery.v1.Capsuleget_coupon_list) 

## /live.xlottery.v1.Capsule/get_detail
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|
|from|否|string| 来源 h5 web room|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  普通扭蛋信息
        "normal": {
            "status": true,
            //  扭蛋数量
            "coin": 0,
            //  变化值
            "change": 0,
            //  进度
            "progress": {
                //  当前进度
                "now": 0,
                //  最大进度
                "max": 0
            },
            //  规则
            "rule": "",
            //  奖品列表
            "gift": [
                {
                    //  礼物名称
                    "name": "",
                    //  礼物图片
                    "image": "",
                    //  用法
                    "usage": {
                        //  用法描述
                        "text": "",
                        //  跳转链接
                        "url": ""
                    },
                    //  web礼物图片
                    "web_image": "",
                    //  mobile礼物图片
                    "mobile_image": ""
                }
            ],
            //  历史获奖列表
            "list": [
                {
                    //  数量
                    "num": 0,
                    //  礼物名称
                    "gift": "",
                    //  时间
                    "date": "",
                    //  用户名
                    "name": ""
                }
            ]
        },
        //  梦幻扭蛋信息，若梦幻扭蛋status=false，则无coin、change、process、gift、list字段
        "colorful": {
            "status": true,
            //  扭蛋数量
            "coin": 0,
            //  变化值
            "change": 0,
            //  进度
            "progress": {
                //  当前进度
                "now": 0,
                //  最大进度
                "max": 0
            },
            //  规则
            "rule": "",
            //  奖品列表
            "gift": [
                {
                    //  礼物名称
                    "name": "",
                    //  礼物图片
                    "image": "",
                    //  用法
                    "usage": {
                        //  用法描述
                        "text": "",
                        //  跳转链接
                        "url": ""
                    },
                    //  web礼物图片
                    "web_image": "",
                    //  mobile礼物图片
                    "mobile_image": ""
                }
            ],
            //  历史获奖列表
            "list": [
                {
                    //  数量
                    "num": 0,
                    //  礼物名称
                    "gift": "",
                    //  时间
                    "date": "",
                    //  用户名
                    "name": ""
                }
            ]
        }
    }
}
```


## /live.xlottery.v1.Capsule/open_capsule
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|
|type|否|string| 扭蛋类型|
|count|否|integer| 扭的个数|
|platform|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  扭蛋币扣除状态
        "status": true,
        //  奖品文案
        "text": [
            ""
        ],
        //  是否包含实物奖品
        "isEntity": true,
        //  用户扭蛋币拥有状态
        "info": {
            //  普通扭蛋币
            "normal": {
                //  拥有的币
                "coin": 0,
                //  变化值
                "change": 0,
                //  进度
                "progress": {
                    //  当前进度
                    "now": 0,
                    //  最大进度
                    "max": 0
                }
            },
            //  梦幻扭蛋币
            "colorful": {
                //  拥有的币
                "coin": 0,
                //  变化值
                "change": 0,
                //  进度
                "progress": {
                    //  当前进度
                    "now": 0,
                    //  最大进度
                    "max": 0
                }
            }
        },
        //  头衔? 恒为空字符串, 忽略之
        "showTitle": "",
        //  奖品列表
        "awards": [
            {
                //  奖品名字
                "name": "",
                //  奖品数量
                "num": 0,
                //  奖品 X 数量
                "text": "",
                //  奖品图片
                "img": "",
                //  奖品用法说明
                "usage": {
                    //  用法描述
                    "text": "",
                    //  跳转链接
                    "url": ""
                },
                //  web礼物图片
                "web_image": "",
                //  mobile礼物图片
                "mobile_image": ""
            }
        ]
    }
}
```


## /live.xlottery.v1.Capsule/get_coin_list
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|是|integer||
|page_size|是|integer||

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
                "id": 0,
                "title": "",
                "change_num": 0,
                "start_time": 0,
                "end_time": 0,
                "status": 0,
                "gift_type": 0,
                "gift_config": "",
                "area_ids": [
                    {
                        "parent_id": 0,
                        "is_all": 0,
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


## /live.xlottery.v1.Capsule/update_coin_config
### 无标题

#### 方法：GET

#### 请求参数

```javascript
{
    "id": 0,
    "title": "",
    "change_num": 0,
    "start_time": 0,
    "end_time": 0,
    "status": 0,
    "gift_type": 0,
    "gift_ids": [
        0
    ],
    "area_ids": [
        {
            "parent_id": 0,
            "is_all": 0,
            "list": [
                0
            ]
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/update_coin_status
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||
|status|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/delete_coin
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/get_pool_list
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|是|integer||
|page_size|是|integer||

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
                "id": 0,
                "coin_id": 0,
                "title": "",
                "coin_title": "",
                "start_time": 0,
                "end_time": 0,
                "status": 0,
                "rule": ""
            }
        ]
    }
}
```


## /live.xlottery.v1.Capsule/update_pool
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|coin_id|是|integer||
|title|是|string||
|start_time|是|integer||
|end_time|是|integer||
|rule|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/delete_pool
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/update_pool_status
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||
|status|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/get_pool_prize
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|pool_id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "id": 0,
                "pool_id": 0,
                "type": 0,
                "num": 0,
                "object_id": 0,
                "web_url": "",
                "mobile_url": "",
                "description": "",
                "jump_url": "",
                "pro_type": 0,
                "chance": "",
                "loop": 0,
                "limit": 0,
                "name": "",
                "weight": 0,
                "white_uids": [
                    0
                ],
                "expire": 0
            }
        ]
    }
}
```


## /live.xlottery.v1.Capsule/get_prize_type
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


## /live.xlottery.v1.Capsule/get_prize_expire
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


## /live.xlottery.v1.Capsule/update_pool_prize
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|pool_id|否|integer||
|type|是|integer||
|num|是|integer||
|object_id|否|integer||
|expire|否|integer||
|web_url|是|string||
|mobile_url|是|string||
|description|是|string||
|jump_url|否|string||
|pro_type|是|integer||
|chance|否|integer||
|loop|否|integer||
|limit|否|integer||
|weight|否|integer||
|white_uids|否|多个integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true,
        "prize_id": 0
    }
}
```


## /live.xlottery.v1.Capsule/delete_pool_prize
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": true
    }
}
```


## /live.xlottery.v1.Capsule/get_capsule_info
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|
|type|否|string| 类型|
|from|否|string| 来源 h5 web room|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  扭蛋数量
        "coin": 0,
        //  规则
        "rule": "",
        //  奖品列表，包含数量
        "gift_list": [
            {
                //  礼物id
                "id": 0,
                //  礼物名称
                "name": "",
                //  礼物数量
                "num": 0,
                //  权重
                "weight": 0,
                //  礼物图片
                "web_url": "",
                //  礼物图片
                "mobile_url": "",
                //  用法
                "usage": {
                    //  用法描述
                    "text": "",
                    //  跳转链接
                    "url": ""
                },
                //  奖品类型 2 头衔
                "type": 0,
                //  过期时间
                "expire": ""
            }
        ],
        //  奖品列表，不包含数量，同一类别只有一条
        "gift_filter": [
            {
                //  礼物id
                "id": 0,
                //  礼物名称
                "name": "",
                //  礼物图片
                "web_url": "",
                //  礼物图片
                "mobile_url": "",
                //  用法
                "usage": {
                    //  用法描述
                    "text": "",
                    //  跳转链接
                    "url": ""
                }
            }
        ]
    }
}
```


## /live.xlottery.v1.Capsule/open_capsule_by_type
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|
|type|否|string| 扭蛋类型，week：周星|
|count|否|integer| 扭的个数 1 10 100|
|platform|否|string| 平台|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  扭蛋币扣除状态
        "status": true,
        //  是否包含实物奖品
        "isEntity": true,
        //  用户扭蛋币拥有状态
        "info": {
            //  拥有的币
            "coin": 0
        },
        //  奖品列表
        "awards": [
            {
                //  奖品id
                "id": 0,
                //  奖品名字
                "name": "",
                //  奖品数量
                "num": 0,
                //  奖品 X 数量
                "text": "",
                //  礼物图片
                "web_url": "",
                //  礼物图片
                "mobile_url": "",
                //  奖品用法说明
                "usage": {
                    //  用法描述
                    "text": "",
                    //  跳转链接
                    "url": ""
                },
                //  奖品权重
                "weight": 0,
                //  奖品类型 2 头衔
                "type": 0,
                //  过期时间
                "expire": ""
            }
        ],
        //  奖品列表
        "text": [
            ""
        ]
    }
}
```


## /live.xlottery.v1.Capsule/get_coupon_list
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer||

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

