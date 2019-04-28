##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/getMultiConfigs`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|values|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": {
            "1": {
                "keys": {
                    "mapKey": ""
                }
            }
        }
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/getServiceConfig`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|tree_id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": {
            "mapKey": ""
        }
    }
}
```

##

`POST http://api.live.bilibili.com/xlive/internal/resource/v1/titans/setServiceConfig`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|tree_name|否|string||
|tree_path|否|string||
|tree_id|否|integer||
|service|否|string||
|keyword|否|string||
|template|否|integer||
|name|否|string||
|value|否|string||
|status|否|integer||
|id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "id": 0
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/getServiceConfigList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|tree_name|否|string||
|tree_id|否|integer||
|service|否|string||
|keyword|否|string||
|page|否|integer||
|page_size|否|integer||
|name|否|string||
|status|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                // Id
                "id": 0,
                // tree_name
                "tree_name": "",
                "tree_path": "",
                "tree_id": 0,
                "service": "",
                // 索引名称
                "template": 0,
                "keyword": "",
                // 配置值
                "value": "",
                // 配置解释
                "name": "",
                // 创建时间
                "ctime": "",
                // 最近更新时间
                "mtime": "",
                // 状态
                "status": 0
            }
        ],
        "total_num": 0
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/getTreeIds`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|tree_name|否|string||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            0
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/getEasyList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||
|page|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "tree_name": "",
                "tree_path": "",
                "tree_id": 0,
                "keyword": "",
                "name": ""
            }
        ]
    }
}
```

##

`POST http://api.live.bilibili.com/xlive/internal/resource/v1/titans/setEasyList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "eId": 0
    }
}
```

