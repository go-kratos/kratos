CREATE TABLE `figure_user_[00-99]` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户信用分模型',
`mid` int(11) NOT NULL COMMENT '用户mid',
`score` smallint(6) NOT NULL DEFAULT '0' COMMENT '信用分值',
`lawful_score` smallint(6) NOT NULL DEFAULT '0' COMMENT '守序',
`wide_score` smallint(6) NOT NULL DEFAULT '0' COMMENT '博览',
`friendly_score` smallint(6) NOT NULL DEFAULT '0' COMMENT '友爱',
`bounty_score` smallint(6) NOT NULL DEFAULT '0' COMMENT '慷慨',
`creativity_score` smallint(6) NOT NULL DEFAULT '0' COMMENT '创造',
`ver` int(11) NOT NULL DEFAULT '0' COMMENT '版本',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创造世界',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_mid` (`mid`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=697525 DEFAULT CHARSET=utf8 COMMENT='信用用户表';

CREATE TABLE `figure_rank` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
`score_from` smallint(6) NOT NULL DEFAULT '0' COMMENT '起始分段分值，包含',
`score_to` smallint(6) NOT NULL DEFAULT '0' COMMENT '结束分段分值，包含',
`percentage` smallint(6) NOT NULL DEFAULT '0' COMMENT '排名百分比',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_percentage` (`percentage`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='当前信用分排名表'

CREATE TABLE `figure_rank_history` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
`score_from` smallint(6) NOT NULL DEFAULT '0' COMMENT '起始分段分值，包含',
`score_to` smallint(6) NOT NULL DEFAULT '0' COMMENT '结束分段分值，包含',
`percentage` smallint(6) NOT NULL DEFAULT '0' COMMENT '排名百分比',
`ver` int(11) NOT NULL DEFAULT '0' COMMENT '版本',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='信用分历史排名表';