##Add 添加资源接口

`GET http://api.live.bilibili.com/xlive/resource/v2/userResource/Add`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|title|是|string|名称|
|url|是|string|URL|
|weight|是|integer|权重|
|creator|是|string|创建人|

```json
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

##Edit 编辑现有资源

`GET http://api.live.bilibili.com/xlive/resource/v2/userResource/Edit`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|custom_id|是|integer|资源ID|
|title|否|string|名称|
|url|否|string|URL|
|weight|否|integer|权重|

```json
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

##Query 请求单个资源

`GET http://api.live.bilibili.com/xlive/resource/v2/userResource/Query`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|custom_id|是|integer|资源ID|

```json
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

##List 获取资源列表

`GET http://api.live.bilibili.com/xlive/resource/v2/userResource/List`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|page|否|integer|页码|
|page_size|否|integer|每页数据量|

```json
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

##SetStatus 更改资源状态

`GET http://api.live.bilibili.com/xlive/resource/v2/userResource/SetStatus`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|res_type|是|integer|资源类型|
|custom_id|是|integer|页码|
|status|是|integer|每页数据量|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

