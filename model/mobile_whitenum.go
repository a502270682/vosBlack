package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// MobileWhitenum 每个企业自己的白名单号码
type MobileWhitenum struct {
	NID      int       `json:"nid,omitempty" gorm:"column:nID"`
	EnID     int       `json:"en_id,omitempty" gorm:"column:en_id"` // 请求企业ID
	WhiteNum string    `json:"white_num" gorm:"column:whitenum"`    // 白名单号码
	WnName   string    `json:"wn_name" gorm:"column:wn_name"`       // 白名单名称
	IStatus  int       `json:"i_status" gorm:"column:i_status"`     // 状态，1启用，0停用，-1删除
	JoinDt   time.Time `json:"join_dt" gorm:"column:join_dt"`
	JoinUser string    `json:"join_user" gorm:"column:join_user"`
	Remark   string    `json:"remark" gorm:"column:remark"`
}

func (MobileWhitenum) TableName() string {
	return "mobile_whitenum"
}

type MobileWhitenumImpl struct {
	DB *gorm.DB
}

var mobileWhitenumImpl *MobileWhitenumImpl

type MobileWhitenumRepo interface {
	GetByWhiteNum(ctx context.Context, enID int, whiteNum string) (*MobileWhitenum, error)
}

func InitMobileWhitenumRepo(d *gorm.DB) {
	mobileWhitenumImpl = &MobileWhitenumImpl{
		DB: d,
	}
}

func GetMobileWhitenumImpl() MobileWhitenumRepo {
	return mobileWhitenumImpl
}

func (m *MobileWhitenumImpl) GetByWhiteNum(ctx context.Context, enID int, whiteNum string) (*MobileWhitenum, error) {
	res := &MobileWhitenum{}
	err := m.DB.WithContext(ctx).Where("en_id = ?", enID).
		Where("whitenum = ?", whiteNum).First(res).Error
	return res, err
}
