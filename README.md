以下是优化后的文档，旨在使其更易于理解、使用和维护。

---

# 本项目参考了 code-innovator-zyx/wechat-gptbot 项目

欢迎来到 **星期五v0.1 微信 GPT 机器人** 项目！这个项目可以帮助你在微信上创建一个智能机器人，用于与孩子或宠物互动。你可以利用它发送消息、查询天气、获取每日新闻等。

> 项目地址: [https://github.com/code-innovator-zyx/wechat-gptbot](https://github.com/code-innovator-zyx/wechat-gptbot)


    

1运行服务,你可以选择两种运行方式：

   ```shell
   # 1本地运行
   make local
   
   # 2 docker运行
   make docker
   ```

首次执行时，机器人会提示扫码登录微信。

- 登录完成后，系统会生成一个 `token.json` 文件，用于保存当前的微信登录状态，避免每次运行都需要

本服务需要依赖mysql redis apollo 
配置如下
```shell
CREATE TABLE `group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bot_id` bigint(20) DEFAULT NULL,
  `bot_name` varchar(255) DEFAULT NULL,
  `group_username` varchar(255) DEFAULT NULL,
  `group_nickname` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4

CREATE TABLE `group_message` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `bot_id` bigint(20) DEFAULT NULL,
  `group_username` varchar(255) DEFAULT NULL COMMENT '发送群ID',
  `username` varchar(255) DEFAULT NULL COMMENT '发送人',
  `msg_id` varchar(128) DEFAULT NULL,
  `content` varchar(2048) DEFAULT NULL COMMENT '发送内容',
  `msg_type` int(11) DEFAULT NULL COMMENT '消息类型',
  `msg_create_time` bigint(20) DEFAULT NULL COMMENT '消息时间戳',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_deletedAt` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=utf8mb4

CREATE TABLE `group_sign` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bot_id` bigint(20) DEFAULT NULL,
  `group_username` varchar(255) DEFAULT NULL,
  `username` varchar(255) DEFAULT NULL,
  `sign_time` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_index` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4

CREATE TABLE `group_userinfo` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bot_id` varchar(255) DEFAULT NULL,
  `group_username` varchar(255) DEFAULT NULL,
  `username` varchar(255) DEFAULT NULL,
  `display_name` varchar(255) DEFAULT NULL,
  `nickname` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uqi_b_g_u` (`bot_id`,`group_username`,`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4635 DEFAULT CHARSET=utf8mb4

```

仅需要配置好数据库和redis的连接信息即可
需要安装配置 携程项目apollo
请按照 https://github.com/apolloconfig/apollo-quick-start 快速开启配置 配置内容详情请见
文件夹内 wechat-bot-application.json wechat-bot-tempConfig.json.json  
如果需要自然语言回答处理能力，请配置好self_gpt的相关信息
需要安装 飞致云的 maxkb 项目 以及羊驼大模型 我自己用的mac m1 pro 跑的32b的通义千问模型

如需crawler_domain film_domain shorturl_domain moyu_domain 请配置好相关信息 请联系我获取相关信息
由于南北向并未鉴权开发 故先不提供相关信息 如需要请自行联系我
zhangyi@murphyyi.com


```shell
apollo:
  enable: true
  appId: wechat-bot
  cluster: DEV
  namespaceName: application,tempConfig.json
  endpoint: http://xxxx
  secret: xxx
  dynamic: true
  interval: 5
  isBackupConfig: true
  mustStart: false

redis:
  enable: true
  addr: xxxxx:6379
  password: xxxxxxx
  db: 1

mysql:
  enable: true
  prefix: ""
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: wechat_bot_dev
  username: xxxx
  password: xxxxxxx
  path: vpn.enjoye.top
  engine: ""
  log-mode: error
  max-idle-conns: 10
  max-open-conns: 100
  singular: false
  log-zap: false

self_gpt:
  base_url: http://xxxxxx:11435/api
  keyword_authorization: xxxxxxx
  gpt_authorization: xxxxx
crawler_domain: http://xxxx:18891
film_domain: https://xxxxx/film
shorturl_domain: https://xxxxx/shorturl
moyu_domain: https://xxxxx/moyu
```

