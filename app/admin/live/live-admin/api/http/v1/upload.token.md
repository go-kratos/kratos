<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/token/new](#xliveinternallive-adminv1tokennew)  Request for a token for upload.

## /xlive/internal/live-admin/v1/token/new
### Request for a token for upload.

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|bucket|是|string| 上传到 BFS 的 bucket|
|dir|否|string| 上传到指定的 BFS 目录（可以用来区分业务）|
|operator|是|string| 操作人（mlive通过dashboard授权获取到的操作人）|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  授予的 token
        "token": ""
    }
}
```

