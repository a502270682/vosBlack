package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type GatewayApplyHourList struct {
	NID            int       `json:"nID" gorm:"column:nID"`
	GwID           int       `json:"gw_id" gorm:"column:gw_id"`
	DayReport      time.Time `json:"day_report" gorm:"column:day_report"`
	RepYear        int       `json:"rep_year" gorm:"column:rep_year"`
	RepMonth       int       `json:"rep_month" gorm:"column:rep_month"`
	RepDay         int       `json:"rep_day" gorm:"column:rep_day"`
	RepHour        int       `json:"rep_hour" gorm:"column:rep_hour"`
	MbRequestCount int64     `json:"mb_request_count" gorm:"column:mb_request_count"` // 黑名单请求次数
	MbHitCount     int64     `json:"mb_hit_count" gorm:"column:mb_hit_count"`         // 黑名单命中次数
	JoinDt         time.Time `json:"join_dt" gorm:"column:join_dt"`                   // 添加时间
	Remark         string    `json:"remark" gorm:"column:remark"`                     // 备注
}

func (e GatewayApplyHourList) TableName() string {
	return "sys_gateway_applyhourlist"
}

var gatewayApplyHourList *GateWayApplyHourListImpl

type GateWayApplyHourListImpl struct {
	DB *gorm.DB
}

type GateWayApplyHourListRepo interface {
	GetDBForTransaction() *gorm.DB
	GetsByEnID(ctx context.Context, id int) ([]*GatewayApplyHourList, error)
	GetLatestByEnID(ctx context.Context, id int) (*GatewayApplyHourList, error)
	Upsert(ctx context.Context, entity *GatewayApplyHourList) error
}

func InitGateWayApplyHourListRepo(d *gorm.DB) {
	gatewayApplyHourList = &GateWayApplyHourListImpl{
		DB: d,
	}
}

func GetGateWayApplyHourListImpl() GateWayApplyHourListRepo {
	return gatewayApplyHourList
}

func (a *GateWayApplyHourListImpl) GetDBForTransaction() *gorm.DB {
	return a.DB
}

func (a *GateWayApplyHourListImpl) GetsByEnID(ctx context.Context, id int) ([]*GatewayApplyHourList, error) {
	var res []*GatewayApplyHourList
	err := a.DB.WithContext(ctx).Where("en_id = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *GateWayApplyHourListImpl) GetLatestByEnID(ctx context.Context, id int) (*GatewayApplyHourList, error) {
	var res []*GatewayApplyHourList
	err := a.DB.WithContext(ctx).Where("en_id = ?", id).Order("join_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func (a *GateWayApplyHourListImpl) Upsert(ctx context.Context, entity *GatewayApplyHourList) error {
	return a.DB.WithContext(ctx).Save(entity).Error
}
