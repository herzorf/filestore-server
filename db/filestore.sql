# ************************************************************
# Sequel Ace SQL dump
# 版本号： 20044
#
# https://sequel-ace.com/
# https://github.com/Sequel-Ace/Sequel-Ace
#
# 主机: 127.0.0.1 (MySQL 8.0.31)
# 数据库: filestore
# 生成时间: 2023-01-30 04:10:20 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE='NO_AUTO_VALUE_ON_ZERO', SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# 转储表 file
# ------------------------------------------------------------

CREATE DATABASE `filestore`;

USE `filestore`;


DROP TABLE IF EXISTS `file`;

CREATE TABLE `file` (
                        `id` int NOT NULL AUTO_INCREMENT,
                        `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
                        `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
                        `file_size` bigint DEFAULT '0' COMMENT '文件大小',
                        `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
                        `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
                        `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日期',
                        `status` int NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
                        `ext1` int DEFAULT '0' COMMENT '备用字段1',
                        `ext2` text COMMENT '备用字段2',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_file_hash` (`file_sha1`),
                        KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=71 DEFAULT CHARSET=utf8mb3;

LOCK TABLES `file` WRITE;
/*!40000 ALTER TABLE `file` DISABLE KEYS */;

INSERT INTO `file` (`id`, `file_sha1`, `file_name`, `file_size`, `file_addr`, `create_at`, `update_at`, `status`, `ext1`, `ext2`)
VALUES
    (62,'5f16aa2fe17459ce29ab6daa66e42421dd833e2f','萌兔兔.png',143373,'temp/萌兔兔.png','2023-01-28 13:46:09','2023-01-28 13:46:09',1,0,NULL),
    (63,'005af40df5ccc16cf6203b3561f2b328fe1e0dda','兔耳朵.png',65269,'temp/兔耳朵.png','2023-01-28 13:47:38','2023-01-28 13:47:38',1,0,NULL),
    (64,'d2ed3d37f46d8a877847f57eafab1c855cdd58dc','下载.jpeg',3841,'temp/下载.jpeg','2023-01-28 14:02:03','2023-01-28 14:02:03',1,0,NULL),
    (65,'c4e5695a2c58163e92116bcbd6b28858c7e70a9d','flower.jpeg',2690,'temp/flower.jpeg','2023-01-28 16:11:17','2023-01-28 16:11:17',1,0,NULL),
    (69,'9d8a8e015fb03fe5b6fe4821e7a4384ae165e96d','istockphoto-1386867464-612x612.jpg',22867,'temp/istockphoto-1386867464-612x612.jpg','2023-01-29 12:12:35','2023-01-29 12:12:35',1,0,NULL);

/*!40000 ALTER TABLE `file` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
                        `id` int NOT NULL AUTO_INCREMENT,
                        `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
                        `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
                        `email` varchar(64) DEFAULT '' COMMENT '邮箱',
                        `phone` varchar(128) DEFAULT '' COMMENT '手机号',
                        `email_validated` tinyint(1) DEFAULT '0' COMMENT '邮箱是否已验证',
                        `phone_validated` tinyint(1) DEFAULT '0' COMMENT '手机号是否已验证',
                        `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
                        `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
                        `profile` text COMMENT '用户属性',
                        `status` int NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_username` (`user_name`),
                        KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `user_name`, `user_pwd`, `email`, `phone`, `email_validated`, `phone_validated`, `signup_at`, `last_active`, `profile`, `status`)
VALUES
    (17,'hezhongfeng','5f8016401d91b4a81b3b786bac0ace40512f50f3','','',0,0,'2023-01-15 07:08:08','2023-01-15 07:08:08',NULL,0),
    (18,'admin1111111','5edda127dbef108cb0fc490317f881e2996c6802','','',0,0,'2023-01-17 12:45:34','2023-01-17 12:45:34',NULL,0),
    (20,'admin123456','25ee7fc896a5c6b12d7a5421ecea2c80296a760e','','',0,0,'2023-01-21 09:25:21','2023-01-21 09:25:21',NULL,0),
    (25,'1111111','a6e40b9bd5082199a58251ef561e3bdfba1be15f','','',0,0,'2023-01-21 11:01:27','2023-01-21 11:01:27',NULL,0),
    (26,'222222','45c4681fea82a3ca6060f7fac5f379cfe693d8cd','','',0,0,'2023-01-25 14:25:27','2023-01-25 14:25:27',NULL,0),
    (35,'333333','e7e58427bd1907d48367109d31139a2583984567','','',0,0,'2023-01-25 15:42:27','2023-01-25 15:42:27',NULL,0);

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;


# 转储表 user_file
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_file`;

CREATE TABLE `user_file` (
                             `id` int NOT NULL AUTO_INCREMENT,
                             `user_name` varchar(64) NOT NULL,
                             `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
                             `file_size` bigint DEFAULT '0' COMMENT '文件大小',
                             `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
                             `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
                             `last_update` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                             `status` int NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
                             PRIMARY KEY (`id`),
                             UNIQUE KEY `idx_user_file` (`user_name`,`file_sha1`),
                             KEY `idx_status` (`status`),
                             KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=78 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



# 转储表 user_token
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_token`;

CREATE TABLE `user_token` (
                              `id` int NOT NULL AUTO_INCREMENT,
                              `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
                              `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登录token',
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=67 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

LOCK TABLES `user_token` WRITE;
/*!40000 ALTER TABLE `user_token` DISABLE KEYS */;

INSERT INTO `user_token` (`id`, `user_name`, `user_token`)
VALUES
    (66,'222222','3863dc0261eab9fc496977f45280e0a263d6a19f');

/*!40000 ALTER TABLE `user_token` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
