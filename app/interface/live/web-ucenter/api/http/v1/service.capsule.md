<!-- package=live.webucenter.v1 -->
- [/xlive/web-ucenter/v1/capsule/get_detail](#xliveweb-ucenterv1capsuleget_detail) 
- [/xlive/web-ucenter/v1/capsule/open_capsule](#xliveweb-ucenterv1capsuleopen_capsule) 
- [/xlive/web-ucenter/v1/capsule/get_capsule_info](#xliveweb-ucenterv1capsuleget_capsule_info) 
- [/xlive/web-ucenter/v1/capsule/open_capsule_by_type](#xliveweb-ucenterv1capsuleopen_capsule_by_type) 

## /xlive/web-ucenter/v1/capsule/get_detail
### 无标题

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
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


## /xlive/web-ucenter/v1/capsule/open_capsule
### 无标题

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|是|string| 扭蛋类型|
|count|是|integer| 扭的个数|
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


## /xlive/web-ucenter/v1/capsule/get_capsule_info
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|是|string| 扭蛋类型|
|from|是|string| 来源 h5 web room|

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


## /xlive/web-ucenter/v1/capsule/open_capsule_by_type
### 无标题

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|是|string| 扭蛋类型|
|count|是|integer| 扭的个数|
|platform|否|string||

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

