CREATE TABLE `filter_area` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID, 帖子的评论和回复ID',
`area` varchar(50) NOT NULL COMMENT '业务方',
`typeid` smallint(6) NOT NULL COMMENT '分区id',
`filterid` int(11) NOT NULL COMMENT '过滤内容id',
`level` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '业务过滤等级',
`is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除：0:未删除 1:已删除',
`ctime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
`mtime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_area_filterid_typeid` (`area`,`filterid`,`typeid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='敏感词对应业务表';

CREATE TABLE `filter_content` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID, 过滤id',
`mode` tinyint(4) NOT NULL COMMENT '过滤模式 0-正则 ，1-string',
`filter` varchar(255) NOT NULL COMMENT '过滤内容',
`level` tinyint(4) NOT NULL COMMENT '过滤等级',
`comment` varchar(128) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '过滤备注',
`ctime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
`mtime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
`stime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '生效时间',
`etime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '失效时间',
`key` varchar(30) NOT NULL DEFAULT '' COMMENT '业务方指定key',
`state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除，0:未删除 1:已删除',
`type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '类型：0 其他违禁类;1 政治宗教;2 色情;3 低俗;4 血腥暴力;5 赌博诈骗;6 运营规避',
`source` tinyint(4) NOT NULL DEFAULT '0' COMMENT '来源：0 上级部门;2 审核规避;4 运营规避;8 审核提示',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_filter_key` (`filter`,`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='过滤内容';

CREATE TABLE `filter_key` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
`key` varchar(30) NOT NULL COMMENT '业务方指定key',
`area` varchar(10) NOT NULL COMMENT '过滤区域',
`filterid` int(11) NOT NULL COMMENT '过滤内容id',
`state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除，0:未删除 1:已删除',
`ctime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
`mtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
PRIMARY KEY (`id`),
UNIQUE KEY `ux_key_area_filterid` (`key`,`area`,`filterid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='key过滤表';

CREATE TABLE `filter_key_log` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
`key` varchar(20) NOT NULL DEFAULT '' COMMENT '业务方指定key',
`adid` int(11) NOT NULL COMMENT '管理员id',
`name` varchar(16) NOT NULL COMMENT 'name',
`comment` varchar(50) NOT NULL COMMENT '操作原因',
`state` tinyint(4) NOT NULL COMMENT '操作类型 0-添加,1-编辑, 2-删除',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
PRIMARY KEY (`id`),
KEY `ix_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='key操作日志表';

CREATE TABLE `filter_log` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID, 过滤id',
`adid` int(11) NOT NULL COMMENT '管理员id',
`comment` varchar(50) NOT NULL COMMENT '操作原因',
`state` tinyint(4) NOT NULL COMMENT '操作类型 0-添加,1-编辑, 2-删除',
`ctime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
`mtime` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
`filterid` int(11) NOT NULL DEFAULT '0' COMMENT '敏感词id',
`name` varchar(16) NOT NULL COMMENT '名字',
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='操作日志';

CREATE TABLE `filter_white_area` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
`area` varchar(20) NOT NULL COMMENT '业务方',
`tpid` int(11) NOT NULL COMMENT '分区id',
`content_id` int(11) NOT NULL COMMENT '内容id',
`state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，1：删除',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_area_tyid_contentid` (`content_id`,`tpid`,`area`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='过滤白名单业务关系表';

CREATE TABLE `filter_white_content` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
`content` varchar(64) NOT NULL COMMENT '过滤规则',
`mode` tinyint(4) NOT NULL COMMENT '模式，0：正则，1：字符串',
`comment` varchar(50) NOT NULL COMMENT '说明',
`state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0：正常，1：删除',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_content` (`content`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='过滤白名单内容表';

CREATE TABLE `filter_white_log` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID, 过滤id',
`content_id` int(11) NOT NULL COMMENT '内容id',
`adid` int(11) NOT NULL COMMENT '管理员id',
`name` varchar(20) NOT NULL COMMENT '用户昵称',
`comment` varchar(50) NOT NULL COMMENT '操作原因',
`state` tinyint(4) NOT NULL COMMENT '操作类型 0-添加,1-编辑, 2-删除',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
PRIMARY KEY (`id`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='过滤日志操作记录表';
 
CREATE TABLE `filter_area_type` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`name` VARCHAR(50) NOT NULL COMMENT '业务名称',
`showname` VARCHAR(16) NOT NULL COMMENT '业务显示名称',
`groupid` INT(11) NOT NULL COMMENT '业务分组id',
`common_flag` TINYINT(4) NOT NULL COMMENT '是否过滤基础库',
`is_delete` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '是否删除：0:未删除 1:已删除',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_area` (`name`),
INDEX `ix_groupid` (`groupid`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='敏感词业务类型表';
 
CREATE TABLE `filter_area_group` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`name` VARCHAR(16) NOT NULL COMMENT '业务分组名称',
`is_delete` TINYINT(4) NOT NULL DEFAULT '0' COMMENT '是否删除：0:未删除 1:已删除',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_name` (`name`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='敏感词业务分组表';
 
CREATE TABLE `filter_area_type_log` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`areaid` INT(11) NOT NULL COMMENT '业务id',
`state` TINYINT(4) NOT NULL COMMENT '操作',
`adid` INT(11) NOT NULL COMMENT '管理员id',
`ad_name` VARCHAR(16) NOT NULL COMMENT '管理员名称',
`comment` VARCHAR(50) NOT NULL COMMENT '变动理由',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
INDEX `ix_areaid` (`areaid`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='敏感词业务日志';
 
CREATE TABLE `filter_area_group_log` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`groupid` INT(11) NOT NULL COMMENT '业务组id',
`state` TINYINT(4) NOT NULL COMMENT '操作',
`adid` INT(11) NOT NULL COMMENT '管理员id',
`ad_name` VARCHAR(16) NOT NULL COMMENT '管理员名称',
`comment` VARCHAR(50) NOT NULL COMMENT '变动理由',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
INDEX `ix_groupid` (`groupid`),
KEY `ix_mtime` (`mtime`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='敏感词业务组日志';