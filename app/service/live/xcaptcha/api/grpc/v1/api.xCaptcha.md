## 创建验证码

`GET http://api.live.bilibili.com/xlive/xcaptcha/v1/xCaptcha/create`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|否|integer||
|client_type|否|string||
|height|否|integer||
|width|否|integer||
|uid|否|integer||
|client_ip|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "type": 0,
        "geetest": {
            "gt": "",
            "challenge": ""
        },
        "image": {
            "tips": "",
            "token": "",
            "content": ""
        }
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/xcaptcha/v1/xCaptcha/verify`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|_anti|是|string||
|uid|是|integer||
|client_ip|否|string||
|room_id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/xcaptcha/v1/xCaptcha/check`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|anti|是|string||
|uid|是|integer||
|client_ip|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "type": 0,
        "token": ""
    }
}
```

