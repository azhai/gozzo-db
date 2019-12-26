-- MySQL dump 10.16  Distrib 10.1.40-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: db_test
-- ------------------------------------------------------
-- Server version	10.1.40-MariaDB


--
-- Dumping data for table `t_access`
--

INSERT INTO `t_access` VALUES (1,'superuser','menu','*',512,'all','2019-12-22 08:44:07',NULL),(2,'visitor','menu','/dashboard',2,'view','2019-12-22 08:44:07',NULL),(3,'visitor','menu','/error/404',2,'view','2019-12-22 08:44:07',NULL);

--
-- Dumping data for table `t_group`
--


--
-- Dumping data for table `t_menu`
--

INSERT INTO `t_menu` VALUES (1,1,2,1,'/dashboard','面板','dashboard',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(2,3,6,1,'/permission','权限','lock',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(3,4,5,2,'role','角色权限',NULL,NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(4,7,12,1,'/table','Table','table',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(5,8,9,2,'complex-table','复杂Table',NULL,NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(6,10,11,2,'inline-edit-table','内联编辑',NULL,NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(7,13,18,1,'/excel','Excel','excel',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(8,14,15,2,'export-selected-excel','选择导出',NULL,NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(9,16,17,2,'upload-excel','上传Excel',NULL,NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(10,19,20,1,'/theme/index','主题','theme',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(11,21,22,1,'/error/404','404错误','404',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL),(12,23,24,1,'https://cn.vuejs.org/','外部链接','link',NULL,'2019-12-24 11:51:43','2019-12-24 11:51:43',NULL);

--
-- Dumping data for table `t_role`
--

INSERT INTO `t_role` VALUES (1,'superuser','超级用户，无上权限的超级管理员。','2019-11-30 19:12:00','2019-11-30 19:12:00',NULL),(2,'member','普通用户，除权限外的其他页面。','2019-11-30 19:12:00','2019-11-30 19:12:00',NULL),(3,'visitor','基本用户，只能看到面板页。','2019-11-30 19:12:00','2019-11-30 19:12:00',NULL);

--
-- Dumping data for table `t_user`
--

INSERT INTO `t_user` VALUES (1,'6kff25twcor76222','admin','09e8ff53$x1KWXASXqGRzA7YwipQhibg/0LMtkoU39VfW8EYtxAI=','管理员',NULL,NULL,'',NULL,'/avatars/avatar-admin.jpg','不受限的超管账号。','2019-11-30 19:12:00','2019-11-30 19:12:00',NULL),(2,'6kff25u4cor76223','demo','acfd1f8b$o6ySKi7yaMmZrKIaT4O/oGUoei6n/xKOXik4PtXuvwk=','演示用户',NULL,NULL,'',NULL,'/avatars/avatar-demo.jpg','演示和测试账号。','2019-11-30 19:12:00','2019-11-30 19:12:00',NULL);

--
-- Dumping data for table `t_user_role`
--

INSERT INTO `t_user_role` VALUES (1,'6kff25twcor76222','superuser'),(2,'6kff25u4cor76223','member');


-- Dump completed on 2019-12-26 11:06:46
