/*
 Navicat Premium Data Transfer

 Source Server         : docker_mysql
 Source Server Type    : MySQL
 Source Server Version : 80020
 Source Host           : localhost:3309
 Source Schema         : evn

 Target Server Type    : MySQL
 Target Server Version : 80020
 File Encoding         : 65001

 Date: 23/12/2023 15:32:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for lv_article_classification
-- ----------------------------
DROP TABLE IF EXISTS `lv_article_classification`;
CREATE TABLE `lv_article_classification`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '分类id',
  `a_id` int NOT NULL DEFAULT 0 COMMENT '上级id 0表示顶级',
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类名称',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文章分类' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_article_contribution
-- ----------------------------
DROP TABLE IF EXISTS `lv_article_contribution`;
CREATE TABLE `lv_article_contribution`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户ID',
  `classification_id` int NULL DEFAULT NULL COMMENT '分类ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文章标题',
  `cover` json NOT NULL COMMENT '文章封面',
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标签定义，分割	',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '富文本存储值',
  `content_storage_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '富文本存储类型',
  `is_comments` int NOT NULL DEFAULT 1 COMMENT '是否可以评论0否1是',
  `heat` int NULL DEFAULT 0 COMMENT '文章热度',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_article_contribution_comments
-- ----------------------------
DROP TABLE IF EXISTS `lv_article_contribution_comments`;
CREATE TABLE `lv_article_contribution_comments`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `article_id` int NOT NULL COMMENT '文章id',
  `context` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `comment_id` int NULL DEFAULT NULL COMMENT '评论上级id',
  `comment_user_id` int NOT NULL DEFAULT 0 COMMENT '评论上级的用户id (0为空)',
  `comment_first_id` int NOT NULL DEFAULT 0 COMMENT '一级评论(顶层评论)id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 103 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_article_contribution_like
-- ----------------------------
DROP TABLE IF EXISTS `lv_article_contribution_like`;
CREATE TABLE `lv_article_contribution_like`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `article_id` int NOT NULL COMMENT '文章id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_home_rotograph
-- ----------------------------
DROP TABLE IF EXISTS `lv_home_rotograph`;
CREATE TABLE `lv_home_rotograph`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `cover` json NOT NULL COMMENT '封面',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '轮播图标题',
  `color` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '图片主题颜色',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '跳转类型 : video跳转视频 article跳转文章 live跳转直播',
  `to_id` int NOT NULL COMMENT '跳转id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '首页轮播图' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_live_info
-- ----------------------------
DROP TABLE IF EXISTS `lv_live_info`;
CREATE TABLE `lv_live_info`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `uid` int UNSIGNED NULL DEFAULT NULL COMMENT '绑定的用户id',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '直播间标题',
  `img` json NULL COMMENT '直播间封面',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uid`(`uid` ASC) USING BTREE,
  INDEX `uid_2`(`uid` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_upload_method
-- ----------------------------
DROP TABLE IF EXISTS `lv_upload_method`;
CREATE TABLE `lv_upload_method`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `interface` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求接口名',
  `method` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '上传方法',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '规定图片存储路径',
  `quality` decimal(2, 1) NOT NULL DEFAULT 1.0 COMMENT '图片质量',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users
-- ----------------------------
DROP TABLE IF EXISTS `lv_users`;
CREATE TABLE `lv_users`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `openid` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '微信小程序openid',
  `email` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '邮箱',
  `username` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '密码',
  `photo` json NULL COMMENT '头像',
  `gender` smallint NULL DEFAULT NULL COMMENT '性别 0男 1女 2保密',
  `birth_date` datetime(3) NULL DEFAULT NULL COMMENT '出生时间',
  `is_visible` smallint NULL DEFAULT NULL COMMENT '是否可见 0 false 1 true',
  `signature` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '个性签名',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '登陆时间',
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `Id`(`id` ASC, `openid` ASC) USING BTREE,
  INDEX `openid`(`openid` ASC) USING BTREE,
  INDEX `idx_lv_users_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 48 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '用户信息' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_attention
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_attention`;
CREATE TABLE `lv_users_attention`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户的id',
  `attention_id` int NOT NULL COMMENT '关注id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 55 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户关注列表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_chat_list
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_chat_list`;
CREATE TABLE `lv_users_chat_list`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `tid` int NOT NULL COMMENT '聊天对象id',
  `unread` int NOT NULL DEFAULT 0 COMMENT '未读消息数量',
  `last_message` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '最后一条消息',
  `created_at` datetime(3) NOT NULL COMMENT '初次聊天时间',
  `last_at` datetime(3) NULL DEFAULT NULL COMMENT '最后聊天时间',
  `updated_at` datetime(3) NOT NULL COMMENT '最后进入时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 69 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_chat_msg
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_chat_msg`;
CREATE TABLE `lv_users_chat_msg`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '消息id',
  `uid` int NOT NULL COMMENT '发送者id',
  `tid` int NOT NULL COMMENT '接收者id',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'text' COMMENT '消息类型',
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息内容',
  `created_at` datetime(3) NOT NULL COMMENT '发送时间',
  `updated_at` datetime(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 46 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_collect
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_collect`;
CREATE TABLE `lv_users_collect`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `favorites_id` int NOT NULL COMMENT '收藏夹id',
  `video_id` int NOT NULL COMMENT '视频iid',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 78 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户收藏表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_favorites
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_favorites`;
CREATE TABLE `lv_users_favorites`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '收藏夹创建者id',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '收藏夹标题',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '收藏夹内容介绍',
  `cover` json NULL COMMENT '收藏夹封面',
  `max` int NOT NULL DEFAULT 1000 COMMENT '最大存储量',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '收藏时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户收藏夹表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_notice
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_notice`;
CREATE TABLE `lv_users_notice`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '通知id',
  `uid` int NOT NULL COMMENT '通知用户id',
  `cid` int NOT NULL COMMENT '操作者id',
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '视频评论videoComment  视频点赞 videoLike  文章评论articleComment  文章点赞 articleLike',
  `to_id` int NOT NULL COMMENT '对于类型id',
  `is_read` int NOT NULL DEFAULT 0 COMMENT '消息是否已读 0未读 1已读',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知内容',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  UNIQUE INDEX `id`(`id` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 98 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户通知表' ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_users_record
-- ----------------------------
DROP TABLE IF EXISTS `lv_users_record`;
CREATE TABLE `lv_users_record`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型 : video跳转视频 article跳转文章 live跳转直播',
  `to_id` int NULL DEFAULT NULL COMMENT '关联id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 172 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_video_contribution
-- ----------------------------
DROP TABLE IF EXISTS `lv_video_contribution`;
CREATE TABLE `lv_video_contribution`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '视频标题',
  `video` json NULL COMMENT '视频(默认1080p)',
  `video_720p` json NULL COMMENT '720p视频',
  `video_480p` json NULL COMMENT '480p视频',
  `video_360p` json NULL COMMENT '360p视频',
  `media_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '媒资ID 阿里云oss 确定视频需要',
  `cover` json NOT NULL COMMENT '视频封面',
  `video_duration` int NOT NULL COMMENT '视频时长',
  `reprinted` int NOT NULL DEFAULT 0 COMMENT '是否转载0不是1是',
  `label` varchar(400) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '标签定义，分割',
  `introduce` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '视频介绍',
  `heat` int NOT NULL DEFAULT 0 COMMENT '视频热度',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3153 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_video_contribution_barrage
-- ----------------------------
DROP TABLE IF EXISTS `lv_video_contribution_barrage`;
CREATE TABLE `lv_video_contribution_barrage`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `video_id` int NOT NULL COMMENT '视频id',
  `time` float(10, 6) NOT NULL DEFAULT 0.000000 COMMENT '出现时间',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `text` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '弹幕内容',
  `type` int NULL DEFAULT 0 COMMENT '弹幕位置',
  `color` int NOT NULL COMMENT '弹幕颜色（十进制）',
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 68 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_video_contribution_comments
-- ----------------------------
DROP TABLE IF EXISTS `lv_video_contribution_comments`;
CREATE TABLE `lv_video_contribution_comments`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `video_id` int NOT NULL COMMENT '视频id',
  `context` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '评论内容',
  `comment_id` int NULL DEFAULT NULL COMMENT '评论上级id',
  `comment_user_id` int NOT NULL DEFAULT 0 COMMENT '评论上级的用户id (0为空)',
  `comment_first_id` int NOT NULL DEFAULT 0 COMMENT '一级评论(顶层评论)id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 66 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for lv_video_contribution_like
-- ----------------------------
DROP TABLE IF EXISTS `lv_video_contribution_like`;
CREATE TABLE `lv_video_contribution_like`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL COMMENT '用户id',
  `video_id` int NOT NULL COMMENT '视频id',
  `created_at` datetime(3) NOT NULL COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 83 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
