package data

import (
	"context"
	"github.com/yi-nology/sdk/conf"
	"gorm.io/gorm"
	"wechat-gptbot/data/proto"
)

var GroupUserInfoSingleton = GroupUserInfo{}

type GroupUserInfo struct {
}

func (GroupUserInfo) Create(ctx context.Context, data proto.GroupUserinfo) error {
	return conf.MysqlClient.WithContext(ctx).Create(&data).Error
}

func (GroupUserInfo) BatchCreate(ctx context.Context, data []proto.GroupUserinfo) error {
	return conf.MysqlClient.WithContext(ctx).CreateInBatches(&data, len(data)).Error
}

func (GroupUserInfo) Update(ctx context.Context, data proto.GroupUserinfo) error {
	return conf.MysqlClient.Save(&data).Error
}

func (GroupUserInfo) Get(ctx context.Context, id int64) (*proto.GroupUserinfo, error) {
	var data *proto.GroupUserinfo
	err := conf.MysqlClient.WithContext(ctx).First(&data, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return data, nil
}

func (GroupUserInfo) Delete(ctx context.Context, id int64) error {
	return conf.MysqlClient.Delete(&proto.GroupUserinfo{}, id).Error
}

func (GroupUserInfo) DeleteByGroupUsername(ctx context.Context, groupUsername, username string) error {
	return conf.MysqlClient.Where("group_username = ? and username = ?", groupUsername, username).Delete(&proto.GroupUserinfo{}).Error
}

func (GroupUserInfo) List(ctx context.Context, page, pageSize int) ([]proto.GroupUserinfo, error) {
	var data []proto.GroupUserinfo
	err := conf.MysqlClient.WithContext(ctx).Offset(page * pageSize).Limit(pageSize).Find(&data).Error
	return data, err
}

func (GroupUserInfo) GetByGroupUsername(ctx context.Context, groupUsername string) ([]proto.GroupUserinfo, error) {
	var data []proto.GroupUserinfo
	rx := conf.MysqlClient.Where("group_username = ? ", groupUsername).Find(&data)
	if rx.Error != nil {
		return nil, rx.Error
	}
	if rx.RowsAffected == 0 {
		return nil, nil
	}
	return data, nil
}

func (GroupUserInfo) GetByBotId(ctx context.Context, botId string) ([]proto.GroupUserinfo, error) {
	var data []proto.GroupUserinfo
	err := conf.MysqlClient.Where("bot_id = ?", botId).Find(&data).Error
	return data, err
}
