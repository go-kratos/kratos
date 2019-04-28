<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/gaea/get_config_by_keyword](#xliveinternallive-adminv1gaeaget_config_by_keyword) 
- [/xlive/internal/live-admin/v1/gaea/set_config_by_keyword](#xliveinternallive-adminv1gaeaset_config_by_keyword) 
- [/xlive/internal/live-admin/v1/gaea/get_configs_by_params](#xliveinternallive-adminv1gaeaget_configs_by_params) 
- [/xlive/internal/live-admin/v1/gaea/get_configs_by_team](#xliveinternallive-adminv1gaeaget_configs_by_team) 
- [/xlive/internal/live-admin/v1/gaea/get_configs_by_keyword](#xliveinternallive-adminv1gaeaget_configs_by_keyword) 
- [/xlive/internal/live-admin/v1/gaea/get_configs_by_teams](#xliveinternallive-adminv1gaeaget_configs_by_teams) 

## /xlive/internal/live-admin/v1/gaea/get_config_by_keyword
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|否|integer|team|
|keyword|否|string|索引名称|
|id|否|integer|id|

#### 响应

```javascript
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
        // id
        "id": 0
    }
}
```


## /xlive/internal/live-admin/v1/gaea/set_config_by_keyword
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|否|integer|team|
|keyword|是|string|索引名称|
|value|是|string|配置值|
|name|否|string|配置解释|
|id|否|integer|编辑时id|
|status|否|integer|状态|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "id": 0
    }
}
```


## /xlive/internal/live-admin/v1/gaea/get_configs_by_params
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|否|integer||
|keyword|否|string||
|name|否|string||
|status|否|integer||
|page|是|integer||
|page_size|是|integer|页量|
|id|否|integer|id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
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


## /xlive/internal/live-admin/v1/gaea/get_configs_by_team
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|是|string|team|
|page|是|integer|页码 从1开始|
|page_size|是|integer|页量|

#### 响应

```javascript
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


## /xlive/internal/live-admin/v1/gaea/get_configs_by_keyword
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|keyword|是|string|索引名称|
|page|是|integer|页码|
|page_size|是|integer|页量|

#### 响应

```javascript
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


## /xlive/internal/live-admin/v1/gaea/get_configs_by_teams
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|team|是|多个integer|team|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": ""
    }
}
```

