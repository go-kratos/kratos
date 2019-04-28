<!-- package=live.xroomfeed.v1 -->
- [/live.xroomfeed.v1.RecPool/GetList](#live.xroomfeed.v1.RecPoolGetList)  根据模块位置获取投放列表 position=>RoomItem

## /live.xroomfeed.v1.RecPool/GetList
### 根据模块位置获取投放列表 position=>RoomItem

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|module_type|是|integer| 投放模块|
|position_num|是|integer| 投放模块位置数|
|page_num|否|integer| 投放模块页数 不传或传0、1都按一页算(暂时没用)|
|module_exist_rooms|否|string| 当前模块已存在的位置房间（逗号分隔、有序），1~position*N（内部去重,保证同一个房间优先出现在好位置）|
|other_exist_rooms|否|string| 其它模块已存在的位置房间（逗号分隔、有序），1~position*N（内部去重,保证同一个房间优先出现在好位置）|
|from|否|string| 请求来源|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  主播position => 房间信息(依赖计算的)
        "list": {
            "1": {
                // 房间id
                "room_id": 0,
                // 主播uid
                "uid": 0,
                // 房间标题
                "title": "",
                // 人气
                "popularity_count": 0,
                // 关键帧
                "keyframe": "",
                // 封面
                "cover": "",
                // 二级分区id
                "area_id": 0,
                // 一级分区id
                "parent_area_id": 0,
                // 二级分区名称
                "area_name": "",
                // 一级分区名称
                "parent_area_name": "",
                // 推荐规则 10000+rule_id
                "rec_type": 0
            }
        }
    }
}
```

