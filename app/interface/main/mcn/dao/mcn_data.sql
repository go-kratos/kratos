#投稿数及昨日增量
CREATE TABLE `dm_con_mcn_archive_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `up_all` int(11) NOT NULL DEFAULT '0' COMMENT '绑定up主数',
  `archive_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总投稿数',
  `archive_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '投稿数昨日增量',
  `play_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总播放数',
  `play_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '播放数昨日增量',
  `danmu_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总弹幕数',
  `danmu_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '弹幕数昨日增量',
  `reply_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总评论数',
  `reply_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '评论数昨日增量',
  `share_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总分享数',
  `share_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '分享数昨日增量',
  `coin_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总硬币数',
  `coin_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '硬币数昨日增量',
  `fav_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总收藏数',
  `fav_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '收藏数昨日增量',
  `like_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总点赞数',
  `like_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '点赞数昨日增量',
  `fans_all` bigint(20) NOT NULL DEFAULT '0' COMMENT '总粉丝数',
  `fans_inc` bigint(20) NOT NULL DEFAULT '0' COMMENT '昨日粉丝增量',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn稿件汇总指标';


#播放/弹幕/评论/分享/硬币/收藏/点赞数每日增量
CREATE TABLE `dm_con_mcn_index_inc_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `value` bigint(20) NOT NULL DEFAULT '0' COMMENT '当日播放/弹幕/评论/分享/硬币/收藏/点赞数',
  `type` varchar(20) NOT NULL COMMENT '分区类型，play、danmu、reply、share、coin、fav、like',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='播放/弹幕/评论/分享/硬币/收藏/点赞数每日增量';


#mcn播放/弹幕/评论/分享/硬币/收藏/点赞来源分区
CREATE TABLE `dm_con_mcn_index_source_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `type_id` int(11) NOT NULL COMMENT '一级分区ID',
  `rank` int(11) NOT NULL COMMENT '排名',
  `value` bigint(20) NOT NULL DEFAULT '0' COMMENT '一级分区',
  `type` varchar(20) NOT NULL COMMENT '分区类型，play、danmu、reply、share、coin、fav、like',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`type_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn播放/弹幕/评论/分享/硬币/收藏/点赞来源分区';


#mcn稿件播放来源占比
CREATE TABLE `dm_con_mcn_play_source_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `iphone` bigint(20) NOT NULL COMMENT 'iphone播放量',
  `andriod` bigint(20) NOT NULL COMMENT 'andriod播放量',
  `pc` bigint(20) NOT NULL COMMENT 'pc播放量',
  `h5` bigint(20) NOT NULL COMMENT 'h5播放量',
  `other` bigint(20) NOT NULL COMMENT 'other播放量',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn稿件播放来源占比';


#游客/粉丝性别占比
CREATE TABLE `dm_con_mcn_fans_sex_w` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `male` bigint(20) NOT NULL COMMENT '男性人数',
  `female` bigint(20) NOT NULL COMMENT '女性人数',
  `type` varchar(20) NOT NULL COMMENT '粉丝类型，guest、fans',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='游客/粉丝性别占比';


#游客/粉丝年龄分布
CREATE TABLE `dm_con_mcn_fans_age_w` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `a` bigint(20) NOT NULL COMMENT '0-16岁人数',
  `b` bigint(20) NOT NULL COMMENT '16-25岁人数',
  `c` bigint(20) NOT NULL COMMENT '25-40岁人数',
  `d` bigint(20) NOT NULL COMMENT '40岁以上人数',
  `type` varchar(20) NOT NULL COMMENT '粉丝类型，guest、fans',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='游客/粉丝年龄分布';


#游客/粉丝观看途径
CREATE TABLE `dm_con_mcn_fans_play_way_w` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `app` bigint(20) NOT NULL COMMENT 'app观看人数',
  `pc` bigint(20) NOT NULL COMMENT 'pc观看人数',
  `outside` bigint(20) NOT NULL COMMENT '站外观看人数',
  `other` bigint(20) NOT NULL COMMENT '其他观看人数',
  `type` varchar(20) NOT NULL COMMENT '粉丝类型，guest、fans',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='游客/粉丝观看途径';


#游客/粉丝地区分布
CREATE TABLE `dm_con_mcn_fans_area_w` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `province` varchar(200) NOT NULL COMMENT '省份',
  `user` bigint(20) NOT NULL COMMENT '人数',
  `type` varchar(20) NOT NULL COMMENT '粉丝类型，guest、fans',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`province`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='游客/粉丝地区分布';


#游客/粉丝倾向分布
CREATE TABLE `dm_con_mcn_fans_type_w` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `type_id` int(11) NOT NULL COMMENT '二级分区ID',
  `user` bigint(20) NOT NULL COMMENT '人数',
  `type` varchar(20) NOT NULL COMMENT '粉丝类型，guest、fans',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`type_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='游客/粉丝倾向分布';


