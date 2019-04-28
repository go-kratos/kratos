##获取有效banner配置

`GET http://api.live.bilibili.com/xlive/resource/v1/banner/GetBlinkBanner`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string||
|build|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "id": 0,
        "title": "",
        "jumpPath": "",
        "jumpTime": 0,
        "jumpPathType": 0,
        "imageUrl": ""
    }
}
```

##获取有效banner配置

`GET http://api.live.bilibili.com/xlive/resource/v1/banner/GetBanner`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string||
|build|否|integer||
|type|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "id": 0,
                "title": "",
                "jumpPath": "",
                "jumpTime": 0,
                "jumpPathType": 0,
                "imageUrl": ""
            }
        ]
    }
}
```

