package model

import (
	"context"
	"gorm.io/gorm"
)

type EnterpriseUserlist struct {
	NID         int    `json:"nid" gorm:"column:nid"`
	EnID        int    `json:"en_id" gorm:"column:en_id"`
	UserID      string `json:"userid" gorm:"column:userid"`           // 用户账户
	UserPass    string `json:"userpass" gorm:"column:userpass"`       // 用户密码
	LoginType   int    `json:"logintype" gorm:"column:logintype"`     //
	MobilePhone string `json:"mobilephone" gorm:"column:mobilephone"` // 手机号
	UserEmail   string `json:"useremail" gorm:"column:useremail"`     // 邮箱
	IStatus     string `json:"i_status" gorm:"column:i_status"`       //状态，1启用，0停用，9暂停，-1删除
	JoinDt      string `json:"join_dt" gorm:"column:join_dt"`         //
	Remark      string `json:"remark" gorm:"column:remark"`
}

func (EnterpriseUserlist) TableName() string {
	return "t_enterprise_userlist"
}

var enterpriseUserlist *EnterpriseUserlistImpl

type EnterpriseUserlistImpl struct {
	DB *gorm.DB
}

type EnterpriseUserlistRepo interface {
	GetByUserID(ctx context.Context, userID string) (*EnterpriseUserlist, error)
}

func InitEnterpriseUserlistRepo(d *gorm.DB) {
	enterpriseUserlist = &EnterpriseUserlistImpl{
		DB: d,
	}
}

func GetEnterpriseUserlistImpl() EnterpriseUserlistRepo {
	return enterpriseUserlist
}

func (e *EnterpriseUserlistImpl) GetByUserID(ctx context.Context, userID string) (*EnterpriseUserlist, error) {
	res := &EnterpriseUserlist{}
	err := e.DB.WithContext(ctx).Where("userid = ?", userID).First(res).Error
	return res, err
}
