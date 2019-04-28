##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/get_config_by_keyword`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|否|integer|team|
|keyword|否|string|索引名称|
|id|否|integer|id|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        // team
        "team": 0,
        // 索引名称
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
        "status": 0,
        // 状态
        "id": 0
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/set_config_by_keyword`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|否|integer|team|
|keyword|是|string|索引名称|
|value|是|string|配置值|
|name|否|string|配置解释|
|id|否|integer|编辑时id|
|status|否|integer|记录状态 新增时默认为0|

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

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/get_configs_by_params`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|否|integer||
|keyword|否|string||
|name|否|string||
|status|否|integer||
|page|是|integer||
|page_size|是|integer|页量|
|id|否|integer|id|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        // 数据列表
        "list": [
            {
                // Id
                "id": 0,
                // team
                "team": 0,
                // 索引名称
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
        // 记录总数
        "total_num": 0
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/getByTreeId`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|tree_id|是|integer||

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

`GET http://api.live.bilibili.com/xlive/internal/resource/v1/titans/get_configs_by_likes`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|params|是|多个string||

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

