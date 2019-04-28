use bilibili_apm_v2;

CREATE TABLE `overlord_appid` ( 
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键', 
    `tree_id` int(11) NOT NULL DEFAULT '0' COMMENT '服务树id', 
    `app_id` varchar(50) NOT NULL DEFAULT '' COMMENT '业务appid', 
    `cid` int(11) NOT NULL DEFAULT '0' COMMENT '关联集群id', 
    `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', 
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间', 
    PRIMARY KEY (`id`), 
    UNIQUE KEY `uk_appids_name` (`tree_id`, `cid`), 
    KEY `ix_appid` (`app_id`),
    KEY `ix_mtime` (`mtime`) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='业务关联集群';

CREATE TABLE `overlord_cluster` ( 
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键', 
    `name` varchar(50) NOT NULL DEFAULT '' COMMENT '集群名字', 
    `type` varchar(20) NOT NULL DEFAULT '' COMMENT '集群类型', 
    `zone` varchar(20) NOT NULL DEFAULT 'sh001' COMMENT '机房', 
    `hash_method` varchar(20) NOT NULL DEFAULT '' COMMENT 'hash方法', 
    `hash_distribution` varchar(20) NOT NULL DEFAULT '' COMMENT '分布策略', 
    `hashtag` varchar(10) NOT NULL DEFAULT '' COMMENT 'hash tag', 
    `listen_proto` varchar(10) NOT NULL DEFAULT 'tcp' COMMENT '监听协议', 
    `listen_addr` varchar(30) NOT NULL DEFAULT '' COMMENT '监听地址', 
    `nodeconn` int(11) NOT NULL DEFAULT '1' COMMENT '跟节点连接数', 
    `dial` int(11) NOT NULL DEFAULT '1000' COMMENT 'dial 超时', 
    `read` int(11) NOT NULL DEFAULT '1000' COMMENT 'read超时', 
    `write` int(11) NOT NULL DEFAULT '1000' COMMENT 'write 超时', 
    `ping_fail_limit` int(11) NOT NULL DEFAULT '3' COMMENT '失败检测次数', 
    `auto_eject` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'auto eject', 
    `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', 
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间', 
    PRIMARY KEY (`id`), 
    UNIQUE KEY `uk_name` (`name`), 
    KEY `ix_mtime` (`mtime`) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='集群';

CREATE TABLE `overlord_node` ( 
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键', 
    `cid` int(11) NOT NULL DEFAULT '0' COMMENT '关联集群id', 
    `alias` varchar(50) NOT NULL DEFAULT '' COMMENT '节点别名', 
    `addr` varchar(50) NOT NULL DEFAULT '' COMMENT '节点地址', 
    `weight` int(11) NOT NULL DEFAULT '1' COMMENT '节点权重', 
    `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', 
    `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间', 
    PRIMARY KEY (`id`), 
    UNIQUE KEY `uk_cid_alias` (`cid`,`alias`),
    KEY `ix_mtime` (`mtime`) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点';
