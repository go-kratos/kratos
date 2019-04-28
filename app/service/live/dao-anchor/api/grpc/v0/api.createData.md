<!-- package=live.daoanchor.v0 -->
- [/live.daoanchor.v0.CreateData/CreateCacheList](#live.daoanchor.v0.CreateDataCreateCacheList)  CreateCacheList 生成历史数据缓存列表
- [/live.daoanchor.v0.CreateData/CreateLiveCacheList](#live.daoanchor.v0.CreateDataCreateLiveCacheList)  CreateLiveCacheList 生成开播历史数据缓存列表
- [/live.daoanchor.v0.CreateData/GetContentMap](#live.daoanchor.v0.CreateDataGetContentMap)  GetContentMap 获取需要生成历史数据的对象列表
- [/live.daoanchor.v0.CreateData/CreateDBData](#live.daoanchor.v0.CreateDataCreateDBData) 

## /live.daoanchor.v0.CreateData/CreateCacheList
### CreateCacheList 生成历史数据缓存列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|否|多个integer||
|content|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.daoanchor.v0.CreateData/CreateLiveCacheList
### CreateLiveCacheList 生成开播历史数据缓存列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|否|多个integer||
|content|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.daoanchor.v0.CreateData/GetContentMap
### GetContentMap 获取需要生成历史数据的对象列表

#### 方法：GET

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": {
            "mapKey": 0
        }
    }
}
```


## /live.daoanchor.v0.CreateData/CreateDBData
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|否|多个integer||
|content|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

