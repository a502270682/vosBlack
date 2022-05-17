package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type EnterpriseInfo struct {
	NID           int       `json:"nid" gorm:"column:nid"`
	EnName        string    `json:"en_name" gorm:"en_name"`
	Licensenum    string    `json:"licensenum" gorm:"licensenum"`
	EnKind        string    `json:"en_kind" gorm:"en_kind"`
	Regpeople     string    `json:"regpeople" gorm:"regpeople"`
	RegpeopleSfz  string    `json:"regpeople_sfz" gorm:"regpeople_sfz"`
	Contacts      string    `json:"contacts" gorm:"contacts"`
	ContactsSfz   string    `json:"contacts_sfz" gorm:"contacts_sfz"`
	ContactsTel   string    `json:"contacts_tel" gorm:"contacts_tel"`
	ContactsEmail string    `json:"contacts_email" gorm:"contacts_email"`
	IsPattern     int       `json:"is_pattern" gorm:"is_pattern"`
	IsWhitenum    int       `json:"is_whitenum" gorm:"is_whitenum"`
	IsBlacknum    int       `json:"is_blacknum" gorm:"is_blacknum"`
	IsGateway     int       `json:"is_gateway" gorm:"is_gateway"`
	IStatus       IStatus   `json:"i_status" gorm:"i_status"`
	JoinDt        time.Time `json:"join_dt" gorm:"join_dt"`
	JoinUser      string    `json:"join_user" gorm:"join_user"`
	Remark        string    `json:"remark" gorm:"remark"`
}

func (e EnterpriseInfo) TableName() string {
	return "t_enterprise_infolist"
}

var enterpriseInfoImpl *EnterpriseInfoImpl

type EnterpriseInfoImpl struct {
	DB *gorm.DB
}

func InitEnterpriseInfoImplRepo(d *gorm.DB) {
	enterpriseInfoImpl = &EnterpriseInfoImpl{
		DB: d,
	}
}

func GetEnterpriseInfoImpl() EnterpriseInfoRepo {
	return enterpriseInfoImpl
}

type EnterpriseInfoRepo interface {
	GetByEnID(ctx context.Context, id int) (*EnterpriseInfo, error)
	GetAllActiveEnterprise(ctx context.Context) ([]*EnterpriseInfo, error)
}

func (e *EnterpriseInfoImpl) GetAllActiveEnterprise(ctx context.Context) ([]*EnterpriseInfo, error) {
	var res []*EnterpriseInfo
	err := e.DB.WithContext(ctx).Where("i_status = ?", IStatusActive).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *EnterpriseInfoImpl) GetByEnID(ctx context.Context, id int) (*EnterpriseInfo, error) {
	var res *EnterpriseInfo
	err := e.DB.WithContext(ctx).Where("nID = ?", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
