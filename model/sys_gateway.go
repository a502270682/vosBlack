package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SysGatewayInfo struct {
	NID         int       `json:"nid,omitempty" gorm:"column:nid"`
	GwName      string    `json:"gw_name" gorm:"column:gw_name"`           // 网关名称
	GwUrl       string    `json:"gw_uri" gorm:"column:gw_uri"`             //网关调用地址
	GwType      GwType    `json:"gw_type" gorm:"column:gw_type"`           // 网关类型
	Priority    int       `json:"priority" gorm:"column:priority"`         // 优先级
	BlackPrefix string    `json:"black_prefix" gorm:"column:black_prefix"` // 接口其他参数
	IStatus     int       `json:"i_status" gorm:"i_status"`                //状态，0停用，1启动，-1删除
	MbLevel     int       `json:"mb_level" gorm:"mb_level"`                // 黑名单等级
	JoinDt      time.Time `json:"join_dt" gorm:"join_dt"`                  // 添加时间
	Remark      string    `json:"remark" gorm:"remark"`                    // 备注
	GwEnID      string    `json:"gw_enid" gorm:"gw_enid"`                  // 第三方接口的企业id
	GwPass      string    `json:"gw_pass" gorm:"gw_pass"`                  // 第三方接口的企业密码
	GwAk        string    `json:"gw_ak" gorm:"gw_ak"`
}

func (SysGatewayInfo) TableName() string {
	return "sys_gateway"
}

var sysGatewayImpl *SysGatewayImpl

type SysGatewayImpl struct {
	DB *gorm.DB
}

func InitSysGatewayImplRepo(d *gorm.DB) {
	sysGatewayImpl = &SysGatewayImpl{
		DB: d,
	}
}

func GetSysGatewayImpl() SysGatewayRepo {
	return sysGatewayImpl
}

type SysGatewayRepo interface {
	GetByEnID(ctx context.Context, id int) (*SysGatewayInfo, error)
	GetAllActiveGateWay(ctx context.Context) ([]*SysGatewayInfo, error)
}

func (s *SysGatewayImpl) GetByEnID(ctx context.Context, id int) (*SysGatewayInfo, error) {
	var res *SysGatewayInfo
	err := s.DB.WithContext(ctx).Where("nID = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *SysGatewayImpl) GetAllActiveGateWay(ctx context.Context) ([]*SysGatewayInfo, error) {
	var res []*SysGatewayInfo
	err := s.DB.WithContext(ctx).Where("i_status = ?", IStatusActive).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
