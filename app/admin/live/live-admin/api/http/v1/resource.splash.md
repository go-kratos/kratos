<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/splash/getInfo](#xliveinternallive-adminv1splashgetInfo) 获取有效闪屏配置

## /xlive/internal/live-admin/v1/splash/getInfo
###获取有效闪屏配置

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|build|是|integer||

#### 响应

```javascript
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

