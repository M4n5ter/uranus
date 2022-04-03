/*
 Navicat Premium Data Transfer

 Source Server         : localMysql8
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 192.168.141.4:3306
 Source Schema         : flight

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 03/04/2022 22:19:53
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for wallet
-- ----------------------------
DROP TABLE IF EXISTS `wallet`;
CREATE TABLE `wallet`  (
  `id` bigint(0) NOT NULL,
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0),
  `delete_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0),
  `del_state` tinyint(0) NOT NULL,
  `version` bigint(0) NOT NULL,
  `user_id` bigint(0) NOT NULL COMMENT '用户id\r\n',
  `money` bigint(0) NOT NULL COMMENT '余额(分)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
