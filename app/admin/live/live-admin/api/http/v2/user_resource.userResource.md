<!-- package=live.liveadmin.v2 -->
- [/xlive/internal/live-admin/v2/userResource/add](#xliveinternallive-adminv2userResourceadd) Add 添加资源接口
- [/xlive/internal/live-admin/v2/userResource/edit](#xliveinternallive-adminv2userResourceedit) Edit 编辑现有资源
- [/xlive/internal/live-admin/v2/userResource/get](#xliveinternallive-adminv2userResourceget) List 获取资源列表
- [/xlive/internal/live-admin/v2/userResource/setStatus](#xliveinternallive-adminv2userResourcesetStatus) SetStatus 更改资源状态
- [/xlive/internal/live-admin/v2/userResource/getSingle](#xliveinternallive-adminv2userResourcegetSingle) Query 请求单个资源    

## /xlive/internal/live-admin/v2/userResource/add
###Add 添加资源接口

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|title|是|string|名称|
|url|是|string|URL|
|weight|是|integer|权重|
|creator|是|string|创建人|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // ID
        "id": 0,
        // 资源ID
        "custom_id": 0
    }
}
```


## /xlive/internal/live-admin/v2/userResource/edit
###Edit 编辑现有资源

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|custom_id|是|integer|资源ID|
|title|否|string|名称|
|url|否|string|URL|
|weight|否|integer|权重|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v2/userResource/get
###List 获取资源列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|page|否|integer|页码|
|page_size|否|integer|每页数据量|

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
                // ID
                "id": 0,
                // 资源类型
                "res_type": 0,
                // 资源ID
                "custom_id": 0,
                // 名称
                "title": "",
                // URL
                "url": "",
                // 权重
                "weight": 0,
                // 创建人
                "creator": "",
                // "状态1.上线中2.下线"
                "status": 0,
                // 创建时刻
                "ctime": "",
                // 修改时刻
                "mtime": ""
            }
        ]
    }
}
```


## /xlive/internal/live-admin/v2/userResource/setStatus
###SetStatus 更改资源状态

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|custom_id|是|integer|页码|
|status|是|integer|每页数据量|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/internal/live-admin/v2/userResource/getSingle
###Query 请求单个资源    

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|custom_id|是|integer|资源ID|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        // ID
        "id": 0,
        // 资源类型
        "res_type": 0,
        // 资源ID
        "custom_id": 0,
        // 名称
        "title": "",
        // URL
        "url": "",
        // 权重
        "weight": 0,
        // 创建人
        "creator": "",
        // "状态1.上线中2.下线"
        "status": 0,
        // 创建时刻
        "ctime": "",
        // 修改时刻
        "mtime": ""
    }
}
```

