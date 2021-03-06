package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// EnterpriseCalltimelist 呼叫时段
type EnterpriseCalltimelist struct {
	NID         int       `json:"nid,omitempty" gorm:"column:nid"`
	EnID        int       `json:"en_id,omitempty" gorm:"column:en_id"`   // 请求企业ID
	BlackID     int       `json:"black_id" gorm:"column:black_id"`       // 黑名单规则id  t_enterprise_blacklist的nid
	TimeName    string    `json:"timename" gorm:"column:timename"`       // 时间段名称
	BeginHour   int       `json:"beginhour" gorm:"column:beginhour"`     // 开始小时
	BeginMinute int       `json:"beginninute" gorm:"column:beginninute"` //开始分钟
	EndHour     int       `json:"endhour" gorm:"column:endhour"`         //结束时间
	Edminute    int       `json:"edminute" gorm:"column:edminute"`       // 结束分钟
	JoinDt      time.Time `json:"join_dt" gorm:"column:join_dt"`         // 添加时间
	Remark      string    `json:"remark" gorm:"column:remark"`           // 备注
}

func (EnterpriseCalltimelist) TableName() string {
	return "t_enterprise_calltimelist"
}

var enterpriseCalltimelistImpl *EnterpriseCalltimelistImpl

type EnterpriseCalltimelistImpl struct {
	DB *gorm.DB
}

func InitEnterpriseCalltimelistRepo(db *gorm.DB) {
	enterpriseCalltimelistImpl = &EnterpriseCalltimelistImpl{
		DB: db,
	}
}

type EnterpriseCalltimelistRepo interface {
	GetByEnID(ctx context.Context, enID int, blackID int) ([]*EnterpriseCalltimelist, error)
	GetAll(ctx context.Context) ([]*EnterpriseCalltimelist, error)
}

func GetEnterpriseCalltimelistImpl() EnterpriseCalltimelistRepo {
	return enterpriseCalltimelistImpl
}

func (e *EnterpriseCalltimelistImpl) GetByEnID(ctx context.Context, enID int, blackID int) ([]*EnterpriseCalltimelist, error) {
	res := []*EnterpriseCalltimelist{}
	err := e.DB.WithContext(ctx).
		Where("en_id = ?", enID).
		Where("black_id = ?", blackID).
		Find(&res).Error
	return res, err
}

func (e *EnterpriseCalltimelistImpl) GetAll(ctx context.Context) ([]*EnterpriseCalltimelist, error) {
	res := make([]*EnterpriseCalltimelist, 0)
	err := e.DB.WithContext(ctx).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
