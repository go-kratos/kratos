##客户端获取能否直播接口

`GET http://api.live.bilibili.com/xlive/resource/v1/liveCheck/LiveCheck`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string| 平台|
|system|否|string| 操作系统|
|mobile|否|string| 设备|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "is_live": 0
    }
}
```

##后台查询所有配置设备黑名单

`GET http://api.live.bilibili.com/xlive/resource/v1/liveCheck/GetLiveCheckList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  android
        "android": [
            {
                //  os
                "system": "",
                //  device
                "mobile": [
                    ""
                ]
            }
        ],
        //  ios
        "ios": [
            {
                //  os
                "system": "",
                //  device
                "mobile": [
                    ""
                ]
            }
        ]
    }
}
```

##后台添加能否直播设备黑名单

`GET http://api.live.bilibili.com/xlive/resource/v1/liveCheck/AddLiveCheck`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|live_check|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

