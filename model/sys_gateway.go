package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type SysGatewayInfo struct {
	NID         int
	GwName      string
	GwUrl       string
	GwType      GwType
	Priority    int
	BlackPrefix string
	IStatus     IStatus
	MbLevel     int
	JoinDt      time.Time
	Remark      string
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
