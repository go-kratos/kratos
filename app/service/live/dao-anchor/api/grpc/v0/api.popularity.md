<!-- package=live.daoanchor.v0 -->
- [/live.daoanchor.v0.Popularity/GetAnchorGradeList](#live.daoanchor.v0.PopularityGetAnchorGradeList)  GetAnchorGradeList 获取人气值主播评级列表
- [/live.daoanchor.v0.Popularity/EditAnchorGrade](#live.daoanchor.v0.PopularityEditAnchorGrade)  EditAnchorGrade  编辑主播评级对应的人气值数据
- [/live.daoanchor.v0.Popularity/GetContentList](#live.daoanchor.v0.PopularityGetContentList)  GetContentList  人气内容系数列表
- [/live.daoanchor.v0.Popularity/AddContent](#live.daoanchor.v0.PopularityAddContent)  AddContent 添加内容系数
- [/live.daoanchor.v0.Popularity/EditContent](#live.daoanchor.v0.PopularityEditContent)  EditContent 编辑内容系数
- [/live.daoanchor.v0.Popularity/DeleteContent](#live.daoanchor.v0.PopularityDeleteContent)  DeleteContent 删除内容系数

## /live.daoanchor.v0.Popularity/GetAnchorGradeList
### GetAnchorGradeList 获取人气值主播评级列表

#### 方法：GET

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                // 主播评级id 1=>S;2=>A;3=>B...
                "grade_id": 0,
                // 主播评级名称 S;A;B...
                "grade_name": "",
                // 人数基数
                "online_base": 0,
                // 人气倍数
                "popularity_rate": 0
            }
        ]
    }
}
```


## /live.daoanchor.v0.Popularity/EditAnchorGrade
### EditAnchorGrade  编辑主播评级对应的人气值数据

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|grade_id|是|integer||
|online_base|是|integer||
|popularity_rate|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.daoanchor.v0.Popularity/GetContentList
### GetContentList  人气内容系数列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|page|否|integer||
|page_size|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "page": 0,
        "page_size": 0,
        "total_count": 0,
        "list": [
            {
                // 父分区id
                "parent_area_id": 0,
                // 父分区名称
                "parent_area_name": "",
                // 二级分区信息<id,name>
                "area_list": {
                    "1": ""
                },
                // 人气倍率系数
                "popularity_rate": 0
            }
        ]
    }
}
```


## /live.daoanchor.v0.Popularity/AddContent
### AddContent 添加内容系数

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|parent_id|是|integer||
|list|是|多个integer||
|popularity_rate|是|integer||
|is_all|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.daoanchor.v0.Popularity/EditContent
### EditContent 编辑内容系数

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||
|popularity_rate|否|integer||
|list|否|多个integer||
|parent_id|否|integer||
|is_all|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.daoanchor.v0.Popularity/DeleteContent
### DeleteContent 删除内容系数

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

