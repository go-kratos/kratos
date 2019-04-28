CREATE TABLE `up_base_info` (
  `id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '自增ID',
  `mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  `name` varchar(36) NOT NULL DEFAULT '' COMMENT '昵称',
  `sex` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别 0:保密 1:男 2:女',
  `join_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '注册时间',
  `first_up_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '第一次投稿时间',
  `level` smallint(6) NOT NULL DEFAULT '0' COMMENT '等级',
  `fans_count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝量',
  `account_state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '账号状态，0：正常，1：封禁',
  `activity` int(11) NOT NULL DEFAULT '0' COMMENT '活跃度',
  `article_count_30day` int(11) NOT NULL DEFAULT '0' COMMENT '30天内投稿量(所有业务)',
  `article_count_accumulate` int(11) NOT NULL DEFAULT '0' COMMENT '累计投稿量(各业务累加)',
  `verify_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '认证类型,0-个人，1-企业',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `business_type` tinyint(4) NOT NULL COMMENT '业务类型, 1视频/2音频/3专栏',
  `credit_score` int(11) NOT NULL DEFAULT '500' COMMENT '信用分',
  `pr_score` int(11) NOT NULL DEFAULT '0' COMMENT '影响分',
  `quality_score` int(11) NOT NULL DEFAULT '0' COMMENT '质量分',
  `active_tid` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '最多稿件分区',
  `attr` int(11) NOT NULL DEFAULT '0' COMMENT '属性，以位区分',
  `birthday` date NOT NULL DEFAULT '0000-00-00' COMMENT '生日',
  `active_province` varchar(32) NOT NULL DEFAULT '' COMMENT '省份',
  `active_city` varchar(32) NOT NULL DEFAULT '' COMMENT '城市',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mid_type` (`mid`,`business_type`),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_uptime` (`first_up_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='up基本信息表';

CREATE TABLE `up_play_info` (
  `id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '自增ID',
  `mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  `business_type` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '业务类型, 1视频/2音频/3专栏',
  `play_count_accumulate` bigint(20) NOT NULL DEFAULT '0' COMMENT '累计播放次数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `article_count` int(11) NOT NULL DEFAULT '0' COMMENT '总稿件数',
  `play_count_90day` bigint(20) NOT NULL DEFAULT '0' COMMENT '90天内稿件总播放次数',
  `play_count_30day` bigint(20) NOT NULL DEFAULT '0' COMMENT '30天内稿件总播放次数',
  `play_count_7day` bigint(20) NOT NULL DEFAULT '0' COMMENT '7天内稿件总播放次数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mid_bus` (`mid`,`business_type`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='up基本播放信息表';

CREATE TABLE `up_rank` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `mid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  `type` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT '排行榜类型',
  `value` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '排行榜数值,根据type不同，代表的含义不同',
  `generate_date` date NOT NULL DEFAULT '0000-00-00' COMMENT '排行榜日',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `value2` int(11) NOT NULL DEFAULT '0' COMMENT '分数2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_date_type_mid` (`generate_date`,`type`,`mid`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=35001 DEFAULT CHARSET=utf8 COMMENT='Up蹿升榜';

CREATE TABLE `task_info` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `generate_date` date NOT NULL DEFAULT '0000-00-00' COMMENT '计算日',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `start_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '结束时间',
  `task_state` smallint(6) NOT NULL DEFAULT '0' COMMENT '任务状态, 0表示初始化, 1表示结束，其他状态根据task_type定义',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_date_type` (`generate_date`,`task_type`),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_date` (`generate_date`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8 COMMENT='计算任务信息表';

-- score_section_history: table
CREATE TABLE `score_section_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `generate_date` date NOT NULL COMMENT '生成日期',
  `score_type` smallint(6) NOT NULL DEFAULT '0' COMMENT '类型, 1质量分，2影响分，3信用分',
  `section_0` int(11) NOT NULL DEFAULT '0' COMMENT '0~100的人数',
  `section_1` int(11) NOT NULL DEFAULT '0' COMMENT '101~200',
  `section_2` int(11) NOT NULL DEFAULT '0' COMMENT '201~300',
  `section_3` int(11) NOT NULL DEFAULT '0' COMMENT '301~400',
  `section_4` int(11) NOT NULL DEFAULT '0' COMMENT '401~500',
  `section_5` int(11) NOT NULL DEFAULT '0' COMMENT '501~600',
  `section_6` int(11) NOT NULL DEFAULT '0' COMMENT '601~700',
  `section_7` int(11) NOT NULL DEFAULT '0' COMMENT '701~800',
  `section_8` int(11) NOT NULL DEFAULT '0' COMMENT '801~900',
  `section_9` int(11) NOT NULL DEFAULT '0' COMMENT '901~1000',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_date_type` (`generate_date`,`score_type`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8 COMMENT='up分数段人数分布表'
;

-- up_stats_history: table
CREATE TABLE `up_stats_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1活跃、2新增（可以通过累计来计算）、3累计up主人数',
  `sub_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '子类',
  `value1` int(11) NOT NULL DEFAULT '0' COMMENT '分数',
  `value2` int(11) NOT NULL DEFAULT '0' COMMENT '分数2,备用',
  `generate_date` date NOT NULL DEFAULT '0000-00-00' COMMENT '生成日期',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_type_subtype_date` (`generate_date`,`type`,`sub_type`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8 COMMENT='up主总数表'
;

-- up_scores_history_00: table
CREATE TABLE `up_scores_history_00` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `mid` int(11) unsigned NOT NULL COMMENT 'up主id',
  `score_type` tinyint(4) NOT NULL COMMENT '1内容分，2影响分，3信用分',
  `score` int(11) unsigned NOT NULL COMMENT '分数',
  `generate_date` date NOT NULL COMMENT '生成日期',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mid_scoretype_date` (`mid`,`score_type`,`generate_date`),
  KEY `ix_date` (`generate_date`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=101239 DEFAULT CHARSET=utf8 COMMENT='up分数表'
;

-- credit_log_00: table
CREATE TABLE `credit_log_00` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `type` smallint(6) unsigned NOT NULL COMMENT '日志类型',
  `op_type` smallint(6) NOT NULL COMMENT '操作类型',
  `reason` smallint(6) NOT NULL COMMENT '原因类型',
  `business_type` smallint(6) unsigned NOT NULL COMMENT '业务类型',
  `mid` int(11) NOT NULL COMMENT 'mid',
  `oid` int(11) unsigned NOT NULL COMMENT 'oid',
  `uid` smallint(6) unsigned NOT NULL COMMENT '操作人员id',
  `content` varchar(255) NOT NULL COMMENT '操作内容',
  `extra` varchar(2000) NOT NULL DEFAULT '' COMMENT '额外信息',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mtime` (`mtime`),
  KEY `ix_mid` (`mid`)
) ENGINE=InnoDB AUTO_INCREMENT=394475 DEFAULT CHARSET=utf8 COMMENT='信用分日志表'
;
CREATE TABLE credit_log_00 LIKE credit_log_00;
CREATE TABLE credit_log_01 LIKE credit_log_00;
CREATE TABLE credit_log_02 LIKE credit_log_00;
CREATE TABLE credit_log_03 LIKE credit_log_00;
CREATE TABLE credit_log_04 LIKE credit_log_00;
CREATE TABLE credit_log_05 LIKE credit_log_00;
CREATE TABLE credit_log_06 LIKE credit_log_00;
CREATE TABLE credit_log_07 LIKE credit_log_00;
CREATE TABLE credit_log_08 LIKE credit_log_00;
CREATE TABLE credit_log_09 LIKE credit_log_00;
CREATE TABLE credit_log_10 LIKE credit_log_00;
CREATE TABLE credit_log_11 LIKE credit_log_00;
CREATE TABLE credit_log_12 LIKE credit_log_00;
CREATE TABLE credit_log_13 LIKE credit_log_00;
CREATE TABLE credit_log_14 LIKE credit_log_00;
CREATE TABLE credit_log_15 LIKE credit_log_00;
CREATE TABLE credit_log_16 LIKE credit_log_00;
CREATE TABLE credit_log_17 LIKE credit_log_00;
CREATE TABLE credit_log_18 LIKE credit_log_00;
CREATE TABLE credit_log_19 LIKE credit_log_00;
CREATE TABLE credit_log_20 LIKE credit_log_00;
CREATE TABLE credit_log_21 LIKE credit_log_00;
CREATE TABLE credit_log_22 LIKE credit_log_00;
CREATE TABLE credit_log_23 LIKE credit_log_00;
CREATE TABLE credit_log_24 LIKE credit_log_00;
CREATE TABLE credit_log_25 LIKE credit_log_00;
CREATE TABLE credit_log_26 LIKE credit_log_00;
CREATE TABLE credit_log_27 LIKE credit_log_00;
CREATE TABLE credit_log_28 LIKE credit_log_00;
CREATE TABLE credit_log_29 LIKE credit_log_00;
CREATE TABLE credit_log_30 LIKE credit_log_00;
CREATE TABLE credit_log_31 LIKE credit_log_00;
CREATE TABLE credit_log_32 LIKE credit_log_00;
CREATE TABLE credit_log_33 LIKE credit_log_00;
CREATE TABLE credit_log_34 LIKE credit_log_00;
CREATE TABLE credit_log_35 LIKE credit_log_00;
CREATE TABLE credit_log_36 LIKE credit_log_00;
CREATE TABLE credit_log_37 LIKE credit_log_00;
CREATE TABLE credit_log_38 LIKE credit_log_00;
CREATE TABLE credit_log_39 LIKE credit_log_00;
CREATE TABLE credit_log_40 LIKE credit_log_00;
CREATE TABLE credit_log_41 LIKE credit_log_00;
CREATE TABLE credit_log_42 LIKE credit_log_00;
CREATE TABLE credit_log_43 LIKE credit_log_00;
CREATE TABLE credit_log_44 LIKE credit_log_00;
CREATE TABLE credit_log_45 LIKE credit_log_00;
CREATE TABLE credit_log_46 LIKE credit_log_00;
CREATE TABLE credit_log_47 LIKE credit_log_00;
CREATE TABLE credit_log_48 LIKE credit_log_00;
CREATE TABLE credit_log_49 LIKE credit_log_00;
CREATE TABLE credit_log_50 LIKE credit_log_00;
CREATE TABLE credit_log_51 LIKE credit_log_00;
CREATE TABLE credit_log_52 LIKE credit_log_00;
CREATE TABLE credit_log_53 LIKE credit_log_00;
CREATE TABLE credit_log_54 LIKE credit_log_00;
CREATE TABLE credit_log_55 LIKE credit_log_00;
CREATE TABLE credit_log_56 LIKE credit_log_00;
CREATE TABLE credit_log_57 LIKE credit_log_00;
CREATE TABLE credit_log_58 LIKE credit_log_00;
CREATE TABLE credit_log_59 LIKE credit_log_00;
CREATE TABLE credit_log_60 LIKE credit_log_00;
CREATE TABLE credit_log_61 LIKE credit_log_00;
CREATE TABLE credit_log_62 LIKE credit_log_00;
CREATE TABLE credit_log_63 LIKE credit_log_00;
CREATE TABLE credit_log_64 LIKE credit_log_00;
CREATE TABLE credit_log_65 LIKE credit_log_00;
CREATE TABLE credit_log_66 LIKE credit_log_00;
CREATE TABLE credit_log_67 LIKE credit_log_00;
CREATE TABLE credit_log_68 LIKE credit_log_00;
CREATE TABLE credit_log_69 LIKE credit_log_00;
CREATE TABLE credit_log_70 LIKE credit_log_00;
CREATE TABLE credit_log_71 LIKE credit_log_00;
CREATE TABLE credit_log_72 LIKE credit_log_00;
CREATE TABLE credit_log_73 LIKE credit_log_00;
CREATE TABLE credit_log_74 LIKE credit_log_00;
CREATE TABLE credit_log_75 LIKE credit_log_00;
CREATE TABLE credit_log_76 LIKE credit_log_00;
CREATE TABLE credit_log_77 LIKE credit_log_00;
CREATE TABLE credit_log_78 LIKE credit_log_00;
CREATE TABLE credit_log_79 LIKE credit_log_00;
CREATE TABLE credit_log_80 LIKE credit_log_00;
CREATE TABLE credit_log_81 LIKE credit_log_00;
CREATE TABLE credit_log_82 LIKE credit_log_00;
CREATE TABLE credit_log_83 LIKE credit_log_00;
CREATE TABLE credit_log_84 LIKE credit_log_00;
CREATE TABLE credit_log_85 LIKE credit_log_00;
CREATE TABLE credit_log_86 LIKE credit_log_00;
CREATE TABLE credit_log_87 LIKE credit_log_00;
CREATE TABLE credit_log_88 LIKE credit_log_00;
CREATE TABLE credit_log_89 LIKE credit_log_00;
CREATE TABLE credit_log_90 LIKE credit_log_00;
CREATE TABLE credit_log_91 LIKE credit_log_00;
CREATE TABLE credit_log_92 LIKE credit_log_00;
CREATE TABLE credit_log_93 LIKE credit_log_00;
CREATE TABLE credit_log_94 LIKE credit_log_00;
CREATE TABLE credit_log_95 LIKE credit_log_00;
CREATE TABLE credit_log_96 LIKE credit_log_00;
CREATE TABLE credit_log_97 LIKE credit_log_00;
CREATE TABLE credit_log_98 LIKE credit_log_00;
CREATE TABLE credit_log_99 LIKE credit_log_00;

CREATE TABLE up_scores_history_00 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_01 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_02 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_03 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_04 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_05 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_06 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_07 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_08 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_09 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_10 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_11 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_12 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_13 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_14 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_15 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_16 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_17 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_18 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_19 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_20 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_21 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_22 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_23 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_24 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_25 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_26 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_27 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_28 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_29 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_30 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_31 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_32 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_33 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_34 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_35 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_36 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_37 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_38 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_39 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_40 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_41 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_42 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_43 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_44 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_45 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_46 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_47 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_48 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_49 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_50 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_51 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_52 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_53 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_54 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_55 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_56 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_57 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_58 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_59 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_60 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_61 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_62 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_63 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_64 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_65 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_66 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_67 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_68 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_69 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_70 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_71 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_72 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_73 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_74 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_75 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_76 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_77 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_78 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_79 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_80 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_81 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_82 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_83 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_84 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_85 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_86 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_87 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_88 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_89 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_90 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_91 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_92 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_93 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_94 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_95 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_96 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_97 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_98 LIKE up_scores_history_00;
CREATE TABLE up_scores_history_99 LIKE up_scores_history_00;
--------- 以上是建表语句 -----------

--------- 修改 up_play_info表字段，有数据溢出 -----------
# (fat 1, uat 1, prod 1)
alter table up_play_info modify play_count_accumulate BIGINT(20) NOT NULL DEFAULT '0' COMMENT '累计播放次数';
alter table up_play_info modify play_count_90day BIGINT(20) NOT NULL DEFAULT '0' COMMENT '90天内稿件总播放次数';
alter table up_play_info modify play_count_30day BIGINT(20) NOT NULL DEFAULT '0' COMMENT '30天内稿件总播放次数';
alter table up_play_info modify play_count_7day BIGINT(20) NOT NULL DEFAULT '0' COMMENT '7天内稿件总播放次数';

--------- 修改 自增ID为bigint unsigned -----------
# (fat 1, uat 1, prod 1)
alter table up_play_info modify id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '自增ID';
alter table up_base_info modify id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '自增ID';
