package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type EnterpriseApplyHourList struct {
	NID            int       `json:"nID" gorm:"column:nID"`
	EnID           int       `json:"en_id" gorm:"column:en_id"`
	DayReport      time.Time `json:"day_report" gorm:"column:day_report"`
	RepYear        int       `json:"rep_year" gorm:"column:rep_year"`
	RepMonth       int       `json:"rep_month" gorm:"column:rep_month"`
	RepDay         int       `json:"rep_day" gorm:"column:rep_day"`
	RepHour        int       `json:"rep_hour" gorm:"column:rep_hour"`
	MbRequestCount int64     `json:"mb_request_count" gorm:"column:mb_request_count"` // 黑名单请求次数
	MbHitCount     int64     `json:"mb_hit_count" gorm:"column:mb_hit_count"`         // 黑名单命中次数
	WnHitCount     int64     `json:"wn_hit_count" gorm:"column:wn_hit_count"`         // 白名单命中次数
	MpRequestCount int64     `json:"mp_request_count" gorm:"column:mp_request_count"` // 靓号请求次数
	MpHitCount     int64     `json:"mp_hit_count" gorm:"column:mp_hit_count"`         // 靓号命中次数
	GwRequestCount int64     `json:"gw_request_count" gorm:"column:gw_request_count"` // 外部接口请求次数
	GwHitCount     int64     `json:"gw_hit_count" gorm:"column:gw_hit_count"`         // 外部接口命中次数
	FqRequestCount int64     `json:"fq_request_count" gorm:"column:fq_request_count"` // 频次请求次数
	FqHitCount     int64     `json:"fq_hit_count" gorm:"column:fq_hit_count"`         // 频次命中次数
	JoinDt         time.Time `json:"join_dt" gorm:"column:join_dt"`                   // 添加时间
	Remark         string    `json:"remark" gorm:"column:remark"`                     // 备注
}

func (e EnterpriseApplyHourList) TableName() string {
	return "t_enterprise_applyhourlist"
}

var applyHourList *ApplyHourListImpl

type ApplyHourListImpl struct {
	DB *gorm.DB
}

type ApplyHourListRepo interface {
	GetsByEnID(ctx context.Context, id int) ([]*EnterpriseApplyHourList, error)
	GetLatestByEnID(ctx context.Context, id int) (*EnterpriseApplyHourList, error)
	Upsert(ctx context.Context, entity *EnterpriseApplyHourList) error
}

func InitApplyHourListRepo(d *gorm.DB) {
	applyHourList = &ApplyHourListImpl{
		DB: d,
	}
}

func GetInitApplyHourListImpl() ApplyHourListRepo {
	return applyHourList
}

func (a *ApplyHourListImpl) GetsByEnID(ctx context.Context, id int) ([]*EnterpriseApplyHourList, error) {
	var res []*EnterpriseApplyHourList
	err := a.DB.WithContext(ctx).Where("en_id = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApplyHourListImpl) GetLatestByEnID(ctx context.Context, id int) (*EnterpriseApplyHourList, error) {
	var res []*EnterpriseApplyHourList
	err := a.DB.WithContext(ctx).Where("en_id = ?", id).Order("join_at desc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func (a *ApplyHourListImpl) Upsert(ctx context.Context, entity *EnterpriseApplyHourList) error {
	return a.DB.WithContext(ctx).Save(entity).Error
}
