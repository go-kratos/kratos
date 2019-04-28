<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/payGoods/add](#xliveinternallive-adminv1payGoodsadd) 
- [/xlive/internal/live-admin/v1/payGoods/update](#xliveinternallive-adminv1payGoodsupdate) 
- [/xlive/internal/live-admin/v1/payGoods/getList](#xliveinternallive-adminv1payGoodsgetList) 
- [/xlive/internal/live-admin/v1/payGoods/close](#xliveinternallive-adminv1payGoodsclose) 
- [/xlive/internal/live-admin/v1/payGoods/open](#xliveinternallive-adminv1payGoodsopen) 

## /xlive/internal/live-admin/v1/payGoods/add
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string| 平台|
|title|是|string| 商品名称|
|type|是|integer| 商品类型 2 付费直播门票|
|price|是|integer| 商品价格(分)|
|start_time|是|string| 开始时间|
|end_time|是|string| 结束时间|
|ip_limit|否|integer| ip限制，0不限制，1仅限大陆，2仅限港澳台，3大陆+港澳台|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/payGoods/update
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer| 购票id|
|platform|否|string| 平台|
|title|否|string| 商品名称|
|type|否|integer| 商品类型 2 付费直播门票|
|price|否|integer| 商品价格(分)|
|start_time|否|string| 开始时间|
|end_time|否|string| 结束时间|
|ip_limit|否|integer| ip限制，0不限制，1仅限大陆，2仅限港澳台，3大陆+港澳台|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/payGoods/getList
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer| 购票id|
|platform|否|string| 平台|
|title|否|string| 商品名称|
|type|否|integer| 商品类型 2 付费直播门票|
|ip_limit|否|integer| ip限制，0不限制，1仅限大陆，2仅限港澳台，3大陆+港澳台|
|page_num|否|integer| 页号，0开始|
|page_size|否|integer| 每页个数|

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
                //  购票id
                "id": 0,
                //  标题
                "title": "",
                //  平台
                "platform": "",
                //  类型，2为付费直播
                "type": 0,
                //  价格，分
                "price": 0,
                //  开始购票时间
                "start_time": "",
                //  结束购票时间
                "end_time": "",
                //  ip限制
                "ip_limit": 0,
                //  购票状态，0关闭，1购票中，2未开始
                "status": 0
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v1/payGoods/close
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer| 购票id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/payGoods/open
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer| 购票id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

