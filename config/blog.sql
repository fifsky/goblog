# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.6.22)
# Database: blog
# Generation Time: 2017-08-22 02:24:44 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table cates
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cates`;

CREATE TABLE `cates` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `desc` varchar(255) NOT NULL DEFAULT '',
  `domain` varchar(100) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `un_domain` (`domain`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `cates` WRITE;
/*!40000 ALTER TABLE `cates` DISABLE KEYS */;

INSERT INTO `cates` (`id`, `name`, `desc`, `domain`, `created_at`, `updated_at`)
VALUES
	(1,'默认分类','默认分类','default','2017-08-18 15:21:56','2017-08-18 15:21:56'),
	(4,'技术笔记','技术笔记,PHP,REDIS,LINUX,MYSQL,GO','notes','2017-08-18 15:21:56','2017-08-18 15:21:56');

/*!40000 ALTER TABLE `cates` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table links
# ------------------------------------------------------------

DROP TABLE IF EXISTS `links`;

CREATE TABLE `links` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `url` varchar(200) NOT NULL DEFAULT '',
  `desc` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

LOCK TABLES `links` WRITE;
/*!40000 ALTER TABLE `links` DISABLE KEYS */;

INSERT INTO `links` (`id`, `name`, `url`, `desc`, `created_at`)
VALUES
	(1,'fifsky','http://fifsky.com','fifsky','2017-08-18 15:21:56');

/*!40000 ALTER TABLE `links` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table moods
# ------------------------------------------------------------

DROP TABLE IF EXISTS `moods`;

CREATE TABLE `moods` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL DEFAULT '',
  `user_id` int(10) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `moods` WRITE;
/*!40000 ALTER TABLE `moods` DISABLE KEYS */;

INSERT INTO `moods` (`id`, `content`, `user_id`, `created_at`)
VALUES
	(1,'Hi,fifsky!',1,'2017-08-18 15:21:56');

/*!40000 ALTER TABLE `moods` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table options
# ------------------------------------------------------------

DROP TABLE IF EXISTS `options`;

CREATE TABLE `options` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `option_key` varchar(100) NOT NULL DEFAULT '',
  `option_value` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `option_name` (`option_key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `options` WRITE;
/*!40000 ALTER TABLE `options` DISABLE KEYS */;

INSERT INTO `options` (`id`, `option_key`, `option_value`)
VALUES
	(1,'site_name','無處告別'),
	(2,'site_desc','回首往事，珍重眼前人'),
	(3,'site_keyword','fifsky,rita,生活,博客,豆豆'),
	(4,'post_num','10');

/*!40000 ALTER TABLE `options` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table posts
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posts`;

CREATE TABLE `posts` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `cate_id` int(11) unsigned NOT NULL DEFAULT '1',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1:文章,2:页面',
  `user_id` int(11) unsigned NOT NULL,
  `title` varchar(200) NOT NULL DEFAULT '',
  `url` varchar(100) NOT NULL DEFAULT '' COMMENT '页面缩略名',
  `content` longtext NOT NULL,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `post_author` (`user_id`),
  KEY `type_status_date` (`id`,`type`),
  KEY `post_name` (`url`) USING BTREE,
  KEY `post_title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;

INSERT INTO `posts` (`id`, `cate_id`, `type`, `user_id`, `title`, `url`, `content`, `created_at`, `updated_at`)
VALUES
	(37,1,1,1,'fifsky blog for go','','<p>\r\n	https://github.com/fifsky/goblog\r\n</p>\r\n<p>\r\n	<br />\r\n</p>','2017-08-22 10:17:23','2017-08-22 10:17:23');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table comments
# ------------------------------------------------------------
DROP TABLE IF EXISTS `comments`;

CREATE TABLE `comments` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `post_id` int(11) NOT NULL COMMENT '文章PID',
  `pid` int(11) NOT NULL COMMENT '回复评论ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `content` tinytext NOT NULL COMMENT '内容',
  `ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'IP',
  `created_at` datetime NOT NULL COMMENT '评论时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;


# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `password` varchar(100) NOT NULL DEFAULT '',
  `nick_name` varchar(100) NOT NULL DEFAULT '',
  `email` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1正常，2删除',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1:管理员,2:编辑',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `un_user_name` (`name`),
  UNIQUE KEY `uix_users_email` (`email`),
  UNIQUE KEY `uix_users_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`id`, `name`, `password`, `nick_name`, `email`, `status`, `type`, `created_at`, `updated_at`)
VALUES
	(1,'test','e10adc3949ba59abbe56e057f20f883e','test','test@test.com',1,1,'2017-08-18 15:21:56','2017-08-18 15:21:56');

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table reminds
# ------------------------------------------------------------

DROP TABLE IF EXISTS `reminds`;

CREATE TABLE `reminds` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` int(11) NOT NULL COMMENT '0固定，1每分钟，2每个小时，3每周，4，每天，5，每月，6，每年',
  `at` varchar(20) NOT NULL DEFAULT '' COMMENT '@手机号',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
  `remind_date` datetime NOT NULL COMMENT '提醒日期',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