#游客/粉丝标签地图分布
CREATE TABLE `dm_con_mcn_fans_tag_w` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `tag_id` int(11) NOT NULL COMMENT '标签ID',
  `user` bigint(20) NOT NULL COMMENT '人数',
  `type` varchar(20) NOT NULL COMMENT '粉丝类型，guest、fans',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`,`tag_id`,`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='游客/粉丝标签地图分布';


#mcn粉丝数相关
CREATE TABLE `dm_con_mcn_fans_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `fans_all` bigint(20) NOT NULL COMMENT 'mcn总粉丝数',
  `fans_inc` bigint(20) NOT NULL COMMENT 'mcn粉丝数昨日增量',
  `act_fans` bigint(20) NOT NULL COMMENT 'mcn活跃粉丝数',
  `fans_dec_all` bigint(20) NOT NULL COMMENT 'mcn取关粉丝总数',
  `fans_dec` bigint(20) NOT NULL COMMENT 'mcn昨日取关粉丝数',
  `view_fans_rate` float(3,2) NOT NULL COMMENT '观看活跃度',
  `act_fans_rate` float(3,2) NOT NULL COMMENT '互动活跃度',
  `reply_fans_rate` float(3,2) NOT NULL COMMENT '评论活跃度',
  `danmu_fans_rate` float(3,2) NOT NULL COMMENT '弹幕活跃度',
  `coin_fans_rate` float(3,2) NOT NULL COMMENT '投币活跃度',
  `like_fans_rate` float(3,2) NOT NULL COMMENT '点赞活跃度',
  `fav_fans_rate` float(3,2) NOT NULL COMMENT '收藏活跃度',
  `share_fans_rate` float(3,2) NOT NULL COMMENT '分享活跃度',
  `live_gift_fans_rate` float(3,2) NOT NULL COMMENT '直播礼物活跃度',
  `live_danmu_fans_rate` float(3,2) NOT NULL COMMENT '直播弹幕活跃度',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn粉丝数相关';


#mcn粉丝按天增量
CREATE TABLE `dm_con_mcn_fans_inc_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `fans_inc` bigint(20) NOT NULL COMMENT '当日新增粉丝数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn粉丝按天增量';


#mcn粉丝按天取关数
CREATE TABLE `dm_con_mcn_fans_dec_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `fans_dec` bigint(20) NOT NULL COMMENT '当日取关粉丝数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn粉丝按天取关数';


#mcn粉丝关注渠道
CREATE TABLE `dm_con_mcn_fans_attention_way_d` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `sign_id` bigint(20) NOT NULL COMMENT 'mcn签约ID',
  `mcn_mid` bigint(20) NOT NULL COMMENT 'mcn的mid',
  `log_date` date NOT NULL COMMENT '日期',
  `homepage` bigint(20) NOT NULL COMMENT '主站个人空间关注用户数',
  `video` bigint(20) NOT NULL COMMENT '主站视频页关注用户数',
  `article` bigint(20) NOT NULL COMMENT '专栏关注用户数',
  `music` bigint(20) NOT NULL COMMENT '音频关注用户数',
  `other` bigint(20) NOT NULL COMMENT '其他关注用户数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_idx` (`sign_id`,`log_date`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='mcn粉丝关注渠道';