##
> 需要登录

`GET http://api.live.bilibili.com/xlive/web-room/v1/captcha/create`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|否|integer||
|client_type|否|string||
|height|否|integer||
|width|否|integer||

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
> 需要登录

`POST http://api.live.bilibili.com/xlive/web-room/v1/captcha/verify`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|anti|否|string||

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

