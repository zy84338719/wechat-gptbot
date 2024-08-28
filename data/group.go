package data

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"gorm.io/gorm"
	"wechat-gptbot/data/proto"
)

var GroupSingleton = Group{}

type Group struct {
}

func (Group) Create(ctx context.Context, data proto.Group) error {
	return conf.MysqlClient.WithContext(ctx).Create(&data).Error
}

func (Group) Update(ctx context.Context, data proto.Group) error {
	return conf.MysqlClient.WithContext(ctx).Save(&data).Error
}

func (Group) GetByUsername(ctx context.Context, botId int64, username string) (*proto.Group, error) {
	var g = proto.Group{}
	tx := conf.MysqlClient.WithContext(ctx).Where("bot_id = ? and group_username = ?", botId, username).Find(&proto.Group{}).Find(&g)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return &g, nil
}

func (Group) Get(ctx context.Context, id int64) (*proto.Group, error) {
	var data *proto.Group
	err := conf.MysqlClient.WithContext(ctx).First(&data, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return data, err
}

func (Group) Delete(ctx context.Context, id int64) error {
	return conf.MysqlClient.WithContext(ctx).Delete(&proto.Group{}, id).Error
}

func (Group) List(ctx context.Context, page, pageSize int) ([]proto.Group, error) {
	var data []proto.Group
	tx := conf.MysqlClient.WithContext(ctx).Offset(page * pageSize).Limit(pageSize).Find(&data)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return data, nil
}

func (Group) GetByBotId(ctx context.Context, botId string) ([]proto.Group, error) {
	var data []proto.Group
	tx := conf.MysqlClient.WithContext(ctx).Where("bot_id = ?", botId).Find(&data)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return data, nil
}
