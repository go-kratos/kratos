## 获取主播拥有的的所有房管, 无需登录态
 `method:"GET"

`GET http://api.live.bilibili.com/xlive/web-room/v1/roomAdmin/get_by_room`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer| 页数|
|roomid|是|integer| 房间号|
|page_size|否|integer| 每页数量|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "page": {
            //  当前页码
            "page": 0,
            //  每页大小
            "page_size": 0,
            //  总页数
            "total_page": 0,
            //  总记录数
            "total_count": 0
        },
        "data": [
            {
                //  用户id
                "uid": 0,
                //  用户名
                "uname": "",
                //  用户头像
                "face": "",
                //  上任时间
                "ctime": ""
            }
        ]
    }
}
```

