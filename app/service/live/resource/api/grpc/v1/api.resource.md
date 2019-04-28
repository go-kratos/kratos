##Add 添加资源接口

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/Add`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|title|是|string||
|jumpPath|否|string||
|jumpTime|否|integer||
|type|是|string||
|device|是|string||
|startTime|是|string||
|endTime|是|string||
|imageUrl|是|string||
|jumpPathType|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "id": [
            0
        ]
    }
}
```

##Add 添加资源接口(不限制位置和平台)

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/AddEx`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|title|是|string||
|jumpPath|否|string||
|jumpTime|否|integer||
|type|是|string||
|device|是|string||
|startTime|是|string||
|endTime|是|string||
|imageUrl|是|string||
|jumpPathType|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "id": [
            0
        ]
    }
}
```

##Edit 编辑资源接口

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/Edit`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|id|是|integer||
|title|否|string||
|jumpPath|否|string||
|jumpTime|否|integer||
|startTime|否|string||
|endTime|否|string||
|imageUrl|否|string||
|jumpPathType|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

##Offline 下线资源接口

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/Offline`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|id|是|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

##GetList 获取资源列表

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/GetList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|page|否|integer||
|pageSize|否|integer||
|type|是|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "currentPage": 0,
        "totalCount": 0,
        "list": [
            {
                "id": 0,
                "title": "",
                "jumpPath": "",
                "device_platform": "",
                "device_build": 0,
                "startTime": "",
                "endTime": "",
                "status": 0,
                "device_limit": 0,
                "imageUrl": "",
                "jumpPathType": 0,
                "jumpTime": 0
            }
        ]
    }
}
```

##获取平台列表

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/GetPlatformList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "platform": [
            ""
        ]
    }
}
```

##GetListEx 获取资源列表

`GET http://api.live.bilibili.com/xlive/resource/v1/resource/GetListEx`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|page|否|integer||
|pageSize|否|integer||
|type|是|多个string||
|device_platform|否|string||
|status|否|string||
|startTime|否|string||
|endTime|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "currentPage": 0,
        "totalCount": 0,
        "list": [
            {
                "id": 0,
                "title": "",
                "jumpPath": "",
                "device_platform": "",
                "device_build": 0,
                "startTime": "",
                "endTime": "",
                "status": 0,
                "device_limit": 0,
                "imageUrl": "",
                "jumpPathType": 0,
                "jumpTime": 0,
                "type": ""
            }
        ]
    }
}
```

