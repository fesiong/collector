/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MariaDB
 Source Server Version : 100410
 Source Host           : localhost:3306
 Source Schema         : collector

 Target Server Type    : MariaDB
 Target Server Version : 100410
 File Encoding         : 65001

 Date: 11/07/2020 10:27:20
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for fe_article
-- ----------------------------
DROP TABLE IF EXISTS `fe_article`;
CREATE TABLE `fe_article` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `source_id` int(11) NOT NULL DEFAULT 0,
  `title` varchar(250) NOT NULL DEFAULT '',
  `keywords` varchar(250) NOT NULL DEFAULT '',
  `description` varchar(250) NOT NULL DEFAULT '',
  `article_type` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '文章类型：0默认，1政策，2新闻，3',
  `origin_url` varchar(250) NOT NULL DEFAULT '',
  `author` varchar(100) NOT NULL DEFAULT '',
  `views` int(10) unsigned NOT NULL DEFAULT 0,
  `status` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '状态，0不对外显示，1显示，99删除',
  `publish_time` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '发布时间',
  `created_time` int(10) unsigned NOT NULL DEFAULT 0,
  `updated_time` int(10) unsigned NOT NULL DEFAULT 0,
  `deleted_time` int(10) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `created_time` (`created_time`),
  KEY `updated_time` (`updated_time`),
  KEY `origin_url` (`origin_url`),
  KEY `title` (`title`),
  KEY `article_type` (`article_type`),
  KEY `views` (`views`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for fe_article_data
-- ----------------------------
DROP TABLE IF EXISTS `fe_article_data`;
CREATE TABLE `fe_article_data` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `content` longtext DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for fe_article_source
-- ----------------------------
DROP TABLE IF EXISTS `fe_article_source`;
CREATE TABLE `fe_article_source` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `url` varchar(250) NOT NULL DEFAULT '' COMMENT '列表源',
  `url_type` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '0,政策，1新闻',
  `error_times` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '出错次数',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `url` (`url`(191)) USING BTREE,
  KEY `error_times` (`error_times`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
