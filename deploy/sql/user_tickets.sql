/*
 Navicat Premium Data Transfer

 Source Server         : test2
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 192.168.75.132:3306
 Source Schema         : flight

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 29/03/2022 18:38:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_tickets
-- ----------------------------
DROP TABLE IF EXISTS `user_tickets`;
CREATE TABLE `user_tickets`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `del_state` tinyint(1) NOT NULL COMMENT '是否已经删除',
  `version` bigint NOT NULL COMMENT '版本号',
  `auth_key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户平台唯一id',
  `ticket_id` bigint NOT NULL COMMENT '票id',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_key_ticket`(`auth_key`, `ticket_id`) USING BTREE,
  UNIQUE INDEX `idx_auth_key`(`auth_key`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;
