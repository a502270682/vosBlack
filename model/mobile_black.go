package model

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MobileBlack struct {
	NID       int       `json:"nid" gorm:"column:nid"`
	Mobile    string    `json:"mobile" gorm:"column:mobile"`         // 手机号码，后8位
	MobileAll string    `json:"mobile_all" gorm:"column:mobile_all"` // 完整手机号码
	MbLevel   int       `json:"mb_level" gorm:"column:mb_level"`     // 黑名单级别
	GwId      int       `json:"gw_id" gorm:"column:gw_id"`           // 调用网关id
	Hit       int       `json:"hit" gorm:"column:hit"`               // 累计请求次数
	EnID      int       `json:"en_id" gorm:"column:en_id"`           // 请求企业ID
	JoinDt    time.Time `json:"join_dt" gorm:"column:join_dt"`       // 添加时间
	GwRewrite string    `json:"gw_rewrite" gorm:"column:gw_rewrite"` // 第三方网关返回值
	Caller    string    `json:"caller" gorm:"column:caller"`
	Remark    string    `json:"remark" gorm:"column:remark"` // 号码备注
}

var mobileBlackApi *MobileBlackImpl

type MobileBlackImpl struct {
	DB *gorm.DB
}

type MobileBlackQueryCondition struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Prefix string `json:"prefix"`
}

type MobileBlackRepo interface {
	Insert(ctx context.Context, black *MobileBlack, prefix string) error
	Del(ctx context.Context, id int, prefix string) error
	GetListByCondition(ctx context.Context, query MobileBlackQueryCondition) ([]*MobileBlack, int64, error)
	GetOneByMobileAll(ctx context.Context, prefix string, mobileAll string) (*MobileBlack, error)
	GetOneByMobile(ctx context.Context, prefix string, mobile string) (*MobileBlack, error)
}

func InitAMobileBlackRepo(d *gorm.DB) {
	mobileBlackApi = &MobileBlackImpl{
		DB: d,
	}
}

func GetMobileBlackApi() MobileBlackRepo {
	return mobileBlackApi
}

func (m *MobileBlackImpl) Insert(ctx context.Context, black *MobileBlack, prefix string) error {
	tableName := fmt.Sprintf("mobile_black_%s", prefix)
	err := m.DB.WithContext(ctx).Table(tableName).Create(black).Error
	return err
}

func (m *MobileBlackImpl) Del(ctx context.Context, id int, prefix string) error {
	tableName := fmt.Sprintf("mobile_black_%s", prefix)
	return m.DB.WithContext(ctx).Table(tableName).Delete(&MobileBlack{}, id).Error
}

func (m *MobileBlackImpl) GetListByCondition(ctx context.Context, query MobileBlackQueryCondition) ([]*MobileBlack, int64, error) {
	tableName := fmt.Sprintf("mobile_black_%s", query.Prefix)
	var res []*MobileBlack
	var total int64
	db := m.DB
	// 增加where 条件
	err := db.WithContext(ctx).Table(tableName).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.WithContext(ctx).Table(tableName).Limit(query.Limit).Offset(query.Offset).Find(&res).Error
	if err != nil {
		return nil, 0, err
	}
	return res, total, nil

}

func (m *MobileBlackImpl) GetOneByMobileAll(ctx context.Context, prefix string, mobileAll string) (*MobileBlack, error) {
	tableName := fmt.Sprintf("mobile_black_%s", prefix)
	res := &MobileBlack{}
	err := m.DB.WithContext(ctx).Table(tableName).Where("mobile_all = ?", mobileAll).First(res).Error
	return res, err
}

func (m *MobileBlackImpl) GetOneByMobile(ctx context.Context, prefix string, mobile string) (*MobileBlack, error) {
	tableName := fmt.Sprintf("mobile_black_%s", prefix)
	var res *MobileBlack
	err := m.DB.WithContext(ctx).Table(tableName).Where("mobile = ?", mobile).First(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
