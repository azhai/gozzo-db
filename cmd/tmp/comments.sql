ALTER TABLE `t_access`
	CHANGE `resource_type` `resource_type` varchar(50) NOT NULL DEFAULT '' COMMENT '资源类型',
	CHANGE `resource_args` `resource_args` varchar(255) DEFAULT NULL COMMENT '资源参数',
	CHANGE `perm_code` `perm_code` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '权限码',
	CHANGE `actions` `actions` varchar(50) NOT NULL DEFAULT '' COMMENT '允许的操作',
	CHANGE `granted_at` `granted_at` timestamp DEFAULT NULL COMMENT '授权时间',
	CHANGE `revoked_at` `revoked_at` timestamp DEFAULT NULL COMMENT '撤销时间',
	CHANGE `role_name` `role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
 COMMENT = '权限控制';

ALTER TABLE `t_group`
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `gid` `gid` char(16) NOT NULL DEFAULT '' COMMENT '唯一ID',
	CHANGE `title` `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
	CHANGE `remark` `remark` text DEFAULT NULL COMMENT '说明备注',
 COMMENT = '用户组';

ALTER TABLE `t_menu`
	CHANGE `deleted_at` `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
	CHANGE `path` `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',
	CHANGE `updated_at` `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
	CHANGE `rgt` `rgt` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '右边界',
	CHANGE `depth` `depth` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '高度',
	CHANGE `title` `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
	CHANGE `icon` `icon` varchar(30) DEFAULT NULL COMMENT '图标',
	CHANGE `remark` `remark` text DEFAULT NULL COMMENT '说明备注',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `lft` `lft` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '左边界',
 COMMENT = '菜单';

ALTER TABLE `t_role`
	CHANGE `name` `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
	CHANGE `remark` `remark` text DEFAULT NULL COMMENT '说明备注',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `updated_at` `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
	CHANGE `deleted_at` `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
 COMMENT = '角色';

ALTER TABLE `t_user`
	CHANGE `uid` `uid` char(16) NOT NULL DEFAULT '' COMMENT '唯一ID',
	CHANGE `vice_gid` `vice_gid` char(16) DEFAULT NULL COMMENT '次用户组',
	CHANGE `created_at` `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
	CHANGE `username` `username` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
	CHANGE `password` `password` varchar(60) NOT NULL DEFAULT '' COMMENT '密码',
	CHANGE `realname` `realname` varchar(20) DEFAULT NULL COMMENT '昵称/称呼',
	CHANGE `avatar` `avatar` varchar(100) DEFAULT NULL COMMENT '头像',
	CHANGE `deleted_at` `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
	CHANGE `email` `email` varchar(50) DEFAULT NULL COMMENT '电子邮箱',
	CHANGE `introduction` `introduction` varchar(500) DEFAULT NULL COMMENT '介绍说明',
	CHANGE `updated_at` `updated_at` timestamp DEFAULT NULL COMMENT '更新时间',
	CHANGE `mobile` `mobile` varchar(20) DEFAULT NULL COMMENT '手机号码',
	CHANGE `prin_gid` `prin_gid` char(16) NOT NULL DEFAULT '' COMMENT '主用户组',
 COMMENT = '用户';

ALTER TABLE `t_user_role`
	CHANGE `user_uid` `user_uid` char(16) NOT NULL DEFAULT '' COMMENT '用户ID',
	CHANGE `role_name` `role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
 COMMENT = '用户角色';

