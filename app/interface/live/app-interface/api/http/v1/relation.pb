syntax = "proto3";

package live.appinterface.v1;

option go_package = "v1";
import "github.com/gogo/protobuf/gogoproto/gogo.proto";

// Index 相关服务
service Relation {
    // [app端关注二级页][全量]正在直播接口
    // `midware:"guest"`
    rpc liveAnchor (LiveAnchorReq) returns (LiveAnchorResp);
    // [app端关注二级页][分页]暂未开播接口
    // `midware:"guest"`
    rpc unliveAnchor (UnLiveAnchorReq) returns (UnLiveAnchorResp);
}

// liveAnchor请求
message LiveAnchorReq {
    // 调试咒语
    string buyaofangqizhiliao = 1;
    // 平台
    string platform = 2;
    // 设备
    string device = 3;
    // 版本号
    string build = 4;
    // 排序类型
    int64 sortRule = 5;
    // 筛选类型
    int64 filterRule = 6;
    // 清晰度
    int64 quality = 7;

}

// liveAnchor响应
message LiveAnchorResp {
    repeated Rooms rooms = 1;
    message Rooms {
        // 房间id
        int64 roomid = 1;
        // 用户id
        int64 uid = 2;
        // 用户昵称
        string uname = 3;
        // 用户头像
        string face = 4;
        // 直播间标题
        string title = 5;
        // 直播间标签
        string live_tag_name = 6;
        // 开始直播时间
        int64 live_time = 7;
        // 人气值
        int64 online = 8;
        // 秒开url
        string playurl = 9;
        // 可选清晰度
        repeated int64  accept_quality = 10;
        // 当前清晰度
        int64 current_quality = 11;
        // pk_id
        int64 pk_id = 12;
        // 特别关注标志
        int64 special_attention = 13;
        // 老的分区id
        int64 area = 14;
        // 老的分区名
        string area_name = 15;
        // 子分区id
        int64 area_v2_id = 16;
        // 子分区名
        string area_v2_name = 17;
        // 父分区名
        string area_v2_parent_name = 18;
        // 父分区id
        int64 area_v2_parent_id = 19;
        // 广播适配标志
        int64 broadcast_type = 20;
        // 官方认证标志
        int64 official_verify = 21;
        // 直播间跳转链接
        string link = 22;
        // 直播间封面
        string cover = 23;
        // 角标文字
        string pendent_ru = 24;
        // 角标颜色
        string pendent_ru_color = 25;
        // 角标背景图
        string pendent_ru_pic = 26;
        string play_url_h265 = 27;
    }
    int64 total_count = 2;
    int64 card_type = 3;
    int64 big_card_type = 4;
}


// unLiveAnchor请求
message UnLiveAnchorReq {
    // 调试咒语
    string buyaofangqizhiliao = 1;
    // 分页号
    int64 page = 2;
    // 页大小
    int64 pagesize = 3;
}

// unLiveAnchor响应
message UnLiveAnchorResp {
    repeated Rooms rooms = 1;
    message Rooms {
        // 上次开播描述
        string live_desc = 1;
        // 房间id
        int64 roomid = 2;
        // 用户id
        int64 uid = 3;
        // 用户昵称
        string uname = 4;
        // 用户头像
        string face = 5;
        // 特别关注标志
        int64 special_attention = 6;
        // 官方认证标志
        int64 official_verify = 7;
        // 直播状态标志
        int64 live_status = 8;
        // 广播适配标志
        int64 broadcast_type = 9;
        // 老的分区id
        int64 area = 10;
        // 粉丝数
        int64 attentions = 11;
        // 老的分区名
        string area_name = 12;
        // 子分区id
        int64 area_v2_id = 13;
        // 子分区名
        string area_v2_name = 14;
        // 父分区名
        string area_v2_parent_name = 15;
        // 父分区id
        int64 area_v2_parent_id = 16;
        // 直播间跳转链接
        string link = 17;
        // 房间页公告
        string announcement_content = 18;
        // 房间页公告发布时间
        string announcement_time = 19;
    }
    int64 total_count = 2;
    int64 no_room_count = 3;
    int64 has_more = 4;
}
