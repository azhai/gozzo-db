-- MySQL dump 10.16  Distrib 10.1.40-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: db_test
-- ------------------------------------------------------
-- Server version	10.1.40-MariaDB


--
-- Table structure for table `t_access`
--

DROP TABLE IF EXISTS `t_access`;
CREATE TABLE `t_access` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
  `resource_type` varchar(50) NOT NULL DEFAULT '' COMMENT '资源类型',
  `resource_args` varchar(255) DEFAULT NULL COMMENT '资源参数',
  `perm_code` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '权限码',
  `actions` varchar(50) NOT NULL DEFAULT '' COMMENT '允许的操作',
  `granted_at` timestamp NULL DEFAULT NULL COMMENT '授权时间',
  `revoked_at` timestamp NULL DEFAULT NULL COMMENT '撤销时间',
  PRIMARY KEY (`id`),
  KEY `idx_t_access_role_name` (`role_name`),
  KEY `idx_t_access_revoked_at` (`revoked_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限控制';

--
-- Table structure for table `t_group`
--

DROP TABLE IF EXISTS `t_group`;
CREATE TABLE `t_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `gid` char(16) NOT NULL DEFAULT '' COMMENT '唯一ID',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `remark` text COMMENT '说明备注',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_t_group_gid` (`gid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户组';

--
-- Table structure for table `t_menu`
--

DROP TABLE IF EXISTS `t_menu`;
CREATE TABLE `t_menu` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `lft` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '左边界',
  `rgt` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '右边界',
  `depth` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '高度',
  `path` varchar(100) NOT NULL DEFAULT '' COMMENT '路径',
  `title` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `icon` varchar(30) DEFAULT NULL COMMENT '图标',
  `remark` text COMMENT '说明备注',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_t_menu_rgt` (`rgt`),
  KEY `idx_t_menu_depth` (`depth`),
  KEY `idx_t_menu_path` (`path`),
  KEY `idx_t_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单';

--
-- Table structure for table `t_role`
--

DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `remark` text COMMENT '说明备注',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_t_role_name` (`name`),
  KEY `idx_t_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色';

--
-- Table structure for table `t_user`
--

DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` char(16) NOT NULL DEFAULT '' COMMENT '唯一ID',
  `username` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(60) NOT NULL DEFAULT '' COMMENT '密码',
  `realname` varchar(20) DEFAULT NULL COMMENT '昵称/称呼',
  `mobile` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `email` varchar(50) DEFAULT NULL COMMENT '电子邮箱',
  `prin_gid` char(16) NOT NULL DEFAULT '' COMMENT '主用户组',
  `vice_gid` char(16) DEFAULT NULL COMMENT '次用户组',
  `avatar` varchar(100) DEFAULT NULL COMMENT '头像',
  `introduction` varchar(500) DEFAULT NULL COMMENT '介绍说明',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_t_user_uid` (`uid`),
  KEY `idx_t_user_username` (`username`),
  KEY `idx_t_user_mobile` (`mobile`),
  KEY `idx_t_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户';

--
-- Table structure for table `t_user_role`
--

DROP TABLE IF EXISTS `t_user_role`;
CREATE TABLE `t_user_role` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_uid` char(16) NOT NULL DEFAULT '' COMMENT '用户ID',
  `role_name` varchar(50) NOT NULL DEFAULT '' COMMENT '角色名',
  PRIMARY KEY (`id`),
  KEY `idx_t_user_role_user_uid` (`user_uid`),
  KEY `idx_t_user_role_role_name` (`role_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色';


-- Dump completed on 2019-12-26 11:06:46
