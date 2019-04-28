<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/resource/add](#xliveinternallive-adminv1resourceadd) Add 添加资源接口
- [/xlive/internal/live-admin/v1/resource/addEx](#xliveinternallive-adminv1resourceaddEx) AddEx 添加资源接口(不限制位置和平台)
- [/xlive/internal/live-admin/v1/resource/edit](#xliveinternallive-adminv1resourceedit) Edit 编辑资源接口
- [/xlive/internal/live-admin/v1/resource/offline](#xliveinternallive-adminv1resourceoffline) Offline 下线资源接口
- [/xlive/internal/live-admin/v1/resource/getList](#xliveinternallive-adminv1resourcegetList) GetList 获取资源列表
- [/xlive/internal/live-admin/v1/resource/getPlatformList](#xliveinternallive-adminv1resourcegetPlatformList) 获取平台列表
- [/xlive/internal/live-admin/v1/resource/getListEx](#xliveinternallive-adminv1resourcegetListEx) GetListEx 获取资源列表

## /xlive/internal/live-admin/v1/resource/add
###Add 添加资源接口

#### 方法：POST

#### 请求参数

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

#### 响应

```javascript
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


## /xlive/internal/live-admin/v1/resource/addEx
###AddEx 添加资源接口(不限制位置和平台)

#### 方法：POST

#### 请求参数

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

#### 响应

```javascript
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


## /xlive/internal/live-admin/v1/resource/edit
###Edit 编辑资源接口

#### 方法：POST

#### 请求参数

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

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/resource/offline
###Offline 下线资源接口

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v1/resource/getList
###GetList 获取资源列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|page|否|integer||
|pageSize|否|integer||
|type|是|string||

#### 响应

```javascript
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


## /xlive/internal/live-admin/v1/resource/getPlatformList
###获取平台列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|type|否|integer||

#### 响应

```javascript
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


## /xlive/internal/live-admin/v1/resource/getListEx
###GetListEx 获取资源列表

#### 方法：GET

#### 请求参数

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

#### 响应

```javascript
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

