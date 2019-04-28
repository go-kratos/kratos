<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/payLive/add](#xliveinternallive-adminv1payLiveadd) 
- [/xlive/internal/live-admin/v1/payLive/update](#xliveinternallive-adminv1payLiveupdate) 
- [/xlive/internal/live-admin/v1/payLive/getList](#xliveinternallive-adminv1payLivegetList) 
- [/xlive/internal/live-admin/v1/payLive/close](#xliveinternallive-adminv1payLiveclose) 
- [/xlive/internal/live-admin/v1/payLive/open](#xliveinternallive-adminv1payLiveopen) 

## /xlive/internal/live-admin/v1/payLive/add
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string| 平台|
|room_id|是|integer| 商品名称|
|title|是|string| 商品名称|
|status|否|integer| 鉴权状态，1开，0关|
|start_time|是|string| 开始时间|
|end_time|是|string| 结束时间|
|live_end_time|否|string| 正片结束时间|
|live_pic|是|string| 正片保底图|
|ad_pic|是|string| 广告图|
|goods_link|是|string| 购买链接|
|goods_id|是|string| 门票id，逗号分隔|
|ip_limit|否|integer| ip限制，0不限制，1仅限大陆，2仅限港澳台，3大陆+港澳台|
|buy_goods_id|是|integer| 购买门票id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/payLive/update
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|live_id|是|integer| id|
|platform|否|string| 平台|
|room_id|是|integer| 商品名称|
|title|是|string| 商品名称|
|status|否|integer| 鉴权状态，1开，0关|
|start_time|是|string| 开始时间|
|end_time|是|string| 结束时间|
|live_end_time|是|string| 正片结束时间|
|live_pic|是|string| 正片保底图|
|ad_pic|是|string| 广告图|
|goods_link|是|string| 购买链接|
|goods_id|是|string| 门票id，逗号分隔|
|ip_limit|否|integer| ip限制，0不限制，1仅限大陆，2仅限港澳台，3大陆+港澳台|
|buy_goods_id|是|integer| 购买门票id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/payLive/getList
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|否|integer| 房间id|
|title|否|string| 商品名称|
|ip_limit|否|integer| ip限制|
|page_num|否|integer| 页号，0开始|
|page_size|是|integer| 每页个数|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "page_info": {
            //  记录总数
            "total_count": 0,
            //  当前页号
            "page_num": 0
        },
        "goods_info": [
            {
                //  房间id
                "room_id": 0,
                //  付费直播id
                "live_id": 0,
                //  标题
                "title": "",
                //  平台
                "platform": "",
                //  生效状态，1生效，0未生效
                "pay_live_status": 0,
                //  开始购票时间
                "start_time": "",
                //  结束购票时间
                "end_time": "",
                //  正片结束
                "live_end_time": "",
                //  正片保底图
                "live_pic": "",
                //  广告图
                "ad_pic": "",
                //  购票链接
                "goods_link": "",
                //  购票id
                "goods_id": "",
                //  ip限制
                "ip_limit": 0,
                //  鉴权状态，0关闭，1开启
                "status": 0,
                //  引导购票id
                "buy_goods_id": 0
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/payLive/close
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|live_id|是|integer| 直播id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/payLive/open
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|live_id|是|integer| 直播id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

