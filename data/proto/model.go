package proto

import (
	"github.com/eatmoreapple/openwechat"
	"gorm.io/gorm"
	"time"
)

// CREATE TABLE `group_userinfo` (
//  `id` int(11) NOT NULL,
//  `bot_id` varchar(255) DEFAULT NULL,
//  `bot_name` varchar(255) DEFAULT NULL,
//  `group_username` varchar(255) DEFAULT NULL,
//  `username` varchar(255) DEFAULT NULL,
//  `display_name` varchar(255) DEFAULT NULL,
//  `nikename` varchar(255) DEFAULT NULL,
//  `createAt` datetime DEFAULT NULL,
//  `updateAt` datetime DEFAULT NULL,
//  `deleteAt` datetime DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

//CREATE TABLE `group_message` (
//  `id` int(11) NOT NULL COMMENT '主键',
//  `bot_id` varchar(255) DEFAULT NULL,
//  `group_username` varchar(255) DEFAULT NULL COMMENT '发送群ID',
//  `send_username` varchar(255) DEFAULT NULL COMMENT '发送人',
//  `content` varchar(2048) DEFAULT NULL COMMENT '发送内容',
//  `sub_msg_type` int(11) DEFAULT NULL COMMENT '消息类型',
//  `msg_create_time` bigint(20) DEFAULT NULL COMMENT '消息时间戳',
//  `createdAt` datetime DEFAULT NULL,
//  `updatedAt` datetime DEFAULT NULL,
//  `deletedAt` datetime DEFAULT NULL,
//  PRIMARY KEY (`id`),
//  KEY `idx_deletedAt` (`deletedAt`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

//CREATE TABLE `group` (
//  `ID` int(11) NOT NULL,
//  `bot_id` int(11) DEFAULT NULL,
//  `bot_name` varchar(255) DEFAULT NULL,
//  `group_username` varchar(255) DEFAULT NULL,
//  `group_nikename` varchar(255) DEFAULT NULL,
//  `CreateAt` datetime DEFAULT NULL,
//  `UpdateAt` datetime DEFAULT NULL,
//  `DeleteAt` datetime DEFAULT NULL,
//  PRIMARY KEY (`ID`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

// CREATE TABLE `group_sign` (
//  `id` bigint(20) NOT NULL AUTO_INCREMENT,
//  `bot_id` bigint(20) DEFAULT NULL,
//  `group_username` varchar(255) DEFAULT NULL,
//  `username` varchar(255) DEFAULT NULL,
//  `sigin_time` datetime DEFAULT NULL,
//  `create_at` datetime DEFAULT NULL,
//  `update_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
//  `delete_at` datetime DEFAULT NULL,
//  PRIMARY KEY (`id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

// 将上面的表结构转化为Go的结构体
type GroupUserinfo struct {
	gorm.Model
	BotID         int64  `gorm:"column:bot_id"`
	GroupUsername string `gorm:"column:group_username"`
	Username      string `gorm:"column:username"`
	DisplayName   string `gorm:"column:display_name"`
	Nickname      string `gorm:"column:nickname"`
	IsQuit        bool   `gorm:"column:is_quit"`
}

func (u GroupUserinfo) TableName() string {
	return "group_userinfo"
}

type GroupMessage struct {
	gorm.Model
	BotID         int64                  `gorm:"column:bot_id"`
	GroupUsername string                 `gorm:"column:group_username"`
	Username      string                 `gorm:"column:username"`
	MsgId         string                 `gorm:"column:msg_id"`
	Content       string                 `gorm:"column:content"`
	MsgType       openwechat.MessageType `gorm:"column:msg_type"`
	MsgCreateTime int64                  `gorm:"column:msg_create_time"`
}

func (m GroupMessage) TableName() string {
	return "group_message"
}

type Group struct {
	gorm.Model
	Version       int64  `gorm:"column:version"`
	BotID         int64  `gorm:"column:bot_id"`
	BotName       string `gorm:"column:bot_name"`
	GroupUsername string `gorm:"column:group_username"`
	GroupNickname string `gorm:"column:group_nickname"`
}

func (g Group) TableName() string {
	return "group"
}

type GroupSign struct {
	gorm.Model
	BotID         int64     `gorm:"column:bot_id"`
	GroupUsername string    `gorm:"column:group_username"`
	Username      string    `gorm:"column:username"`
	SignTime      time.Time `gorm:"column:sign_time"`
}

func (g *GroupSign) TableName() string {
	return "group_sign"
}

type Groups []Group

func (g Groups) Len() int {
	return len(g)
}

func (g Groups) Less(i, j int) bool {
	return g[i].CreatedAt.Sub(g[j].CreatedAt) > 0
}

func (g Groups) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g Groups) GetSlice() []Group {
	return g

}
