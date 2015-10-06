/*
 Navicat Premium Data Transfer

 Source Server         : foo
 Source Server Type    : SQLite
 Source Server Version : 3008004
 Source Database       : main

 Target Server Type    : SQLite
 Target Server Version : 3008004
 File Encoding         : utf-8

 Date: 03/24/2015 11:30:22 AM
*/

PRAGMA foreign_keys = false;

-- ----------------------------
--  Table structure for good_friends
-- ----------------------------
DROP TABLE IF EXISTS "good_friends";
CREATE TABLE "good_friends" (
	 "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "first_name" TEXT NOT NULL,
	 "last_name" TEXT NOT NULL
);

-- ----------------------------
--  Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
	 "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	 "name" TEXT,
	 "alive" integer(1,0),
	 "created_at" timestamp NOT NULL,
	 "updated_at" timestamp NOT NULL,
	 "birth_date" timestamp,
	 "bio" TEXT,
	 "price" FLOAT
);

PRAGMA foreign_keys = true;
