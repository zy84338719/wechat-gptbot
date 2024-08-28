package data

import (
	"context"
	"github.com/yi-nology/common/utils/xnow"
	"github.com/yi-nology/sdk/conf"
	"wechat-gptbot/data/proto"
)

var GroupSignSingleton = new(GroupSign)

type GroupSign struct {
}

func (g GroupSign) Create(data proto.GroupSign) error {
	return conf.MysqlClient.Create(&data).Error
}

func (g GroupSign) GetByGroupUsernameAndNowDay(groupUsername, username string) (*proto.GroupSign, error) {
	var data *proto.GroupSign
	// 事件大于今天

	rx := conf.MysqlClient.Where("group_username = ? and username = ? and sign_time > ?", groupUsername, username, xnow.BeginningOfDay()).Find(&data)
	if rx.Error != nil {
		return nil, rx.Error
	}
	if rx.RowsAffected == 0 {
		return nil, nil
	}
	return data, nil
}

func (g GroupSign) GetByGroupUsernameAndDayCount(groupUsername string) (int64, error) {
	var count int64
	err := conf.MysqlClient.Model(proto.GroupSign{}).Where("group_username = ? and sign_time > ?", groupUsername, xnow.BeginningOfDay()).Count(&count).Error
	return count, err
}

// 本月签到次数
func (g GroupSign) GetByGroupUsernameAndMonthCount(username string) (int64, error) {
	var count int64
	err := conf.MysqlClient.Model(proto.GroupSign{}).Where("username = ? and sign_time > ?", username, xnow.BeginningOfMonth()).Count(&count).Error
	return count, err
}

func (g GroupSign) Count(ctx context.Context) int {
	var count int64
	err := conf.MysqlClient.WithContext(ctx).Model(&proto.GroupSign{}).Count(&count).Error
	if err != nil {
		return 0
	}
	return int(count)
}

// list
func (g GroupSign) List(ctx context.Context, page, pageSize int) ([]proto.GroupSign, error) {
	var data []proto.GroupSign
	err := conf.MysqlClient.WithContext(ctx).Offset(page * pageSize).Limit(pageSize).Find(&data).Error
	return data, err
}

// update
func (g GroupSign) Update(ctx context.Context, data proto.GroupSign) error {
	return conf.MysqlClient.WithContext(ctx).Save(&data).Error
}
