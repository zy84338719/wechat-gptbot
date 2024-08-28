package data

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"gorm.io/gorm"
	"wechat-gptbot/data/proto"
)

var GroupMessageSingleton = GroupMessage{}

type GroupMessage struct {
}

func (GroupMessage) Create(ctx context.Context, data proto.GroupMessage) error {
	return conf.MysqlClient.WithContext(ctx).Create(&data).Error
}

func (GroupMessage) Update(ctx context.Context, data proto.GroupMessage) error {
	tx := conf.MysqlClient.WithContext(ctx).Model(proto.GroupMessage{}).Where("id = ?", data.ID).Updates(&data)
	return tx.Error
}

func (GroupMessage) Get(ctx context.Context, id int64) (*proto.GroupMessage, error) {
	var data *proto.GroupMessage
	err := conf.MysqlClient.WithContext(ctx).First(&data, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return data, err
}

func (GroupMessage) Delete(ctx context.Context, id int64) error {
	return conf.MysqlClient.Delete(&proto.GroupMessage{}, id).Error
}

// count
func (GroupMessage) Count(ctx context.Context) (int64, error) {
	var count int64
	err := conf.MysqlClient.WithContext(ctx).Model(&proto.GroupMessage{}).Count(&count).Error
	return count, err
}

func (GroupMessage) List(ctx context.Context, page, pageSize int) ([]proto.GroupMessage, error) {
	var data []proto.GroupMessage
	tx := conf.MysqlClient.Offset(page * pageSize).Limit(pageSize).Find(&data)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return data, nil
}

func (GroupMessage) GetByGroupUsername(ctx context.Context, username string) ([]proto.GroupMessage, error) {
	var data []proto.GroupMessage
	tx := conf.MysqlClient.WithContext(ctx).Where("group_username = ?", username).Find(&data)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return data, nil
}

func (GroupMessage) GetByUsername(ctx context.Context, username string) ([]proto.GroupMessage, error) {
	var data []proto.GroupMessage
	err := conf.MysqlClient.WithContext(ctx).Where("send_username = ?", username).Find(&data).Error
	return data, err
}
