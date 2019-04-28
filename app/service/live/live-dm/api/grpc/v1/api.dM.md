##

`GET http://api.live.bilibili.com/xlive/live-dm/v1/dM/SendMsg`

### 请求参数

```json
{
    "uid": 0,
    "roomid": 0,
    "msg": "",
    "rnd": "",
    "ip": "",
    "fontsize": 0,
    "mode": 0,
    "platform": "",
    "msgtype": 0,
    "bubble": 0,
    "lancer": {
        "buvid": "",
        "userAgent": "",
        "refer": "",
        "cookie": "",
        "build": 0
    }
}
```

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "islimit": true,
        "limitmsg": "",
        "code": 0
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/live-dm/v1/dM/GetHistory`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "room": [
            ""
        ],
        "admin": [
            ""
        ]
    }
}
```

