CREATE TABLE keywords(
id INT(11) UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键id',
area INT(11) NOT NULL DEFAULT 0 COMMENT '业务类型',
sender_id INT(11) NOT NULL DEFAULT 0 COMMENT '发送者的id',
content VARCHAR(40) NOT NULL COMMENT '关键字内容',
regexp_name VARCHAR(40) NOT NULL COMMENT '该关键字命中正则名称',
regexp_content VARCHAR(500) NOT NULL COMMENT '正则内容',
tag tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:limit, 1:restrict, 2: whitelist, 3: blacklist',
hit_counts INT(11) NOT NULL DEFAULT 0 COMMENT '命中关键字次数',
state tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:default, 1:deleted',
origin_content VARCHAR(1500) NOT NULL COMMENT '过滤前的内容',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
UNIQUE KEY `uk_area_content` (`area`, `content`),
KEY `ix_mtime` (`mtime`),
KEY `ix_area_state_ctime` (`area`,`state`, `ctime`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='过滤限制关键字表';


CREATE TABLE regexps(
id INT(11) UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键id',
name VARCHAR(20) NOT NULL COMMENT 'name',
area INT(11) NOT NULL DEFAULT 0 COMMENT '业务类型 1: reply, 2: imessage',
operation tinyint(4) NOT NULL DEFAULT 0 COMMENT '0: limit, 1: put into whitelist, 2: restrict limit, 3: ignore',
content VARCHAR(200) NOT NULL comment '正则内容',
state tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:default, 1:deleted',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
KEY `ix_mtime` (`mtime`),
UNIQUE KEY `uk_area_content` (`area`, `content`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='过滤限制正则表';


CREATE TABLE rate_limit_rules(
id INT(11) UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键id',
area INT(11) NOT NULL DEFAULT 0 COMMENT '业务类型',
limit_type tinyint(4) NOT NULL DEFAULT 0 COMMENT '0: default, 1: strict',
limit_scope tinyint(4) NOT NULL DEFAULT 0 COMMENT '0: local, 1: global',
dur_sec int(11) NOT NULL DEFAULT 0 COMMENT '持续时间',
allowed_counts int(11) NOT NULL DEFAULT 0 COMMENT '允许发送次数',
state tinyint(4) NOT NULL DEFAULT 0 COMMENT '0:default, 1:deleted',
`ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
KEY `ix_mtime` (`mtime`)
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='频率规则表';
CREATE UNIQUE INDEX uk_area_limit_type_limit_scope ON rate_limit_rules (area, limit_type, limit_scope);
