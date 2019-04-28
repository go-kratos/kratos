##

`POST http://api.live.bilibili.com/xlive/web-room/v1/dM/SendMsg`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer||
|msg|是|string||
|rnd|是|string||
|fontsize|是|integer||
|mode|否|integer||
|color|是|integer||
|bubble|否|integer||
|anti|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

##

`POST http://api.live.bilibili.com/xlive/web-room/v1/dM/GetHistory`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "Room": [
            ""
        ],
        "Admin": [
            ""
        ]
    }
}
```

