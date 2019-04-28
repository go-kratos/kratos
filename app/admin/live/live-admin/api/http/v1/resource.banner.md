<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/banner/getBlinkBanner](#xliveinternallive-adminv1bannergetBlinkBanner) 获取有效banner配置
- [/xlive/internal/live-admin/v1/banner/getBanner](#xliveinternallive-adminv1bannergetBanner) 获取有效banner配置

## /xlive/internal/live-admin/v1/banner/getBlinkBanner
###获取有效banner配置

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


## /xlive/internal/live-admin/v1/banner/getBanner
###获取有效banner配置

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|build|是|integer||
|type|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "id": 0,
                "title": "",
                "jumpPath": "",
                "jumpTime": 0,
                "jumpPathType": 0,
                "imageUrl": ""
            }
        ]
    }
}
```

