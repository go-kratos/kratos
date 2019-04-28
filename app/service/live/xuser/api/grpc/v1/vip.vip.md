## Info 返回用户vip信息

`GET http://api.live.bilibili.com/xlive/xuser/v1/vip/Info`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "info": {
            "vip": 0,
            "svip": 0,
            "vip_time": "",
            "svip_time": ""
        }
    }
}
```

## Buy 购买月费/年费姥爷

`GET http://api.live.bilibili.com/xlive/xuser/v1/vip/Buy`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|order_id|是|string||
|uid|是|integer||
|good_id|是|integer||
|good_num|是|integer||
|platform|是|integer||
|source|是|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "status": 0
    }
}
```

