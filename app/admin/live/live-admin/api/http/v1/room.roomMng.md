<!-- package=live.liveadmin.v1 -->
- [/xlive/internal/live-admin/v1/roomMng/getSecondVerifyListWithPics](#xliveinternallive-adminv1roomMnggetSecondVerifyListWithPics)  获取带有图片地址的二次审核列表

## /xlive/internal/live-admin/v1/roomMng/getSecondVerifyListWithPics
### 获取带有图片地址的二次审核列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|否|integer| 房间id|
|area|否|string| 分区id多个|
|page|是|integer| 页数|
|pagesize|否|integer| 页码|
|biz|否|string| 业务，0直播监控1直播鉴黄2房间举报|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "result": [
            {
                //  日志id
                "id": 0,
                //  当天切断记录
                "recent_cut_times": 0,
                //  当天警告记录
                "recent_warn_times": 0,
                //  用户名
                "uname": "",
                //  房间号
                "room_id": 0,
                //  主播id
                "uid": 0,
                //  房间标题
                "title": "",
                //  分区名
                "area_v2_name": "",
                //  粉丝数
                "fc": 0,
                //  警告理由
                "warn_reason": "",
                //  图片列表
                "pics": [
                    ""
                ],
                //  警告时间
                "break_time": "",
                //  共计警告时间
                "warn_times": 0
            }
        ],
        //  总数
        "count": 0,
        //  页码
        "page": 0,
        //  分页大小
        "pagesize": 0
    }
}
```

