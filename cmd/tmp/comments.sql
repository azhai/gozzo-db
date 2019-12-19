ALTER TABLE `t_access`
 COMMENT = '';

ALTER TABLE `t_group`
	CHANGE `gid` `gid` varchar(16) NOT NULL DEFAULT '' COMMENT '唯一ID',
	CHANGE `title` `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
	CHANGE `remark` `remark` varchar(255) DEFAULT NULL COMMENT '说明备注',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
 COMMENT = '用户组';

ALTER TABLE `t_menu`
	CHANGE `deleted_at` `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
	CHANGE `left_edge` `left_edge` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '左边界',
	CHANGE `seq_no` `seq_no` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '次序',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `updated_at` `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
	CHANGE `title` `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
	CHANGE `remark` `remark` varchar(255) DEFAULT NULL COMMENT '说明备注',
	CHANGE `right_edge` `right_edge` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '右边界',
 COMMENT = '菜单';

ALTER TABLE `t_role`
	CHANGE `remark` `remark` varchar(255) DEFAULT NULL COMMENT '说明备注',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `updated_at` `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
	CHANGE `deleted_at` `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
 COMMENT = '角色';

ALTER TABLE `t_user`
	CHANGE `vice_gid` `vice_gid` varchar(16) DEFAULT NULL COMMENT '次用户组',
	CHANGE `deleted_at` `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
	CHANGE `username` `username` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
	CHANGE `mobile` `mobile` varchar(20) DEFAULT NULL COMMENT '手机号码',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `updated_at` `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
	CHANGE `password` `password` varchar(60) NOT NULL DEFAULT '' COMMENT '密码',
	CHANGE `realname` `realname` varchar(20) DEFAULT NULL COMMENT '昵称/称呼',
	CHANGE `email` `email` varchar(50) DEFAULT NULL COMMENT '电子邮箱',
	CHANGE `prin_gid` `prin_gid` varchar(16) NOT NULL DEFAULT '' COMMENT '主用户组',
	CHANGE `uid` `uid` varchar(16) NOT NULL DEFAULT '' COMMENT '唯一ID',
 COMMENT = '用户';

ALTER TABLE `t_user_role`
 COMMENT = '用户角色';

