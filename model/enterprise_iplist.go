package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// 企业IP表
type EnterpriseIplist struct {
	NID       int       `json:"nid" gorm:"column:nid"`
	EnID      int       `json:"en_id" gorm:"column:en_id"`
	IpType    int       `json:"ip_type" gorm:"column:ip_type"`
	IpAll     string    `json:"ip_all" gorm:"ip_all"`
	Ip_1      int       `json:"ip_1" gorm:"column:ip_1"`
	Ip_2      int       `json:"ip_2" gorm:"column:ip_2"`
	Ip_3      int       `json:"ip_3" gorm:"column:ip_3"`
	Ip_4      int       `json:"ip_4" gorm:"column:ip_4"`
	Ip_5      int       `json:"ip_5" gorm:"column:ip_5"`
	Ip_6      int       `json:"ip_6" gorm:"column:ip_6"`
	Ip_7      int       `json:"ip_7" gorm:"column:ip_7"`
	Ip_8      int       `json:"ip_8" gorm:"column:ip_8"`
	Inputtype int       `json:"inputtype" gorm:"column:inputtype"`
	IStatus   IStatus   `json:"i_status" gorm:"column:i_status"`
	JoinDt    time.Time `json:"join_dt" gorm:"column:join_dt"`
	Remark    string    `json:"remark" gorm:"remark"`
}

func (EnterpriseIplist) TableName() string {
	return "t_enterprise_iplist"
}

var enterpriseIplistImpl *EnterpriseIplistImpl

type EnterpriseIplistImpl struct {
	DB *gorm.DB
}

func InitEnterpriseIplistRepo(d *gorm.DB) {
	enterpriseIplistImpl = &EnterpriseIplistImpl{
		DB: d,
	}
}

type EnterpriseIplistRepo interface {
	GetOneByIP(ctx context.Context, ip string) (*EnterpriseIplist, error)
}

func GetEnterpriseIplistImpl() EnterpriseIplistRepo {
	return enterpriseIplistImpl
}

func (e *EnterpriseIplistImpl) GetOneByIP(ctx context.Context, ip string) (*EnterpriseIplist, error) {
	res := &EnterpriseIplist{}
	err := e.DB.WithContext(ctx).Where("ip_all = ?", ip).First(res).Error
	return res, err
}
